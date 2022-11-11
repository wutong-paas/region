package webapi

import (
	"log"
	"os"
	"strings"

	"github.com/wutong-paas/region/internal/routes"

	"github.com/gin-gonic/gin"
)

type server struct {
	Router routes.Router
	*gin.Engine
	port     string
	pathBase string
}

func NewServer() *server {
	port := os.Getenv("APP_PORT")
	pathBase := os.Getenv("APP_PATH_BASE")
	if port == "" {
		port = "8081"
	}
	engine := gin.Default()
	engine.SetTrustedProxies([]string{"0.0.0.0"})
	var rootRouterGroup *gin.RouterGroup
	if pathBase != "" && !strings.HasPrefix(pathBase, "/") {
		pathBase = "/" + pathBase
	}
	rootRouterGroup = engine.Group(pathBase).Group("/api")

	return &server{
		//Router: routes.Router{engine, rootRouterGroup},
		Router:   routes.Router{RouterGroup: rootRouterGroup},
		port:     port,
		Engine:   engine,
		pathBase: pathBase,
	}
}

func (srv *server) Start() {
	log.Println("region webapi server started...")
	srv.Router.Register()

	if err := srv.Run(":" + srv.port); err != nil {
		log.Fatal(err)
	}
}
