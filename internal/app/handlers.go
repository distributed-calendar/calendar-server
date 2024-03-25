package app

import (
	"log/slog"
	"net/http"
)

func (a *App) pingHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("pong"))
		if err != nil {
			slog.Error("cannot write to ping", err)
		}
	})
}
