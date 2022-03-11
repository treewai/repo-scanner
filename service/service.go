package service

import (
	"secret-scanner/db"
	"secret-scanner/models"

	"github.com/gofiber/fiber/v2"
)

type Service interface {
	AddRepository(c *fiber.Ctx) error
	GetRepositories(c *fiber.Ctx) error
	GetRepositoryById(c *fiber.Ctx) error
	UpdateRepository(c *fiber.Ctx) error
	DeleteRepository(c *fiber.Ctx) error

	ScanRepository(c *fiber.Ctx) error
	GetScanResults(c *fiber.Ctx) error

	UpdateScanResult(result *models.Result) error
}

type service struct {
	db       db.Database
	register RegisterFunc
}

func NewService(db db.Database, f RegisterFunc) Service {
	return &service{
		db:       db,
		register: f,
	}
}

type RegisterFunc func(*models.Result) error

func Register(scanc chan<- *models.Result) RegisterFunc {
	return func(req *models.Result) error {
		scanc <- req
		return nil
	}
}
