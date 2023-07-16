package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/astaxie/session"
	"tamago-kake-gohan.github.io/tosho-kanri-ghost/src/model"
)

type bookDetail struct {
	OwnerName    string  `json:"owner_name"`
	BorrowerName string  `json:"borrower_name"`
	Title        string  `json:"title"`
	ISBN         string  `json:"isbn"`
	State        string  `json:"state"` //available, lending, unavailable
	Rating       float64 `json:"rating"`
}

type GetBookDetailResponse struct {
	Message string     `json:"message"`
	Status  string     `json:"status"`
	Data    bookDetail `json:"data"`
}

type GetBookDetailHandler struct {
	sess *session.Manager
	db   *sql.DB
}

func NewGetBookDetailHandler(sess *session.Manager, db *sql.DB) *GetBookDetailHandler {
	return &GetBookDetailHandler{
		sess: sess,
		db:   db,
	}
}

func (h *GetBookDetailHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response := GetBookDetailResponse{}
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
	userBookId := r.URL.Query().Get("user_book_id")
	var book bookDetail
	err := h.db.QueryRow(
		"SELECT `Owner`.`Name` as `OwnerName`, `Book`.`Title`,`Book`.`ISBN`, `UserBook`.`State`, `Borrower`.`Name` as `BorrowerName`, `Review`.`Rating`"+
			"FROM `UserBook`"+
			"INNER JOIN `User` as `Owner` ON `UserBook`.`UserId` = `Owner`.`Id` "+
			"INNER JOIN `Book` ON `UserBook`.`BookId` = `Book`.`Id`"+
			"LEFT JOIN `Review` ON `Owner`.`Id` = `Review`.`UserId` AND `Review`.`BookId` = `Book`.`Id`"+
			"LEFT JOIN `UserLendBook` ON `UserLendBook`.`OwnerId` = `Owner`.`Id` AND `UserLendBook`.`UserBookId` = `UserBook`.`Id` AND `UserLendBook`.`Status` = 'accepted'"+
			"LEFT JOIN `User` as `Borrower` ON `UserLendBook`.`BorrowerId` = `Borrower`.`Id`"+
			"WHERE `UserBook`.`Id` = ?", userBookId).Scan(&book.OwnerName, &book.Title, &book.State, &book.BorrowerName, &book.Rating, &book.ISBN)
	if nil != err {
		response.Message = "データの取得に失敗しました"
		response.Status = "error"
		json.NewEncoder(w).Encode(response)
		return
	}
	response.Message = ""
	response.Status = "success"
	response.Data = book
	json.NewEncoder(w).Encode(response)
}
