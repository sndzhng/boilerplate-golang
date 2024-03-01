package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sndzhng/gin-template/internal/config"
	"github.com/sndzhng/gin-template/internal/controller/route"
	"github.com/sndzhng/gin-template/internal/datastore"
)

func main() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	config.InitialConfig(os.Args)
	config.InitialTimeZone()

	// datastore.ConnectCloudStorage()
	// datastore.ConnectMongodb()
	// defer datastore.DisconnectMongodb()
	datastore.ConnectPostgresql()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Server.Port),
		Handler: route.SetupRouter(),
	}

	go startServer(server)
	shutdownServer(server)
}

func startServer(server *http.Server) {
	log.Printf("\nRun on %s environment", config.Environment)
	log.Printf("Listening and serving HTTP on %s\n", config.Server.Port)

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error listening and serving HTTP on %s\n", err)
	}
}

func shutdownServer(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Print("Shutdown Server ...")

	contextTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := server.Shutdown(contextTimeout)
	if err != nil {
		log.Fatal("Server Shutdown ", err)
	}
}
