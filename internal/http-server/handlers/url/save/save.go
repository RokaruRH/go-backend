package save

import (
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json: "alias,omitempty"`
}

const aliasLength = 6

type URLSaver interface {
	SaveURL(urtToSave string, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}
		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {

			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid requrest", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}
		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}
	}
}
