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

	query := `SELECT 
				j.task,
				j.id,
				c.course_ID, 
				c.course_name, 
				c.ta_allocation, 
				c.work_hour,
				c.class_start,
				c.class_end,
				c.location,
				c.course_program,
				cd.class_day_value, 
				p.firstname,
				p.lastname,
				s.semester_value,
				st.status_value,
				g.grade_value
			FROM ta_job_posting AS j
			LEFT JOIN courses AS c
				ON j.course_ID = c.id
			LEFT JOIN class_days AS cd
				ON c.class_day_ID = cd.class_day_ID 
			LEFT JOIN professors AS p
				ON c.professor_ID = p.professor_ID
			LEFT JOIN semester AS s
				ON c.semester_ID = s.semester_ID
			LEFT JOIN status AS st
				ON j.status_ID = st.status_ID
			LEFT JOIN grades AS g
				ON j.grade_ID = g.grade_ID`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	//garantees that connection is released back to the pool ,prevent leak
	defer rows.Close()

	var courses []response.Course
	for rows.Next() {
		var course response.Course
		var firstname string
		var lastname string
		err := rows.Scan(
			&course.Task,
			&course.JobPostID,
			&course.CourseID,
			&course.CourseName,
			&course.TaAllocation,
			&course.WorkHour,
			&course.ClassStart,
			&course.ClassEnd,
			&course.Location,
			&course.CourseProgram,
			&course.Classday,
			&firstname,
			&lastname,
			&course.Semester,
			&course.Status,
			&course.Grade,
		)
		if err != nil {
			return nil, err
		}
		course.ProfessorName = firstname + " " + lastname
		courses = append(courses, course)
	}

	return courses, nil
}

func (r CourseRepositoryImplementation) GetAllCourseByStudentId(studentId int) ([]response.Course, error) {

	query := `SELECT 
				j.task,
				j.id,
				c.course_ID, 
				c.course_name, 
				c.ta_allocation, 
				c.work_hour,
				c.class_start,
				c.class_end,
				c.location,
				c.course_program,
				cd.class_day_value, 
				p.firstname,
				p.lastname,
				s.semester_value,
				st.status_value,
				g.grade_value
			FROM ta_job_posting AS j
			LEFT JOIN courses AS c
				ON j.course_ID = c.id
			LEFT JOIN class_days AS cd
				ON c.class_day_ID = cd.class_day_ID 
			LEFT JOIN professors AS p
				ON c.professor_ID = p.professor_ID
			LEFT JOIN semester AS s
				ON c.semester_ID = s.semester_ID
			LEFT JOIN status AS st
				ON j.status_ID = st.status_ID
			LEFT JOIN grades AS g
				ON j.grade_ID = g.grade_ID
			WHERE NOT EXISTS(
				SELECT 1 
				FROM ta_application as ta
				WHERE ta.job_post_ID = j.id
				AND ta.student_ID = $1)
			`

	rows, err := r.db.Query(query, studentId)
	if err != nil {
		return nil, err
	}
	//garantees that connection is released back to the pool ,prevent leak
	defer rows.Close()

	var courses []response.Course
	for rows.Next() {
		var course response.Course
		var firstname string
		var lastname string
		err := rows.Scan(
			&course.Task,
			&course.JobPostID,
			&course.CourseID,
			&course.CourseName,
			&course.TaAllocation,
			&course.WorkHour,
			&course.ClassStart,
			&course.ClassEnd,
			&course.Location,
			&course.CourseProgram,
			&course.Classday,
			&firstname,
			&lastname,
			&course.Semester,
			&course.Status,
			&course.Grade,
		)
		if err != nil {
			return nil, err
		}
		course.ProfessorName = firstname + " " + lastname
		courses = append(courses, course)
	}

	return courses, nil
}

func (r CourseRepositoryImplementation) GetProfessorCourse(professorId int) ([]response.Course, error) {

	query := `SELECT course_ID, 
				course_name, 
				ta_allocation, 
				work_hour 
			FROM courses
				WHERE professor_ID=$1`

	rows, err := r.db.Query(query, professorId)
	if err != nil {
		return nil, err
	}
	//garantees that connection is released back to the pool ,prevent leak
	defer rows.Close()

	var courses []response.Course
	for rows.Next() {
		var course response.Course

		err := rows.Scan(&course.CourseID, &course.CourseName, &course.TaAllocation, &course.WorkHour)
		if err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}

	return courses, nil
}

