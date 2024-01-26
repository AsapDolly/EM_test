package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "net/http/pprof"

	"github.com/AsapDolly/EM_test/internal/server"
	_ "github.com/lib/pq"
)

func main() {
	var dbConnection *sql.DB
	defer dbConnection.Close()

	cfg, router := server.GetServer(dbConnection)

	var srv = &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: router,
	}

	idleConnsClosed := make(chan struct{})

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		<-sigint
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
	<-idleConnsClosed
	os.Exit(0)

}
