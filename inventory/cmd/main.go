package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	partRepository "github.com/alexis871aa/microservices-rocket-factory/inventory/internal/repository/part"
	inventoryV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50051

type inventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer

	mu    sync.RWMutex
	parts map[string]*inventoryV1.Part
}

func (s *inventoryService) GetPart(_ context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	part, ok := s.parts[req.GetUuid()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", req.GetUuid())
	}

	return &inventoryV1.GetPartResponse{
		Part: part,
	}, nil
}

func (s *inventoryService) ListParts(_ context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var parts []*inventoryV1.Part
	filter := req.GetFilter()

	for _, part := range s.parts {
		if matchesFilter(part, filter) {
			parts = append(parts, &inventoryV1.Part{
				Uuid:          part.Uuid,
				Name:          part.Name,
				Description:   part.Description,
				Price:         part.Price,
				StockQuantity: part.StockQuantity,
				Category:      part.Category,
				Dimensions:    copyDimensions(part.Dimensions),
				Manufacturer:  copyManufacturer(part.Manufacturer),
				Tags:          append([]string(nil), part.Tags...),
				Metadata:      copyMetadata(part.Metadata),
				CreatedAt:     part.CreatedAt,
			})
		}
	}

	return &inventoryV1.ListPartsResponse{
		Parts: parts,
	}, nil
}

func copyDimensions(d *inventoryV1.Dimensions) *inventoryV1.Dimensions {
	if d == nil {
		return nil
	}
	return &inventoryV1.Dimensions{
		Length: d.Length,
		Width:  d.Width,
		Height: d.Height,
		Weight: d.Weight,
	}
}

func copyManufacturer(m *inventoryV1.Manufacturer) *inventoryV1.Manufacturer {
	if m == nil {
		return nil
	}
	return &inventoryV1.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

func copyMetadata(original map[string]*inventoryV1.Value) map[string]*inventoryV1.Value {
	if original == nil {
		return nil
	}

	result := make(map[string]*inventoryV1.Value)

	for key, value := range original {
		if value == nil {
			result[key] = nil
			continue
		}

		result[key] = copyValue(value)
	}

	return result
}

func copyValue(original *inventoryV1.Value) *inventoryV1.Value {
	if original == nil {
		return nil
	}

	switch kind := original.Kind.(type) {
	case *inventoryV1.Value_StringValue:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_StringValue{
				StringValue: kind.StringValue,
			},
		}
	case *inventoryV1.Value_Int64Value:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_Int64Value{
				Int64Value: kind.Int64Value,
			},
		}
	case *inventoryV1.Value_DoubleValue:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_DoubleValue{
				DoubleValue: kind.DoubleValue,
			},
		}
	case *inventoryV1.Value_BoolValue:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_BoolValue{
				BoolValue: kind.BoolValue,
			},
		}
	default:
		return &inventoryV1.Value{}
	}
}

func matchesFilter(part *inventoryV1.Part, filter *inventoryV1.PartsFilter) bool {
	return matchesBy(part, filter.Uuids, func(part *inventoryV1.Part) string { return part.Uuid }) &&
		matchesBy(part, filter.Names, func(part *inventoryV1.Part) string { return part.Name }) &&
		matchesBy(part, filter.Categories, func(part *inventoryV1.Part) inventoryV1.Category { return part.Category }) &&
		matchesBy(part, filter.ManufacturerCountries, func(part *inventoryV1.Part) string { return part.Manufacturer.Country })
}

func matchesBy[T comparable](part *inventoryV1.Part, filter []T, get func(part *inventoryV1.Part) T) bool {
	if len(filter) == 0 {
		return true
	}

	for _, filterValue := range filter {
		if get(part) == filterValue {
			return true
		}
	}
	return false
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

	// Регистрируем наш сервис
	repo := partRepository.NewRepository()

	service := &inventoryService{
		parts: make(map[string]*inventoryV1.Part),
	}
	service.initParts()

	inventoryV1.RegisterInventoryServiceServer(s, service)

	reflection.Register(s)

	go func() {
		log.Printf("starting InventoryService server on port %d", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Fatalf("failed to serve: %v", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down InventoryService server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}

func (s *inventoryService) initParts() {
	parts := generateParts()

	for _, part := range parts {
		s.parts[part.Uuid] = part
	}
}

func generateParts() []*inventoryV1.Part {
	names := []string{
		"Main Engine",
		"Reserve Engine",
		"Thruster",
		"Fuel Tank",
		"Left Wing",
		"Right Wing",
		"Window A",
		"Window B",
		"Control Module",
		"Stabilizer",
	}

	descriptions := []string{
		"Primary propulsion unit",
		"Backup propulsion unit",
		"Thruster for fine adjustments",
		"Main fuel tank",
		"Left aerodynamic wing",
		"Right aerodynamic wing",
		"Front viewing window",
		"Side viewing window",
		"Flight control module",
		"Stabilization fin",
	}

	var parts []*inventoryV1.Part
	for i := 0; i < gofakeit.Number(1, 50); i++ {
		idx := gofakeit.Number(0, len(names)-1)
		parts = append(parts, &inventoryV1.Part{
			Uuid:          uuid.NewString(),
			Name:          names[idx],
			Description:   descriptions[idx],
			Price:         roundTo(gofakeit.Float64Range(100, 10_000)),
			StockQuantity: int64(gofakeit.Number(1, 100)),
			Category:      inventoryV1.Category(gofakeit.Number(1, 4)), //nolint:gosec // safe: gofakeit.Number returns 1..4
			Dimensions:    generateDimensions(),
			Manufacturer:  generateManufacturer(),
			Tags:          generateTags(),
			Metadata:      generateMetadata(),
			CreatedAt:     timestamppb.Now(),
		})
	}

	return parts
}

func generateDimensions() *inventoryV1.Dimensions {
	return &inventoryV1.Dimensions{
		Length: roundTo(gofakeit.Float64Range(1, 1000)),
		Width:  roundTo(gofakeit.Float64Range(1, 1000)),
		Height: roundTo(gofakeit.Float64Range(1, 1000)),
		Weight: roundTo(gofakeit.Float64Range(1, 1000)),
	}
}

func generateManufacturer() *inventoryV1.Manufacturer {
	return &inventoryV1.Manufacturer{
		Name:    gofakeit.Name(),
		Country: gofakeit.Country(),
		Website: gofakeit.URL(),
	}
}

func generateTags() []string {
	var tags []string
	for i := 0; i < gofakeit.Number(1, 10); i++ {
		tags = append(tags, gofakeit.EmojiTag())
	}

	return tags
}

func generateMetadata() map[string]*inventoryV1.Value {
	metadata := make(map[string]*inventoryV1.Value)

	for i := 0; i < gofakeit.Number(1, 10); i++ {
		metadata[gofakeit.Word()] = generateMetadataValue()
	}

	return metadata
}

func generateMetadataValue() *inventoryV1.Value {
	switch gofakeit.Number(0, 3) {
	case 0:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_StringValue{
				StringValue: gofakeit.Word(),
			},
		}

	case 1:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_Int64Value{
				Int64Value: int64(gofakeit.Number(1, 100)),
			},
		}

	case 2:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_DoubleValue{
				DoubleValue: roundTo(gofakeit.Float64Range(1, 100)),
			},
		}

	case 3:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_BoolValue{
				BoolValue: gofakeit.Bool(),
			},
		}

	default:
		return nil
	}
}

func roundTo(x float64) float64 {
	return math.Round(x*100) / 100
}
