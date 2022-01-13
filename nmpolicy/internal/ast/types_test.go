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

package ast_test

import (
	"testing"

	assert "github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"

	"github.com/nmstate/nmpolicy/nmpolicy/internal/ast"
)

func strPtr(str string) *string {
	return &str
}

func TestYAML(t *testing.T) {
	astYAML := `
pos: 1
eqfilter: 
- pos: 2 
  path: 
  - pos: 3
    identity: currentState
- pos: 4
  path:
  - pos: 5
    identity: routes
  - pos: 6 
    identity: running
  - pos: 7 
    identity: destination
- pos: 8
  string: 0.0.0.0/0
`

	obtainedAST := &ast.Node{}
	assert.NoError(t, yaml.Unmarshal([]byte(astYAML), obtainedAST))

	expectedAST := &ast.Node{
		Meta: ast.Meta{Position: 1},
		EqFilter: &ast.TernaryOperator{
			{Meta: ast.Meta{Position: 2}, Path: &ast.VariadicOperator{
				{Meta: ast.Meta{Position: 3}, Terminal: ast.CurrentStateIdentity()},
			}},
			{Meta: ast.Meta{Position: 4}, Path: &ast.VariadicOperator{
				{Meta: ast.Meta{Position: 5}, Terminal: ast.Terminal{Identity: strPtr("routes")}},
				{Meta: ast.Meta{Position: 6}, Terminal: ast.Terminal{Identity: strPtr("running")}},
				{Meta: ast.Meta{Position: 7}, Terminal: ast.Terminal{Identity: strPtr("destination")}},
			}},
			{Meta: ast.Meta{Position: 8}, Terminal: ast.Terminal{Str: strPtr("0.0.0.0/0")}},
		},
	}
	assert.Equal(t, expectedAST, obtainedAST)
}
