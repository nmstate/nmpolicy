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

package tests

import (
	"testing"
	"time"

	assert "github.com/stretchr/testify/require"

	"github.com/nmstate/nmpolicy/nmpolicy"
	"github.com/nmstate/nmpolicy/nmpolicy/types"
)

func TestBasicPolicy(t *testing.T) {
	t.Run("Basic policy", func(t *testing.T) {
		testEmptyPolicy(t)
		testPolicyWithOnlyDesiredState(t)
		testPolicyWithCachedCapturesAndNoDesiredStateRef(t)
	})
}

func testEmptyPolicy(t *testing.T) {
	t.Run("is empty", func(t *testing.T) {
		s, err := nmpolicy.GenerateState(types.PolicySpec{}, nil, types.NoCache())

		assert.NoError(t, err)

		expectedEmptyState := types.GeneratedState{MetaInfo: types.MetaInfo{Version: "0"}}
		assert.NotEqual(t, time.Time{}, s.MetaInfo.TimeStamp)
		assert.Equal(t, expectedEmptyState, resetTimeStamp(s))
	})
}

func testPolicyWithOnlyDesiredState(t *testing.T) {
	// When a basic input with only the desired state is provided,
	// the policy just passes it as is to the output with no modifications.
	t.Run("with only desired state", func(t *testing.T) {
		stateData := []byte(`this is not a legal yaml format!`)
		policySpec := types.PolicySpec{
			DesiredState: stateData,
		}

		s, err := nmpolicy.GenerateState(policySpec, nil, types.NoCache())

		assert.NoError(t, err)
		expectedState := types.GeneratedState{
			DesiredState: stateData,
			MetaInfo:     types.MetaInfo{Version: "0"},
		}
		assert.Equal(t, expectedState, resetTimeStamp(s))
	})
}

func testPolicyWithCachedCapturesAndNoDesiredStateRef(t *testing.T) {
	t.Run("with all captures cached and desired state that has no ref", func(t *testing.T) {
		stateData := []byte(`this is not a legal yaml format!`)
		const capID0 = "cap0"
		policySpec := types.PolicySpec{
			Capture: map[types.CaptureID]types.Expression{
				capID0: "my expression",
			},
			DesiredState: stateData,
		}

		cacheState := types.CachedState{
			Capture: map[types.CaptureID]types.CaptureState{capID0: {State: []byte("some captured state")}},
		}
		s, err := nmpolicy.GenerateState(
			policySpec,
			nil,
			cacheState)

		assert.NoError(t, err)
		expectedState := types.GeneratedState{
			Cache:        cacheState,
			DesiredState: stateData,
			MetaInfo:     types.MetaInfo{Version: "0"},
		}
		assert.Equal(t, expectedState, resetTimeStamp(s))
	})
}

func resetTimeStamp(s types.GeneratedState) types.GeneratedState {
	s.MetaInfo.TimeStamp = time.Time{}
	return s
}
