package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/north-fy/golang-restapi-todolist/internal/config"
)

const (
	pathToConfig string = "./config/server/config.yaml"
	serverHost   string = "localhost"
)

func main() {
	cfg := config.MustLoadConfig(pathToConfig)

	server := http.Server{
		Addr: fmt.Sprintf("%s:%d", serverHost, cfg.ServerCfg.Port),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
