package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

type HealthResponse struct {
	Message string `json:"message"`
}

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	data := HealthResponse{}
	data.Message = "OK"
	err := json.NewEncoder(w).Encode(data)
	if nil != err {
		log.Println(err)
	}
}
