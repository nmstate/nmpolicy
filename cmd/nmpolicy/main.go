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

package main

import (
	"fmt"
	"log"
	"os"

	flag "github.com/spf13/pflag"
	"sigs.k8s.io/yaml"

	"github.com/nmstate/nmpolicy/nmpolicy"
	"github.com/nmstate/nmpolicy/nmpolicy/types"
)

var (
	policyPath               string
	currentStatePath         string
	capturedStatesInputPath  string
	capturedStatesOutputPath string
	desiredStatePath         string
	help                     bool
)

func main() {
	flag.StringVar(&policyPath, "policy", "policy.yaml", "input file path to NNPolicy")
	flag.StringVar(&currentStatePath, "current-state", "current-state.yaml", "input file path to current NMState")
	flag.StringVar(&capturedStatesInputPath, "captured-states-input", "captured-states-input.yaml",
		"input file path to the captured NMStates from previous run")
	flag.StringVar(&desiredStatePath, "desired-state", "desired-state.yaml", "output file for the generated NMState")
	flag.StringVar(&capturedStatesOutputPath, "captured-states-output", "captured-states-output.yaml",
		"output file path to the captured NMStates from previous run")
	flag.BoolVarP(&help, "help", "h", false, "")

	flag.Parse()

	if help {
		flag.Usage()
	}

	policy := types.PolicySpec{}
	if err := unmarshalFromPath(policyPath, &policy); err != nil {
		log.Fatalf("failed reading policy: %v", err)
	}
	capturedStatesInput := types.CachedState{}
	if err := unmarshalFromPath(capturedStatesInputPath, &capturedStatesInput); err != nil {
		log.Fatalf("failed reading captured states: %v", err)
	}

	currentState, err := os.ReadFile(currentStatePath)
	if err != nil {
		log.Fatalf("failed reading current state: %v", err)
	}

	generatedState, err := nmpolicy.GenerateState(policy, currentState, capturedStatesInput)
	if err != nil {
		log.Fatalf("failed generating desired state: %v", err)
	}

	if err := marshalToPath(desiredStatePath, generatedState.DesiredState); err != nil {
		log.Fatalf("failed writing desired state: %v", err)
	}

	if err := marshalToPath(capturedStatesOutputPath, generatedState.Cache); err != nil {
		log.Fatalf("failed writing captured states: %v", err)
	}
}

func unmarshalFromPath(path string, unmarshaled interface{}) error {
	marshaled, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed reading file to unmarshal it: %w", err)
	}

	if err := yaml.Unmarshal(marshaled, unmarshaled); err != nil {
		return fmt.Errorf("failed unmarshaling %s: %w", path, err)
	}
	return nil
}

func marshalToPath(path string, unmarshaled interface{}) error {
	marshaled, err := yaml.Marshal(unmarshaled)
	if err != nil {
		return fmt.Errorf("failed marshaling %s: %w", path, err)
	}
	if err := os.WriteFile(path, marshaled, 0); err != nil {
		return fmt.Errorf("failed writing file with marshaled data: %w", err)
	}
	return nil
}
