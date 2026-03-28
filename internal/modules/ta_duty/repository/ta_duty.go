package repository

import (
	"TA-management/internal/modules/ta_duty/dto/request"
	"TA-management/internal/modules/ta_duty/dto/response"
)

type TaDutyRepository interface {
	GetTADutyRoadmap(courseID int, studentID int) (*[]response.DutyChecklistItem, error)
	MarkDutyAsDone(courseID int, studentID int, dutyDate string) error
	GetTADutyDataExportPayment(courseID int, month int) (*[]request.CreatePaymentData, *request.CourseDutyData, error)
	GetTADutyDataExportSignature(courseID int, month int) (*request.CreateSignatureSheet, *request.CourseDutyData, error)
}
