package response

import "time"

type Course struct {
	CourseID   string `json:"courseID"`
	CourseName string `json:"courseName"`
}

type Application struct {
	CourseID    string    `json:"courseID"`
	CourseName  string    `json:"courseName"`
	StudentID   int       `json:"studentID"`
	StatusID    int       `json:"statusID"`
	StatusCode  string    `json:"statusCode"`
	CreatedDate time.Time `json:"createdDate"`
}

type ApplicationTrancript struct {
	FileName   string
	Transcript []byte
}
