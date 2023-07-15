package model

type UserLendBook struct {
	Id         int    `json:"id"`
	OwnerId    int    `json:"owner_id"`
	BorrowerId int    `json:"borrower_id"`
	UserBookId int    `json:"user_book_id"`
	Status     string `json:"status"` //requested, accepted, rejected, returned
}
