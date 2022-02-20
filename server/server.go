package server

import (
	"secret-scanner/db"
	"secret-scanner/service"

	"github.com/gofiber/fiber/v2"
)

const (
	ServerName = "Secret-Scanner"
)

type Server interface {
	Start() error
	Stop()
	StopNotify() <-chan struct{}
}

type server struct {
	cfg     *Config
	app     *fiber.App
	service service.Service

	// a channel to listen stop message
	stopc chan struct{}
	// a channel to notify done serving
	done chan struct{}
}

func NewServer(db db.DB, cfg *Config) Server {
	s := &server{
		cfg: cfg,
		app: fiber.New(fiber.Config{
			AppName:      ServerName,
			ServerHeader: cfg.Name,
			Prefork:      cfg.Prefork,
		}),
		service: service.NewService(db),
	}

	s.registerServices()
	return s
}

// Start starts the server and listen on given url
func (s *server) Start() error {
	go func() {
		select {
		case <-s.stopc:
			s.app.Shutdown()
			close(s.done)
		case <-s.done:
		}
	}()

	return s.app.Listen(s.cfg.ClientURL)
}

// Stop stops the server gracefully
func (s *server) Stop() {
	select {
	case s.stopc <- struct{}{}:
	case <-s.done:
	}
}

// StopNotify return done channel
func (s *server) StopNotify() <-chan struct{} {
	return s.done
}
