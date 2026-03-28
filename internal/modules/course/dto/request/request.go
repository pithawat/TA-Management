package request

import "time"

type RejectApplication struct {
	ApplicationId int    `json:"applicationId"`
	RejectReason  string `json:"rejectReason"`
}

type CreateCourse struct {
	CourseName      string    `json:"courseName"`
	CourseCode      string    `json:"courseCode"`
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
	WorkHour        int       `json:"workHour"`
	CreatedDate     time.Time `json:"-"`
}

type CreateJobPost struct {
	CourseID     int       `json:"courseID"`
	ProfessorID  int       `json:"professorID"`
	Location     string    `json:"location"`
	TaAllocation int       `json:"taAllocation"`
	GradeID      int       `json:"gradeID"`
	Task         string    `json:"task"`
	CreatedDate  time.Time `json:"-"`
}

type UpdateJobPost struct {
	Id           int     `json:"id"`
	ProfessorID  *int    `json:"professorID"`
	Location     *string `json:"location"`
	TaAllocation *int    `json:"taAllocation"`
	GradeID      *int    `json:"gradeID"`
	Task         *string `json:"task"`
}

type UpdateCourse struct {
	CourseName      *string    `json:"courseName"`
	CourseCode      *string    `json:"courseCode"`
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
	StudentID        int    `form:"studentID"`
	FirstName        string `form:"firstName"`
	LastName         string `form:"lastName"`
	Grade            string `form:"grade"`
	Purpose          string `form:"purpose"`
	Experience       string `form:"experience"`
	AttachNewPDF     bool   `form:"attachNewPDF"`
	JobPostID        *int
	TranscriptBytes  *[]byte
	TranscriptName   *string
	BankAccountBytes *[]byte
	BankAccountName  *string
	StudentCardBytes *[]byte
	StudentCardName  *string
	PhoneNumber      *string `form:"phoneNumber"`
	FirstnameThai    *string `form:"firstname_thai"`
	LastnameThai     *string `form:"lastname_thai"`
}

type UpdateCourseDiscord struct {
	RoleID      string `json:"role_id"`
	ChannelID   string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
}
