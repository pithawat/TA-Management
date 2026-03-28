package service

import (
	"TA-management/internal/constants"
	generalresponse "TA-management/internal/modules/shared/dto/response"
	"TA-management/internal/modules/ta_duty/dto/request"
	"TA-management/internal/modules/ta_duty/dto/response"
	"TA-management/internal/modules/ta_duty/repository"
	"TA-management/internal/utils"
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"html/template"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

//go:embed templates/signature.html
var htmlTemplateString string

type TaDutyServiceImplementation struct {
	repo   repository.TaDutyRepository
	logger *zap.SugaredLogger
}

func NewTaDutyServiceImplementation(repo repository.TaDutyRepository, logger *zap.SugaredLogger) TaDutyServiceImplementation {
	return TaDutyServiceImplementation{
		repo:   repo,
		logger: logger}
}

func (s TaDutyServiceImplementation) GetTADutyRoadmap(courseID int, studentID int) (*[]response.DutyChecklistItem, error) {
	result, err := s.repo.GetTADutyRoadmap(courseID, studentID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s TaDutyServiceImplementation) MarkDutyAsDone(courseID int, studentID int, dutyDate string) (*generalresponse.GeneralResponse, error) {
	// 1. Load Thailand Timezone
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		s.logger.Errorf("Failed to load timezone: %v", err)
		// Fallback to UTC if timezone loading fails on server
		loc = time.UTC
	}

	// 2. Parse the input date IN THE CONTEXT of Thailand
	// This creates: YYYY-MM-DD 00:00:00 +0700 ICT
	parsedDate, err := time.ParseInLocation("2006-01-02", dutyDate, loc)
	if err != nil {
		return nil, errors.New("invalid date format")
	}

	// 3. Get "Now" in Thailand and strip the time (get 00:00:00)
	now := time.Now().In(loc)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	if parsedDate.After(today) {
		fmt.Println("Date is in the future")
		return nil, errors.New("cannot check off future duties")
	}

	err = s.repo.MarkDutyAsDone(courseID, studentID, dutyDate)
	if err != nil {
		s.logger.Errorf("Failed to save duty to DB: %v", err)
		return nil, err
	}

	return &generalresponse.GeneralResponse{
		Message: "MarkDuty as Done successfully",
	}, nil
}

func (s TaDutyServiceImplementation) ExportPaymentReport(rq request.ExportPaymentReportRequest) (*bytes.Buffer, *request.CourseDutyData, error) {
	TAdutyData, courseData, err := s.repo.GetTADutyDataExportPayment(rq.CourseID, rq.Month)
	if err != nil {
		s.logger.Errorf("Failed to Get Ta Duty Data :%v", err)
		return nil, nil, err
	}

	courseData.MonthName = utils.GetThaiMonthName(rq.Month)
	courseData.Year = fmt.Sprintf("%d", rq.Year+543)

	fileBytes, err := s.GeneratePaymentExcel(*TAdutyData, *courseData, rq.HourlyRate)
	if err != nil {
		s.logger.Errorf("Failed on GenearatePaymentExcel")
		return nil, nil, err
	}

	return fileBytes, courseData, nil
}

func (s TaDutyServiceImplementation) GeneratePaymentExcel(students []request.CreatePaymentData, courseData request.CourseDutyData, hourlyRate int) (*bytes.Buffer, error) {
	f, err := excelize.OpenFile("./prototype/payment-template.xlsx")
	if err != nil {
		s.logger.Errorf("cannot open file : %v", err)
		return nil, err
	}
	defer f.Close()

	sheetName := "Sheet1"

	f.SetCellValue(sheetName, "A1", constants.PaymentReportTitle+fmt.Sprintf(" %s %s กลุ่ม(%s) ประจำภาคเรียนที่ %s", courseData.CourseCode, courseData.CourseName, courseData.Sec, courseData.Semester))
	f.SetCellValue(sheetName, "A2", constants.PaymentReportSubTitle+fmt.Sprintf("%s %s", courseData.MonthName, courseData.Year))

	//Fill duty date and Timerange
	if len(students) > 0 {
		for i, duty := range students[0].Duty {

			col, _ := excelize.ColumnNumberToName(4 + i)
			if i >= 32 {
				break
			}

			f.SetCellValue(sheetName, fmt.Sprintf("%s5", col), duty.TimeRange)
			f.SetCellValue(sheetName, fmt.Sprintf("%s6", col), duty.Date)
		}
	}

	//Fill student data
	grandTotal := 0
	for i, student := range students {
		rowIdx := 7 + i
		if rowIdx > 20 {
			break
		}

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowIdx), i+1)

		//Fill student name
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowIdx), student.StudentName)

		//Fill workHour Rate
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowIdx), hourlyRate)

		totalHours := 0

		// work hour check
		for dutyIdx, duty := range student.Duty {
			colName, _ := excelize.ColumnNumberToName(4 + dutyIdx)
			if dutyIdx >= 32 {
				break
			}

			if duty.IsChecked {
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", colName, rowIdx), student.WorkHour)
				totalHours += student.WorkHour
			} else {
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", colName, rowIdx), "-")
			}

		}

		// Fill total hour and total payment on that student
		f.SetCellValue(sheetName, fmt.Sprintf("AL%d", rowIdx), totalHours)
		f.SetCellValue(sheetName, fmt.Sprintf("AM%d", rowIdx), totalHours*hourlyRate)
		grandTotal += totalHours * hourlyRate
	}
	toatalThaiText := utils.ThaiBahtText(grandTotal)
	f.SetCellValue(sheetName, "A22", fmt.Sprintf("%s %s", constants.PaymentReportTotalThaiText, toatalThaiText))

	var buf bytes.Buffer
	if _, err := f.WriteTo(&buf); err != nil {
		return nil, err
	}
	return &buf, nil

}

