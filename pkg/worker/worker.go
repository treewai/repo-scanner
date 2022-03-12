package worker

import (
	"secret-scanner/pkg/scanner"
)

type Worker interface {
	Start()
	Stop()
	Do(j *scanner.Job)
}

type worker struct {
	cfg      *Config
	scanners []scanner.Scanner
	jobc     chan *scanner.Job
	done     chan struct{}
}

func NewWorker(cfg *Config) Worker {
	var scanners []scanner.Scanner
	for i := 0; i < cfg.MaxWorker; i++ {
		scanners = append(scanners, scanner.NewScanner(cfg.WorkerDir))
	}

	return &worker{
		cfg:      cfg,
		scanners: scanners,
		jobc:     make(chan *scanner.Job, cfg.MaxWorker),
		done:     make(chan struct{}),
	}
}

func (w *worker) Start() {
	for _, sc := range w.scanners {
		go func(s scanner.Scanner) {
			for {
				select {
				case j := <-w.jobc:
					s.Scan(j)
				case <-w.done:
				}
			}
		}(sc)
	}
}

func (w *worker) Stop() {
	close(w.done)
}

func (w *worker) Do(j *scanner.Job) {
	select {
	case w.jobc <- j:
	case <-j.Ctx.Done():
	case <-w.done:
	}
}
