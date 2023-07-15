package router

import (
	"database/sql"
	"net/http"

	"github.com/astaxie/session"
	"tamago-kake-gohan.github.io/tosho-kanri-ghost/src/handler"
	"tamago-kake-gohan.github.io/tosho-kanri-ghost/src/utils"
)

func NewRouter(db *sql.DB, session *session.Manager) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/session_test", utils.CORS(handler.NewSessionTestHandler(session)))
	mux.Handle("/health", utils.CORS(handler.NewHealthHandler()))
	mux.Handle("/api/v1/login", utils.CORS(handler.NewLoginHandler(session, db)))
	mux.Handle("/api/v1/register", utils.CORS(handler.NewRegisterHandler(session, db)))
	mux.Handle("/api/v1/logout", utils.CORS(handler.NewLogoutHandler(session)))
	mux.Handle("/api/v1/get_teams", utils.CORS(handler.NewGetTeamsHandler(session, db)))
	mux.Handle("/api/v1/create_teams", utils.CORS(handler.NewCreateTeamsHandler(session, db)))
	mux.Handle("/api/v1/get_books", utils.CORS(handler.NewGetBooksHandler(session, db)))
	mux.Handle("/api/v1/add_book", utils.CORS(handler.NewAddBookHandler(session, db)))
	return mux
}
