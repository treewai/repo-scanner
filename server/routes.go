package server

import "github.com/gofiber/fiber/v2/middleware/logger"

func (s *server) registerServices() {
	s.app.Use(logger.New())

	v1 := s.app.Group("/v1")
	{
		v1.Post("/repositories", s.service.AddRepository)
		v1.Get("/repositories", s.service.GetRepositories)
		v1.Get("/repositories/:repoId", s.service.GetRepositoryById)
		v1.Patch("/repositories/:repoId", s.service.UpdateRepository)
		v1.Delete("/repositories/:repoId", s.service.DeleteRepository)

		v1.Post("/scans/:repoId", s.service.ScanRepository)
		v1.Get("/scans", s.service.GetScanResults)
	}
}
