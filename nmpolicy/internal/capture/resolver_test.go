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
	"testing"

	assert "github.com/stretchr/testify/require"

	"github.com/nmstate/nmpolicy/nmpolicy/internal/capture"
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
		captureResolver := capture.NewResolver(lexerStub{}, parserStub{}, resolverStub{})
		result, err := captureResolver.Resolve(
			map[string]string{},
			map[string]types.CapturedState{"cap0": {State: []byte("some captured state")}},
			[]byte("some state"),
		)
		assert.NoError(t, err)

		assert.Nil(t, result.CapturedStates())
	})
}

func testNoCacheAndState(t *testing.T) {
	t.Run("resolve with no cache and state", func(t *testing.T) {
		captureResolver := capture.NewResolver(lexerStub{}, parserStub{}, resolverStub{})
		result, err := captureResolver.Resolve(
			map[string]string{"cap0": "my expression"},
			map[string]types.CapturedState{},
			[]byte{},
		)
		assert.NoError(t, err)

		assert.Nil(t, result.CapturedStates())
	})
}

func testAllCapturesCached(t *testing.T) {
	t.Run("resolve with all captures cached", func(t *testing.T) {
		capCache := map[string]types.CapturedState{
			"cap0": {State: []byte("some captured state")},
			"cap1": {State: []byte("another captured state")},
		}

		captureResolver := capture.NewResolver(lexerStub{}, parserStub{}, resolverStub{})
		result, err := captureResolver.Resolve(
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
		const captureEntryName = "cap0"

		captureResolver := capture.NewResolver(lexerStub{}, parserStub{}, resolverStub{})
		result, err := captureResolver.Resolve(
			map[string]string{captureEntryName: "my expression"},
			map[string]types.CapturedState{},
			[]byte("some state"),
		)
		assert.NoError(t, err)

		assert.Equal(t, map[string]types.CapturedState{
			captureEntryName: {State: []byte("resolver: parser: lexer: my expression")},
		}, result.CapturedStates())
	})
}

func testExpressionsWithPartialCache(t *testing.T) {
	t.Run("resolve with expressions and partial cache", func(t *testing.T) {
		const captureEntryName0 = "cap0"
		const captureEntryName1 = "cap1"

		capturedStatesCache := map[string]types.CapturedState{captureEntryName0: {State: []byte("some captured state")}}
		captureResolver := capture.NewResolver(lexerStub{}, parserStub{}, resolverStub{})

		result, err := captureResolver.Resolve(
			map[string]string{
				captureEntryName0: "my expression",
				captureEntryName1: "another expression",
			},
			capturedStatesCache,
			[]byte("some state"),
		)
		assert.NoError(t, err)

		expectedCapturedStates := capturedStatesCache
		expectedCapturedStates[captureEntryName1] = types.CapturedState{State: []byte("resolver: parser: lexer: another expression")}
		assert.Equal(t, expectedCapturedStates, result.CapturedStates())
	})
}

func testExpressionsWithOverCache(t *testing.T) {
	t.Run("resolve with cache that is not included in the expressions", func(t *testing.T) {
		const captureEntryName0 = "cap0"
		const captureEntryName1 = "cap1"

		capturedStatesCache := map[string]types.CapturedState{
			captureEntryName0: {State: []byte("some captured state")},
			captureEntryName1: {State: []byte("another captured state")},
		}
		captureResolver := capture.NewResolver(lexerStub{}, parserStub{}, resolverStub{})

		result, err := captureResolver.Resolve(
			map[string]string{
				captureEntryName0: "my expression",
			},
			capturedStatesCache,
			[]byte("some state"),
		)
		assert.NoError(t, err)

		expectedCapturedStates := map[string]types.CapturedState{
			captureEntryName0: {State: []byte("some captured state")},
		}
		assert.Equal(t, expectedCapturedStates, result.CapturedStates())
	})
}

func testLexFailure(t *testing.T) {
	t.Run("resolve fails due to lexing", func(t *testing.T) {
		captureResolver := capture.NewResolver(lexerStub{failLex: true}, parserStub{}, resolverStub{})
		_, err := captureResolver.Resolve(
			map[string]string{"cap0": "my expression"},
			map[string]types.CapturedState{},
			[]byte("some state"),
		)
		assert.Error(t, err)
	})
}

func testParseFailure(t *testing.T) {
	t.Run("resolve fails due to parsing", func(t *testing.T) {
		captureResolver := capture.NewResolver(lexerStub{}, parserStub{failParse: true}, resolverStub{})
		_, err := captureResolver.Resolve(
			map[string]string{"cap0": "my expression"},
			map[string]types.CapturedState{},
			[]byte("some state"),
		)
		assert.Error(t, err)
	})
}

func testResolveFailure(t *testing.T) {
	t.Run("resolve fails due to resolving", func(t *testing.T) {
		captureResolver := capture.NewResolver(lexerStub{}, parserStub{}, resolverStub{failResolve: true})
		_, err := captureResolver.Resolve(
			map[string]string{"cap0": "my expression"},
			map[string]types.CapturedState{},
			[]byte("some state"),
		)
		assert.Error(t, err)
	})
}
