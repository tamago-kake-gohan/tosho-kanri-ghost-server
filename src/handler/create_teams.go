package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/astaxie/session"
	"tamago-kake-gohan.github.io/tosho-kanri-ghost/src/model"
	"tamago-kake-gohan.github.io/tosho-kanri-ghost/src/utils"
)

type CreateTeamsResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	TeamId  int    `json:"team_id"`
}

type CreateTeamsHandler struct {
	sess *session.Manager
	db   *sql.DB
}

type CreateTeamsBody struct {
	Name   string   `json:"name"`
	Emails []string `json:"emails"`
}

type MissingUserErrorResponse struct {
	Message string   `json:"message"`
	Status  string   `json:"status"`
	Data    []string `json:"data"`
}

func NewCreateTeamsHandler(sess *session.Manager, db *sql.DB) *CreateTeamsHandler {
	return &CreateTeamsHandler{
		sess: sess,
		db:   db,
	}
}

func (h *CreateTeamsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	sess := h.sess.SessionStart(w, r)
	userId := sess.Get("user_id")
	if nil == userId {
		w.WriteHeader(http.StatusForbidden)
		response := GetTeamsResponse{}
		response.Message = "ログインしてください"
		response.Status = "error"
		json.NewEncoder(w).Encode(response)
		return
	}
	var body CreateTeamsBody
	json.Unmarshal(utils.GetRequestBody(r), &body)
	var missingUser []string
	users := []int{userId.(int)}
	for _, email := range body.Emails {
		user := &model.User{}
		err := h.db.QueryRow("SELECT * FROM `User` WHERE `Email` = ?", email).Scan(&user.Id, &user.Email, &user.Name, &user.Password)
		if nil != err {
			missingUser = append(missingUser, email)
		}
		users = append(users, user.Id)
	}
	if len(missingUser) > 0 {
		response := MissingUserErrorResponse{}
		response.Message = "ユーザーが見つかりませんでした"
		response.Status = "error"
		response.Data = missingUser
		json.NewEncoder(w).Encode(response)
		return
	}
	res, _ := h.db.Exec("INSERT INTO `Team` (`Name`, `Owner`) VALUES (?, ?)", body.Name, userId)
	teamId, _ := res.LastInsertId()
	for _, user := range users {
		h.db.Exec("INSERT INTO `UserTeam` (`UserId`, `TeamId`) VALUES (?, ?)", user, teamId)
	}
	response := CreateTeamsResponse{}
	response.Message = "チームを作成しました"
	response.Status = "success"
	response.TeamId = int(teamId)
	json.NewEncoder(w).Encode(response)
}
