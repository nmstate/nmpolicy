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
	captureExpressions map[string]string,
	capturedStatesCache map[string]types.CapturedState,
	state []byte) (Result, error) {
	if len(captureExpressions) == 0 || len(state) == 0 && len(capturedStatesCache) == 0 {
		return Result{}, nil
	}

	capturedStates := filterCacheBasedOnExprCaptures(capturedStatesCache, captureExpressions)
	captureExpressions = filterOutExprBasedOnCachedCaptures(captureExpressions, capturedStatesCache)

	captureASTPool := map[string]ast.Node{}
	for captureEntryName, captureEntryExpression := range captureExpressions {
		captureEntryTokens, err := c.lexer.Lex(captureEntryExpression)
		if err != nil {
			return Result{}, fmt.Errorf("failed to resolve capture expression, err: %v", err)
		}

		captureEntryAST, err := c.parser.Parse(captureEntryTokens)
		if err != nil {
			return Result{}, fmt.Errorf("failed to resolve capture expression, err: %v", err)
		}

		captureASTPool[captureEntryName] = captureEntryAST
	}

	resolverResult, err := c.resolver.Resolve(captureASTPool, state)
	if err != nil {
		return Result{}, fmt.Errorf("failed to resolve capture expression, err: %v", err)
	}

	for captureEntryName, capturedState := range capturedStates {
		resolverResult.Marshaled[captureEntryName] = capturedState
	}
	return NewResult(c.lexer, c.parser, c.resolver, resolverResult), nil
}

func filterOutExprBasedOnCachedCaptures(captureExpressions map[string]string,
	capturedStates map[string]types.CapturedState) map[string]string {
	for captureEntryName := range capturedStates {
		delete(captureExpressions, captureEntryName)
	}
	return captureExpressions
}

func filterCacheBasedOnExprCaptures(capturedStates map[string]types.CapturedState,
	captureExpressions map[string]string) map[string]types.CapturedState {
	filteredCapturedStates := map[string]types.CapturedState{}

	for captureEntryName := range captureExpressions {
		if capturedState, ok := capturedStates[captureEntryName]; ok {
			state := append([]byte{}, capturedState.State...)

			filteredCapturedStates[captureEntryName] = types.CapturedState{
				State:    state,
				MetaInfo: capturedState.MetaInfo,
			}
		}
	}
	return filteredCapturedStates
}
