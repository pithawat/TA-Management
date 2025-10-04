package router

import (
	"TA-management/config"
	"TA-management/internal/modules/course/controller"
	"TA-management/internal/modules/course/repository"
	"TA-management/internal/modules/course/service"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	db := config.ConnectDatabase()

	courseRepo := repository.NewCourseRepository(db)
	courseService := service.NewCourseService(courseRepo)
	// courseController := controller.NewCourseController(courseService)

	baseRouter := r.Group("/TA-management")

	baseRouter.Use()
	{
		courseRouter := baseRouter.Group("/course")
		controller.InitializeController(courseService, courseRouter)
	}

	return r

}
