package getall

import (
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/sqlite"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	resp.Response
	URLs []sqlite.URL `json:"urls"`
}

type URLAllGetter interface {
	GetAll() ([]sqlite.URL, error)
}

func New(log *slog.Logger, urlAllGetter URLAllGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.getall.New"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		urls, err := urlAllGetter.GetAll()
		if err != nil {
			log.Error("failed to get data", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to get data"))

			return
		}

		responseOK(w, r, urls)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, urls []sqlite.URL) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		URLs:     urls,
	})
}
