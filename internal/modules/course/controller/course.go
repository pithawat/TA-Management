package controller

import (
	"TA-management/internal/modules/course/dto/request"
	"TA-management/internal/modules/course/service"
	"TA-management/internal/utils"
	"fmt"
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
		r.GET("/student/:studentId", c.GetAllCourseByStudentId)
		r.GET("/professor/:professorId", c.findProfessorCourse)
		r.POST("", c.createCourse)
		r.PATCH("/:courseId", c.updateCourse)
		r.DELETE("/:courseId", c.deleteCourse)
		r.POST("/apply/:jobPostId", c.applyJobPost)
		r.GET("/application/student/:studentId", c.getApplicationByStudentId)
		r.GET("/application/course/:courseId", c.getApplicationBycourseId)
		r.GET("/application/professor/:professorId", c.getApplicationByProfessorId)
		r.GET("/application/:applilcationId", c.getApplicationDetail)
		r.GET("/application/transcript/:applicationId", c.getApplicationtranscriptPdf)
		r.GET("/applicton/bankaccount/:applicationId", c.getApplicationbankAccountPdf)
		r.GET("/application/studentcard/:applicationId", c.getApplicationstudentCardPdf)
		r.POST("/application/approve/:applicationId", c.approveApplication)
	}
}

func (controller CourseController) findAllCourse(ctx *gin.Context) {
	//validate
	result, err := controller.service.GetAllCourse()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (controller CourseController) GetAllCourseByStudentId(ctx *gin.Context) {

	studentId, ok := utils.ValidateParam(ctx, "studentId")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validation Param Failed"})
		return
	}

	result, err := controller.service.GetAllCourseByStudentId(studentId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
	}

	ctx.JSON(http.StatusOK, result)
}

func (controller CourseController) findProfessorCourse(ctx *gin.Context) {

	//validate
	professorId, ok := utils.ValidateParam(ctx, "professorId")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validation Param Failed"})
		return
	}

	result, err := controller.service.GetProfessorCourse(professorId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (controller CourseController) createCourse(ctx *gin.Context) {
	var request request.CreateCourse
	if err := ctx.ShouldBindJSON(&request); err != nil {
		fmt.Println(err)
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validation Param Failed"})
		return
	}

	rq := request.UpdateCourse{}
	if err := ctx.ShouldBindJSON(&rq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "body request not valid"})
		return
	}

	rq.Id = courseId
	result, err := controller.service.UpdateCourse(rq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}
	ctx.JSON(http.StatusNoContent, result)
}

func (controller CourseController) deleteCourse(ctx *gin.Context) {
	id, ok := utils.ValidateParam(ctx, "courseId")

	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validate Param Failed"})
		return
	}

	result, err := controller.service.DeleteCourse(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}
	ctx.JSON(http.StatusNoContent, result)

}

func (controller CourseController) applyJobPost(ctx *gin.Context) {
	rq := request.ApplyJobPost{}
	id, ok := utils.ValidateParam(ctx, "jobPostId")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validate Param Failed."})
		return
	}

	transcriptName, transcriptBytes, err := utils.GetFileData(ctx, "Transcript")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	bankAccountName, bankAccountBytes, err := utils.GetFileData(ctx, "BankAccount")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	studentCardName, studentCardBytes, err := utils.GetFileData(ctx, "StudentCard")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if err := ctx.ShouldBind(&rq); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "body request not valid"})
		return
	}

	rq.JobPostID = &id
	rq.TranscriptName = &transcriptName
	rq.TranscriptBytes = transcriptBytes
	rq.BankAccountName = &bankAccountName
	rq.BankAccountBytes = bankAccountBytes
	rq.StudentCardName = &studentCardName
	rq.StudentCardBytes = studentCardBytes

	result, err := controller.service.ApplyJobPost(rq)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, result)
		return
	}

	ctx.JSON(http.StatusCreated, result)

}

func (controller CourseController) getApplicationByStudentId(ctx *gin.Context) {
	id, ok := utils.ValidateParam(ctx, "studentId")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validate Param Failed."})
		return
	}

	result, err := controller.service.GetApplicationByStudentId(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, result)
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (controller CourseController) getApplicationByProfessorId(ctx *gin.Context) {
	id, ok := utils.ValidateParam(ctx, "professorId")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validate Param Failed."})
		return
	}

	result, err := controller.service.GetApplicationByProfessorId(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, result)
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (controller CourseController) getApplicationBycourseId(ctx *gin.Context) {
	id, ok := utils.ValidateParam(ctx, "courseId")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validate Param Failed."})
		return
	}

	result, err := controller.service.GetApplicationByCourseId(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, result)
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (controller CourseController) getApplicationDetail(ctx *gin.Context) {
	id, ok := utils.ValidateParam(ctx, "applicationId")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validate Param Failed."})
		return
	}

	result, err := controller.service.GetApplicationDetail(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, result)
		return
	}

	ctx.JSON(http.StatusOK, result)

}

func (controller CourseController) getApplicationtranscriptPdf(ctx *gin.Context) {
	id, ok := utils.ValidateParam(ctx, "applicationId")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validate Param Failed."})
		return
	}

	result, err := controller.service.GetApplicationTranscriptPdf(id)
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

func (controller CourseController) getApplicationbankAccountPdf(ctx *gin.Context) {
	id, ok := utils.ValidateParam(ctx, "applicationId")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validate Param Failed."})
		return
	}
	result, err := controller.service.GetApplicationBankAccountPdf(id)
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
		ctx.JSON(http.StatusNotFound, gin.H{"error": "NO Bank Account data found"})
	}
}

func (controller CourseController) getApplicationstudentCardPdf(ctx *gin.Context) {
	id, ok := utils.ValidateParam(ctx, "applicationId")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validate Param Failed."})
		return
	}
	result, err := controller.service.GetApplicationStudentCardPdf(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, result)
		return
	}
	if result != nil {
		fileName := "student_card_" + result.FileName + "pdf."
		ctx.Header("Content-Disposition", "inline; filename="+fileName)
		ctx.Header("Content-Type", "application/pdf")
		ctx.Data(http.StatusOK, "application/pdf", result.FileBytes)
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "NO Student Card data found"})
	}
}

func (controller CourseController) approveApplication(ctx *gin.Context) {
	id, ok := utils.ValidateParam(ctx, "applicationId")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validate Param Failed."})
		return
	}

	result, err := controller.service.ApproveApplication(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, result)
		return
	}
	ctx.JSON(http.StatusCreated, result)
}
