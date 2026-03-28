package response

import "time"

type Course struct {
	CourseID      int    `json:"courseID"`
	CourseCode    string `json:"courseCode"`
	CourseName    string `json:"courseName"`
	CourseProgram string `json:"courseProgram"`
	WorkHour      int    `json:"workHour"`
	ClassStart    string `json:"classStart"`
	ClassEnd      string `json:"classEnd"`
	Classday      string `json:"classday"`
	ProfessorName string `json:"professorName"`
	Semester      string `json:"semester"`
	Status        string `json:"status"`
	Section       string `json:"section"`
	SemesterStart string `json:"semesterStart"`
	SemesterEnd   string `json:"semesterEnd"`
	DiscordRoleID string `json:"discordRoleID"`
}

type JobPost struct {
	CourseCode         string `json:"courseCode"`
	CourseName         string `json:"courseName"`
	CourseProgram      string `json:"courseProgram"`
	TaAllocation       int    `json:"taAllocation"`
	RemainingPositions int    `json:"remainingPositions"`
	WorkHour           int    `json:"workHour"`
	ClassStart         string `json:"classStart"`
	ClassEnd           string `json:"classEnd"`
	Location           string `json:"location"`
	Grade              string `json:"grade"`
	Task               string `json:"task"`
	Classday           string `json:"classday"`
	ProfessorName      string `json:"professorName"`
	Semester           string `json:"semester"`
	Status             string `json:"status"`
	StatusID           int    `json:"statusID"`
	JobPostID          int    `json:"jobPostID"`
	CourseID           int    `json:"courseID"`
	Section            string `json:"sec"`
}

type Application struct {
	ApplicationId  int       `json:"applicationId"`
	CourseID       string    `json:"courseID"`
	CourseCode     string    `json:"courseCode"`
	CourseName     string    `json:"courseName"`
	StudentID      int       `json:"studentID"`
	StatusID       int       `json:"statusID"`
	StatusCode     string    `json:"statusCode"`
	Grade          string    `json:"grade"`
	StudentName    string    `json:"studentName"`
	StudentNameTH  string    `json:"studentNameTH"`
	PhoneNumber    string    `json:"phoneNumber"`
	Classday       string    `json:"classDay"`
	ProfessorName  string    `json:"professorName"`
	ClassStart     string    `json:"classStart"`
	ClassEnd       string    `json:"classEnd"`
	CreatedDate    time.Time `json:"createdDate"`
	HasTranscript  bool      `json:"hasTranscript"`
	HasBankAccount bool      `json:"hasBankAccount"`
	HasStudentCard bool      `json:"hasStudentCard"`
	Location       string    `json:"location"`
	RejectReason   string    `json:"rejectReason"`
	JobPostID      int       `json:"jobPostID"`
}

type PdfFile struct {
	FileName  string
	FileBytes []byte
}
