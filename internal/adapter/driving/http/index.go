package http

import (
	"auth_service/internal/config"
	"auth_service/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Server struct {
	router *gin.Engine
	cfg    *config.Config
	uc     *usecase.UseCases
}

const httpServerReadHeaderTimeout = 70 * time.Second

func New(
	cfg *config.Config,
	uc *usecase.UseCases,
) *http.Server {
	r := gin.New()

	srv := &Server{
		router: r,
		cfg:    cfg,
		uc:     uc,
	}

	srv.endpoints()

	httpServer := &http.Server{
		Addr:              cfg.HTTPPort,
		Handler:           srv,
		ReadHeaderTimeout: httpServerReadHeaderTimeout,
	}

	// srv.log.Info(fmt.Sprintf("HTTP server is initialized on port: %v", cfg.HTTPPort))

	return httpServer
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
