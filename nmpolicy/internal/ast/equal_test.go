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

package ast_test

import (
	"testing"

	assert "github.com/stretchr/testify/require"

	"github.com/nmstate/nmpolicy/nmpolicy/internal/ast"
)

func TestTeminalDeepEqual(t *testing.T) {
	literalA := "literalA"
	literalB := "literalB"

	newPtr := func(literal string) *string {
		return &literal
	}

	tests := []struct {
		lhs      ast.Terminal
		rhs      ast.Terminal
		expected bool
	}{
		{ast.Terminal{Identity: &literalA}, ast.Terminal{Identity: &literalA}, true},
		{ast.Terminal{Str: &literalA}, ast.Terminal{Str: &literalA}, true},
		{
			ast.Terminal{Str: &literalA, Identity: &literalB},
			ast.Terminal{Str: &literalA, Identity: &literalB},
			true,
		},
		{
			ast.Terminal{Str: nil},
			ast.Terminal{Str: nil},
			true,
		},
		{
			ast.Terminal{Str: nil},
			ast.Terminal{Str: &literalA},
			false,
		},
		{
			ast.Terminal{Str: &literalA},
			ast.Terminal{Str: &literalB},
			false,
		},
		{
			ast.Terminal{Str: &literalA},
			ast.Terminal{Str: newPtr(literalA)},
			true,
		},
		{
			ast.Terminal{Str: &literalA},
			ast.Terminal{Str: newPtr(literalB)},
			false,
		},
		{
			ast.Terminal{Identity: nil},
			ast.Terminal{Identity: &literalA},
			false,
		},
		{
			ast.Terminal{Identity: &literalA},
			ast.Terminal{Identity: &literalB},
			false,
		},
		{
			ast.Terminal{Identity: &literalA},
			ast.Terminal{Identity: newPtr(literalA)},
			true,
		},
		{
			ast.Terminal{Identity: &literalA},
			ast.Terminal{Identity: newPtr(literalB)},
			false,
		},
		{
			ast.Terminal{Identity: &literalA, Str: &literalA},
			ast.Terminal{Identity: &literalB, Str: newPtr(literalA)},
			false,
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, test.lhs.DeepEqual(test.rhs))
	}
}
