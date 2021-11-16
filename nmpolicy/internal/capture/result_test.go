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
	"github.com/nmstate/nmpolicy/nmpolicy/internal/resolver"
)

func TestCaptureResolverResult(t *testing.T) {
	t.Run("CaptureResolver result", func(t *testing.T) {
		testResolveCaptureEntryPathSuccess(t)

		testResolveCaptureEntryPathWithLexFailure(t)
		testResolveCaptureEntryPathWithParseFailure(t)
		testResolveCaptureEntryPathWithResolveFailure(t)
	})
}

func testResolveCaptureEntryPathSuccess(t *testing.T) {
	t.Run("ResolveCaptureEntryPath success", func(t *testing.T) {
		result := capture.NewResult(lexerStub{}, parserStub{}, resolverStub{}, resolver.Result{})
		obtainedValue, err := result.ResolveCaptureEntryPath("my expression")
		assert.NoError(t, err)
		assert.Equal(t, "resolver: parser: lexer: my expression", obtainedValue)
	})
}

func testResolveCaptureEntryPathWithLexFailure(t *testing.T) {
	t.Run("ResolveCaptureEntryPath lex failure", func(t *testing.T) {
		result := capture.NewResult(lexerStub{failLex: true}, parserStub{}, resolverStub{}, resolver.Result{})
		_, err := result.ResolveCaptureEntryPath("my expression")
		assert.Error(t, err)
	})
}

func testResolveCaptureEntryPathWithParseFailure(t *testing.T) {
	t.Run("ResolveCaptureEntryPath parser failure", func(t *testing.T) {
		result := capture.NewResult(lexerStub{}, parserStub{failParse: true}, resolverStub{}, resolver.Result{})
		_, err := result.ResolveCaptureEntryPath("my expression")
		assert.Error(t, err)
	})
}

func testResolveCaptureEntryPathWithResolveFailure(t *testing.T) {
	t.Run("ResolveCaptureEntryPath resolver failure", func(t *testing.T) {
		result := capture.NewResult(lexerStub{}, parserStub{}, resolverStub{failResolve: true}, resolver.Result{})
		_, err := result.ResolveCaptureEntryPath("my expression")
		assert.Error(t, err)
	})
}
