package getFull

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/KKitsun/link-shortener/internal/handlers/response"
)

type FullURLResponse struct {
	URL string `json:"url"`
	response.Response
}

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func GetFull(urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			render.JSON(w, r, response.ErrorResponse("Invalid alias request", nil))
			return
		}

		urlFromDB, err := urlGetter.GetURL(alias)
		if err != nil {
			render.JSON(w, r, response.ErrorResponse("Error getting url", err))
			return
		}

		// http.Redirect(w, r, urlFromDB, http.StatusFound)
		render.JSON(w, r, FullURLResponse{
			URL: urlFromDB,
			Response: response.Response{
				Status: "Success",
			},
		})
	}
}
