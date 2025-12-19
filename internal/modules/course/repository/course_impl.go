package repository

import (
	"TA-management/internal/modules/course/dto/request"
	"TA-management/internal/modules/course/dto/response"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

type CourseRepositoryImplementation struct {
	db *sql.DB
}

func NewCourseRepository(DB *sql.DB) CourseRepositoryImplementation {
	return CourseRepositoryImplementation{db: DB}
}

func (r CourseRepositoryImplementation) GetAllCourse() ([]response.Course, error) {

	query := "SELECT course_ID, course_name FROM courses"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	//garantees that connection is released back to the pool ,prevent leak
	defer rows.Close()

	var courses []response.Course
	for rows.Next() {
		var course response.Course

		err := rows.Scan(&course.CourseID, &course.CourseName)
		if err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}

	return courses, nil
}

func (r CourseRepositoryImplementation) CreateCourse(body request.CreateCourse) (int, error) {

	queryCheck := "SELECT COUNT(*) FROM courses WHERE course_ID=$1 AND course_program_ID=$2 AND sec=$3 AND semester_ID=$4 "

	var count int

	row := r.db.QueryRow(queryCheck,
		body.CourseID,
		body.CourseProgramID,
		body.Sec,
		body.SemesterID,
	)

	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to scan duplicate check result:%w", err)
	}

	if count > 0 {
		return 0, errors.New("course already exists")
	}

	query := `INSERT INTO courses(course_ID, course_name,
	professor_ID, course_program_ID, course_program, sec, 
	semester_ID, semester, class_day_ID, class_day, 
	class_start, class_end) 
	values ($1,$2,$3, $4, $5 ,$6 ,$7 ,$8 ,$9 ,$10 ,$11 ,$12)
	RETURNING id`

	var lastInsertId int

	err = r.db.QueryRow(query,
		body.CourseID,
		body.CourseName,
		body.ProfessorID,
		body.CourseProgramID,
		body.CourseProgram,
		body.Sec,
		body.SemesterID,
		body.Semester,
		body.ClassdayID,
		body.Classday,
		body.ClassStart,
		body.ClassEnd,
	).Scan(&lastInsertId)

	if err != nil {
		return 0, err
	}

	return lastInsertId, nil

}

func (r CourseRepositoryImplementation) UpdateCourse(body request.UpdateCourse) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	query := "UPDATE courses SET "
	params := []interface{}{}
	placeholderID := 1 // PostgreSQL needs $1, $2...

	// Helper to add fields and increment the counter
	addField := func(columnName string, value interface{}) {
		query += fmt.Sprintf("%s = $%d, ", columnName, placeholderID)
		params = append(params, value)
		placeholderID++
	}

	// 1. Check each field
	if body.CourseName != nil {
		addField("course_name", body.CourseName)
	}
	if body.CourseID != nil {
		addField("course_id", body.CourseID)
	}
	if body.CourseProgramID != nil {
		addField("course_program_id", body.CourseProgramID)
	}
	if body.CourseProgram != nil {
		addField("course_program", body.CourseProgram)
	}
	if body.Sec != nil {
		addField("sec", body.Sec)
	}
	if body.SemesterID != nil {
		addField("semester_id", body.SemesterID)
	}
	if body.Semester != nil {
		addField("semester", body.Semester)
	}
	if body.ClassdayID != nil {
		addField("class_day_id", body.ClassdayID)
	}
	if body.Classday != nil {
		addField("class_day", body.Classday)
	}
	if body.ClassStart != nil {
		addField("class_start", body.ClassStart)
	}
	if body.ClassEnd != nil {
		addField("class_end", body.ClassEnd)
	}

	// 2. Safety Check: Did the user send ANY data?
	if len(params) == 0 {
		tx.Rollback()
		return nil // Or return an error "nothing to update"
	}

	// 3. Finalize the query string
	query = strings.TrimSuffix(query, ", ")
	query += fmt.Sprintf(" WHERE id = $%d;", placeholderID)
	params = append(params, body.Id)

	// 4. Execute
	_, err = tx.Exec(query, params...)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r CourseRepositoryImplementation) DeleteCourse(id int) error {
	query := "UPDATE courses SET deleted_date = $1 WHERE id = $2"
	_, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return err
	}
	return nil
}

