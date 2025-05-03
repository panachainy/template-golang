package server

import (
	"fmt"
	"net/http"
	"template-golang/config"
	"template-golang/database"

	cockroachHandlers "template-golang/modules/cockroach/handlers"
	cockroachRepositories "template-golang/modules/cockroach/repositories"
	cockroachUsecases "template-golang/modules/cockroach/usecases"

	"github.com/gin-gonic/gin"
)

type ginServer struct {
	router *gin.Engine
	db     database.Database
	conf   *config.Config
}

func NewGinServer(conf *config.Config, db database.Database) Server {
	router := gin.Default()
	return &ginServer{
		router: router,
		db:     db,
		conf:   conf,
	}
}

func (s *ginServer) Start() {
	s.router.GET("/v1/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	s.initializeCockroachHttpHandler()

	serverUrl := fmt.Sprintf(":%d", s.conf.Server.Port)
	s.router.Run(serverUrl)
}

func (s *ginServer) initializeCockroachHttpHandler() {
	cockroachPostgresRepository := cockroachRepositories.NewCockroachPostgresRepository(s.db)
	cockroachFCMMessaging := cockroachRepositories.NewCockroachFCMMessaging()
	cockroachUsecase := cockroachUsecases.NewCockroachUsecaseImpl(
		cockroachPostgresRepository,
		cockroachFCMMessaging,
	)
	cockroachHttpHandler := cockroachHandlers.NewCockroachHttpHandler(cockroachUsecase)

	cRouters := s.router.Group("/v1/cockroach")
	cRouters.POST("", cockroachHttpHandler.DetectCockroach)
}
