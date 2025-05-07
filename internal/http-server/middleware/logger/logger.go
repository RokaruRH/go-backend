package logger

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func New(log *slog.logger) func(next http.Handler) http.Handler{
	return func(next http.Handler) http.Handler{
		log = log.With(
			slog.String("component", "moddleware/logger"),
		)

		log.Info("logger middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request){
			entry := log.With(
				slog.String("method",         r.Method    ),
				slog.String("path",           r.URL.Path  ),
				slog.String("remote_addr",    r.RemoteAddr),
				slog.String("user_agent",    r.UserAgent()),
				slog.String("request_id",  middleware.GetReqID(r.Context())),
			)
			ww := middleware.NewWrapResponseWriter(w,r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				entry.Info("request complated",
					slog.Int("status", ww.Status()),
				)
			}
		}
	}
}
