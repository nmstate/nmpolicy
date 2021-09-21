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

package lexer_test

import (
	"testing"

	assert "github.com/stretchr/testify/require"

	"github.com/nmstate/nmpolicy/nmpolicy/internal/lexer"
)

func TestLexer(t *testing.T) {
	type expected struct {
		tokens []lexer.Token
		err    string
	}

	var tests = []struct {
		expression string
		expected   expected
	}{
		{"    ", expected{tokens: []lexer.Token{
			{3, lexer.EOF, ""}},
		}},
		{"    31    03   ", expected{tokens: []lexer.Token{
			{4, lexer.NUMBER, "31"},
			{10, lexer.NUMBER, "03"},
			{14, lexer.EOF, ""}},
		}},
		{` "foobar1" "foo 1 bar"    " foo bar - " ' bar foo' 789 "" `, expected{tokens: []lexer.Token{
			{1, lexer.STRING, "foobar1"},
			{11, lexer.STRING, "foo 1 bar"},
			{26, lexer.STRING, " foo bar - "},
			{40, lexer.STRING, " bar foo"},
			{51, lexer.NUMBER, "789"},
			{55, lexer.STRING, ""},
			{57, lexer.EOF, ""}},
		}},
		{` "bar1" "foo dar`, expected{
			err: `invalid string format (missing " terminator), pos=15`,
		}},
		{` "bar1" 'foo dar`, expected{
			err: "invalid string format (missing ' terminator), pos=15",
		}},
		{"155 -44", expected{
			err: "illegal rune -",
		}},
		{"255 1.3", expected{
			err: "invalid number format (. is not a digit)",
		}},
		{"355 1e3", expected{
			err: "invalid number format (e is not a digit)",
		}},
		{"455 0xEA", expected{
			err: "invalid number format (x is not a digit)",
		}},
		{"555 2.3-4", expected{
			err: "invalid number format (. is not a digit)",
		}},
		{"655 3333_444_333", expected{
			err: "invalid number format (_ is not a digit)",
		}},
		{"755 33 44 -.3", expected{
			err: "illegal rune -",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.expression, func(t *testing.T) {
			obtainedTokens, obtainedErr := lexer.New().Lex(tt.expression)
			if tt.expected.err != "" {
				assert.EqualError(t, obtainedErr, tt.expected.err)
			} else {
				assert.NoError(t, obtainedErr)
				assert.Equal(t, tt.expected.tokens, obtainedTokens)
			}
		})
	}
}
