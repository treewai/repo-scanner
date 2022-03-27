package server

import (
	"net"
	"sync"
	"testing"
	"time"
)

const listenURL = ":9999"

func getServer() Server {
	return NewServer(nil, &Config{
		ClientURL: listenURL,
	})
}

func TestServerNew(t *testing.T) {
	s := getServer()
	if s == nil {
		t.Error("failed to start server")
	}
	if _, ok := s.(*server); !ok {
		t.Error("failed to start server type")
	}
}

func TestServerStartSuccess(t *testing.T) {
	s := getServer()

	go func() {
		<-time.After(time.Second)
		conn, err := net.DialTimeout("tcp", listenURL, time.Second)
		if err != nil {
			t.Errorf("failed to serve on %v, err: %v", listenURL, err)
		}
		conn.Close()
		s.Stop()
	}()

	if err := s.Start(); err != nil {
		t.Errorf("failed to start server, err: %v", err)
	}
}

func TestServerGracefullStopSuccess(t *testing.T) {
	s := getServer()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		wg.Done()
		if err := s.Start(); err != nil {
			t.Errorf("failed to start server, err: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		wg.Done()
		select {
		case <-s.StopNotify():
		case <-time.After(5 * time.Second):
			t.Errorf("failed to gracefully stop server")
		}
	}()

	wg.Wait()
	s.Stop()
}
