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
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	"github.com/nmstate/nmpolicy/nmpolicy"
	"github.com/nmstate/nmpolicy/nmpolicy/types"
)

// genCmd represents the gen command
var (
	currentStateFile         string
	capturedStatesInputFile  string
	capturedStatesOutputFile string
)

func genCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen policy.yaml",
		Short: "Generates NMState by policy filename",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			policySpec, err := readPolicySpec(args[0])
			if err != nil {
				return fmt.Errorf("failed reading policy spec: %v", err)
			}

			currentState, err := readCurrentState()
			if err != nil {
				return fmt.Errorf("failed reading current state: %v", err)
			}

			capturedStates, err := readCapturedStates()
			if err != nil {
				return fmt.Errorf("failed reading captured states: %v", err)
			}

			generatedState, err := nmpolicy.GenerateState(policySpec, currentState, types.CachedState{Capture: capturedStates})
			if err != nil {
				return err
			}

			if err := writeCapturedStates(generatedState.Cache.Capture); err != nil {
				return fmt.Errorf("failed writing captured states: %v", err)
			}

			fmt.Printf("%s", generatedState.DesiredState)
			return nil
		},
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	cmd.Flags().StringVarP(&currentStateFile, "current-state", "s", "",
		"input file path to current NMState. If not specified, STDIN is used.")
	cmd.Flags().StringVarP(&capturedStatesInputFile, "captured-states-input", "i", "",
		"input file path for already resolved captured states.")
	cmd.Flags().StringVarP(&capturedStatesOutputFile, "captured-states-output", "o",
		filepath.Join(homeDir, ".cache", "nmpolicy", "captured-states.yaml"),
		"output file path to the emitted captured states.")
	return cmd
}

func readStdin() ([]byte, error) {
	var buf bytes.Buffer
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanBytes)
	for {
		if !scanner.Scan() {
			if scanner.Err() != nil {
				return nil, scanner.Err()
			}
			break
		}
		buf.Write(scanner.Bytes())
	}
	return buf.Bytes(), nil
}

func readCurrentState() ([]byte, error) {
	if currentStateFile == "" {
		return readStdin()
	}
	return os.ReadFile(currentStateFile)
}

func readPolicySpec(policySpecFile string) (types.PolicySpec, error) {
	policySpecMarshaled, err := os.ReadFile(policySpecFile)
	if err != nil {
		return types.PolicySpec{}, err
	}

	policySpec := types.PolicySpec{}
	if err := yaml.Unmarshal(policySpecMarshaled, &policySpec); err != nil {
		return types.PolicySpec{}, err
	}
	return policySpec, nil
}

func readCapturedStates() (map[string]types.CaptureState, error) {
	if capturedStatesInputFile == "" {
		return map[string]types.CaptureState{}, nil
	}
	capturedStatesMarshaled, err := os.ReadFile(capturedStatesInputFile)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		return map[string]types.CaptureState{}, nil
	}

	capturedStates := map[string]types.CaptureState{}
	if err := yaml.Unmarshal(capturedStatesMarshaled, &capturedStates); err != nil {
		return nil, err
	}
	return capturedStates, nil
}

func writeCapturedStates(capturedStates map[string]types.CaptureState) error {
	marshaledCapturedStates := []byte{}
	if len(capturedStates) > 0 {
		var err error
		marshaledCapturedStates, err = yaml.Marshal(capturedStates)
		if err != nil {
			return err
		}
	}
	if err := os.MkdirAll(filepath.Dir(capturedStatesOutputFile), os.ModePerm); err != nil {
		return err
	}
	if err := os.WriteFile(capturedStatesOutputFile, marshaledCapturedStates, os.ModePerm); err != nil {
		return err
	}
	return nil
}
