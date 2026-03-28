package controller

import (
	"TA-management/internal/modules/lookup/dto/request"
	"TA-management/internal/modules/lookup/service"
	"fmt"
	"net/http"
	"strconv"

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

func InitializePublicController(lookupService service.LookupService, r *gin.RouterGroup) {
	c := NewLookupController(lookupService)
	{
		r.GET("/course-program", c.getCourseProgram)
		r.GET("/classday", c.getClassday)
		r.GET("/semester", c.getSemester)
		r.GET("/semester-dropdown", c.getSemesterDropdown)
		r.GET("/grade", c.getGrade)
		r.GET("/professors", c.getProfessors)
		r.GET("/holiday", c.GetHolidays)
		r.GET("/available-months", c.GetAvailableMonths)
	}
}

func InitializeProtectedController(lookupService service.LookupService, r *gin.RouterGroup) {
	c := NewLookupController(lookupService)
	{
		r.POST("/add-semester", c.addSemester)
		r.PATCH("/semester", c.updateSemester)
		r.POST("/semester-active/:semesterID", c.setSemesterActive)
		r.POST("/holiday", c.AddSpecialHoliday)
		r.DELETE("/holiday/:id", c.DeleteHoliday)
		r.GET("/ta", c.getTA)
		r.GET("/transcript", c.GetTranscript)
		r.GET("/bank-account", c.GetBankAccount)
		r.GET("/student-card", c.GetStudentCard)
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

func (controller LookupController) getSemesterDropdown(ctx *gin.Context) {
	result, err := controller.service.GetSemesterDropdown()
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

func (controller LookupController) GetHolidays(ctx *gin.Context) {
	monthStr := ctx.Query("month")
	yearStr := ctx.Query("year")

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid month"})
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year"})
		return
	}

	holidays, err := controller.service.GetHolidays(month, year)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch holidays: %v", err)})
		return
	}

	ctx.JSON(http.StatusOK, holidays)
}

func (controller LookupController) AddSpecialHoliday(ctx *gin.Context) {
	var req request.CreateHoliday
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	req.Type = "special" // Force type to be special

	if err := controller.service.AddSpecialHoliday(req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to add holiday: %v", err)})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Holiday added successfully"})
}

func (controller LookupController) DeleteHoliday(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := controller.service.DeleteHoliday(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete holiday: %v", err)})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Holiday deleted successfully"})
}

func (controller LookupController) getTA(ctx *gin.Context) {

	searchVal := ctx.Query("searchVal")
	if searchVal == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
	}

	result, err := controller.service.GetTA(searchVal)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	ctx.JSON(http.StatusOK, result)

}

func (controller LookupController) GetAvailableMonths(ctx *gin.Context) {
	month, err := strconv.Atoi(ctx.Query("month"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed on get query param"})
		return
	}

	result, err := controller.service.GetAvailableMonths(month)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (controller LookupController) GetTranscript(ctx *gin.Context) {
	studentID, err := strconv.Atoi(ctx.Query("studentID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	result, err := controller.service.GetTranscript(studentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, result)
		return
	}

	if result != nil {
		fileName := "transcript_" + result.FileName + "pdf."
		ctx.Header("Content-Disposition", "inline; filename="+fileName)
		ctx.Header("Content-Type", "application/pdf")

		ctx.Data(http.StatusOK, "application/pdf", result.FileBytes)
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "NO Transcript data found"})
	}

}

func (controller LookupController) GetBankAccount(ctx *gin.Context) {
	studentID, err := strconv.Atoi(ctx.Query("studentID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	result, err := controller.service.GetBankAccount(studentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, result)
		return
	}

	if result != nil {
		fileName := "bank_account_" + result.FileName + "pdf."
		ctx.Header("Content-Disposition", "inline; filename="+fileName)
		ctx.Header("Content-Type", "application/pdf")

		ctx.Data(http.StatusOK, "application/pdf", result.FileBytes)
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "NO Transcript data found"})
	}

}

func (controller LookupController) GetStudentCard(ctx *gin.Context) {
	studentID, err := strconv.Atoi(ctx.Query("studentID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	result, err := controller.service.GetStudentCard(studentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, result)
		return
	}

	if result != nil {
		fileName := "studetn_card_" + result.FileName + "pdf."
		ctx.Header("Content-Disposition", "inline; filename="+fileName)
		ctx.Header("Content-Type", "application/pdf")

		ctx.Data(http.StatusOK, "application/pdf", result.FileBytes)
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "NO Transcript data found"})
	}
}

func (controller LookupController) addSemester(ctx *gin.Context) {
	var rq request.CreateSemester

	if err := ctx.ShouldBindJSON(&rq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := controller.service.AddSemester(rq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Already added semester!"})
}

func (controller LookupController) updateSemester(ctx *gin.Context) {
	var rq request.UpdateSemester

	if err := ctx.ShouldBindJSON(&rq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	result, err := controller.service.UpdateSemester(rq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (controller LookupController) setSemesterActive(ctx *gin.Context) {
	idStr := ctx.Param("semesterID")
	semesterID, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	err = controller.service.SetSemesterActive(semesterID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "success"})
}
