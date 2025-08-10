package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	paymentV1API "github.com/alexis871aa/microservices-rocket-factory/payment/internal/api/payment/v1"
	"github.com/alexis871aa/microservices-rocket-factory/payment/internal/config"
	paymentService "github.com/alexis871aa/microservices-rocket-factory/payment/internal/service/payment"
	paymentV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/payment/v1"
)

const configPath = "./deploy/compose/payment/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	lis, err := net.Listen("tcp", config.AppConfig().PaymentGRPC.Address())
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

	service := paymentService.NewService()
	api := paymentV1API.NewAPI(service)

	paymentV1.RegisterPaymentServiceServer(s, api)

	reflection.Register(s)

	go func() {
		log.Printf("starting PaymentService server on port %s", config.AppConfig().PaymentGRPC.Port())
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
