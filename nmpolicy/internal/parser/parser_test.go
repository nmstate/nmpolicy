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

package parser_test

import (
	"testing"

	assert "github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"

	"github.com/nmstate/nmpolicy/nmpolicy/internal/ast"
	"github.com/nmstate/nmpolicy/nmpolicy/internal/lexer"
	"github.com/nmstate/nmpolicy/nmpolicy/internal/parser"
)

func TestParser(t *testing.T) {
	testParseFailures(t)
	testParseSuccess(t)
}

func testParseFailures(t *testing.T) {
	var tests = []test{
		expectError("invalid expression: unexpected token `.`",
			fromTokens(
				dot(),
				eof(),
			),
		),
		expectError(`invalid path: missing identity after dot`,
			fromTokens(
				identity("routes"),
				dot(),
				eof(),
			),
		),
		expectError(`invalid path: missing dot`,
			fromTokens(
				identity("routes"),
				identity("destination"),
				eof(),
			),
		),
		expectError(`invalid path: missing identity after dot`,
			fromTokens(
				identity("routes"),
				dot(),
				dot(),
				identity("destination"),
				eof(),
			),
		),
	}
	runTest(t, tests)
}

func testParseSuccess(t *testing.T) {
	var tests = []test{
		expectAST(t, `
pos: 0
path:
- pos: 0
  identity: routes
- pos: 7
  identity: running
- pos: 15
  identity: destination`,
			fromTokens(
				identity("routes"),
				dot(),
				identity("running"),
				dot(),
				identity("destination"),
				eof(),
			),
		),
		expectError(`invalid equality filter: missing left hand argument`,
			fromTokens(
				eqfilter(),
				str("0.0.0.0/0"),
				eof(),
			),
		),
		expectError(`invalid equality filter: left hand argument is not a path`,
			fromTokens(
				str("foo"),
				eqfilter(),
				str("0.0.0.0/0"),
				eof(),
			),
		),
		expectError(`invalid equality filter: right hand argument is not a string`,
			fromTokens(
				identity("routes"),
				dot(),
				identity("running"),
				dot(),
				identity("destination"),
				eqfilter(),
				identity("this"),
				dot(),
				identity("is"),
				dot(),
				identity("wrong"),
				eof(),
			),
		),

		expectAST(t, `
pos: 26
eqfilter: 
- pos: 0
  identity: currentState
- pos: 0
  path: 
  - pos: 0
    identity: routes
  - pos: 7
    identity: running
  - pos: 15
    identity: destination
- pos: 28 
  string: 0.0.0.0/0`,
			fromTokens(
				identity("routes"),
				dot(),
				identity("running"),
				dot(),
				identity("destination"),
				eqfilter(),
				str("0.0.0.0/0"),
				eof(),
			),
		),
	}
	runTest(t, tests)
}

func runTest(t *testing.T, tests []test) {
	for _, tt := range tests {
		t.Run(description(tt), func(t *testing.T) {
			p := parser.New()
			obtainedAST, obtainedErr := p.Parse(tt.tokens)
			if tt.expected.err != "" {
				assert.EqualError(t, obtainedErr, tt.expected.err)
			} else {
				assert.NoError(t, obtainedErr)
				assert.Equal(t, *tt.expected.ast, obtainedAST)
			}
		})
	}
}

type expected struct {
	ast *ast.Node
	err string
}

type test struct {
	expression string
	tokens     []lexer.Token
	expected   expected
}

func description(tst test) string {
	if tst.expected.err != "" {
		return tst.expected.err
	}
	return tst.expression
}

func fromTokens(tokens ...lexer.Token) test {
	t := test{tokens: tokens}
	for i := range t.tokens {
		t.tokens[i].Position = len(t.expression)
		t.expression += t.tokens[i].Literal
	}
	return t
}

func expectAST(t *testing.T, astYAML string, tst test) test {
	a := &ast.Node{}
	err := yaml.Unmarshal([]byte(astYAML), a)
	assert.NoError(t, err)
	tst.expected.ast = a
	return tst
}

func expectError(err string, tst test) test {
	tst.expected.err = err
	return tst
}

func identity(literal string) lexer.Token {
	return lexer.Token{Type: lexer.IDENTITY, Literal: literal}
}

func str(literal string) lexer.Token {
	return lexer.Token{Type: lexer.STRING, Literal: literal}
}

func dot() lexer.Token {
	return lexer.Token{Type: lexer.DOT, Literal: "."}
}

func eof() lexer.Token {
	return lexer.Token{Type: lexer.EOF, Literal: ""}
}

func eqfilter() lexer.Token {
	return lexer.Token{Type: lexer.EQFILTER, Literal: "=="}
}