func (r CourseRepositoryImplementation) ApplyCourse(body request.ApplyCourse) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var fileId int
	query := "INSERT INTO transcript_storage(file_bytes,file_name) VALUES($1, $2) RETURNING transcript_ID"
	err = tx.QueryRow(query, body.FileBytes, body.FileName).Scan(&fileId)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed on insert to transcript file : %v", err)
	}
	fmt.Println(body.StudentID)
	fmt.Println(body.StatusID)
	var applicationId int
	query = `INSERT INTO ta_application(transcript_ID, student_ID, 
			status_ID, course_ID, created_date)
			VALUES($1, $2, $3, $4, $5)
			RETURNING id`

	err = tx.QueryRow(query,
		fileId,
		body.StudentID,
		body.StatusID,
		body.CourseID,
		time.Now(),
	).Scan(&applicationId)

	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed on insert to ta_application: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed on commit transaction")
	}

	return applicationId, nil

}

func (r CourseRepositoryImplementation) GetApplicationByStudentId(studentId int) ([]response.Application, error) {
	query := `SELECT 
					ta.student_ID, 
					ta.status_ID, 
					ta.course_ID, 
					ta.created_date,
					st.status_value
				FROM ta_application AS ta 
				LEFT JOIN status AS st
					ON ta.status_ID = st.status_ID
				WHERE student_ID = $1`

	rows, err := r.db.Query(query, studentId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var applications []response.Application
	for rows.Next() {
		var application response.Application

		err := rows.Scan(
			&application.StudentID,
			&application.StatusID,
			&application.CourseID,
			&application.CreatedDate,
			&application.StatusCode,
		)
		if err != nil {
			return nil, err
		}
		applications = append(applications, application)
	}

	return applications, nil
}

func (r CourseRepositoryImplementation) GetApplicationByCourseId(courseId int) ([]response.Application, error) {
	query := `SELECT 
					ta.student_ID, 
					ta.status_ID, 
					ta.course_ID, 
					ta.created_date,
					st.status_value
				FROM ta_application AS ta 
				LEFT JOIN status AS st
					ON ta.status_ID = st.status_ID
				WHERE course_ID = $1`

	rows, err := r.db.Query(query, courseId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var applications []response.Application
	for rows.Next() {
		var application response.Application

		err := rows.Scan(
			&application.StudentID,
			&application.StatusID,
			&application.CourseID,
			&application.CreatedDate,
			&application.StatusCode,
		)
		if err != nil {
			return nil, err
		}
		applications = append(applications, application)
	}

	return applications, nil
}

func (r CourseRepositoryImplementation) GetApplicationDetail(ApplicationId int) (*response.Application, error) {
	query := `SELECT 
					ta.student_ID, 
					ta.status_ID, 
					ta.course_ID, 
					ta.created_date,
					st.status_value
				FROM ta_application AS ta 
				LEFT JOIN status AS st
					ON ta.status_ID = st.status_ID
				WHERE id = $1`

	var application response.Application
	err := r.db.QueryRow(query, ApplicationId).Scan(
		&application.StudentID,
		&application.StatusID,
		&application.CourseID,
		&application.CreatedDate,
		&application.StatusCode,
	)
	if err != nil {
		return nil, err
	}

	return &application, nil
}

func (r CourseRepositoryImplementation) GetApplicationPdf(ApplicationId int) (*response.ApplicationTrancript, error) {
	// 	query := `SELECT
	// 				ta.transcript_ID,
	// 				fs.file_name,
	// 				fs.file_bytes
	// 			FROM ta_application as ta
	// 			LEFT JOIN file_storage as fs
	// 			ON ta.transcript_ID = fs.transcript_ID
	// 			WHERE ta.id = $1`

	// 	var
	// 	err := r.db.QueryRow(query,ApplicationId).Scan()
	return nil, nil
}
