/*
 * Copyright 2021 NMPolicy Authors.
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

package parser

import "fmt"

type InvalidExpressionError struct {
	msg string
}

func (e *InvalidExpressionError) Error() string {
	return fmt.Sprintf("invalid expression: %s", e.msg)
}

type InvalidPathError struct {
	msg string
}

func (e *InvalidPathError) Error() string {
	return fmt.Sprintf("invalid path: %s", e.msg)
}

type InvalidEqualityFilter struct {
	msg string
}

func (e *InvalidEqualityFilter) Error() string {
	return fmt.Sprintf("invalid equality filter: %s", e.msg)
}