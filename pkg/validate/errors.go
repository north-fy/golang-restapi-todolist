package validate

import "errors"

var (
	errFieldRequired error = errors.New("field with name is empty")
	errCorrectLength error = errors.New("length is not correct")
	errCorrectNumber error = errors.New("number phone is not correct")
)
