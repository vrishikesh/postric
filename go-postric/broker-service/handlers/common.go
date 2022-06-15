package handlers

import (
	"broker/util"
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

func Home(w http.ResponseWriter, r *http.Request) render.Renderer {
	return util.Response("Hit the broker")
}

func Ping(w http.ResponseWriter, r *http.Request) render.Renderer {
	return util.Response("pong")
}

func HandleSubmission(w http.ResponseWriter, r *http.Request) render.Renderer {
	requestPayload := new(util.RequestPayload)

	if err := render.Bind(r, requestPayload); err != nil {
		return util.ErrBadRequest(err)
	}

	switch requestPayload.Action {
	case "auth":
		return Authenticate(requestPayload.Auth)
	default:
		return util.ErrNotFound(errors.New("unknown action"))
	}
}
