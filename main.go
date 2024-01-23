package main

import (
	"errors"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

var urls = make(map[string]string)

func handleShorten(w http.ResponseWriter, r *http.Request) {

	data := &LinkRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	shortKey := generateShortKey()
	urls[shortKey] = data.Link.URL
	link := &Link{URL: shortKey}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, &LinkResponse{Link: link})
}

func handleFullLinkRequest(w http.ResponseWriter, r *http.Request) {

	data := &LinkRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	originalURL, found := urls[data.Link.URL]
	if !found {
		http.Error(w, "Shortened key not found", http.StatusNotFound)
		return
	}

	link := &Link{URL: originalURL}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, &LinkResponse{Link: link})

}

func handleRedirect(w http.ResponseWriter, r *http.Request) {

	if shortenedURL := chi.URLParam(r, "shortenedURL"); shortenedURL != "" {
		originalURL, found := urls[shortenedURL]
		if !found {
			http.Error(w, "Shortened key not found", http.StatusNotFound)
			return
		}
		http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
	} else {
		http.Error(w, "Shortened key is missing", http.StatusBadRequest)
		return
	}

}

func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rnd.Intn(len(charset))]
	}
	return string(shortKey)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hi"))
		})
		r.Post("/shorten", handleShorten)
		r.Post("/full", handleFullLinkRequest)
		r.Route("/sl", func(r chi.Router) {
			r.Get("/{shortenedURL}", handleRedirect)
		})
	})

	http.ListenAndServe(":3030", r)
}

type Link struct {
	URL string `json:"url"`
}

type LinkRequest struct {
	*Link
}

func (a *LinkRequest) Bind(r *http.Request) error {
	if a.Link == nil {
		return errors.New("missing required Article fields.")
	}

	return nil
}

type LinkResponse struct {
	*Link
}

func (rd *LinkResponse) Render(w http.ResponseWriter, r *http.Request) error {

	return nil
}

type ErrResponse struct {
	Err            error `json:"-"`
	HTTPStatusCode int   `json:"-"`

	StatusText string `json:"status"`
	AppCode    int64  `json:"code,omitempty"`
	ErrorText  string `json:"error,omitempty"` 
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

