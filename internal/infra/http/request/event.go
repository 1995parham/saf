package request

import (
	"fmt"

	"github.com/1995parham/saf/internal/domain/model/service"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// Event represents a event request payload.
// By providing an identification for event request
// you can remove duplicate events.
type Event struct {
	Subject string
	ID      string
	Service string
	Payload []byte
}

// Validate event request payload.
func (r Event) Validate() error {
	services := make([]any, 0)
	for _, t := range service.TypeNames() {
		services = append(services, t)
	}

	if err := validation.ValidateStruct(&r,
		validation.Field(&r.Subject, is.Alphanumeric, validation.Required),
		validation.Field(&r.Service, is.Alphanumeric, validation.Required, validation.In(services...)),
	); err != nil {
		return fmt.Errorf("event request validation failed: %w", err)
	}

	return nil
}
