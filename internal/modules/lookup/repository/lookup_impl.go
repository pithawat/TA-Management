package repository

import (
	"TA-management/internal/modules/lookup/dto/request"
	"TA-management/internal/modules/lookup/dto/response"
	"TA-management/internal/modules/ta_duty/entity"
	"database/sql"
	"fmt"
	"strings"
)

type LookupRepositoryImplementation struct {
	db *sql.DB
}

func NewLookupRepository(DB *sql.DB) LookupRepositoryImplementation {
	return LookupRepositoryImplementation{db: DB}
}

func (r LookupRepositoryImplementation) GetCourseProgram() (*[]response.LookupResponse, error) {

	query := `SELECT 
				course_program_id, 
				course_program_value_thai 
			FROM course_programs`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	var coursePrograms []response.LookupResponse
	for rows.Next() {
		var courseProgram response.LookupResponse
		err := rows.Scan(&courseProgram.Id, &courseProgram.Value)
		if err != nil {
			return nil, err
		}
		coursePrograms = append(coursePrograms, courseProgram)
	}
	return &coursePrograms, nil
}

func (r LookupRepositoryImplementation) GetClassday() (*[]response.LookupResponse, error) {

	query := `SELECT 
				class_day_id, 
				class_day_value_thai 
			FROM class_days`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	var classDays []response.LookupResponse
	for rows.Next() {
		var classDay response.LookupResponse
		err := rows.Scan(&classDay.Id, &classDay.Value)
		if err != nil {
			return nil, err
		}
		classDays = append(classDays, classDay)
	}
	return &classDays, nil
}

func (r LookupRepositoryImplementation) GetSemester() (*[]response.SemesterResponse, error) {

	query := `SELECT 
				semester_id, 
				semester_value,
                start_date,
                end_date,
				is_active
			FROM semester
            WHERE end_date >= CURRENT_DATE OR is_active = TRUE
			ORDER BY start_date ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	var semesters []response.SemesterResponse
	for rows.Next() {
		var semester response.SemesterResponse
		err := rows.Scan(&semester.Id, &semester.Semester, &semester.StartDate, &semester.EndDate, &semester.IsActive)
		if err != nil {
			return nil, err
		}
		semesters = append(semesters, semester)
	}
	return &semesters, nil
}

func (r LookupRepositoryImplementation) GetSemesterDropdown() (*[]response.LookupResponse, error) {

	query := `SELECT 
				semester_id, 
				semester_value
			FROM semester
			ORDER BY start_date DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	var semesters []response.LookupResponse
	for rows.Next() {
		var semester response.LookupResponse
		err := rows.Scan(&semester.Id, &semester.Value)
		if err != nil {
			return nil, err
		}
		semesters = append(semesters, semester)
	}
	return &semesters, nil
}

func (r LookupRepositoryImplementation) GetGrade() (*[]response.LookupResponse, error) {

	query := `SELECT 
				grade_id, 
				grade_value 
			FROM grades`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	var grades []response.LookupResponse
	for rows.Next() {
		var grade response.LookupResponse
		err := rows.Scan(&grade.Id, &grade.Value)
		if err != nil {
			return nil, err
		}
		grades = append(grades, grade)
	}
	return &grades, nil
}

func (r LookupRepositoryImplementation) GetProfessors() (*[]response.LookupResponse, error) {
	query := `SELECT 
				professor_ID, 
				CONCAT(prefix, ' ', firstname_thai, ' ', lastname_thai) as fullname 
			FROM professors`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	var professors []response.LookupResponse
	for rows.Next() {
		var professor response.LookupResponse
		err := rows.Scan(&professor.Id, &professor.Value)
		if err != nil {
			return nil, err
		}
		professors = append(professors, professor)
	}
	return &professors, nil
}

func (r LookupRepositoryImplementation) SyncOfficialHoliday(holidays []request.CreateHoliday) error {
	if len(holidays) == 0 {
		return nil
	}

	values := []interface{}{}
	query := `insert into holidays (holiday_date, name_eng, name_thai, category) VALUES`

	for i, h := range holidays {
		//calculate placeholder
		p1 := i*4 + 1
		p2 := i*4 + 2
		p3 := i*4 + 3
		p4 := i*4 + 4

		query += fmt.Sprintf("($%d, $%d, $%d, $%d)", p1, p2, p3, p4)
		if i < len(holidays)-1 {
			query += ","
		}
		values = append(values, h.Date, h.NameEng, h.NameThai, h.Type)
	}

	query += ` ON CONFLICT (holiday_date) DO UPDATE SET 
				name_eng = EXCLUDED.name_eng,
				name_thai = EXCLUDED.name_thai,
				category = EXCLUDED.category`

	_, err := r.db.Exec(query, values...)
	if err != nil {
		fmt.Printf("Failed to insert holidays: %v\n", err)
		return err
	}

	return nil
}

func (r LookupRepositoryImplementation) GetHolidaysByMonth(month int, year int) ([]entity.Holiday, error) {
	query := `SELECT id, holiday_date, name_thai, category FROM holidays 
	          WHERE EXTRACT(MONTH FROM holiday_date) = $1 AND EXTRACT(YEAR FROM holiday_date) = $2`
	rows, err := r.db.Query(query, month, year)
	if err != nil {
		return nil, fmt.Errorf("failed to get holidays: %v", err)
	}
	defer rows.Close()

	var holidays []entity.Holiday
	for rows.Next() {
		var h entity.Holiday
		if err := rows.Scan(&h.ID, &h.Date, &h.Name, &h.Type); err != nil {
			return nil, fmt.Errorf("failed to scan holiday: %v", err)
		}
		holidays = append(holidays, h)
	}
	return holidays, nil
}

