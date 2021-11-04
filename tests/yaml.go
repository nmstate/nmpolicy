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

package tests

import (
	"sigs.k8s.io/yaml"

	"github.com/nmstate/nmpolicy/nmpolicy/types"
)

func formatYAMLs(generatedState types.GeneratedState) (types.GeneratedState, error) {
	for captureID, captureState := range generatedState.Cache.Capture {
		formatedYAML, err := formatYAML(captureState.State)
		if err != nil {
			return types.GeneratedState{}, err
		}
		captureState.State = formatedYAML
		generatedState.Cache.Capture[captureID] = captureState
	}
	formatedYAML, err := formatYAML(generatedState.DesiredState)
	if err != nil {
		return types.GeneratedState{}, err
	}
	generatedState.DesiredState = formatedYAML
	return generatedState, nil
}

func formatYAML(unformatedYAML []byte) ([]byte, error) {
	unmarshaled := map[string]interface{}{}

	err := yaml.Unmarshal(unformatedYAML, &unmarshaled)
	if err != nil {
		return nil, err
	}

	marshaled, err := yaml.Marshal(unmarshaled)
	if err != nil {
		return nil, err
	}
	return marshaled, nil
}
