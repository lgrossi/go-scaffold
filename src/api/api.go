package api

import (
	"database/sql"
	"errors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lgrossi/go-scaffold/src/api/limiter"
	"github.com/lgrossi/go-scaffold/src/api/middlewares"
	"github.com/lgrossi/go-scaffold/src/configs"
	"github.com/lgrossi/go-scaffold/src/database"
	"github.com/lgrossi/go-scaffold/src/logger"
	"github.com/lgrossi/go-scaffold/src/network"
	"google.golang.org/grpc"
	"net/http"
	"sync"
)

type Api struct {
	Router         *gin.Engine
	DB             *sql.DB
	GrpcConnection *grpc.ClientConn
	network.ServerInterface
}

func (_api *Api) Initialize(gConfigs configs.GlobalConfigs) error {
	_api.DB = database.PullConnection(gConfigs)

	ipLimiter := &limiter.IPRateLimiter{
		Visitors: make(map[string]*limiter.Visitor),
		Mu:       &sync.RWMutex{},
	}

	ipLimiter.Init()

	gin.SetMode(gin.ReleaseMode)

	_api.Router = gin.New()
	_api.Router.Use(logger.LogRequest())
	_api.Router.Use(gin.Recovery())
	_api.Router.Use(ipLimiter.Limit())

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://127.0.0.1:8080"}
	config.AllowCredentials = true
	_api.Router.Use(cors.New(config))

	_api.initializeRoutes()

	var err error
	/* Generate HTTP/GRPC reverse proxy */
	_api.GrpcConnection, err = grpc.Dial(gConfigs.ServerConfigs.Grpc.Format(), grpc.WithInsecure())
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
	_api.Router.POST("/login", _api.login)

	authorized := _api.Router.Group("/")
	authorized.Use(middlewares.VerifyToken(_api.DB))
	{
		authorized.POST("/protected", _api.protectedExample)
		authorized.GET("/logout", _api.logout)
		authorized.POST("/grpc", _api.grpcExample)
	}
}
