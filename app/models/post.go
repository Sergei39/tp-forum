package models

type Post struct {
	Id       int64  `json:"id"`
	Parent   int64  `json:"parent"`
	Author   string `json:"author"`
	Message  string `json:"message"`
	IsEdited bool   `json:"isEdited"`
	Forum    string `json:"forum"`
	Thread   int    `json:"thread"`
	Created  string `json:"created"`
}

type RequestPost struct {
	Id      int    `json:"id"`
	Related string `json:"related"`
}

type InfoPost struct {
	Post   Post   `json:"post"`
	User   User   `json:"user"`
	Forum  Forum  `json:"forum"`
	Thread Thread `json:"thread"`
}

type MessagePostRequest struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
}
