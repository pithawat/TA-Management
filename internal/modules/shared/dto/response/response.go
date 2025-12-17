package response

type GeneralResponse struct {
	Message string `json:"message"`
}

type RequestDataResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type CreateResponse struct {
	Message string `json:"message"`
	Id      int    `json:"id"`
}
