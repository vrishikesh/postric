package util

import "net/http"

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

// func (p *RequestPayload) Render(w http.ResponseWriter, r *http.Request) error {
// 	return nil
// }

func (p *RequestPayload) Bind(r *http.Request) error {
	return nil
}

type ContextKey string
