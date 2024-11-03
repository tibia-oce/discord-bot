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

	/* Generate HTTP/GRPC reverse proxy */
	var err error
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

	// Group creates a new router group. You should add all the routes that have common middlewares or the same path prefix.
	// For example, all the routes that use a common middleware for authorization could be grouped.
	authorized := _api.Router.Group("/")
	authorized.Use(middlewares.VerifyTokenHandler(_api.DB))
	{
		authorized.POST("/grpc", _api.grpcExample)
	}

}
