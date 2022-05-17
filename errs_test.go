package errs_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/bjaus/errs"
)

func TestLog(t *testing.T) {
	err := fmt.Errorf("external error")

	t.Run("no logging requested", func(t *testing.T) {
		err := errs.Internal(err)
		if err.ToLog() {
			t.Errorf("got %t, want %t", err.ToLog(), false)
		}
	})

	t.Run("logging requested", func(t *testing.T) {
		err := errs.Internal(err).Log()
		if !err.ToLog() {
			t.Errorf("got %t, want %t", err.ToLog(), false)
		}
	})
}

func TestNew(t *testing.T) {
	err := fmt.Errorf("external error")
	err = errs.New(errs.Code(-1), err)
	e := errs.From(err)
	if e == nil {
		t.Fatal("should not be nil")
	}
	want := errs.CodeInternal
	if e.Kind() != want {
		t.Errorf("code: got %q, want %q", e.Kind(), want)
	}
}

func TestTrace(t *testing.T) {

	t.Run("stack trace messages", func(t *testing.T) {
		a := func() error {
			return fmt.Errorf("external error")
		}

		b := func() error {
			err := a()
			return errs.Trace(err)
		}

		c := func() error {
			err := b()
			return errs.Conflict(err, "new error message")
		}

		d := func() error {
			err := c()
			return errs.Trace(err)
		}

		err := d()

		want := 2
		got := strings.Count(err.Error(), "external error")
		if got != want {
			t.Errorf("trace count: got %d, want %d", got, want)
		}

		want = 1
		got = strings.Count(err.Error(), "new error message")
		if got != want {
			t.Errorf("trace count: got %d, want %d", got, want)
		}
	})

	t.Run("trace nil error", func(t *testing.T) {
		err := errs.Trace(nil)
		if err != nil {
			t.Errorf("tracing nil should produce a nil error")
		}
	})
}

func TestMessageWithTooManyArgs(t *testing.T) {
	want := fmt.Sprintf("message with %d argument", 1)
	err := errs.Internal("message with %d argument", 1, "oops")

	if !strings.Contains(err.Error(), want) {
		t.Errorf("error: got %q, should contain %q", err.Error(), want)
	}

	if err.Message() != want {
		t.Errorf("message: got %q, want %q", err.Message(), want)
	}
}

func TestMessageWithTooFewArgs(t *testing.T) {
	want := fmt.Sprintf("message with %d argument, plus...", 1)
	err := errs.Internal("message with %d argument, plus %s", 1)

	if !strings.Contains(err.Error(), want) {
		t.Errorf("error: got %q, should contain %q", err.Error(), want)
	}

	if err.Message() != want {
		t.Errorf("message: got %q, want %q", err.Message(), want)
	}
}

func TestMessageWithArgs(t *testing.T) {
	want := fmt.Sprintf("message with %d argument", 1)
	x := fmt.Errorf("external error")
	err := errs.Internal(x, "message with %d argument", 1)

	if !strings.Contains(err.Error(), want) {
		t.Errorf("error: got %q, should contain %q", err.Error(), want)
	}

	if err.Message() != want {
		t.Errorf("message: got %q, want %q", err.Message(), want)
	}
}

func TestExternalWithMessage(t *testing.T) {
	x := fmt.Errorf("external error")
	msg := "internal message"
	err := errs.Internal(x, msg)

	if !strings.Contains(err.Error(), x.Error()) {
		t.Errorf("error: got %q, should contain %q", err.Error(), x.Error())
	}
	if !strings.Contains(err.Error(), msg) {
		t.Errorf("error: got %q, should contain %q", err.Error(), msg)
	}

	if err.Message() != msg {
		t.Errorf("message: got %q, want %q", err.Message(), msg)
	}
}

func TestExternalOnly(t *testing.T) {
	x := fmt.Errorf("external error")
	err := errs.Internal(x)

	want := x.Error()
	parts := strings.Split(err.Error(), errs.Sep)
	if len(parts) != 2 {
		t.Fatal("parts should have found a seperator in message")
	}

	if !strings.Contains(parts[1], want) {
		t.Errorf("error: got %q, should contain %q", parts[1], want)
	}

	if err.Message() != want {
		t.Errorf("message: got %q, want %q", err.Message(), want)
	}
}

func TestNilWithMessage(t *testing.T) {
	want := "internal message"
	err := errs.Internal(nil, want)
	if err == nil {
		t.Error("error should not be nil")
	}
	if !strings.Contains(err.Error(), want) {
		t.Errorf("error: got %q, should contain %q", err.Error(), want)
	}
	if err.Message() != want {
		t.Errorf("message: got %q, want %q", err.Message(), want)
	}
}

