package controller

import (
	"TA-management/internal/modules/authen/service"
	"TA-management/internal/utils"

	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type AuthController struct {
	service           service.AuthenService
	googleOAuthConfig *oauth2.Config
}

func NewAuthenController(authenService service.AuthenService, config *oauth2.Config) *AuthController {
	return &AuthController{
		service:           authenService,
		googleOAuthConfig: config,
	}
}

func InitializeController(authenService service.AuthenService, googleOAuthConfig *oauth2.Config, r *gin.RouterGroup) {
	c := NewAuthenController(authenService, googleOAuthConfig)
	r.Use()
	{
		r.GET("/google", c.handleLogin)
		r.GET("google/callback", c.handleCallback)

	}
}

func (controller AuthController) handleCallback(ctx *gin.Context) {
	queryState := ctx.Query("state")
	code := ctx.Query("code")

	if queryState == "" || code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters: state or code missing."})
		return
	}

	cookieState, err := ctx.Cookie("oauth_state")
	if err != nil || cookieState == "" || cookieState != queryState {
		// Log the failure reason internally (optional but recommended)
		fmt.Printf("State check failed. Cookie Error: %v, Cookie State: %s, Query State: %s\n", err, cookieState, queryState)

		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "State mismatch or session expired. Invalid CSRF protection."})
		return
	}

	signedJWT, user, err := controller.service.HandleGoogleCallback(ctx, code)

	if err != nil {
		// Map service errors to appropriate HTTP status codes
		status := http.StatusInternalServerError
		if err.Error() == "email not verified" {
			status = http.StatusForbidden
		} else if err.Error() == "code exchange failed" {
			status = http.StatusBadRequest
		}
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}
	ctx.SetCookie("auth_token", signedJWT, 3600*24*7, "/", "localhost", false, true) // Cookie lasts 7 days
	ctx.JSON(http.StatusOK, user)
}

func (controller AuthController) handleLogin(ctx *gin.Context) {
	state, err := utils.RandState(24)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create state"})
	}
	utils.SetStateCookie(ctx.Writer, state)

	url := controller.googleOAuthConfig.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("prompt", "consent"),
	)
	ctx.JSON(http.StatusOK, gin.H{"url": url})
}
