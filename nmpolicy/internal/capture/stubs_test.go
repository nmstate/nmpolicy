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

	"sigs.k8s.io/yaml"

	"github.com/nmstate/nmpolicy/nmpolicy/internal/ast"
	"github.com/nmstate/nmpolicy/nmpolicy/internal/lexer"
	"github.com/nmstate/nmpolicy/nmpolicy/internal/types"
	"github.com/nmstate/nmpolicy/nmpolicy/internal/types/typestest"
)

type lexerStub struct {
	failLex bool
}

func (l lexerStub) Lex(expression string) ([]lexer.Token, error) {
	if l.failLex {
		return nil, fmt.Errorf("lex failed")
	}
	literal := fmt.Sprintf(`{"lexer": %q}`, expression)
	return []lexer.Token{{Literal: literal}}, nil
}

type parserStub struct {
	failParse bool
}

func (p parserStub) Parse(expression string, tokens []lexer.Token) (ast.Node, error) {
	if p.failParse {
		return ast.Node{}, fmt.Errorf("parse failed")
	}
	literal := fmt.Sprintf(`{"parser": %s}`, tokens[0].Literal)
	return ast.Node{Terminal: ast.Terminal{String: &literal}}, nil
}

type resolverStub struct {
	failResolve bool
}

func (r resolverStub) Resolve(captureASTPool types.CaptureASTPool,
	state types.NMState, capturedStates types.CapturedStates) (types.CapturedStates, error) {
	if r.failResolve {
		return nil, fmt.Errorf("resolve stub failed")
	}
	capsState := types.CapturedStates{}
	for id, entry := range captureASTPool {
		state := types.NMState{}
		marshaled := fmt.Sprintf(`{"resolver": %s}`, *entry.String)
		if err := yaml.Unmarshal([]byte(marshaled), &state); err != nil {
			return nil, fmt.Errorf("resolve stub failed: unmarshaling `%s`: %v", marshaled, err)
		}
		capsState[id] = types.CapturedState{State: state}
	}
	for id, capturedState := range capturedStates {
		capsState[id] = capturedState
	}
	return capsState, nil
}

func (r resolverStub) ResolveCaptureEntryPath(captureEntryPathAST ast.Node,
	capturedStates types.CapturedStates) (interface{}, error) {
	if r.failResolve {
		return nil, fmt.Errorf("resolve capture entry path failed")
	}
	return fmt.Sprintf(`{"resolver": %s}`, *captureEntryPathAST.String), nil
}

func defaultStubCapturedState(t *testing.T, expression string) types.CapturedState {
	return types.CapturedState{
		State: typestest.ToNMState(t, defaultStubValue(expression)),
	}
}

func defaultStubValue(expression string) string {
	return fmt.Sprintf(`{"resolver": {"parser": {"lexer": %q}}}`, expression)
}
