package model

type UserBook struct {
	Id     int    `json:"id"`
	UserId int    `json:"user_id"`
	BookId int    `json:"book_id"`
	State  string `json:"state"` //available, lending, unavailable
}
