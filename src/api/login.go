package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tibia-oce/discord-bot/src/api/middlewares"
	"github.com/tibia-oce/discord-bot/src/database"
	exampleProtoMessages "github.com/tibia-oce/discord-bot/src/grpc/example_proto_defs"
	"github.com/tibia-oce/discord-bot/src/logger"
)

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
