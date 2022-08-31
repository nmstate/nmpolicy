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

package tests

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	assert "github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"

	"github.com/nmstate/nmpolicy/nmpolicy/types"
)

func TestNmstatectl(t *testing.T) {
	err := build()
	assert.NoError(t, err)
	t.Run("Basic policy", func(t *testing.T) {
		testEmptyPolicy(t)
		testPolicyWithOnlyDesiredState(t)
		testPolicyWithCachedCaptureAndDesiredStateWithoutRef(t)
		testPolicyWithoutCache(t)
		testPolicyWithFullCache(t)
		testPolicyWithPartialCache(t)
		testGenerateUniqueTimestamps(t)

		testFailureLexer(t)
		testFailureParser(t)
		testFailureResolver(t)
	})
}

func testEmptyPolicy(t *testing.T) {
	t.Run("is empty cmd", func(t *testing.T) {
		capturedStatesOutput := capturedStates(t)
		obtainedState, err := nmpolicyctl("", "gen", "testdata/policy/empty.yaml", "-o", capturedStatesOutput)
		assert.NoError(t, err)
		assert.Empty(t, obtainedState)
		assert.Empty(t, file(t, capturedStatesOutput))
	})
}

func testPolicyWithOnlyDesiredState(t *testing.T) {
	t.Run("with only desired state", func(t *testing.T) {
		capturedStatesOutput := capturedStates(t)
		obtainedState, err := nmpolicyctl("", "gen", "testdata/policy/dummy-no-capture.yaml", "-o", capturedStatesOutput)
		assert.NoError(t, err)
		assert.YAMLEq(t, file(t, "testdata/state/dummy.yaml"), string(obtainedState))
		assert.Empty(t, file(t, capturedStatesOutput))
	})
}

func testPolicyWithCachedCaptureAndDesiredStateWithoutRef(t *testing.T) {
	t.Run("with all captures cached and desired state that has no ref", func(t *testing.T) {
		capturedStatesOutput := capturedStates(t)
		obtainedState, err := nmpolicyctl("", "gen", "testdata/policy/dummy.yaml",
			"-i", "testdata/cache/dummy.yaml", "-o", capturedStatesOutput)
		assert.NoError(t, err)
		assert.YAMLEq(t, file(t, "testdata/state/dummy.yaml"), string(obtainedState))
		assert.YAMLEq(t, file(t, "testdata/cache/dummy.yaml"), file(t, capturedStatesOutput))
	})
}

func testPolicyWithoutCache(t *testing.T) {
	t.Run("without cache", func(t *testing.T) {
		capturedStatesOutput := capturedStates(t)
		obtainedState, err := nmpolicyctl(file(t, "testdata/state/main.yaml"), "gen", "testdata/policy/linux-bridge-default-gw-no-cache.yaml",
			"-i", "testdata/cache/dummy.yaml", "-o", capturedStatesOutput)
		assert.NoError(t, err)
		assert.YAMLEq(t, file(t, "testdata/state/linux-bridge-default-gw.yaml"), string(obtainedState))
		assert.YAMLEq(t, resetTimeStampFromCache(t,
			file(t, "testdata/cache/linux-bridge-default-gw.yaml")), resetTimeStampFromCache(t, file(t, capturedStatesOutput)))
	})
}

func testPolicyWithFullCache(t *testing.T) {
	t.Run("with full cache", func(t *testing.T) {
		capturedStatesOutput := capturedStates(t)
		obtainedState, err := nmpolicyctl(file(t, "testdata/state/main.yaml"), "gen", "testdata/policy/linux-bridge-default-gw-full-cache.yaml",
			"-i", "testdata/cache/linux-bridge-default-gw.yaml", "-o", capturedStatesOutput)
		assert.NoError(t, err)
		assert.YAMLEq(t, file(t, "testdata/state/linux-bridge-default-gw.yaml"), string(obtainedState))
		assert.YAMLEq(t, resetTimeStampFromCache(t,
			file(t, "testdata/cache/base-iface-and-bridge-routes.yaml")), resetTimeStampFromCache(t, file(t, capturedStatesOutput)))
	})
}

func testPolicyWithPartialCache(t *testing.T) {
	t.Run("with partial cache", func(t *testing.T) {
		capturedStatesOutput := capturedStates(t)
		obtainedState, err := nmpolicyctl(file(t, "testdata/state/main.yaml"), "gen",
			"testdata/policy/linux-bridge-default-gw-partial-cache.yaml",
			"-i", "testdata/cache/default-gw.yaml", "-o", capturedStatesOutput)
		assert.NoError(t, err)
		assert.YAMLEq(t, file(t, "testdata/state/linux-bridge-default-gw.yaml"), string(obtainedState))
		assert.YAMLEq(t, resetTimeStampFromCache(t,
			file(t, "testdata/cache/linux-bridge-default-gw.yaml")), resetTimeStampFromCache(t, file(t, capturedStatesOutput)))
	})
}

