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

package lexer

type TokenType int

const (
	EOF TokenType = iota
	IDENTITY
	NUMBER
	STRING

	DOT // .

	operatorsBegin
	PIPE     // |
	REPLACE  // :=
	EQFILTER // ==
	NEFILTER // !=
	MERGE    // +
	operatorsEnd

	TRUE  // true
	FALSE // false
)

var tokens = []string{
	EOF:      "EOF",
	IDENTITY: "IDENTITY",
	NUMBER:   "NUMBER",
	STRING:   "STRING",

	DOT:  "DOT",
	PIPE: "PIPE",

	REPLACE:  "REPLACE",
	EQFILTER: "EQFILTER",
	NEFILTER: "NEQFILTER",
	MERGE:    "MERGE",

	TRUE:  "TRUE",
	FALSE: "FALSE",
}

func (t TokenType) String() string {
	return tokens[t]
}

func (t TokenType) IsOperator() bool {
	return t > operatorsBegin && t < operatorsEnd
}

type Token struct {
	Position int
	Type     TokenType
	Literal  string
}
