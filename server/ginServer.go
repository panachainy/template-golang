package server

import (
	"fmt"
	"net/http"
	"template-golang/config"
	"template-golang/database"

	cockroachHandlers "template-golang/modules/cockroach/handlers"
	cockroachRepositories "template-golang/modules/cockroach/repositories"
	cockroachUsecases "template-golang/modules/cockroach/usecases"

	docs "template-golang/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

const (
	apiV1Path = "/api/v1"
)

type ginServer struct {
	router *gin.Engine
	db     database.Database
	conf   *config.Config
}

func NewGinServer(conf *config.Config, db database.Database) Server {
	r := gin.Default()

	return &ginServer{
		router: r,
		db:     db,
		conf:   conf,
	}
}

func (s *ginServer) Start() {
	docs.SwaggerInfo.BasePath = apiV1Path

	v1 := s.router.Group(apiV1Path)

	v1.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	s.initializeCockroachHttpHandler()
	s.initSwagger()
	s.router.GET("/example/helloworld", Helloworld)

	serverUrl := fmt.Sprintf(":%d", s.conf.Server.Port)

	fmt.Printf("Mode: %v\n", gin.Mode())

	if gin.Mode() == gin.DebugMode {
		fmt.Println()
		fmt.Printf("API URL: http://localhost:%d%s\n", s.conf.Server.Port, apiV1Path)

		fmt.Printf("Swagger UI URL: http://localhost:%d/swagger/index.html\n", s.conf.Server.Port)
		fmt.Printf("Swagger JSON URL: http://localhost:%d/swagger/doc.json\n", s.conf.Server.Port)
		fmt.Println()
	}

	s.router.Run(serverUrl)
}

func (s *ginServer) initSwagger() {
	ginSwagger.WrapHandler(swaggerfiles.Handler,
		ginSwagger.URL("http://localhost:8080/swagger/doc.json"),
		ginSwagger.DefaultModelsExpandDepth(-1))

	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

func (s *ginServer) initializeCockroachHttpHandler() {

	cockroachPostgresRepository := cockroachRepositories.NewCockroachPostgresRepository(s.db)
	cockroachFCMMessaging := cockroachRepositories.NewCockroachFCMMessaging()
	cockroachUsecase := cockroachUsecases.NewCockroachUsecaseImpl(
		cockroachPostgresRepository,
		cockroachFCMMessaging,
	)
	cockroachHttpHandler := cockroachHandlers.NewCockroachHttpHandler(cockroachUsecase)

	v1 := s.router.Group(apiV1Path)
	cockroachRouters := v1.Group("/cockroach")
	cockroachRouters.POST("", cockroachHttpHandler.DetectCockroach)
}
