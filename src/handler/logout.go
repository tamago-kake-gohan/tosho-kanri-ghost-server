package handler

import (
	"encoding/json"
	"net/http"

	"github.com/astaxie/session"
)

type LogoutResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type LogoutHandler struct {
	sess *session.Manager
}

func NewLogoutHandler(sess *session.Manager) *LogoutHandler {
	return &LogoutHandler{
		sess: sess,
	}
}

func (h *LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sess := h.sess.SessionStart(w, r)
	sess.Delete("user_id")
	response := LogoutResponse{}
	response.Message = "ログアウトしました"
	response.Status = "success"
	json.NewEncoder(w).Encode(response)
}
