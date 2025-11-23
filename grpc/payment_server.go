package grpc

import (
	"context"
	"p3-graded-challenge-2-ziancarlos/models"
	pb "p3-graded-challenge-2-ziancarlos/proto/payment"
	"p3-graded-challenge-2-ziancarlos/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaymentServer struct {
	pb.UnimplementedPaymentServiceServer
	service service.PaymentService
}

func NewPaymentServer(service service.PaymentService) *PaymentServer {
	return &PaymentServer{
		service: service,
	}
}

func (s *PaymentServer) CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.PaymentResponse, error) {
	paymentReq := &models.PaymentRequest{
		Amount: req.Amount,
	}

	payment, err := s.service.CreatePayment(ctx, paymentReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create payment: %v", err)
	}

	return &pb.PaymentResponse{
		Id:     payment.ID,
		Amount: payment.Amount,
	}, nil
}

func (s *PaymentServer) GetAllPayments(ctx context.Context, req *pb.GetAllPaymentsRequest) (*pb.GetAllPaymentsResponse, error) {
	payments, err := s.service.GetAllPayments(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get payments: %v", err)
	}

	var pbPayments []*pb.PaymentResponse
	for _, payment := range payments {
		pbPayments = append(pbPayments, &pb.PaymentResponse{
			Id:     payment.ID,
			Amount: payment.Amount,
		})
	}

	return &pb.GetAllPaymentsResponse{
		Payments: pbPayments,
	}, nil
}

func (s *PaymentServer) GetPaymentByID(ctx context.Context, req *pb.GetPaymentByIDRequest) (*pb.PaymentResponse, error) {
	payment, err := s.service.GetPaymentByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "payment not found: %v", err)
	}

	return &pb.PaymentResponse{
		Id:     payment.ID,
		Amount: payment.Amount,
	}, nil
}

func (s *PaymentServer) DeletePayment(ctx context.Context, req *pb.DeletePaymentRequest) (*pb.DeletePaymentResponse, error) {
	err := s.service.DeletePayment(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete payment: %v", err)
	}

	return &pb.DeletePaymentResponse{
		Message: "Payment deleted successfully",
	}, nil
}
