package errors

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/alessandra1408/goqrlog/pkg/log"
	"github.com/alessandra1408/goqrlog/pkg/model"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var (
	ErrInternalServer = errors.New("internal server error")
	ErrBadRequest     = errors.New("invalid request")
	ErrNotFound       = errors.New("resource not found")
	ErrUnauthorized   = errors.New("not authorized")
	ErrForbidden      = errors.New("forbidden")

	ErrValidatorRequired  = errors.New("field is required")
	ErrValidatorMaxLenght = errors.New("field exceeded the max length")

	ErrBindingContentType = errors.New("unexpected content-type header")
	ErrBindingRequest     = errors.New("unable to process request body")

	ErrBindingMessage    = "error binding request params: %v"
	ErrValidatingMessage = "error validating request params: %#v"

	ErrRequestCanceledByClient = errors.New("unable to process: request canceled")
	ErrRequestTimedOut         = errors.New("unable to process: request timeout")
)

func ValidatorErrors(err error) []model.ValidationError {
	errorsArr := []model.ValidationError{}

	errs := err.(validator.ValidationErrors)

	for _, e := range errs {
		errorsArr = append(errorsArr, model.ValidationError{
			Field:  e.Field(),
			Reason: validatorErrorTranslation(e.Tag()).Error(),
		})
	}

	return errorsArr
}

func BindingError(expectedContentType string, r *http.Request) error {
	ctype := r.Header.Get("Content-Type")

	if expectedContentType != ctype {
		return ErrBindingContentType
	}

	return ErrBindingRequest
}

func validatorErrorTranslation(tag string) error {
	switch tag {
	case "required":
		return ErrValidatorRequired
	case "max":
		return ErrValidatorMaxLenght
	default:
		return ErrInternalServer
	}
}

func ContextCanceledError(err error) bool {
	return strings.Contains(err.Error(), context.Canceled.Error())
}

func ContextDeadlineError(err error) bool {
	return strings.Contains(err.Error(), context.DeadlineExceeded.Error())
}

func BindingErrorResponse(c echo.Context, handlerLog log.Log, err error) error {
	handlerLog.Warnf(ErrBindingMessage, err)
	return c.JSON(http.StatusBadRequest, model.ErrorResponse{
		Message: BindingError(echo.MIMEApplicationForm, c.Request()).Error(),
	})
}

func ValidationResponse(c echo.Context, handlerLog log.Log, err error) error {
	l := log.LogWithUserAgent(handlerLog, c.Request().UserAgent())

	errs := ValidatorErrors(err)
	l.Warnf(ErrValidatingMessage, errs)

	return c.JSON(http.StatusBadRequest, model.ValidationResponse{
		Message: ErrBadRequest.Error(),
		Errors:  errs,
	})
}
