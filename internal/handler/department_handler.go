package handler

import (
	"strconv"

	"hos_schedule/internal/pkg/response"
	"hos_schedule/internal/service"

	"github.com/gin-gonic/gin"
)

type DepartmentHandler struct {
	service *service.DepartmentService
}

func NewDepartmentHandler(service *service.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{service: service}
}

func (h *DepartmentHandler) List(c *gin.Context) {
	campusID, _ := strconv.ParseInt(c.Query("campus_id"), 10, 64)
	hospitalID, _ := strconv.ParseInt(c.Query("hospital_id"), 10, 64)

	var departments interface{}
	var err error

	if campusID > 0 {
		departments, err = h.service.ListByCampus(campusID)
	} else if hospitalID > 0 {
		departments, err = h.service.ListByHospital(hospitalID)
	} else {
		response.BadRequest(c, "campus_id or hospital_id required")
		return
	}

	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, departments)
}

func (h *DepartmentHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid department ID")
		return
	}

	department, err := h.service.GetByID(id)
	if err != nil {
		response.NotFound(c, "Department not found")
		return
	}
	response.Success(c, department)
}
