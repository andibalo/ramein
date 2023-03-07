package v1

import (
	"github.com/andibalo/ramein/orion/internal/apperr"
	"github.com/andibalo/ramein/orion/internal/config"
	"github.com/andibalo/ramein/orion/internal/constants"
	"github.com/andibalo/ramein/orion/internal/httpresp"
	"github.com/andibalo/ramein/orion/internal/request"
	"github.com/andibalo/ramein/orion/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

	var data request.CreateTemplateReq

	if err := c.BindJSON(&data); err != nil {
		httpresp.HttpRespError(c, apperr.ErrBadRequest)
		return
	}

	err := h.templateService.CreateTemplate(data)

	if err != nil {
		h.cfg.Logger().Error("[CreateTemplate] Error creating template", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
}

func (h *TemplateController) GetTemplateByID(c *gin.Context) {

	template, err := h.templateService.GetTemplateByID(c.Param("id"))

	if err != nil {
		h.cfg.Logger().Error("[CreateTemplate] Error creating template", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, template, nil)
}
