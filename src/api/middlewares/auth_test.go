package middlewares

import (
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestVerifyToken(t *testing.T) {
	type args struct {
		db *sql.DB
	}
	tests := []struct {
		name string
		args args
		want gin.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := VerifyToken(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VerifyToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createRefreshToken(t *testing.T) {
	type args struct {
		expireAt int64
	}
	tests := []struct {
		name string
		args args
	}{{
		name: "Creates refresh token with desired expiration",
		args: args{expireAt: time.Now().Unix() + 89102304589},
	}, {
		name: "Creates refresh token expiring now",
		args: args{expireAt: time.Now().Unix()},
	}, {
		name: "Creates refresh token expiring in the past",
		args: args{expireAt: time.Now().Unix() - 91284},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := jwt.StandardClaims{}
			jwt.ParseWithClaims(createRefreshToken(tt.args.expireAt), &session, func(token *jwt.Token) (interface{}, error) {
				return []byte(TokenString), nil
			})
			assert.Equal(t, tt.args.expireAt, session.ExpiresAt)
		})
	}
}

func Test_extractToken(t *testing.T) {
	tests := []struct {
		name   string
		cookie string
		want   string
	}{{
		name:   "Should parse valid token from cookie",
		cookie: "jwt=my_auth_token",
		want:   "my_auth_token",
	}, {
		name:   "Should return empty string for empty cookie",
		cookie: "",
		want:   "",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = &http.Request{}

			if tt.cookie != "" {
				c.Request.Header = http.Header{
					"Cookie": {fmt.Sprintf("%s", tt.cookie)},
				}
			}
			assert.Equal(t, tt.want, ExtractToken(c))
		})
	}
}

func Test_parseToken(t *testing.T) {
	type args struct {
		token *jwt.Token
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{{
		name:    "Returns error for invalid signing method",
		args:    args{token: &jwt.Token{Method: &jwt.SigningMethodRSA{}}},
		want:    nil,
		wantErr: true,
	}, {
		name:    "Returns token string for valid signing methods",
		args:    args{token: &jwt.Token{Method: &jwt.SigningMethodHMAC{}}},
		want:    []byte(TokenString),
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseToken(tt.args.token)
			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
