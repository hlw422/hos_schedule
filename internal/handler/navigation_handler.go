package handler

import (
	"strconv"

	"hos_schedule/internal/pkg/response"
	"hos_schedule/internal/service"

	"github.com/gin-gonic/gin"
)

type NavigationHandler struct {
	hospitalService *service.HospitalService
}

func NewNavigationHandler(hospitalService *service.HospitalService) *NavigationHandler {
	return &NavigationHandler{hospitalService: hospitalService}
}

type NavigationInfo struct {
	CampusID  int64   `json:"campus_id"`
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Phone     string  `json:"phone"`
}

func (h *NavigationHandler) GetNavigationInfo(c *gin.Context) {
	campusID, err := strconv.ParseInt(c.Param("campus_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid campus ID")
		return
	}

	campus, err := h.hospitalService.GetCampusByID(campusID)
	if err != nil {
		response.NotFound(c, "Campus not found")
		return
	}

	response.Success(c, NavigationInfo{
		CampusID:  campus.ID,
		Name:      campus.Name,
		Address:   campus.Address,
		Latitude:  campus.Latitude,
		Longitude: campus.Longitude,
		Phone:     campus.Phone,
	})
}

func (h *NavigationHandler) GetNearbyCampuses(c *gin.Context) {
	lat, _ := strconv.ParseFloat(c.Query("lat"), 64)
	lng, _ := strconv.ParseFloat(c.Query("lng"), 64)

	_ = lat
	_ = lng

	campuses, err := h.hospitalService.GetAllCampuses()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, campuses)
}
