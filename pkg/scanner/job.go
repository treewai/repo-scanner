package scanner

import (
	"secret-scanner/models"
)

type Job struct {
	Repo   *models.Repository
	Result *models.Result
	Err    error
	Done   chan struct{}
}
