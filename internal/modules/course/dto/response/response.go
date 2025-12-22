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
	ApplicationId  int       `json:"applicationId"`
	CourseID       string    `json:"courseID"`
	CourseName     string    `json:"courseName"`
	StudentID      int       `json:"studentID"`
	StatusID       int       `json:"statusID"`
	StatusCode     string    `json:"statusCode"`
	Grade          string    `json:"grade"`
	StudentName    string    `json:"studentName"`
	PhoneNumber    string    `json:"phoneNumber"`
	Classday       string    `json:"classDay"`
	ClassStart     string    `json:"classStart"`
	ClassEnd       string    `json:"classEnd"`
	CreatedDate    time.Time `json:"createdDate"`
	HasTranscript  bool      `json:"hasTranscript"`
	HasBankAccount bool      `json:"hasBankAccount"`
	HasStudentCard bool      `json:"hasStudentCard"`
}

type PdfFile struct {
	FileName  string
	FileBytes []byte
}
