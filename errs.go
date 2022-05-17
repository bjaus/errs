// errs package is intended for better error communication/handling between system
// components and includes user-facing error messages. This package was inspired by:
// - https://middlemost.com/failure-is-your-domain/
// - https://github.com/upspin/upspin
package errs

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rewardStyle/creator-account-service/pkg/caller"
)

var Sep string = " :: "

type Error interface {
	error
	Kind() Code
	ToLog() bool
	Message() string
	Operation() caller.Caller
}

func New(code Code, args ...interface{}) *impl {
	return newError(code, args...)
}

func Trace(err error) *impl {
	if err == nil {
		return nil
	}
	prev, ok := err.(*impl)
	if ok {
		return newError(prev.code, err, prev.msg)
	} else {
		return newError(CodeInternal, err)
	}
}

func Is(err error) bool {
	var i Error
	return errors.As(err, &i)
}

func From(err error) Error {
	if i := Error(nil); errors.As(err, &i) {
		return i
	}
	return nil
}

func Kind(err error) Code {
	var i Error
	if errors.As(err, &i) {
		return i.Kind()
	}
	return CodeInternal
}

func newError(code Code, args ...interface{}) *impl {
	if !code.isValid() {
		code = CodeInternal
	}

	i := impl{
		code: code,
		op:   caller.Parse(2),
	}

	// Handle invalid calls with a panic
	if len(args) == 0 {
		panic(fmt.Sprintf("errs: no arguments provided from %s", i.op))
	}

	// Handle external error passed in alone but it's nil
	if len(args) == 1 && args[0] == nil {
		return nil
	}

	var msg string
	var index int

OUT:
	// Iterate through args and assign to impl accordingly
	for idx, arg := range args {
		switch v := arg.(type) {
		case nil:
			continue
		case string:
			msg = v
			index = idx + 1
			break OUT
		case Error:
			copy := &v
			i.err = *copy
		case error:
			i.err = v
		default:
			panic(fmt.Sprintf("errs: invalid argument from %s", i.op))
		}
	}

	// Handle fmt.Sprintf style messages with cleaner output given bad argument count
	if msg != "" {
		args := args[index:]

		if len(args) == 0 {
			i.msg = msg
		} else {
			count := strings.Count(msg, "%") - strings.Count(msg, "%%")

			if len(args) == count {
				i.msg = fmt.Sprintf(msg, args...)
			} else if count < len(args) {
				args = args[:count]
				i.msg = fmt.Sprintf(msg, args...)
			} else {
				var b strings.Builder

				mark := '%'
				num := count - len(args)

				var found int

				for _, r := range msg {
					if r == mark {
						if found == num {
							break
						}
						found++
					}
					b.WriteRune(r)
				}
				v := strings.TrimSpace(b.String())
				msg = fmt.Sprintf("%s...", v)
				i.msg = fmt.Sprintf(msg, args...)
			}
		}
	}

	prev, ok := i.err.(*impl)

	// Ensure that a message is present
	if i.msg == "" {
		if ok {
			i.msg = prev.msg
		} else if i.err != nil {
			i.msg = i.err.Error()
		}
	}

	// Carry through log request
	if ok {
		if !i.log && prev.log {
			i.log = prev.log
		}
	}

	return &i
}
