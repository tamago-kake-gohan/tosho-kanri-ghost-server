package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"tamago-kake-gohan.github.io/tosho-kanri-ghost/src/model"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := model.HealthResponse{}
	data.Message = "OK"
	err := json.NewEncoder(w).Encode(data)
	if nil != err {
		log.Println(err)
	}
}