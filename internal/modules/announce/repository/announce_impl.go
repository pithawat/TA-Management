package repository

import (
	"TA-management/internal/modules/announce/dto/request"
	"TA-management/internal/modules/announce/dto/response"
	"database/sql"
	"time"
)

type AnnouncementRepoImplementation struct {
	db *sql.DB
}

func NewAnnouncementRepository(db *sql.DB) *AnnouncementRepoImplementation {
	return &AnnouncementRepoImplementation{db: db}
}

func (r AnnouncementRepoImplementation) GetStudentEmailByCourseID(courseID int) (*request.EmailRequest, *response.CourseDetail, error) {

	query := `SELECT 
				st.email,
				c.course_name,
				c.course_code
				FROM ta_courses AS tc
				LEFT JOIN students AS st ON st.student_ID=tc.student_ID
				LEFT JOIN courses AS c ON c.course_ID = tc.course_ID
				WHERE tc.course_ID=$1`

	var emailRequest request.EmailRequest
	var courseDetail response.CourseDetail
	var email string
	rows, err := r.db.Query(query, courseID)
	if err != nil {
		return nil, nil, err
	}

	for rows.Next() {
		err := rows.Scan(&email, &courseDetail.CourseName, &courseDetail.CourseCode)
		if err != nil {
			rows.Close()
			return nil, nil, err
		}
		emailRequest.To = append(emailRequest.To, email)
	}
	rows.Close()

	return &emailRequest, &courseDetail, nil
}

func (r AnnouncementRepoImplementation) GetStudentEmailByCourseIDs() (*request.EmailRequest, error) {

	query := `SELECT DISTINCT
				st.email
				FROM ta_courses AS tc
				JOIN students AS st ON st.student_ID=tc.student_ID`

	var emailRequest request.EmailRequest

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			return nil, err
		}
		emailRequest.To = append(emailRequest.To, email)
	}

	return &emailRequest, nil
}

func (r AnnouncementRepoImplementation) GetStudentEmailByStudentID(studentID int) (*request.EmailRequest, error) {
	query := `SELECT email FROM students WHERE student_ID = $1`

	var emailRequest request.EmailRequest
	var email string
	err := r.db.QueryRow(query, studentID).Scan(&email)
	if err != nil {
		return nil, err
	}

	emailRequest.To = append(emailRequest.To, email)

	return &emailRequest, nil
}

func (r AnnouncementRepoImplementation) SaveEmailHistory(rq request.CreateEmailHistory) error {

	query := `INSERT INTO email_history (
		subject, 
		body, 
		received_name,
		n_received, 
		status_ID, 
		created_date) VALUES($1, $2, $3, $4, $5, $6) `

	_, err := r.db.Exec(query,
		rq.Subject,
		rq.Body,
		rq.ReceivedName,
		rq.NReceived,
		rq.StatusID,
		time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (r AnnouncementRepoImplementation) GetEmailHistory() (*[]response.EmailHistory, error) {

	query := `SELECT
				eh.id,
				eh.subject,
				eh.body,
				eh.received_name,
				eh.n_received,
				s.status_value,
				eh.created_date
			FROM email_history eh
			LEFT JOIN status AS s ON eh.status_ID = s.status_ID
			ORDER BY eh.created_date desc limit 10
			`
	var emailHistorys []response.EmailHistory
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var emailHistory response.EmailHistory
		err := rows.Scan(
			&emailHistory.Id,
			&emailHistory.Subject,
			&emailHistory.Body,
			&emailHistory.ReceivedName,
			&emailHistory.NReceived,
			&emailHistory.Status,
			&emailHistory.CreatedDate)
		if err != nil {
			return nil, err
		}
		emailHistorys = append(emailHistorys, emailHistory)
	}

	return &emailHistorys, err

}

func (r AnnouncementRepoImplementation) GetDiscordRoleID(courseID int) (string, error) {
	query := `SELECT 
				role_id
				FROM discord_channels
				WHERE course_ID = $1
	`
	var roleID string
	err := r.db.QueryRow(query, courseID).Scan(&roleID)
	if err != nil {
		return "", err
	}

	return roleID, nil
}

func (r AnnouncementRepoImplementation) CreateNewDiscordChannel(roleID string, channelID string, channelName string, courseID int) error {

	query := `INSERT INTO discord_channels (
				channel_ID, 
				channel_name,
				role_ID, 
				course_ID) VALUES($1, $2, $3, $4)`

	_, err := r.db.Exec(query,
		channelID,
		channelName,
		roleID,
		courseID)

	if err != nil {
		return err
	}
	return nil

}
