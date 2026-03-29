package postgres

import (
	"testing"

	"github.com/north-fy/golang-restapi-todolist/internal/domain/models"
	"github.com/north-fy/golang-restapi-todolist/internal/storage/postgres/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPostgres_CreateUser(t *testing.T) {
	tests := []struct {
		name    string
		data    interface{}
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
			require: nil, // valid error ?
		},
		{
			name:    "empty fields",
			data:    models.User{},
			require: nil, // empty err
		},
		{
			name: "empty once field",
			data: models.User{
				FirstName:   "",
				LastName:    "kovalev",
				NumberPhone: "+89128114119",
			},
			require: nil, // ??
		},
	}

	storage := mocks.NewStorageUser(t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			user, ok := test.data.(models.User)
			require.Equal(t, ok, true)

			storage.On("CreateUser", t.Context(), user.FirstName,
				user.LastName, user.NumberPhone).Return(mock.Anything, test.require).Once()

			id, err := storage.CreateUser(t.Context(), user.FirstName,
				user.LastName, user.NumberPhone)
			assert.EqualError(t, err, test.require.Error())
			_ = id
		})
	}
}
