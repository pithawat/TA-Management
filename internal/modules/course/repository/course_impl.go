package repository

import (
	"TA-management/internal/modules/course/dto/request"
	"TA-management/internal/modules/course/dto/response"
	"database/sql"
	"errors"
	"fmt"
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

func (r CourseRepositoryImplementation) CreateCourse(body request.CreateCourse) error {

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
		return fmt.Errorf("failed to scan duplicate check result:%w", err)
	}

	if count > 0 {
		return errors.New("course already exists")
	}

	query := "INSERT INTO courses(course_ID, course_name,professor_ID, course_program_ID, course_program, sec, semester_ID, semester, class_day_ID, class_day, class_start, class_end) values ($1,$2,$3, $4, $5 ,$6 ,$7 ,$8 ,$9 ,$10 ,$11 ,$12)"
	_, err = r.db.Exec(query,
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
	)

	if err != nil {
		return err
	}

	return nil

}
