package restapi

import (
	"net/http"

	"github.com/north-fy/golang-restapi-todolist/internal/config"
	statshandler "github.com/north-fy/golang-restapi-todolist/internal/handler/stats"
	taskhandler "github.com/north-fy/golang-restapi-todolist/internal/handler/task"
	userhandler "github.com/north-fy/golang-restapi-todolist/internal/handler/user"
	"github.com/north-fy/golang-restapi-todolist/internal/service/stats"
	"github.com/north-fy/golang-restapi-todolist/internal/service/task"
	"github.com/north-fy/golang-restapi-todolist/internal/service/user"
	"github.com/north-fy/golang-restapi-todolist/internal/storage/postgres"
	redisch "github.com/north-fy/golang-restapi-todolist/internal/storage/redis"
	"github.com/sirupsen/logrus"
)

type RestAPIServer struct {
	log    *logrus.Logger
	router *Route
}

func NewRestAPIServer(log *logrus.Logger, cfg config.StorageConfig, cfgCache config.RedisConfig) *RestAPIServer {
	storage := postgres.NewStorage(cfg)

	serv := RestAPIServer{
		router: newRouter(),
	}

	redis := redisch.NewStorage(cfgCache)

	userHandler := userhandler.NewHandlerUser(log, user.NewServiceUser(log, storage))
	taskHandler := taskhandler.NewHandlerTask(log, task.NewServiceTask(log, storage))
	statsHandler := statshandler.NewHandlerStats(log, stats.NewServiceStats(log, storage, redis))

	serv.router.ConfigureRouter(userHandler, taskHandler, statsHandler)

	return &serv
}

func (r RestAPIServer) Run(addr string) {
	if err := http.ListenAndServe(addr, r.router.mu); err != nil {
		panic(err)
	}
}
