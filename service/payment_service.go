package service

import (
	"context"
	"fmt"
	"p3-graded-challenge-2-ziancarlos/models"
	"p3-graded-challenge-2-ziancarlos/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentService interface {
	CreatePayment(ctx context.Context, req *models.PaymentRequest) (*models.PaymentResponse, error)
	GetAllPayments(ctx context.Context) ([]models.PaymentResponse, error)
	GetPaymentByID(ctx context.Context, id string) (*models.PaymentResponse, error)
	DeletePayment(ctx context.Context, id string) error
}

type paymentService struct {
	repo repository.PaymentRepository
}

func NewPaymentService(repo repository.PaymentRepository) PaymentService {
	return &paymentService{
		repo: repo,
	}
}

func (s *paymentService) CreatePayment(ctx context.Context, req *models.PaymentRequest) (*models.PaymentResponse, error) {
	if req.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than 0")
	}

	payment := &models.Payment{
		Amount: req.Amount,
	}

	err := s.repo.Create(ctx, payment)
	if err != nil {
		return nil, err
	}

	return &models.PaymentResponse{
		ID:     payment.ID.Hex(),
		Amount: payment.Amount,
	}, nil
}

func (s *paymentService) GetAllPayments(ctx context.Context) ([]models.PaymentResponse, error) {
	payments, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var responses []models.PaymentResponse
	for _, payment := range payments {
		responses = append(responses, models.PaymentResponse{
			ID:     payment.ID.Hex(),
			Amount: payment.Amount,
		})
	}

	return responses, nil
}

func (s *paymentService) GetPaymentByID(ctx context.Context, id string) (*models.PaymentResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid payment ID: %w", err)
	}

	payment, err := s.repo.FindByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	return &models.PaymentResponse{
		ID:     payment.ID.Hex(),
		Amount: payment.Amount,
	}, nil
}

func (s *paymentService) DeletePayment(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid payment ID: %w", err)
	}

	return s.repo.Delete(ctx, objectID)
}

