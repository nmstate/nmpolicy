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

import (
	"fmt"

	"github.com/nmstate/nmpolicy/nmpolicy/internal/ast"

	"github.com/nmstate/nmpolicy/nmpolicy/internal/lexer"
)

type Parser struct {
	tokens          []lexer.Token
	currentTokenIdx int
	lastNode        *ast.Node
}

func New() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(tokens []lexer.Token) (ast.Node, error) {
	p.reset(tokens)
	node, err := p.parse()
	if err != nil {
		return ast.Node{}, err
	}
	return node, nil
}

func (p *Parser) parse() (ast.Node, error) {
	for {
		if p.currentToken() == nil {
			return ast.Node{}, nil
		} else if p.currentToken().Type == lexer.EOF {
			break
		} else if p.currentToken().Type == lexer.STRING {
			str := p.parseString()
			p.lastNode = &str
		} else if p.currentToken().Type == lexer.IDENTITY {
			identity := p.parseIdentity()
			p.lastNode = &identity
			path, err := p.parsePath()
			if err != nil {
				return ast.Node{}, err
			}
			p.lastNode = path
		} else if p.currentToken().Type == lexer.EQFILTER {
			eqfilter, err := p.parseEqFilter()
			if err != nil {
				return ast.Node{}, err
			}
			p.lastNode = eqfilter
		} else {
			return ast.Node{}, &InvalidExpressionError{fmt.Sprintf("unexpected token `%+v`", p.currentToken().Literal)}
		}
		p.nextToken()
	}
	return *p.lastNode, nil
}

func (p *Parser) nextToken() {
	if len(p.tokens) == 0 {
		return
	}
	if p.currentTokenIdx >= len(p.tokens)-1 {
		p.currentTokenIdx = len(p.tokens) - 1
	} else {
		p.currentTokenIdx++
	}
}

func (p *Parser) prevToken() {
	if len(p.tokens) == 0 {
		return
	}
	if p.currentTokenIdx > 0 {
		p.currentTokenIdx--
	}
	if p.currentTokenIdx >= len(p.tokens)-1 {
		p.currentTokenIdx = len(p.tokens) - 1
	}
}

func (p *Parser) currentToken() *lexer.Token {
	return &p.tokens[p.currentTokenIdx]
}

func (p *Parser) parseIdentity() ast.Node {
	return ast.Node{
		Meta:     ast.Meta{Position: p.currentToken().Position},
		Terminal: ast.Terminal{Identity: &p.currentToken().Literal},
	}
}

func (p *Parser) parseString() ast.Node {
	return ast.Node{
		Meta:     ast.Meta{Position: p.currentToken().Position},
		Terminal: ast.Terminal{String: &p.currentToken().Literal},
	}
}

func (p *Parser) parsePath() (*ast.Node, error) {
	operator := &ast.Node{
		Meta: ast.Meta{Position: p.currentToken().Position},
		Path: &ast.VariadicOperator{*p.lastNode},
	}
	for {
		p.nextToken()
		if p.currentToken().Type == lexer.DOT {
			p.nextToken()
			if p.currentToken().Type == lexer.IDENTITY {
				path := append(*operator.Path, p.parseIdentity())
				operator.Path = &path
			} else {
				return nil, &InvalidPathError{"missing identity after dot"}
			}
		} else if p.currentToken().Type != lexer.EOF && p.currentToken().Type != lexer.EQFILTER {
			return nil, &InvalidPathError{"missing dot"}
		} else {
			// Token has not being consumed let's go back.
			p.prevToken()
			break
		}
	}
	return operator, nil
}

func (p *Parser) parseEqFilter() (*ast.Node, error) {
	operator := &ast.Node{
		Meta:     ast.Meta{Position: p.currentToken().Position},
		EqFilter: &ast.TernaryOperator{},
	}
	if p.lastNode == nil {
		return nil, &InvalidEqualityFilter{"missing left hand argument"}
	}
	if p.lastNode.Path == nil {
		return nil, &InvalidEqualityFilter{"left hand argument is not a path"}
	}
	operator.EqFilter[0].Terminal = ast.CurrentStateIdentity()
	operator.EqFilter[1] = *p.lastNode

	p.nextToken()

	if p.currentToken().Type == lexer.STRING {
		operator.EqFilter[2] = p.parseString()
	} else if p.currentToken().Type != lexer.EOF {
		return nil, &InvalidEqualityFilter{"right hand argument is not a string"}
	}
	return operator, nil
}

func (p *Parser) reset(tokens []lexer.Token) {
	*p = *New()
	p.tokens = tokens
}
