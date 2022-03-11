package service

import (
	"secret-scanner/models"

	"github.com/gofiber/fiber/v2"
)

func (s *service) ScanRepository(c *fiber.Ctx) error {
	repo, err := s.db.GetRepositoryById(c.Params("repoId"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if repo == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	res := models.Result{
		Status: models.StatusQueued,
		RepoID: repo.ID,
		//Repository: repo,
	}
	result, err := s.db.AddResult(&res)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	result.Repository = repo

	if err := s.register(result); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusOK)
}
