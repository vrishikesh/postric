package handlers

import (
	"broker/util"
	"log"
	"net/http"

	"github.com/go-chi/render"
)

func HandleSubmission(w http.ResponseWriter, r *http.Request) {
	requestPayload := new(util.RequestPayload)

	if err := render.Bind(r, requestPayload); err != nil {
		render.JSON(w, r, util.ErrInvalidRequest(err))
		return
	}

	log.Print(requestPayload)
}
