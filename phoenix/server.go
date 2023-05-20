package phoenix

import (
	"context"
	"github.com/andibalo/ramein/phoenix/internal/api"
	v1 "github.com/andibalo/ramein/phoenix/internal/api/v1"
	"github.com/andibalo/ramein/phoenix/internal/config"
	"github.com/andibalo/ramein/phoenix/internal/repository"
	"github.com/andibalo/ramein/phoenix/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Server struct {
	gin *gin.Engine
}

func NewServer(ctx context.Context, cfg config.Config, db neo4j.DriverWithContext) *Server {

	router := gin.Default()

	userRepo := repository.NewUserRepo(ctx, cfg, db)

	userService := service.NewUserService(cfg, userRepo)

	userController := v1.NewUserController(cfg, userService)

	registerHandlers(router, &api.HealthCheck{}, userController)

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
