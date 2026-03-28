package request

type EmailRequest struct {
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}

type MailForAllCourse struct {
	CourseID []int  `json:"courseID"`
	Subject  string `json:"subject"`
	Body     string `json:"body"`
}

type MailForCourse struct {
	CourseID int    `json:"courseID"`
	Subject  string `json:"subject"`
	Body     string `json:"body"`
}

type MailForTA struct {
	StudentID int    `json:"studentID"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
}

type CreateEmailHistory struct {
	Subject      string
	Body         string
	ReceivedName string
	NReceived    int
	StatusID     int
}

type CreateDiscordChannel struct {
	CourseID   int    `json:"courseID"`
	CourseCode string `json:"courseCode"`
	CourseName string `json:"courseName"`
	Semester   string `json:"semester"`
	Sec        string `json:"sec"`
	GuildID    string
}
