package repository

import (
	"TA-management/internal/modules/course/dto/request"
	"TA-management/internal/modules/course/dto/response"
	"database/sql"
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

	query := "INSERT INTO courses(course_ID, course_name, created_date, created_by) values ($1,$2,$3,$4)"
	_, err := r.db.Exec(query,
		body.CourseID,
		body.CourseName,
		body.CreatedDate,
		"linkky",
	)

	if err != nil {
		return err
	}

	return nil

}
