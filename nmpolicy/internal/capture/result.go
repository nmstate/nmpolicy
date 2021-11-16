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

	"github.com/nmstate/nmpolicy/nmpolicy/internal/resolver"
	"github.com/nmstate/nmpolicy/nmpolicy/types"
)

type Result struct {
	lexer          Lexer
	parser         Parser
	resolver       Resolver
	resolverResult resolver.Result
}

func NewResult(l Lexer, p Parser, r Resolver, resolverResult resolver.Result) Result {
	return Result{
		lexer:          l,
		parser:         p,
		resolver:       r,
		resolverResult: resolverResult,
	}
}

func (r Result) ResolveCaptureEntryPath(captureEntryPathExpression string) (interface{}, error) {
	captureEntryPathTokens, err := r.lexer.Lex(captureEntryPathExpression)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve capture entry path expression: %v", err)
	}

	captureEntryPathAST, err := r.parser.Parse(captureEntryPathTokens)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve capture entry path expression: %v", err)
	}

	resolvedCaptureEntryPath, err := r.resolver.ResolveCaptureEntryPath(captureEntryPathAST, r.resolverResult.Unmarshaled)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve capture entry path expression: %v", err)
	}

	return resolvedCaptureEntryPath, nil
}

func (r Result) CapturedStates() map[string]types.CaptureState {
	return r.resolverResult.Marshaled
}
