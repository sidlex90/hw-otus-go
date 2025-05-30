package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	Request struct {
		FakeURL string `validate:"len:invalid"`
	}

	ManyValidationStruct struct {
		Field string `validate:"len:20|regexp:^[0-9]*$"`
	}
)

func TestValidate(t *testing.T) {
	var (
		DefID     = "123e4567-e89b-12d3-a456-426614174000"
		DefName   = "John Doe"
		DefAge    = 30
		DefEmail  = "johndoe@example.com"
		DefRole   = UserRole("admin")
		DefPhones = []string{"12345678901", "09876543210"}
	)

	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{Response{}, ErrInvalidIn},
		{Response{200, ""}, nil},
		{Response{404, "{\"b\": 55}"}, nil},
		{Response{500, ""}, nil},
		{Response{405, ""}, ErrInvalidIn},
		{App{"valid"}, nil},
		{App{"in valid"}, ErrInvalidStringLen},
		{App{"min"}, ErrInvalidStringLen},
		{Token{}, nil},
		{
			User{
				ID:     DefID,
				Name:   DefName,
				Age:    3,
				Email:  DefEmail,
				Role:   DefRole,
				Phones: DefPhones,
			},
			ErrInvalidMin,
		},
		{
			User{ID: DefID, Name: DefName, Age: DefAge, Email: "lkjdflkdjf", Role: DefRole, Phones: DefPhones},
			ErrInvalidStringRegexp,
		},
		{
			User{ID: DefID, Name: DefName, Age: 100, Email: DefEmail, Role: DefRole, Phones: DefPhones},
			ErrInvalidMax,
		},
		{
			User{
				ID:     DefID,
				Name:   DefName,
				Age:    DefAge,
				Email:  DefEmail,
				Role:   DefRole,
				Phones: []string{"1234567", "09876543210"},
			},
			ErrInvalidStringLen,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			result, err := Validate(tt.in)
			assert.Nil(t, err)
			if tt.expectedErr != nil {
				assert.Len(t, result, 1)
				assert.ErrorIs(t, result[0].Err, tt.expectedErr)
			} else {
				assert.Len(t, result, 0)
			}
			_ = tt
		})
	}

	// many test cases
	user := User{ID: "invalid", Phones: []string{"09876543210", "1234567890"}}
	result, err := Validate(user)
	assert.Nil(t, err)
	assert.Len(t, result, 5)

	// invalid validator format
	_, err = Validate(Request{FakeURL: "https://test.ru"})
	assert.ErrorIs(t, err, ErrValidatorCompilationError)

	// test if input not struct
	tests2 := []struct {
		name        string
		input       interface{}
		expectedErr error
	}{
		{"Nil input", nil, ErrNotStruct},
		{"String input", "not a struct", ErrNotStruct},
		{"Int input", 123, ErrNotStruct},
		{"Slice input", []string{"value1", "value2"}, ErrNotStruct},
	}

	for _, tt := range tests2 {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Validate(tt.input)
			assert.ErrorIs(t, err, tt.expectedErr)
			assert.Len(t, result, 0)
		})
	}

	// multi validations test
	mv := ManyValidationStruct{Field: "8VI4567890"}
	mvResult, err := Validate(mv)
	assert.Nil(t, err)
	assert.Len(t, mvResult, 2)
	for _, v := range mvResult {
		assert.True(t, errors.Is(v.Err, ErrInvalidStringLen) || errors.Is(v.Err, ErrInvalidStringRegexp))
	}
}
