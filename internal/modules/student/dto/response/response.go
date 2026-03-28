package response

type StudentProfile struct {
	StudentID           int    `json:"studentId" db:"student_id"`
	FirstnameThai       string `json:"firstnameThai" db:"firstname_thai"`
	LastnameThai        string `json:"lastnameThai" db:"lastname_thai"`
	Email               string `json:"email" db:"email"`
	PhoneNumber         string `json:"phoneNumber" db:"phone_number"`
	HasTranscript       bool   `json:"hasTranscript" db:"-"`
	TranscriptFileName  string `json:"transcriptFileName" db:"-"`
	HasBankAccount      bool   `json:"hasBankAccount" db:"-"`
	BankAccountFileName string `json:"bankAccountFileName" db:"-"`
	HasStudentCard      bool   `json:"hasStudentCard" db:"-"`
	StudentCardFileName string `json:"studentCardFileName" db:"-"`
}
