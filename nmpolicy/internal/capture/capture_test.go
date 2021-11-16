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

package capture_test

import (
	"fmt"
	"testing"

	assert "github.com/stretchr/testify/require"

	"github.com/nmstate/nmpolicy/nmpolicy/internal/ast"
	"github.com/nmstate/nmpolicy/nmpolicy/internal/capture"
	"github.com/nmstate/nmpolicy/nmpolicy/internal/lexer"
	"github.com/nmstate/nmpolicy/nmpolicy/internal/resolver"
	"github.com/nmstate/nmpolicy/nmpolicy/types"
)

func TestBasicPolicy(t *testing.T) {
	t.Run("Capture", func(t *testing.T) {
		testNoExpressions(t)
		testNoCacheAndState(t)
		testAllCapturesCached(t)
		testResolvingExpressions(t)
		testExpressionsWithPartialCache(t)
		testExpressionsWithOverCache(t)

		testLexFailure(t)
		testParseFailure(t)
		testResolveFailure(t)
	})
}

func testNoExpressions(t *testing.T) {
	t.Run("resolve with no expression", func(t *testing.T) {
		capCtrl := capture.New(lexerStub{}, parserStub{}, resolverStub{})
		result, err := capCtrl.Resolve(
			map[string]string{},
			map[string]types.CaptureState{"cap0": {State: []byte("some captured state")}},
			[]byte("some state"),
		)
		assert.NoError(t, err)

		assert.Nil(t, result.CapturedStates())
	})
}

func testNoCacheAndState(t *testing.T) {
	t.Run("resolve with no cache and state", func(t *testing.T) {
		capCtrl := capture.New(lexerStub{}, parserStub{}, resolverStub{})
		result, err := capCtrl.Resolve(
			map[string]string{"cap0": "my expression"},
			map[string]types.CaptureState{},
			[]byte{},
		)
		assert.NoError(t, err)

		assert.Nil(t, result.CapturedStates())
	})
}

func testAllCapturesCached(t *testing.T) {
	t.Run("resolve with all captures cached", func(t *testing.T) {
		capCache := map[string]types.CaptureState{
			"cap0": {State: []byte("some captured state")},
			"cap1": {State: []byte("another captured state")},
		}

		capCtrl := capture.New(lexerStub{}, parserStub{}, resolverStub{})
		result, err := capCtrl.Resolve(
			map[string]string{
				"cap0": "my expression",
				"cap1": "another expression",
			},
			capCache,
			[]byte{},
		)
		assert.NoError(t, err)

		assert.Equal(t, capCache, result.CapturedStates())
	})
}

func testResolvingExpressions(t *testing.T) {
	t.Run("resolve expressions", func(t *testing.T) {
		const capID = "cap0"

		capCtrl := capture.New(lexerStub{}, parserStub{}, resolverStub{})
		result, err := capCtrl.Resolve(
			map[string]string{capID: "my expression"},
			map[string]types.CaptureState{},
			[]byte("some state"),
		)
		assert.NoError(t, err)

		assert.Equal(t, map[string]types.CaptureState{capID: {State: []byte("resolver: parser: lexer: my expression")}}, result.CapturedStates())
	})
}

func testExpressionsWithPartialCache(t *testing.T) {
	t.Run("resolve with expressions and partial cache", func(t *testing.T) {
		const capID0 = "cap0"
		const capID1 = "cap1"

		capCache := map[string]types.CaptureState{capID0: {State: []byte("some captured state")}}
		capCtrl := capture.New(lexerStub{}, parserStub{}, resolverStub{})

		result, err := capCtrl.Resolve(
			map[string]string{
				capID0: "my expression",
				capID1: "another expression",
			},
			capCache,
			[]byte("some state"),
		)
		assert.NoError(t, err)

		expectedCaps := capCache
		expectedCaps[capID1] = types.CaptureState{State: []byte("resolver: parser: lexer: another expression")}
		assert.Equal(t, expectedCaps, result.CapturedStates())
	})
}

func testExpressionsWithOverCache(t *testing.T) {
	t.Run("resolve with cache that is not included in the expressions", func(t *testing.T) {
		const capID0 = "cap0"
		const capID1 = "cap1"

		capCache := map[string]types.CaptureState{
			capID0: {State: []byte("some captured state")},
			capID1: {State: []byte("another captured state")},
		}
		capCtrl := capture.New(lexerStub{}, parserStub{}, resolverStub{})

		result, err := capCtrl.Resolve(
			map[string]string{
				capID0: "my expression",
			},
			capCache,
			[]byte("some state"),
		)
		assert.NoError(t, err)

		expectedCaps := map[string]types.CaptureState{
			capID0: {State: []byte("some captured state")},
		}
		assert.Equal(t, expectedCaps, result.CapturedStates())
	})
}

func testLexFailure(t *testing.T) {
	t.Run("resolve fails due to lexing", func(t *testing.T) {
		capCtrl := capture.New(lexerStub{failLex: true}, parserStub{}, resolverStub{})
		_, err := capCtrl.Resolve(
			map[string]string{"cap0": "my expression"},
			map[string]types.CaptureState{},
			[]byte("some state"),
		)
		assert.Error(t, err)
	})
}

func testParseFailure(t *testing.T) {
	t.Run("resolve fails due to parsing", func(t *testing.T) {
		capCtrl := capture.New(lexerStub{}, parserStub{failParse: true}, resolverStub{})
		_, err := capCtrl.Resolve(
			map[string]string{"cap0": "my expression"},
			map[string]types.CaptureState{},
			[]byte("some state"),
		)
		assert.Error(t, err)
	})
}

func testResolveFailure(t *testing.T) {
	t.Run("resolve fails due to resolving", func(t *testing.T) {
		capCtrl := capture.New(lexerStub{}, parserStub{}, resolverStub{failResolve: true})
		_, err := capCtrl.Resolve(
			map[string]string{"cap0": "my expression"},
			map[string]types.CaptureState{},
			[]byte("some state"),
		)
		assert.Error(t, err)
	})
}

type lexerStub struct {
	failLex bool
}

func (l lexerStub) Lex(expression string) ([]lexer.Token, error) {
	if l.failLex {
		return nil, fmt.Errorf("lex failed")
	}
	literal := fmt.Sprintf("lexer: %s", expression)
	return []lexer.Token{{Literal: literal}}, nil
}

type parserStub struct {
	failParse bool
}

func (p parserStub) Parse(tokens []lexer.Token) (ast.Node, error) {
	if p.failParse {
		return ast.Node{}, fmt.Errorf("parse failed")
	}
	literal := fmt.Sprintf("parser: %s", tokens[0].Literal)
	return ast.Node{Terminal: ast.Terminal{String: &literal}}, nil
}

type resolverStub struct {
	failResolve bool
}

func (r resolverStub) Resolve(astPool map[string]ast.Node, state []byte) (resolver.Result, error) {
	if r.failResolve {
		return resolver.Result{}, fmt.Errorf("resolve failed")
	}

	capsState := map[string]types.CaptureState{}
	for id, entry := range astPool {
		capsState[id] = types.CaptureState{State: []byte(fmt.Sprintf("resolver: %s", *entry.String))}
	}

	return resolver.Result{Marshaled: capsState}, nil
}

func (r resolverStub) ResolveCaptureEntryPath(captureEntryPathAST ast.Node,
	capturedStates map[string]map[string]interface{}) (interface{}, error) {
	if r.failResolve {
		return nil, fmt.Errorf("resolve capture entry path failed")
	}
	return fmt.Sprintf("resolver: %s", *captureEntryPathAST.String), nil
}
