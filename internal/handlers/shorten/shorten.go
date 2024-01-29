package shorten

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/go-playground/validator"

	"github.com/KKitsun/link-shortener/internal/handlers/response"
)

type URLRequest struct {
	URL string `json:"url" validate:"required,url"`
}

type AliasResponse struct {
	Alias string `json:"alias"`
	response.Response
}

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

func Shorten(urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req URLRequest

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, response.ErrorResponse("Error decoding request", err))
			return
		}

		if err := validator.New().Struct(req); err != nil {
			render.JSON(w, r, response.ErrorResponse("Invalid request", err))
			return
		}

		alias := generateAlias()

		_, err = urlSaver.SaveURL(req.URL, alias)
		if err != nil {
			render.JSON(w, r, response.ErrorResponse("Error saving url to the database", err))
			return
		}

		render.JSON(w, r, AliasResponse{
			Alias: alias,
			Response: response.Response{
				Status: "Success",
			},
		})

	}
}

func generateAlias() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const codeLength = 6

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	alias := make([]byte, codeLength)
	for i := range alias {
		alias[i] = charset[rnd.Intn(len(charset))]
	}
	return string(alias)
}
