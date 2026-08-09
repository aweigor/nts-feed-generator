package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/aweigor/nts-feed-generator/config"
	"github.com/aweigor/nts-feed-generator/internal/heartbeat"
	"github.com/aweigor/nts-feed-generator/pkg/middleware"
)

func NewApp() http.Handler {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	router := http.NewServeMux()
	heartbeat.NewHeartbeatHandler(router, conf)

	mwStack := middleware.Chain(middleware.CORS, middleware.Logging)

	return mwStack(router)
}

func main() {
	app := NewApp()
	server := http.Server{Addr: ":9000", Handler: app}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		fmt.Println("Server is up on :9000")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	<-quit
	fmt.Println("Shutting down...")
	server.Shutdown(context.Background())
}
