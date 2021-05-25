package middlewares

import (
	"database/sql"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/lgrossi/go-scaffold/src/database"
	"github.com/lgrossi/go-scaffold/src/logger"
	"net/http"
	"time"
)

const (
	TokenString            = "MY_CUSTOM_SIGNED_STRING"
	AccessTokenExpiration  = 15 * 60
	RefreshTokenExpiration = 7 * 24 * 60 * 60
)

type TokenSession struct {
	Email        string
	RefreshToken string
	jwt.StandardClaims
}

func GenerateAccessToken(db *sql.DB, c *gin.Context, email string, refreshToken ...string) *database.ActiveToken {
	if len(refreshToken) < 1 {
		refreshToken = []string{createRefreshToken(time.Now().Unix() + RefreshTokenExpiration)}
	}

	session := TokenSession{
		Email:          email,
		RefreshToken:   refreshToken[0],
		StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Unix() + AccessTokenExpiration},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, session)
	tokenStr, err := token.SignedString([]byte(TokenString))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		logger.Panic(err)
	}

	activeToken := database.ActiveToken{TokenStr: tokenStr, Email: email}
	database.InsertActiveToken(db, activeToken)

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("jwt", tokenStr, AccessTokenExpiration, "", "", true, false)

	return &activeToken
}

func VerifyToken(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := ExtractToken(c)
		session := TokenSession{}

		token, _ := jwt.ParseWithClaims(tokenString, &session, parseToken)

		if token == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		activeToken := database.ActiveToken{Email: session.Email, TokenStr: tokenString}
		if !database.IsActiveToken(db, &activeToken) {
			if activeToken.ExpiresAt > 0 {
				database.DeactivateToken(db, activeToken)
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		if !token.Valid {
			database.DeactivateToken(db, activeToken)

			refreshedToken, err := refreshAccessToken(db, c, &session)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				return
			}

			activeToken = *refreshedToken
		}

		c.Set("jwt", activeToken.TokenStr)
		c.Set("session", &session)
		c.Next()
	}
}

func parseToken(token *jwt.Token) (interface{}, error) {
	//Make sure that the token method conform to "SigningMethodHMAC"
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return []byte(TokenString), fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return []byte(TokenString), nil
}

func refreshAccessToken(db *sql.DB, c *gin.Context, session *TokenSession) (*database.ActiveToken, error) {
	refreshTokenStr := session.RefreshToken

	_, err := jwt.Parse(refreshTokenStr, parseToken)
	if err != nil {
		return nil, err
	}

	return GenerateAccessToken(db, c, session.Email, refreshTokenStr), nil
}

func createRefreshToken(expiresAt int64) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{ExpiresAt: expiresAt})
	tokenStr, err := token.SignedString([]byte(TokenString))

	if err != nil {
		logger.Panic(err)
	}

	return tokenStr
}

func ExtractToken(c *gin.Context) string {
	token, err := c.Cookie("jwt")

	if err == nil {
		return token
	}

	return ""
}
