package models

type Task struct {
	Id      string `json:"id,omitempty"`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title,omitempty"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

type Tasks struct {
	Tasks []Task `json:"tasks"`
}

type AddId struct {
	Id int64 `json:"id,omitempty"`
}

type Err struct {
	Err string `json:"error,omitempty"`
}

type Authentication struct {
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}
