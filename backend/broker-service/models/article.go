package models

import (
	"errors"
	"net/http"
	"strings"
)

type Article struct {
	ID     string `json:"id"`
	UserID uint64 `json:"user_id"`
	Title  string `json:"title"`
	Slug   string `json:"slug"`
}

type ArticleRequest struct {
	*Article

	User *UserPayload `json:"user,omitempty"`

	ProtectedID string `json:"id"`
}

func (a *ArticleRequest) Bind(r *http.Request) error {
	if a.Article == nil {
		return errors.New("missing required Article fields")
	}

	a.ProtectedID = ""                                 // unset the protected ID
	a.Article.Title = strings.ToLower(a.Article.Title) // as an example, we down-case
	return nil
}

type ArticleResponse struct {
	*Article

	User *UserPayload `json:"user,omitempty"`

	Elapsed int64 `json:"elapsed"`
}

func (rd *ArticleResponse) Render(w http.ResponseWriter, r *http.Request) error {
	rd.Elapsed = 10
	return nil
}
