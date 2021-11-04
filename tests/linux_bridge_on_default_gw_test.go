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

package tests

import (
	"testing"

	assert "github.com/stretchr/testify/require"

	"github.com/nmstate/nmpolicy/nmpolicy"
	"github.com/nmstate/nmpolicy/nmpolicy/types"
)

func TestLinuxBridgeAtDefaultGw(t *testing.T) {
	t.Run("Linux Bridge on default gw test", func(t *testing.T) {
		testLinuxBridgeAtDefaultGw(t)
	})
}

var (
	linuxBridgeAtDefaultGwCurrentState = []byte(`
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
  config:
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
    - ip: 169.254.1.0
      prefix-length: 16
    dhcp: false
    enabled: true
- name: eth2
  type: ethernet
  state: down
  ipv4:
    address:
    - ip: 1.2.3.4
      prefix-length: 24
    dhcp: false
    enabled: false
`)
	linuxBridgeAtDefaultGwDesiredState = []byte(`
interfaces:
- name: br1
  description: Linux bridge with base interface as a port
  type: linux-bridge
  state: up
  ipv4: {{ capture.base-iface.interfaces.0.ipv4 }}
  bridge:
    options:
      stp:
        enabled: false
    port:
    - name: {{ capture.base-iface.interfaces.0.name }}
routes:
  config: {{ capture.bridge-routes-takeover.running }}
`)

	defaultGwCaptureName          = "default-gw"
	baseInterfaceRouteCaptureName = "base-iface-route"
	baseInterfaceCaptureName      = "base-iface"
	bridgeRoutes                  = "bridge-routes"
	deleteBaseIfaceRoutes         = "delete-base-iface-routes"
	bridgeRoutesTakeover          = "bridge-routes-takeover"

	linuxBridgeAtDefaultGwExpected = types.GeneratedState{
		MetaInfo: types.MetaInfo{
			Version: "0",
		},
		DesiredState: []byte(`
interfaces:
- name: br1
  description: Linux bridge with base interface as a port
  type: linux-bridge
  state: up
  ipv4:
    address:
    - ip: 10.244.0.1
      prefix-length: 24
    - ip: 169.254.1.0
      prefix-length: 16
    dhcp: false
    enabled: true
  bridge:
    options:
      stp:
        enabled: false
    port:
    - name: eth1
routes:
  config:
  - destination: 0.0.0.0/0
    next-hop-address: 192.168.100.1
    next-hop-interface: br1
    table-id: 254
  - destination: 1.1.1.0/24
    next-hop-address: 192.168.100.1
    next-hop-interface: br1
    table-id: 254
  - destination: 0.0.0.0/0
    next-hop-address: 192.168.100.1
    next-hop-interface: eth1
    state: absent
    table-id: 254
  - destination: 1.1.1.0/24
    next-hop-address: 192.168.100.1
    next-hop-interface: eth1
    state: absent
    table-id: 254
`),
		Cache: types.CachedState{
			Capture: map[string]types.CaptureState{
				defaultGwCaptureName: {State: []byte(`
routes:
  running:
  - destination: 0.0.0.0/0
    next-hop-address: 192.168.100.1
    next-hop-interface: eth1
    table-id: 254
`)},
				baseInterfaceRouteCaptureName: {State: []byte(`
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
`)},
				baseInterfaceCaptureName: {State: []byte(`
interfaces:
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
      enabled: true
`)},
				bridgeRoutes: {State: []byte(`
routes:
  running:
  - destination: 0.0.0.0/0
    next-hop-address: 192.168.100.1
    next-hop-interface: br1
    table-id: 254
  - destination: 1.1.1.0/24
    next-hop-address: 192.168.100.1
    next-hop-interface: br1
    table-id: 254
`)},
				deleteBaseIfaceRoutes: {State: []byte(`
routes:
  running:
  - destination: 0.0.0.0/0
    next-hop-address: 192.168.100.1
    next-hop-interface: eth1
    state: absent
    table-id: 254
  - destination: 1.1.1.0/24
    next-hop-address: 192.168.100.1
    next-hop-interface: eth1
    state: absent
    table-id: 254
`)},
				bridgeRoutesTakeover: {State: []byte(`
routes:
  running:
  - destination: 0.0.0.0/0
    next-hop-address: 192.168.100.1
    next-hop-interface: br1
    table-id: 254
  - destination: 1.1.1.0/24
    next-hop-address: 192.168.100.1
    next-hop-interface: br1
    table-id: 254
  - destination: 0.0.0.0/0
    next-hop-address: 192.168.100.1
    next-hop-interface: eth1
    state: absent
    table-id: 254
  - destination: 1.1.1.0/24
    next-hop-address: 192.168.100.1
    next-hop-interface: eth1
    state: absent
    table-id: 254

`)},
			},
		},
	}
)

func testLinuxBridgeAtDefaultGw(t *testing.T) {
	t.Run("with a linux bridge on top of the default gw interface without DHCP", func(t *testing.T) {
		policySpec := types.PolicySpec{
			Capture: map[string]string{
				defaultGwCaptureName:          `routes.running.destination=="0.0.0.0/0"`,
				baseInterfaceRouteCaptureName: `routes.running.next-hop-interface==capture.default-gw.routes.running.0.next-hop-interface`,
				baseInterfaceCaptureName:      `interfaces.name==capture.default-gw.routes.running.0.next-hop-interface`,
				bridgeRoutes:                  `capture.base-iface-routes | routes.running.next-hop-interface:="br1"`,
				deleteBaseIfaceRoutes:         `capture.base-iface-route | routes.running.state:="absent"`,
				bridgeRoutesTakeover:          `capture.delete-base-iface-routes.routes.running + capture.bridge-routes.routes.running`,
			},
			DesiredState: linuxBridgeAtDefaultGwDesiredState,
		}
		obtained, err := nmpolicy.GenerateState(
			policySpec,
			linuxBridgeAtDefaultGwCurrentState,
			types.CachedState{})
		assert.NoError(t, err)

		obtained = resetTimeStamp(obtained)

		obtained, err = formatYAMLs(obtained)
		assert.NoError(t, err)

		expected, err := formatYAMLs(linuxBridgeAtDefaultGwExpected)
		assert.NoError(t, err)

		assert.Equal(t, expected, obtained)
	})
}
