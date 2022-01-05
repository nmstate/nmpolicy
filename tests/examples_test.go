/*
 * Copyright 2022 NMPolicy Authors.
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

package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	assert "github.com/stretchr/testify/require"

	"sigs.k8s.io/yaml"

	"github.com/nmstate/nmpolicy/nmpolicy"
	"github.com/nmstate/nmpolicy/nmpolicy/types"
	"github.com/nmstate/nmpolicy/nmpolicy/types/typestest"
)

const (
	examplesPath = "../docs/examples"
)

type example struct {
	name           string
	policy         types.PolicySpec
	currentState   []byte
	generatedState []byte
	capturedStates map[string]types.CaptureState
}

func (e *example) readFile(fileName string) ([]byte, error) {
	return os.ReadFile(filepath.Join(examplesPath, e.name, fileName))
}

func (e *example) load() error {
	policyMarshaled, err := e.readFile("policy.md")
	if err != nil {
		return err
	}
	policyMarshaled = convertToPolicyYAML(policyMarshaled)
	if err = yaml.Unmarshal(policyMarshaled, &e.policy); err != nil {
		return fmt.Errorf("failed unmarshaling example policy: %v", err)
	}

	e.currentState, err = e.readFile("current.yaml")
	if err != nil {
		return err
	}

	e.generatedState, err = e.readFile("generated.yaml")
	if err != nil {
		return err
	}

	capturedStatesMarshaled, err := e.readFile("captured.yaml")
	if err != nil {
		return err
	}
	e.capturedStates = map[string]types.CaptureState{}
	if err := yaml.Unmarshal(capturedStatesMarshaled, &e.capturedStates); err != nil {
		return fmt.Errorf("failed unmarshaling example captured states: %v", err)
	}
	return nil
}

func (e *example) run() (types.GeneratedState, error) {
	obtained, err := nmpolicy.GenerateState(e.policy, e.currentState, types.CachedState{})
	if err != nil {
		return types.GeneratedState{}, err
	}
	return obtained, nil
}

func TestExamples(t *testing.T) {
	examplesEntries, err := os.ReadDir(examplesPath)
	assert.NoError(t, err)
	for _, examplesEntry := range examplesEntries {
		if !examplesEntry.IsDir() {
			continue
		}
		e := example{name: examplesEntry.Name()}
		t.Run(e.name, func(t *testing.T) {
			err := e.load()
			assert.NoError(t, err, "should successfully load the example")
			obtained, err := e.run()
			assert.NoError(t, err, "should successfully run the example")
			assert.YAMLEq(t, string(e.generatedState), string(obtained.DesiredState))
			expectedCapturedStates, err := formatCapturedStates(e.capturedStates)
			assert.NoError(t, err)
			obtainedCaptuerdStates, err := formatCapturedStates(obtained.Cache.Capture)
			assert.NoError(t, err)
			assert.Equal(t, resetCapturedStatesTimeStamp(expectedCapturedStates), resetCapturedStatesTimeStamp(obtainedCaptuerdStates))
		})
	}
}

func convertToPolicyYAML(policyMD []byte) []byte {
	replaced := strings.ReplaceAll(string(policyMD), "{% raw %}", "")
	replaced = strings.ReplaceAll(replaced, "{% endraw %}", "")
	return []byte(replaced)
}

func formatCapturedStates(capturedStates map[string]types.CaptureState) (map[string]types.CaptureState, error) {
	for captureID, captureState := range capturedStates {
		formatedYAML, err := typestest.FormatYAML(captureState.State)
		if err != nil {
			return nil, err
		}
		captureState.State = formatedYAML
		capturedStates[captureID] = captureState
	}
	return capturedStates, nil
}
