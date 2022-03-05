package db

import (
	"fmt"
	"secret-scanner/models"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database interface {
	AddRepository(repos []*models.Repository) error
	GetRepositories() ([]*models.Repository, error)
	GetRepositoryById(repoId string) (*models.Repository, error)
	UpdateRepository(repo *models.Repository) error
	DeleteRepository(repoId string) error
}

type database struct {
	conn *gorm.DB
}

func NewClient(cfg *Config) (Database, error) {
	time.Sleep(10 * time.Second)
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
