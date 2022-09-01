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
	"fmt"
	"testing"

	assert "github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"

	"github.com/nmstate/nmpolicy/nmpolicy/internal/ast"
	"github.com/nmstate/nmpolicy/nmpolicy/internal/lexer"
	"github.com/nmstate/nmpolicy/nmpolicy/internal/parser"
)

func TestParser(t *testing.T) {
	testParsePath(t)
	testParseEqFilter(t)
	testParseReplace(t)
	testParseReplaceWithPath(t)
	testParseCapturePipeReplace(t)

	testParseBasicFailures(t)
	testParsePathFailures(t)
	testParseEqFilterFailure(t)
	testParseReplaceFailure(t)

	testParserReuse(t)
}

func testParseBasicFailures(t *testing.T) {
	var tests = []test{
		expectError("invalid expression: unexpected token `.`"+`
| .
| ^`,
			fromTokens(
				dot(),
				eof(),
			),
		),
	}
	runTest(t, tests)
}

func testParsePathFailures(t *testing.T) {
	var tests = []test{
		expectError(`invalid path: missing identity or number after dot
| routes.
| ......^`,
			fromTokens(
				identity("routes"),
				dot(),
				eof(),
			),
		),
		expectError(`invalid path: missing dot
| routesdestination
| ......^`,
			fromTokens(
				identity("routes"),
				identity("destination"),
				eof(),
			),
		),
		expectError(`invalid path: missing identity or number after dot
| routes..destination
| .......^`,
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

func testParseEqFilterFailure(t *testing.T) {
	var tests = []test{
		expectError(`invalid equality filter: missing left hand argument
| ==0.0.0.0/0
| ^`,
			fromTokens(
				eqfilter(),
				str("0.0.0.0/0"),
				eof(),
			),
		),
		expectError(`invalid equality filter: left hand argument is not a path
| foo==0.0.0.0/0
| ...^`,
			fromTokens(
				str("foo"),
				eqfilter(),
				str("0.0.0.0/0"),
				eof(),
			),
		),
		expectError(`invalid equality filter: missing right hand argument
| routes.running.destination==
| ...........................^`,
			fromTokens(
				identity("routes"),
				dot(),
				identity("running"),
				dot(),
				identity("destination"),
				eqfilter(),
				eof(),
			),
		),

		expectError(`invalid equality filter: right hand argument is not a string or identity
| routes.running.destination====
| ............................^`,
			fromTokens(
				identity("routes"),
				dot(),
				identity("running"),
				dot(),
				identity("destination"),
				eqfilter(),
				eqfilter(),
				eof(),
			),
		),
		expectError(`invalid pipe: missing pipe in expression
| |routes.running.next-hop-interface:=br1
| ^`,
			fromTokens(
				pipe(),
				identity("routes"),
				dot(),
				identity("running"),
				dot(),
				identity("next-hop-interface"),
				replace(),
				str("br1"),
				eof(),
			),
		),
		expectError(`invalid pipe: missing pipe out expression
| capture.default-gw|
| ..................^`,
			fromTokens(
				identity("capture"),
				dot(),
				identity("default-gw"),
				pipe(),
				eof(),
			),
		),

		expectError(`invalid pipe: only paths can be piped in
| foo|routes.running.next-hop-interface:=br1
| ...^`,
			fromTokens(
				str("foo"),
				pipe(),
				identity("routes"),
				dot(),
				identity("running"),
				dot(),
				identity("next-hop-interface"),
				replace(),
				str("br1"),
				eof(),
			),
		),
	}
	runTest(t, tests)
}

func testParseReplaceFailure(t *testing.T) {
	var tests = []test{
		expectError(`invalid replace: missing left hand argument
| :=0.0.0.0/0
| ^`,
			fromTokens(
				replace(),
				str("0.0.0.0/0"),
				eof(),
			),
		),
		expectError(`invalid replace: left hand argument is not a path
| foo:=0.0.0.0/0
| ...^`,
			fromTokens(
				str("foo"),
				replace(),
				str("0.0.0.0/0"),
				eof(),
			),
		),
		expectError(`invalid replace: missing right hand argument
| routes.running.destination:=
| ...........................^`,
			fromTokens(
				identity("routes"),
				dot(),
				identity("running"),
				dot(),
				identity("destination"),
				replace(),
				eof(),
			),
		),

		expectError(`invalid replace: right hand argument is not a string or identity
| routes.running.destination:=:=
| ............................^`,
			fromTokens(
				identity("routes"),
				dot(),
				identity("running"),
				dot(),
				identity("destination"),
				replace(),
				replace(),
				eof(),
			),
		),
	}
	runTest(t, tests)
}

func testParsePath(t *testing.T) {
	var tests = []test{
		expectEmptyAST(fromTokens()),
		expectEmptyAST(fromTokens(eof())),
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
	}
	runTest(t, tests)
}

func testParseEqFilter(t *testing.T) {
	var tests = []test{
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
		expectAST(t, `
pos: 33
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
    identity: next-hop-interface
- pos: 35
  path:
  - pos: 35 
    identity: capture
  - pos: 43
    identity: default-gw
  - pos: 54
    identity: routes
  - pos: 61
    number: 0
  - pos: 63
    identity: next-hop-interface
`,
			fromTokens(
				identity("routes"),
				dot(),
				identity("running"),
				dot(),
				identity("next-hop-interface"),
				eqfilter(),
				identity("capture"),
				dot(),
				identity("default-gw"),
				dot(),
				identity("routes"),
				dot(),
				number(0),
				dot(),
				identity("next-hop-interface"),
				eof(),
			),
		),
	}
	runTest(t, tests)
}

func testParseReplace(t *testing.T) {
	var tests = []test{
		expectAST(t, `
pos: 33
replace:
- pos: 0
  identity: currentState
- pos: 0 
  path: 
  - pos: 0 
    identity: routes
  - pos: 7
    identity: running
  - pos: 15
    identity: next-hop-interface
- pos: 35
  string: br1
`,
			fromTokens(
				identity("routes"),
				dot(),
				identity("running"),
				dot(),
				identity("next-hop-interface"),
				replace(),
				str("br1"),
				eof(),
			),
		),
		expectAST(t, `
pos: 23
replace:
- pos: 0
  identity: currentState
- pos: 0
  path:
  - pos: 0
    identity: interfaces
  - pos: 11
    identity: lldp
  - pos: 16
    identity: enabled
- pos: 25
  boolean: true
`,
			fromTokens(
				identity("interfaces"),
				dot(),
				identity("lldp"),
				dot(),
				identity("enabled"),
				replace(),
				boolean(true),
				eof(),
			),
		),
	}
	runTest(t, tests)
}

func testParseReplaceWithPath(t *testing.T) {
	var tests = []test{
		expectAST(t, `
pos: 33
replace:
- pos: 0
  identity: currentState
- pos: 0 
  path: 
  - pos: 0 
    identity: routes
  - pos: 7
    identity: running
  - pos: 15
    identity: next-hop-interface
- pos: 35
  path:
  - pos: 35
    identity: capture
  - pos: 43
    identity: primary-nic
  - pos: 55
    identity: interfaces
  - pos: 66
    number: 0
  - pos: 68
    identity: name
`,
			fromTokens(
				identity("routes"),
				dot(),
				identity("running"),
				dot(),
				identity("next-hop-interface"),
				replace(),
				identity("capture"),
				dot(),
				identity("primary-nic"),
				dot(),
				identity("interfaces"),
				dot(),
				number(0),
				dot(),
				identity("name"),
				eof(),
			),
		),
	}
	runTest(t, tests)
}

func testParseCapturePipeReplace(t *testing.T) {
	var tests = []test{
		expectAST(t, `
pos: 52
replace:
- pos: 0
  path:
  - pos: 0
    identity: capture
  - pos: 8
    identity: default-gw
- pos: 19 
  path: 
  - pos: 19
    identity: routes
  - pos: 26
    identity: running
  - pos: 34
    identity: next-hop-interface
- pos: 54
  string: br1
`,
			fromTokens(
				identity("capture"),
				dot(),
				identity("default-gw"),
				pipe(),
				identity("routes"),
				dot(),
				identity("running"),
				dot(),
				identity("next-hop-interface"),
				replace(),
				str("br1"),
				eof(),
			),
		),
	}
	runTest(t, tests)
}

func testParserReuse(t *testing.T) {
	p := parser.New()
	testToRun1 := expectAST(t, `
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
	)
	testToRun2 := expectAST(t, `
pos: 14
eqfilter: 
- pos: 0
  identity: currentState
- pos: 0
  path: 
  - pos: 0
    identity: routes
  - pos: 7
    identity: running
- pos: 16
  string: foo`,
		fromTokens(
			identity("routes"),
			dot(),
			identity("running"),
			eqfilter(),
			str("foo"),
			eof(),
		),
	)
	runTestWithParser(t, testToRun1, p)
	runTestWithParser(t, testToRun2, p)
}

func runTest(t *testing.T, tests []test) {
	for _, tt := range tests {
		t.Run(description(tt), func(t *testing.T) {
			runTestWithParser(t, tt, parser.New())
		})
	}
}

func runTestWithParser(t *testing.T, testToRun test, p parser.Parser) {
	obtainedAST, obtainedErr := p.Parse(testToRun.expression, testToRun.tokens)
	if testToRun.expected.err != "" {
		assert.EqualError(t, obtainedErr, testToRun.expected.err)
	} else {
		assert.NoError(t, obtainedErr)
		assert.Equal(t, *testToRun.expected.ast, obtainedAST)
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

func expectEmptyAST(tst test) test {
	tst.expected.ast = &ast.Node{}
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

func number(literal int) lexer.Token {
	return lexer.Token{Type: lexer.NUMBER, Literal: fmt.Sprintf("%d", literal)}
}

func boolean(literal bool) lexer.Token {
	return lexer.Token{Type: lexer.BOOLEAN, Literal: fmt.Sprintf("%t", literal)}
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

func replace() lexer.Token {
	return lexer.Token{Type: lexer.REPLACE, Literal: ":="}
}

func pipe() lexer.Token {
	return lexer.Token{Type: lexer.PIPE, Literal: "|"}
}
