package v1

import (
	"github.com/andibalo/ramein/orion/internal/config"
	"github.com/andibalo/ramein/orion/internal/constants"
	"github.com/andibalo/ramein/orion/internal/service"
	"github.com/gin-gonic/gin"
)

const (
	templateBasePath = "/template"
)

type TemplateController struct {
	cfg             config.Config
	templateService service.TemplateService
}

func NewTemplateController(cfg config.Config, templateService service.TemplateService) *TemplateController {

	return &TemplateController{
		cfg:             cfg,
		templateService: templateService,
	}
}

func (h *TemplateController) AddRoutes(r *gin.Engine) {
	tr := r.Group(constants.V1BasePath + templateBasePath)

	tr.POST("/", h.CreateTemplate)
	tr.GET("/:id", h.GetTemplateByID)
}

func (h *TemplateController) CreateTemplate(c *gin.Context) {

}

func (h *TemplateController) GetTemplateByID(c *gin.Context) {

}
