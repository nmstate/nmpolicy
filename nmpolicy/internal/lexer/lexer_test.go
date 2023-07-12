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

type expected struct {
	tokens []lexer.Token
	err    string
}

type test struct {
	expression string
	expected   expected
}

func TestLexer(t *testing.T) {
	testBasicExpressions(t)
	testFailures(t)
	testLinuxBridgeAtDefaultGwScenario(t)
}

func testBasicExpressions(t *testing.T) {
	t.Run("basic expressions", func(t *testing.T) {
		runTest(t, []test{
			{"    ", expected{tokens: []lexer.Token{
				{3, lexer.EOF, ""}},
			}},
			{"    31    03   ", expected{tokens: []lexer.Token{
				{4, lexer.NUMBER, "31"},
				{10, lexer.NUMBER, "03"},
				{14, lexer.EOF, ""}},
			}},
			{` "foobar1" "foo 1 bar"    " foo bar - " ' bar foo' 789 "" true false "true" "false" truse truefoo falsefoo`,
				expected{tokens: []lexer.Token{
					{1, lexer.STRING, "foobar1"},
					{11, lexer.STRING, "foo 1 bar"},
					{26, lexer.STRING, " foo bar - "},
					{40, lexer.STRING, " bar foo"},
					{51, lexer.NUMBER, "789"},
					{55, lexer.STRING, ""},
					{58, lexer.BOOLEAN, "true"},
					{63, lexer.BOOLEAN, "false"},
					{69, lexer.STRING, "true"},
					{76, lexer.STRING, "false"},
					{84, lexer.IDENTITY, "truse"},
					{90, lexer.IDENTITY, "truefoo"},
					{98, lexer.IDENTITY, "falsefoo"},
					{105, lexer.EOF, ""}},
				}},
			{" foo f1-o-o fo-o-o1  ", expected{tokens: []lexer.Token{
				{1, lexer.IDENTITY, "foo"},
				{5, lexer.IDENTITY, "f1-o-o"},
				{12, lexer.IDENTITY, "fo-o-o1"},
				{20, lexer.EOF, ""}},
			}},
			{" . foo1.dar1.0.dar2:=foo3 . dar3 ... moo3+boo3|doo3", expected{tokens: []lexer.Token{
				{1, lexer.DOT, "."},
				{3, lexer.IDENTITY, "foo1"},
				{7, lexer.DOT, "."},
				{8, lexer.IDENTITY, "dar1"},
				{12, lexer.DOT, "."},
				{13, lexer.NUMBER, "0"},
				{14, lexer.DOT, "."},
				{15, lexer.IDENTITY, "dar2"},
				{19, lexer.REPLACE, ":="},
				{21, lexer.IDENTITY, "foo3"},
				{26, lexer.DOT, "."},
				{28, lexer.IDENTITY, "dar3"},
				{33, lexer.DOT, "."},
				{34, lexer.DOT, "."},
				{35, lexer.DOT, "."},
				{37, lexer.IDENTITY, "moo3"},
				{41, lexer.MERGE, "+"},
				{42, lexer.IDENTITY, "boo3"},
				{46, lexer.PIPE, "|"},
				{47, lexer.IDENTITY, "doo3"},
				{50, lexer.EOF, ""}},
			}},
			{" . foo1.dar1:=foo2 . dar2 ... moo3+boo3|doo3 == := := !=", expected{tokens: []lexer.Token{
				{1, lexer.DOT, "."},
				{3, lexer.IDENTITY, "foo1"},
				{7, lexer.DOT, "."},
				{8, lexer.IDENTITY, "dar1"},
				{12, lexer.REPLACE, ":="},
				{14, lexer.IDENTITY, "foo2"},
				{19, lexer.DOT, "."},
				{21, lexer.IDENTITY, "dar2"},
				{26, lexer.DOT, "."},
				{27, lexer.DOT, "."},
				{28, lexer.DOT, "."},
				{30, lexer.IDENTITY, "moo3"},
				{34, lexer.MERGE, "+"},
				{35, lexer.IDENTITY, "boo3"},
				{39, lexer.PIPE, "|"},
				{40, lexer.IDENTITY, "doo3"},
				{45, lexer.EQFILTER, "=="},
				{48, lexer.REPLACE, ":="},
				{51, lexer.REPLACE, ":="},
				{54, lexer.NEFILTER, "!="},
				{55, lexer.EOF, ""}},
			}},
			{"foo1.3|foo2", expected{tokens: []lexer.Token{
				{0, lexer.IDENTITY, "foo1"},
				{4, lexer.DOT, "."},
				{5, lexer.NUMBER, "3"},
				{6, lexer.PIPE, "|"},
				{7, lexer.IDENTITY, "foo2"},
				{10, lexer.EOF, ""}},
			}},
		})
	})
}

