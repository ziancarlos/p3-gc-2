package service

import (
	"context"
	"fmt"
	"p3-graded-challenge-2-ziancarlos/models"
	"p3-graded-challenge-2-ziancarlos/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductService interface {
	CreateProduct(ctx context.Context, req *models.ProductRequest) (*models.ProductResponse, error)
	GetAllProducts(ctx context.Context) ([]models.ProductResponse, error)
	GetProductByID(ctx context.Context, id string) (*models.ProductResponse, error)
	UpdateProduct(ctx context.Context, id string, req *models.ProductRequest) (*models.ProductResponse, error)
	DeleteProduct(ctx context.Context, id string) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{
		repo: repo,
	}
}

func (s *productService) CreateProduct(ctx context.Context, req *models.ProductRequest) (*models.ProductResponse, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if req.Price <= 0 {
		return nil, fmt.Errorf("price must be greater than 0")
	}

	product := &models.Product{
		Name:  req.Name,
		Price: req.Price,
	}

	err := s.repo.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	return &models.ProductResponse{
		ID:    product.ID.Hex(),
		Name:  product.Name,
		Price: product.Price,
	}, nil
}

func (s *productService) GetAllProducts(ctx context.Context) ([]models.ProductResponse, error) {
	products, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var responses []models.ProductResponse
	for _, product := range products {
		responses = append(responses, models.ProductResponse{
			ID:    product.ID.Hex(),
			Name:  product.Name,
			Price: product.Price,
		})
	}

	return responses, nil
}

func (s *productService) GetProductByID(ctx context.Context, id string) (*models.ProductResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %w", err)
	}

	product, err := s.repo.FindByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	return &models.ProductResponse{
		ID:    product.ID.Hex(),
		Name:  product.Name,
		Price: product.Price,
	}, nil
}

func (s *productService) UpdateProduct(ctx context.Context, id string, req *models.ProductRequest) (*models.ProductResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %w", err)
	}

	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if req.Price <= 0 {
		return nil, fmt.Errorf("price must be greater than 0")
	}

	product := &models.Product{
		Name:  req.Name,
		Price: req.Price,
	}

	err = s.repo.Update(ctx, objectID, product)
	if err != nil {
		return nil, err
	}

	return &models.ProductResponse{
		ID:    id,
		Name:  product.Name,
		Price: product.Price,
	}, nil
}

func (s *productService) DeleteProduct(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid product ID: %w", err)
	}

	return s.repo.Delete(ctx, objectID)
}

