package utils

import (
	"database/sql"

	"tamago-kake-gohan.github.io/tosho-kanri-ghost/src/model"
)

func IsUserInTeam(db *sql.DB, userId int, teamId int) bool {
	var userTeam = model.UserTeam{}
	err := db.QueryRow("SELECT * FROM `UserTeam` WHERE `UserId` = ? AND `TeamId` = ?", userId, teamId).Scan(&userTeam.Id, &userTeam.UserId, &userTeam.TeamId)
	return nil == err
}
