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
	Subject string `json:"subject,omitempty"`
	ID      string `json:"id,omitempty"`
	Service string `json:"service,omitempty"`
	Payload []byte `json:"payload,omitempty"`
}

// Validate event request payload.
func (r Event) Validate() error {
	services := make([]any, len(service.TypeNames()))
	for i, t := range service.TypeNames() {
		services[i] = t
	}

	err := validation.ValidateStruct(&r,
		validation.Field(&r.Subject, is.Alphanumeric, validation.Required),
		validation.Field(&r.Service, is.Alphanumeric, validation.Required, validation.In(services...)),
	)
	if err != nil {
		return fmt.Errorf("event request validation failed: %w", err)
	}

	return nil
}
