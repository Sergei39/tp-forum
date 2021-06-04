package models

type Thread struct {
	Id      int     `json:"id"`
	Title   string  `json:"title"`
	Author  string  `json:"author"`
	Forum   string  `json:"forum"`
	Message string  `json:"message"`
	Votes   *int    `json:"votes"`
	Slug    string  `json:"slug"`
	Created *string `json:"created"`
}

type ThreadPosts struct {
	SlugOrId string `json:"slug"`
	Limit    string `json:"limit"`
	Since    string `json:"since"`
	Sort     string `json:"sort"`
	Desc     bool   `json:"desc"`
	ThreadId int
}
