package handler

import (
	"strconv"

	"hos_schedule/internal/model"
	"hos_schedule/internal/pkg/response"
	"hos_schedule/internal/service"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	hospitalService    *service.HospitalService
	departmentService  *service.DepartmentService
	doctorService      *service.DoctorService
	scheduleService    *service.ScheduleService
	appointmentService *service.AppointmentService
}

func NewAdminHandler(
	hospitalService *service.HospitalService,
	departmentService *service.DepartmentService,
	doctorService *service.DoctorService,
	scheduleService *service.ScheduleService,
	appointmentService *service.AppointmentService,
) *AdminHandler {
	return &AdminHandler{
		hospitalService:    hospitalService,
		departmentService:  departmentService,
		doctorService:      doctorService,
		scheduleService:    scheduleService,
		appointmentService: appointmentService,
	}
}

func (h *AdminHandler) UpdateHospital(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid hospital ID")
		return
	}

	var hospital model.Hospital
	if err := c.ShouldBindJSON(&hospital); err != nil {
		response.BadRequest(c, "Invalid request")
		return
	}
	hospital.ID = id
	// TODO: 调用 service 更新医院信息
	response.Success(c, hospital)
}

func (h *AdminHandler) AddCampus(c *gin.Context) {
	var campus model.Campus
	if err := c.ShouldBindJSON(&campus); err != nil {
		response.BadRequest(c, "Invalid request")
		return
	}
	// TODO: 调用 service 添加院区
	response.Success(c, campus)
}

func (h *AdminHandler) CreateDepartment(c *gin.Context) {
	var dept model.Department
	if err := c.ShouldBindJSON(&dept); err != nil {
		response.BadRequest(c, "Invalid request")
		return
	}
	// TODO: 调用 service 创建
	response.Success(c, dept)
}

func (h *AdminHandler) UpdateDepartment(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid department ID")
		return
	}

	var dept model.Department
	if err := c.ShouldBindJSON(&dept); err != nil {
		response.BadRequest(c, "Invalid request")
		return
	}
	dept.ID = id
	// TODO: 调用 service 更新
	response.Success(c, dept)
}

func (h *AdminHandler) DeleteDepartment(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid department ID")
		return
	}
	// TODO: 调用 service 删除
	response.Success(c, gin.H{"id": id})
}

func (h *AdminHandler) CreateDoctor(c *gin.Context) {
	var doc model.Doctor
	if err := c.ShouldBindJSON(&doc); err != nil {
		response.BadRequest(c, "Invalid request")
		return
	}
	// TODO: 调用 service 创建
	response.Success(c, doc)
}

func (h *AdminHandler) UpdateDoctor(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid doctor ID")
		return
	}

	var doc model.Doctor
	if err := c.ShouldBindJSON(&doc); err != nil {
		response.BadRequest(c, "Invalid request")
		return
	}
	doc.ID = id
	// TODO: 调用 service 更新
	response.Success(c, doc)
}

func (h *AdminHandler) UpdateDoctorStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid doctor ID")
		return
	}

	var req struct {
		Status int8 `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request")
		return
	}
	// TODO: 调用 service 启用/停用医生
	response.Success(c, gin.H{"id": id, "status": req.Status})
}

func (h *AdminHandler) CreateSchedule(c *gin.Context) {
	var sched model.Schedule
	if err := c.ShouldBindJSON(&sched); err != nil {
		response.BadRequest(c, "Invalid request")
		return
	}
	// TODO: 调用 service 创建，初始化 Redis 号源
	response.Success(c, sched)
}

func (h *AdminHandler) BatchCreateSchedule(c *gin.Context) {
	var req struct {
		DoctorID   int64    `json:"doctor_id"`
		CampusID   int64    `json:"campus_id"`
		Dates      []string `json:"dates"`
		TimePeriod string   `json:"time_period"`
		TotalCount int      `json:"total_count"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request")
		return
	}
	// TODO: 批量创建排班
	response.Success(c, req)
}

func (h *AdminHandler) UpdateSchedule(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid schedule ID")
		return
	}

	var sched model.Schedule
	if err := c.ShouldBindJSON(&sched); err != nil {
		response.BadRequest(c, "Invalid request")
		return
	}
	sched.ID = id
	// TODO: 调用 service 更新
	response.Success(c, sched)
}

func (h *AdminHandler) DeleteSchedule(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid schedule ID")
		return
	}
	// TODO: 调用 service 删除，释放 Redis 号源
	response.Success(c, gin.H{"id": id})
}

func (h *AdminHandler) ListAppointments(c *gin.Context) {
	doctorID, _ := strconv.ParseInt(c.Query("doctor_id"), 10, 64)
	date := c.Query("date")
	status := c.Query("status")

	_ = doctorID
	_ = date
	_ = status
	// TODO: 调用 service 查询预约列表
	response.Success(c, []interface{}{})
}

func (h *AdminHandler) GetAppointmentStats(c *gin.Context) {
	// TODO: 统计今日预约量、取消量、到诊率
	response.Success(c, gin.H{
		"today_total":   0,
		"today_cancel":  0,
		"today_visited": 0,
	})
}
