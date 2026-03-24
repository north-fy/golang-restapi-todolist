package main

import (
	"fmt"

	"github.com/north-fy/golang-restapi-todolist/internal/app/restapi"
	"github.com/north-fy/golang-restapi-todolist/internal/config"
	"github.com/sirupsen/logrus"
)

const (
	pathToConfig string = "./config/server/config.yaml"
	serverHost   string = "localhost"
)

func main() {
	cfg := config.MustLoadConfig(pathToConfig)
	log := logrus.New()

	serv := restapi.NewRestAPIServer(log, cfg.StorageCfg)
	addr := fmt.Sprintf("%s:%d", serverHost, cfg.ServerCfg.Port)

	log.Info("Server is created on localhost:8080")
	serv.Run(addr)
}
