package param

import (
	"errors"
	"strings"
)

func ValidateParamValues(param string, values *[]string, validValues []string) error {

	if values != nil && len(*values) != 0 {

		validMap := make(map[string]bool)
		for _, v := range toLowerSlice(validValues) {
			validMap[v] = true
		}

		var invalidValues []string
		for _, val := range toLowerSlice(*values) {
			if !validMap[val] {
				invalidValues = append(invalidValues, val)
			}
		}

		if len(invalidValues) != 0 {
			joined := strings.Join(invalidValues, ", ")
			return errors.New("Invalid param '" + param + "' value/s '" + joined + "', valid values:\n" + strings.Join(validValues, "\n"))
		}
	}

	return nil
}

func toLowerSlice(slice []string) []string {
	for i, s := range slice {
		slice[i] = ToSimpleParamName(s)
	}
	return slice
}

func ValidateAll(params ...func() error) error {
	for _, paramFunc := range params {
		if err := paramFunc(); err != nil {
			return err
		}
	}
	return nil
}
