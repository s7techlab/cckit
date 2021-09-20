package router

import (
	"errors"
	"fmt"
)

var ErrInvalidRequest = errors.New(`invalid request`)

type (
	Validator interface {
		Validate() error
	}
)

// ValidateRequest use Validator interface and create error, allow to use error.Is(ErrInvalidRequest)
func ValidateRequest(request Validator) error {
	if err := request.Validate(); err != nil {
		return fmt.Errorf(`%s: %w`, err.Error(), ErrInvalidRequest)
	}

	return nil
}
