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

package state

import (
	"fmt"
	"testing"

	assert "github.com/stretchr/testify/require"
	yaml "sigs.k8s.io/yaml"
)

func TestExpanderCapturesAreMapValues(t *testing.T) {
	desiredState := `
interfaces:
- name: br1
  description: Linux bridge with base interface as a port
  type: linux-bridge
  state: up
  ipv4: "{{ capture.base-iface.interfaces[0].ipv4 }}"
  bridge:
    options:
      stp:
        enabled: false
    port:
    - name: "{{ capture.base-iface.interfaces[0].name }}"
routes:
  config: "{{ capture.bridge-routes-takeover.running }}"
`
	expectedExandedState := `interfaces:
- bridge:
    options:
      stp:
        enabled: false
    port:
    - name: eth1
  description: Linux bridge with base interface as a port
  ipv4: 1.2.3.4
  name: br1
  state: up
  type: linux-bridge
routes:
  config:
  - destination: 0.0.0.0/0
    next-hop-address: 192.168.100.1
    next-hop-interface: eth1
    table-id: 254
  - destination: 1.1.1.0/24
    next-hop-address: 192.168.100.1
    next-hop-interface: eth1
    table-id: 254`

	routes := `
  - destination: 0.0.0.0/0
    next-hop-address: 192.168.100.1
    next-hop-interface: eth1
    table-id: 254
  - destination: 1.1.1.0/24
    next-hop-address: 192.168.100.1
    next-hop-interface: eth1
    table-id: 254
`
	unmarshaledRoutes, err := unmarshalState([]byte(routes))
	assert.NoError(t, err)
	capturerStub := pathCapturerStub{failResolve: false,
		pathResults: map[string]interface{}{"capture.base-iface.interfaces[0].ipv4": "1.2.3.4", "capture.base-iface.interfaces[0].name": "eth1",
			"capture.bridge-routes-takeover.running": unmarshaledRoutes},
	}
	expandedState, err := NewExpander(capturerStub).Expand([]byte(desiredState))
	assert.NoError(t, err)
	verifyResult(t, expectedExandedState, expandedState)
}

func TestExpanderCaptureIsTopLevel(t *testing.T) {
	desiredState := `
"{{ capture.base-iface }}"
`
	expectedExandedState := `interfaces:
  - name: eth1
    type: ethernet
    state: up
    ipv4:
      address:
      - ip: 10.244.0.1
        prefix-length: 24
      - ip: 169.254.1.0
        prefix-length: 16
      dhcp: false
      enabled: true`

	unmarshaledInterfaces, err := unmarshalState([]byte(expectedExandedState))
	assert.NoError(t, err)

	capturerStub := pathCapturerStub{failResolve: false,
		pathResults: map[string]interface{}{"capture.base-iface": unmarshaledInterfaces},
	}
	expandedState, err := NewExpander(capturerStub).Expand([]byte(desiredState))
	assert.NoError(t, err)
	verifyResult(t, expectedExandedState, expandedState)
}

func TestExpanderResolveCaptureFails(t *testing.T) {
	desiredState := `
"{{ capture.enabled-iface }}"
`
	capturerStub := pathCapturerStub{failResolve: true}
	expandedState, err := NewExpander(capturerStub).Expand([]byte(desiredState))

	assert.Error(t, err)
	assert.Nil(t, expandedState)
}

func verifyResult(t *testing.T, expectedExandedState string, expandedState []byte) {
	expectedState := make(map[string]interface{})
	actualState := make(map[string]interface{})
	assert.NoError(t, yaml.Unmarshal([]byte(expectedExandedState), &expectedState))
	assert.NoError(t, yaml.Unmarshal(expandedState, &actualState))
	assert.Equal(t, expectedState, actualState)
}

type pathCapturerStub struct {
	failResolve bool
	pathResults map[string]interface{}
}

func (c pathCapturerStub) ResolveCaptureEntryPath(capturePath string) (interface{}, error) {
	if c.failResolve {
		return nil, fmt.Errorf("resolved failed")
	}

	result, found := c.pathResults[capturePath]
	if !found {
		return nil, fmt.Errorf("couldn't find capture path")
	}

	return result, nil
}
