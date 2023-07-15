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
	mux.Handle("/api/v1/login", handler.NewLoginHandler(session, db))
	mux.Handle("/api/v1/register", handler.NewRegisterHandler(session, db))
	return mux
}
