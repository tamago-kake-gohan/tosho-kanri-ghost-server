package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/astaxie/session"
)

type indexBook struct {
	Id     int     `json:"id"`
	Title  string  `json:"title"`
	State  string  `json:"state"` //available, lending, unavailable
	Rating float64 `json:"rating"`
}

type GetBooksResponse struct {
	Message string       `json:"message"`
	Status  string       `json:"status"`
	Data    []*indexBook `json:"data"`
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
	rows, err := h.db.Query("SELECT `UserBook`.`Id`,`Book`.`Title`,`UserBook`.`State` FROM `UserBook` INNER JOIN `Book` ON `UserBook`.`BookId` = `Book`.`Id` WHERE `UserBook`.`UserId` = ?", userId)
	if nil != err {
		response.Message = "データの取得に失敗しました"
		response.Status = "error"
		json.NewEncoder(w).Encode(response)
		return
	}
	result := make([]*indexBook, 0)
	for rows.Next() {
		book := &indexBook{}
		if err := rows.Scan(&book.Id, &book.Title, &book.State, &book.Title); err != nil {
			log.Printf("getRows rows.Scan error err:%v", err)
			continue
		}
		h.db.QueryRow("SELECT `Review`.`Rating` FROM `Review` WHERE `BookId` = ? AND `UserId` = ?", book.Id, userId).Scan(&book.Rating)
		result = append(result, book)
	}
	response.Message = ""
	response.Status = "success"
	response.Data = result
	json.NewEncoder(w).Encode(response)
}
