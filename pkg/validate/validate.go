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

func ValidateName(data string) error {
	if data == "" && isRequiredName {
		return errFieldRequired
	}

	if minLengthName >= len(data) || len(data) >= maxLengthName {
		return errCorrectLength
	}

	return nil
}

func ValidateNumber(data string) error {
	if data == "" && isRequiredNumber {
		return errFieldRequired
	}

	if minLengthNumber >= len(data) || len(data) >= maxLengthNumber {
		return errCorrectLength
	}

	if !strings.HasPrefix(data, "+") {
		return errCorrectNumber
	}

	return nil
}
