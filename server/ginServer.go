package server

import (
	"fmt"
	"net/http"
	"template-golang/config"
	"template-golang/modules/auth"
	"template-golang/modules/cockroach"
	"time"

	docs "template-golang/docs"

	"github.com/gin-contrib/cors"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

const (
	apiV1Path = "/api/v1"
)

type Modules struct {
	cockroach *cockroach.Cockroach
	auth      *auth.Auth
}

type ginServer struct {
	router  *gin.Engine
	conf    *config.Config
	modules Modules
}

func Provide(
	conf *config.Config,
	cockroach *cockroach.Cockroach,
	auth *auth.Auth,
) *ginServer {

	config := cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})

	r := gin.Default()

	r.Use(config)

	return &ginServer{
		router: r,
		conf:   conf,
		modules: Modules{
			cockroach: cockroach,
			auth:      auth,
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

	s.modules.auth.Handler.Routes(v1)

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
	cockroachRouters.POST("", s.modules.cockroach.Handler.DetectCockroach)
}
