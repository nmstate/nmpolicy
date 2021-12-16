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

package typestest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	yaml "sigs.k8s.io/yaml"

	"github.com/nmstate/nmpolicy/nmpolicy/internal/types"
)

func ToNMState(t *testing.T, marshaled string) (nmState types.NMState) {
	assert.NoError(t, yaml.Unmarshal([]byte(marshaled), &nmState))
	return nmState
}

func ToIface(t *testing.T, marshaled string) (iface interface{}) {
	assert.NoError(t, yaml.Unmarshal([]byte(marshaled), &iface))
	return iface
}

func ToCaptureExpressions(t *testing.T, marshaled string) (captureExpressions types.CaptureExpressions) {
	assert.NoError(t, yaml.Unmarshal([]byte(marshaled), &captureExpressions))
	return captureExpressions
}

func ToCapturedStates(t *testing.T, marshaled string) (capturedStates types.CapturedStates) {
	assert.NoError(t, yaml.Unmarshal([]byte(marshaled), &capturedStates))
	return capturedStates
}

func ToCaptureASTPool(t *testing.T, marshaled string) (captureASTPool types.CaptureASTPool) {
	assert.NoError(t, yaml.Unmarshal([]byte(marshaled), &captureASTPool))
	return captureASTPool
}
