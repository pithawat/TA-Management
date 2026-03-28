package repository

import (
	"database/sql"
	"time"
)

// SoftDeleteExpiredSemesterData soft-deletes course-related data for all
// semesters whose end_date has passed. Must be run in safe cascade order:
// ta_application → ta_job_posting → ta_courses → courses.
func SoftDeleteExpiredSemesterData(db *sql.DB) error {
	now := time.Now()

	// 1. Soft-delete ta_application records linked to expired semesters
	_, err := db.Exec(`
		UPDATE ta_application
		SET deleted_date = $1
		WHERE deleted_date IS NULL
		  AND job_post_ID IN (
		      SELECT j.id FROM ta_job_posting j
		      INNER JOIN courses c ON j.course_ID = c.course_ID
		      INNER JOIN semester s ON c.semester_ID = s.semester_ID
		      WHERE s.end_date < $2
		  )
	`, now, now)
	if err != nil {
		return err
	}

	// 2. Soft-delete ta_job_posting records linked to expired semesters
	_, err = db.Exec(`
		UPDATE ta_job_posting
		SET deleted_date = $1
		WHERE deleted_date IS NULL
		  AND course_ID IN (
		      SELECT c.course_ID FROM courses c
		      INNER JOIN semester s ON c.semester_ID = s.semester_ID
		      WHERE s.end_date < $2
		  )
	`, now, now)
	if err != nil {
		return err
	}

	// 3. Soft-delete ta_courses records linked to expired semesters
	_, err = db.Exec(`
		UPDATE ta_courses
		SET deleted_date = $1
		WHERE deleted_date IS NULL
		  AND course_ID IN (
		      SELECT c.course_ID FROM courses c
		      INNER JOIN semester s ON c.semester_ID = s.semester_ID
		      WHERE s.end_date < $2
		  )
	`, now, now)
	if err != nil {
		return err
	}

	// 4. Soft-delete courses linked to expired semesters
	_, err = db.Exec(`
		UPDATE courses
		SET deleted_date = $1
		WHERE deleted_date IS NULL
		  AND semester_ID IN (
		      SELECT semester_ID FROM semester
		      WHERE end_date < $2
		  )
	`, now, now)
	if err != nil {
		return err
	}

	return nil
}
