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

package capture

import (
	"fmt"

	"github.com/nmstate/nmpolicy/nmpolicy/internal/ast"
	"github.com/nmstate/nmpolicy/nmpolicy/internal/lexer"
	"github.com/nmstate/nmpolicy/nmpolicy/internal/resolver"
	"github.com/nmstate/nmpolicy/nmpolicy/types"
)

type CaptureResolver struct {
	lexer    Lexer
	parser   Parser
	resolver Resolver
}

type Lexer interface {
	Lex(expression string) ([]lexer.Token, error)
}

type Parser interface {
	Parse([]lexer.Token) (ast.Node, error)
}

type Resolver interface {
	Resolve(astPool map[string]ast.Node, state []byte) (resolver.Result, error)
	ResolveCaptureEntryPath(captureEntryPathAST ast.Node, capturdStates map[string]map[string]interface{}) (interface{}, error)
}

func NewResolver(l Lexer, p Parser, r Resolver) CaptureResolver {
	return CaptureResolver{
		lexer:    l,
		parser:   p,
		resolver: r,
	}
}

func (c CaptureResolver) Resolve(
	capturesExpr map[string]string,
	capturesCache map[string]types.CapturedState,
	state []byte) (Result, error) {
	if len(capturesExpr) == 0 || len(state) == 0 && len(capturesCache) == 0 {
		return Result{}, nil
	}

	capturesState := filterCacheBasedOnExprCaptures(capturesCache, capturesExpr)
	capturesExpr = filterOutExprBasedOnCachedCaptures(capturesExpr, capturesCache)

	astPool := map[string]ast.Node{}
	for capID, capExpr := range capturesExpr {
		tokens, err := c.lexer.Lex(capExpr)
		if err != nil {
			return Result{}, fmt.Errorf("failed to resolve capture expression, err: %v", err)
		}

		astRoot, err := c.parser.Parse(tokens)
		if err != nil {
			return Result{}, fmt.Errorf("failed to resolve capture expression, err: %v", err)
		}

		astPool[capID] = astRoot
	}

	resolverResult, err := c.resolver.Resolve(astPool, state)
	if err != nil {
		return Result{}, fmt.Errorf("failed to resolve capture expression, err: %v", err)
	}

	for capID, capState := range capturesState {
		resolverResult.Marshaled[capID] = capState
	}
	return NewResult(c.lexer, c.parser, c.resolver, resolverResult), nil
}

func filterOutExprBasedOnCachedCaptures(capturesExpr map[string]string, capturesCache map[string]types.CapturedState) map[string]string {
	for capID := range capturesCache {
		delete(capturesExpr, capID)
	}
	return capturesExpr
}

func filterCacheBasedOnExprCaptures(capsState map[string]types.CapturedState, capsExpr map[string]string) map[string]types.CapturedState {
	caps := map[string]types.CapturedState{}

	for capID := range capsExpr {
		if capState, ok := capsState[capID]; ok {
			state := append([]byte{}, capState.State...)

			caps[capID] = types.CapturedState{
				State:    state,
				MetaInfo: capState.MetaInfo,
			}
		}
	}
	return caps
}
