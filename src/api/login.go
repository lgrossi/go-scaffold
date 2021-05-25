package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lgrossi/go-scaffold/src/api/middlewares"
	"github.com/lgrossi/go-scaffold/src/database"
	exampleProtoMessages "github.com/lgrossi/go-scaffold/src/grpc/example_proto_defs"
	"github.com/lgrossi/go-scaffold/src/logger"
	"net/http"
)

func (_api *Api) login(c *gin.Context) {
	user := &database.User{}
	if err := c.ShouldBindJSON(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logger.Panic(err)
	}

	if !database.Login(_api.DB, user) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	middlewares.GenerateAccessToken(_api.DB, c, user.Email)
	c.JSON(http.StatusCreated, gin.H{"session": gin.H{"user": user}})
}

func (_api *Api) logout(c *gin.Context) {
	session := getSession(c)
	tokenStr, _ := c.Get("jwt")

	database.DeactivateToken(_api.DB, database.ActiveToken{TokenStr: tokenStr.(string), Email: session.Email})

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
