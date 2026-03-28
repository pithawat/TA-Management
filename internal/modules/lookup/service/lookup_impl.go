package service

import (
	"TA-management/internal/modules/lookup/dto/request"
	"TA-management/internal/modules/lookup/dto/response"
	"TA-management/internal/modules/lookup/repository"
	"TA-management/internal/modules/ta_duty/entity"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type LookupServiceImplementation struct {
	repo repository.LookupRepository
}

func NewLookupService(repo repository.LookupRepository) LookupServiceImplementation {
	return LookupServiceImplementation{repo: repo}
}

func (s LookupServiceImplementation) GetCourseProgram() (*[]response.LookupResponse, error) {
	result, err := s.repo.GetCourseProgram()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s LookupServiceImplementation) GetClassday() (*[]response.LookupResponse, error) {
	result, err := s.repo.GetClassday()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s LookupServiceImplementation) GetSemester() (*[]response.SemesterResponse, error) {
	result, err := s.repo.GetSemester()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s LookupServiceImplementation) GetSemesterDropdown() (*[]response.LookupResponse, error) {
	result, err := s.repo.GetSemesterDropdown()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s LookupServiceImplementation) GetGrade() (*[]response.LookupResponse, error) {
	result, err := s.repo.GetGrade()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return result, nil
}

func (s LookupServiceImplementation) GetProfessors() (*[]response.LookupResponse, error) {
	result, err := s.repo.GetProfessors()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return result, nil
}

func (s LookupServiceImplementation) SyncOfficialHoliday(apiKey string, url string) error {
	BOTData, err := s.FetchFromBOT(apiKey, url)
	if err != nil {
		fmt.Printf("Failed to fetch holidays from BOT: %v\n", err)
		return err
	}

	err = s.repo.SyncOfficialHoliday(BOTData)
	if err != nil {
		fmt.Printf("Repo sync failed: %v\n", err)
		return err
	}

	return nil
}

func (s LookupServiceImplementation) FetchFromBOT(apiKey string, url string) ([]request.CreateHoliday, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", apiKey)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call http request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api return non-OK status: %v", resp.Status)
	}

	var botResponse struct {
		Result struct {
			Data []struct {
				HolidayDescriptionThai string `json:"HolidayDescriptionThai"`
				HolidayDescription     string `json:"HolidayDescription"`
				Date                   string `json:"Date"`
			} `json:"data"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&botResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	var holidays []request.CreateHoliday
	for _, item := range botResponse.Result.Data {
		date, err := time.Parse("2006-01-02", item.Date)
		if err != nil {
			fmt.Printf("Failed to parse date %s: %v\n", item.Date, err)
			continue // Skip invalid dates but continue processing others
		}

		holidays = append(holidays, request.CreateHoliday{
			Date:     date,
			NameThai: item.HolidayDescriptionThai,
			NameEng:  item.HolidayDescription,
			Type:     "official",
		})
	}

	return holidays, nil
}

func (s LookupServiceImplementation) FetchFromGoogle(apiKey string) ([]request.CreateHoliday, error) {
	calendarID := "en.th%23holiday@group.v.calendar.google.com"
	url := fmt.Sprintf("https://www.googleapis.com/calendar/v3/calendars/%s/events?key=%s", calendarID, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call http request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api return non-OK status: %v", resp.Status)
	}

	var data struct {
		Items []struct {
			Summary string `json:"summary"`
			Start   struct {
				Date string `json:"date"`
			} `json:"start"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var holidays []request.CreateHoliday
	for _, item := range data.Items {
		date, err := time.Parse("2006-01-02", item.Start.Date)
		if err != nil {
			return nil, fmt.Errorf("failed to parse time: %v", err)
		}

		holidays = append(holidays, request.CreateHoliday{
			Date:    date,
			NameEng: item.Summary,
			Type:    "official",
		})
	}
	return holidays, nil
}

func (s LookupServiceImplementation) GetHolidays(month int, year int) ([]entity.Holiday, error) {
	holidays, err := s.repo.GetHolidaysByMonth(month, year)
	if err != nil {
		fmt.Printf("Failed to get holidays: %v\n", err)
		return nil, err
	}
	return holidays, nil
}

func (s LookupServiceImplementation) AddSpecialHoliday(req request.CreateHoliday) error {
	err := s.repo.AddSpecialHoliday(req)
	if err != nil {
		fmt.Printf("Failed to add special holiday: %v\n", err)
		return err
	}
	return nil
}

func (s LookupServiceImplementation) DeleteHoliday(id int) error {
	err := s.repo.DeleteHoliday(id)
	if err != nil {
		fmt.Printf("Failed to delete holiday: %v\n", err)
		return err
	}
	return nil
}

func (s LookupServiceImplementation) GetTA(searchVal string) (*[]response.TaDetail, error) {
	result, err := s.repo.GetTA(searchVal)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return result, err
}

func (s LookupServiceImplementation) GetAvailableMonths(courseId int) (*[]response.AvailableMonth, error) {
	result, err := s.repo.GetAvailableMonths(courseId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return result, err
}

func (s LookupServiceImplementation) GetTranscript(studentID int) (*response.PdfFile, error) {
	result, err := s.repo.GetTranscript(studentID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return result, nil
}

func (s LookupServiceImplementation) GetBankAccount(studentID int) (*response.PdfFile, error) {
	result, err := s.repo.GetBankAccount(studentID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return result, nil
}

func (s LookupServiceImplementation) GetStudentCard(studentID int) (*response.PdfFile, error) {
	result, err := s.repo.GetStudentCard(studentID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return result, nil
}

func (s LookupServiceImplementation) AddSemester(rq request.CreateSemester) error {

	err := s.repo.AddSemester(rq)
	if err != nil {
		fmt.Printf("failed to insert semester: %v \n", err)
		return err
	}

	return nil
}

func (s LookupServiceImplementation) UpdateSemester(rq request.UpdateSemester) (*[]response.SemesterResponse, error) {

	err := s.repo.UpdateSemester(rq)
	if err != nil {
		fmt.Printf("failed to update semester: %v", err)
		return nil, err
	}

	result, err := s.repo.GetSemester()
	if err != nil {
		fmt.Printf("failed to get semester: %v", err)
		return nil, err
	}

	return result, nil
}

func (s LookupServiceImplementation) SetSemesterActive(semesterID int) error {

	err := s.repo.SetSemesterActive(semesterID)
	if err != nil {
		fmt.Printf("failed to set semester active: %v", err)
	}
	return nil
}
