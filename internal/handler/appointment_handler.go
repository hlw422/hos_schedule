package handler

import (
	"strconv"

	"hos_schedule/internal/pkg/response"
	"hos_schedule/internal/service"

	"github.com/gin-gonic/gin"
)

type AppointmentHandler struct {
	service *service.AppointmentService
}

func NewAppointmentHandler(service *service.AppointmentService) *AppointmentHandler {
	return &AppointmentHandler{service: service}
}

func (h *AppointmentHandler) Create(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var req service.CreateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request")
		return
	}

	appointment, err := h.service.Create(c.Request.Context(), userID, &req)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, appointment)
}

func (h *AppointmentHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid appointment ID")
		return
	}

	appointment, err := h.service.GetByID(id)
	if err != nil {
		response.NotFound(c, "Appointment not found")
		return
	}
	response.Success(c, appointment)
}

func (h *AppointmentHandler) List(c *gin.Context) {
	userID := c.GetInt64("user_id")
	status := c.Query("status")

	appointments, err := h.service.ListByUser(userID, status)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, appointments)
}

func (h *AppointmentHandler) Cancel(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid appointment ID")
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&req)

	if err := h.service.Cancel(c.Request.Context(), id, req.Reason); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, nil)
}
