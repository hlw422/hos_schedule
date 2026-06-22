package handler

import (
	"strconv"

	"hos_schedule/internal/model"
	"hos_schedule/internal/pkg/response"
	"hos_schedule/internal/service"

	"github.com/gin-gonic/gin"
)

type PatientHandler struct {
	service *service.PatientService
}

func NewPatientHandler(service *service.PatientService) *PatientHandler {
	return &PatientHandler{service: service}
}

func (h *PatientHandler) List(c *gin.Context) {
	userID := c.GetInt64("user_id")
	patients, err := h.service.ListByUser(userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, patients)
}

func (h *PatientHandler) Create(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var patient model.Patient
	if err := c.ShouldBindJSON(&patient); err != nil {
		response.BadRequest(c, "Invalid request")
		return
	}

	patient.UserID = userID
	if err := h.service.Create(&patient); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, patient)
}

func (h *PatientHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid patient ID")
		return
	}

	var patient model.Patient
	if err := c.ShouldBindJSON(&patient); err != nil {
		response.BadRequest(c, "Invalid request")
		return
	}

	patient.ID = id
	if err := h.service.Update(&patient); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, patient)
}

func (h *PatientHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid patient ID")
		return
	}

	if err := h.service.Delete(id); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *PatientHandler) SetDefault(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid patient ID")
		return
	}

	if err := h.service.SetDefault(userID, id); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, nil)
}
