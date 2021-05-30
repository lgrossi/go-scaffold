package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lgrossi/go-scaffold/src/api/external_clients"
	"github.com/lgrossi/go-scaffold/src/api/middlewares"
	"github.com/lgrossi/go-scaffold/src/database"
	exampleProtoMessages "github.com/lgrossi/go-scaffold/src/grpc/example_proto_defs"
	"github.com/lgrossi/go-scaffold/src/logger"
	"net/http"
)

func (_api *Api) login(c *gin.Context) {
	request := &database.AuthRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logger.Panic(err)
	}

	user := database.Login(_api.DB, request)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	middlewares.GenerateTokens(c, user)
	c.JSON(http.StatusCreated, gin.H{"session": gin.H{"user": user}})
}

func (_api *Api) register(c *gin.Context) {
	user := &database.User{}
	if err := c.ShouldBindJSON(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logger.Panic(err)
	}

	user = database.CreateUser(_api.DB, user)
	if user == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "This email is already in use."})
		return
	}

	middlewares.GenerateTokens(c, user)
	c.JSON(http.StatusCreated, gin.H{"session": gin.H{"user": user}})
}

func (_api *Api) refresh(c *gin.Context) {
	middlewares.RefreshToken(_api.DB, c)
}

func (_api *Api) resetPassword(c *gin.Context) {
	request := &struct{ Email string }{}
	if err := c.ShouldBindJSON(request); err != nil {
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

func (_api *Api) logout(c *gin.Context) {
	c.SetCookie(middlewares.RefreshTokenCookieKey, "", -1, "", "", true, true)
	c.SetCookie(middlewares.AccessTokenCookieKey, "", -1, "", "", true, true)

	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (_api *Api) protectedExample(c *gin.Context) {
	session := getSession(c)

	c.JSON(
		http.StatusOK,
		gin.H{
			"status": fmt.Sprintf(
				"Logged in as %s",
				session.Email,
			),
		},
	)
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

func getSession(c *gin.Context) *middlewares.TokenSession {
	session, _ := c.Get("session")
	return session.(*middlewares.TokenSession)
}
