package repository

import (
	"TA-management/internal/modules/announce/dto/request"
	"TA-management/internal/modules/announce/dto/response"
)

type AnnouncementRepository interface {
	GetStudentEmailByCourseIDs() (*request.EmailRequest, error)
	GetStudentEmailByCourseID(courseID int) (*request.EmailRequest, *response.CourseDetail, error)
	GetStudentEmailByStudentID(studentID int) (*request.EmailRequest, error)
	SaveEmailHistory(rq request.CreateEmailHistory) error
	GetEmailHistory() (*[]response.EmailHistory, error)
	GetDiscordRoleID(courseID int) (string, error)
	CreateNewDiscordChannel(roleID string, channelID string, channelName string, courseID int) error
}