func (r LookupRepositoryImplementation) AddSpecialHoliday(holiday request.CreateHoliday) error {
	query := `INSERT INTO holidays (holiday_date, name_thai, category) VALUES ($1, $2, 'special')`
	_, err := r.db.Exec(query, holiday.Date, holiday.NameThai)
	if err != nil {
		return fmt.Errorf("failed to add special holiday: %v", err)
	}
	return nil
}

func (r LookupRepositoryImplementation) DeleteHoliday(id int) error {
	query := `DELETE FROM holidays WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete holiday: %v", err)
	}
	return nil
}

func (r LookupRepositoryImplementation) GetTA(searchVal string) (*[]response.TaDetail, error) {

	searchVal = "%" + searchVal + "%"
	var TAs []response.TaDetail

	// query := `SELECT
	// 			student_ID,
	// 			firstname || ' ' || lastname AS name
	// 			FROM students
	// 			WHERE CAST(student_ID as TEXT) LIKE $1
	// 			OR (firstname || ' ' || lastname) LIKE $2
	// 		`

	//search first then filter is robust than add where clause in query
	query := `
				WITH searchable_students AS(
					SELECT
						s.student_ID,
						s.firstname_thai || ' ' || s.lastname_thai AS name
					FROM ta_courses tc
					LEFT JOIN students s ON s.student_id = tc.student_id
				)
				SELECT * FROM searchable_students
				WHERE student_ID::TEXT ILIKE $1
				OR name ILIKE $1  `
	rows, err := r.db.Query(query, searchVal)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var TA response.TaDetail
		if err := rows.Scan(&TA.Id, &TA.Name); err != nil {
			return nil, err
		}
		TAs = append(TAs, TA)
	}

	return &TAs, nil
}

func (r LookupRepositoryImplementation) GetAvailableMonths(courseId int) (*[]response.AvailableMonth, error) {

	query := `
			SELECT DISTINCT
				EXTRACT(MONTH FROM series_date) AS month_id,
				TO_CHAR(series_date, 'Month') AS month_name,
				EXTRACT(YEAR FROM series_date) AS year_val
			FROM(
				SELECT generate_series(start_date, end_date, '1 month'::interval)::date AS series_date
				FROM semester s
				JOIN courses c ON c.semester_ID=s.semester_ID
				WHERE c.course_ID = $1
			)sub
			ORDER BY year_val, month_id;
		`
	rows, err := r.db.Query(query, courseId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var months []response.AvailableMonth
	for rows.Next() {
		var month response.AvailableMonth
		err := rows.Scan(&month.MonthID, &month.MonthName, &month.Year)
		if err != nil {
			return nil, err
		}
		months = append(months, month)
	}

	return &months, nil
}

func (r LookupRepositoryImplementation) GetStudentCard(studentID int) (*response.PdfFile, error) {

	query := `SELECT file_name, file_bytes FROM student_card_storage WHERE student_ID=$1`

	var result response.PdfFile
	err := r.db.QueryRow(query, studentID).Scan(&result.FileName, &result.FileBytes)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}

func (r LookupRepositoryImplementation) GetTranscript(studentID int) (*response.PdfFile, error) {

	query := `SELECT file_name, file_bytes FROM transcript_storage WHERE student_ID=$1`

	var result response.PdfFile
	err := r.db.QueryRow(query, studentID).Scan(&result.FileName, &result.FileBytes)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}

func (r LookupRepositoryImplementation) GetBankAccount(studentID int) (*response.PdfFile, error) {

	query := `SELECT file_name, file_bytes FROM bank_account_storage WHERE student_ID=$1`

	var result response.PdfFile
	err := r.db.QueryRow(query, studentID).Scan(&result.FileName, &result.FileBytes)
	if err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &result, nil
}

func (r LookupRepositoryImplementation) AddSemester(rq request.CreateSemester) error {

	query := `INSERT INTO semester (
				semester_value,
				start_date,
				end_date) VALUES($1, $2, $3)
	`
	_, err := r.db.Exec(query, fmt.Sprintf("%s/%s", rq.Semester, rq.Year), rq.StartDate, rq.EndDate)
	if err != nil {
		return err
	}

	return nil
}

func (r LookupRepositoryImplementation) UpdateSemester(rq request.UpdateSemester) error {

	query := `UPDATE semester SET `
	params := []interface{}{}
	placeholderID := 1

	addFeild := func(columnName string, value interface{}) {
		query += fmt.Sprintf("%s = $%d, ", columnName, placeholderID)
		params = append(params, value)
		placeholderID++
	}

	if rq.Semester != nil {
		addFeild("semester_value", rq.Semester)
	}

	if rq.StartDate != nil {
		addFeild("start_date", rq.StartDate)
	}

	if rq.EndDate != nil {
		addFeild("end_date", rq.EndDate)
	}

	if len(params) == 0 {
		return nil // Or return an error "nothing to update"
	}

	query = strings.TrimSuffix(query, ", ")
	query += fmt.Sprintf(" WHERE semester_ID = $%d", placeholderID)

	fmt.Println(query)
	params = append(params, rq.ID)

	_, err := r.db.Exec(query, params...)
	if err != nil {
		return err
	}

	return nil
}

func (r LookupRepositoryImplementation) SetSemesterActive(semesterID int) error {

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	query := `UPDATE semester SET is_active=FALSE`
	_, err = tx.Exec(query)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to set all semester to default")
	}

	query = `UPDATE semester SET is_active=TRUE WHERE semester_ID=$1`
	_, err = tx.Exec(query, semesterID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to set active semester")
	}

	return tx.Commit()

}
