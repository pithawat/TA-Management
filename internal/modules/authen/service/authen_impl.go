package service

import (
	authresponse "TA-management/internal/modules/authen/dto"
	"TA-management/internal/modules/authen/repository"
	"TA-management/internal/modules/shared/dto/response"

	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

type AuthenServiceImplementation struct {
	repo              repository.AuthenRepository
	googleOAuthConfig *oauth2.Config
	jwtSecret         []byte
}

func NewAuthenService(repo repository.AuthenRepository, config *oauth2.Config, secret []byte) AuthenServiceImplementation {
	return AuthenServiceImplementation{
		repo:              repo,
		googleOAuthConfig: config,
		jwtSecret:         secret,
	}
}

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verifiedEmail"`
	Name          string `json:"name"`
}

type AppClaims struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

func (s AuthenServiceImplementation) HandleGoogleCallback(ctx *gin.Context, code string) (string, *response.RequestDataResponse, error) {

	token, err := s.googleOAuthConfig.Exchange(ctx, code)
	if err != nil {
		return "", nil, errors.New("code exchange failed")
	}

	client := s.googleOAuthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return "", nil, errors.New("failed to query userinfo")
	}
	defer resp.Body.Close()

	var gu GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&gu); err != nil {
		return "", nil, errors.New("failed to decode userinfo")
	}

	// if !gu.VerifiedEmail {
	// 	return "", nil, errors.New("email not verified")
	// }

	role, err := s.repo.CheckUserRole(gu.Name)
	if err != nil {
		fmt.Print(err)
		return "", nil, errors.New("failed to check user role")
	}

	now := time.Now()
	claims := AppClaims{
		Sub:   gu.ID,
		Email: gu.Email,
		Name:  gu.Name,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "example.com/googlelogin",
			Subject:   gu.ID,
			Audience:  []string{"TA-mangement"},
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	j := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := j.SignedString(s.jwtSecret)
	if err != nil {
		fmt.Print(err)
		return "", nil, errors.New("failed to sign jwt")
	}

	data := authresponse.LoginResponse{
		Id:    gu.ID,
		Email: gu.Email,
		Name:  gu.Name,
		Role:  role,
	}

	response := response.RequestDataResponse{
		Data:    data,
		Message: ("Login Success"),
	}

	return signed, &response, nil
}

// func (s AuthenServiceImplementation) CheckUserRole(name string) (string, error) {

// 	role,err := s.CheckUserRole(name)
// }
