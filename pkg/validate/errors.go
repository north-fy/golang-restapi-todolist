package validate

import "errors"

var (
	ErrFieldRequired error = errors.New("field with name is empty")
	ErrCorrectLength error = errors.New("length is not correct")
	ErrCorrectNumber error = errors.New("number phone is not correct")
)
