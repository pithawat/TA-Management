package repository

import "TA-management/internal/modules/authen/dto/request"

type AuthenRepository interface {
	CheckUserRole(name string) (string, error)
	AddStudent(rq request.CreateStudent) error
	GetUserIDByName(name string, role string) (string, error)
}
