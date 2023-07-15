package router

import (
	"database/sql"
	"net/http"

	"github.com/astaxie/session"
	"tamago-kake-gohan.github.io/tosho-kanri-ghost/src/handler"
)

func NewRouter(db *sql.DB, session *session.Manager) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/session_test", handler.NewSessionTestHandler(session))
	mux.Handle("/health", handler.NewHealthHandler())
	return mux
}
