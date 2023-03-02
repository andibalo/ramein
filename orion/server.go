package orion

import (
	"entgo.io/ent/entc/integration/ent"
	"github.com/andibalo/ramein/orion/internal/api"
	"github.com/andibalo/ramein/orion/internal/config"
	"github.com/gin-gonic/gin"
)

type Server struct {
	gin *gin.Engine
}

func NewServer(cfg config.Config, db *ent.Client) *Server {

	router := gin.Default()

	registerHandlers(router, &api.HealthCheck{})

	return &Server{
		gin: router,
	}
}

func (s *Server) Start(addr string) error {
	return s.gin.Run(addr)
}

func registerHandlers(g *gin.Engine, handlers ...api.Handler) {
	for _, handler := range handlers {
		handler.AddRoutes(g)
	}
}
