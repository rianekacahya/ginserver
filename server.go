package ginserver

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"

	"github.com/rianekacahya/config"
	"github.com/rianekacahya/ginserver/middleware"
)

var (
	server *gin.Engine
	mutex  sync.Once
)

func GetServer() *gin.Engine {
	mutex.Do(func() {
		server = newServer()
	})
	return server
}

func newServer() *gin.Engine {
	return gin.New()
}

func InitServer() {

	// Set debug status parameter
	if config.GetGinServerDebug() {
		gin.SetMode(gin.ReleaseMode)
	}

	// init default middleware
	GetServer().Use(gin.Recovery())
	GetServer().Use(cors.Default())
	GetServer().Use(middleware.Recovery())
	GetServer().Use(middleware.Headers())
	GetServer().Use(middleware.Logger())

	// healthCheck endpoint
	GetServer().GET("/infrastructure/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, "OK")
	})
}

func StartServer(ctx context.Context) {
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", config.GetGinServerPort()),
		Handler: GetServer(),
		ReadTimeout:  time.Duration(config.GetHTTPServerReadTimeout()) * time.Second,
		WriteTimeout: time.Duration(config.GetHTTPServerWriteTimeout()) * time.Second,
		IdleTimeout:  time.Duration(config.GetHTTPServerIdleTimeout()) * time.Second,
	}

	select {
	case <-ctx.Done():
		if err := srv.Shutdown(ctx); err != nil {
			panic(err)
		}
	default:
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}
}
