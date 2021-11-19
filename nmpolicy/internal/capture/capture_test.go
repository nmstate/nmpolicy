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
	"github.com/nmstate/nmpolicy/nmpolicy/internal/types"
	"github.com/nmstate/nmpolicy/nmpolicy/internal/types/typestest"
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
		resolvedCaps, err := captureResolverWithDefaultStubs().Resolve(
			types.CaptureExpressions{},
			types.CapturedStates{"cap0": {State: typestest.ToNMState(t, "name: some captured state")}},
			typestest.ToNMState(t, "name: some state"),
		)
		assert.NoError(t, err)

		assert.Nil(t, resolvedCaps)
	})
}

func testNoCacheAndState(t *testing.T) {
	t.Run("resolve with no cache and state", func(t *testing.T) {
		resolvedCaps, err := captureResolverWithDefaultStubs().Resolve(
			types.CaptureExpressions{"cap0": "my expression"},
			types.CapturedStates{},
			types.NMState{},
		)
		assert.NoError(t, err)

		assert.Nil(t, resolvedCaps)
	})
}

func testAllCapturesCached(t *testing.T) {
	t.Run("resolve with all captures cached", func(t *testing.T) {
		capCache := types.CapturedStates{
			"cap0": {State: typestest.ToNMState(t, "name: some captured state")},
			"cap1": {State: typestest.ToNMState(t, "name: another captured state")},
		}

		resolvedCaps, err := captureResolverWithDefaultStubs().Resolve(
			types.CaptureExpressions{
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

		resolvedCaps, err := captureResolverWithDefaultStubs().Resolve(
			types.CaptureExpressions{capID: "my expression"},
			types.CapturedStates{},
			typestest.ToNMState(t, "name: some state"),
		)
		assert.NoError(t, err)

		assert.Equal(t, types.CapturedStates{capID: defaultStubCapturedState(t, "my expression")}, resolvedCaps)
	})
}

func testExpressionsWithPartialCache(t *testing.T) {
	t.Run("resolve with expressions and partial cache", func(t *testing.T) {
		const capID0 = "cap0"
		const capID1 = "cap1"

		capCache := types.CapturedStates{capID0: {State: typestest.ToNMState(t, "name: some captured state")}}

		resolvedCaps, err := captureResolverWithDefaultStubs().Resolve(
			types.CaptureExpressions{
				capID0: "my expression",
				capID1: "another expression",
			},
			capCache,
			typestest.ToNMState(t, "name: some state"),
		)
		assert.NoError(t, err)

		expectedCaps := capCache
		expectedCaps[capID1] = defaultStubCapturedState(t, "another expression")
		assert.Equal(t, expectedCaps, resolvedCaps)
	})
}

func testExpressionsWithOverCache(t *testing.T) {
	t.Run("resolve with cache that is not included in the expressions", func(t *testing.T) {
		const capID0 = "cap0"
		const capID1 = "cap1"

		capCache := types.CapturedStates{
			capID0: {State: typestest.ToNMState(t, "name: some captured state")},
			capID1: {State: typestest.ToNMState(t, "name: another captured state")},
		}

		resolvedCaps, err := captureResolverWithDefaultStubs().Resolve(
			types.CaptureExpressions{
				capID0: "my expression",
			},
			capCache,
			typestest.ToNMState(t, "name: some state"),
		)
		assert.NoError(t, err)

		expectedCaps := types.CapturedStates{
			capID0: {State: typestest.ToNMState(t, "name: some captured state")},
		}
		assert.Equal(t, expectedCaps, resolvedCaps)
	})
}

func testLexFailure(t *testing.T) {
	t.Run("resolve fails due to lexing", func(t *testing.T) {
		capCtrl := capture.New(lexerStub{failLex: true}, parserStub{}, resolverStub{})
		_, err := capCtrl.Resolve(
			types.CaptureExpressions{"cap0": "my expression"},
			types.CapturedStates{},
			typestest.ToNMState(t, "name: some state"),
		)
		assert.Error(t, err)
	})
}

func testParseFailure(t *testing.T) {
	t.Run("resolve fails due to parsing", func(t *testing.T) {
		capCtrl := capture.New(lexerStub{}, parserStub{failParse: true}, resolverStub{})
		_, err := capCtrl.Resolve(
			types.CaptureExpressions{"cap0": "my expression"},
			types.CapturedStates{},
			typestest.ToNMState(t, "name: some state"),
		)
		assert.Error(t, err)
	})
}

func testResolveFailure(t *testing.T) {
	t.Run("resolve fails due to resolving", func(t *testing.T) {
		capCtrl := capture.New(lexerStub{}, parserStub{}, resolverStub{failResolve: true})
		_, err := capCtrl.Resolve(
			types.CaptureExpressions{"cap0": "my expression"},
			types.CapturedStates{},
			typestest.ToNMState(t, "name: some state"),
		)
		assert.Error(t, err)
	})
}

func captureResolverWithDefaultStubs() capture.Capture {
	return capture.New(lexerStub{}, parserStub{}, resolverStub{})
}
