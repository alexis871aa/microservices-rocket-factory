package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	customMiddleware "github.com/alexis871aa/microservices-rocket-factory/order/internal/middleware"
	orderV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/payment/v1"
)

const (
	httpPort      = "8080"
	inventoryAddr = "localhost:50051"
	paymentAddr   = "localhost:50052"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
	// Таймауты для внешних сервисов
	inventoryTimeout = 5 * time.Second
	paymentTimeout   = 3 * time.Second
)

// OrderStorage представляет потокобезопасное хранилище данных о заказах
type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*orderV1.OrderDto
}

// NewOrderStorage создаёт новое хранилище данных для заказов
func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*orderV1.OrderDto),
	}
}

func (s *OrderStorage) CreateOrder(order *orderV1.OrderDto) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[order.OrderUUID] = order
	return nil
}

func (s *OrderStorage) GetOrder(orderUUID string) (*orderV1.OrderDto, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[orderUUID]
	if !ok {
		return nil, fmt.Errorf("заказ не найден")
	}
	return order, nil
}

func (s *OrderStorage) UpdateOrder(orderUUID string, order *orderV1.OrderDto) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.orders[orderUUID]; !ok {
		return fmt.Errorf("заказ не найден")
	}
	s.orders[orderUUID] = order
	return nil
}

func (s *OrderStorage) DeleteOrder(orderUUID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.orders[orderUUID]; !ok {
		return fmt.Errorf("заказ не найден")
	}
	delete(s.orders, orderUUID)
	return nil
}

// OrderHandler реализует интерфейс orderV1.Handler для обработки запросов к API OrderService
type OrderHandler struct {
	storage         *OrderStorage
	inventoryClient inventoryV1.InventoryServiceClient
	paymentClient   paymentV1.PaymentServiceClient
}

// NewOrderHandler создаёт новый обработчик запросов к API OrderService
func NewOrderHandler(storage *OrderStorage, inventoryClient inventoryV1.InventoryServiceClient, paymentClient paymentV1.PaymentServiceClient) *OrderHandler {
	return &OrderHandler{
		storage:         storage,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, inventoryTimeout)
	defer cancel()

	resp, err := h.inventoryClient.ListParts(ctxWithTimeout, &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			Uuids: req.PartUuids,
		},
	})
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Ошибка при получении информации о деталях",
		}, nil
	}

	if len(resp.Parts) != len(req.PartUuids) {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "Не все необходимые детали найдены",
		}, nil
	}

	var totalPrice float32
	for _, part := range resp.Parts {
		totalPrice += float32(part.Price)
	}

	orderUuid := uuid.New().String()
	order := &orderV1.OrderDto{
		OrderUUID:  orderUuid,
		UserUUID:   req.UserUUID,
		PartUuids:  req.PartUuids,
		TotalPrice: totalPrice,
		Status:     orderV1.OrderStatusPENDINGPAYMENT,
	}

	err = h.storage.CreateOrder(order)
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Ошибка при создании заказа",
		}, nil
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID:  orderUuid,
		TotalPrice: totalPrice,
	}, nil
}

func convertPaymentMethod(orderMethod orderV1.PaymentMethod) paymentV1.PaymentMethod {
	switch orderMethod {
	case orderV1.PaymentMethodUNKNOWN:
		return paymentV1.PaymentMethod_UNKNOWN
	case orderV1.PaymentMethodCARD:
		return paymentV1.PaymentMethod_CARD
	case orderV1.PaymentMethodSBP:
		return paymentV1.PaymentMethod_SBP
	case orderV1.PaymentMethodCREDITCARD:
		return paymentV1.PaymentMethod_CREDIT_CARD
	case orderV1.PaymentMethodINVESTORMONEY:
		return paymentV1.PaymentMethod_INVESTOR_MONEY
	default:
		return paymentV1.PaymentMethod_UNKNOWN
	}
}

func (h *OrderHandler) PaymentOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PaymentOrderParams) (orderV1.PaymentOrderRes, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, paymentTimeout)
	defer cancel()

	order, err := h.storage.GetOrder(params.OrderUUID)
	if err != nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Заказ не найден",
		}, nil
	}

	resp, err := h.paymentClient.PayOrder(ctxWithTimeout, &paymentV1.PayOrderRequest{
		OrderUuid:     order.OrderUUID,
		UserUuid:      order.UserUUID,
		PaymentMethod: convertPaymentMethod(req.PaymentMethod),
	})
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Ошибка при обработке платежа",
		}, nil
	}

	order.Status = orderV1.OrderStatusPAID
	order.TransactionUUID = orderV1.NewOptNilString(resp.TransactionUuid)
	order.PaymentMethod = &orderV1.NilOrderDtoPaymentMethod{
		Value: orderV1.OrderDtoPaymentMethod(req.PaymentMethod),
	}

	cerr := h.storage.UpdateOrder(params.OrderUUID, order)
	if cerr != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Ошибка при обработке платежа",
		}, nil
	}

	return &orderV1.PayOrderResponse{
		TransactionUUID: resp.TransactionUuid,
	}, nil
}

func (h *OrderHandler) GetOrderById(_ context.Context, params orderV1.GetOrderByIdParams) (orderV1.GetOrderByIdRes, error) {
	order, err := h.storage.GetOrder(params.OrderUUID)
	if err != nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Заказ не найден",
		}, nil
	}

	return &orderV1.GetOrderResponse{
		OrderDto: *order,
	}, nil
}

func (h *OrderHandler) CancelOrder(_ context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	order, err := h.storage.GetOrder(params.OrderUUID)
	if err != nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "заказ не найден",
		}, nil
	}

	if order.Status == orderV1.OrderStatusPAID {
		return &orderV1.ConflictError{
			Code:    409,
			Message: "заказ уже оплачен и не может быть отменён",
		}, nil
	}

	if order.Status == orderV1.OrderStatusPENDINGPAYMENT {
		order.Status = orderV1.OrderStatusCANCELLED
		err = h.storage.UpdateOrder(params.OrderUUID, order)
		if err != nil {
			return &orderV1.InternalServerError{
				Code:    500,
				Message: "ошибка при отмене заказа",
			}, nil
		}

		return &orderV1.CancelOrderNoContent{}, nil
	}

	return &orderV1.ConflictError{
		Code:    409,
		Message: "заказ не может быть отменён",
	}, nil
}

func (h *OrderHandler) NewError(_ context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    orderV1.NewOptInt(http.StatusInternalServerError),
			Message: orderV1.NewOptString(err.Error()),
		},
	}
}

func main() {
	inventoryConn, err := grpc.NewClient(
		inventoryAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return
	}
	defer func() {
		if ierr := inventoryConn.Close(); ierr != nil {
			log.Printf("failed to close inventory connection: %v\n", ierr)
		}
	}()
	inventoryClient := inventoryV1.NewInventoryServiceClient(inventoryConn)

	paymentConn, err := grpc.NewClient(
		paymentAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return
	}
	defer func() {
		if perr := paymentConn.Close(); perr != nil {
			log.Printf("failed to close payment connection: %v\n", perr)
		}
	}()
	paymentClient := paymentV1.NewPaymentServiceClient(paymentConn)

	storage := NewOrderStorage()
	orderHandler := NewOrderHandler(storage, inventoryClient, paymentClient)

	orderServer, err := orderV1.NewServer(orderHandler)
	if err != nil {
		log.Printf("ошибка создания сервера OpenAPI: %v", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(2 * time.Second))
	r.Use(customMiddleware.RequestLogger)

	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска http сервера: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы http сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке http сервера: %v\n", err)
	}

	log.Println("✅ Http сервер остановлен")
}
