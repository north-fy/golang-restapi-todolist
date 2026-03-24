package restapi

import (
	"net/http"

	"github.com/north-fy/golang-restapi-todolist/internal/config"
	userhandler "github.com/north-fy/golang-restapi-todolist/internal/handler/user"
	"github.com/north-fy/golang-restapi-todolist/internal/service/user"
	"github.com/north-fy/golang-restapi-todolist/internal/storage/postgres"
	"github.com/sirupsen/logrus"
)

type Storage struct {
	user user.StorageUser
	// task task.StorageTask
}

type RestAPIServer struct {
	log    *logrus.Logger
	router *Route
}

func NewRestAPIServer(log *logrus.Logger, cfg config.StorageConfig) *RestAPIServer {
	st := Storage{
		user: postgres.NewStorage(cfg),
	}
	//TODO: Implement !!
	_ = st

	serv := RestAPIServer{
		router: newRouter(),
	}

	userHandler := userhandler.NewHandlerUser(log, user.NewServiceUser(log, st.user))

	serv.router.ConfigureRouter(userHandler)

	return &serv
}

func (r RestAPIServer) Run(addr string) {
	if err := http.ListenAndServe(addr, r.router.mu); err != nil {
		panic(err)
	}
}
