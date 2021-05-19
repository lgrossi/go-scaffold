package api

import (
	"context"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	example_proto_messages "github.com/lgrossi/go-scaffold/src/grpc/example_proto_defs"
	"net/http"
)

type RequestPayload struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
}

func (_api *Api) example(c *gin.Context) {
	payload := &RequestPayload{}
	if err := c.ShouldBindJSON(payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grpcClient := example_proto_messages.NewExampleServiceClient(_api.GrpcConnection)

	_, err := grpcClient.HelloWorld(
		context.Background(),
		&example_proto_messages.HelloRequest{Email: payload.Email, Password: payload.Password},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}
