package service

import (
	"TA-management/internal/modules/shared/dto/response"

	"github.com/gin-gonic/gin"
)

type AuthenService interface {
	HandleGoogleCallback(ctx *gin.Context, code string) (string, *response.RequestDataResponse, error)
	// CheckUserRole(name string) (response.GeneralResponse, error)
}
