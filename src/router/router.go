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
	mux.Handle("/api/v1/get_team_books", utils.CORS(handler.NewGetTeamBooksHandler(session, db)))
	mux.Handle("/api/v1/create_teams", utils.CORS(handler.NewCreateTeamsHandler(session, db)))
	mux.Handle("/api/v1/get_books", utils.CORS(handler.NewGetBooksHandler(session, db)))
	mux.Handle("/api/v1/add_book", utils.CORS(handler.NewAddBookHandler(session, db)))
	mux.Handle("/api/v1/get_book_detail", utils.CORS(handler.NewGetBookDetailHandler(session, db)))
	mux.Handle("/api/v1/update_user_book_state", utils.CORS(handler.NewUpdateUserBookStateHandler(session, db)))
	mux.Handle("/api/v1/request_rental", utils.CORS(handler.NewRequestRentalHandler(session, db)))
	mux.Handle("/api/v1/try_auth", utils.CORS(handler.NewTryAuthHandler(session)))
	mux.Handle("/api/v1/process_request", utils.CORS(handler.NewProcessRequestHandler(session, db)))
	mux.Handle("/api/v1/process_return", utils.CORS(handler.NewProcessReturnHandler(session, db)))
	mux.Handle("/api/v1/get_requests", utils.CORS(handler.NewGetRequestsHandler(session, db)))
	mux.Handle("/api/v1/set_rate", utils.CORS(handler.NewSetRateHandler(session, db)))
	mux.Handle("/api/v1/get_team_users", utils.CORS(handler.NewGetTeamUsersHandler(session, db)))
	return mux
}
