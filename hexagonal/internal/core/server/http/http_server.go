package http

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"user-service/internal/core/config"
)

const defaultHost = "0.0.0.0"

type HttpServer interface {
	Start()
	io.Closer
}

type httpServer struct {
	Port   uint
	server *http.Server
}

func NewHttpServer(router *gin.Engine, config config.HttpServerConfig) HttpServer {
	return &httpServer{
		Port: config.Port,
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", defaultHost, config.Port),
			Handler: router,
		},
	}
}

func (httpServer httpServer) Start() {
	go func() {
		if err := httpServer.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf(
				"failed to stater HttpServer listen port %d, err=%s\n",
				httpServer.Port, err.Error(),
			)
		}
	}()
	log.Printf("Start Service with port %d\n", httpServer.Port)
}

func (httpServer httpServer) Close() error {
	ctx, cancel := context.WithTimeout(
		context.Background(), time.Duration(3)*time.Second,
	)
	defer cancel()

	return httpServer.server.Shutdown(ctx)
}
