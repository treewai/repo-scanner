package models

import "time"

type Repository struct {
	ID         string    `json:"repoId" validate:"required" gorm:"default:uuid_generate_v3()"`
	Name       string    `json:"name" validate:"required"`
	Url        string    `json:"url" validate:"required"`
	CreatedAt  time.Time `json:"-"`
	ModifiedAt time.Time `json:"-"`
}
