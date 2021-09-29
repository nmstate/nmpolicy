package source

import (
	"fmt"
)

// Error add the Position and Snippet to the normal message error, the Snippet
// is generated with `Decorate` function
type Error struct {
	Position int
	Message  string
	Snippet  string
}

// Decorate add the `Snippet` to `Error`, using the Position and error and
// source it will compose a string with the expression and a pointer.
func (e *Error) Decorate(src Source) *Error {
	e.Snippet = src.Snippet(e.Position)
	return e
}

// Error mark this struct as a golang error.
func (e *Error) Error() string {
	return e.format()
}

// NewError construct a nmpolicy error that wraps golagn error with a position
// on the expression where the error is ocurring
func NewError(err error, pos int) *Error {
	return &Error{
		Position: pos,
		Message:  fmt.Sprintf("%v", err),
	}
}

func (e *Error) format() string {
	return fmt.Sprintf(
		"%s, pos=%d\n%s",
		e.Message,
		e.Position,
		e.Snippet,
	)
}
