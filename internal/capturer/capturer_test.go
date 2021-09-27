package capturer

import (
	"fmt"
	"testing"

	"github.com/nmstate/nmpolicy/internal/ast"
	"github.com/nmstate/nmpolicy/internal/state"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestCapturer(t *testing.T) {

	toCapturedState := func(capturedStatesYaml map[string]string) (CapturedStateByName, error) {
		capturedStates := CapturedStateByName{}
		for k, v := range capturedStatesYaml {
			capturedState := state.State{}
			err := yaml.Unmarshal([]byte(v), &capturedState)
			if err != nil {
				return nil, err
			}
			capturedStates[k] = capturedState
		}
		return capturedStates, nil
	}

	toCommands := func(commandsYaml map[string]string) (CommandByName, error) {
		commands := CommandByName{}
		for k, v := range commandsYaml {
			command, err := ast.FromYAMLString(v)
			if err != nil {
				return nil, err
			}
			commands[k] = *command
		}
		return commands, nil
	}

	var tests = []struct {
		commands      map[string]string
		currentState  string
		capturedState map[string]string
		err           string
	}{
		{
			commands: map[string]string{
				"default-gw": `
equal: 
  - path:
    - id: routes
    - id: running
    - id: destination
  - string: 0.0.0.0/0
`},

			currentState: `
routes:
  running:
  - destination: 0.0.0.0/0
    next-hop-address: 192.168.100.1
    next-hop-interface: eth1
    table-id: 254
  - destination: 1.1.1.0/24
    next-hop-address: 192.168.100.1
    next-hop-interface: eth1
    table-id: 254
interfaces:
  - name: eth1
    type: ethernet
    state: up
    ipv4:
      address:
      - ip: 10.244.0.1
        prefix-length: 24
      dhcp: false
      enabled: true
`,

			capturedState: map[string]string{
				"default-gw": `
routes:
  running:
  - destination: 0.0.0.0/0
    next-hop-address: 192.168.100.1
    next-hop-interface: eth1
    table-id: 254
`,
			},
		},
	}

	for ti, tt := range tests {
		t.Run(fmt.Sprintf("%d", ti+1), func(t *testing.T) {
			commands, err := toCommands(tt.commands)
			assert.NoError(t, err)
			cachedCapturedStates := CapturedStateByName{}
			capturer := New(commands, cachedCapturedStates)
			currentState := state.State{}
			err = yaml.Unmarshal([]byte(tt.currentState), &currentState)
			assert.NoError(t, err)
			obtainedCapturedState, err := capturer.Capture(currentState)
			if tt.err != "" {
				assert.EqualError(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				expectedCaturedState, err := toCapturedState(tt.capturedState)
				assert.NoError(t, err)
				assert.Equal(t, expectedCaturedState, obtainedCapturedState)
			}
		})
	}
}
