package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/astaxie/session"
	"tamago-kake-gohan.github.io/tosho-kanri-ghost/src/model"
)

type GetTeamUsersResponse struct {
	Message string      `json:"message"`
	Status  string      `json:"status"`
	Users   []*TeamUser `json:"users"`
}

type GetTeamUsersHandler struct {
	sess *session.Manager
	db   *sql.DB
}

type TeamUser struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewGetTeamUsersHandler(sess *session.Manager, db *sql.DB) *GetTeamUsersHandler {
	return &GetTeamUsersHandler{
		sess: sess,
		db:   db,
	}
}

func (h *GetTeamUsersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response := GetTeamUsersResponse{}
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
	rows, _ := h.db.Query("SELECT `User`.`Id`,`User`.`Name`, `User`.`Email` FROM `UserTeam` INNER JOIN `User` ON `User`.`Id` = `UserTeam`.`UserId` WHERE `UserTeam`.`TeamId` = ?", teamId)
	users := make([]*TeamUser, 0)
	for rows.Next() {
		user := &TeamUser{}
		if err := rows.Scan(&user.Id, &user.Name, &user.Email); err != nil {
			log.Printf("getRows rows.Scan error err:%v", err)
		}
		users = append(users, user)
	}

	response.Message = ""
	response.Status = "success"
	response.Users = users
	json.NewEncoder(w).Encode(response)
}
