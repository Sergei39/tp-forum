package models

type Vote struct {
	Id     int    `json:"id"`
	User   string `json:"nickname"`
	Thread int    `json:"thread"`
	Voice  int    `json:"voice"`
}
