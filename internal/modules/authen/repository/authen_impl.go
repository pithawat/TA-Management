package repository

import (
	"database/sql"
	"fmt"
	"strings"
)

type AuthenRepositoryImplementation struct {
	db *sql.DB
}

func NewAuthenRepository(DB *sql.DB) AuthenRepositoryImplementation {
	return AuthenRepositoryImplementation{db: DB}
}

func (r AuthenRepositoryImplementation) CheckUserRole(name string) (string, error) {
	parts := strings.Split(name, " ")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid name format")
	}
	firstname := parts[0]
	lastname := parts[1]
	fmt.Print(firstname, lastname)

	var isProf bool
	query := "SELECT EXISTS(SELECT 1 FROM professors WHERE firstname = $1 AND lastname = $2)"
	err := r.db.QueryRow(query, firstname, lastname).Scan(&isProf)
	if err != nil {
		return "", err
	}
	if isProf {
		return "professor", nil
	}

	var isAccount bool
	query = "SELECT EXISTS(SELECT 1 FROM accountants WHERE firstname = $1 AND lastname = $2)"
	err = r.db.QueryRow(query, firstname, lastname).Scan(&isProf)
	if err != nil {
		return "", err
	}
	if isAccount {
		return "account", nil
	}

	return "student", nil
}
