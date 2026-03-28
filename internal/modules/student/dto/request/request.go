package request

type UpdateProfile struct {
	FirstnameThai string `json:"firstnameThai" binding:"required"`
	LastnameThai  string `json:"lastnameThai" binding:"required"`
	PhoneNumber   string `json:"phoneNumber" binding:"required"`
}
