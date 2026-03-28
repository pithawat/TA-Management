package request

import "time"

type CreateSemester struct {
	Semester  string    `json:"semester"`
	Year      string    `json:"year"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}

type UpdateSemester struct {
	ID        int        `json:"id"`
	Semester  *string    `json:"semester"`
	Year      *string    `json:"year"`
	StartDate *time.Time `json:"startDate"`
	EndDate   *time.Time `json:"endDate"`
}

type CreateHoliday struct {
	Date     time.Time `json:"date"`
	NameThai string    `json:"nameThai"`
	NameEng  string    `json:"nameEng"`
	Type     string    `json:"type"`
}
