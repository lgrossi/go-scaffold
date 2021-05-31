package middlewares

import (
	"database/sql"
	"encoding/base32"
	"errors"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/lgrossi/go-scaffold/src/database"
	"github.com/lgrossi/go-scaffold/src/logger"
	"net/http"
	"time"
)

const (
	TokenString               = "MY_CUSTOM_SIGNED_STRING"
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
		user, err := VerifyToken(db, extractToken(c, AccessTokenCookieKey))

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		if user == nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.Set("user", user)
	}
}

func VerifyToken(db *sql.DB, signedToken string) (*database.User, error) {
	session := TokenSession{}
	token, _ := jwt.ParseWithClaims(signedToken, &session, parseToken)

	if token == nil || !token.Valid {
		return nil, errors.New("invalid access token")
	}

	user := database.GetUserByEmail(db, session.Email)
	if user == nil {
		return nil, nil
	}

	return user, nil
}

func RefreshToken(db *sql.DB, c *gin.Context) {
	refreshTokenSession := TokenSession{}
	signedRefreshToken := extractToken(c, RefreshTokenCookieKey)
	refreshToken, _ := jwt.ParseWithClaims(signedRefreshToken, &refreshTokenSession, parseToken)

	if refreshToken == nil || !refreshToken.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	user := database.GetUserById(db, refreshTokenSession.UserId)
	if user == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	generateAccessToken(c, user, AccessTokenLifetime)
	c.JSON(http.StatusCreated, gin.H{"session": gin.H{"user": user}})
}

func GenerateTokens(c *gin.Context, user *database.User) []string {
	return []string{
		generateAccessToken(c, user, AccessTokenLifetime),
		generateRefreshToken(c, user, RefreshTokenLifetime),
	}
}

func GenerateEmailVerificationLink(email string) string {
	signedToken := signToken(
		TokenSession{
			Email:          email,
			StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Unix() + VerificationTokenLifetime},
		},
	)

	encoded := base32.StdEncoding.EncodeToString([]byte(signedToken))

	return fmt.Sprintf("%s%s", "http://127.0.0.1:80/user/verification/", encoded)
}

func generateAccessToken(c *gin.Context, user *database.User, lifetime int) string {
	signedToken := signToken(
		TokenSession{
			UserId:         user.ID,
			Email:          user.Email,
			StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Unix() + int64(lifetime)},
		},
	)

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(AccessTokenCookieKey, signedToken, lifetime, "", "", true, true)

	return signedToken
}

func generateRefreshToken(c *gin.Context, user *database.User, lifetime int) string {
	signedToken := signToken(
		TokenSession{
			UserId:         user.ID,
			StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Unix() + int64(lifetime)},
		},
	)

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(RefreshTokenCookieKey, signedToken, lifetime, "", "", true, true)

	return signedToken
}

func signToken(claims jwt.Claims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(TokenString))

	if err != nil {
		logger.Panic(err)
	}

	return signedToken
}

func parseToken(token *jwt.Token) (interface{}, error) {
	//Make sure that the token method conform to "SigningMethodHMAC"
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return []byte(TokenString), fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return []byte(TokenString), nil
}

func extractToken(c *gin.Context, key string) string {
	token, err := c.Cookie(key)

	if err == nil {
		return token
	}

	return ""
}
