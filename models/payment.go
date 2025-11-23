package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payment struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Amount float64            `json:"amount" bson:"amount" validate:"required,gt=0"`
}

type PaymentRequest struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
}

type PaymentResponse struct {
	ID     string  `json:"id"`
	Amount float64 `json:"amount"`
}
