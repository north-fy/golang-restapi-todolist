package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/north-fy/golang-restapi-todolist/internal/domain/models"
	"github.com/pkg/errors"
)

//go:generate mockery --name=StorageUser --filename=user_repository_mock.go --output=./mocks --outpkg=mocks --dry-run=false
type StorageUser interface {
	CreateUser(ctx context.Context, firstName, lastName, numberPhone string) (int, error)
	GetUser(ctx context.Context, id int) (models.User, error)
	GetUsersWithPagination(ctx context.Context, pt models.Pagination) ([]models.User, error)
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
		return 0, fmt.Errorf("%s: %s", op, err.Error())
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
		if errors.As(err, &sql.ErrNoRows) {
			return models.User{}, models.ErrNoRows
		}

		return models.User{}, fmt.Errorf("%s: %s", op, err.Error())
	}

	return user, nil
}

func (s *Storage) GetUsersWithPagination(ctx context.Context, pt models.Pagination) ([]models.User, error) {
	query := `
	SELECT
	    id, first_name, last_name, number_phone
	FROM users
	LIMIT $1
	OFFSET $2
	`

	users := make([]models.User, 0)
	rows, err := s.conn.Query(ctx, query, pt.Limit, pt.Offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", op, err.Error())
	}

	for rows.Next() {
		var oneUser models.User
		if err := rows.Scan(&oneUser.ID, &oneUser.FirstName, &oneUser.LastName, &oneUser.NumberPhone); err != nil {
			return nil, fmt.Errorf("%s: %s", op, err.Error())
		}

		users = append(users, oneUser)
	}

	return users, nil
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
		return fmt.Errorf("%s: %s", op, err.Error())
	}

	if count := ct.RowsAffected(); count == 0 {
		return models.ErrNoRows
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
		return fmt.Errorf("%s: %s", op, err.Error())
	}

	if count := ct.RowsAffected(); count == 0 {
		return models.ErrNoRows
	}

	return nil
}
