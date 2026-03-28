package response

import (
	"time"
)

type LookupResponse struct {
	Id    int    `json:"id"`
	Value string `json:"value"`
}

type TaDetail struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type AvailableMonth struct {
	MonthID   int    `json:"monthID"`
	MonthName string `json:"monthName"`
	Year      int    `json:"year"`
}

type PdfFile struct {
	FileName  string
	FileBytes []byte
}

type SemesterResponse struct {
	Id        int       `json:"id"`
	Semester  string    `json:"semester"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	IsActive  bool      `json:"isActive"`
}
