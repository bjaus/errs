package errs

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/rewardStyle/creator-account-service/pkg/caller"
)

var _ Error = new(impl)

type impl struct {
	code Code
	msg  string
	log  bool
	err  error
	op   caller.Caller
}

func (i *impl) Log() *impl {
	i.log = true
	return i
}

func (i *impl) ToLog() bool {
	return i.log
}

func (i *impl) Kind() Code {
	return i.code
}

func (i *impl) Message() string {
	return i.msg
}

func (i *impl) Operation() caller.Caller {
	return i.op
}

func (i *impl) Error() string {
	var b bytes.Buffer
	msg := i.msg
	code := i.code

	fmt.Fprintf(&b, "%s [%s] %s", i.op, i.code, i.msg)

	var x *impl
	e := errors.Unwrap(i)

	for {
		if e == nil {
			break
		}

		s := b.String()
		b.Reset()

		if !errors.As(e, &x) {
			fmt.Fprintf(&b, "%s%s%s", e.Error(), Sep, s)
		} else {
			if x.msg == msg {
				if x.code == code {
					fmt.Fprintf(&b, "%s%s%s", x.op, Sep, s)
				} else {
					fmt.Fprintf(&b, "%s [%s]%s%s", x.op, x.code, Sep, s)
				}
			} else {
				fmt.Fprintf(&b, "%s [%s] %s%s%s", x.op, x.code, x.msg, Sep, s)
			}
			msg = x.msg
			code = x.code
		}

		e = errors.Unwrap(e)
	}

	return b.String()
}

func (i *impl) Unwrap() error {
	return i.err
}
