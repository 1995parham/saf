package request_test

import (
	"testing"

	"github.com/1995parham/saf/internal/http/request"
)

func TestEventValidation(t *testing.T) {
	t.Parallel()

	cases := []struct {
		request request.Event
		isValid bool
	}{
		{
			request: request.Event{
				Subject: "",
				Payload: []byte{},
			},
			isValid: false,
		},
		{
			request: request.Event{
				Subject: "hello",
				Payload: []byte{},
			},
			isValid: true,
		},
		{
			request: request.Event{
				Subject: "hello",
				Payload: []byte("Hello World"),
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
