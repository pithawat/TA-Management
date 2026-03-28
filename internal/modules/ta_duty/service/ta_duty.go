package service

import (
	generalresponse "TA-management/internal/modules/shared/dto/response"
	"TA-management/internal/modules/ta_duty/dto/request"
	"TA-management/internal/modules/ta_duty/dto/response"
	"bytes"
)

type TaDutyService interface {
	GetTADutyRoadmap(courseID int, studentID int) (*[]response.DutyChecklistItem, error)
	MarkDutyAsDone(courseID int, studentID int, dutyDate string) (*generalresponse.GeneralResponse, error)
	ExportPaymentReport(rq request.ExportPaymentReportRequest) (*bytes.Buffer, *request.CourseDutyData, error)
	ExportSignatureSheet(rq request.ExportSignatureSheet) (*bytes.Buffer, *request.CourseDutyData, error)
}
