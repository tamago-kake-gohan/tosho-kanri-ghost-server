package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/astaxie/session"
	"tamago-kake-gohan.github.io/tosho-kanri-ghost/src/model"
	"tamago-kake-gohan.github.io/tosho-kanri-ghost/src/utils"
)

type ProcessReturnResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type ProcessReturnHandler struct {
	sess *session.Manager
	db   *sql.DB
}

type ProcessReturnBody struct {
	UserBookId int `json:"user_book_id"`
}

func NewProcessReturnHandler(sess *session.Manager, db *sql.DB) *ProcessReturnHandler {
	return &ProcessReturnHandler{
		sess: sess,
		db:   db,
	}
}

func (h *ProcessReturnHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	var body ProcessReturnBody
	json.Unmarshal(utils.GetRequestBody(r), &body)

	var userBook = model.UserBook{}
	err := h.db.QueryRow("SELECT * FROM `UserLendBook` WHERE `UserId` = ? AND `Id` = ? AND `Status` = 'accepted'", userId, body.UserBookId).Scan(&userBook.Id, &userBook.UserId, &userBook.BookId, &userBook.State)
	if nil != err {
		w.WriteHeader(http.StatusNotFound)
		response := ProcessReturnResponse{}
		response.Message = "該当の本が見つかりませんでした"
		response.Status = "error"
		json.NewEncoder(w).Encode(response)
	}
	h.db.Exec("UPDATE `UserLendBook` SET `Status` = 'returned' WHERE `Id` = ?", body.UserBookId)
	response := ProcessReturnResponse{}
	response.Message = "返却しました"
	response.Status = "success"
	json.NewEncoder(w).Encode(response)
}
