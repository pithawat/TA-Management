package controller

import (
	"TA-management/internal/modules/course/dto/request"
	"TA-management/internal/modules/course/service"
	"TA-management/internal/utils"
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
		r.PATCH("/:courseId", c.updateCourse)
		r.DELETE("/:courseId", c.deleteCourse)
		// r.POST("/apply/:courseId", c.applyCourse)
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
	result, err := controller.service.CreateCourse(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	ctx.JSON(http.StatusCreated, result)
}

func (controller CourseController) updateCourse(ctx *gin.Context) {
	courseId, ok := utils.ValidateParam(ctx, "courseId")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Validation Param Failed"})
		return
	}

	rq := request.UpdateCourse{}
	if err := ctx.ShouldBindJSON(&rq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "body request not valid"})
		return
	}

	rq.Id = courseId
	result, err := controller.service.UpdateCourse(rq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "something went wrong"})
		return
	}
	ctx.JSON(http.StatusNoContent, result)
}

func (controller CourseController) deleteCourse(ctx *gin.Context) {
	id, ok := utils.ValidateParam(ctx, "courseId")

	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Validate Param Failed"})
		return
	}

	result, err := controller.service.DeleteCourse(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "something went wrong"})
		return
	}
	ctx.JSON(http.StatusNoContent, result)

}
