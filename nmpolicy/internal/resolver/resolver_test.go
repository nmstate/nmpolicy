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

package resolver_test

import (
	"testing"

	assert "github.com/stretchr/testify/require"

	"github.com/nmstate/nmpolicy/nmpolicy/internal/resolver"
	"github.com/nmstate/nmpolicy/nmpolicy/internal/types/typestest"
)

var sourceYAML string = `
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
  - destination: 2.2.2.0/24
    next-hop-address: 192.168.200.1
    next-hop-interface: eth2
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
`

type test struct {
	captureASTPool         string
	capturedStatesCache    string
	expectedCapturedStates string
	err                    string
}

func runTest(t *testing.T, testToRun test) {
	captureASTPool := typestest.ToCaptureASTPool(t, testToRun.captureASTPool)
	currentState := typestest.ToNMState(t, sourceYAML)
	capturedStatesCache := typestest.ToCapturedStates(t, testToRun.capturedStatesCache)
	obtaintedCapturedStates, err := resolver.New().Resolve(captureASTPool, currentState, capturedStatesCache)
	if testToRun.err == "" {
		assert.NoError(t, err)
		expectedCapturedState := typestest.ToCapturedStates(t, testToRun.expectedCapturedStates)
		assert.Equal(t, expectedCapturedState, obtaintedCapturedStates)
	} else {
		assert.EqualError(t, err, testToRun.err)
	}
}

func TestFilter(t *testing.T) {
	t.Run("Resolve Filter", func(t *testing.T) {
		testFilterMapListOnSecondPathIdentity(t)
		testFilterMapListOnFirstPathIdentity(t)
		testFilterList(t)
		testFilterCaptureRef(t)
		testFilterCaptureRefWithoutCapturedState(t)
		testFilterCaptureRefNotFound(t)
		testFilterBadCaptureRef(t)
		testFilterCaptureRefPathNotFoundMap(t)
		testFilterCaptureRefPathNotFoundSlice(t)
		testFilterCaptureRefInvalidStateForPathMap(t)
		testFilterCaptureRefInvalidStateForPathSlice(t)
		testFilterInvalidTypeOnPath(t)
		testFilterInvalidPath(t)
		testFilterNonCaptureRefPathAtThirdArg(t)
		testReplaceCurrentState(t)
		testReplaceCapturedState(t)
	})
}

func testFilterMapListOnSecondPathIdentity(t *testing.T) {
	t.Run("Filter map, list on second path identity", func(t *testing.T) {
		testToRun := test{
			captureASTPool: `
default-gw:
  pos: 1
  eqfilter:
  - pos: 2
    identity: currentState
  - pos: 3
    path:
    - pos: 4
      identity: routes
    - pos: 5
      identity: running
    - pos: 6
      identity: destination
  - pos: 7
    string: 0.0.0.0/0
`,

			expectedCapturedStates: `
default-gw:
  state:
    routes:
      running:
      - destination: 0.0.0.0/0
        next-hop-address: 192.168.100.1
        next-hop-interface: eth1
        table-id: 254
`,
		}
		runTest(t, testToRun)
	})
}

func testFilterMapListOnFirstPathIdentity(t *testing.T) {
	t.Run("Filter map, list on first path identity", func(t *testing.T) {
		testToRun := test{
			captureASTPool: `
up-interfaces: 
  pos: 1
  eqfilter:
  - pos: 2
    identity: currentState
  - pos: 3
    path:
    - pos: 4
      identity: interfaces
    - pos: 5
      identity: state
  - pos: 6
    string: down
`,
			expectedCapturedStates: `
up-interfaces:
  state: 
    interfaces:
      - name: eth2
        type: ethernet
        state: down
        ipv4:
          address:
          - ip: 1.2.3.4
            prefix-length: 24
          dhcp: false
          enabled: false
`,
		}
		runTest(t, testToRun)
	})
}

func testFilterList(t *testing.T) {
	t.Run("Filter list", func(t *testing.T) {
		testToRun := test{
			captureASTPool: `
specific-ipv4:
  pos: 1
  eqfilter:
  - pos: 2
    identity: currentState
  - pos: 3
    path:
    - pos: 4
      identity: interfaces
    - pos: 5
      identity: ipv4
    - pos: 6
      identity: address
    - pos: 7
      identity: ip
  - pos: 8
    string: 10.244.0.1
`,
			expectedCapturedStates: `
specific-ipv4:
  state: 
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
`,
		}
		runTest(t, testToRun)
	})
}

