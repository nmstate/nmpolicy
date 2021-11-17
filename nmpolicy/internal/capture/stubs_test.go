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

	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"

	"github.com/nmstate/nmpolicy/nmpolicy/internal/ast"
	"github.com/nmstate/nmpolicy/nmpolicy/internal/lexer"
	"github.com/nmstate/nmpolicy/nmpolicy/types"
	"github.com/nmstate/nmpolicy/nmpolicy/types/typestest"
)

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

func (r resolverStub) Resolve(astPool map[string]ast.Node,
	state []byte, capturedStates map[string]map[string]interface{}) (map[string]types.CaptureState, error) {
	if r.failResolve {
		return nil, fmt.Errorf("resolve failed")
	}
	capsState := map[string]types.CaptureState{}
	for id, entry := range astPool {
		capsState[id] = types.CaptureState{State: []byte(fmt.Sprintf("resolver: %s", *entry.String))}
	}
	marshaledCapturedStates, err := marshalCapturedStates(capturedStates)
	if err != nil {
		return nil, err
	}
	for id, capturedState := range marshaledCapturedStates {
		capsState[id] = capturedState
	}
	return capsState, nil
}

func (r resolverStub) ResolveCaptureEntryPath(captureEntryPathAST ast.Node,
	capturedStates map[string]map[string]interface{}) (interface{}, error) {
	if r.failResolve {
		return nil, fmt.Errorf("resolve capture entry path failed")
	}
	return fmt.Sprintf("resolver: %s", *captureEntryPathAST.String), nil
}

func defaultStubCapturedState(expression string) types.CaptureState {
	return types.CaptureState{
		State: []byte(defaultStubValue(expression)),
	}
}

func defaultStubValue(expression string) string {
	return fmt.Sprintf("resolver: parser: lexer: %s", expression)
}

func marshalCapturedStates(capturedStates map[string]map[string]interface{}) (map[string]types.CaptureState, error) {
	marshaledCapturedStates := map[string]types.CaptureState{}
	for id, capturedState := range capturedStates {
		marshaledCapturedState := types.CaptureState{}
		var err error
		marshaledCapturedState.State, err = yaml.Marshal(capturedState)
		if err != nil {
			return nil, err
		}
		marshaledCapturedStates[id] = marshaledCapturedState
	}
	return marshaledCapturedStates, nil
}

func formatYAML(t *testing.T, unformatedYAML string) []byte {
	formatted, err := typestest.FormatYAML([]byte(unformatedYAML))
	assert.NoError(t, err)
	return formatted
}
