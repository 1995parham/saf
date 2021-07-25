package request_test

import (
	"testing"

	"github.com/1995parham/fandogh/internal/http/request"
)

// nolint: funlen
func TestRegisterValidation(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		password string
		email    string
		isValid  bool
	}{
		{
			isValid: false,
		},
		{
			name:     "Parham Alvani",
			password: "1234567",
			email:    "parham.alvani@gmail.com",
			isValid:  true,
		},
		{
			name:     "Parham Alvani",
			password: "1234567",
			email:    "parham.alvani_gmail",
			isValid:  false,
		},
		{
			name:     "Parham Alvani",
			password: "1234567",
			email:    "parham.alvani@gmail",
			isValid:  false,
		},
		{
			name:     "Parham Alvani",
			password: "1234567",
			email:    "parham.alvani@aut.ac.ir",
			isValid:  true,
		},
		{
			name:     "Parham Alvani",
			password: "12345",
			email:    "parham.alvani@aut.ac.ir",
			isValid:  false,
		},
		{
			name:     "Parham Alvani",
			password: "123456",
			email:    "parham.alvani@aut.ac.ir",
			isValid:  true,
		},
		{
			name:     "پرهام الوانی",
			password: "123456",
			email:    "parham.alvani@aut.ac.ir",
			isValid:  true,
		},
	}

	for _, c := range cases {
		rq := request.Register{
			Name:     c.name,
			Email:    c.email,
			Password: c.password,
		}

		err := rq.Validate()

		if c.isValid && err != nil {
			t.Fatalf("valid request %+v has error %s", rq, err)
		}

		if !c.isValid && err == nil {
			t.Fatalf("invalid request %+v has no error", rq)
		}
	}
}

func TestLoginValidation(t *testing.T) {
	t.Parallel()

	cases := []struct {
		password string
		email    string
		isValid  bool
	}{
		{
			isValid: false,
		},
		{
			password: "1234567",
			email:    "parham.alvani@gmail.com",
			isValid:  true,
		},
		{
			password: "1234567",
			email:    "parham.alvani_gmail",
			isValid:  false,
		},
		{
			password: "1234567",
			email:    "parham.alvani@gmail",
			isValid:  false,
		},
		{
			password: "1234567",
			email:    "parham.alvani@aut.ac.ir",
			isValid:  true,
		},
		{
			password: "12345",
			email:    "parham.alvani@aut.ac.ir",
			isValid:  false,
		},
		{
			password: "123456",
			email:    "parham.alvani@aut.ac.ir",
			isValid:  true,
		},
	}

	for _, c := range cases {
		rq := request.Login{
			Email:    c.email,
			Password: c.password,
		}

		err := rq.Validate()

		if c.isValid && err != nil {
			t.Fatalf("valid request %+v has error %s", rq, err)
		}

		if !c.isValid && err == nil {
			t.Fatalf("invalid request %+v has no error", rq)
		}
	}
}
