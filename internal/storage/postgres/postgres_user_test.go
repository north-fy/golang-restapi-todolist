package postgres

import (
	"context"
	"testing"

	"github.com/north-fy/golang-restapi-todolist/internal/domain/models"
	"github.com/north-fy/golang-restapi-todolist/internal/storage/postgres/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPostgres_CreateUser(t *testing.T) {
	tests := []struct {
		name    string
		data    models.User
		require error
	}{
		{
			name: "success",
			data: models.User{
				FirstName:   "ivan",
				LastName:    "kovalev",
				NumberPhone: "+89128114119",
			},
			require: nil,
		},
		{
			name: "not enough length",
			data: models.User{
				FirstName:   "iv",
				LastName:    "kov",
				NumberPhone: "+89128114119",
			},
			require: nil,
		},
		{
			name:    "empty fields",
			data:    models.User{},
			require: nil,
		},
		{
			name: "empty once field",
			data: models.User{
				FirstName:   "ivan",
				LastName:    "",
				NumberPhone: "+89128114119",
			},
			require: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(mocks.StorageUser)
			mockStorage.On("CreateUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1, nil)

			id, err := mockStorage.CreateUser(context.Background(), tt.data.FirstName, tt.data.LastName, tt.data.NumberPhone)

			assert.Equal(t, tt.require, err)
			assert.Equal(t, 1, id)
			mockStorage.AssertExpectations(t)
		})
	}
}
