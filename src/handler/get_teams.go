package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/astaxie/session"
	"tamago-kake-gohan.github.io/tosho-kanri-ghost/src/model"
)

type GetTeamsResponse struct {
	Message string        `json:"message"`
	Status  string        `json:"status"`
	Data    []*model.Team `json:"data"`
}

type GetTeamsHandler struct {
	sess *session.Manager
	db   *sql.DB
}

func NewGetTeamsHandler(sess *session.Manager, db *sql.DB) *GetTeamsHandler {
	return &GetTeamsHandler{
		sess: sess,
		db:   db,
	}
}

func (h *GetTeamsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response := GetTeamsResponse{}
	sess := h.sess.SessionStart(w, r)
	userId := sess.Get("user_id")
	if nil == userId {
		w.WriteHeader(http.StatusForbidden)
		response.Message = "ログインしてください"
		response.Status = "error"
		json.NewEncoder(w).Encode(response)
		return
	}
	rows, err := h.db.Query("SELECT `Team`.* FROM `UserTeam` INNER JOIN `Team` ON `UserTeam`.`TeamId` = `Team`.`Id` WHERE `UserTeam`.`UserId` = ?", userId)
	if nil != err {
		response.Message = "データの取得に失敗しました"
		response.Status = "error"
		json.NewEncoder(w).Encode(response)
		return
	}
	result := make([]*model.Team, 0)
	for rows.Next() {
		team := &model.Team{}
		if err := rows.Scan(&team.Id, &team.Name, &team.Owner); err != nil {
			log.Fatalf("getRows rows.Scan error err:%v", err)
		}
		result = append(result, team)
	}
	response.Message = ""
	response.Status = "success"
	response.Data = result
	json.NewEncoder(w).Encode(response)
}