func testFailures(t *testing.T) {
	t.Run("failures", func(t *testing.T) {
		runTest(t, []test{
			{"foo=bar", expected{
				err: `invalid EQFILTER operation format (b is not equal char)
| foo=bar
| ....^`,
			}},
			{" foo 1foo ", expected{
				err: `invalid number format (f is not a digit)
|  foo 1foo 
| ......^`,
			}},
			{" foo -foo ", expected{
				err: `illegal rune -
|  foo -foo 
| .....^`,
			}},
			{` "bar1" "foo dar`, expected{
				err: `invalid string format (missing " terminator)
|  "bar1" "foo dar
| ...............^`,
			}},
			{` "bar1" 'foo dar`, expected{
				err: `invalid string format (missing ' terminator)
|  "bar1" 'foo dar
| ...............^`,
			}},
			{"155 -44", expected{
				err: `illegal rune -
| 155 -44
| ....^`,
			}},
			{"255 1,3", expected{
				err: `invalid number format (, is not a digit)
| 255 1,3
| .....^`,
			}},
			{"355 1e3", expected{
				err: `invalid number format (e is not a digit)
| 355 1e3
| .....^`,
			}},
			{"455 0xEA", expected{
				err: `invalid number format (x is not a digit)
| 455 0xEA
| .....^`,
			}},
			{"555 2,3-4", expected{
				err: `invalid number format (, is not a digit)
| 555 2,3-4
| .....^`,
			}},
			{"655 3333_444_333", expected{
				err: `invalid number format (_ is not a digit)
| 655 3333_444_333
| ........^`,
			}},
			{"755 33 44 -.3", expected{
				err: `illegal rune -
| 755 33 44 -.3
| ..........^`,
			}},
		})
	})
}

