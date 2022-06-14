package handlers

import (
	"broker/util"
	"net/http"

	"github.com/go-chi/render"
)

func Broker(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, util.JsonResponse{
		Error:   false,
		Message: "Hit the broker",
	})
}

func Ping(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, util.JsonResponse{
		Error:   false,
		Message: "pong",
	})
}
