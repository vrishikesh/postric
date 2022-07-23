package util

import (
	"net/http"

	"github.com/go-chi/render"
)

type JsonResponse struct {
	Err            error `json:"-"`
	HTTPStatusCode int   `json:"-"`

	StatusText string `json:"status"`
	AppCode    uint64 `json:"code,omitempty"`
	Error      bool   `json:"error"`
	Message    string `json:"message,omitempty"`
	Data       any    `json:"data,omitempty"`
}

func (e *JsonResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func Write(isError bool, statusCode int, err error, data any) render.Renderer {
	var message string
	if err != nil {
		message = err.Error()
	}

	return &JsonResponse{
		Err:            err,
		HTTPStatusCode: statusCode,

		StatusText: http.StatusText(statusCode),
		Error:      isError,
		Message:    message,
		Data:       data,
	}
}

func ErrResponse(statusCode int, err error) render.Renderer {
	return Write(true, statusCode, err, nil)
}

func Response(statusCode int, data any) render.Renderer {
	return Write(false, statusCode, nil, data)
}

func RouteHandler(callback func(http.ResponseWriter, *http.Request) render.Renderer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, callback(w, r))
	}
}
