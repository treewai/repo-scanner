package db

import (
	"fmt"
	"secret-scanner/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Database inteface for database
type Database interface {
	AddRepository(repo *models.Repository) error
	GetRepositories() ([]*models.Repository, error)
	GetRepositoryById(repoId string) (*models.Repository, error)
	UpdateRepository(repo *models.Repository) error
	DeleteRepository(repoId string) error

	AddResult(res *models.Result) (*models.Result, error)
	GetResults() ([]*models.Result, error)
	UpdateResult(res *models.Result) error
}

type database struct {
	conn *gorm.DB
}

// NewClient creates and returns database client instance
func NewClient(cfg *Config) (Database, error) {
	dsn := fmt.Sprintf("host=%s port=%s dbname=%s sslmode=disable user=%s password=%s",
		cfg.Hostname, cfg.Port, cfg.Name, cfg.User, cfg.Password)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &database{
		conn: db,
	}, nil
}
