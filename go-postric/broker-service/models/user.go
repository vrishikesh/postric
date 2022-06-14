package models

import "net/http"

type User struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type UserPayload struct {
	*User
	Role string `json:"role"`
}

func (u *UserPayload) Bind(r *http.Request) error {
	return nil
}

func (u *UserPayload) Render(w http.ResponseWriter, r *http.Request) error {
	u.Role = "collaborator"
	return nil
}
