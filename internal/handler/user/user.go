package user

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/north-fy/golang-restapi-todolist/internal/domain/models"
	"github.com/north-fy/golang-restapi-todolist/pkg/write"
	"github.com/sirupsen/logrus"
)

const op = "handler/user/user"

type ServiceUser interface {
	CreateUser(ctx context.Context, firstName, lastName, numberPhone string) (int, error)
	GetUser(ctx context.Context, id int) (models.User, error)
	GetTasksByUser(ctx context.Context, id int) ([]models.Task, error)
	EditInfoUser(ctx context.Context, user models.User) error
	DeleteUser(ctx context.Context, id int) error
}

type HandlerUser struct {
	log     *logrus.Logger
	service ServiceUser
}

func NewHandlerUser(log *logrus.Logger, user ServiceUser) *HandlerUser {
	return &HandlerUser{
		log:     log,
		service: user,
	}
}

// user created: first name, last name, number phone
func (h *HandlerUser) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	requestUser := models.User{}
	if err := json.NewDecoder(r.Body).Decode(&requestUser); err != nil {
		write.WriteError(w, http.StatusInternalServerError, "bad request for create user")
		h.log.Errorf("%s: %s", op, err.Error())
		return
	}

	id, err := h.service.CreateUser(context.TODO(), requestUser.FirstName, requestUser.LastName, requestUser.NumberPhone)
	if err != nil {
		// TODO: организовать обработку конкретных ошибок из сервиса
		write.WriteError(w, http.StatusInternalServerError, "bad request for create user")
		h.log.Errorf("%s: %s", op, err.Error())
		return
	}

	write.WriteJSON(w, http.StatusCreated, map[string]int{"id": id})
}

func (h *HandlerUser) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/users/"))
	if err != nil {
		write.WriteError(w, http.StatusBadRequest, "bad request for get user")
		h.log.Errorf("%s: %s", op, err.Error())
		return
	}

	userModel, err := h.service.GetUser(context.TODO(), userID)
	if err != nil {
		write.WriteError(w, http.StatusInternalServerError, "internal error")
		h.log.Errorf("%s: %s", op, err.Error())
		return
	}

	write.WriteJSON(w, http.StatusOK, userModel)
}

func (h *HandlerUser) HandleGetTasks(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		write.WriteError(w, http.StatusBadRequest, "bad request for get user")
		h.log.Errorf("%s: %s", op, err.Error())
		return
	}

	tasks, err := h.service.GetTasksByUser(context.TODO(), userID)
	if err != nil {
		write.WriteError(w, http.StatusInternalServerError, "internal error")
		h.log.Errorf("%s: %s", op, err.Error())
		return
	}

	write.WriteJSON(w, http.StatusOK, tasks)
}

func (h *HandlerUser) HandleEditUser(w http.ResponseWriter, r *http.Request) {
	requestUser := models.User{}
	if err := json.NewDecoder(r.Body).Decode(&requestUser); err != nil {
		write.WriteError(w, http.StatusInternalServerError, "bad request for create user")
		h.log.Errorf("%s: %s", op, err.Error())
		return
	}

	userID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		write.WriteError(w, http.StatusBadRequest, "bad request for get user")
		h.log.Errorf("%s: %s", op, err.Error())
		return
	}

	requestUser.ID = userID

	if err = h.service.EditInfoUser(context.TODO(), requestUser); err != nil {
		write.WriteError(w, http.StatusInternalServerError, "internal error")
		h.log.Errorf("%s: %s", op, err.Error())
		return
	}

	write.WriteJSON(w, http.StatusOK, nil)
}

func (h *HandlerUser) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		write.WriteError(w, http.StatusBadRequest, "bad request for get user")
		h.log.Errorf("%s: %s", op, err.Error())
		return
	}

	if err = h.service.DeleteUser(context.TODO(), userID); err != nil {
		write.WriteError(w, http.StatusInternalServerError, "internal error")
		h.log.Errorf("%s: %s", op, err.Error())
		return
	}

	write.WriteJSON(w, http.StatusOK, nil)
}
