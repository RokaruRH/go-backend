package save

import (
	"net/http"
	"strings"

	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
)

// initialize a single validator instance
var validate = validator.New()

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias" validate:"required,alphanum,min=3,max=30"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias"`
}

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

// New returns an HTTP handler that expects a JSON body with URL and Alias.
// It validates the input, enforces a request size limit, delegates saving
// to URLSaver, and returns a structured JSON response with proper status codes.
func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"
		ctx := r.Context()
		logger := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(ctx)),
		)

		// Limit JSON body size to prevent abuse
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB
		defer r.Body.Close()

		var req Request
		// Decode JSON
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			// Check for body-too-large by matching error text
			if strings.Contains(err.Error(), "request body too large") {
				logger.Error("request body too large", sl.Err(err))
				render.Status(r, http.StatusRequestEntityTooLarge)
				render.JSON(w, r, resp.Error("request body too large"))
				return
			}
			logger.Error("decode JSON failed", sl.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid JSON payload"))
			return
		}

		// Validate required fields and constraints
		if err := validate.StructCtx(ctx, &req); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			logger.Error("validation failed", sl.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.ValidationError(validateErrs))
			return
		}

		// Delegate saving the URL
		id, err := urlSaver.SaveURL(req.URL, req.Alias)
		if err != nil {
			logger.Error("service SaveURL failed", sl.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("could not save URL"))
			return
		}

		logger.Info("URL saved", slog.Int64("id", id), slog.String("alias", req.Alias))

		// Return the successful response
		render.Status(r, http.StatusCreated)
		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias:    req.Alias,
		})
	}
}
