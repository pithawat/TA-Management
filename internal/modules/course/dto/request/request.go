package request

import "time"

type CreateCourse struct {
	CourseName  string    `json:"courseName"`
	CourseID    int       `json:"courseID"`
	CreatedDate time.Time `json:"-"`
	CreatedBy   string    `json:"-"`
}
