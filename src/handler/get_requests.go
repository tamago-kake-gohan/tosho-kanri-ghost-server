package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/astaxie/session"
	"tamago-kake-gohan.github.io/tosho-kanri-ghost/src/model"
)

type Request struct {
	Id           int    `json:"id"`
	Title        string `json:"title"`
	BorrowerName string `json:"borrower_name"`
	Status       string `json:"status"`
}

type GetRequestsResponse struct {
	Message string     `json:"message"`
	Status  string     `json:"status"`
	Data    []*Request `json:"data"`
}

type GetRequestsHandler struct {
	sess *session.Manager
	db   *sql.DB
}

func NewGetRequestsHandler(sess *session.Manager, db *sql.DB) *GetRequestsHandler {
	return &GetRequestsHandler{
		sess: sess,
		db:   db,
	}
}

func (h *GetRequestsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response := GetRequestsResponse{}
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
	rows, err := h.db.Query("SELECT `UserLendBook`.`Id`,`Book`.`Title`, `Borrower`.`Name` as `BorrowerName`,`UserLendBook`.`Status` FROM `UserLendBook` INNER JOIN `UserBook` ON `UserLendBook`.`UserBookId` = `UserBook`.`Id` INNER JOIN `Book` ON `UserBook`.`BookId` = `Book`.`Id` INNER JOIN `User` as `Borrower` ON `UserLendBook`.`BorrowerId` = `Borrower`.`Id`", userId)
	if nil != err {
		response.Message = "データの取得に失敗しました"
		response.Status = "error"
		json.NewEncoder(w).Encode(response)
		return
	}
	result := make([]*Request, 0)
	for rows.Next() {
		request := &Request{}
		if err := rows.Scan(&request.Id, &request.Title, &request.BorrowerName, &request.Status); err != nil {
			log.Printf("getRows rows.Scan error err:%v", err)
			continue
		}
		result = append(result, request)
	}
	response.Message = ""
	response.Status = "success"
	response.Data = result
	json.NewEncoder(w).Encode(response)
}
