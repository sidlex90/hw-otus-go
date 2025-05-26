package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

var (
	//validator params errors
	ErrValidatorCompilationError = errors.New("validator compilation error: ")

	//validation error
	ErrInvalidStringLen    = errors.New("invalid string length: ")
	ErrInvalidStringRegexp = errors.New("invalid string format: ")
	ErrInvalidMin          = errors.New("invalid min value: ")
	ErrInvalidMax          = errors.New("invalid max value: ")
	ErrInvalidIn           = errors.New("value is not in the allowed set: ")
)

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}

func Validate(v interface{}) (ValidationErrors, error) {

	var validationErrors []ValidationError

	// Ensure the input is a struct
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Struct {
		typ := val.Type()
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)

			// Get the "validate" tag
			tag := field.Tag.Get("validate")
			if tag == "" {
				continue
			}

			fieldValue := val.Field(i)
			validations := strings.Split(tag, "|")
			for _, validation := range validations {
				if vErr, err := eachValidation(validation, field.Name, fieldValue); err == nil && vErr != nil {
					validationErrors = append(validationErrors, ValidationError{
						Field: field.Name,
						Err:   vErr,
					})
				} else if err != nil {
					return validationErrors, err
				}
			}
		}
	}

	return validationErrors, nil
}

func eachValidation(validation, fieldName string, fieldValue reflect.Value) (error, error) {
	if fieldValue.Kind() == reflect.Slice {
		for i := 0; i < fieldValue.Len(); i++ {
			if vErr, err := runValidation(validation, fieldName, fieldValue.Index(i)); err != nil || vErr != nil {
				return vErr, err
			}
		}
		return nil, nil
	}

	return runValidation(validation, fieldName, fieldValue)
}

func runValidation(validation, fieldName string, fieldValue reflect.Value) (error, error) {
	parts := strings.SplitN(validation, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("%w invalid validation format: %s", ErrValidatorCompilationError, validation)
	}

	rule, param := parts[0], parts[1]

	switch rule {
	case "len":
		return validateStringLen(param, fieldValue.String())
	case "min":
		return validateMin(fieldValue, param)
	case "max":
		return validateMax(fieldValue, param)
	case "regexp":
		return validateRegex(param, fieldValue.String())
	case "in":
		return validateIn(param, fieldValue)
	default:
		return nil, fmt.Errorf("unknown validation rule: %s", rule)
	}
}

func validateStringLen(param string, value string) (error, error) {
	expectedLen, err := strconv.Atoi(param)
	if err != nil {
		return nil, fmt.Errorf("%w invalid length value: %s", ErrValidatorCompilationError, param)
	}

	if len(value) != expectedLen {
		return fmt.Errorf("%w string length must be %d, got %d", ErrInvalidStringLen, expectedLen, len(value)), nil
	}
	return nil, nil
}

func validateRegex(regex string, value string) (error, error) {
	re, err := regexp.Compile(regex)
	if err != nil {
		return nil, fmt.Errorf("%w invalid regex: %s", ErrValidatorCompilationError, regex)
	}

	if !re.MatchString(value) {
		return fmt.Errorf("%w value does not match regex: %s", ErrInvalidStringRegexp, regex), nil
	}
	return nil, nil
}

func validateMin(fieldValue reflect.Value, param string) (error, error) {
	min, err := strconv.Atoi(param)
	if err != nil {
		return nil, fmt.Errorf("%w invalid min value: %s", ErrValidatorCompilationError, param)
	}

	if fieldValue.Kind() == reflect.Int && int(fieldValue.Int()) < min {
		return fmt.Errorf("%wvalue must be at least %d", ErrInvalidMin, min), nil
	}
	return nil, nil
}

func validateMax(fieldValue reflect.Value, param string) (error, error) {
	max, err := strconv.Atoi(param)
	if err != nil {
		return nil, fmt.Errorf("%w invalid max value: %s", ErrValidatorCompilationError, param)
	}

	if fieldValue.Kind() == reflect.Int && int(fieldValue.Int()) > max {
		return fmt.Errorf("%w value must be at most %d", ErrInvalidMax, max), nil
	}
	return nil, nil
}

func validateIn(param string, fieldValue reflect.Value) (error, error) {
	allowedValues := strings.Split(param, ",")
	switch fieldValue.Kind() {
	case reflect.String:
		value := fieldValue.String()
		for _, allowed := range allowedValues {
			if value == allowed {
				return nil, nil
			}
		}
		return fmt.Errorf("%w value must be one of %v", ErrInvalidIn, allowedValues), nil
	case reflect.Int:
		value := int(fieldValue.Int())
		for _, allowed := range allowedValues {
			if intVal, err := strconv.Atoi(allowed); err == nil && value == intVal {
				return nil, nil
			}
		}
		return fmt.Errorf("%w value must be one of %v", ErrInvalidIn, allowedValues), nil
	default:
		return nil, fmt.Errorf("%w unsupported type for 'in' validation", ErrValidatorCompilationError)
	}
}
