package controllers

import (
	"net/http"
	"p3-graded-challenge-2-ziancarlos/models"
	"p3-graded-challenge-2-ziancarlos/service"

	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	service service.PaymentService
}

func NewPaymentController(service service.PaymentService) *PaymentController {
	return &PaymentController{
		service: service,
	}
}

// CreatePayment godoc
// @Summary Create a new payment
// @Description Create a new payment with the provided amount
// @Tags payments
// @Accept json
// @Produce json
// @Param payment body models.PaymentRequest true "Payment Request"
// @Success 201 {object} models.PaymentResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /payments [post]
func (c *PaymentController) CreatePayment(ctx *gin.Context) {
	var req models.PaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := c.service.CreatePayment(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, payment)
}

// GetAllPayments godoc
// @Summary Get all payments
// @Description Get a list of all payments
// @Tags payments
// @Produce json
// @Success 200 {array} models.PaymentResponse
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /payments [get]
func (c *PaymentController) GetAllPayments(ctx *gin.Context) {
	payments, err := c.service.GetAllPayments(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, payments)
}

// GetPaymentByID godoc
// @Summary Get payment by ID
// @Description Get a payment by its ID
// @Tags payments
// @Produce json
// @Param id path string true "Payment ID"
// @Success 200 {object} models.PaymentResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /payments/{id} [get]
func (c *PaymentController) GetPaymentByID(ctx *gin.Context) {
	id := ctx.Param("id")

	payment, err := c.service.GetPaymentByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, payment)
}

// DeletePayment godoc
// @Summary Delete payment by ID
// @Description Delete a payment by its ID
// @Tags payments
// @Produce json
// @Param id path string true "Payment ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /payments/{id} [delete]
func (c *PaymentController) DeletePayment(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.service.DeletePayment(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Payment deleted successfully"})
}
