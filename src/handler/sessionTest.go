package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/astaxie/session"
)

type SessionTestResponse struct {
	Count int `json:"count"`
}

type SessionTestHandler struct {
	sess *session.Manager
}

func NewSessionTestHandler(sess *session.Manager) *SessionTestHandler {
	return &SessionTestHandler{
		sess: sess,
	}
}

func (h *SessionTestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	sess := h.sess.SessionStart(w, r)
	ct := sess.Get("countnum")
	if ct == nil {
		sess.Set("countnum", 1)
		ct = 1
	} else {
		sess.Set("countnum", (ct.(int) + 1))
		ct = (ct.(int) + 1)
	}
	data := SessionTestResponse{}
	data.Count = ct.(int)
	err := json.NewEncoder(w).Encode(data)
	if nil != err {
		log.Println(err)
	}
}
