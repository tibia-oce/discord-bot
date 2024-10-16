package api

import (
	"database/sql"
	"errors"
	"net/http"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tibia-oce/discord-bot/src/api/middlewares"
	"github.com/tibia-oce/discord-bot/src/configs"
	"github.com/tibia-oce/discord-bot/src/database"
	"github.com/tibia-oce/discord-bot/src/logger"
	"github.com/tibia-oce/discord-bot/src/network"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Api struct {
	Router         *gin.Engine
	DB             *sql.DB
	GrpcConnection *grpc.ClientConn
	network.ServerInterface
}

func (_api *Api) Initialize(gConfigs configs.GlobalConfigs) error {
	_api.DB = database.PullConnection(gConfigs)

	ipLimiter := &middlewares.IPRateLimiter{
		Visitors: make(map[string]*middlewares.Visitor),
		Mu:       &sync.RWMutex{},
	}

	ipLimiter.Init()

	gin.SetMode(gin.ReleaseMode)

	_api.Router = gin.New()
	_api.Router.Use(logger.LogRequest())
	_api.Router.Use(gin.Recovery())
	_api.Router.Use(ipLimiter.Limit())

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://127.0.0.1:3000", "http://localhost:3000"}
	config.AllowCredentials = true
	_api.Router.Use(cors.New(config))
	_api.Router.HandleMethodNotAllowed = true

	_api.initializeRoutes()

	var err error
	/* Generate HTTP/GRPC reverse proxy */
	_api.GrpcConnection, err = grpc.NewClient(gConfigs.ServerConfigs.Grpc.Format(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error(errors.New("couldn't start GRPC reverse proxy"))
		return err
	}

	return nil
}

func (_api *Api) Run(gConfigs configs.GlobalConfigs) error {
	err := http.ListenAndServe(gConfigs.ServerConfigs.Http.Format(), _api.Router)

	/* Make sure we free the reverse proxy connection */
	if _api.GrpcConnection != nil {
		closeErr := _api.GrpcConnection.Close()
		if closeErr != nil {
			logger.Error(closeErr)
		}
	}

	return err
}

func (_api *Api) GetName() string {
	return "api"
}

func (_api *Api) initializeRoutes() {
	_api.Router.POST("/user/login", _api.login)
	_api.Router.POST("/user/resetPassword", _api.resetPassword)
	_api.Router.POST("/user/register", _api.register)
	_api.Router.GET("/user/refresh", _api.refresh)
	_api.Router.GET("/user/verification/:token", _api.verifyEmail)

	authorized := _api.Router.Group("/")
	authorized.Use(middlewares.VerifyTokenHandler(_api.DB))
	{
		authorized.POST("/grpc", _api.grpcExample)
		authorized.GET("/user/logout", _api.logout)
	}

	authorized.Use(isVerified)
	{
		authorized.POST("/grpc2", _api.grpcExample)
	}
}

func isVerified(c *gin.Context) {
	user := getUserSession(c)

	if !user.Verified {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Email address not yet verified"})
	}
}
