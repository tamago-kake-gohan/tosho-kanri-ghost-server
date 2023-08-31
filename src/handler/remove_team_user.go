package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/astaxie/session"
	"tamago-kake-gohan.github.io/tosho-kanri-ghost/src/model"
)

type RemoveTeamUserResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type RemoveTeamUserHandler struct {
	sess *session.Manager
	db   *sql.DB
}

func NewRemoveTeamUserHandler(sess *session.Manager, db *sql.DB) *RemoveTeamUserHandler {
	return &RemoveTeamUserHandler{
		sess: sess,
		db:   db,
	}
}

func (h *RemoveTeamUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response := RemoveTeamUserResponse{}
	sess := h.sess.SessionStart(w, r)
	userId := sess.Get("user_id")
	if nil == userId {
		w.WriteHeader(http.StatusForbidden)
		response := model.ForbiddenResponse{}
		response.Message = "ログインしてください"
		response.Status = "error"
		response.Code = 403
		json.NewEncoder(w).Encode(response)
		return
	}
	teamId := r.URL.Query().Get("team_id")
	var userTeam model.UserTeam
	err := h.db.QueryRow("SELECT * FROM `Team` WHERE `Id` = ? AND `Owner` = ?", teamId, userId).Scan(&userTeam.Id, &userTeam.UserId, &userTeam.TeamId)
	if nil != err {
		w.WriteHeader(http.StatusForbidden)
		response.Message = "権限がありません"
		response.Status = "error"
		json.NewEncoder(w).Encode(response)
		return
	}
	user := &model.User{}
	email := r.URL.Query().Get("email")
	err = h.db.QueryRow("SELECT * FROM `User` WHERE `Email` = ?", email).Scan(&user.Id, &user.Email, &user.Name, &user.Password)
	if nil != err {
		w.WriteHeader(http.StatusForbidden)
		response.Message = "ユーザーが見つかりませんでした"
		response.Status = "error"
		json.NewEncoder(w).Encode(response)
		return
	}
	h.db.Exec("DELETE FROM `UserTeam` WHERE `UserId` = ? AND `TeamId` = ?", user.Id, teamId)

	response.Message = ""
	response.Status = "success"
	json.NewEncoder(w).Encode(response)
}
