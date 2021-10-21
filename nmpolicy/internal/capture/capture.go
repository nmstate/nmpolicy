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

package capture

import (
	"fmt"

	"github.com/nmstate/nmpolicy/nmpolicy/internal/ast"
	"github.com/nmstate/nmpolicy/nmpolicy/internal/lexer"
	"github.com/nmstate/nmpolicy/nmpolicy/types"
)

type CapsExpressions = map[types.CaptureID]types.Expression
type CapsState = map[types.CaptureID]types.CaptureState

type Capture struct {
	astPool  AstPooler
	lexer    Lexer
	parser   Parser
	resolver Resolver
}

type AstPooler interface {
	Add(id types.CaptureID, ast ast.Node)
	Range() ast.Pool
}

type Lexer interface {
	Lex(expression types.Expression) ([]lexer.Token, error)
}

type Parser interface {
	Parse([]lexer.Token) (ast.Node, error)
}

type Resolver interface {
	Resolve(astPool AstPooler, state types.NMState) (CapsState, error)
}

func New(astPool AstPooler, leXer Lexer, parser Parser, resolver Resolver) Capture {
	return Capture{
		astPool:  astPool,
		lexer:    leXer,
		parser:   parser,
		resolver: resolver,
	}
}

func (c Capture) Resolve(capturesExpr CapsExpressions, capturesCache CapsState, state types.NMState) (CapsState, error) {
	if len(capturesExpr) == 0 || len(state) == 0 && len(capturesCache) == 0 {
		return nil, nil
	}

	capturesExpr = filterOutCachedCaptures(capturesExpr, capturesCache)
	capturesState := newCapturesState(capturesCache)

	for capID, capExpr := range capturesExpr {
		tokens, err := c.lexer.Lex(capExpr)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve capture expression, err: %v", err)
		}

		astRoot, err := c.parser.Parse(tokens)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve capture expression, err: %v", err)
		}

		c.astPool.Add(capID, astRoot)
	}

	resolvedCapsState, err := c.resolver.Resolve(c.astPool, state)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve capture expression, err: %v", err)
	}

	for capID, capExpr := range resolvedCapsState {
		capturesState[capID] = capExpr
	}

	return capturesState, nil
}

func filterOutCachedCaptures(capturesExpr CapsExpressions, capturesCache CapsState) CapsExpressions {
	for capID := range capturesCache {
		delete(capturesExpr, capID)
	}
	return capturesExpr
}

func newCapturesState(c CapsState) CapsState {
	caps := CapsState{}

	for capID, capState := range c {
		state := append(types.NMState{}, capState.State...)

		caps[capID] = types.CaptureState{
			State:    state,
			MetaInfo: capState.MetaInfo,
		}
	}
	return caps
}
