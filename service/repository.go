package service

import (
	"secret-scanner/models"

	"github.com/gofiber/fiber/v2"
)

func (s *service) AddRepository(c *fiber.Ctx) error {
	var repo models.Repository
	if err := c.BodyParser(&repo); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	repos := []*models.Repository{
		&repo,
	}
	if err := s.db.AddRepository(repos); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusCreated)
}

func (s *service) GetRepositories(c *fiber.Ctx) error {
	repos, err := s.db.GetRepositories()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(repos)
}

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

func (s *service) UpdateRepository(c *fiber.Ctx) error {
	var repo models.Repository
	if err := c.BodyParser(&repo); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	repo.ID = c.Params("repoId")
	if err := s.db.UpdateRepository(&repo); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusOK)
}

func (s *service) DeleteRepository(c *fiber.Ctx) error {
	if err := s.db.DeleteRepository(c.Params("repoId")); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusOK)
}