func testFilterCaptureRef(t *testing.T) {
	t.Run("Filter list with capture reference", func(t *testing.T) {
		testToRun := test{
			capturedStatesCache: `
default-gw:
  state: 
    routes:
      running:
      - destination: 0.0.0.0/0
        next-hop-address: 192.168.100.1
        next-hop-interface: eth1
        table-id: 254
`,
			captureASTPool: `
base-iface-routes:
  pos: 1
  eqfilter:
  - pos: 2
    identity: currentState
  - pos: 3
    path:
    - pos: 4
      identity: routes
    - pos: 5
      identity: running
    - pos: 6
      identity: next-hop-interface
  - pos: 7
    path:
    - pos: 8
      identity: capture
    - pos: 9
      identity: default-gw
    - pos: 10
      identity: routes
    - pos: 11
      identity: running
    - pos: 12
      number: 0
    - pos: 13
      identity: next-hop-interface
`,

			expectedCapturedStates: `
default-gw:
  state: 
    routes:
      running:
      - destination: 0.0.0.0/0
        next-hop-address: 192.168.100.1
        next-hop-interface: eth1
        table-id: 254
base-iface-routes:
  state:
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
`,
		}
		runTest(t, testToRun)
	})
}

func testFilterCaptureRefWithoutCapturedState(t *testing.T) {
	t.Run("Filter list with capture reference", func(t *testing.T) {
		testToRun := test{
			captureASTPool: `
default-gw:
  pos: 1
  eqfilter:
  - pos: 2
    identity: currentState
  - pos: 3
    path:
    - pos: 4
      identity: routes
    - pos: 5
      identity: running
    - pos: 6
      identity: destination
  - pos: 7
    string: 0.0.0.0/0
base-iface-routes:
  pos: 1
  eqfilter:
  - pos: 2
    identity: currentState
  - pos: 3
    path:
    - pos: 4
      identity: routes
    - pos: 5
      identity: running
    - pos: 6
      identity: next-hop-interface
  - pos: 7
    path:
    - pos: 8
      identity: capture
    - pos: 9
      identity: default-gw
    - pos: 10
      identity: routes
    - pos: 11
      identity: running
    - pos: 12
      number: 0
    - pos: 13
      identity: next-hop-interface
`,
			expectedCapturedStates: `
default-gw:
  state: 
    routes:
      running:
      - destination: 0.0.0.0/0
        next-hop-address: 192.168.100.1
        next-hop-interface: eth1
        table-id: 254
base-iface-routes:
  state: 
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
`,
		}
		runTest(t, testToRun)
	})
}

func testFilterInvalidTypeOnPath(t *testing.T) {
	t.Run("Filter invalid type on path", func(t *testing.T) {
		testToRun := test{
			captureASTPool: `
invalid-path-type:
  pos: 1
  eqfilter:
  - pos: 2
    identity: currentState
  - pos: 3
    path:
    - pos: 4
      identity: interfaces
    - pos: 5
      identity: ipv4
    - pos: 6
      identity: address
  - pos: 7
    string: 10.244.0.1
`,
			err: "resolve error: eqfilter error: failed applying operation on the path: " +
				"error comparing the expected and obtained values : " +
				"the value [map[ip:10.244.0.1 prefix-length:24] map[ip:169.254.1.0 prefix-length:16]] of type []interface {} " +
				"not supported,curretly only string values are supported",
		}
		runTest(t, testToRun)
	})
}

func testFilterInvalidPath(t *testing.T) {
	t.Run("Filter invalid path", func(t *testing.T) {
		testToRun := test{
			captureASTPool: `
invalid-path-identity:
  pos: 1
  eqfilter:
  - pos: 2
    identity: currentState
  - pos: 3
    path:
    - pos: 4
      identity: interfaces
    - pos: 5
      identity: name-invalid-path
  - pos: 6
    string: eth0
`,
			err: "resolve error: eqfilter error: failed applying operation on the path: cannot find key name-invalid-path in " +
				"map[ipv4:map[address:[map[ip:10.244.0.1 prefix-length:24] " +
				"map[ip:169.254.1.0 prefix-length:16]] dhcp:false enabled:true] name:eth1 state:up type:ethernet]",
		}
		runTest(t, testToRun)
	})
}

