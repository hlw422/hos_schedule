package router

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"hos_schedule/internal/config"
	"hos_schedule/internal/handler"
	"hos_schedule/internal/middleware"
	"hos_schedule/internal/pkg/response"
	redisutil "hos_schedule/internal/pkg/redis"
	"hos_schedule/internal/pkg/wechat"
	"hos_schedule/internal/repository"
	"hos_schedule/internal/service"
)

func Register(r *gin.Engine, db *gorm.DB, rdb *redis.Client, cfg *config.Config) {
	authService := service.NewAuthService(db, cfg)
	authHandler := handler.NewAuthHandler(authService)

	api := r.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			response.Success(c, gin.H{"status": "ok"})
		})

		api.POST("/auth/login", authHandler.Login)

		auth := api.Group("")
		auth.Use(middleware.Auth(cfg))
		{
			auth.GET("/me", func(c *gin.Context) {
				userID, _ := c.Get("user_id")
				response.Success(c, gin.H{"user_id": userID})
			})
		}

		hospitalRepo := repository.NewHospitalRepo(db)
		hospitalService := service.NewHospitalService(hospitalRepo)
		hospitalHandler := handler.NewHospitalHandler(hospitalService)

		departmentRepo := repository.NewDepartmentRepo(db)
		departmentService := service.NewDepartmentService(departmentRepo)
		departmentHandler := handler.NewDepartmentHandler(departmentService)

		api.GET("/hospitals", hospitalHandler.List)
		api.GET("/hospitals/:id", hospitalHandler.GetByID)
		api.GET("/hospitals/:id/campuses", hospitalHandler.GetCampuses)

		navigationHandler := handler.NewNavigationHandler(hospitalService)
		api.GET("/navigation/:campus_id", navigationHandler.GetNavigationInfo)
		api.GET("/navigation/nearby", navigationHandler.GetNearbyCampuses)

		api.GET("/departments", departmentHandler.List)
		api.GET("/departments/:id", departmentHandler.GetByID)

		doctorRepo := repository.NewDoctorRepo(db)
		doctorService := service.NewDoctorService(doctorRepo)
		doctorHandler := handler.NewDoctorHandler(doctorService)

		scheduleRepo := repository.NewScheduleRepo(db)
		scheduleService := service.NewScheduleService(scheduleRepo)
		scheduleHandler := handler.NewScheduleHandler(scheduleService)

		api.GET("/doctors", doctorHandler.List)
		api.GET("/doctors/:id", doctorHandler.GetByID)

		api.GET("/schedules", scheduleHandler.List)
		api.GET("/schedules/:id", scheduleHandler.GetByID)

		recommendHandler := handler.NewRecommendHandler(doctorService, departmentService)
		api.GET("/recommend/departments", recommendHandler.GetHotDepartments)
		api.GET("/recommend/doctors", recommendHandler.GetRecommendedDoctors)

		patientRepo := repository.NewPatientRepo(db)
		patientService := service.NewPatientService(patientRepo)
		patientHandler := handler.NewPatientHandler(patientService)

		auth.GET("/patients", patientHandler.List)
		auth.POST("/patients", patientHandler.Create)
		auth.PUT("/patients/:id", patientHandler.Update)
		auth.DELETE("/patients/:id", patientHandler.Delete)
		auth.PUT("/patients/:id/default", patientHandler.SetDefault)

		slotManager := redisutil.NewSlotManager(rdb)

		appointmentRepo := repository.NewAppointmentRepo(db)
		appointmentService := service.NewAppointmentService(db, appointmentRepo, scheduleRepo, slotManager)
		appointmentHandler := handler.NewAppointmentHandler(appointmentService)

		auth.POST("/appointments", appointmentHandler.Create)
		auth.GET("/appointments", appointmentHandler.List)
		auth.GET("/appointments/:id", appointmentHandler.GetByID)
		auth.PUT("/appointments/:id/cancel", appointmentHandler.Cancel)

		payClient := wechat.NewPayClient(cfg.Wechat.AppID, cfg.Wechat.MchID, cfg.Wechat.APIKey, cfg.Wechat.NotifyURL)
		paymentHandler := handler.NewPaymentHandler(appointmentService, payClient)
		api.POST("/payments/callback", paymentHandler.WechatPayCallback)
		auth.POST("/payments/:id/create", paymentHandler.CreatePayment)

		notificationRepo := repository.NewNotificationRepo(db)
		wechatClient := wechat.NewClient(&cfg.Wechat)
		notificationService := service.NewNotificationService(notificationRepo, wechatClient, &cfg.SMS)
		notificationHandler := handler.NewNotificationHandler(notificationService)

		auth.POST("/notifications/subscribe", notificationHandler.Subscribe)

		auth.GET("/doctor/schedules", doctorHandler.GetMySchedules)
		auth.GET("/doctor/appointments", doctorHandler.GetTodayAppointments)

		adminHandler := handler.NewAdminHandler(
			hospitalService, departmentService, doctorService, scheduleService, appointmentService, slotManager,
		)

		admin := api.Group("/admin")
		admin.Use(middleware.Auth(cfg))
		admin.Use(middleware.AdminRequired())
		{
			admin.PUT("/hospitals/:id", adminHandler.UpdateHospital)
			admin.POST("/hospitals/campuses", adminHandler.AddCampus)

			admin.POST("/departments", adminHandler.CreateDepartment)
			admin.PUT("/departments/:id", adminHandler.UpdateDepartment)
			admin.DELETE("/departments/:id", adminHandler.DeleteDepartment)

			admin.POST("/doctors", adminHandler.CreateDoctor)
			admin.PUT("/doctors/:id", adminHandler.UpdateDoctor)
			admin.PUT("/doctors/:id/status", adminHandler.UpdateDoctorStatus)

			admin.POST("/schedules", adminHandler.CreateSchedule)
			admin.POST("/schedules/batch", adminHandler.BatchCreateSchedule)
			admin.PUT("/schedules/:id", adminHandler.UpdateSchedule)
			admin.DELETE("/schedules/:id", adminHandler.DeleteSchedule)

			admin.GET("/appointments", adminHandler.ListAppointments)
			admin.GET("/appointments/stats", adminHandler.GetAppointmentStats)
		}
	}
}
