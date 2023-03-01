package httpresp

import (
	"errors"
	"github.com/andibalo/ramein/core/internal/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"net/http"
	"reflect"
	"time"
)

type Meta struct {
	Path       string `json:"path"`
	StatusCode int    `json:"statusCode"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	Error      string `json:"error" swaggerignore:"true"`
	Timestamp  string `json:"timestamp"`
}

type Response struct {
	Data       interface{} `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Success    string      `json:"success"`
}

type Pagination struct {
	CurrentPage     int64   `json:"current_page"`
	CurrentElements int64   `json:"current_elements"`
	TotalPages      int64   `json:"total_pages"`
	TotalElements   int64   `json:"total_elements"`
	SortBy          string  `json:"sort_by"`
	CursorStart     *string `json:"cursor_start,omitempty"`
	CursorEnd       *string `json:"cursor_end,omitempty"`
}

// HTTPErrResp http error response
type HTTPErrResp struct {
	Meta Meta `json:"metadata"`
}

func HttpRespError(c *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	// Retrieve the custom status code if it's a *fiber.Error
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	jsonErrResp := &HTTPErrResp{
		Meta: Meta{
			Path:       c.Path(),
			StatusCode: code,
			Status:     http.StatusText(code),
			Message:    e.Message,
			Error:      e.Error(),
			Timestamp:  time.Now().Format(time.RFC3339),
		},
	}

	l := logger.InitLogger()

	l.Error("An Error Occured", zap.Any("Error :", e))

	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	return c.Status(code).JSON(jsonErrResp)
}

func HttpRespSuccess(c *fiber.Ctx, data interface{}, pagination *Pagination) error {

	//check typenya slice / array , soalnya kalo bukan slice / array ga perlu di dikosongin datanya soalnya udah kena error di handler
	kind := reflect.ValueOf(data).Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		//check kalo data nya nil / kosong
		if data == nil || reflect.ValueOf(data).IsNil() {
			//kalo data arraynya kosong returnnya "data": []
			data = []interface{}{}
		}
	}

	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	return c.Status(fiber.StatusOK).JSON(Response{
		Data:       data,
		Pagination: pagination,
		Success:    "success",
	})
}

func ResetPagination() *Pagination {

	return &Pagination{
		CurrentPage:     1,
		CurrentElements: 0,
		TotalPages:      0,
		TotalElements:   0,
	}
}
