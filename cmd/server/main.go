package main

import (
	"fmt"

	"github.com/north-fy/golang-restapi-todolist/internal/app/restapi"
	"github.com/north-fy/golang-restapi-todolist/internal/config"
	"github.com/sirupsen/logrus"
)

const (
	envName = ".env"
)

func main() {
	cfg := config.MustLoadConfig(envName)
	log := logrus.New()

	serv := restapi.NewRestAPIServer(log, cfg.StorageCfg, cfg.RedisCfg)
	addr := fmt.Sprintf("%s:%d", cfg.ServerCfg.Host, cfg.ServerCfg.Port)

	log.Infof("Server is created on %s:%d", cfg.ServerCfg.Host, cfg.ServerCfg.Port)
	serv.Run(addr)
}
