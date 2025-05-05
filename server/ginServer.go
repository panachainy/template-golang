package server

import (
	"fmt"
	"net/http"
	"template-golang/config"
	"template-golang/modules/cockroach"

	docs "template-golang/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

const (
	apiV1Path = "/api/v1"
)

type Handlers struct {
	cockroach *cockroach.Cockroach
}

type ginServer struct {
	router   *gin.Engine
	conf     *config.Config
	handlers Handlers
}

func NewGinServer(conf *config.Config, cockroach *cockroach.Cockroach) *ginServer {
	r := gin.Default()

	return &ginServer{
		router: r,
		conf:   conf,
		handlers: Handlers{
			cockroach: cockroach,
		},
	}
}

func (s *ginServer) Start() {
	docs.SwaggerInfo.BasePath = apiV1Path

	v1 := s.router.Group(apiV1Path)

	v1.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	s.initializeCockroachHttpHandler()
	if gin.Mode() == gin.DebugMode {
		s.initSwagger()
	}

	serverUrl := fmt.Sprintf(":%d", s.conf.Server.Port)

	s.router.Run(serverUrl)
}

func (s *ginServer) initSwagger() {

	ginSwagger.WrapHandler(swaggerfiles.Handler,
		ginSwagger.URL("http://localhost:8080/swagger/doc.json"),
		ginSwagger.DefaultModelsExpandDepth(-1))

	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// swagger info
	fmt.Println()
	fmt.Printf("API URL: http://localhost:%d%s\n", s.conf.Server.Port, apiV1Path)

	fmt.Printf("Swagger UI URL: http://localhost:%d/swagger/index.html\n", s.conf.Server.Port)
	fmt.Printf("Swagger JSON URL: http://localhost:%d/swagger/doc.json\n", s.conf.Server.Port)
	fmt.Println()
}

func (s *ginServer) initializeCockroachHttpHandler() {
	v1 := s.router.Group(apiV1Path)
	cockroachRouters := v1.Group("/cockroach")
	cockroachRouters.POST("", s.handlers.cockroach.Handler.DetectCockroach)
}
