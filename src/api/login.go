package api

import (
	"context"
	"encoding/base32"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lgrossi/go-scaffold/src/api/external_clients"
	"github.com/lgrossi/go-scaffold/src/api/middlewares"
	"github.com/lgrossi/go-scaffold/src/database"
	exampleProtoMessages "github.com/lgrossi/go-scaffold/src/grpc/example_proto_defs"
	"github.com/lgrossi/go-scaffold/src/logger"
	"github.com/lgrossi/go-scaffold/src/security"
	"net/http"
)

func (_api *Api) register(c *gin.Context) {
	request := &database.RegisterRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logger.Panic(err)
	}

	user := database.CreateUser(_api.DB, request)
	if user == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "This email is already in use."})
		return
	}

	go external_clients.SendEmail(
		user.Email,
		"Email Verification",
		fmt.Sprintf(
			`Hello %s,<br><br>Thank you for registrate with us<br>Please click <a href="%s">here</a> to validate your email.`,
			user.Name,
			middlewares.GenerateEmailVerificationURL(user),
		),
	)

	middlewares.GenerateAccessToken(c, user)
	middlewares.GenerateRefreshToken(c, user)
	c.JSON(http.StatusCreated, gin.H{"session": gin.H{"user": user}})
}

func (_api *Api) verifyEmail(c *gin.Context) {
	decoded, _ := base32.StdEncoding.DecodeString(c.Param("token"))

	manager := security.TokenManager{Signed: string(decoded)}
	user := manager.VerifyToken(_api.DB)

	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verification expired"})
		return
	}

	if !database.SetUserEmailAsVerified(_api.DB, user.Email) {
		c.JSON(http.StatusAlreadyReported, gin.H{"status": "Already verified"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (_api *Api) resetPassword(c *gin.Context) {
	request := database.User{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logger.Panic(err)
	}

	password := database.ResetPassword(_api.DB, request.Email)
	if password == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Email doesn't match any registered account"})
		return
	}

	go external_clients.SendEmail(
		request.Email,
		"Password Recovery",
		fmt.Sprintf("Your new password is %s.<br>If you didn't request password recovery, please contact us.", password),
	)

	c.JSON(http.StatusOK, gin.H{"message": "Recovery email sent"})
}

func (_api *Api) login(c *gin.Context) {
	request := &database.LoginRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logger.Panic(err)
	}

	user := database.Login(_api.DB, request)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	middlewares.GenerateAccessToken(c, user)
	middlewares.GenerateRefreshToken(c, user)
	c.JSON(http.StatusCreated, gin.H{"session": gin.H{"user": user}})
}

func (_api *Api) refresh(c *gin.Context) {
	manager := security.ExtractToken(c, middlewares.RefreshTokenCookieKey)
	user := manager.VerifyToken(_api.DB)

	if manager.Error != nil || user == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	middlewares.GenerateAccessToken(c, user)
	c.JSON(http.StatusCreated, gin.H{"session": gin.H{"user": user}})
}

func (_api *Api) logout(c *gin.Context) {
	c.SetCookie(middlewares.RefreshTokenCookieKey, "", -1, "", "", true, true)
	c.SetCookie(middlewares.AccessTokenCookieKey, "", -1, "", "", true, true)

	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (_api *Api) grpcExample(c *gin.Context) {
	payload := &database.User{}
	if err := c.ShouldBindJSON(payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logger.Panic(err)
	}

	grpcClient := exampleProtoMessages.NewExampleServiceClient(_api.GrpcConnection)

	_, err := grpcClient.HelloWorld(
		context.Background(),
		&exampleProtoMessages.HelloRequest{Email: payload.Email, Password: payload.Password},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		logger.Panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func getUserSession(c *gin.Context) *database.User {
	user, _ := c.Get("user")
	return user.(*database.User)
}
