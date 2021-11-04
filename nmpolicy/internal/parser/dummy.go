// Copyright 2021 The NMPolicy Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package parser

import (
	"reflect"

	"github.com/nmstate/nmpolicy/nmpolicy/internal/ast"
	"github.com/nmstate/nmpolicy/nmpolicy/internal/lexer"
)

var (
	zero         = 0
	DefaultGwAST = ast.Node{
		Meta: ast.Meta{Position: lexer.DefaultGwTokens[5].Position},
		EqFilter: &ast.TernaryOperator{
			ast.Node{
				Meta:     ast.Meta{Position: 0},
				Terminal: ast.CurrentStateIdentity()},
			ast.Node{
				Meta: ast.Meta{Position: 0},
				Path: &ast.VariadicOperator{
					ast.Node{
						Meta:     ast.Meta{Position: lexer.DefaultGwTokens[0].Position},
						Terminal: ast.Terminal{Identity: &lexer.DefaultGwTokens[0].Literal},
					},
					ast.Node{
						Meta:     ast.Meta{Position: lexer.DefaultGwTokens[2].Position},
						Terminal: ast.Terminal{Identity: &lexer.DefaultGwTokens[2].Literal},
					},
					ast.Node{
						Meta:     ast.Meta{Position: lexer.DefaultGwTokens[4].Position},
						Terminal: ast.Terminal{Identity: &lexer.DefaultGwTokens[4].Literal},
					},
				},
			},
			ast.Node{
				Meta:     ast.Meta{Position: lexer.DefaultGwTokens[6].Position},
				Terminal: ast.Terminal{String: &lexer.DefaultGwTokens[6].Literal},
			},
		},
	}

	BaseIfaceRoutesAST = ast.Node{
		Meta: ast.Meta{Position: lexer.BaseIfaceRoutesTokens[5].Position},
		EqFilter: &ast.TernaryOperator{
			ast.Node{
				Meta:     ast.Meta{Position: 0},
				Terminal: ast.CurrentStateIdentity()},
			ast.Node{
				Meta: ast.Meta{Position: 0},
				Path: &ast.VariadicOperator{
					ast.Node{
						Meta:     ast.Meta{Position: lexer.BaseIfaceRoutesTokens[0].Position},
						Terminal: ast.Terminal{Identity: &lexer.BaseIfaceRoutesTokens[0].Literal},
					},
					ast.Node{
						Meta:     ast.Meta{Position: lexer.BaseIfaceRoutesTokens[2].Position},
						Terminal: ast.Terminal{Identity: &lexer.BaseIfaceRoutesTokens[2].Literal},
					},
					ast.Node{
						Meta:     ast.Meta{Position: lexer.BaseIfaceRoutesTokens[4].Position},
						Terminal: ast.Terminal{Identity: &lexer.BaseIfaceRoutesTokens[4].Literal},
					},
				},
			},
			ast.Node{
				Meta: ast.Meta{Position: lexer.BaseIfaceRoutesTokens[6].Position},
				Path: &ast.VariadicOperator{
					ast.Node{
						Meta:     ast.Meta{Position: lexer.BaseIfaceRoutesTokens[6].Position},
						Terminal: ast.Terminal{Identity: &lexer.BaseIfaceRoutesTokens[6].Literal},
					},
					ast.Node{
						Meta:     ast.Meta{Position: lexer.BaseIfaceRoutesTokens[8].Position},
						Terminal: ast.Terminal{Identity: &lexer.BaseIfaceRoutesTokens[8].Literal},
					},
					ast.Node{
						Meta:     ast.Meta{Position: lexer.BaseIfaceRoutesTokens[10].Position},
						Terminal: ast.Terminal{Identity: &lexer.BaseIfaceRoutesTokens[10].Literal},
					},
					ast.Node{
						Meta:     ast.Meta{Position: lexer.BaseIfaceRoutesTokens[12].Position},
						Terminal: ast.Terminal{Identity: &lexer.BaseIfaceRoutesTokens[12].Literal},
					},
					ast.Node{
						Meta:     ast.Meta{Position: lexer.BaseIfaceRoutesTokens[14].Position},
						Terminal: ast.Terminal{Number: &zero},
					},
					ast.Node{
						Meta:     ast.Meta{Position: lexer.BaseIfaceRoutesTokens[16].Position},
						Terminal: ast.Terminal{Identity: &lexer.BaseIfaceRoutesTokens[16].Literal},
					},
				},
			},
		},
	}
	BaseIfaceAST = ast.Node{
		Meta: ast.Meta{Position: lexer.BaseIfaceTokens[3].Position},
		EqFilter: &ast.TernaryOperator{
			ast.Node{
				Meta:     ast.Meta{Position: 0},
				Terminal: ast.CurrentStateIdentity()},
			ast.Node{
				Meta: ast.Meta{Position: 0},
				Path: &ast.VariadicOperator{
					ast.Node{
						Meta:     ast.Meta{Position: lexer.BaseIfaceTokens[0].Position},
						Terminal: ast.Terminal{Identity: &lexer.BaseIfaceTokens[0].Literal},
					},
					ast.Node{
						Meta:     ast.Meta{Position: lexer.BaseIfaceTokens[2].Position},
						Terminal: ast.Terminal{Identity: &lexer.BaseIfaceTokens[2].Literal},
					},
				},
			},
			ast.Node{
				Meta: ast.Meta{Position: lexer.BaseIfaceTokens[4].Position},
				Path: &ast.VariadicOperator{
					ast.Node{
						Meta:     ast.Meta{Position: lexer.BaseIfaceTokens[4].Position},
						Terminal: ast.Terminal{Identity: &lexer.BaseIfaceTokens[4].Literal},
					},
					ast.Node{
						Meta:     ast.Meta{Position: lexer.BaseIfaceTokens[6].Position},
						Terminal: ast.Terminal{Identity: &lexer.BaseIfaceTokens[6].Literal},
					},
					ast.Node{
						Meta:     ast.Meta{Position: lexer.BaseIfaceTokens[8].Position},
						Terminal: ast.Terminal{Identity: &lexer.BaseIfaceTokens[8].Literal},
					},
					ast.Node{
						Meta:     ast.Meta{Position: lexer.BaseIfaceTokens[10].Position},
						Terminal: ast.Terminal{Identity: &lexer.BaseIfaceTokens[10].Literal},
					},
					ast.Node{
						Meta:     ast.Meta{Position: lexer.BaseIfaceTokens[12].Position},
						Terminal: ast.Terminal{Number: &zero},
					},
					ast.Node{
						Meta:     ast.Meta{Position: lexer.BaseIfaceTokens[14].Position},
						Terminal: ast.Terminal{Identity: &lexer.BaseIfaceTokens[14].Literal},
					},
				},
			},
		},
	}

	BridgeRoutesAST = ast.Node{
		Meta: ast.Meta{Position: lexer.BridgeRoutesTokens[9].Position},
		Replace: &ast.TernaryOperator{
			ast.Node{
				Meta: ast.Meta{Position: lexer.BridgeRoutesTokens[3].Position},
				Pipe: &ast.UnaryOperator{
					Meta: ast.Meta{Position: lexer.BridgeRoutesTokens[4].Position},
					Path: &ast.VariadicOperator{
						ast.Node{
							Meta:     ast.Meta{Position: lexer.BridgeRoutesTokens[0].Position},
							Terminal: ast.Terminal{Identity: &lexer.BridgeRoutesTokens[0].Literal},
						},
						ast.Node{
							Meta:     ast.Meta{Position: lexer.BridgeRoutesTokens[2].Position},
							Terminal: ast.Terminal{Identity: &lexer.BridgeRoutesTokens[2].Literal},
						},
					},
				},
			},
			ast.Node{
				Meta: ast.Meta{Position: 0},
				Path: &ast.VariadicOperator{
					ast.Node{
						Meta:     ast.Meta{Position: lexer.BridgeRoutesTokens[4].Position},
						Terminal: ast.Terminal{Identity: &lexer.BridgeRoutesTokens[4].Literal},
					},
					ast.Node{
						Meta:     ast.Meta{Position: lexer.BridgeRoutesTokens[6].Position},
						Terminal: ast.Terminal{Identity: &lexer.BridgeRoutesTokens[6].Literal},
					},
					ast.Node{
						Meta:     ast.Meta{Position: lexer.BridgeRoutesTokens[8].Position},
						Terminal: ast.Terminal{Identity: &lexer.BridgeRoutesTokens[8].Literal},
					},
				},
			},
			ast.Node{
				Meta:     ast.Meta{Position: lexer.BridgeRoutesTokens[10].Position},
				Terminal: ast.Terminal{String: &lexer.BridgeRoutesTokens[10].Literal},
			},
		},
	}
	DeleteBaseIfaceRoutesAST = ast.Node{
		Meta: ast.Meta{Position: lexer.DeleteBaseIfaceRoutesTokens[9].Position},
		Replace: &ast.TernaryOperator{
			ast.Node{
				Meta: ast.Meta{Position: lexer.DeleteBaseIfaceRoutesTokens[3].Position},
				Pipe: &ast.UnaryOperator{
					Meta: ast.Meta{Position: lexer.DeleteBaseIfaceRoutesTokens[4].Position},
					Path: &ast.VariadicOperator{
						ast.Node{
							Meta:     ast.Meta{Position: lexer.DeleteBaseIfaceRoutesTokens[0].Position},
							Terminal: ast.Terminal{Identity: &lexer.DeleteBaseIfaceRoutesTokens[0].Literal},
						},
						ast.Node{
							Meta:     ast.Meta{Position: lexer.DeleteBaseIfaceRoutesTokens[2].Position},
							Terminal: ast.Terminal{Identity: &lexer.DeleteBaseIfaceRoutesTokens[2].Literal},
						},
					},
				},
			},
			ast.Node{
				Path: &ast.VariadicOperator{
					ast.Node{
						Meta:     ast.Meta{Position: lexer.DeleteBaseIfaceRoutesTokens[4].Position},
						Terminal: ast.Terminal{Identity: &lexer.DeleteBaseIfaceRoutesTokens[4].Literal},
					},
					ast.Node{
						Meta:     ast.Meta{Position: lexer.DeleteBaseIfaceRoutesTokens[6].Position},
						Terminal: ast.Terminal{Identity: &lexer.DeleteBaseIfaceRoutesTokens[6].Literal},
					},
					ast.Node{
						Meta:     ast.Meta{Position: lexer.DeleteBaseIfaceRoutesTokens[8].Position},
						Terminal: ast.Terminal{Identity: &lexer.DeleteBaseIfaceRoutesTokens[8].Literal},
					},
				},
			},
			ast.Node{
				Meta:     ast.Meta{Position: lexer.DeleteBaseIfaceRoutesTokens[10].Position},
				Terminal: ast.Terminal{String: &lexer.DeleteBaseIfaceRoutesTokens[10].Literal},
			},
		},
	}
	BridgeRoutesTakeoverAST = ast.Node{
		Meta: ast.Meta{Position: lexer.BridgeRoutesTakeoverTokens[7].Position},
		Merge: &ast.BinaryOperator{
			ast.Node{
				Meta: ast.Meta{Position: 0},
				Path: &ast.VariadicOperator{
					ast.Node{
						Meta:     ast.Meta{Position: lexer.BridgeRoutesTakeoverTokens[0].Position},
						Terminal: ast.Terminal{Identity: &lexer.BridgeRoutesTakeoverTokens[0].Literal},
					},
					ast.Node{
						Meta:     ast.Meta{Position: lexer.BridgeRoutesTakeoverTokens[2].Position},
						Terminal: ast.Terminal{Identity: &lexer.BridgeRoutesTakeoverTokens[2].Literal},
					},
					ast.Node{
						Meta:     ast.Meta{Position: lexer.BridgeRoutesTakeoverTokens[4].Position},
						Terminal: ast.Terminal{Identity: &lexer.BridgeRoutesTakeoverTokens[4].Literal},
					},
					ast.Node{
						Meta:     ast.Meta{Position: lexer.BridgeRoutesTakeoverTokens[6].Position},
						Terminal: ast.Terminal{Identity: &lexer.BridgeRoutesTakeoverTokens[6].Literal},
					},
				},
			},
			ast.Node{
				Meta: ast.Meta{Position: lexer.BridgeRoutesTakeoverTokens[8].Position},
				Path: &ast.VariadicOperator{
					ast.Node{
						Meta:     ast.Meta{Position: lexer.BridgeRoutesTakeoverTokens[8].Position},
						Terminal: ast.Terminal{Identity: &lexer.BridgeRoutesTakeoverTokens[8].Literal},
					},
					ast.Node{
						Meta:     ast.Meta{Position: lexer.BridgeRoutesTakeoverTokens[10].Position},
						Terminal: ast.Terminal{Identity: &lexer.BridgeRoutesTakeoverTokens[10].Literal},
					},
					ast.Node{
						Meta:     ast.Meta{Position: lexer.BridgeRoutesTakeoverTokens[12].Position},
						Terminal: ast.Terminal{Identity: &lexer.BridgeRoutesTakeoverTokens[12].Literal},
					},
					ast.Node{
						Meta:     ast.Meta{Position: lexer.BridgeRoutesTakeoverTokens[14].Position},
						Terminal: ast.Terminal{Identity: &lexer.BridgeRoutesTakeoverTokens[14].Literal},
					},
				},
			},
		},
	}
)

func dummy(tokens []lexer.Token) *ast.Node {
	if reflect.DeepEqual(tokens, lexer.DefaultGwTokens) {
		return &DefaultGwAST
	}
	if reflect.DeepEqual(tokens, lexer.BaseIfaceRoutesTokens) {
		return &BaseIfaceRoutesAST
	}
	if reflect.DeepEqual(tokens, lexer.BaseIfaceTokens) {
		return &BaseIfaceAST
	}
	if reflect.DeepEqual(tokens, lexer.BridgeRoutesTokens) {
		return &BridgeRoutesAST
	}
	if reflect.DeepEqual(tokens, lexer.DeleteBaseIfaceRoutesTokens) {
		return &DeleteBaseIfaceRoutesAST
	}
	if reflect.DeepEqual(tokens, lexer.BridgeRoutesTakeoverTokens) {
		return &BridgeRoutesTakeoverAST
	}
	return nil
}
