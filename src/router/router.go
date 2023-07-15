package router

import (
	"database/sql"
	"net/http"

	"tamago-kake-gohan.github.io/tosho-kanri-ghost/src/handler"
)

func NewRouter(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/health", handler.NewHealthHandler())

	return mux
}
