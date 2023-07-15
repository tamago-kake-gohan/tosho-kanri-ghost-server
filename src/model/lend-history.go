package model

import "time"

type LendHistory struct {
	Id     int       `json:"id"`
	LendId int       `json:"lend_id"`
	Status string    `json:"status"` //requested, accepted, rejected, returned
	Date   time.Time `json:"date"`
}