func testGenerateUniqueTimestamps(t *testing.T) {
	t.Run("with no cache all the timestamps should be the same", func(t *testing.T) {
		beforeGenerate := time.Now()
		capturedStatesOutput := capturedStates(t)
		_, err := nmpolicyctl(file(t, "testdata/state/main.yaml"), "gen", "testdata/policy/linux-bridge-default-gw-no-cache.yaml",
			"-i", "testdata/cache/dummy.yaml", "-o", capturedStatesOutput)
		assert.NoError(t, err)

		capturedStates := map[string]types.CaptureState{}
		err = yaml.Unmarshal([]byte(file(t, capturedStatesOutput)), &capturedStates)
		assert.NoError(t, err)
		assert.NotEmpty(t, capturedStates)
		emptyTime := time.Time{}
		previousTimeStamp := emptyTime
		for captureEntryName, capturedState := range capturedStates {
			assert.Greaterf(t, capturedState.MetaInfo.TimeStamp.Sub(beforeGenerate), time.Duration(0),
				"captured state for %s should have non zero timestamp", captureEntryName)
			if previousTimeStamp != emptyTime {
				assert.Equal(t, previousTimeStamp, capturedState.MetaInfo.TimeStamp)
			}
			previousTimeStamp = capturedState.MetaInfo.TimeStamp
		}
	})
}

func testFailureLexer(t *testing.T) {
	t.Run("with lexer error", func(t *testing.T) {
		capturedStatesOutput := capturedStates(t)
		_, err := nmpolicyctl(file(t, "testdata/state/main.yaml"), "gen", "testdata/policy/bad-lexer.yaml",
			"-i", "testdata/cache/dummy.yaml", "-o", capturedStatesOutput)
		assert.EqualError(t, err,
			`'' 'Error: failed to generate state, err: failed to resolve capture expression, err: illegal rune -
| routes.running.destination==-"0.0.0.0/0"
| ............................^
': exit status 1`)
	})
}

func testFailureParser(t *testing.T) {
	t.Run("with parser error", func(t *testing.T) {
		capturedStatesOutput := capturedStates(t)
		_, err := nmpolicyctl(file(t, "testdata/state/main.yaml"), "gen", "testdata/policy/bad-parser.yaml",
			"-i", "testdata/cache/dummy.yaml", "-o", capturedStatesOutput)

		assert.EqualError(t, err,
			"'' 'Error: failed to generate state, err: failed to resolve capture expression, "+
				"err: invalid pipe: only paths can be piped in"+`
| routes.running.destination=="0.0.0.0/0" |
| ........................................^
': exit status 1`)
	})
}

func testFailureResolver(t *testing.T) {
	t.Run("with resolver error", func(t *testing.T) {
		capturedStatesOutput := capturedStates(t)
		_, err := nmpolicyctl(file(t, "testdata/state/main.yaml"), "gen", "testdata/policy/bad-resolver.yaml",
			"-i", "testdata/cache/dummy.yaml", "-o", capturedStatesOutput)

		assert.EqualError(t, err,
			"'' 'Error: failed to generate state, err: failed to resolve capture expression, "+
				"err: resolve error: eqfilter error: invalid path input source (Path=[Identity=interfaces]), only capture reference is supported"+`
| interfaces | routes.running.destination=="0.0.0.0/0"
| ^
': exit status 1`)
	})
}

func resetCapturedStatesTimeStamp(capturedStates map[string]types.CaptureState) map[string]types.CaptureState {
	for captureID, captureState := range capturedStates {
		captureState.MetaInfo.TimeStamp = time.Time{}
		capturedStates[captureID] = captureState
	}
	return capturedStates
}

func resetTimeStampFromCache(t *testing.T, marshaled string) string {
	capturedStates := map[string]types.CaptureState{}
	err := yaml.Unmarshal([]byte(marshaled), &capturedStates)
	assert.NoError(t, err)
	capturedStates = resetCapturedStatesTimeStamp(capturedStates)
	marshaledBytes, err := yaml.Marshal(capturedStates)
	assert.NoError(t, err)
	return string(marshaledBytes)
}

func nmpolicyctl(input string, args ...string) ([]byte, error) {
	nmpolicyctl_path := os.Getenv("NMPOLICYCTL")
	if nmpolicyctl_path == "" {
		nmpolicyctl_path = "../.out/nmpolicyctl"
	}
	cmd := exec.Command(nmpolicyctl_path, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	if input != "" {
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return nil, err
		}
		go func() {
			defer stdin.Close()
			_, err = io.WriteString(stdin, input)
			if err != nil {
				panic(err)
			}
		}()
	}
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("'%s' '%s': %v", stdout.String(), stderr.String(), err)
	}
	return stdout.Bytes(), nil
}

func build() error {
	cmd := exec.Command("make", "-C", "..", "build")
	var stdout, stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("'%s' '%s': %v", stdout.String(), stderr.String(), err)
	}
	return nil
}

func file(t *testing.T, filePath string) string {
	content, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	return string(content)
}

func capturedStates(t *testing.T) string {
	return filepath.Join(t.TempDir(), "captured-states.yaml")
}