func TestNilOnly(t *testing.T) {
	err := errs.New(errs.CodeInternal, nil)
	if err != nil {
		t.Error("error should be nil")
	}
}

func TestNoArguments(t *testing.T) {
	defer func() {
		if rec := recover(); rec == nil {
			t.Error("should have produced panic")
		}
	}()

	errs.Internal() // nolint
}

func TestInvalidArguments(t *testing.T) {
	defer func() {
		if rec := recover(); rec == nil {
			t.Error("should have produced panic")
		}
	}()

	errs.Internal(1, 2.0, true) // nolint
}

func TestIs(t *testing.T) {
	var err error

	err = fmt.Errorf("external error")
	if errs.Is(err) {
		t.Errorf("should not be errs.Error")
	}

	err = errs.Unauthorized(err, "internal error")
	if !errs.Is(err) {
		t.Errorf("should be errs.Error")
	}

	err = errs.Forbidden(err)
	if !errs.Is(err) {
		t.Errorf("should be errs.Error")
	}
}

func TestFrom(t *testing.T) {
	var err error

	err = fmt.Errorf("external error")
	if e := errs.From(err); e != nil {
		t.Error("non-errs.Error should be nil")
	}

	err = errs.Internal(err, "some message")
	if e := errs.From(err); e == nil {
		t.Errorf("errs.Error should not be nil")
	}
}

func TestOperation(t *testing.T) {
	err := errs.Internal("some message")
	op := err.Operation()

	want := op.String()
	if !strings.Contains(err.Error(), want) {
		t.Errorf("error: got %q, should contain %q", err.Error(), want)
	}

	// NOTE: this is the current test function's name
	// If you change the function name and not the line
	// below then you'll break the test.
	want = "TestOperation"
	if op.Function() != want {
		t.Errorf("func: got %q, want %q", op.Function(), want)
	}

	want = "errs_test"
	if op.Package() != want {
		t.Errorf("pkg: got %q, want %q", op.Package(), want)
	}

	if op.LineNumber() <= 0 {
		t.Errorf("line: should be greater than 0")
	}
}

func TestCode(t *testing.T) {
	err := fmt.Errorf("external error")

	cases := []struct {
		err  errs.Error
		code errs.Code
		name string
	}{
		{err: errs.Conflict(err), code: errs.CodeConflict, name: "conflict"},
		{err: errs.Expired(err), code: errs.CodeExpired, name: "expired"},
		{err: errs.Forbidden(err), code: errs.CodeForbidden, name: "forbidden"},
		{err: errs.Internal(err), code: errs.CodeInternal, name: "internal"},
		{err: errs.Invalid(err), code: errs.CodeInvalid, name: "invalid"},
		{err: errs.NotFound(err), code: errs.CodeNotFound, name: "not_found"},
		{err: errs.NotImplemented(err), code: errs.CodeNotImplemented, name: "not_implemented"},
		{err: errs.Temporary(err), code: errs.CodeTemporary, name: "temporary"},
		{err: errs.Timeout(err), code: errs.CodeTimeout, name: "timeout"},
		{err: errs.Unauthorized(err), code: errs.CodeUnauthorized, name: "unauthorized"},
		{err: errs.Unprocessable(err), code: errs.CodeUnprocessable, name: "unprocessable"},
	}

	for _, tc := range cases {
		code := errs.Kind(tc.err)
		t.Run(code.String(), func(t *testing.T) {
			got := tc.err.Kind()
			if got != tc.code {
				t.Errorf("got %q, want %q", got, tc.code)
			}
			if got.String() != tc.name {
				t.Errorf("got %q, want %q", got, tc.name)
			}

			got = errs.Kind(tc.err)
			if got != tc.code {
				t.Errorf("got %q, want %q", got, tc.code)
			}
		})
	}

	code := errs.Kind(err)
	want := errs.CodeInternal
	if code != want {
		t.Errorf("external errors should produce %q error codes", want)
	}
	name := "internal"
	if code.String() != name {
		t.Errorf("name: got %q, want %q", code.String(), name)
	}

	if errs.Code(0).String() != "unknown" {
		t.Errorf("code zero value: got %q, want %q", errs.Code(0), "uknown")
	}
}

func TestUnwrap(t *testing.T) {
	e1 := fmt.Errorf("external error")
	e2 := errs.Internal(e1, "some message")

	got := errors.Unwrap(e2)
	if got != e1 {
		t.Errorf("got %v, want %v", got, e1)
	}
}
