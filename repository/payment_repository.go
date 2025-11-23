package repository

import (
	"context"
	"fmt"
	"p3-graded-challenge-2-ziancarlos/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment *models.Payment) error
	FindAll(ctx context.Context) ([]models.Payment, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.Payment, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type paymentRepository struct {
	collection *mongo.Collection
}

func NewPaymentRepository(collection *mongo.Collection) PaymentRepository {
	return &paymentRepository{
		collection: collection,
	}
}

func (r *paymentRepository) Create(ctx context.Context, payment *models.Payment) error {
	result, err := r.collection.InsertOne(ctx, payment)
	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}
	payment.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *paymentRepository) FindAll(ctx context.Context) ([]models.Payment, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find payments: %w", err)
	}
	defer cursor.Close(ctx)

	var payments []models.Payment
	if err := cursor.All(ctx, &payments); err != nil {
		return nil, fmt.Errorf("failed to decode payments: %w", err)
	}

	return payments, nil
}

func (r *paymentRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Payment, error) {
	var payment models.Payment
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&payment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, fmt.Errorf("failed to find payment: %w", err)
	}
	return &payment, nil
}

func (r *paymentRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete payment: %w", err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("payment not found")
	}
	return nil
}
