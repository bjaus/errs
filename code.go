package errs

type Code int

const (
	_ Code = iota
	CodeInternal
	CodeConflict
	CodeExpired
	CodeForbidden
	CodeInvalid
	CodeNotFound
	CodeNotImplemented
	CodeTemporary
	CodeTimeout
	CodeUnauthorized
	CodeUnprocessable
)

var codes map[Code]string = map[Code]string{
	CodeInternal:       "internal",
	CodeConflict:       "conflict",
	CodeExpired:        "expired",
	CodeForbidden:      "forbidden",
	CodeInvalid:        "invalid",
	CodeNotFound:       "not_found",
	CodeNotImplemented: "not_implemented",
	CodeTemporary:      "temporary",
	CodeTimeout:        "timeout",
	CodeUnauthorized:   "unauthorized",
	CodeUnprocessable:  "unprocessable",
}

func (k Code) String() string {
	tag, ok := codes[k]
	if ok {
		return tag
	}
	return "unknown"
}

func (k Code) isValid() bool {
	_, ok := codes[k]
	return ok
}

// Conflict...
func Conflict(args ...interface{}) *impl {
	return newError(CodeConflict, args...)
}

// Forbidden...
func Forbidden(args ...interface{}) *impl {
	return newError(CodeForbidden, args...)
}

// Internal...
func Internal(args ...interface{}) *impl {
	return newError(CodeInternal, args...)
}

// Invalid...
func Invalid(args ...interface{}) *impl {
	return newError(CodeInvalid, args...)
}

// NotFound...
func NotFound(args ...interface{}) *impl {
	return newError(CodeNotFound, args...)
}

// NotImplemented...
func NotImplemented(args ...interface{}) *impl {
	return newError(CodeNotImplemented, args...)
}

// Timeout...
func Timeout(args ...interface{}) *impl {
	return newError(CodeTimeout, args...)
}

// Temporary...
func Temporary(args ...interface{}) *impl {
	return newError(CodeTemporary, args...)
}

// Unauthorized...
func Unauthorized(args ...interface{}) *impl {
	return newError(CodeUnauthorized, args...)
}

// Unprocessable...
func Unprocessable(args ...interface{}) *impl {
	return newError(CodeUnprocessable, args...)
}

// Expired...
func Expired(args ...interface{}) *impl {
	return newError(CodeExpired, args...)
}
