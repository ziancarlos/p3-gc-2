package repository

import (
	"context"
	"fmt"
	"p3-graded-challenge-2-ziancarlos/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepository interface {
	Create(ctx context.Context, product *models.Product) error
	FindAll(ctx context.Context) ([]models.Product, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.Product, error)
	Update(ctx context.Context, id primitive.ObjectID, product *models.Product) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type productRepository struct {
	collection *mongo.Collection
}

func NewProductRepository(collection *mongo.Collection) ProductRepository {
	return &productRepository{
		collection: collection,
	}
}

func (r *productRepository) Create(ctx context.Context, product *models.Product) error {
	result, err := r.collection.InsertOne(ctx, product)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	product.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *productRepository) FindAll(ctx context.Context) ([]models.Product, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find products: %w", err)
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, fmt.Errorf("failed to decode products: %w", err)
	}

	return products, nil
}

func (r *productRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Product, error) {
	var product models.Product
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to find product: %w", err)
	}
	return &product, nil
}

func (r *productRepository) Update(ctx context.Context, id primitive.ObjectID, product *models.Product) error {
	update := bson.M{
		"$set": bson.M{
			"name":  product.Name,
			"price": product.Price,
		},
	}
	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("product not found")
	}
	return nil
}

func (r *productRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("product not found")
	}
	return nil
}

