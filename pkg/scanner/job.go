package scanner

import (
	"context"
	"secret-scanner/models"

	"gorm.io/datatypes"
)

type Job struct {
	Repo     *models.Repository
	Findings []datatypes.JSONMap
	Ctx      context.Context
	Err      error
	Done     chan struct{}
}
