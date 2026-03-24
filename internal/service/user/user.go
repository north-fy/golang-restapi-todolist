package user

import (
	"context"

	"github.com/north-fy/golang-restapi-todolist/internal/domain/models"
	"github.com/north-fy/golang-restapi-todolist/pkg/validate"
	"github.com/sirupsen/logrus"
)

const op = "service/user/user"

/*
	type ServiceUser interface {
		CreateUser(ctx context.Context, firstName, lastName, numberPhone string) (int, error)
		GetUser(ctx context.Context, id int) (models.User, error)
		GetTasksByUser(ctx context.Context, id int) ([]models.Task, error)
		EditInfoUser(ctx context.Context, user models.User) error
		DeleteUser(ctx context.Context, id int) error
	}
*/
type StorageUser interface {
	CreateUser(ctx context.Context, firstName, lastName, numberPhone string) (int, error)
	GetUser(ctx context.Context, id int) (models.User, error)
	GetTasks(ctx context.Context, id int) ([]models.Task, error)
	UpdateUser(ctx context.Context, user models.User) error
	DeleteUser(ctx context.Context, id int) error
}

type ServiceUser struct {
	log     *logrus.Logger
	storage StorageUser
}

func NewServiceUser(log *logrus.Logger, storage StorageUser) *ServiceUser {
	return &ServiceUser{
		log:     log,
		storage: storage,
	}
}

func (s *ServiceUser) CreateUser(ctx context.Context, firstName, lastName, numberPhone string) (int, error) {
	if err := ValidateInfo(models.User{FirstName: firstName, LastName: lastName, NumberPhone: numberPhone}); err != nil {
		s.log.Errorf("%s: %s", op, err.Error())
		return 0, err
	}

	id, err := s.storage.CreateUser(ctx, firstName, lastName, numberPhone)
	if err != nil {
		s.log.Errorf("%s: %s", op, err.Error())
		return 0, err
	}

	return id, nil
}

func (s *ServiceUser) GetUser(ctx context.Context, id int) (models.User, error) {
	user, err := s.storage.GetUser(ctx, id)
	if err != nil {
		s.log.Errorf("%s: %s", op, err.Error())
		return models.User{}, err
	}

	return user, nil
}

func (s *ServiceUser) GetTasksByUser(ctx context.Context, id int) ([]models.Task, error) {
	tasks, err := s.storage.GetTasks(ctx, id)
	if err != nil {
		s.log.Errorf("%s: %s", op, err.Error())
		return nil, err
	}

	return tasks, nil
}

func (s *ServiceUser) EditInfoUser(ctx context.Context, user models.User) error {
	if err := ValidateInfo(user); err != nil {
		s.log.Errorf("%s: %s", op, err.Error())
		return err
	}

	if err := s.storage.UpdateUser(ctx, user); err != nil {
		s.log.Errorf("%s: %s", op, err.Error())
		return err
	}

	return nil
}

func (s *ServiceUser) DeleteUser(ctx context.Context, id int) error {
	if err := s.storage.DeleteUser(ctx, id); err != nil {
		s.log.Errorf("%s: %s", op, err.Error())
		return err
	}

	return nil
}

func ValidateInfo(u models.User) error {
	if err := validate.ValidateName(u.FirstName); err != nil {
		return err
	}

	if err := validate.ValidateName(u.LastName); err != nil {
		return err
	}

	if err := validate.ValidateNumber(u.NumberPhone); err != nil {
		return err
	}

	return nil
}
