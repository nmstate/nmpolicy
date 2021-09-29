/*
 * Copyright 2001 NMPolicy Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at:
 *
 *	  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package expression

import (
	"fmt"
)

// Error add the Position and Snippet to the normal message error, the Snippet
// is generated with `Decorate` function
type Error struct {
	position int
	inner    error
	snippet  string
}

// Decorate add the `Snippet` to `Error`, using the Position and error and
// source it will compose a string with the expression and a pointer.
func (e *Error) Decorate(expression string) *Error {
	e.snippet = snippet(expression, e.position)
	return e
}

// Error mark this struct as a golang error.
func (e *Error) Error() string {
	return e.format()
}

func (e *Error) Unwrap() error {
	return e.inner
}

// NewError construct a nmpolicy error that wraps golang error with a position
// on the expression where the error is ocurring
func WrapError(err error, pos int) *Error {
	return &Error{
		position: pos,
		inner:    err,
	}
}

func (e *Error) format() string {
	errMsg := fmt.Sprintf(
		"%v, pos=%d",
		e.inner,
		e.position,
	)
	if e.snippet != "" {
		errMsg = fmt.Sprintf(
			"%s\n%s",
			errMsg,
			e.snippet,
		)
	}
	return errMsg
}
