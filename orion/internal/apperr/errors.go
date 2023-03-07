package apperr

import (
	"errors"
	"github.com/andibalo/ramein/orion/internal/response"

	"fmt"
	"net/http"
)

var (
	ErrBadRequest          = errors.New("request is incomplete or invalid")
	ErrUnprocessable       = fmt.Errorf("request is complete but invalid - %w", ErrBadRequest)
	ErrInvalidParam        = errors.New("invalid input param(s)")
	ErrUnauthorized        = errors.New("authorization missing or bad")
	ErrPermissionDenied    = errors.New("permission denied")
	ErrNotFound            = errors.New("record not found")
	ErrDuplicateUser       = errors.New("user already exists")
	ErrInternalServerError = errors.New("internal server error")
	ErrDuplicateRequest    = errors.New("duplicate request")
	ErrTimeout             = errors.New("timeout")
	ErrConnectTimeout      = errors.New("connect timeout")
	ErrRateLimited         = errors.New("rate limited")
	ErrUserError           = errors.New("user error")
	ErrUnexpectedError     = errors.New("unexpected error")
	ErrTransientError      = errors.New("transient application error")
)

func MapResponseCodeToErrors(code response.Code) error {
	switch code {
	case response.Success:
		return nil

	case response.BadRequest:
		return ErrBadRequest

	case response.Unauthorized:
		return ErrUnauthorized

	case response.Forbidden:
		return ErrPermissionDenied

	case response.ServerError:
		return ErrTransientError

	case response.GatewayTimeout:
		return ErrTimeout

	default:
		return ErrUnexpectedError
	}
}

func MapStatusCodeToErrors(httpCode int) error {
	switch httpCode {
	case http.StatusCreated:
		return nil

	case http.StatusPaymentRequired:
		return ErrUserError

	case http.StatusConflict:
		return ErrDuplicateRequest

	case http.StatusUnprocessableEntity, http.StatusBadRequest:
		return ErrBadRequest

	case http.StatusUnauthorized:
		return ErrUnauthorized

	case http.StatusForbidden:
		return ErrPermissionDenied

	case http.StatusInternalServerError, http.StatusServiceUnavailable, http.StatusBadGateway:
		return ErrTransientError

	case http.StatusGatewayTimeout:
		return ErrTimeout

	case http.StatusTooManyRequests:
		return ErrRateLimited

	default:
		return ErrUnexpectedError
	}
}

func MapErrorsToCode(err error) response.Code {
	switch {
	case errors.Is(err, ErrUnauthorized):
		return response.Unauthorized

	case errors.Is(err, ErrNotFound):
		return response.NotFound

	case errors.Is(err, ErrTimeout):
		return response.GatewayTimeout

	case errors.Is(err, ErrBadRequest):
		return response.BadRequest

	default:
		return response.ServerError
	}
}

func MapErrorsToStatusCode(err error) int {
	switch {
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized

	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound

	case errors.Is(err, ErrTimeout):
		return http.StatusGatewayTimeout

	case errors.Is(err, ErrBadRequest), errors.Is(err, ErrDuplicateUser):
		return http.StatusBadRequest

	default:
		return http.StatusInternalServerError
	}
}
