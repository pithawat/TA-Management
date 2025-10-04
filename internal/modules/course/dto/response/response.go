package response

type GeneralResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type Course struct {
	CourseID   string `json:"courseID"`
	CourseName string `json:"courseName"`
}