func (r CourseRepositoryImplementation) CreateCourse(body request.CreateCourse) (int, error) {

	queryCheck := `SELECT COUNT(*) FROM courses 
					WHERE course_ID=$1 
					AND course_program_ID=$2 
					AND sec=$3 
					AND semester_ID=$4 
					AND deleted_date IS NULL`

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

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	query := `INSERT INTO courses(
	course_ID, 
	course_name,
	professor_ID, 
	course_program_ID, 
	course_program, 
	sec, 
	semester_ID, 
	semester, 
	class_day_ID, 
	class_day, 
	class_start, 
	class_end,
	work_hour,
	ta_allocation,
	location,
	created_date) 
	values ($1,$2,$3, $4, $5 ,$6 ,$7 ,$8 ,$9 ,$10 ,$11 ,$12, $13, $14, $15, $16)
	RETURNING id`

	var lastInsertId int

	err = tx.QueryRow(query,
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
		body.WorkHour,
		body.TaAllocation,
		body.Location,
		time.Now(),
	).Scan(&lastInsertId)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	statusID := 1
	//Insert ta_jobposting
	query = `INSERT INTO ta_job_posting(
	professor_ID, 
	course_ID,
	grade_ID,
	task,
	ta_allocation,
	status_ID,
	created_date)
	values ($1,$2,$3, $4, $5 ,$6 ,$7 )
	RETURNING id`

	_, err = tx.Exec(query,
		body.ProfessorID,
		lastInsertId,
		body.GradeID,
		body.Task,
		body.TaAllocation,
		statusID,
		time.Now(),
	)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
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

	query = "UPDATE ta_job_posting SET deleted_date = $1 WHERE course_id = $2"
	_, err = r.db.Exec(query, time.Now(), id)
	if err != nil {
		return err
	}

	return nil
}

func (r CourseRepositoryImplementation) ApplyJobPost(body request.ApplyJobPost) (int, error) {

	//check student cannot make duplicate apply on same job_post
	var count int
	queryCheck := `SELECT COUNT(*) FROM ta_application 
					WHERE job_post_id = $1
					AND student_id = $2`

	err := r.db.QueryRow(queryCheck, body.JobPostID, body.StudentID).Scan(&count)
	if err != nil {
		return 0, err
	}

	if count > 0 {
		return 0, fmt.Errorf("student:%d already apply to this job", body.StudentID)
	}

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

	statusId := 5
	var applicationId int
	query = `INSERT INTO ta_application(
				transcript_ID, 
				student_ID, 
				status_ID, 
				job_post_ID,
				grade,
				purpose,
				created_date)
			VALUES($1, $2, $3, $4, $5, $6 ,$7)
			RETURNING id`

	err = tx.QueryRow(query,
		fileId,
		body.StudentID,
		statusId,
		body.JobPostID,
		body.Grade,
		body.JobPostID,
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
					tp.course_ID, 
					ta.created_date,
					st.status_value
				FROM ta_application AS ta 
				LEFT JOIN ta_job_posting AS tp
					ON ta.job_post_ID = tp.id
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
	query := `SELECT
				fs.file_name,
				fs.file_bytes
			FROM ta_application as ta
			LEFT JOIN transcript_storage as fs
			ON ta.transcript_ID = fs.transcript_ID
			WHERE ta.id = $1`

	var application response.ApplicationTrancript
	err := r.db.QueryRow(query, ApplicationId).Scan(
		&application.FileName,
		&application.Transcript,
	)
	if err != nil {
		return nil, err
	}
	return &application, nil
}

func (r CourseRepositoryImplementation) ApproveApplication(ApplicationId int) error {

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	//get courseId,studentId
	pendingStatus := 3
	var courseId int
	var studentId int
	query := `SELECT 
					jp.course_ID, 
					ta.student_ID 
				FROM ta_application as ta
				LEFT JOIN ta_job_posting as jp
					ON ta.job_post_ID = jp.id
				WHERE ta.id =$1 AND ta.status_id =$2`

	err = tx.QueryRow(query,
		ApplicationId,
		pendingStatus).Scan(&courseId, &studentId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed get courseId : %v", err)
	}

	//update status on ta_application
	approveStatus := 5
	query = `UPDATE ta_application SET status_ID = $1 WHERE id = $2`

	_, err = tx.Exec(query, approveStatus, ApplicationId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed update ta_application: %v", err)
	}

	//Insert new ta_course
	query = `INSERT INTO ta_courses(
				student_ID,
				course_ID,
				created_date)
				VALUES($1, $2, $3)`
	_, err = tx.Exec(query, studentId, courseId, time.Now())
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to Insert new ta_course : %v", err)
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed on commit transaction")
	}
	return nil
}

func (r CourseRepositoryImplementation) GetApplicationByProfessorId(professorId int) ([]response.Application, error) {
	query := `SELECT 
					ta.id,
					ta.student_ID, 
					ta.status_ID, 
					jp.course_ID, 
					ta.created_date,
					st.status_value,
					c.course_name,
					ta.grade,
					stu.firstname,
					stu.lastname
				FROM ta_application AS ta 
				LEFT JOIN status AS st
					ON ta.status_ID = st.status_ID
				LEFT JOIN ta_job_posting AS jp
					ON ta.job_post_ID = jp.id
				LEFT JOIN courses AS c
					ON jp.course_ID = c.id
				LEFT JOIN students AS stu
					ON ta.student_id = stu.student_id
				WHERE c.professor_ID = $1`

	rows, err := r.db.Query(query, professorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var applications []response.Application
	for rows.Next() {
		var application response.Application
		var firstname string
		var lastname string

		err := rows.Scan(
			&application.ApplicationId,
			&application.StudentID,
			&application.StatusID,
			&application.CourseID,
			&application.CreatedDate,
			&application.StatusCode,
			&application.CourseName,
			&application.Grade,
			&firstname,
			&lastname,
		)
		if err != nil {
			return nil, err
		}

		application.StudentName = firstname + " " + lastname
		applications = append(applications, application)
	}

	return applications, nil
}
