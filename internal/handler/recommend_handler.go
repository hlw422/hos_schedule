package handler

import (
	"hos_schedule/internal/pkg/response"
	"hos_schedule/internal/service"

	"github.com/gin-gonic/gin"
)

type RecommendHandler struct {
	doctorService     *service.DoctorService
	departmentService *service.DepartmentService
}

func NewRecommendHandler(doctorService *service.DoctorService, departmentService *service.DepartmentService) *RecommendHandler {
	return &RecommendHandler{
		doctorService:     doctorService,
		departmentService: departmentService,
	}
}

func (h *RecommendHandler) GetHotDepartments(c *gin.Context) {
	departments, err := h.departmentService.ListByHospital(1)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	if len(departments) > 8 {
		departments = departments[:8]
	}
	response.Success(c, departments)
}

func (h *RecommendHandler) GetRecommendedDoctors(c *gin.Context) {
	limit := 10
	doctors, err := h.doctorService.ListRecommended(limit)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, doctors)
}
