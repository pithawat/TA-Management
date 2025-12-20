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
	ClassStart      time.Time `json:"classStart"`
	ClassEnd        time.Time `json:"classEnd"`
	TaAllocation    int       `json:"taAllocation"`
	GradeID         int       `json:"gradeID"`
	Task            string    `json:"task"`
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

type ApplyCourse struct {
	StudentID int `form:"studentID"`
	StatusID  int `form:"statusID"`
	CourseID  *int
	FileBytes *[]byte
	FileName  *string
}
