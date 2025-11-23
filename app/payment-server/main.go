package main

import (
	"fmt"
	"log"
	"net"
	"p3-graded-challenge-2-ziancarlos/config"
	grpcServer "p3-graded-challenge-2-ziancarlos/grpc"
	"p3-graded-challenge-2-ziancarlos/middleware"
	pb "p3-graded-challenge-2-ziancarlos/proto/payment"
	"p3-graded-challenge-2-ziancarlos/repository"
	"p3-graded-challenge-2-ziancarlos/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize JWT
	middleware.InitJWT(cfg.JWTSecret)

	// Connect to MongoDB
	client, err := config.ConnectDB(cfg.MongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Setup repositories
	paymentCollection := config.GetCollection(client, cfg.PaymentDBName, "payments")
	paymentRepo := repository.NewPaymentRepository(paymentCollection)

	// Setup services
	paymentService := service.NewPaymentService(paymentRepo)

	// Create gRPC server with JWT interceptor
	grpcServerInstance := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.UnaryInterceptor),
	)

	// Register payment service
	paymentServer := grpcServer.NewPaymentServer(paymentService)
	pb.RegisterPaymentServiceServer(grpcServerInstance, paymentServer)

	// Enable reflection for gRPC tools like grpcurl
	reflection.Register(grpcServerInstance)

	// Start listening
	address := fmt.Sprintf(":%s", cfg.PortPayment)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Payment gRPC server listening on %s", address)
	if err := grpcServerInstance.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

