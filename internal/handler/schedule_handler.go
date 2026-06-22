package handler

import (
	"strconv"

	"hos_schedule/internal/pkg/response"
	"hos_schedule/internal/service"

	"github.com/gin-gonic/gin"
)

type ScheduleHandler struct {
	service *service.ScheduleService
}

func NewScheduleHandler(service *service.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{service: service}
}

func (h *ScheduleHandler) List(c *gin.Context) {
	doctorID, _ := strconv.ParseInt(c.Query("doctor_id"), 10, 64)
	departmentID, _ := strconv.ParseInt(c.Query("department_id"), 10, 64)
	date := c.Query("date")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if doctorID > 0 {
		schedules, err := h.service.ListByDoctor(doctorID, startDate, endDate)
		if err != nil {
			response.InternalError(c, err.Error())
			return
		}
		response.Success(c, schedules)
		return
	}

	if departmentID > 0 && date != "" {
		schedules, err := h.service.ListByDepartment(departmentID, date)
		if err != nil {
			response.InternalError(c, err.Error())
			return
		}
		response.Success(c, schedules)
		return
	}

	response.BadRequest(c, "doctor_id or department_id+date required")
}

func (h *ScheduleHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid schedule ID")
		return
	}

	schedule, err := h.service.GetByID(id)
	if err != nil {
		response.NotFound(c, "Schedule not found")
		return
	}
	response.Success(c, schedule)
}
