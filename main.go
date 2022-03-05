package main

import (
	"log"
	"os"
	"os/signal"
	"secret-scanner/db"
	"secret-scanner/server"
	"syscall"
)

func main() {
	db, err := db.NewClient(&db.Config{
		Hostname: "fullstack-postgres",
		Port:     "5432",
		Name:     "fullstack_api",
		User:     "steven",
		Password: "password",
	})
	if err != nil {
		log.Fatalln(err)
	}

	srvCfg := &server.Config{
		Name:      "SecretScanner",
		ClientURL: ":8080",
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
