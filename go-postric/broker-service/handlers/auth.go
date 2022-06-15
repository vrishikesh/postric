package handlers

import (
	"broker/util"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

func Authenticate(a util.AuthPayload) render.Renderer {
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	request, err := http.NewRequest(
		"POST",
		"http://auth-service/v1/signin",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return util.ErrBadRequest(err)
	}

	client := new(http.Client)
	resp, err := client.Do(request)
	if err != nil {
		return util.ErrBadRequest(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return util.ErrBadRequest(errors.New("invalid credentials"))
	}

	if resp.StatusCode != http.StatusOK {
		return util.ErrBadRequest(errors.New("error calling auth service"))
	}

	jsonFromService := new(util.JsonResponse)
	err = json.NewDecoder(resp.Body).Decode(jsonFromService)
	if err != nil {
		return util.ErrBadRequest(err)
	}

	if jsonFromService.Error {
		return util.ErrResponse(err, http.StatusUnauthorized)
	}

	return util.Response(jsonFromService.Data)
}
