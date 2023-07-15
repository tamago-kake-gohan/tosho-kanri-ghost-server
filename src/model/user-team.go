package model

type UserTeam struct {
	Id     int `json:"id"`
	UserId int `json:"user_id"`
	TeamId int `json:"team_id"`
}
