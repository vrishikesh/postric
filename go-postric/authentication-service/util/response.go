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

func Write(err error, statusCode int, data any) render.Renderer {
	isError, message := false, ""
	if err != nil {
		isError = true
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

func ErrResponse(err error, statusCode ...int) render.Renderer {
	code := http.StatusBadRequest

	if len(statusCode) > 0 {
		code = statusCode[0]
	}

	return Write(err, code, nil)

	// return &JsonResponse{
	// 	Err:            err,
	// 	HTTPStatusCode: code,

	// 	StatusText: http.StatusText(code),
	// 	Error:      true,
	// 	Message:    err.Error(),
	// }
}

func ErrBadRequest(err error) render.Renderer {
	return ErrResponse(err, http.StatusBadRequest)
}

func ErrNotFound(err error) render.Renderer {
	return ErrResponse(err, http.StatusNotFound)
}

func ErrUnauthorized(err error) render.Renderer {
	return ErrResponse(err, http.StatusUnauthorized)
}

func ErrInternalServer(err error) render.Renderer {
	return ErrResponse(err, http.StatusInternalServerError)
}

func Response(data any, statusCode ...int) render.Renderer {
	code := http.StatusOK

	if len(statusCode) > 0 {
		code = statusCode[0]
	}

	return Write(nil, code, data)

	// return &JsonResponse{
	// 	HTTPStatusCode: code,

	// 	StatusText: http.StatusText(code),
	// 	Error:      false,
	// 	Data:       data,
	// }
}

func RouteHandler(callback func(http.ResponseWriter, *http.Request) render.Renderer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, callback(w, r))
	}
}
