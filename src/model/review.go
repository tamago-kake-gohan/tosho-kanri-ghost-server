package model

type Review struct {
	Id      int     `json:"id"`
	BookId  int     `json:"book_id"`
	UserId  int     `json:"user_id"`
	Comment string  `json:"comment"`
	Rating  float64 `json:"rating"`
}
