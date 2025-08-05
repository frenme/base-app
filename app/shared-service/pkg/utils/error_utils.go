package utils

import (
	"errors"
	"fmt"
	"net/http"
	"shared/pkg/dto"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

var ErrorNotFound = errors.New("not found")
var ErrorBadRequest = errors.New("bad request")

func HandleQueryRequestData(c *gin.Context, obj any) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid data: "+err.Error())
		return err
	}
	validate := validator.New()
	if err := validate.Struct(obj); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid data: "+err.Error())
		return err
	}
	return nil
}

func HandleBodyRequestData(c *gin.Context, obj any) error {
	if err := c.ShouldBindJSON(&obj); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid data: "+err.Error())
		return err
	}
	validate := validator.New()
	if err := validate.Struct(obj); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid data: "+err.Error())
		return err
	}
	return nil
}

func HandleError(c *gin.Context, err error) error {
	if err == nil {
		return nil
	}

	var errorStatus *ErrorStatus
	switch {
	case errors.As(err, &errorStatus) && errors.Is(errorStatus.Base, ErrorBadRequest):
		RespondWithError(c, http.StatusBadRequest, err.Error())
		return err
	case errors.As(err, &errorStatus) && errors.Is(errorStatus.Base, ErrorNotFound):
		RespondWithError(c, http.StatusNotFound, err.Error())
		return err
	default:
		RespondWithError(c, http.StatusInternalServerError, err.Error())
		return err
	}
}

func RespondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, dto.ErrorResponse{
		StatusCode: code,
		Message:    message,
	})
}

type ErrorStatus struct {
	Base error
	Msg  string
}

func (e *ErrorStatus) Error() string {
	return fmt.Sprintf("%s: %s", e.Base.Error(), e.Msg)
}

func (e *ErrorStatus) Unwrap() error {
	return e.Base
}
