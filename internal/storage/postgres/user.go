package postgres

import (
	"context"

	"github.com/north-fy/golang-restapi-todolist/internal/domain/models"
	"github.com/pkg/errors"
)

//go:generate mockery --name=StorageUser --filename=user_repository_mock.go --output=./mocks --outpkg=mocks --dry-run=false
type StorageUser interface {
	CreateUser(ctx context.Context, firstName, lastName, numberPhone string) (int, error)
	GetUser(ctx context.Context, id int) (models.User, error)
	GetTasks(ctx context.Context, id int) ([]models.Task, error)
	UpdateUser(ctx context.Context, user models.User) error
	DeleteUser(ctx context.Context, id int) error
}

func (s *Storage) CreateUser(ctx context.Context, firstName, lastName, numberPhone string) (int, error) {
	query := `
	INSERT INTO users(first_name, last_name, number_phone)
	VALUES ($1, $2, $3)
	RETURNING id
	`

	var id int
	if err := s.conn.QueryRow(ctx, query, firstName, lastName, numberPhone).Scan(&id); err != nil {
		return 0, errors.Errorf("%s: %s", op, err.Error())
	}

	return id, nil
}

func (s *Storage) GetUser(ctx context.Context, id int) (models.User, error) {
	query := `
	SELECT first_name, last_name, number_phone
	FROM users
	WHERE id = $1
	`

	user := models.User{}
	if err := s.conn.QueryRow(ctx, query, id).Scan(&user.FirstName, &user.LastName, &user.NumberPhone); err != nil {
		return models.User{}, errors.Errorf("%s: %s", op, err.Error())
	}

	return user, nil
}

func (s *Storage) GetTasks(ctx context.Context, id int) ([]models.Task, error) {
	query := `
	SELECT
		task.id, task.user_id, task.title, task.description, 
		task.completed, task.created_at, task.completed_at
	FROM
	    users INNER JOIN task
		ON users.id = task.user_id
	WHERE users.id = $1
	`

	tasks := []models.Task{}
	rows, err := s.conn.Query(ctx, query, id)
	if err != nil {
		return nil, errors.Errorf("%s: %s", op, err.Error())
	}

	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Description,
			&task.Completed, &task.CreatedAt, &task.CompletedAt); err != nil {
			return nil, errors.Errorf("%s: %s", op, err.Error())
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *Storage) UpdateUser(ctx context.Context, user models.User) error {
	query := `
	UPDATE users
	SET first_name = COALESCE(NULLIF($1, ''), first_name),
	    last_name = COALESCE(NULLIF($2, ''), last_name),
	    number_phone = COALESCE(NULLIF($3, ''), number_phone)
	WHERE id = $4
	`

	ct, err := s.conn.Exec(ctx, query, user.FirstName, user.LastName, user.NumberPhone, user.ID)
	if err != nil {
		return errors.Errorf("%s: %s", op, err.Error())
	}

	if count := ct.RowsAffected(); count == 0 {
		return errors.Errorf("%s: %s", op, "id not found")
	}

	return nil
}

func (s *Storage) DeleteUser(ctx context.Context, id int) error {
	query := `
	DELETE FROM users
	WHERE id = $1
	`

	ct, err := s.conn.Exec(ctx, query, id)
	if err != nil {
		return errors.Errorf("%s: %s", op, err.Error())
	}

	if count := ct.RowsAffected(); count == 0 {
		return errors.Errorf("%s: %s", op, "id not found")
	}

	return nil
}
