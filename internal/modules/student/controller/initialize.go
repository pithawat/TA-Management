package controller

import (
	"TA-management/internal/modules/student/service"

	"github.com/gin-gonic/gin"
)

func InitializeController(studentSvc service.StudentService, router *gin.RouterGroup) {
	controller := NewStudentController(studentSvc)

	// Student profile routes
	router.GET("/:id", controller.GetStudentProfile)
	router.PUT("/:id", controller.UpdateStudentProfile)

	// Document routes
	router.POST("/:id/documents/:type", controller.UploadDocument)
	router.GET("/:id/documents/:type", controller.GetDocument)
	router.DELETE("/:id/documents/:type", controller.DeleteDocument)
}
