package response

import (
	"github.com/gin-gonic/gin"
	"github.com/rianekacahya/errors"
	"net/http"
)

type response struct {
	Message interface{} `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Error(c *gin.Context, err error) {
	var(
		status int
		response = new(response)
	)

	if err != nil {
		// Mapping Status
		errorStatus := errors.GetStatus(err)
		errorMessage := errors.GetError(err)
		switch errorStatus {
		case errors.GENERIC:
			status = http.StatusInternalServerError
		case errors.FORBIDDEN:
			status = http.StatusForbidden
		case errors.BADREQUEST:
			status = http.StatusBadRequest
		case errors.NOTFOUND:
			status = http.StatusNotFound
		case errors.UNAUTHORIZED:
			status = http.StatusUnauthorized
		default:
			status = http.StatusInternalServerError
			response.Message = err.Error()
		}

		if errorStatus != errors.NOTYPE {
			if e, ok := err.(error); ok {
				response.Message = e.Error()
			}else{
				response.Message = errorMessage
			}
		}
	}

	c.AbortWithStatusJSON(status, response)
}

func Render(c *gin.Context, status int, data interface{}) {
	var response = new(response)

	response.Message = "success"
	response.Data = data

	c.JSON(http.StatusOK, response)
}
