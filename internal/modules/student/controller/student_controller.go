package controller

import (
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"TA-management/internal/modules/student/dto/request"
	"TA-management/internal/modules/student/service"

	"github.com/gin-gonic/gin"
)

type StudentController struct {
	service service.StudentService
}

func NewStudentController(service service.StudentService) *StudentController {
	return &StudentController{service: service}
}

// GetStudentProfile godoc
// @Summary Get student profile
// @Description Get student profile by ID
// @Tags student
// @Accept json
// @Produce json
// @Param id path int true "Student ID"
// @Success 200 {object} response.StudentProfile
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /student/{id} [get]
func (c *StudentController) GetStudentProfile(ctx *gin.Context) {
	studentID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	profile, err := c.service.GetStudentProfile(studentID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, profile)
}

// UpdateStudentProfile godoc
// @Summary Update student profile
// @Description Update student Thai name and phone number
// @Tags student
// @Accept json
// @Produce json
// @Param id path int true "Student ID"
// @Param request body request.UpdateProfile true "Update Profile Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /student/{id} [put]
func (c *StudentController) UpdateStudentProfile(ctx *gin.Context) {
	studentID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	var req request.UpdateProfile
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.UpdateStudentProfile(studentID, &req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

// UploadDocument godoc
// @Summary Upload document
// @Description Upload or replace student document
// @Tags student
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Student ID"
// @Param type path string true "Document type" Enums(transcript, bank-account, student-card)
// @Param file formData file true "Document file"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /student/{id}/documents/{type} [post]
func (c *StudentController) UploadDocument(ctx *gin.Context) {
	studentID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	docType := ctx.Param("type")
	// Convert URL param format to internal format
	if docType == "bank-account" {
		docType = "bank_account"
	} else if docType == "student-card" {
		docType = "student_card"
	}

	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	defer file.Close()

	// Validate file type (PDF only)
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".pdf" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Only PDF files are allowed"})
		return
	}

	if err := c.service.UploadDocument(studentID, docType, file, header.Filename); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Document uploaded successfully"})
}

// GetDocument godoc
// @Summary Get document
// @Description View student document
// @Tags student
// @Produce application/pdf
// @Param id path int true "Student ID"
// @Param type path string true "Document type" Enums(transcript, bank-account, student-card)
// @Success 200 {file} binary
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /student/{id}/documents/{type} [get]
func (c *StudentController) GetDocument(ctx *gin.Context) {
	studentID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	docType := ctx.Param("type")
	// Convert URL param format to internal format
	if docType == "bank-account" {
		docType = "bank_account"
	} else if docType == "student-card" {
		docType = "student_card"
	}

	fileBytes, fileName, err := c.service.GetDocument(studentID, docType)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("Content-Disposition", "inline; filename="+fileName)
	ctx.Header("Content-Type", "application/pdf")
	ctx.Data(http.StatusOK, "application/pdf", fileBytes)
}

// DeleteDocument godoc
// @Summary Delete document
// @Description Delete student document
// @Tags student
// @Accept json
// @Produce json
// @Param id path int true "Student ID"
// @Param type path string true "Document type" Enums(transcript, bank-account, student-card)
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /student/{id}/documents/{type} [delete]
func (c *StudentController) DeleteDocument(ctx *gin.Context) {
	studentID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	docType := ctx.Param("type")
	// Convert URL param format to internal format
	if docType == "bank-account" {
		docType = "bank_account"
	} else if docType == "student-card" {
		docType = "student_card"
	}

	if err := c.service.DeleteDocument(studentID, docType); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})
}
