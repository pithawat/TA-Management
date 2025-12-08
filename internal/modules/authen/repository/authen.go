package repository

type AuthenRepository interface {
	CheckUserRole(name string) (string, error)
}
