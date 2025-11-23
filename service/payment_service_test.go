package service

import (
	"context"
	"errors"
	"p3-graded-challenge-2-ziancarlos/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockPaymentRepository is a mock implementation of PaymentRepository
type MockPaymentRepository struct {
	mock.Mock
}

func (m *MockPaymentRepository) Create(ctx context.Context, payment *models.Payment) error {
	args := m.Called(ctx, payment)
	if args.Get(0) == nil {
		// Simulate setting the ID after creation
		payment.ID = primitive.NewObjectID()
		return nil
	}
	return args.Error(0)
}

func (m *MockPaymentRepository) FindAll(ctx context.Context) ([]models.Payment, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Payment), args.Error(1)
}

func (m *MockPaymentRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Payment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Payment), args.Error(1)
}

func (m *MockPaymentRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreatePayment_Success(t *testing.T) {
	mockRepo := new(MockPaymentRepository)
	service := NewPaymentService(mockRepo)

	ctx := context.Background()
	req := &models.PaymentRequest{
		Amount: 100.50,
	}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*models.Payment")).Return(nil)

	result, err := service.CreatePayment(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 100.50, result.Amount)
	assert.NotEmpty(t, result.ID)
	mockRepo.AssertExpectations(t)
}

func TestCreatePayment_InvalidAmount(t *testing.T) {
	mockRepo := new(MockPaymentRepository)
	service := NewPaymentService(mockRepo)

	ctx := context.Background()
	req := &models.PaymentRequest{
		Amount: -10.0,
	}

	result, err := service.CreatePayment(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "amount must be greater than 0")
}

func TestGetAllPayments_Success(t *testing.T) {
	mockRepo := new(MockPaymentRepository)
	service := NewPaymentService(mockRepo)

	ctx := context.Background()
	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()

	expectedPayments := []models.Payment{
		{ID: id1, Amount: 100.0},
		{ID: id2, Amount: 200.0},
	}

	mockRepo.On("FindAll", ctx).Return(expectedPayments, nil)

	result, err := service.GetAllPayments(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, id1.Hex(), result[0].ID)
	assert.Equal(t, 100.0, result[0].Amount)
	mockRepo.AssertExpectations(t)
}

func TestGetPaymentByID_Success(t *testing.T) {
	mockRepo := new(MockPaymentRepository)
	service := NewPaymentService(mockRepo)

	ctx := context.Background()
	id := primitive.NewObjectID()
	expectedPayment := &models.Payment{
		ID:     id,
		Amount: 150.0,
	}

	mockRepo.On("FindByID", ctx, id).Return(expectedPayment, nil)

	result, err := service.GetPaymentByID(ctx, id.Hex())

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, id.Hex(), result.ID)
	assert.Equal(t, 150.0, result.Amount)
	mockRepo.AssertExpectations(t)
}

func TestGetPaymentByID_InvalidID(t *testing.T) {
	mockRepo := new(MockPaymentRepository)
	service := NewPaymentService(mockRepo)

	ctx := context.Background()

	result, err := service.GetPaymentByID(ctx, "invalid-id")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid payment ID")
}

func TestDeletePayment_Success(t *testing.T) {
	mockRepo := new(MockPaymentRepository)
	service := NewPaymentService(mockRepo)

	ctx := context.Background()
	id := primitive.NewObjectID()

	mockRepo.On("Delete", ctx, id).Return(nil)

	err := service.DeletePayment(ctx, id.Hex())

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeletePayment_NotFound(t *testing.T) {
	mockRepo := new(MockPaymentRepository)
	service := NewPaymentService(mockRepo)

	ctx := context.Background()
	id := primitive.NewObjectID()

	mockRepo.On("Delete", ctx, id).Return(errors.New("payment not found"))

	err := service.DeletePayment(ctx, id.Hex())

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