func testFilterBadCaptureRef(t *testing.T) {
	t.Run("Filter list with non existing capture reference", func(t *testing.T) {
		testToRun := test{
			captureASTPool: `
base-iface-routes:
  pos: 1
  eqfilter:
  - pos: 2
    identity: currentState
  - pos: 3
    path:
    - pos: 4
      identity: routes
    - pos: 5
      identity: running
    - pos: 6
      identity: next-hop-interface
  - pos: 7
    path:
    - pos: 8
      identity: capture
`,
			err: "resolve error: eqfilter error: path capture ref is missing capture entry name",
		}
		runTest(t, testToRun)
	})
}

func testFilterCaptureRefNotFound(t *testing.T) {
	t.Run("Filter list with non existing capture reference", func(t *testing.T) {
		testToRun := test{
			captureASTPool: `
base-iface-routes:
  pos: 1
  eqfilter:
  - pos: 2
    identity: currentState
  - pos: 3
    path:
    - pos: 4
      identity: routes
    - pos: 5
      identity: running
    - pos: 6
      identity: next-hop-interface
  - pos: 7
    path:
    - pos: 8
      identity: capture
    - pos: 9
      identity: default-gw
    - pos: 11
      identity: routes
`,
			err: "resolve error: eqfilter error: capture entry 'default-gw' not found",
		}
		runTest(t, testToRun)
	})
}

func testFilterCaptureRefInvalidStateForPathMap(t *testing.T) {
	t.Run("Filter list with capture reference and invalid identity path step", func(t *testing.T) {
		testToRun := test{
			capturedStatesCache: `
default-gw:
  state:
    routes:
       running:
       - destination: 0.0.0.0/0
         next-hop-address: 192.168.100.1
         next-hop-interface: eth1
         table-id: 254
`,
			captureASTPool: `
base-iface-routes:
  pos: 1
  eqfilter:
  - pos: 2
    identity: currentState
  - pos: 3
    path:
    - pos: 4
      identity: routes
    - pos: 5
      identity: running
    - pos: 6
      identity: next-hop-interface
  - pos: 7
    path:
    - pos: 8
      identity: capture
    - pos: 9
      identity: default-gw
    - pos: 11
      identity: routes
    - pos: 12
      identity: running
    - pos: 13
      identity: badfield
    - pos: 14
      identity: next-hop-interface
`,
			err: "resolve error: eqfilter error: failed walking non map state " +
				"'[map[destination:0.0.0.0/0 next-hop-address:192.168.100.1 next-hop-interface:eth1 table-id:254]]' " +
				"with path '[routes running badfield]'",
		}
		runTest(t, testToRun)
	})
}

func testFilterCaptureRefInvalidStateForPathSlice(t *testing.T) {
	t.Run("Filter list with capture reference and invalid numeric path step", func(t *testing.T) {
		testToRun := test{
			capturedStatesCache: `
default-gw:
  state: 
    routes:
      running:
       - destination: 0.0.0.0/0
         next-hop-address: 192.168.100.1
         next-hop-interface: eth1
         table-id: 254
`,
			captureASTPool: `
base-iface-routes:
  pos: 1
  eqfilter:
  - pos: 2
    identity: currentState
  - pos: 3
    path:
    - pos: 4
      identity: routes
    - pos: 5
      identity: running
    - pos: 6
      identity: next-hop-interface
  - pos: 7
    path:
    - pos: 8
      identity: capture
    - pos: 9
      identity: default-gw
    - pos: 11
      identity: routes
    - pos: 12
      number: 1
    - pos: 13
      number: 0
    - pos: 14
      identity: next-hop-interface
`,
			err: "resolve error: eqfilter error: failed walking non slice state " +
				"'map[running:[map[destination:0.0.0.0/0 next-hop-address:192.168.100.1 next-hop-interface:eth1 table-id:254]]]' " +
				"with path '[routes 1]'",
		}
		runTest(t, testToRun)
	})
}

func testFilterCaptureRefPathNotFoundMap(t *testing.T) {
	t.Run("Filter list with capture reference and path with not found identity step", func(t *testing.T) {
		testToRun := test{
			capturedStatesCache: `
default-gw:
  state: 
    routes:
      running:
      - destination: 0.0.0.0/0
        next-hop-address: 192.168.100.1
        next-hop-interface: eth1
        table-id: 254
`,
			captureASTPool: `
base-iface-routes:
  pos: 1
  eqfilter:
  - pos: 2
    identity: currentState
  - pos: 3
    path:
    - pos: 4
      identity: routes
    - pos: 5
      identity: running
    - pos: 6
      identity: next-hop-interface
  - pos: 7
    path:
    - pos: 8
      identity: capture
    - pos: 9
      identity: default-gw
    - pos: 11
      identity: routes
    - pos: 12
      identity: badfield
`,
			err: "resolve error: eqfilter error: step 'badfield' from path '[routes badfield]' not found at map state " +
				"'map[running:[map[destination:0.0.0.0/0 next-hop-address:192.168.100.1 next-hop-interface:eth1 table-id:254]]]'",
		}
		runTest(t, testToRun)
	})
}

