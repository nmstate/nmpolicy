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

func TestCaptureEntry(t *testing.T) {
	t.Run("test CaptureEntry", func(t *testing.T) {
		testResolveCaptureEntryPathSuccess(t)

		testResolveCaptureEntryPathWithLexFailure(t)
		testResolveCaptureEntryPathWithParseFailure(t)
		testResolveCaptureEntryPathWithResolveFailure(t)
	})
}

func testResolveCaptureEntryPathSuccess(t *testing.T) {
	t.Run("ResolveCaptureEntryPath success", func(t *testing.T) {
		capturedStates := map[string]types.CaptureState{}
		captureEntryResolver, err := captureEntryResolverWithDefaultStubs(capturedStates)
		assert.NoError(t, err)
		obtainedValue, err := captureEntryResolver.ResolveCaptureEntryPath("my expression")
		assert.NoError(t, err)
		assert.Equal(t, defaultStubValue("my expression"), obtainedValue)
	})
}

func testResolveCaptureEntryPathWithLexFailure(t *testing.T) {
	t.Run("ResolveCaptureEntryPath lex failure", func(t *testing.T) {
		capturedStates := map[string]types.CaptureState{}
		captureEntryResolver, err := capture.NewCaptureEntryWithLexerParserResolver(capturedStates,
			lexerStub{failLex: true}, parserStub{}, resolverStub{})
		assert.NoError(t, err)
		_, err = captureEntryResolver.ResolveCaptureEntryPath("my expression")
		assert.Error(t, err)
	})
}

func testResolveCaptureEntryPathWithParseFailure(t *testing.T) {
	t.Run("ResolveCaptureEntryPath parser failure", func(t *testing.T) {
		capturedStates := map[string]types.CaptureState{}
		captureEntryResolver, err := capture.NewCaptureEntryWithLexerParserResolver(capturedStates,
			lexerStub{}, parserStub{failParse: true}, resolverStub{})
		assert.NoError(t, err)
		_, err = captureEntryResolver.ResolveCaptureEntryPath("my expression")
		assert.Error(t, err)
	})
}

func testResolveCaptureEntryPathWithResolveFailure(t *testing.T) {
	t.Run("ResolveCaptureEntryPath resolver failure", func(t *testing.T) {
		capturedStates := map[string]types.CaptureState{}
		captureEntryResolver, err := capture.NewCaptureEntryWithLexerParserResolver(capturedStates,
			lexerStub{}, parserStub{}, resolverStub{failResolve: true})
		assert.NoError(t, err)
		_, err = captureEntryResolver.ResolveCaptureEntryPath("my expression")
		assert.Error(t, err)
	})
}

func captureEntryResolverWithDefaultStubs(capturedStates map[string]types.CaptureState) (capture.CaptureEntry, error) {
	return capture.NewCaptureEntryWithLexerParserResolver(capturedStates, lexerStub{}, parserStub{}, resolverStub{})
}
