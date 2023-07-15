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

type LoginResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type LoginHandler struct {
	sess *session.Manager
	db   *sql.DB
}

type LoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewLoginHandler(sess *session.Manager, db *sql.DB) *LoginHandler {
	return &LoginHandler{
		sess: sess,
		db:   db,
	}
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response := LoginResponse{}
	var body LoginBody
	json.Unmarshal(utils.GetRequestBody(r), &body)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	user := &model.User{}
	err := h.db.QueryRow("SELECT * FROM User WHERE email = ?", body.Email).Scan(&user.Id, &user.Email, &user.Name, &user.Password)
	if nil != err {
		log.Println("ユーザーが存在しません", err)
		response.Message = "メールアドレスまたはパスワードが間違っています"
		response.Status = "error"
		json.NewEncoder(w).Encode(response)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if nil != err {
		log.Println("メールアドレスまたはパスワードが間違っています", err)
		response.Message = "メールアドレスまたはパスワードが間違っています"
		response.Status = "error"
		json.NewEncoder(w).Encode(response)
		return
	}

	sess := h.sess.SessionStart(w, r)
	sess.Set("user_id", user.Id)
	response.Message = "ログインしました"
	response.Status = "success"
	json.NewEncoder(w).Encode(response)
}
