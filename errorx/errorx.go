package errorx

import (
	"errors"
	"fmt"
)

// apiCode is the interface for Code feature.
type apiCode interface {
	Error() string // It should be an error.
	Code() int
}

// apiStack is the interface for Stack feature.
type apiStack interface {
	Error() string // It should be an error.
	Stack() string
}

// apiCause is the interface for Cause feature.
type apiCause interface {
	Error() string // It should be an error.
	Cause() error
}

// apiCurrent is the interface for Current feature.
type apiCurrent interface {
	Error() string // It should be an error.
	Current() error
}

// apiNext is the interface for Next feature.
type apiNext interface {
	Error() string // It should be an error.
	Next() error
}

// New creates and returns an error which is formatted from given text.
func New(msgAndArgs ...interface{}) error {
	return &Error{
		stack: callers(),
		text:  messageFromMsgAndArgs(msgAndArgs...),
		code:  -1,
	}
}

// NewSkip creates and returns an error which is formatted from given text.
// The parameter <skip> specifies the stack callers skipped amount.
func NewSkip(skip int, msgAndArgs ...interface{}) error {
	return &Error{
		stack: callers(skip),
		text:  messageFromMsgAndArgs(msgAndArgs...),
		code:  -1,
	}
}

// Wrap wraps error with text.
// It returns nil if given err is nil.
func Wrap(err error, msgAndArgs ...interface{}) error {
	if err == nil {
		return nil
	}

	return &Error{
		error: err,
		stack: callers(),
		text:  messageFromMsgAndArgs(msgAndArgs...),
		code:  Code(err),
	}
}

// WrapSkip wraps error with text.
// It returns nil if given err is nil.
// The parameter <skip> specifies the stack callers skipped amount.
func WrapSkip(skip int, err error, msgAndArgs ...interface{}) error {
	if err == nil {
		return nil
	}
	return &Error{
		error: err,
		stack: callers(skip),
		text:  messageFromMsgAndArgs(msgAndArgs...),
		code:  Code(err),
	}
}

// NewCode creates and returns an error that has error code and given text.
func NewCode(code int, msgAndArgs ...interface{}) error {
	return &Error{
		stack: callers(),
		text:  messageFromMsgAndArgs(msgAndArgs...),
		code:  code,
	}
}

// NewCodeSkip creates and returns an error which has error code and is formatted from given text.
// The parameter <skip> specifies the stack callers skipped amount.
func NewCodeSkip(code, skip int, msgAndArgs ...interface{}) error {
	return &Error{
		stack: callers(skip),
		text:  messageFromMsgAndArgs(msgAndArgs...),
		code:  code,
	}
}

// WrapCode wraps error with code and text.
// It returns nil if given err is nil.
func WrapCode(err error, code int, msgAndArgs ...interface{}) error {
	if err == nil {
		return nil
	}

	return &Error{
		error: err,
		stack: callers(),
		text:  messageFromMsgAndArgs(msgAndArgs...),
		code:  code,
	}
}

// WrapCodeSkip wraps error with code and text.
// It returns nil if given err is nil.
// The parameter <skip> specifies the stack callers skipped amount.
func WrapCodeSkip(err error, code, skip int, msgAndArgs ...interface{}) error {
	if err == nil {
		return nil
	}

	return &Error{
		error: err,
		stack: callers(skip),
		text:  messageFromMsgAndArgs(msgAndArgs...),
		code:  code,
	}
}

// Code returns the error code of current error.
// It returns -1 if it has no error code or it does not implements interface Code.
func Code(err error) int {
	if err != nil {
		if e, ok := err.(apiCode); ok {
			return e.Code()
		}
	}

	return -1
}

// Cause returns the root cause error of <err>.
func Cause(err error) error {
	if err != nil {
		if e, ok := err.(apiCause); ok {
			return e.Cause()
		}
	}

	return err
}

// Stack returns the stack callers as string.
// It returns the error string directly if the <err> does not support stacks.
func Stack(err error) string {
	if err == nil {
		return ""
	}

	if e, ok := err.(apiStack); ok {
		return e.Stack()
	}

	return err.Error()
}

// Current creates and returns the current level error.
// It returns nil if current level error is nil.
func Current(err error) error {
	if err == nil {
		return nil
	}

	if e, ok := err.(apiCurrent); ok {
		return e.Current()
	}

	return err
}

// Next returns the next level error.
// It returns nil if current level error or the next level error is nil.
func Next(err error) error {
	if err == nil {
		return nil
	}

	if e, ok := err.(apiNext); ok {
		return e.Next()
	}

	return nil
}

// Is reports whether any error in err's chain matches target.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As finds the first error in err's chain that matches target, and if so, sets
// target to that error value and returns true.
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
func Unwrap(err error) error {
	return Next(err)
}

func messageFromMsgAndArgs(msgAndArgs ...interface{}) string {
	if len(msgAndArgs) == 0 || msgAndArgs == nil {
		return ""
	}

	if len(msgAndArgs) == 1 {
		msg := msgAndArgs[0]
		if msgAsStr, ok := msg.(string); ok {
			return msgAsStr
		}

		return fmt.Sprintf("%+v", msg)
	}

	if len(msgAndArgs) > 1 {
		return fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	}

	return ""
}
