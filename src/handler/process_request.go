package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/astaxie/session"
	"tamago-kake-gohan.github.io/tosho-kanri-ghost/src/model"
	"tamago-kake-gohan.github.io/tosho-kanri-ghost/src/utils"
)

type ProcessRequestResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type ProcessRequestHandler struct {
	sess *session.Manager
	db   *sql.DB
}

type ProcessRequestBody struct {
	UserBookId int  `json:"user_book_id"`
	Accept     bool `json:"accept"`
}

func NewProcessRequestHandler(sess *session.Manager, db *sql.DB) *ProcessRequestHandler {
	return &ProcessRequestHandler{
		sess: sess,
		db:   db,
	}
}

func (h *ProcessRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	sess := h.sess.SessionStart(w, r)
	userId := sess.Get("user_id")
	if nil == userId {
		w.WriteHeader(http.StatusForbidden)
		response := model.ForbiddenResponse{}
		response.Message = "ログインしてください"
		response.Status = "error"
		response.Code = 403
		json.NewEncoder(w).Encode(response)
		return
	}
	var body ProcessRequestBody
	json.Unmarshal(utils.GetRequestBody(r), &body)

	var userBook = model.UserBook{}
	err := h.db.QueryRow("SELECT * FROM `UserLendBook` WHERE `UserId` = ? AND `Id` = ?", userId, body.UserBookId).Scan(&userBook.Id, &userBook.UserId, &userBook.BookId, &userBook.State)
	if nil != err {
		w.WriteHeader(http.StatusNotFound)
		response := ProcessRequestResponse{}
		response.Message = "該当の本が見つかりませんでした"
		response.Status = "error"
		json.NewEncoder(w).Encode(response)
	}
	status := "accepted"
	if !body.Accept {
		status = "rejected"
	}
	h.db.Exec("UPDATE `UserLendBook` SET `Status` = ? WHERE `Id` = ?", status, body.UserBookId)
	response := ProcessRequestResponse{}
	response.Message = "リクエストを送信しました"
	response.Status = "success"
	json.NewEncoder(w).Encode(response)
}
