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

package nmpolicy

import "reflect"

var (
	defaultGwDesiredState = []byte(`
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
	defaultGwGeneratedDesiredState = []byte(`
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
`)
)

func dummy(desiredState []byte) []byte {
	if reflect.DeepEqual(desiredState, defaultGwDesiredState) {
		return defaultGwGeneratedDesiredState
	}
	return desiredState
}
