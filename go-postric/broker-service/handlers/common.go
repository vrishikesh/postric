package handlers

import (
	"broker/util"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-chi/render"
)

var templateUrl = "http://%s-service/%s"

func Home(w http.ResponseWriter, r *http.Request) render.Renderer {
	return util.Response(http.StatusOK, "Hit the broker")
}

func Ping(w http.ResponseWriter, r *http.Request) render.Renderer {
	return util.Response(http.StatusOK, "pong")
}

func Handle(w http.ResponseWriter, r *http.Request) render.Renderer {
	// Get the service name
	re := regexp.MustCompile(`(?i)\/[^\/]+`)
	service := re.FindString(r.URL.String())
	service = strings.Trim(service, "/")

	// Read the whole json body
	rBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("could not read the body: %v", err)
		return util.ErrResponse(http.StatusBadRequest, err)
	}

	// Get the service path
	path := strings.Replace(r.URL.String(), service, "", 1)
	path = strings.Trim(path, "/")

	return ForwardRequest(r.Method, service, path, rBody)
}

func ForwardRequest(method, service, path string, body ...[]byte) render.Renderer {
	log.Printf("calling service: %s", fmt.Sprintf(templateUrl, service, path))

	buffer := new(bytes.Buffer)
	if len(body) > 0 {
		buffer = bytes.NewBuffer(body[0])
	}

	request, err := http.NewRequest(
		method,
		fmt.Sprintf(templateUrl, service, path),
		buffer,
	)
	if err != nil {
		log.Printf("could not create new request: %v", err)
		return util.ErrResponse(http.StatusBadRequest, err)
	}

	client := new(http.Client)
	resp, err := client.Do(request)
	if err != nil {
		log.Printf("error in service response: %v", err)
		return util.ErrResponse(http.StatusBadRequest, err)
	}
	defer resp.Body.Close()

	jsonFromService := new(util.JsonResponse)
	err = json.NewDecoder(resp.Body).Decode(jsonFromService)
	if err != nil {
		log.Printf("could not parse json: %v", err)
		return util.ErrResponse(http.StatusUnprocessableEntity, err)
	}

	if jsonFromService.Error {
		log.Printf("handled error received from service: %v", jsonFromService.Message)
		return util.ErrResponse(resp.StatusCode, errors.New(jsonFromService.Message))
	}

	log.Printf("call to service successful with status code %d and data: %+v", resp.StatusCode, jsonFromService)
	return util.Response(resp.StatusCode, jsonFromService.Data)
}
