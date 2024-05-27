package models

type Repo struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	GitURL  string `json:"git_url"`
	Comment string `json:"comment"`
}
