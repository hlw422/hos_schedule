package handler

import (
	"strconv"

	"hos_schedule/internal/pkg/response"
	"hos_schedule/internal/service"

	"github.com/gin-gonic/gin"
)

type DoctorHandler struct {
	service *service.DoctorService
}

func NewDoctorHandler(service *service.DoctorService) *DoctorHandler {
	return &DoctorHandler{service: service}
}

func (h *DoctorHandler) List(c *gin.Context) {
	departmentID, _ := strconv.ParseInt(c.Query("department_id"), 10, 64)
	if departmentID > 0 {
		doctors, err := h.service.ListByDepartment(departmentID)
		if err != nil {
			response.InternalError(c, err.Error())
			return
		}
		response.Success(c, doctors)
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	doctors, err := h.service.ListRecommended(limit)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, doctors)
}

func (h *DoctorHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid doctor ID")
		return
	}

	doctor, err := h.service.GetByID(id)
	if err != nil {
		response.NotFound(c, "Doctor not found")
		return
	}
	response.Success(c, doctor)
}

func (h *DoctorHandler) GetMySchedules(c *gin.Context) {
	userID := c.GetInt64("user_id")

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	response.Success(c, gin.H{
		"user_id":    userID,
		"start_date": startDate,
		"end_date":   endDate,
	})
}

func (h *DoctorHandler) GetTodayAppointments(c *gin.Context) {
	userID := c.GetInt64("user_id")

	response.Success(c, gin.H{
		"user_id": userID,
	})
}
