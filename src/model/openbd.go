package model

type BookMeta struct {
	Summary struct {
		Title  string `json:"title"`
		Author string `json:"author"`
	}
}
