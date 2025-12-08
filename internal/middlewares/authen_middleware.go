package middleware

import (
	"TA-management/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AppClaims struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Pic   string `json:"pic"`
	jwt.RegisteredClaims
}

func AuthMiddleware(jwtSecret []byte) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, err := ctx.Cookie("auth_token")
		if err != nil || tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Token cookie misssing."})
			ctx.Abort()
			return
		}

		claims, err := utils.DecodeToken(tokenString, jwtSecret)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			ctx.Abort()
			return
		}
		ctx.Set("claims", claims)
		ctx.Next()

	}
}
