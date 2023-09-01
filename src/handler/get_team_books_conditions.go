package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/astaxie/session"
	"tamago-kake-gohan.github.io/tosho-kanri-ghost/src/model"
)

type GetTeamBooksConditionResponse struct {
	Message string      `json:"message"`
	Status  string      `json:"status"`
	Data    []*teamBook `json:"data"`
}

type GetTeamBooksConditionHandler struct {
	sess *session.Manager
	db   *sql.DB
}

func NewGetTeamBooksConditionHandler(sess *session.Manager, db *sql.DB) *GetTeamBooksConditionHandler {
	return &GetTeamBooksConditionHandler{
		sess: sess,
		db:   db,
	}
}

func (h *GetTeamBooksConditionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response := GetTeamBooksResponse{}
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
	err := h.db.QueryRow("SELECT * FROM `UserTeam` WHERE `UserId` = ? AND `TeamId` = ?", userId, teamId).Scan(&userTeam.Id, &userTeam.UserId, &userTeam.TeamId)
	if nil != err {
		w.WriteHeader(http.StatusForbidden)
		response.Message = "権限がありません"
		response.Status = "error"
		json.NewEncoder(w).Encode(response)
		return
	}
	keyword := r.URL.Query().Get("keyword")
	status := r.URL.Query().Get("status")
	statement := "SELECT `UserBook`.`Id`, `Book`.`Title`, `UserBook`.`State`, `User`.`name` as `OwnerName` FROM `UserTeam` INNER JOIN `User` ON `User`.`Id` = `UserTeam`.`UserId` INNER JOIN `UserBook` ON `UserBook`.`UserId` = `User`.`Id` INNER JOIN `Book` ON `Book`.`Id` = `UserBook`.`BookId` WHERE `UserTeam`.`TeamId` = ? AND `Title` LIKE '%" + keyword + "%' "
	if status != "all" {
		statement += "AND State = '" + status + "'"
	}
	rows, _ := h.db.Query(statement, teamId)
	books := make([]*teamBook, 0)
	for rows.Next() {
		book := &teamBook{}
		if err := rows.Scan(&book.Id, &book.Title, &book.State, &book.OwnerName); err != nil {
			log.Printf("getRows rows.Scan error err:%v", err)
		}
		books = append(books, book)
	}

	response.Message = ""
	response.Status = "success"
	response.Data = books
	json.NewEncoder(w).Encode(response)
}
