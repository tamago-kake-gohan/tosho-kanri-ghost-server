package model

type Team struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Owner int    `json:"owner"`
}
