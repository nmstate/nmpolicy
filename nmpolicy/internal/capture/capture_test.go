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

package capture_test

import (
	"fmt"
	"testing"

	assert "github.com/stretchr/testify/require"

	"github.com/nmstate/nmpolicy/nmpolicy/internal/ast"
	"github.com/nmstate/nmpolicy/nmpolicy/internal/capture"
	"github.com/nmstate/nmpolicy/nmpolicy/internal/lexer"
	"github.com/nmstate/nmpolicy/nmpolicy/types"
)

func TestBasicPolicy(t *testing.T) {
	t.Run("Capture", func(t *testing.T) {
		testNoExpressions(t)
		testNoCacheAndState(t)
		testAllCapturesCached(t)
		testResolvingExpressions(t)
		testExpressionsWithPartialCache(t)

		testLexFailure(t)
		testParseFailure(t)
		testResolveFailure(t)
	})
}

func testNoExpressions(t *testing.T) {
	t.Run("resolve with no expression", func(t *testing.T) {
		capCtrl := capture.New(ast.NewPool(),
			lexerFactory(lexerStub{}),
			parserFactory(parserStub{}),
			resolverFactory(resolverStub{}))
		resolvedCaps, err := capCtrl.Resolve(
			capture.CapsExpressions{},
			map[types.CaptureID]types.CaptureState{"cap0": {State: []byte("some captured state")}},
			types.NMState("some state"),
		)
		assert.NoError(t, err)

		assert.Nil(t, resolvedCaps)
	})
}

func testNoCacheAndState(t *testing.T) {
	t.Run("resolve with no cache and state", func(t *testing.T) {
		capCtrl := capture.New(ast.NewPool(),
			lexerFactory(lexerStub{}),
			parserFactory(parserStub{}),
			resolverFactory(resolverStub{}))
		resolvedCaps, err := capCtrl.Resolve(
			map[types.CaptureID]types.Expression{"cap0": "my expression"},
			capture.CapsState{},
			types.NMState{},
		)
		assert.NoError(t, err)

		assert.Nil(t, resolvedCaps)
	})
}

func testAllCapturesCached(t *testing.T) {
	t.Run("resolve with all captures cached", func(t *testing.T) {
		capCache := map[types.CaptureID]types.CaptureState{
			"cap0": {State: []byte("some captured state")},
			"cap1": {State: []byte("another captured state")},
		}

		capCtrl := capture.New(ast.NewPool(),
			lexerFactory(lexerStub{}),
			parserFactory(parserStub{}),
			resolverFactory(resolverStub{}))
		resolvedCaps, err := capCtrl.Resolve(
			map[types.CaptureID]types.Expression{
				"cap0": "my expression",
				"cap1": "another expression",
			},
			capCache,
			types.NMState{},
		)
		assert.NoError(t, err)

		assert.Equal(t, capCache, resolvedCaps)
	})
}

func testResolvingExpressions(t *testing.T) {
	t.Run("resolve expressions", func(t *testing.T) {
		const capID = "cap0"

		capCtrl := capture.New(ast.NewPool(),
			lexerFactory(lexerStub{}),
			parserFactory(parserStub{}),
			resolverFactory(resolverStub{}))
		resolvedCaps, err := capCtrl.Resolve(
			map[types.CaptureID]types.Expression{capID: "my expression"},
			capture.CapsState{},
			types.NMState("some state"),
		)
		assert.NoError(t, err)

		assert.Equal(t, capture.CapsState{capID: types.CaptureState{}}, resolvedCaps)
	})
}

func testExpressionsWithPartialCache(t *testing.T) {
	t.Run("resolve with expressions and partial cache", func(t *testing.T) {
		const capID0 = "cap0"
		const capID1 = "cap1"

		capCache := map[types.CaptureID]types.CaptureState{capID0: {State: []byte("some captured state")}}
		capCtrl := capture.New(ast.NewPool(),
			lexerFactory(lexerStub{}),
			parserFactory(parserStub{}),
			resolverFactory(resolverStub{}))

		resolvedCaps, err := capCtrl.Resolve(
			map[types.CaptureID]types.Expression{
				capID0: "my expression",
				capID1: "another expression",
			},
			capCache,
			types.NMState("some state"),
		)
		assert.NoError(t, err)

		expectedCaps := capCache
		expectedCaps[capID1] = types.CaptureState{}
		assert.Equal(t, expectedCaps, resolvedCaps)
	})
}

func testLexFailure(t *testing.T) {
	t.Run("resolve fails due to lexing", func(t *testing.T) {
		capCtrl := capture.New(ast.NewPool(),
			lexerFactory(lexerStub{failLex: true}),
			parserFactory(parserStub{}),
			resolverFactory(resolverStub{}))
		_, err := capCtrl.Resolve(
			map[types.CaptureID]types.Expression{"cap0": "my expression"},
			capture.CapsState{},
			types.NMState("some state"),
		)
		assert.Error(t, err)
	})
}

func testParseFailure(t *testing.T) {
	t.Run("resolve fails due to parsing", func(t *testing.T) {
		capCtrl := capture.New(ast.NewPool(),
			lexerFactory(lexerStub{failLex: true}),
			parserFactory(parserStub{}),
			resolverFactory(resolverStub{}))
		_, err := capCtrl.Resolve(
			map[types.CaptureID]types.Expression{"cap0": "my expression"},
			capture.CapsState{},
			types.NMState("some state"),
		)
		assert.Error(t, err)
	})
}

func testResolveFailure(t *testing.T) {
	t.Run("resolve fails due to resolving", func(t *testing.T) {
		capCtrl := capture.New(ast.NewPool(),
			lexerFactory(lexerStub{failLex: true}),
			parserFactory(parserStub{}),
			resolverFactory(resolverStub{failResolve: true}))
		_, err := capCtrl.Resolve(
			map[types.CaptureID]types.Expression{"cap0": "my expression"},
			capture.CapsState{},
			types.NMState("some state"),
		)
		assert.Error(t, err)
	})
}

type lexerStub struct {
	failLex    bool
	expression types.Expression
}

func (l lexerStub) Lex() ([]lexer.Token, error) {
	if l.failLex {
		return nil, fmt.Errorf("lex failed")
	}
	return nil, nil
}

func lexerFactory(stub lexerStub) capture.LexerFactory {
	return func(expression types.Expression) capture.Lexer {
		stub.expression = expression
		return &stub
	}
}

type parserStub struct {
	failParse bool
	tokens    []lexer.Token
}

func (p parserStub) Parse() (ast.Node, error) {
	if p.failParse {
		return ast.Node{}, fmt.Errorf("parse failed")
	}
	return ast.Node{}, nil
}

func parserFactory(stub parserStub) capture.ParserFactory {
	return func(tokens []lexer.Token) capture.Parser {
		stub.tokens = tokens
		return &stub
	}
}

type resolverStub struct {
	failResolve bool
	state       types.NMState
	astPool     capture.AstPooler
}

func resolverFactory(stub resolverStub) capture.ResolverFactory {
	return func(state types.NMState, astPool capture.AstPooler) capture.Resolver {
		stub.state = state
		stub.astPool = astPool
		return &stub
	}
}

func (r resolverStub) Resolve() (map[types.CaptureID]types.CaptureState, error) {
	if r.failResolve {
		return nil, fmt.Errorf("resolve failed")
	}

	capsState := capture.CapsState{}
	for id := range r.astPool.Range() {
		capsState[id] = types.CaptureState{}
	}

	return capsState, nil
}
