package response

import "time"

type CourseDetail struct {
	CourseCode string
	CourseName string
}

type EmailHistory struct {
	Id           int       `json:"id"`
	Subject      string    `json:"subject"`
	Body         string    `json:"body"`
	ReceivedName string    `json:"receivedName"`
	NReceived    int       `json:"nReceived"`
	Status       string    `json:"status"`
	CreatedDate  time.Time `json:"createDate"`
}
