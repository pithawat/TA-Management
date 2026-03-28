package service

import (
	"TA-management/internal/modules/announce/dto/request"
	"TA-management/internal/modules/announce/dto/response"
)

type AnnouncementService interface {
	SendMailToAllCourse(rq request.MailForAllCourse)
	SendMailToCourse(rq request.MailForCourse)
	SendMailToTA(rq request.MailForTA)
	GetEmailHistory() (*[]response.EmailHistory, error)
	JoinDiscordChannel(courseID int) (string, error)
	CreateDiscordChannel(rq request.CreateDiscordChannel) error
}
