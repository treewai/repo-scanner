package service

import (
	"secret-scanner/models"

	"github.com/gofiber/fiber/v2"
)

// AddRepository add the given repository to database
func (s *service) AddRepository(c *fiber.Ctx) error {
	if string(c.Request().Header.ContentType()) != fiber.MIMEApplicationJSON {
		return c.SendStatus(fiber.StatusUnsupportedMediaType)
	}
	var repo models.Repository
	if err := c.BodyParser(&repo); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if repo.Name == "" || repo.Url == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if err := s.db.AddRepository(&repo); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusCreated)
}

// GetRepositories search and return all repositories
func (s *service) GetRepositories(c *fiber.Ctx) error {
	repos, err := s.db.GetRepositories()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(repos)
}

// GetRepositoryById search and return repository by given repoId
func (s *service) GetRepositoryById(c *fiber.Ctx) error {
	repo, err := s.db.GetRepositoryById(c.Params("repoId"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if repo == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	return c.Status(fiber.StatusOK).JSON(repo)
}

// UpdateRepository updates the given repository to database
func (s *service) UpdateRepository(c *fiber.Ctx) error {
	if string(c.Request().Header.ContentType()) != fiber.MIMEApplicationJSON {
		return c.SendStatus(fiber.StatusUnsupportedMediaType)
	}
	var repo models.Repository
	if err := c.BodyParser(&repo); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	repoId := c.Params("repoId")
	if repoId == "" || (repo.Name == "" && repo.Url == "") {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	repo.ID = repoId
	if err := s.db.UpdateRepository(&repo); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusOK)
}

// DeleteRepository deletes the given repoId from datatabase
func (s *service) DeleteRepository(c *fiber.Ctx) error {
	if err := s.db.DeleteRepository(c.Params("repoId")); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusOK)
}
