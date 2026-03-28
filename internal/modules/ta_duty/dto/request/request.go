package request

type CreatePaymentData struct {
	StudentName string
	WorkHour    int
	Duty        []DutyChecklistItem
}

type ExportPaymentReportRequest struct {
	CourseID   int `json:"courseID"`
	HourlyRate int `json:"hourlyRate"`
	Month      int `json:"month"`
	Year       int `json:"year"`
}

type ExportSignatureSheet struct {
	CourseID int `json:"courseID"`
	Month    int `json:"month"`
	Year     int `json:"year"`
}

type CourseDutyData struct {
	CourseCode string
	CourseName string
	Semester   string
	Sec        string
	MonthName  string
	Year       string
}

type CreateSignatureSheet struct {
	DutyDate []string
	TAName   []string
}

type DutyChecklistItem struct {
	Date      string `json:"date"`
	TimeRange string `json:"timeRange"`
	Status    string `json:"status"`
	IsChecked bool   `json:"isChecked"`
}

type SignatureSheetData struct {
	CourseCode string
	CourseName string
	Sec        string
	Semester   string
	MonthName  string
	Year       string
	Duties     []DutyGroup
}

type DutyGroup struct {
	Index   int
	Date    string
	TANames []string
}
