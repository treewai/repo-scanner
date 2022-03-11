package models

import "time"

type Repository struct {
	ID         string    `json:"repoId" gorm:"default:uuid_generate_v3()"`
	Name       string    `json:"name"`
	Url        string    `json:"url"`
	CreatedAt  time.Time `json:"-"`
	ModifiedAt time.Time `json:"-"`
}
