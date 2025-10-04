package controller

import (
	"TA-management/internal/modules/course/dto/request"
	"TA-management/internal/modules/course/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CourseController struct {
	service service.CourseService
}

func NewCourseController(courseService service.CourseService) *CourseController {
	return &CourseController{
		service: courseService,
	}
}

func InitializeController(courseService service.CourseService, r *gin.RouterGroup) {
	c := NewCourseController(courseService)
	r.Use()
	{
		r.GET("", c.findAllCourse)
		r.POST("", c.createCourse)
	}
}

func (controller CourseController) findAllCourse(ctx *gin.Context) {
	//validate
	result, err := controller.service.GetAllCourse()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (controller CourseController) createCourse(ctx *gin.Context) {
	var request request.CreateCourse
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request data"})
		return
	}
	err := controller.service.CreateCourse(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"sucess": "created successfully"})
}
