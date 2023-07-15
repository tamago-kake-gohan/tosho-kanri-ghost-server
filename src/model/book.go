package model

type Book struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	ISBN   string `json:"isbn"`
	Author string `json:"author"`
	CCode  string `json:"c_code"`
	State  string `json:"state"` //available, lending, unavailable
}
