package orion

import (
	"entgo.io/ent/entc/integration/ent"
	"github.com/andibalo/ramein/commons/rabbitmq"
	"github.com/andibalo/ramein/orion/internal/api"
	v1 "github.com/andibalo/ramein/orion/internal/api/v1"
	"github.com/andibalo/ramein/orion/internal/config"
	"github.com/andibalo/ramein/orion/internal/pubsub"
	"github.com/andibalo/ramein/orion/internal/repository"
	"github.com/andibalo/ramein/orion/internal/service"
	"github.com/gin-gonic/gin"
)

type Server struct {
	gin *gin.Engine
}

func NewServer(cfg config.Config, db *ent.Client) *Server {

	router := gin.Default()

	rmq := rabbitmq.NewRabitmq(rabbitmq.RabitmqConfiguration{
		URL:    cfg.RabbitMQURL(),
		Enable: true,
	})

	pb := pubsub.NewPubSub(cfg, rmq)

	pb.InitSubscribers()

	templateRepo := repository.NewTemplateRepository(db)

	templateService := service.NewTemplateService(templateRepo)

	templateController := v1.NewTemplateController(cfg, templateService)

	registerHandlers(router, &api.HealthCheck{}, templateController)

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
