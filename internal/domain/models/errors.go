package models

import "errors"

var (
	// server
	ErrInternal   error = errors.New("internal server error")
	ErrBadRequest error = errors.New("bad request")

	// validate
	ErrFieldRequired error = errors.New("required field is empty")
	ErrCorrectLength error = errors.New("length is not correct")
	ErrCorrectNumber error = errors.New("number phone is not correct")

	// handler
	ErrInvalidID          error = errors.New("invalid id")
	ErrInvalidLimitOffset error = errors.New("invalid limit or offset")

	// repository
	ErrNoRows      error = errors.New("not found")
	ErrTargetExist error = errors.New("target already exists")
)

func IsErrValidate(Err error) bool {
	return errors.As(Err, &ErrFieldRequired) || errors.As(Err, &ErrCorrectLength) || errors.As(Err, &ErrCorrectNumber)
}
