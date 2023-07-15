package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/astaxie/session"
)

type GetTeamBooksResponse struct {
	Message string      `json:"message"`
	Status  string      `json:"status"`
	Data    []*TeamBook `json:"data"`
}

type GetTeamBooksHandler struct {
	sess *session.Manager
	db   *sql.DB
}

type TeamBook struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	ISBN      string `json:"isbn"`
	Author    string `json:"author"`
	State     string `json:"state"` //available, lending, unavailable
	Owner     int    `json:"owner"`
	OwnerName string `json:"owner_name"`
}

func NewGetTeamBooksHandler(sess *session.Manager, db *sql.DB) *GetTeamBooksHandler {
	return &GetTeamBooksHandler{
		sess: sess,
		db:   db,
	}
}

func (h *GetTeamBooksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response := GetTeamBooksResponse{}
	sess := h.sess.SessionStart(w, r)
	userId := sess.Get("user_id")
	if nil == userId {
		w.WriteHeader(http.StatusForbidden)
		response.Message = "ログインしてください"
		response.Status = "error"
		json.NewEncoder(w).Encode(response)
		return
	}
	teamId := r.URL.Query().Get("team_id")
	err := h.db.QueryRow("SELECT * FROM `UserTeam` WHERE `UserId` = ? AND `TeamId` = ?", userId, teamId)
	if nil != err {
		w.WriteHeader(http.StatusForbidden)
		response.Message = "権限がありません"
		response.Status = "error"
		json.NewEncoder(w).Encode(response)
		return
	}
	rows, _ := h.db.Query("SELECT `UserBook`.`Id`, `Book`.`Title`, `Book`.`ISBN`, `Book`.`Author`, `UserBook`.`State`, `User`.`Id` as `Owner`, `User`.`name` as `OwnerName` FROM `UserTeam` INNER JOIN `User` ON `User`.`Id` = `UserTeam`.`UserId` INNER JOIN `UserBook` ON `UserBook`.`UserId` = `User`.`Id` INNER JOIN `Book` ON `Book`.`Id` = `UserBook`.`BookId` WHERE `UserTeam`.`TeamId` = ?", teamId)
	books := make([]*TeamBook, 0)
	for rows.Next() {
		book := &TeamBook{}
		if err := rows.Scan(&book.Id, &book.Title, &book.ISBN, &book.Author, &book.State, &book.Owner, &book.OwnerName); err != nil {
			log.Fatalf("getRows rows.Scan error err:%v", err)
		}
		books = append(books, book)
	}

	response.Message = ""
	response.Status = "success"
	response.Data = books
	json.NewEncoder(w).Encode(response)
}
