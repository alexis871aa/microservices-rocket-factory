package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	paymentV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/payment/v1"
)

const grpcPort = 50052

type paymentService struct {
	paymentV1.UnimplementedPaymentServiceServer

	mu sync.RWMutex
}

func (s *paymentService) PayOrder(_ context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	transactionUuid := uuid.NewString()

	log.Printf("Оплата прошла успешно, transaction_uuid: %s", transactionUuid)

	return &paymentV1.PayOrderResponse{
		TransactionUuid: transactionUuid,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v", cerr)
		}
	}()

	s := grpc.NewServer()

	service := &paymentService{}

	paymentV1.RegisterPaymentServiceServer(s, service)

	reflection.Register(s)

	go func() {
		log.Printf("starting PaymentService server on port %d", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Fatalf("failed to serve: %v", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down PaymentService server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
