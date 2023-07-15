package handler

import (
	"encoding/json"
	"net/http"

	"github.com/astaxie/session"
)

type TryAuthResponse struct {
	IsLoggedIn bool `json:"is_logged_in"`
}

type TryAuthHandler struct {
	sess *session.Manager
}

func NewTryAuthHandler(sess *session.Manager) *TryAuthHandler {
	return &TryAuthHandler{
		sess: sess,
	}
}

func (h *TryAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sess := h.sess.SessionStart(w, r)
	response := TryAuthResponse{}
	response.IsLoggedIn = nil != sess.Get("user_id")
	json.NewEncoder(w).Encode(response)
}
