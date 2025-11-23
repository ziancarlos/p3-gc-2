package scheduler

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// CleanupScheduler handles scheduled cleanup tasks
type CleanupScheduler struct {
	paymentCollection *mongo.Collection
	productCollection *mongo.Collection
	interval          time.Duration
}

// NewCleanupScheduler creates a new cleanup scheduler
func NewCleanupScheduler(paymentCollection, productCollection *mongo.Collection, interval time.Duration) *CleanupScheduler {
	return &CleanupScheduler{
		paymentCollection: paymentCollection,
		productCollection: productCollection,
		interval:          interval,
	}
}

// Start begins the scheduled cleanup tasks
func (s *CleanupScheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	log.Printf("Cleanup scheduler started with interval: %v", s.interval)

	for {
		select {
		case <-ctx.Done():
			log.Println("Cleanup scheduler stopped")
			return
		case <-ticker.C:
			s.runCleanup(ctx)
		}
	}
}

// runCleanup performs the cleanup operations
func (s *CleanupScheduler) runCleanup(ctx context.Context) {
	log.Println("Running scheduled cleanup...")

	// Example: Delete payments older than 30 days (if there's a timestamp field)
	// For now, we'll just log the count of documents
	paymentCount, err := s.paymentCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Printf("Error counting payments: %v", err)
	} else {
		log.Printf("Current payment count: %d", paymentCount)
	}

	productCount, err := s.productCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Printf("Error counting products: %v", err)
	} else {
		log.Printf("Current product count: %d", productCount)
	}

	log.Println("Cleanup completed")
}

// RunImmediately executes cleanup task immediately
func (s *CleanupScheduler) RunImmediately(ctx context.Context) {
	log.Println("Running immediate cleanup...")
	s.runCleanup(ctx)
}
