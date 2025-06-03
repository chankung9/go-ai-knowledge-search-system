package validate

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// MustValid validates the given struct using go-playground/validator.
// It panics if validation fails.
func MustValid(s interface{}) {
	validate := validator.New()
	if err := validate.Struct(s); err != nil {
		panic(fmt.Errorf("config validation failed: %w", err))
	}
}
