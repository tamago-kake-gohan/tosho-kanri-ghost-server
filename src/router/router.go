package router

import (
	"net/http"

	"tamago-kake-gohan.github.io/tosho-kanri-ghost/src/handler"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/health", new(handler.HealthHandler))

	return mux
}
