package response

import "time"

// SemesterHistory holds summary data for a past (expired) semester.
type SemesterHistory struct {
	SemesterID    int       `json:"semesterID"`
	SemesterValue string    `json:"semesterValue"`
	StartDate     time.Time `json:"startDate"`
	EndDate       time.Time `json:"endDate"`
	CourseCount   int       `json:"courseCount"`
}
