// Copyright 2021 The NMPolicy Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package nmpolicy

import (
	"time"

	"github.com/nmstate/nmpolicy/nmpolicy/types"
)

// GenerateState creates a NMPolicy state based on the given input data:
// - NMPolicy spec.
// - NMState state, representing a given current state.
// - Cache state which includes (already resolved) named captures.
//
// GenerateState returns a generated state which includes:
// - Desired State: The NMState state which has been built by the policy.
// - Cache: Named NMState states which have been resolved by the policy.
//          Can be saved for use as cache data (passed as input).
// - Meta Info: Extended information about the generated state (e.g. the policy version).
//
// On failure, an error is returned.
func GenerateState(nmpolicy types.PolicySpec, currentState types.NMState, cache types.CachedState) (types.GeneratedState, error) {
	var desiredState types.NMState

	if nmpolicy.DesiredState != nil {
		desiredState = append(desiredState, nmpolicy.DesiredState...)
	}

	return types.GeneratedState{
		Cache:        types.NoCache(),
		DesiredState: desiredState,
		MetaInfo: types.MetaInfo{
			Version:   "0",
			TimeStamp: time.Now().UTC(),
		},
	}, nil
}
