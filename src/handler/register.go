package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/astaxie/session"
	"golang.org/x/crypto/bcrypt"
	"tamago-kake-gohan.github.io/tosho-kanri-ghost/src/model"
	"tamago-kake-gohan.github.io/tosho-kanri-ghost/src/utils"
)

type RegisterResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type RegisterHandler struct {
	sess *session.Manager
	db   *sql.DB
}

type RegisterBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewRegisterHandler(sess *session.Manager, db *sql.DB) *RegisterHandler {
	return &RegisterHandler{
		sess: sess,
		db:   db,
	}
}

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response := RegisterResponse{}
	var body LoginBody
	json.Unmarshal(utils.GetRequestBody(r), &body)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	user := &model.User{}
	err := h.db.QueryRow("SELECT * FROM User WHERE email = ?", body.Email).Scan(&user.Id, &user.Email, &user.Password)
	if nil == err {
		log.Println("すでに登録されているメールアドレスです", err)
		response.Message = "すでに登録されているメールアドレスです"
		response.Status = "error"
		json.NewEncoder(w).Encode(response)
		return
	}
	hadhed, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if nil != err {
		log.Println("不適切なパスワードです", err)
		response.Message = "不適切なパスワードです"
		response.Status = "error"
		json.NewEncoder(w).Encode(response)
		return
	}
	h.db.Exec("INSERT INTO User (email,password) VALUES (?,?)", body.Email, string(hadhed))
	h.db.QueryRow("SELECT * FROM User WHERE email = ?", body.Email).Scan(&user.Id, &user.Email, &user.Password)
	sess := h.sess.SessionStart(w, r)
	sess.Set("user_id", user.Id)
	response.Message = "ログインしました"
	response.Status = "success"
	json.NewEncoder(w).Encode(response)
}
