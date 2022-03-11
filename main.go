package main

import (
	"log"
	"os"
	"os/signal"
	"secret-scanner/db"
	"secret-scanner/server"
	"strconv"
	"syscall"
)

func main() {
	db, err := db.NewClient(&db.Config{
		Hostname: os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Name:     os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		log.Fatalln(err)
	}

	max, err := strconv.Atoi(os.Getenv("APP_SCAN_QUEUE"))
	if err != nil {
		log.Fatalf("invalid APP_SCAN_QUEUE, err: %v", err)
	}
	srvCfg := &server.Config{
		Name:      os.Getenv("APP_NAME"),
		ClientURL: os.Getenv("APP_LISTEN_URL"),
		RepoDir:   os.Getenv("APP_REPO_DIR"),
		ScanQueue: max,
	}
	srv := server.NewServer(db, srvCfg)

	signalc := make(chan os.Signal, 1)
	signal.Notify(signalc, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-signalc
		srv.Stop()
	}()

	srv.Start()

	// Wait server on shutting down
	<-srv.StopNotify()
}