func testFilterCaptureRefPathNotFoundSlice(t *testing.T) {
	t.Run("Filter list with capture reference and path with not found numeric step", func(t *testing.T) {
		testToRun := test{
			capturedStatesCache: `
default-gw:
  state:
    routes:
      running:
      - destination: 0.0.0.0/0
        next-hop-address: 192.168.100.1
        next-hop-interface: eth1
        table-id: 254
`,
			captureASTPool: `
base-iface-routes:
  pos: 1
  eqfilter:
  - pos: 2
    identity: currentState
  - pos: 3
    path:
    - pos: 4
      identity: routes
    - pos: 5
      identity: running
    - pos: 6
      identity: next-hop-interface
  - pos: 7
    path:
    - pos: 8
      identity: capture
    - pos: 9
      identity: default-gw
    - pos: 11
      identity: routes
    - pos: 12
      identity: running
    - pos: 13
      number: 6
`,
			err: "resolve error: eqfilter error: step '6' from path '[routes running 6]' not found at slice state " +
				"'[map[destination:0.0.0.0/0 next-hop-address:192.168.100.1 next-hop-interface:eth1 table-id:254]]'",
		}
		runTest(t, testToRun)
	})
}

func testFilterNonCaptureRefPathAtThirdArg(t *testing.T) {
	t.Run("Filter list with path as third argument without capture reference", func(t *testing.T) {
		testToRun := test{
			captureASTPool: `
base-iface-routes:
  pos: 1
  eqfilter:
  - pos: 2
    identity: currentState
  - pos: 3
    path:
    - pos: 4
      identity: routes
    - pos: 5
      identity: running
    - pos: 6
      identity: next-hop-interface
  - pos: 7
    path:
    - pos: 8
      identity: routes
    - pos: 9
      identity: running
`,
			err: "resolve error: eqfilter error: not supported filtered value path. Only paths with a capture entry reference are supported",
		}
		runTest(t, testToRun)
	})
}

func testReplaceCurrentState(t *testing.T) {
	t.Run("Replace list of structs field from currentState with string value", func(t *testing.T) {
		testToRun := test{
			captureASTPool: `
bridge-routes:
  pos: 1
  replace:
  - pos: 2
    identity: currentState
  - pos: 3
    path:
    - pos: 4
      identity: routes
    - pos: 5
      identity: running
    - pos: 6
      identity: next-hop-interface
  - pos: 7
    string: br1
`,

			expectedCapturedStates: `

bridge-routes:
  state:
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
      - destination: 2.2.2.0/24
        next-hop-address: 192.168.200.1
        next-hop-interface: br1
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
`,
		}
		runTest(t, testToRun)
	})
}

func testReplaceCapturedState(t *testing.T) {
	t.Run("Replace list of structs field from capture reference with string value", func(t *testing.T) {
		testToRun := test{
			capturedStatesCache: `
default-gw:
  state:
    routes:
      running:
      - destination: 0.0.0.0/0
        next-hop-address: 192.168.100.1
        next-hop-interface: eth1
        table-id: 254
`,

			captureASTPool: `
bridge-routes:
  pos: 1
  replace:
  - pos: 2
    path: 
    - pos: 3
      identity: capture
    - pos: 4
      identity: default-gw
  - pos: 3
    path:
    - pos: 4
      identity: routes
    - pos: 5
      identity: running
    - pos: 6
      identity: next-hop-interface
  - pos: 7
    string: br1
`,

			expectedCapturedStates: `
default-gw:
  state:
    routes:
      running:
      - destination: 0.0.0.0/0
        next-hop-address: 192.168.100.1
        next-hop-interface: eth1
        table-id: 254
bridge-routes:
  state:
    routes:
      running:
      - destination: 0.0.0.0/0
        next-hop-address: 192.168.100.1
        next-hop-interface: br1
        table-id: 254
`,
		}
		runTest(t, testToRun)
	})
}
