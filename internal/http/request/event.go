package request

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	PasswordMinLength = 6
	PasswordMaxLength = 0
)

// Event represents a event request payload.
type Register struct {
	Subject string
	Payload []byte
}

// Validate register request payload.
func (r Register) Validate() error {
	if err := validation.ValidateStruct(&r,
		validation.Field(&r.Subject, validation.Required),
		validation.Field(&r.Payload, validation.Required),
	); err != nil {
		return fmt.Errorf("event request validation failed: %w", err)
	}

	return nil
}
