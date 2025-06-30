package request_test

import (
	"testing"

	"github.com/1995parham/saf/internal/infra/http/request"
)

// nolint: funlen
func TestEventValidation(t *testing.T) {
	t.Parallel()

	cases := []struct {
		request request.Event
		isValid bool
	}{
		{
			request: request.Event{
				Subject: "",
				ID:      "",
				Service: "",
				Payload: []byte{},
			},
			isValid: false,
		},
		{
			request: request.Event{
				Subject: "hello",
				ID:      "",
				Service: "OfferService",
				Payload: []byte{},
			},
			isValid: true,
		},
		{
			request: request.Event{
				Subject: "hello",
				ID:      "",
				Service: "",
				Payload: []byte{},
			},
			isValid: false,
		},
		{
			request: request.Event{
				Subject: "hello",
				ID:      "",
				Service: "NewService",
				Payload: []byte{},
			},
			isValid: false,
		},
		{
			request: request.Event{
				Subject: "hello",
				ID:      "",
				Payload: []byte("Hello World"),
				Service: "RideLifeCycleService",
			},
			isValid: true,
		},
	}

	for _, c := range cases {
		rq := c.request

		err := rq.Validate()
		if c.isValid && err != nil {
			t.Fatalf("valid request %+v has error %s", rq, err)
		}

		if !c.isValid && err == nil {
			t.Fatalf("invalid request %+v has no error", rq)
		}
	}
}
