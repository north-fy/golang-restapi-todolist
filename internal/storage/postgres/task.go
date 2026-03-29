package postgres

import (
	"context"

	"github.com/north-fy/golang-restapi-todolist/internal/domain/models"
	"github.com/north-fy/golang-restapi-todolist/internal/handler/task"
	"github.com/pkg/errors"
)

/*
type StorageTask interface {
	InsertTask(ctx context.Context, task models.Task) (int, error)
	SelectTask(ctx context.Context, taskID int) (models.Task, error)
	SelectTasksWithPagination(ctx context.Context, pt task.PaginationTask) ([]models.Task, error)
	SelectTasksByUser(ctx context.Context, userID int) ([]models.Task, error)
	UpdateTask(ctx context.Context, task models.Task) error
	DeleteTask(ctx context.Context, taskID int) error
}
*/

func (s *Storage) InsertTask(ctx context.Context, task models.Task) (int, error) {
	query := `
	INSERT INTO task (title, description, completed, user_id)
	VALUES ($1, $2, false, $3)
	RETURNING id
	`

	var id int
	if err := s.conn.QueryRow(ctx, query, task.Title, task.Description, task.UserID).Scan(&id); err != nil {
		return 0, errors.Errorf("%s: %s", op, err.Error())
	}

	return id, nil
}

func (s *Storage) SelectTask(ctx context.Context, taskID int) (models.Task, error) {
	query := `
	SELECT *
	FROM task
	WHERE id = $1
	`

	// TODO: Исправить скан тупица :D
	
	oneTask := models.Task{}
	if err := s.conn.QueryRow(ctx, query, taskID).Scan(&oneTask); err != nil {
		return models.Task{}, errors.Errorf("%s: %s", op, err.Error())
	}

	return oneTask, nil
}

func (s *Storage) SelectTasksWithPagination(ctx context.Context, pt task.PaginationTask) ([]models.Task, error) {
	query := `
	SELECT *
	FROM task
	ORDER BY id
	LIMIT $1
	OFFSET $2 
	`

	tasks := make([]models.Task, pt.Limit)
	rows, err := s.conn.Query(ctx, query, pt.Limit, pt.Offset)
	if err != nil {
		return nil, errors.Errorf("%s: %s", op, err.Error())
	}

	for rows.Next() {
		var oneTask models.Task
		if err := rows.Scan(&oneTask); err != nil {
			return tasks, errors.Errorf("%s: %s", op, err.Error())
		}

		tasks = append(tasks, oneTask)
	}

	return tasks, nil
}

func (s *Storage) SelectTasksByUser(ctx context.Context, userID int) ([]models.Task, error) {
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
	rows, err := s.conn.Query(ctx, query, userID)
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

func (s *Storage) UpdateTask(ctx context.Context, task models.Task) error {
	query := `
	UPDATE task
	SET title = COALESCE(NULLIF($1, ''), title),
	    description = COALESCE(NULLIF($2, ''), description),
	    completed = COALESCE(NULLIF($3, false), completed),
	    completed_at = COALESCE(NULLIF($3, false), now())
	WHERE id = $4
	`

	_, err := s.conn.Exec(ctx, query, task.Title, task.Description, task.Completed, task.ID)
	if err != nil {
		return errors.Errorf("%s: %s", op, err.Error())
	}

	return nil
}

func (s *Storage) DeleteTask(ctx context.Context, taskID int) error {
	query := `
	DELETE FROM task
	WHERE id = $1
	`

	_, err := s.conn.Exec(ctx, query, taskID)
	if err != nil {
		return errors.Errorf("%s: %s", op, err.Error())
	}

	return nil
}
