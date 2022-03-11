package server

import (
	"secret-scanner/db"
	"secret-scanner/models"
	"secret-scanner/pkg/scanner"
	"secret-scanner/service"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
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

	scanner scanner.Scanner
	scanc   chan *models.Result

	// a channel to listen stop message
	stopc chan struct{}
	// a channel to notify done serving
	done chan struct{}
}

func NewServer(db db.Database, cfg *Config) Server {
	fcfg := fiber.Config{
		AppName:      ServerName,
		ServerHeader: cfg.Name,
		Prefork:      cfg.Prefork,
	}
	scfg := &scanner.Config{
		Path: cfg.RepoDir,
	}
	scanc := make(chan *models.Result, cfg.ScanQueue)

	s := &server{
		cfg:     cfg,
		app:     fiber.New(fcfg),
		service: service.NewService(db, service.Register(scanc)),
		scanner: scanner.NewScanner(scfg),
		scanc:   scanc,
		stopc:   make(chan struct{}),
		done:    make(chan struct{}),
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

	go s.runScanner()

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

func (s *server) runScanner() {
	for {
		select {
		case result := <-s.scanc:
			s.scan(result)
		case <-s.done:
			return
		}
	}
}

func (s *server) scan(result *models.Result) {
	result.Status = models.StatusInProgress
	result.ScanningAt = time.Now()

	s.service.UpdateScanResult(result)

	j := &scanner.Job{
		Repo: result.Repository,
		Done: make(chan struct{}, 1),
	}
	s.scanner.Scan(j)

	select {
	case <-j.Done:
		if j.Err != nil {
			result.Status = models.StatusFailure
		} else {
			result.Status = models.StatusSuccess
			result.Findings = datatypes.JSONMap{
				"findings": j.Findings,
			}
		}
		result.FinishedAt = time.Now()
		s.service.UpdateScanResult(result)
	case <-s.done:
	}
}
