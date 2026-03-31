package task

import (
	"context"

	"github.com/north-fy/golang-restapi-todolist/internal/domain/models"
	"github.com/north-fy/golang-restapi-todolist/internal/handler/task"
	"github.com/north-fy/golang-restapi-todolist/pkg/validate"
	"github.com/sirupsen/logrus"
)

/*
type ServiceTask interface {
	CreateTask(ctx context.Context, task models.Task) (int, error)
	GetTask(ctx context.Context, taskID int) (models.Task, error)
	GetTasksWithPagination(ctx context.Context, pt task.PaginationTask) ([]models.Task, error)
	GetTasksByUser(ctx context.Context, userID int) ([]models.Task, error)
	EditTask(ctx context.Context, task models.Task) error
	DeleteTask(ctx context.Context, taskID int) error
}
*/

const op = "service/task/task"

type StorageTask interface {
	InsertTask(ctx context.Context, task models.Task) (int, error)
	SelectTask(ctx context.Context, taskID int) (models.Task, error)
	SelectTasksWithPagination(ctx context.Context, pt task.PaginationTask) ([]models.Task, error)
	SelectTasksByUser(ctx context.Context, userID int) ([]models.Task, error)
	UpdateTask(ctx context.Context, task models.Task) error
	DeleteTask(ctx context.Context, taskID int) error
}

type ServiceTask struct {
	log     *logrus.Logger
	storage StorageTask
}

func NewServiceTask(log *logrus.Logger, storage StorageTask) *ServiceTask {
	return &ServiceTask{
		log:     log,
		storage: storage,
	}
}

func (s *ServiceTask) CreateTask(ctx context.Context, task models.Task) (int, error) {
	if err := validate.OptValidate(task.Title, true, 1, 100); err != nil {
		s.log.Errorf("%s: %s", op, err.Error())
		return 0, err
	}

	if err := validate.OptValidate(task.Description, false, 1, 1000); err != nil {
		s.log.Errorf("%s: %s", op, err.Error())
		return 0, err
	}

	id, err := s.storage.InsertTask(ctx, task)
	if err != nil {
		s.log.Errorf("%s: %s", op, err.Error())
		return 0, err
	}

	return id, nil
}

func (s *ServiceTask) GetTask(ctx context.Context, taskID int) (models.Task, error) {
	oneTask, err := s.storage.SelectTask(ctx, taskID)
	if err != nil {
		s.log.Errorf("%s: %s", op, err.Error())
		return models.Task{}, err
	}

	return oneTask, nil
}

func (s *ServiceTask) GetTasksWithPagination(ctx context.Context, pt task.PaginationTask) ([]models.Task, error) {
	tasks, err := s.storage.SelectTasksWithPagination(ctx, pt)
	if err != nil {
		s.log.Errorf("%s: %s", op, err.Error())
		return nil, err
	}

	return tasks, nil
}

func (s *ServiceTask) GetTasksByUser(ctx context.Context, userID int) ([]models.Task, error) {
	tasks, err := s.storage.SelectTasksByUser(ctx, userID)
	if err != nil {
		s.log.Errorf("%s: %s", op, err.Error())
		return nil, err
	}

	return tasks, nil
}

func (s *ServiceTask) EditTask(ctx context.Context, task models.Task) error {
	if err := validate.OptValidate(task.Title, false, 1, 100); err != nil {
		s.log.Errorf("%s: %s", op, err.Error())
		return err
	}

	if err := validate.OptValidate(task.Description, false, 1, 1000); err != nil {
		s.log.Errorf("%s: %s", op, err.Error())
		return err
	}

	if err := s.storage.UpdateTask(ctx, task); err != nil {
		s.log.Errorf("%s: %s", op, err.Error())
		return err
	}

	return nil
}

func (s *ServiceTask) DeleteTask(ctx context.Context, taskID int) error {
	if err := s.storage.DeleteTask(ctx, taskID); err != nil {
		s.log.Errorf("%s: %s", op, err.Error())
		return err
	}

	return nil
}
