package handler

import (
	"hos_schedule/internal/pkg/response"
	"hos_schedule/internal/service"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	appointmentService *service.AppointmentService
}

func NewPaymentHandler(appointmentService *service.AppointmentService) *PaymentHandler {
	return &PaymentHandler{appointmentService: appointmentService}
}

func (h *PaymentHandler) WechatPayCallback(c *gin.Context) {
	var req struct {
		PayID   string  `json:"pay_id" binding:"required"`
		Amount  float64 `json:"amount" binding:"required"`
		Success bool    `json:"success"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request")
		return
	}
	if !req.Success {
		response.Error(c, 400, "Payment failed")
		return
	}
	err := h.appointmentService.HandlePaymentCallback(c.Request.Context(), req.PayID, req.Amount)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, gin.H{"status": "ok"})
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	appointmentID := c.Param("id")
	response.Success(c, gin.H{
		"appointment_id": appointmentID,
		"status":         "pending",
		"message":        "WeChat Pay integration pending",
	})
}