func testLinuxBridgeAtDefaultGwScenario(t *testing.T) {
	t.Run("linux bridge at the default gateway scenario", func(t *testing.T) {
		runTest(t, []test{
			{`routes.running.destination=="0.0.0.0/0"`, expected{tokens: []lexer.Token{
				{0, lexer.IDENTITY, "routes"},
				{6, lexer.DOT, "."},
				{7, lexer.IDENTITY, "running"},
				{14, lexer.DOT, "."},
				{15, lexer.IDENTITY, "destination"},
				{26, lexer.EQFILTER, "=="},
				{28, lexer.STRING, "0.0.0.0/0"},
				{38, lexer.EOF, ""}},
			}},
			{`routes.running.next-hop-interface==capturer.default-gw.routes.running.0.next-hop-interface`, expected{tokens: []lexer.Token{
				{0, lexer.IDENTITY, "routes"},
				{6, lexer.DOT, "."},
				{7, lexer.IDENTITY, "running"},
				{14, lexer.DOT, "."},
				{15, lexer.IDENTITY, "next-hop-interface"},
				{33, lexer.EQFILTER, "=="},
				{35, lexer.IDENTITY, "capturer"},
				{43, lexer.DOT, "."},
				{44, lexer.IDENTITY, "default-gw"},
				{54, lexer.DOT, "."},
				{55, lexer.IDENTITY, "routes"},
				{61, lexer.DOT, "."},
				{62, lexer.IDENTITY, "running"},
				{69, lexer.DOT, "."},
				{70, lexer.NUMBER, "0"},
				{71, lexer.DOT, "."},
				{72, lexer.IDENTITY, "next-hop-interface"},
				{89, lexer.EOF, ""}},
			}},
			{`interfaces.name==capturer.default-gw.routes.running.0.next-hop-interface`, expected{tokens: []lexer.Token{
				{0, lexer.IDENTITY, "interfaces"},
				{10, lexer.DOT, "."},
				{11, lexer.IDENTITY, "name"},
				{15, lexer.EQFILTER, "=="},
				{17, lexer.IDENTITY, "capturer"},
				{25, lexer.DOT, "."},
				{26, lexer.IDENTITY, "default-gw"},
				{36, lexer.DOT, "."},
				{37, lexer.IDENTITY, "routes"},
				{43, lexer.DOT, "."},
				{44, lexer.IDENTITY, "running"},
				{51, lexer.DOT, "."},
				{52, lexer.NUMBER, "0"},
				{53, lexer.DOT, "."},
				{54, lexer.IDENTITY, "next-hop-interface"},
				{71, lexer.EOF, ""}},
			}},
			{`capturer.base-iface-routes | routes.running.next-hop-interface:="br1"`, expected{tokens: []lexer.Token{
				{0, lexer.IDENTITY, "capturer"},
				{8, lexer.DOT, "."},
				{9, lexer.IDENTITY, "base-iface-routes"},
				{27, lexer.PIPE, "|"},
				{29, lexer.IDENTITY, "routes"},
				{35, lexer.DOT, "."},
				{36, lexer.IDENTITY, "running"},
				{43, lexer.DOT, "."},
				{44, lexer.IDENTITY, "next-hop-interface"},
				{62, lexer.REPLACE, ":="},
				{64, lexer.STRING, "br1"},
				{68, lexer.EOF, ""}},
			}},
			{`capturer.base-iface-route | routes.running.state:="absent"`, expected{tokens: []lexer.Token{
				{0, lexer.IDENTITY, "capturer"},
				{8, lexer.DOT, "."},
				{9, lexer.IDENTITY, "base-iface-route"},
				{26, lexer.PIPE, "|"},
				{28, lexer.IDENTITY, "routes"},
				{34, lexer.DOT, "."},
				{35, lexer.IDENTITY, "running"},
				{42, lexer.DOT, "."},
				{43, lexer.IDENTITY, "state"},
				{48, lexer.REPLACE, ":="},
				{50, lexer.STRING, "absent"},
				{57, lexer.EOF, ""}},
			}},
			{`capturer.delete-primary-nic-routes.routes.running + capturer.bridge-routes.routes.running`, expected{tokens: []lexer.Token{
				{0, lexer.IDENTITY, "capturer"},
				{8, lexer.DOT, "."},
				{9, lexer.IDENTITY, "delete-primary-nic-routes"},
				{34, lexer.DOT, "."},
				{35, lexer.IDENTITY, "routes"},
				{41, lexer.DOT, "."},
				{42, lexer.IDENTITY, "running"},
				{50, lexer.MERGE, "+"},
				{52, lexer.IDENTITY, "capturer"},
				{60, lexer.DOT, "."},
				{61, lexer.IDENTITY, "bridge-routes"},
				{74, lexer.DOT, "."},
				{75, lexer.IDENTITY, "routes"},
				{81, lexer.DOT, "."},
				{82, lexer.IDENTITY, "running"},
				{88, lexer.EOF, ""}},
			}},
		})
	})
}

func runTest(t *testing.T, tests []test) {
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
