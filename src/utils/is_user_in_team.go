package utils

import "database/sql"

func IsUserInTeam(db *sql.DB, userId int, teamId int) bool {
	err := db.QueryRow("SELECT * FROM `UserTeam` WHERE `UserId` = ? AND `TeamId` = ?", userId, teamId)
	return nil == err
}
