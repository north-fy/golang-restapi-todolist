package task

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/north-fy/golang-restapi-todolist/internal/domain/models"
	"github.com/north-fy/golang-restapi-todolist/pkg/write"
	"github.com/sirupsen/logrus"
)

const op = "handler/task/task"

type PaginationTask struct {
	Limit  int
	Offset int
}

type ServiceTask interface {
	CreateTask(ctx context.Context, task models.Task) (int, error)
	GetTask(ctx context.Context, taskID int) (models.Task, error)
	GetTasksWithPagination(ctx context.Context, pt PaginationTask) ([]models.Task, error)
	GetTasksByUser(ctx context.Context, userID int) ([]models.Task, error)
	EditTask(ctx context.Context, task models.Task) error
	DeleteTask(ctx context.Context, taskID int) error
}

type HandlerTask struct {
	log     *logrus.Logger
	service ServiceTask
}

func NewHandlerTask(log *logrus.Logger, serv ServiceTask) *HandlerTask {
	return &HandlerTask{
		log:     log,
		service: serv,
	}
}

func (h *HandlerTask) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	task := models.Task{}
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		h.log.Errorf("%s: %s", op, err.Error())
		write.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.CreateTask(ctx, task)
	if err != nil {
		h.log.Errorf("%s: %s", op, err.Error())
		write.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	write.WriteJSON(w, http.StatusCreated, id)
}

func (h *HandlerTask) HandleGetTaskByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.log.Errorf("%s: %s", op, err.Error())
		write.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	task, err := h.service.GetTask(ctx, id)
	if err != nil {
		h.log.Errorf("%s: %s", op, err.Error())
		write.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	write.WriteJSON(w, http.StatusOK, task)
}

func (h *HandlerTask) HandleGetPaginationTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	if limitStr == "" || offsetStr == "" {
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		h.log.Errorf("%s: %s", op, err.Error())
		write.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		h.log.Errorf("%s: %s", op, err.Error())
		write.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	pt := PaginationTask{
		Limit:  limit,
		Offset: offset,
	}

	tasks, err := h.service.GetTasksWithPagination(ctx, pt)
	if err != nil {
		h.log.Errorf("%s: %s", op, err.Error())
		write.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	write.WriteJSON(w, http.StatusOK, tasks)
}

func (h *HandlerTask) HandleGetTasksByUserID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.log.Errorf("%s: %s", op, err.Error())
		write.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	tasks, err := h.service.GetTasksByUser(ctx, id)
	if err != nil {
		h.log.Errorf("%s: %s", op, err.Error())
		write.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	write.WriteJSON(w, http.StatusOK, tasks)
}

func (h *HandlerTask) HandleEditTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.log.Errorf("%s: %s", op, err.Error())
		write.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	task := models.Task{}
	if err = json.NewDecoder(r.Body).Decode(&task); err != nil {
		h.log.Errorf("%s: %s", op, err.Error())
		write.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	task.ID = id
	if err = h.service.EditTask(ctx, task); err != nil {
		h.log.Errorf("%s: %s", op, err.Error())
		write.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	write.WriteJSON(w, http.StatusOK, "success")
}

func (h *HandlerTask) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.log.Errorf("%s: %s", op, err.Error())
		write.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.service.DeleteTask(ctx, id); err != nil {
		h.log.Errorf("%s: %s", op, err.Error())
		write.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	write.WriteJSON(w, http.StatusOK, "success")
}
