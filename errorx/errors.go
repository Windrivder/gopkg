package errorx

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"runtime"
	"strings"
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
func New(text string) error {
	return &Error{
		stack: callers(),
		text:  text,
		code:  -1,
	}
}

// Newf returns an error that formats as the given format and args.
func Newf(format string, args ...interface{}) error {
	return &Error{
		stack: callers(),
		text:  fmt.Sprintf(format, args...),
		code:  -1,
	}
}

// NewSkip creates and returns an error which is formatted from given text.
// The parameter <skip> specifies the stack callers skipped amount.
func NewSkip(skip int, text string) error {
	return &Error{
		stack: callers(skip),
		text:  text,
		code:  -1,
	}
}

// NewSkipf returns an error that formats as the given format and args.
// The parameter <skip> specifies the stack callers skipped amount.
func NewSkipf(skip int, format string, args ...interface{}) error {
	return &Error{
		stack: callers(skip),
		text:  fmt.Sprintf(format, args...),
		code:  -1,
	}
}

// Wrap wraps error with text.
// It returns nil if given err is nil.
func Wrap(err error, text string) error {
	if err == nil {
		return nil
	}
	return &Error{
		error: err,
		stack: callers(),
		text:  text,
		code:  Code(err),
	}
}

// Wrapf returns an error annotating err with a stack trace
// at the point Wrapf is called, and the format specifier.
// It returns nil if given <err> is nil.
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &Error{
		error: err,
		stack: callers(),
		text:  fmt.Sprintf(format, args...),
		code:  Code(err),
	}
}

// WrapSkip wraps error with text.
// It returns nil if given err is nil.
// The parameter <skip> specifies the stack callers skipped amount.
func WrapSkip(skip int, err error, text string) error {
	if err == nil {
		return nil
	}
	return &Error{
		error: err,
		stack: callers(skip),
		text:  text,
		code:  Code(err),
	}
}

// WrapSkipf wraps error with text that is formatted with given format and args.
// It returns nil if given err is nil.
// The parameter <skip> specifies the stack callers skipped amount.
func WrapSkipf(skip int, err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &Error{
		error: err,
		stack: callers(skip),
		text:  fmt.Sprintf(format, args...),
		code:  Code(err),
	}
}

// NewCode creates and returns an error that has error code and given text.
func NewCode(code int, text string) error {
	return &Error{
		stack: callers(),
		text:  text,
		code:  code,
	}
}

// NewCodef returns an error that has error code and formats as the given format and args.
func NewCodef(code int, format string, args ...interface{}) error {
	return &Error{
		stack: callers(),
		text:  fmt.Sprintf(format, args...),
		code:  code,
	}
}

// NewCodeSkip creates and returns an error which has error code and is formatted from given text.
// The parameter <skip> specifies the stack callers skipped amount.
func NewCodeSkip(code, skip int, text string) error {
	return &Error{
		stack: callers(skip),
		text:  text,
		code:  code,
	}
}

// NewCodeSkipf returns an error that has error code and formats as the given format and args.
// The parameter <skip> specifies the stack callers skipped amount.
func NewCodeSkipf(code, skip int, format string, args ...interface{}) error {
	return &Error{
		stack: callers(skip),
		text:  fmt.Sprintf(format, args...),
		code:  code,
	}
}

// WrapCode wraps error with code and text.
// It returns nil if given err is nil.
func WrapCode(code int, err error, text string) error {
	if err == nil {
		return nil
	}
	return &Error{
		error: err,
		stack: callers(),
		text:  text,
		code:  code,
	}
}

// WrapCodef wraps error with code and format specifier.
// It returns nil if given <err> is nil.
func WrapCodef(code int, err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &Error{
		error: err,
		stack: callers(),
		text:  fmt.Sprintf(format, args...),
		code:  code,
	}
}

// WrapCodeSkip wraps error with code and text.
// It returns nil if given err is nil.
// The parameter <skip> specifies the stack callers skipped amount.
func WrapCodeSkip(code, skip int, err error, text string) error {
	if err == nil {
		return nil
	}
	return &Error{
		error: err,
		stack: callers(skip),
		text:  text,
		code:  code,
	}
}

