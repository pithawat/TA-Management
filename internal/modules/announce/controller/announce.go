package controller

import (
	"TA-management/internal/modules/announce/dto/request"
	"TA-management/internal/modules/announce/service"
	"TA-management/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AnnouncementController struct {
	service service.AnnouncementService
}

func NewAnnouncementController(service service.AnnouncementService) *AnnouncementController {
	return &AnnouncementController{service: service}
}

func InitializeController(announcementService service.AnnouncementService, r *gin.RouterGroup) {

	c := NewAnnouncementController(announcementService)
	{
		r.POST("/send-mail/all", c.sendMailToAllCourse)
		r.POST("/send-mail/course", c.sendMailToCourse)
		r.POST("/send-mail/individual", c.sendMailToTA)
		r.GET("/email-history", c.getEmailHistory)
		r.POST("/discord/create-channel", c.createDiscordChannel)
		r.GET("/discord/join-channel/:courseID", c.joinDiscordChannel)
	}
}

func (controller AnnouncementController) sendMailToAllCourse(ctx *gin.Context) {

	var rq request.MailForAllCourse
	if err := ctx.ShouldBindJSON(&rq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"erorr": "Invalid request data"})
		return
	}
	controller.service.SendMailToAllCourse(rq)
	ctx.JSON(201, gin.H{"message": "Email are being send."})

}

func (controller AnnouncementController) sendMailToCourse(ctx *gin.Context) {

	var rq request.MailForCourse
	if err := ctx.ShouldBindJSON(&rq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"erorr": "Invalid request data"})
		return
	}
	controller.service.SendMailToCourse(rq)
	ctx.JSON(201, gin.H{"message": "Email are being send."})
}

func (controller AnnouncementController) sendMailToTA(ctx *gin.Context) {

	var rq request.MailForTA
	if err := ctx.ShouldBindJSON(&rq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"erorr": "Invalid request data"})
		return
	}
	controller.service.SendMailToTA(rq)
	ctx.JSON(201, gin.H{"message": "Email are being send."})
}

func (controller AnnouncementController) getEmailHistory(ctx *gin.Context) {

	result, err := controller.service.GetEmailHistory()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong."})
	}

	ctx.JSON(http.StatusOK, result)
}

func (controller AnnouncementController) joinDiscordChannel(ctx *gin.Context) {
	courseID, ok := utils.ValidateParam(ctx, "courseID")

	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validation Param Failed"})
		return
	}

	url, err := controller.service.JoinDiscordChannel(courseID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, url)

}

func (controller AnnouncementController) createDiscordChannel(ctx *gin.Context) {

	var rq request.CreateDiscordChannel
	if err := ctx.ShouldBindJSON(&rq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"erorr": "Invalid request data"})
		return
	}

	err := controller.service.CreateDiscordChannel(rq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "created channel successfully"})
}
