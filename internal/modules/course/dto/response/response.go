package response

import "time"

type Course struct {
	CourseID      string `json:"courseID"`
	CourseName    string `json:"courseName"`
	CourseProgram string `json:"courseProgram"`
	TaAllocation  int    `json:"taAllocation"`
	WorkHour      int    `json:"workHour"`
	ClassStart    string `json:"classStart"`
	ClassEnd      string `json:"classEnd"`
	Location      string `json:"location"`
	Grade         string `json:"grade"`
	Task          string `json:"task"`
	Classday      string `json:"classday"`
	ProfessorName string `json:"professorName"`
	Semester      string `json:"semester"`
	Status        string `json:"status"`
	JobPostID     int    `json:"jobPostID"`
}

type Application struct {
	ApplicationId int       `json:"applicationId"`
	CourseID      string    `json:"courseID"`
	CourseName    string    `json:"courseName"`
	StudentID     int       `json:"studentID"`
	StatusID      int       `json:"statusID"`
	StatusCode    string    `json:"statusCode"`
	CreatedDate   time.Time `json:"createdDate"`
}

type ApplicationTrancript struct {
	FileName   string
	Transcript []byte
}
