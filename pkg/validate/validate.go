package validate

import (
	"strings"

	"github.com/north-fy/golang-restapi-todolist/internal/domain/models"
)

const (
	// Filter by name
	isRequiredName = true
	minLengthName  = 3
	maxLengthName  = 100

	// Filter by number phone
	isRequiredNumber = false
	minLengthNumber  = 10
	maxLengthNumber  = 15
)

// TODO: зарефакторить код в одну функцию
func ValidateName(data string) error {
	if data == "" && isRequiredName {
		return models.ErrFieldRequired
	}

	if minLengthName >= len(data) || len(data) >= maxLengthName {
		return models.ErrCorrectLength
	}

	return nil
}

func ValidateNumber(data string) error {
	if data == "" && isRequiredNumber {
		return models.ErrFieldRequired
	}

	if minLengthNumber >= len(data) || len(data) >= maxLengthNumber {
		return models.ErrCorrectLength
	}

	if !strings.HasPrefix(data, "+") {
		return models.ErrCorrectNumber
	}

	return nil
}

func OptValidate(data string, isRequired bool, minLength, maxLength int) error {
	if data == "" && isRequired {
		return models.ErrFieldRequired
	}

	if !isRequired && data == "" {
		return nil
	}

	if minLength >= len(data) || len(data) >= maxLength {
		return models.ErrCorrectLength
	}

	return nil
}
