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
	return nil, nil
}
