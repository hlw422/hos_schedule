package handler

import (
	"strconv"

	"hos_schedule/internal/pkg/response"
	"hos_schedule/internal/service"

	"github.com/gin-gonic/gin"
)

type HospitalHandler struct {
	service *service.HospitalService
}

func NewHospitalHandler(service *service.HospitalService) *HospitalHandler {
	return &HospitalHandler{service: service}
}

func (h *HospitalHandler) List(c *gin.Context) {
	hospitals, err := h.service.List()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, hospitals)
}

func (h *HospitalHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid hospital ID")
		return
	}

	hospital, err := h.service.GetByID(id)
	if err != nil {
		response.NotFound(c, "Hospital not found")
		return
	}
	response.Success(c, hospital)
}

func (h *HospitalHandler) GetCampuses(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid hospital ID")
		return
	}

	campuses, err := h.service.GetCampuses(id)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, campuses)
}
