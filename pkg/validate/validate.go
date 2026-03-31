package validate

import "strings"

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
		return ErrFieldRequired
	}

	if minLengthName >= len(data) || len(data) >= maxLengthName {
		return ErrCorrectLength
	}

	return nil
}

func ValidateNumber(data string) error {
	if data == "" && isRequiredNumber {
		return ErrFieldRequired
	}

	if minLengthNumber >= len(data) || len(data) >= maxLengthNumber {
		return ErrCorrectLength
	}

	if !strings.HasPrefix(data, "+") {
		return ErrCorrectNumber
	}

	return nil
}

func OptValidate(data string, isRequired bool, minLength, maxLength int) error {
	if data == "" && isRequired {
		return ErrFieldRequired
	}

	if !isRequired && data == "" {
		return nil
	}

	if minLength >= len(data) || len(data) >= maxLength {
		return ErrCorrectLength
	}

	return nil
}
