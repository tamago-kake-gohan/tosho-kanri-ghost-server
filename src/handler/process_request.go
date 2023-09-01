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
	UserLendBookId int  `json:"user_lend_book_id"`
	Accept         bool `json:"accept"`
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
	err := h.db.QueryRow("SELECT `Id` FROM `UserLendBook` WHERE `OwnerId` = ? AND `Id` = ? AND `Status` = 'requested'", userId, body.UserLendBookId).Scan(&userBook.Id)
	if nil != err {
		w.WriteHeader(http.StatusNotFound)
		response := ProcessRequestResponse{}
		response.Message = "該当の本が見つかりませんでした"
		response.Status = "error"
		json.NewEncoder(w).Encode(response)
		return
	}
	status := "accepted"
	if !body.Accept {
		status = "rejected"
	}
	h.db.Exec("UPDATE `UserLendBook` SET `Status` = ? WHERE `Id` = ?", status, body.UserLendBookId)
	response := ProcessRequestResponse{}
	response.Message = "リクエストを処理しました"
	response.Status = "success"
	json.NewEncoder(w).Encode(response)
}
