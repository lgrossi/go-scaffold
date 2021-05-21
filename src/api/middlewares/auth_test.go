package middlewares

import (
	"database/sql"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestCreateAccessToken(t *testing.T) {
	type args struct {
		email        string
		refreshToken []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{{
		name: "Parse includes email in the session",
		args: args{email: "a@a.com"},
		want: "a@a.com",
	}, {
		name: "Parse uses refresh token from argument",
		args: args{email: "a@a.com", refreshToken: []string{"aaa"}},
		want: "a@a.com",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateAccessToken(tt.args.email, tt.args.refreshToken...)

			session := TokenSession{}
			jwt.ParseWithClaims(got, &session, func(token *jwt.Token) (interface{}, error) {
				return []byte(TokenString), nil
			})

			assert.Nil(t, err)
			assert.Equal(t, tt.want, session.Email)
			assert.NotEmpty(t, session.RefreshToken)
			if len(tt.args.refreshToken) > 0 {
				assert.Equal(t, tt.args.refreshToken[0], session.RefreshToken)
			}
		})
	}
}

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
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{{
		name: "Should parse valid token from header",
		args: args{&gin.Context{Request: &http.Request{Header: map[string][]string{"Authorization": {"authorization token"}}}}},
		want: "token",
	}, {
		name: "Should return empty string for empty header",
		args: args{&gin.Context{Request: &http.Request{}}},
		want: "",
	}, {
		name: "Should return empty for invalid header",
		args: args{&gin.Context{Request: &http.Request{Header: map[string][]string{"Authorization": {"authorization my token"}}}}},
		want: "",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ExtractToken(tt.args.c))
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

func Test_refreshAccessToken(t *testing.T) {
	type args struct {
		session *TokenSession
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{{
		name: "Invalid token returns error",
		args: args{session: &TokenSession{
			Email:        "a@a.com",
			RefreshToken: "invalid",
		}},
		want:    "",
		wantErr: true,
	}, {
		name: "Expired token returns error",
		args: args{session: &TokenSession{
			Email:        "a@a.com",
			RefreshToken: createRefreshToken(time.Now().Unix() - 1),
		}},
		want:    "",
		wantErr: true,
	}, {
		name: "Valid refresh token returns new access token",
		args: args{session: &TokenSession{
			Email:        "a@a.com",
			RefreshToken: createRefreshToken(time.Now().Unix() + 10),
		}},
		want:    "a@a.com",
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := refreshAccessToken(tt.args.session)
			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)

			session := TokenSession{}
			jwt.ParseWithClaims(got, &session, func(token *jwt.Token) (interface{}, error) {
				return []byte(TokenString), nil
			})
			assert.Equal(t, tt.want, session.Email)
		})
	}
}
