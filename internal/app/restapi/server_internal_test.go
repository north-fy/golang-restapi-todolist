package restapi

import (
	"testing"
)

func TestRoute_ConfigureRouter(t *testing.T) {
	// TODO: Проверить хендлеры/init router

	tests := []struct {
		name     string
		data     interface{}
		required error
	}{
		{
			name:     "success create user",
			data:     []byte{1},
			required: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

		})
	}
}