// WrapCodeSkipf wraps error with code and text that is formatted with given format and args.
// It returns nil if given err is nil.
// The parameter <skip> specifies the stack callers skipped amount.
func WrapCodeSkipf(code, skip int, err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &Error{
		error: err,
		stack: callers(skip),
		text:  fmt.Sprintf(format, args...),
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

// Error is custom error for additional features.
type Error struct {
	error error  // Wrapped error.
	stack stack  // Stack array, which records the stack information when this error is created or wrapped.
	text  string // Error text, which is created by New* functions.
	code  int    // Error code if necessary.
}

const (
	// Filtering key for current error module paths.
	stackFilterKeyLocal = "/errors/errors"
)

var (
	// goRootForFilter is used for stack filtering purpose.
	// Mainly for development environment.
	goRootForFilter = runtime.GOROOT()
)

func init() {
	if goRootForFilter != "" {
		goRootForFilter = strings.Replace(goRootForFilter, "\\", "/", -1)
	}
}

// Error implements the interface of Error, it returns all the error as string.
func (err *Error) Error() string {
	if err == nil {
		return ""
	}
	errStr := err.text
	if err.error != nil {
		if err.text != "" {
			errStr += ": "
		}
		errStr += err.error.Error()
	}
	return errStr
}

// Code returns the error code.
// It returns -1 if it has no error code.
func (err *Error) Code() int {
	if err == nil {
		return -1
	}
	return err.code
}

// Cause returns the root cause error.
func (err *Error) Cause() error {
	if err == nil {
		return nil
	}
	loop := err
	for loop != nil {
		if loop.error != nil {
			if e, ok := loop.error.(*Error); ok {
				// Internal Error struct.
				loop = e
			} else if e, ok := loop.error.(apiCause); ok {
				// Other Error that implements ApiCause interface.
				return e.Cause()
			} else {
				return loop.error
			}
		} else {
			// return loop
			// To be compatible with Case of https://github.com/pkg/errors.
			return errors.New(loop.text)
		}
	}
	return nil
}

// Format formats the frame according to the fmt.Formatter interface.
//
// %v, %s   : Print all the error string;
// %-v, %-s : Print current level error string;
// %+s      : Print full stack error list;
// %+v      : Print the error string and full stack error list;
func (err *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 's', 'v':
		switch {
		case s.Flag('-'):
			if err.text != "" {
				io.WriteString(s, err.text)
			} else {
				io.WriteString(s, err.Error())
			}
		case s.Flag('+'):
			if verb == 's' {
				io.WriteString(s, err.Stack())
			} else {
				io.WriteString(s, err.Error()+"\n"+err.Stack())
			}
		default:
			io.WriteString(s, err.Error())
		}
	}
}

// Stack returns the stack callers as string.
// It returns an empty string if the <err> does not support stacks.
func (err *Error) Stack() string {
	if err == nil {
		return ""
	}
	var (
		loop   = err
		index  = 1
		buffer = bytes.NewBuffer(nil)
	)
	for loop != nil {
		buffer.WriteString(fmt.Sprintf("%d. %-v\n", index, loop))
		index++
		formatSubStack(loop.stack, buffer)
		if loop.error != nil {
			if e, ok := loop.error.(*Error); ok {
				loop = e
			} else {
				buffer.WriteString(fmt.Sprintf("%d. %s\n", index, loop.error.Error()))
				index++
				break
			}
		} else {
			break
		}
	}
	return buffer.String()
}

// Current creates and returns the current level error.
// It returns nil if current level error is nil.
func (err *Error) Current() error {
	if err == nil {
		return nil
	}
	return &Error{
		error: nil,
		stack: err.stack,
		text:  err.text,
	}
}

// Next returns the next level error.
// It returns nil if current level error or the next level error is nil.
func (err *Error) Next() error {
	if err == nil {
		return nil
	}
	return err.error
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
// Note that do not use pointer as its receiver here.
func (err *Error) MarshalJSON() ([]byte, error) {
	return []byte(`"` + err.Error() + `"`), nil
}

// formatSubStack formats the stack for error.
func formatSubStack(st stack, buffer *bytes.Buffer) {
	index := 1
	space := "  "
	for _, p := range st {
		if fn := runtime.FuncForPC(p - 1); fn != nil {
			file, line := fn.FileLine(p - 1)
			// Custom filtering.
			if strings.Contains(file, stackFilterKeyLocal) {
				continue
			}
			// Avoid stack string like "<autogenerated>"
			if strings.Contains(file, "<") {
				continue
			}
			// Ignore GO ROOT paths.
			if goRootForFilter != "" &&
				len(file) >= len(goRootForFilter) &&
				file[0:len(goRootForFilter)] == goRootForFilter {
				continue
			}
			// Graceful indent.
			if index > 9 {
				space = " "
			}
			buffer.WriteString(fmt.Sprintf(
				"   %d).%s%s\n    \t%s:%d\n",
				index, space, fn.Name(), file, line,
			))
			index++
		}
	}
}
