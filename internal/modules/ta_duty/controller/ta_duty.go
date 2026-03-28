package controller

import (
	"TA-management/internal/modules/ta_duty/dto/request"
	"TA-management/internal/modules/ta_duty/service"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TaDutyController struct {
	service service.TaDutyService
}

func NewTaDutyController(taDutyService service.TaDutyService) TaDutyController {
	return TaDutyController{service: taDutyService}
}

func InitializeController(taDutyService service.TaDutyService, r *gin.RouterGroup) {
	// r.Use() // Middleware if needed
	c := NewTaDutyController(taDutyService)
	{
		r.GET("/duty-roadmap", c.getTADutyRoadmap)
		r.POST("/marked-duty", c.markDutyAsDone)
		r.POST("/export-payment-report", c.exportPaymentReport)
		r.POST("/export-signature-sheet", c.exportSignatureSheet)
	}
}

func (controller TaDutyController) getTADutyRoadmap(ctx *gin.Context) {
	courseID, err := strconv.Atoi(ctx.Query("courseID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get query params : courseID "})
		return
	}

	studentID, err := strconv.Atoi(ctx.Query("studentID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get query params : studentID "})
		return
	}

	result, err := controller.service.GetTADutyRoadmap(courseID, studentID)

	ctx.JSON(http.StatusOK, result)
}

func (controller TaDutyController) markDutyAsDone(ctx *gin.Context) {
	courseID, err := strconv.Atoi(ctx.Query("courseID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get query params : courseID "})
		return
	}

	studentID, err := strconv.Atoi(ctx.Query("studentID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get query params : studentID "})
		return
	}

	dutyDate := ctx.Query("dutyDate")

	result, err := controller.service.MarkDutyAsDone(courseID, studentID, dutyDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}
	ctx.JSON(http.StatusCreated, result)
}

func (controller TaDutyController) exportPaymentReport(ctx *gin.Context) {

	var rq request.ExportPaymentReportRequest
	if err := ctx.ShouldBindJSON(&rq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind query params"})
		return
	}
	fmt.Println(rq)
	buffer, courseData, err := controller.service.ExportPaymentReport(rq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	fileName := fmt.Sprintf("Payment_Report_%s-sec(%s)(%s).xlsx", courseData.CourseName, courseData.Sec, courseData.Semester)
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.sheet")
	ctx.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buffer.Bytes())
}

func (controller TaDutyController) exportSignatureSheet(ctx *gin.Context) {
	var rq request.ExportSignatureSheet

	if err := ctx.ShouldBindJSON(&rq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reqeust data"})
		return
	}

	buffer, courseData, err := controller.service.ExportSignatureSheet(rq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	fileName := fmt.Sprintf("Signature_sheet_%s-sec(%s)(%s).xlsx", courseData.CourseName, courseData.Sec, courseData.Semester)
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.sheet")
	ctx.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buffer.Bytes())

}
