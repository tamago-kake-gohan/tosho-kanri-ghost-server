package handler

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/astaxie/session"
	"tamago-kake-gohan.github.io/tosho-kanri-ghost/src/model"
	"tamago-kake-gohan.github.io/tosho-kanri-ghost/src/utils"
)

type AddBookResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type AddBookHandler struct {
	sess *session.Manager
	db   *sql.DB
}

type AddBookBody struct {
	ISBN string `json:"isbn"`
}

func NewAddBookHandler(sess *session.Manager, db *sql.DB) *AddBookHandler {
	return &AddBookHandler{
		sess: sess,
		db:   db,
	}
}

func (h *AddBookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	sess := h.sess.SessionStart(w, r)
	userId := sess.Get("user_id")
	if nil == userId {
		w.WriteHeader(http.StatusForbidden)
		response := AddBookResponse{}
		response.Message = "ログインしてください"
		response.Status = "error"
		json.NewEncoder(w).Encode(response)
		return
	}
	var body AddBookBody
	json.Unmarshal(utils.GetRequestBody(r), &body)
	var book model.Book
	err := h.db.QueryRow("SELECT * FROM `Book` WHERE `ISBN` = ?", body.ISBN).Scan(&book.Id, &book.ISBN, &book.Title, &book.Author)
	if nil != err {
		resp, _ := http.Get("https://api.openbd.jp/v1/get?isbn=" + body.ISBN)
		var meta []model.BookMeta
		byteArray, _ := io.ReadAll(resp.Body)
		json.Unmarshal(byteArray, &meta)
		h.db.Exec("INSERT INTO `Book` (`ISBN`, `Title`, `Author`) VALUES (?, ?, ?)", body.ISBN, meta[0].Summary.Title, meta[0].Summary.Author)
		h.db.QueryRow("SELECT * FROM `Book` WHERE `ISBN` = ?", body.ISBN).Scan(&book.Id, &book.ISBN, &book.Title, &book.Author)
	}
	_, err = h.db.Exec("INSERT INTO `UserBook` (`UserId`, `BookId`,`State`) VALUES (?, ?,'available')", userId, book.Id)
	if nil != err {
		log.Println(err)
	}
	response := AddBookResponse{}
	response.Message = "本を追加しました"
	response.Status = "success"
	json.NewEncoder(w).Encode(response)
}
