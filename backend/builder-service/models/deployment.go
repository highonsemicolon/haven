package models

type Builder struct {
	Name      string `json:"name" gorm:"unique" notnull:"true"`
	GitURL    string `json:"git_url" notnull:"true"`
	Branch    string `json:"branch" null:"true"`
	HostedURL string `json:"hosted_url" gorm:"default:null"`
	Comment   string `json:"comment" gorm:"default:null"`
}
