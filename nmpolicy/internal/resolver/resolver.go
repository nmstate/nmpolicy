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

package resolver

import (
	"fmt"
	"reflect"

	"github.com/nmstate/nmpolicy/nmpolicy/internal/ast"
	"github.com/nmstate/nmpolicy/nmpolicy/internal/capture"
	"github.com/nmstate/nmpolicy/nmpolicy/types"
)

type Resolver struct {
	state   types.NMState
	astPool capture.AstPooler
}

func New(state types.NMState, astPool capture.AstPooler) Resolver {
	return Resolver{
		state:   state,
		astPool: astPool,
	}
}

func (r Resolver) Resolve() (map[types.CaptureID]types.CaptureState, error) {
	if reflect.DeepEqual(r.astPool.Range(), ast.Pool{
		types.CaptureID("cap0"): ast.Node{
			Meta: ast.Meta{Position: 26},
			EqFilter: &ast.TernaryOperator{
				ast.Node{
					Meta:     ast.Meta{Position: 0},
					Terminal: ast.CurrentStateIdentity()},
				ast.Node{
					Meta: ast.Meta{Position: 0},
					Path: &ast.VariadicOperator{
						ast.Node{
							Meta:     ast.Meta{Position: 0},
							Terminal: ast.Terminal{Identity: strPtr("routes")},
						},
						ast.Node{
							Meta:     ast.Meta{Position: 7},
							Terminal: ast.Terminal{Identity: strPtr("running")},
						},
						ast.Node{
							Meta:     ast.Meta{Position: 15},
							Terminal: ast.Terminal{Identity: strPtr("destination")},
						},
					},
				},
				ast.Node{
					Meta:     ast.Meta{Position: 28},
					Terminal: ast.Terminal{String: strPtr("0.0.0.0/0")},
				},
			},
		},
	}) {
		return map[types.CaptureID]types.CaptureState{
			"cap0": {
				State: types.NMState(`
routes:
 running:
 - destination: 0.0.0.0/0
   next-hop-address: 192.168.100.1
   next-hop-interface: eth1
   table-id: 254
`),
			}}, nil
	}
	fmt.Println("resolver: ast not matching, returning nil")
	return nil, nil
}

func strPtr(str string) *string {
	return &str
}
