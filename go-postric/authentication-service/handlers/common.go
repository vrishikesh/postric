package handlers

import (
	"authentication/util"
	"net/http"

	"github.com/go-chi/render"
)

func Home(w http.ResponseWriter, r *http.Request) render.Renderer {
	return util.Response(http.StatusOK, "Hit the auth")
}

func Ping(w http.ResponseWriter, r *http.Request) render.Renderer {
	return util.Response(http.StatusOK, "pong")
}
