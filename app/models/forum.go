package models

type Forum struct {
	Title   string `json:"title"`
	User    string `json:"user"`
	Slug    string `json:"slug"`
	Posts   int    `json:"posts"`
	Threads int    `json:"threads"`
}

type ForumUsers struct {
	Slug  string `json:"slug"`
	Limit string `json:"limit"`
	Since string `json:"since"`
	Desc  bool   `json:"desc"`
}

type ForumThreads struct {
	Slug  string `json:"slug"`
	Limit int    `json:"limit"`
	Since string `json:"since"`
	Desc  bool   `json:"desc"`
}
