package response

import (
	"net/http"

	"github.com/akubi0w1/chatbot-202201/code"
	"github.com/akubi0w1/chatbot-202201/log"
	"github.com/go-chi/render"
)

func Success(w http.ResponseWriter, r *http.Request, body interface{}) {
	render.Status(r, http.StatusOK)
	render.JSON(w, r, body)
}

func Error(w http.ResponseWriter, r *http.Request, err error) {
	logger := log.New()

	statusCode := code.GetStatusCode(err)

	if statusCode >= 500 {
		logger.Errorf("%s", err.Error())
	} else {
		logger.Warnf("%s", err.Error())
	}

	render.Status(r, statusCode)
	render.JSON(w, r, errorResponse{
		Code:    string(code.GetCode(err)),
		Message: code.GetError(err).Error(),
	})
}

type errorResponse struct {
	StatusCode int    // http status
	Code       string // error code
	Message    string // message
}
