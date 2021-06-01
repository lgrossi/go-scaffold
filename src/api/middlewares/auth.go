package middlewares

import (
	"database/sql"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/lgrossi/go-scaffold/src/database"
	"github.com/lgrossi/go-scaffold/src/security"
	"net/http"
)

const (
	AccessTokenCookieKey      = "JWT_ACCESS_TOKEN"
	RefreshTokenCookieKey     = "JWT_REFRESH_TOKEN"
	AccessTokenLifetime       = 15 * 60
	RefreshTokenLifetime      = 7 * 24 * 60 * 60
	VerificationTokenLifetime = 1 * 60 * 60
)

type TokenSession struct {
	UserId int64
	Email  string
	jwt.StandardClaims
}

func VerifyTokenHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		manager := security.ExtractToken(c, AccessTokenCookieKey)
		user := manager.VerifyToken(db)

		if manager.Error != nil || user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid access token"})
			return
		}

		c.Set("user", user)
	}
}

func GenerateEmailVerificationURL(user *database.User) string {
	return security.InitToken(user, VerificationTokenLifetime).Sign().GenerateVerificationURL()
}

func GenerateAccessToken(c *gin.Context, user *database.User, lifetime ...int64) *security.TokenManager {
	if len(lifetime) > 0 {
		return security.InitToken(user, lifetime[0]).Sign().SetTokenToContext(c, AccessTokenCookieKey)
	}
	return security.InitToken(user, AccessTokenLifetime).Sign().SetTokenToContext(c, AccessTokenCookieKey)
}

func GenerateRefreshToken(c *gin.Context, user *database.User, lifetime ...int64) *security.TokenManager {
	if len(lifetime) > 0 {
		return security.InitToken(user, lifetime[0]).Sign().SetTokenToContext(c, RefreshTokenCookieKey)
	}
	return security.InitToken(user, RefreshTokenLifetime).Sign().SetTokenToContext(c, RefreshTokenCookieKey)
}
