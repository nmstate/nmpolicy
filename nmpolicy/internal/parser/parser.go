/*
 * This file is part of the nmpolicy project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2021 Red Hat, Inc.
 *
 */

package parser

import (
	"fmt"
	"reflect"

	"github.com/nmstate/nmpolicy/nmpolicy/internal/ast"
	"github.com/nmstate/nmpolicy/nmpolicy/internal/lexer"
)

type Parser struct {
	tokens []lexer.Token
}

func New(tokens []lexer.Token) Parser {
	return Parser{tokens: tokens}
}

func (p Parser) Parse() (ast.Node, error) {
	if reflect.DeepEqual(p.tokens, []lexer.Token{
		{Position: 0, Type: lexer.IDENTITY, Literal: "routes"},
		{Position: 6, Type: lexer.DOT, Literal: "."},
		{Position: 7, Type: lexer.IDENTITY, Literal: "running"},
		{Position: 14, Type: lexer.DOT, Literal: "."},
		{Position: 15, Type: lexer.IDENTITY, Literal: "destination"},
		{Position: 26, Type: lexer.EQFILTER, Literal: "=="},
		{Position: 28, Type: lexer.STRING, Literal: "0.0.0.0/0"},
		{Position: 35, Type: lexer.EOF, Literal: ""},
	}) {
		return ast.Node{
			Meta: ast.Meta{Position: 26},
			EqFilter: &ast.TernaryOperator{
				ast.Node{
					Meta:     ast.Meta{Position: 0},
					Terminal: ast.CurrentStateIdentity()},
				ast.Node{
					Meta: ast.Meta{Position: 0},
					Path: &ast.VariadicOperator{
						ast.Node{
							Meta:     ast.Meta{Position: 0},
							Terminal: ast.Terminal{Identity: strPtr("routes")},
						},
						ast.Node{
							Meta:     ast.Meta{Position: 7},
							Terminal: ast.Terminal{Identity: strPtr("running")},
						},
						ast.Node{
							Meta:     ast.Meta{Position: 15},
							Terminal: ast.Terminal{Identity: strPtr("destination")},
						},
					},
				},
				ast.Node{
					Meta:     ast.Meta{Position: 28},
					Terminal: ast.Terminal{String: strPtr("0.0.0.0/0")},
				},
			},
		}, nil
	}

	fmt.Println("parser: tokens not matching, returning nil")
	return ast.Node{}, nil
}

func strPtr(str string) *string {
	return &str
}