func (s TaDutyServiceImplementation) ExportSignatureSheet(rq request.ExportSignatureSheet) (*bytes.Buffer, *request.CourseDutyData, error) {
	TADutydata, courseData, err := s.repo.GetTADutyDataExportSignature(rq.CourseID, rq.Month)
	if err != nil {
		s.logger.Errorf("Failed on get duty data: %v", err)
		return nil, nil, err
	}

	courseData.MonthName = utils.GetThaiMonthName(rq.Month)
	courseData.Year = fmt.Sprintf("%d", rq.Year+543)

	filesBytes, err := s.GenerateSignatureSheet(*TADutydata, *courseData)
	if err != nil {
		s.logger.Errorf("Failed on generate signature sheet %v", err)
		return nil, nil, err
	}

	return filesBytes, courseData, nil
}

func (s TaDutyServiceImplementation) GenerateSignatureSheet(rq request.CreateSignatureSheet, courseData request.CourseDutyData) (*bytes.Buffer, error) {
	f, err := excelize.OpenFile("./prototype/signature-template.xlsx")
	if err != nil {
		s.logger.Errorf("cannot openfile :%v", err)
		return nil, err
	}
	defer f.Close()

	sheetName := "Sheet1"
	f.SetCellValue(sheetName, "A2", fmt.Sprintf("วิชาปฎิบัติการ %s %s (กลุ่ม %s) ประจำภาคเรียนที่ %s", courseData.CourseCode, courseData.CourseName, courseData.Sec, courseData.Semester))
	f.SetCellValue(sheetName, "A4", fmt.Sprintf("ประจำเดือน %s %s", courseData.MonthName, courseData.Year))

	for i, date := range rq.DutyDate {
		rowIdx := 7 + (i * 6)

		//set index
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowIdx), i+1)
		//set date
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowIdx), date)

		//set TA name
		for j, name := range rq.TAName {
			rowIdxx := rowIdx + j + 1
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowIdxx), name)
		}
	}

	var buf bytes.Buffer
	if _, err := f.WriteTo(&buf); err != nil {
		s.logger.Errorf("Failed on write to buffer: %v", err)
		return nil, err
	}

	return &buf, nil
}

func (s TaDutyServiceImplementation) GenerateSignatureSheetPDF(rq request.CreateSignatureSheet, courseData request.CourseDutyData) ([]byte, error) {
	// 1. Setup the data (same as before)
	data := request.SignatureSheetData{
		CourseCode: courseData.CourseCode,
		CourseName: courseData.CourseName,
		Sec:        courseData.Sec,
		Semester:   courseData.Semester,
		MonthName:  courseData.MonthName,
		Year:       courseData.Year,
	}

	for i, date := range rq.DutyDate {
		data.Duties = append(data.Duties, request.DutyGroup{
			Index:   i + 1,
			Date:    date,
			TANames: rq.TAName,
		})
	}

	// 2. Render HTML using the template
	tmpl, err := template.New("signature").Parse(htmlTemplateString)
	if err != nil {
		return nil, err
	}
	var htmlBuf bytes.Buffer
	tmpl.Execute(&htmlBuf, data)

	// 3. Setup Chromedp Context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var pdfBuffer []byte
	err = chromedp.Run(ctx,
		// This loads the HTML string directly into the headless browser
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			return page.SetDocumentContent(frameTree.Frame.ID, htmlBuf.String()).Do(ctx)
		}),
		// This triggers the PDF generation
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().
				WithPrintBackground(true).
				WithPaperWidth(8.27).   // A4 Width in inches
				WithPaperHeight(11.69). // A4 Height in inches
				Do(ctx)
			pdfBuffer = buf
			return err
		}),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %v", err)
	}

	return pdfBuffer, nil
}
