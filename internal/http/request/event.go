package request

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

const (
	PasswordMinLength = 6
	PasswordMaxLength = 0
)

// Event represents a event request payload.
// By providing an identification for event request
// you can remove duplicate events.
type Event struct {
	Subject string
	ID      string
	Payload []byte
}

// Validate event request payload.
func (r Event) Validate() error {
	if err := validation.ValidateStruct(&r,
		validation.Field(&r.Subject, is.Alphanumeric, validation.Required),
	); err != nil {
		return fmt.Errorf("event request validation failed: %w", err)
	}

	return nil
}
