package scanner

import (
	"secret-scanner/models"

	"gorm.io/datatypes"
)

type Job struct {
	Repo     *models.Repository
	Findings []datatypes.JSONMap
	Err      error
	Done     chan struct{}
}
