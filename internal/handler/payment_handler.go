package handler

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"hos_schedule/internal/pkg/response"
	"hos_schedule/internal/pkg/wechat"
	"hos_schedule/internal/service"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	appointmentService *service.AppointmentService
	payClient          *wechat.PayClient
}

func NewPaymentHandler(appointmentService *service.AppointmentService, payClient *wechat.PayClient) *PaymentHandler {
	return &PaymentHandler{
		appointmentService: appointmentService,
		payClient:          payClient,
	}
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	idStr := c.Param("id")
	appointmentID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid appointment ID")
		return
	}

	appointment, err := h.appointmentService.GetByID(appointmentID)
	if err != nil {
		response.NotFound(c, "Appointment not found")
		return
	}

	if appointment.Status != "PENDING_PAY" {
		response.Error(c, 400, "Appointment is not pending payment")
		return
	}

	payID := fmt.Sprintf("PAY%d%d", appointmentID, time.Now().UnixMilli())
	totalFee := int(appointment.PayAmount * 100)

	prepayID, err := h.payClient.CreatePrepayOrder(
		payID,
		fmt.Sprintf("医院预约挂号-%s", appointment.Date),
		"",
		totalFee,
		c.ClientIP(),
	)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	if err := h.appointmentService.UpdatePayID(appointmentID, payID); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"prepay_id": prepayID,
		"pay_id":    payID,
		"amount":    appointment.PayAmount,
	})
}

func (h *PaymentHandler) WechatPayCallback(c *gin.Context) {
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.Error(c, 400, "Invalid request body")
		return
	}

	notification, err := h.payClient.VerifyNotification(data)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	amount := float64(notification.TotalFee) / 100
	if err := h.appointmentService.HandlePaymentCallback(c.Request.Context(), notification.OutTradeNo, amount); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	c.XML(200, map[string]string{
		"return_code": "SUCCESS",
		"return_msg":  "OK",
	})
}
