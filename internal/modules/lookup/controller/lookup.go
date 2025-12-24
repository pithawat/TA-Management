package controller

import (
	"TA-management/internal/modules/lookup/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LookupController struct {
	service service.LookupService
}

func NewLookupController(lookupService service.LookupService) *LookupController {
	return &LookupController{
		service: lookupService,
	}
}

func InitializeController(lookupService service.LookupService, r *gin.RouterGroup) {
	c := NewLookupController(lookupService)
	r.Use()
	{
		r.GET("/course-program", c.getCourseProgram)
		r.GET("/classday", c.getClassday)
		r.GET("/semester", c.getSemester)
		r.GET("/grade", c.getGrade)
		r.GET("/professors", c.getProfessors)
	}
}

func (controller LookupController) getCourseProgram(ctx *gin.Context) {

	result, err := controller.service.GetCourseProgram()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "something went wrong"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (controller LookupController) getClassday(ctx *gin.Context) {
	result, err := controller.service.GetClassday()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "something went wrong"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (controller LookupController) getSemester(ctx *gin.Context) {
	result, err := controller.service.GetSemester()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (controller LookupController) getGrade(ctx *gin.Context) {
	result, err := controller.service.GetGrade()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (controller LookupController) getProfessors(ctx *gin.Context) {
	result, err := controller.service.GetProfessors()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}
