package request

import "time"

// type CreateCourse struct {
// 	CourseName      string    `json:"courseName"`
// 	CourseID        string    `json:"courseID"`
// 	ProfessorID     int       `json:"professorID"`
// 	CourseProgramID int       `json:"courseProgramID"`
// 	CourseProgram   string    `json:"courseProgram"`
// 	Sec             string    `json:"sec"`
// 	SemesterID      int       `json:"semesterID"`
// 	Semester        string    `json:"semester"`
// 	ClassdayID      int       `json:"classdayID"`
// 	Classday        string    `json:"classday"`
// 	ClassStart      time.Time `json:"classStart"`
// 	ClassEnd        time.Time `json:"classEnd"`
// 	CreatedDate     time.Time `json:"-"`
// }

type CreateCourse struct {
	CourseName      string    `json:"courseName"`
	CourseID        string    `json:"courseID"`
	ProfessorID     int       `json:"professorID"`
	CourseProgramID int       `json:"courseProgramID"`
	CourseProgram   string    `json:"courseProgram"`
	Sec             string    `json:"sec"`
	SemesterID      int       `json:"semesterID"`
	Semester        string    `json:"semester"`
	ClassdayID      int       `json:"classdayID"`
	Classday        string    `json:"classday"`
	ClassStart      string    `json:"classStart"`
	ClassEnd        string    `json:"classEnd"`
	TaAllocation    int       `json:"taAllocation"`
	GradeID         int       `json:"gradeID"`
	Task            string    `json:"task"`
	WorkHour        int       `json:"workHour"`
	Location        string    `json:"location"`
	CreatedDate     time.Time `json:"-"`
}

type UpdateCourse struct {
	CourseName      *string    `json:"courseName"`
	CourseID        *string    `json:"courseID"`
	ProfessorID     *int       `json:"professorID"`
	CourseProgramID *int       `json:"courseProgramID"`
	CourseProgram   *string    `json:"courseProgram"`
	Sec             *string    `json:"sec"`
	SemesterID      *int       `json:"semesterID"`
	Semester        *string    `json:"semester"`
	ClassdayID      *int       `json:"classdayID"`
	Classday        *string    `json:"classday"`
	ClassStart      *time.Time `json:"classStart"`
	ClassEnd        *time.Time `json:"classEnd"`
	CreatedDate     *time.Time `json:"-"`
	Id              int        `json:"id"`
}

type ApplyJobPost struct {
	StudentID int    `form:"studentID"`
	Grade     string `form:"grade"`
	Purpose   string `form:"purpose"`
	JobPostID *int
	FileBytes *[]byte
	FileName  *string
}
