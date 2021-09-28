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

import (
	"fmt"
	"strings"

	"github.com/nmstate/nmpolicy/nmpolicy/internal/lexer/scanner"
)

// Lexer struct is used to tokenize values returned by a reader.
type lexer struct {
	scn *scanner.Scanner
}

// NewLexer construct a Lexer using reader as the input.
func New() *lexer {
	return &lexer{}
}

// Lex scans the input for the next token.
// It returns a Token struct with position, type, and the literal value.
func (l *lexer) Lex(expression string) ([]Token, error) {
	l.scn = scanner.New(strings.NewReader(expression))
	// keep looping until we return a token
	tokens := []Token{}
	for {
		token, err := l.lex()
		if err != nil {
			return nil, err
		}
		if token == nil {
			continue
		}
		tokens = append(tokens, *token)
		if token.Type == EOF {
			break
		}
	}
	return tokens, nil
}

func (l *lexer) lex() (*Token, error) {
	for {
		err := l.scn.Next()
		if err != nil {
			return nil, err
		}
		token, err := l.lexCurrentRune()
		if err != nil {
			return nil, err
		}
		if token == nil {
			continue
		}
		return token, nil
	}
}

func (l *lexer) lexCurrentRune() (*Token, error) {
	if l.isEOF() {
		return &Token{l.scn.Position(), EOF, ""}, nil
	} else if l.isSpace() {
		return nil, nil
	} else if l.isDigit() {
		return l.lexNumber()
	}
	return nil, fmt.Errorf("illegal rune %s", string(l.scn.Rune()))
}

func (l *lexer) lexNumber() (*Token, error) {
	token := &Token{l.scn.Position(), NUMBER, string(l.scn.Rune())}
	for {
		if err := l.scn.Next(); err != nil {
			return nil, err
		}

		if l.isEOF() || l.isSpace() {
			// If it's EOF or space we have finish here
			return token, nil
		} else if l.isDigit() {
			token.Literal += string(l.scn.Rune())
		} else {
			// nmpolicy supports only simple numbers with just digist
			return nil, fmt.Errorf("invalid number format (%s is not a digit)", string(l.scn.Rune()))
		}
	}
}
