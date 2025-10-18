package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) endpoints() {
	s.router.GET("/ping", s.Ping)

	authG := s.router.Group("/auth")
	authG.POST("/sign-up", s.SignUp)
	authG.POST("/sign-in", s.SignIn)
}

func (s *Server) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"ping": "pong",
	})
}
