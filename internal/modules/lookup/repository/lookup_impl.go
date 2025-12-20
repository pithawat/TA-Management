package repository

import (
	"TA-management/internal/modules/lookup/dto/response"
	"database/sql"
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
				course_program_value 
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
				class_day_value 
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

func (r LookupRepositoryImplementation) GetSemester() (*[]response.LookupResponse, error) {

	query := `SELECT 
				semester_id, 
				semester_value 
			FROM semester
			ORDER BY start_date ASC`

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
