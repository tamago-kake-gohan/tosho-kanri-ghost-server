package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/astaxie/session"
	"tamago-kake-gohan.github.io/tosho-kanri-ghost/src/model"
	"tamago-kake-gohan.github.io/tosho-kanri-ghost/src/utils"
)

type RequestRentalResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type RequestRentalHandler struct {
	sess *session.Manager
	db   *sql.DB
}

type RequestRentalBody struct {
	UserBookId int `json:"user_book_id"`
	TeamId     int `json:"team_id"`
}

func NewRequestRentalHandler(sess *session.Manager, db *sql.DB) *RequestRentalHandler {
	return &RequestRentalHandler{
		sess: sess,
		db:   db,
	}
}

func (h *RequestRentalHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	sess := h.sess.SessionStart(w, r)
	userId := sess.Get("user_id")
	if nil == userId {
		w.WriteHeader(http.StatusForbidden)
		response := GetTeamsResponse{}
		response.Message = "ログインしてください"
		response.Status = "error"
		json.NewEncoder(w).Encode(response)
		return
	}
	var body RequestRentalBody
	json.Unmarshal(utils.GetRequestBody(r), &body)

	var userBook = model.UserBook{}
	err := h.db.QueryRow("SELECT * FROM `UserBook` WHERE `UserId` = ? AND `Id` = ?", userId, body.UserBookId).Scan(&userBook.Id, &userBook.UserId, &userBook.BookId, &userBook.State)
	if nil != err {
		w.WriteHeader(http.StatusNotFound)
		response := RequestRentalResponse{}
		response.Message = "該当の本が見つかりませんでした"
		response.Status = "error"
		json.NewEncoder(w).Encode(response)
	}
	if !utils.IsUserInTeam(h.db, userBook.BookId, body.TeamId) || !utils.IsUserInTeam(h.db, userId.(int), body.TeamId) {
		w.WriteHeader(http.StatusBadRequest)
		response := RequestRentalResponse{}
		response.Message = "チームに所属していません"
		response.Status = "error"
		json.NewEncoder(w).Encode(response)
		return
	}
	h.db.Exec("INSERT INTO `UserLendBook` (`UserBookId`, `OwnerId`, `BorrowerId`, `Status`) VALUES (?, ?, ?, ?)", userBook.Id, userBook.UserId, userId, "requested")
	response := RequestRentalResponse{}
	response.Message = "リクエストを送信しました"
	response.Status = "success"
	json.NewEncoder(w).Encode(response)
}
