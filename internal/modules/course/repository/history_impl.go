package repository

import (
	"TA-management/internal/modules/course/dto/response"
	sharedResponse "TA-management/internal/modules/shared/dto/response"
	"database/sql"
)

// GetExpiredSemestersFromDB returns all semesters whose end_date has passed,
// with a count of how many courses belong to each.
func GetExpiredSemestersFromDB(db *sql.DB) ([]response.SemesterHistory, error) {
	query := `
		SELECT
			s.semester_ID,
			s.semester_value,
			s.start_date,
			s.end_date,
			COUNT(c.course_ID) AS course_count
		FROM semester s
		LEFT JOIN courses c ON c.semester_ID = s.semester_ID
		WHERE s.end_date < NOW()
		GROUP BY s.semester_ID, s.semester_value, s.start_date, s.end_date
		ORDER BY s.end_date DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []response.SemesterHistory
	for rows.Next() {
		var h response.SemesterHistory
		if err := rows.Scan(&h.SemesterID, &h.SemesterValue, &h.StartDate, &h.EndDate, &h.CourseCount); err != nil {
			return nil, err
		}
		results = append(results, h)
	}
	return results, nil
}

// GetCoursesBySemesterIDFromDB fetches courses for a given semester,
// including soft-deleted ones (so Finance can access historical data).
func GetCoursesBySemesterIDFromDB(db *sql.DB, semesterID int) ([]response.Course, error) {
	query := `
		SELECT
			c.course_ID,
			c.course_code,
			c.course_name,
			cp.course_program_value_thai,
			cd.class_day_value_thai,
			c.class_start,
			c.class_end,
			c.semester,
			c.sec,
			CONCAT(p.prefix, ' ', p.firstname_thai, ' ', p.lastname_thai) as fullname,
			c.work_hour,
			s.start_date,
			s.end_date,
			dc.role_id
		FROM courses AS c
		LEFT JOIN discord_channels AS dc
			ON c.course_ID = dc.course_ID
		LEFT JOIN professors AS p
			ON c.professor_ID = p.professor_ID
		LEFT JOIN class_days AS cd
			ON c.class_day_ID = cd.class_day_ID
		LEFT JOIN course_programs AS cp
			ON c.course_program_ID = cp.course_program_ID
		LEFT JOIN semester AS s
			ON c.semester_ID = s.semester_ID
		WHERE c.semester_ID = $1
		ORDER BY c.course_code ASC
	`

	rows, err := db.Query(query, semesterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []response.Course
	for rows.Next() {
		var course response.Course
		var fullname string
		var roleID sql.NullString
		err := rows.Scan(
			&course.CourseID,
			&course.CourseCode,
			&course.CourseName,
			&course.CourseProgram,
			&course.Classday,
			&course.ClassStart,
			&course.ClassEnd,
			&course.Semester,
			&course.Section,
			&fullname,
			&course.WorkHour,
			&course.SemesterStart,
			&course.SemesterEnd,
			&roleID,
		)
		if err != nil {
			return nil, err
		}
		if roleID.Valid {
			course.DiscordRoleID = roleID.String
		}
		course.ProfessorName = fullname
		courses = append(courses, course)
	}
	return courses, nil
}

// --- CourseRepositoryImplementation interface satisfaction ---

func (r CourseRepositoryImplementation) SoftDeleteExpiredData() error {
	return SoftDeleteExpiredSemesterData(r.db)
}

func (r CourseRepositoryImplementation) GetExpiredSemesters() (*sharedResponse.RequestDataResponse, error) {
	history, err := GetExpiredSemestersFromDB(r.db)
	if err != nil {
		return nil, err
	}
	return &sharedResponse.RequestDataResponse{
		Data:    history,
		Message: "Success",
	}, nil
}

func (r CourseRepositoryImplementation) GetCoursesBySemesterID(semesterID int) (*sharedResponse.RequestDataResponse, error) {
	courses, err := GetCoursesBySemesterIDFromDB(r.db, semesterID)
	if err != nil {
		return nil, err
	}
	return &sharedResponse.RequestDataResponse{
		Data:    courses,
		Message: "Success",
	}, nil
}
