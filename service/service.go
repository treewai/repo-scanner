package service

import (
	"secret-scanner/db"

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
}

type service struct {
	db db.DB
}

func NewService(db db.DB) Service {
	return &service{
		db: db,
	}
}
