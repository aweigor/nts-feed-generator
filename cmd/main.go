package main

import (
	"fmt"
	"log"
	"net/http"

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
	fmt.Println("Server is listening on port 9000")
	server.ListenAndServe()
}
