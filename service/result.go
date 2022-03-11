package service

import (
	"secret-scanner/models"

	"github.com/gofiber/fiber/v2"
)

func (s *service) GetScanResults(c *fiber.Ctx) error {
	results, err := s.db.GetResults()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(results)
}

func (s *service) UpdateScanResult(result *models.Result) error {
	if err := s.db.UpdateResult(result); err != nil {
		return err
	}
	return nil
}
