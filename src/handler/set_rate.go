package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/astaxie/session"
	"tamago-kake-gohan.github.io/tosho-kanri-ghost/src/model"
	"tamago-kake-gohan.github.io/tosho-kanri-ghost/src/utils"
)

type SetRateResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type SetRateHandler struct {
	sess *session.Manager
	db   *sql.DB
}

type SetRateBody struct {
	BookId  int     `json:"book_id"`
	Rate    float64 `json:"rate"`
	Comment string  `json:"comment"`
}

func NewSetRateHandler(sess *session.Manager, db *sql.DB) *SetRateHandler {
	return &SetRateHandler{
		sess: sess,
		db:   db,
	}
}

func (h *SetRateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
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
	var body SetRateBody
	json.Unmarshal(utils.GetRequestBody(r), &body)

	var review = model.Review{}
	err := h.db.QueryRow("SELECT * FROM `Review` WHERE `UserId` = ? AND `BookId` = ?", userId, body.BookId).Scan(&review.Id, &review.UserId, &review.BookId, &review.Rating, &review.Comment)
	if nil != err {
		h.db.Exec("INSERT INTO `Review` (`UserId`, `BookId`, `Rating`,`Comment`) VALUES (?, ?, ?,?)", userId, body.BookId, body.Rate, body.Comment)
	} else {
		h.db.Exec("UPDATE `Review` SET `Rating` = ?, `Comment` = ? WHERE `Id` = ?", body.Rate, body.Comment, review.Id)
	}
	response := SetRateResponse{}
	response.Message = "リクエストを処理しました"
	response.Status = "success"
	json.NewEncoder(w).Encode(response)
}
