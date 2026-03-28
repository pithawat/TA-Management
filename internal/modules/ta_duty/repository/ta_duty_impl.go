package repository

import (
	"TA-management/internal/modules/ta_duty/dto/request"
	"TA-management/internal/modules/ta_duty/dto/response"
	"TA-management/internal/utils"
	"database/sql"
	"time"
)

type TaDutyRepositoryImplentation struct {
	db *sql.DB
}

func NewTaDutyRepository(db *sql.DB) *TaDutyRepositoryImplentation {
	return &TaDutyRepositoryImplentation{db: db}
}

func (r TaDutyRepositoryImplentation) GetTADutyRoadmap(courseID int, studentID int) (*[]response.DutyChecklistItem, error) {
	query := `
		WITH RECURSIVE semester_dates AS (
			SELECT s.start_date, s.end_date
			FROM semester s
			JOIN courses c ON c.semester_ID = s.semester_ID
			WHERE c.course_ID = $1
		),
		all_dates AS (
			SELECT start_date AS d_date FROM semester_dates
			UNION ALL
			SELECT (d_date + INTERVAL '1 day')::date FROM all_dates
			WHERE d_date < (SELECT end_date FROM semester_dates)
		)
		SELECT 
			ad.d_date,
			CASE 
				WHEN h.id IS NOT NULL THEN 'COMPLETED'
				WHEN ad.d_date < CURRENT_DATE THEN 'MISSED'
				WHEN ad.d_date = CURRENT_DATE THEN 'TODAY'
				ELSE 'UPCOMING'
			END as status,
			CASE WHEN h.id IS NOT NULL THEN true ELSE false END as is_checked
		FROM all_dates ad
		JOIN courses c ON c.course_ID = $1
		JOIN class_days cd ON c.class_day_ID = cd.class_day_ID
		LEFT JOIN holidays hol ON ad.d_date = hol.holiday_date
		LEFT JOIN ta_duty_historys h ON h.course_ID = c.course_ID 
			AND h.student_ID = $2 
			AND h.date::date = ad.d_date
		WHERE trim(to_char(ad.d_date, 'Day')) = cd.class_day_value
		AND hol.id IS NULL
		ORDER BY ad.d_date ASC;
    `
	rows, err := r.db.Query(query, courseID, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roadmap []response.DutyChecklistItem
	for rows.Next() {
		var item response.DutyChecklistItem
		var dutyTime time.Time
		if err := rows.Scan(&dutyTime, &item.Status, &item.IsChecked); err != nil {
			return nil, err
		}
		item.Date = dutyTime.Format("2006-01-02")
		roadmap = append(roadmap, item)
	}
	return &roadmap, nil
}

func (r TaDutyRepositoryImplentation) MarkDutyAsDone(courseID int, studentID int, date string) error {
	query := `
		INSERT INTO ta_duty_historys (date, course_ID, student_ID)
		VALUES ($1, $2, $3)

	`
	_, err := r.db.Exec(query, date, courseID, studentID)
	if err != nil {
		return err
	}
	return nil
}

func (r TaDutyRepositoryImplentation) GetTADutyDataExportPayment(courseID int, month int) (*[]request.CreatePaymentData, *request.CourseDutyData, error) {

	var result []request.CreatePaymentData

	query := `SELECT 
					s.student_ID,
					s.firstname_thai || ' ' || s.lastname_thai,
					c.work_hour,
					c.class_start || ' ' || c.class_end || ' น. ' AS time_range,
					c.course_code,
					c.course_name,
					c.semester,
					c.sec
				FROM ta_courses AS tc
				LEFT JOIN students AS s
					ON s.student_ID = tc.student_ID
				LEFT JOIN courses AS c
					ON c.course_ID = tc.course_ID
				where tc.course_ID = $1
				`
	rows, err := r.db.Query(query, courseID)
	if err != nil {
		return nil, nil, err
	}

	type studentBase struct {
		id        int
		name      string
		workhour  int
		timeRange string
	}

	var courseData request.CourseDutyData
	var students []studentBase
	for rows.Next() {
		var student studentBase
		err := rows.Scan(&student.id,
			&student.name,
			&student.workhour,
			&student.timeRange,
			&courseData.CourseCode,
			&courseData.CourseName,
			&courseData.Semester,
			&courseData.Sec,
		)
		if err != nil {
			return nil, nil, err
		}
		students = append(students, student)
	}
	rows.Close()

	query = `
		WITH RECURSIVE semester_dates AS (
			SELECT s.start_date, s.end_date
			FROM semester s
			JOIN courses c ON c.semester_ID = s.semester_ID
			WHERE c.course_ID = $1
		),
		all_dates AS (
			SELECT start_date AS d_date FROM semester_dates
			UNION ALL
			SELECT (d_date + INTERVAL '1 day')::date FROM all_dates
			WHERE d_date < (SELECT end_date FROM semester_dates)
		)
		SELECT 
			ad.d_date,
			CASE 
				WHEN h.id IS NOT NULL THEN 'COMPLETED'
				WHEN ad.d_date < CURRENT_DATE THEN 'MISSED'
				WHEN ad.d_date = CURRENT_DATE THEN 'TODAY'
				ELSE 'UPCOMING'
			END as status,
			CASE WHEN h.id IS NOT NULL THEN true ELSE false END as is_checked
		FROM all_dates ad
		JOIN courses c ON c.course_ID = $1
		JOIN class_days cd ON c.class_day_ID = cd.class_day_ID
		LEFT JOIN holidays hol ON ad.d_date = hol.holiday_date
		LEFT JOIN ta_duty_historys h ON h.course_ID = c.course_ID 
			AND h.student_ID = $2 
			AND h.date::date = ad.d_date
		WHERE trim(to_char(ad.d_date, 'Day')) = cd.class_day_value -- Only show the class days
		AND hol.id IS NULL -- Exclude Holidays
		AND EXTRACT(MONTH FROM ad.d_date) = $3
		ORDER BY ad.d_date ASC;
    `

	for _, student := range students {

		rows, err := r.db.Query(query, courseID, student.id, month)
		if err != nil {
			return nil, nil, err
		}

		var roadmap []request.DutyChecklistItem
		for rows.Next() {
			var item request.DutyChecklistItem
			var dutyTime time.Time

			if err := rows.Scan(&dutyTime, &item.Status, &item.IsChecked); err != nil {
				rows.Close()
				return nil, nil, err
			}

			thaiDate, err := utils.ConvertTOThaiDate(dutyTime.Format("2006-01-02"))
			if err != nil {
				return nil, nil, err
			}

			item.Date = thaiDate
			item.TimeRange = student.timeRange
			roadmap = append(roadmap, item)
		}
		result = append(result, request.CreatePaymentData{
			StudentName: student.name,
			WorkHour:    student.workhour,
			Duty:        roadmap,
		})
		rows.Close()
	}

	return &result, &courseData, nil
}

func (r TaDutyRepositoryImplentation) GetTADutyDataExportSignature(courseId int, month int) (*request.CreateSignatureSheet, *request.CourseDutyData, error) {

	query := `SELECT 
				c.course_code,
				c.course_name,
				c.semester,
				c.sec
			FROM ta_courses AS tc
			LEFT JOIN courses AS c 
				ON c.course_ID = tc.course_ID
			WHERE tc.course_ID = $1
	`
	var courseData request.CourseDutyData
	err := r.db.QueryRow(query, courseId).Scan(&courseData.CourseCode, &courseData.CourseName, &courseData.Semester, &courseData.Sec)
	if err != nil {
		return nil, nil, err
	}

	studentDataQuery := `SELECT 
							s.firstname_thai || ' ' || s.lastname_thai AS name
						FROM ta_courses AS tc
						LEFT JOIN students AS s
							ON s.student_ID = tc.student_ID
						WHERE tc.course_ID = $1						
	`

	var rq request.CreateSignatureSheet
	rows, err := r.db.Query(studentDataQuery, courseId)
	if err != nil {
		return nil, nil, err
	}

	for rows.Next() {
		var studentName string
		err := rows.Scan(&studentName)
		if err != nil {
			return nil, nil, err
		}

		rq.TAName = append(rq.TAName, studentName)
	}
	rows.Close()

	dutyTimeQuery := `
		WITH RECURSIVE semester_dates AS (
			SELECT s.start_date, s.end_date
			FROM semester s
			JOIN courses c ON c.semester_ID = s.semester_ID
			WHERE c.course_ID = $1
		),
		all_dates AS (
			SELECT start_date AS d_date FROM semester_dates
			UNION ALL
			SELECT (d_date + INTERVAL '1 day')::date FROM all_dates
			WHERE d_date < (SELECT end_date FROM semester_dates)
		)
		SELECT 
			ad.d_date
		FROM all_dates ad
		JOIN courses c ON c.course_ID = $1
		JOIN class_days cd ON c.class_day_ID = cd.class_day_ID
		LEFT JOIN holidays hol ON ad.d_date = hol.holiday_date
		WHERE trim(to_char(ad.d_date, 'Day')) = cd.class_day_value -- Only show the class days
		AND hol.id IS NULL -- Exclude Holidays
		AND EXTRACT(MONTH FROM ad.d_date) = $2
		ORDER BY ad.d_date ASC;
    `
	rows, err = r.db.Query(dutyTimeQuery, courseId, month)
	if err != nil {
		return nil, nil, err
	}

	for rows.Next() {
		var dutyDate time.Time
		err := rows.Scan(&dutyDate)
		if err != nil {
			return nil, nil, err
		}

		thaiDate, err := utils.ConvertTOThaiDate(dutyDate.Format("2006-01-02"))
		if err != nil {
			return nil, nil, err
		}

		rq.DutyDate = append(rq.DutyDate, thaiDate)
	}
	rows.Close()

	return &rq, &courseData, nil
}
