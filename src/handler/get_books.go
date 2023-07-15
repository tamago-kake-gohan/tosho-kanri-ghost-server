package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/astaxie/session"
	"tamago-kake-gohan.github.io/tosho-kanri-ghost/src/model"
)

type GetBooksResponse struct {
	Message string        `json:"message"`
	Status  string        `json:"status"`
	Data    []*model.Book `json:"data"`
}

type GetBooksHandler struct {
	sess *session.Manager
	db   *sql.DB
}

func NewGetBooksHandler(sess *session.Manager, db *sql.DB) *GetBooksHandler {
	return &GetBooksHandler{
		sess: sess,
		db:   db,
	}
}

func (h *GetBooksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response := GetBooksResponse{}
	sess := h.sess.SessionStart(w, r)
	userId := sess.Get("user_id")
	if nil == userId {
		w.WriteHeader(http.StatusForbidden)
		response.Message = "ログインしてください"
		response.Status = "error"
		json.NewEncoder(w).Encode(response)
		return
	}
	rows, err := h.db.Query("SELECT `Book`.*,`UserBook`.`State` FROM `UserBook` INNER JOIN `Book` ON `UserBook`.`BookId` = `Book`.`Id` WHERE `UserBook`.`UserId` = ?", userId)
	if nil != err {
		response.Message = "データの取得に失敗しました"
		response.Status = "error"
		json.NewEncoder(w).Encode(response)
		return
	}
	result := make([]*model.Book, 0)
	for rows.Next() {
		book := &model.Book{}
		if err := rows.Scan(&book.Id, &book.Author, &book.CCode, &book.ISBN, &book.State, &book.Title); err != nil {
			log.Fatalf("getRows rows.Scan error err:%v", err)
		}
		result = append(result, book)
	}
	response.Message = ""
	response.Status = "success"
	response.Data = result
	json.NewEncoder(w).Encode(response)
}
