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

package resolver

import (
	"reflect"

	"github.com/nmstate/nmpolicy/nmpolicy/internal/ast"
	"github.com/nmstate/nmpolicy/nmpolicy/internal/parser"
	"github.com/nmstate/nmpolicy/nmpolicy/types"
)

var (
	BaseIfaceRoutesCaptureState = types.CaptureState{
		State: []byte(`
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
`),
	}
	BaseIfaceCaptureState = types.CaptureState{
		State: []byte(`
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
`),
	}
	BridgeRoutesCaptureState = types.CaptureState{
		State: []byte(`
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
`),
	}

	DeleteBaseIfaceRoutesCaptureState = types.CaptureState{
		State: []byte(`
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
`),
	}

	BridgeRoutesTakeoverCaptureState = types.CaptureState{
		State: []byte(`
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
    table-id: 254
    state: absent
  - destination: 1.1.1.0/24
    next-hop-address: 192.168.100.1
    next-hop-interface: eth1
    table-id: 254
    state: absent
`),
	}
)

func dummy(node ast.Node) *types.CaptureState {
	if reflect.DeepEqual(node, parser.BaseIfaceRoutesAST) {
		return &BaseIfaceRoutesCaptureState
	}
	if reflect.DeepEqual(node, parser.BaseIfaceAST) {
		return &BaseIfaceCaptureState
	}
	if reflect.DeepEqual(node, parser.BridgeRoutesAST) {
		return &BridgeRoutesCaptureState
	}
	if reflect.DeepEqual(node, parser.DeleteBaseIfaceRoutesAST) {
		return &DeleteBaseIfaceRoutesCaptureState
	}
	if reflect.DeepEqual(node, parser.BridgeRoutesTakeoverAST) {
		return &BridgeRoutesTakeoverCaptureState
	}
	return nil
}
