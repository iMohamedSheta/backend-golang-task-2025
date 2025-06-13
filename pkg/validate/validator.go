package validate

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validatorInstance *validator.Validate
	once              sync.Once
)

// Validator initializes and returns a singleton instance of the validator
func Validator() *validator.Validate {
	once.Do(func() {
		validatorInstance = validator.New()
	})
	return validatorInstance
}
