package models

type Thread struct {
	Id      int32  `json:"id"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Forum   string `json:"forum"`
	Message string `json:"message"`
	Votes   int32  `json:"votes"`
	Slug    string `json:"slug"`
	Created string `json:"created"`
}
