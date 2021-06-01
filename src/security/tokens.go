package security

import (
	"database/sql"
	"encoding/base32"
	"errors"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/lgrossi/go-scaffold/src/configs"
	"github.com/lgrossi/go-scaffold/src/database"
	"github.com/lgrossi/go-scaffold/src/logger"
	"net/http"
	"time"
)

const (
	EnvKeyTokenSecret = "TOKEN_SECRET"
)

type TokenSession struct {
	UserId int64
	Email  string
	jwt.StandardClaims
}

type TokenManager struct {
	Session  TokenSession
	Signed   string
	Token    *jwt.Token
	Lifetime int64
	Error    error
}

func InitToken(user *database.User, lifetime int64) *TokenManager {
	return &TokenManager{
		Lifetime: lifetime,
		Session: TokenSession{
			UserId:         user.ID,
			Email:          user.Email,
			StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Unix() + lifetime},
		},
	}
}

func (manager *TokenManager) Sign() *TokenManager {
	manager.Token = jwt.NewWithClaims(jwt.SigningMethodHS256, &manager.Session)
	manager.Signed, manager.Error = manager.Token.SignedString([]byte(getTokenSecret()))

	if manager.Error != nil {
		logger.Panic(manager.Error)
	}

	return manager
}

func (manager *TokenManager) SetTokenToContext(c *gin.Context, key string) *TokenManager {
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(key, manager.Signed, int(manager.Lifetime), "", "", true, true)
	return manager
}

func (manager *TokenManager) GenerateVerificationURL() string {
	encoded := base32.StdEncoding.EncodeToString([]byte(manager.Signed))
	return fmt.Sprintf("%s%s", "http://localhost:3000/verify/", encoded)
}

func (manager *TokenManager) VerifyToken(db *sql.DB) *database.User {
	manager.Token, _ = jwt.ParseWithClaims(manager.Signed, &manager.Session, parseToken)

	if manager.Token == nil || !manager.Token.Valid {
		manager.Error = errors.New("invalid token")
		return nil
	}

	user := database.GetUserByEmail(db, manager.Session.Email)
	if user == nil {
		return nil
	}

	return user
}

func ExtractToken(c *gin.Context, key string) *TokenManager {
	token, err := c.Cookie(key)

	if err == nil {
		return &TokenManager{Signed: token}
	}

	return &TokenManager{}
}

func parseToken(token *jwt.Token) (interface{}, error) {
	//Make sure that the token method conform to "SigningMethodHMAC"
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return []byte(getTokenSecret()), fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return []byte(getTokenSecret()), nil
}

func getTokenSecret() string {
	return configs.GetEnvStr(EnvKeyTokenSecret, "MY_CUSTOM_SIGNED_STRING")
}
