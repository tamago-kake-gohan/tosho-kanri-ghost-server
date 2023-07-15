package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/astaxie/session"
	"tamago-kake-gohan.github.io/tosho-kanri-ghost/src/model"
	"tamago-kake-gohan.github.io/tosho-kanri-ghost/src/utils"
)

type UpdateUserBookStateResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type UpdateUserBookStateHandler struct {
	sess *session.Manager
	db   *sql.DB
}

type UpdateUserBookStateBody struct {
	UserBookId int    `json:"user_book_id"`
	State      string `json:"state"` //available, unavailable
}

func NewUpdateUserBookStateHandler(sess *session.Manager, db *sql.DB) *UpdateUserBookStateHandler {
	return &UpdateUserBookStateHandler{
		sess: sess,
		db:   db,
	}
}

func (h *UpdateUserBookStateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	var body UpdateUserBookStateBody
	json.Unmarshal(utils.GetRequestBody(r), &body)

	var userBook = model.UserBook{}
	err := h.db.QueryRow("SELECT * FROM `UserBook` WHERE `UserId` = ? AND `Id` = ?", userId, body.UserBookId).Scan(&userBook.Id, &userBook.UserId, &userBook.BookId, &userBook.State)
	if nil != err {
		w.WriteHeader(http.StatusNotFound)
		response := UpdateUserBookStateResponse{}
		response.Message = "該当の本が見つかりませんでした"
		response.Status = "error"
		json.NewEncoder(w).Encode(response)
	}
	if userBook.State == "lending" {
		w.WriteHeader(http.StatusBadRequest)
		response := UpdateUserBookStateResponse{}
		response.Message = "貸出中の本は状態を変更できません"
		response.Status = "error"
		json.NewEncoder(w).Encode(response)
		return
	}
	if body.State != "available" && body.State != "unavailable" {
		w.WriteHeader(http.StatusBadRequest)
		response := UpdateUserBookStateResponse{}
		response.Message = "値が正しくありません"
		response.Status = "error"
		json.NewEncoder(w).Encode(response)
		return
	}
	h.db.Exec("UPDATE `UserBook` SET `State` = ? WHERE `Id` = ?", body.State, body.UserBookId)
	response := UpdateUserBookStateResponse{}
	response.Message = "更新しました"
	response.Status = "success"
	json.NewEncoder(w).Encode(response)
}
