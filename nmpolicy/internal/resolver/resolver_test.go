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
	yaml "sigs.k8s.io/yaml"

	"github.com/nmstate/nmpolicy/nmpolicy/internal/lexer"
	"github.com/nmstate/nmpolicy/nmpolicy/internal/parser"
	"github.com/nmstate/nmpolicy/nmpolicy/internal/resolver"
	"github.com/nmstate/nmpolicy/nmpolicy/internal/types"
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
    description: "1st ethernet interface"
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
	captureExpressions     string
	captureASTPool         string
	capturedStatesCache    string
	expectedCapturedStates string
	err                    string
}

func withCaptureExpressions(t *testing.T, captureExpressions string) test {
	testToRun := test{}
	captureASTPool := types.CaptureASTPool{}
	for captureEntryName, captureEntryExpression := range typestest.ToCaptureExpressions(t, captureExpressions) {
		l := lexer.New()
		tokens, err := l.Lex(captureEntryExpression)
		assert.NoError(t, err)
		p := parser.New()
		astRoot, err := p.Parse(captureEntryExpression, tokens)
		assert.NoError(t, err)
		captureASTPool[captureEntryName] = astRoot
	}
	captureASTPoolMarshaled, err := yaml.Marshal(captureASTPool)
	assert.NoError(t, err)
	testToRun.captureASTPool = string(captureASTPoolMarshaled)
	testToRun.captureExpressions = captureExpressions
	return testToRun
}

func runTest(t *testing.T, testToRun *test) {
	captureASTPool := typestest.ToCaptureASTPool(t, testToRun.captureASTPool)
	currentState := typestest.ToNMState(t, sourceYAML)
	capturedStatesCache := typestest.ToCapturedStates(t, testToRun.capturedStatesCache)
	captureExpressions := typestest.ToCaptureExpressions(t, testToRun.captureExpressions)
	obtaintedCapturedStates, err := resolver.New().Resolve(captureExpressions, captureASTPool, currentState, capturedStatesCache)
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
		testFilterDifferentTypeOnPath(t)
		testFilterOptionalField(t)
		testFilterNonCaptureRefPathAtThirdArg(t)
		testFilterWithInvalidInputSource(t)
		testFilterWithInvalidTypeInSource(t)
		testFilterBadPath(t)

		testReplaceCurrentState(t)
		testReplaceCapturedState(t)
		testReplaceWithCaptureRef(t)
	})
}

func testFilterMapListOnSecondPathIdentity(t *testing.T) {
	t.Run("Filter map, list on second path identity", func(t *testing.T) {
		testToRun := test{
			captureASTPool: `
default-gw:
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
		runTest(t, &testToRun)
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
		runTest(t, &testToRun)
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
      description: "1st ethernet interface"
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
		runTest(t, &testToRun)
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
		runTest(t, &testToRun)
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
		runTest(t, &testToRun)
	})
}

func testFilterDifferentTypeOnPath(t *testing.T) {
	t.Run("Filter different type on path", func(t *testing.T) {
		testToRun := withCaptureExpressions(t, `
invalid-path-type: interfaces.ipv4.address=="10.244.0.1"
`)
		testToRun.err = `resolve error: eqfilter error: failed applying operation on the path: ` +
			`type missmatch: the value in the path doesn't match the value to filter. ` +
			`"[]interface {}" != "string" -> [map[ip:10.244.0.1 prefix-length:24] map[ip:169.254.1.0 prefix-length:16]] != 10.244.0.1
| interfaces.ipv4.address=="10.244.0.1"
| .......................^`

		runTest(t, &testToRun)
	})
}

func testFilterOptionalField(t *testing.T) {
	t.Run("Filter optional field", func(t *testing.T) {
		testToRun := test{
			captureASTPool: `
description-eth1:
  pos: 1
  eqfilter:
  - pos: 2
    identity: currentState
  - pos: 3
    path:
    - pos: 4
      identity: interfaces
    - pos: 5
      identity: description 
  - pos: 6
    string: 1st ethernet interface 
`,
			expectedCapturedStates: `
description-eth1: 
  state:
    interfaces:
    - name: eth1
      description: "1st ethernet interface"
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
`}
		runTest(t, &testToRun)
	})
}

func testFilterBadCaptureRef(t *testing.T) {
	t.Run("Filter list with non existing capture reference", func(t *testing.T) {
		testToRun := withCaptureExpressions(t, `
base-iface-routes: routes.running.next-hop-interface==capture
`)
		testToRun.err = `resolve error: eqfilter error: path capture ref is missing capture entry name
| routes.running.next-hop-interface==capture
| ...................................^`

		runTest(t, &testToRun)
	})
}

func testFilterCaptureRefNotFound(t *testing.T) {
	t.Run("Filter list with non existing capture reference", func(t *testing.T) {
		testToRun := withCaptureExpressions(t, `
base-iface-routes: routes.running.next-hop-interface==capture.default-gw.routes
`)
		testToRun.err = `resolve error: eqfilter error: capture entry 'default-gw' not found
| routes.running.next-hop-interface==capture.default-gw.routes
| ...................................^`

		runTest(t, &testToRun)
	})
}

func testFilterCaptureRefInvalidStateForPathMap(t *testing.T) {
	t.Run("Filter list with capture reference and invalid identity path step", func(t *testing.T) {
		testToRun := withCaptureExpressions(t, `
base-iface-routes: routes.running.next-hop-interface==capture.default-gw.routes.running.badfield.next-hop-interface
`)

		testToRun.capturedStatesCache = `
default-gw:
  state:
    routes:
       running:
       - destination: 0.0.0.0/0
         next-hop-address: 192.168.100.1
         next-hop-interface: eth1
         table-id: 254
`
		testToRun.err = "resolve error: eqfilter error: invalid path: failed walking non map state " +
			"'[map[destination:0.0.0.0/0 next-hop-address:192.168.100.1 next-hop-interface:eth1 table-id:254]]' " +
			"with path '[routes running badfield]'" + `
| routes.running.next-hop-interface==capture.default-gw.routes.running.badfield.next-hop-interface
| .....................................................................^`

		runTest(t, &testToRun)
	})
}

func testFilterCaptureRefInvalidStateForPathSlice(t *testing.T) {
	t.Run("Filter list with capture reference and invalid numeric path step", func(t *testing.T) {
		testToRun := withCaptureExpressions(t, `
base-iface-routes: routes.running.next-hop-interface==capture.default-gw.routes.1.0.next-hop-interface
`)
		testToRun.capturedStatesCache = `
default-gw:
  state: 
    routes:
      running:
       - destination: 0.0.0.0/0
         next-hop-address: 192.168.100.1
         next-hop-interface: eth1
         table-id: 254
`
		testToRun.err = "resolve error: eqfilter error: invalid path: failed walking non slice state " +
			"'map[running:[map[destination:0.0.0.0/0 next-hop-address:192.168.100.1 next-hop-interface:eth1 table-id:254]]]' " +
			"with path '[routes 1]'" + `
| routes.running.next-hop-interface==capture.default-gw.routes.1.0.next-hop-interface
| .............................................................^`

		runTest(t, &testToRun)
	})
}

func testFilterCaptureRefPathNotFoundMap(t *testing.T) {
	t.Run("Filter list with capture reference and path with not found identity step", func(t *testing.T) {
		testToRun := withCaptureExpressions(t, `
base-iface-routes: routes.running.next-hop-interface==capture.default-gw.routes.badfield
`)
		testToRun.capturedStatesCache = `
default-gw:
  state: 
    routes:
      running:
      - destination: 0.0.0.0/0
        next-hop-address: 192.168.100.1
        next-hop-interface: eth1
        table-id: 254
`
		testToRun.err = "resolve error: eqfilter error: invalid path: step 'badfield' from path '[routes badfield]' not found at map state " +
			"'map[running:[map[destination:0.0.0.0/0 next-hop-address:192.168.100.1 next-hop-interface:eth1 table-id:254]]]'" + `
| routes.running.next-hop-interface==capture.default-gw.routes.badfield
| .............................................................^`

		runTest(t, &testToRun)
	})
}

func testFilterCaptureRefPathNotFoundSlice(t *testing.T) {
	t.Run("Filter list with capture reference and path with not found numeric step", func(t *testing.T) {
		testToRun := withCaptureExpressions(t, `
base-interface-routes: routes.running.next-hop-interface==capture.default-gw.routes.running.6
`)
		testToRun.capturedStatesCache = `
default-gw:
  state:
    routes:
      running:
      - destination: 0.0.0.0/0
        next-hop-address: 192.168.100.1
        next-hop-interface: eth1
        table-id: 254
`
		testToRun.err = "resolve error: eqfilter error: invalid path: step '6' from path '[routes running 6]' not found at slice state " +
			"'[map[destination:0.0.0.0/0 next-hop-address:192.168.100.1 next-hop-interface:eth1 table-id:254]]'" + `
| routes.running.next-hop-interface==capture.default-gw.routes.running.6
| .....................................................................^`

		runTest(t, &testToRun)
	})
}

func testFilterNonCaptureRefPathAtThirdArg(t *testing.T) {
	t.Run("Filter list with path as third argument without capture reference", func(t *testing.T) {
		testToRun := withCaptureExpressions(t, `
base-iface-routes: routes.running.next-hop-interface==routes.running
`)

		testToRun.err = `resolve error: eqfilter error: not supported filtered value path. Only paths with a capture entry reference are supported
| routes.running.next-hop-interface==routes.running
| ...................................^`
		runTest(t, &testToRun)
	})
}

func testFilterWithInvalidInputSource(t *testing.T) {
	t.Run("Filter list with invalid input source", func(t *testing.T) {
		testToRun := withCaptureExpressions(t, `
base-iface-routes: invalidInputSource | routes.running.next-hop-interface=='eth1'
`)

		testToRun.err =
			`resolve error: eqfilter error: invalid path input source (Path=[Identity=invalidInputSource]), only capture reference is supported
| invalidInputSource | routes.running.next-hop-interface=='eth1'
| ^`
		runTest(t, &testToRun)
	})
}

func testFilterWithInvalidTypeInSource(t *testing.T) {
	oldSourceYaml := sourceYAML
	t.Run("Filter list with invalid input source", func(t *testing.T) {
		sourceYAML = `
routes:
   running:
`
		testToRun := withCaptureExpressions(t, `
base-iface-routes: routes.running.next-hop-interface=='eth1'
`)

		testToRun.err =
			`resolve error: eqfilter error: failed applying operation on the path: ` +
				`invalid path: invalid type <nil> for identity step 'Identity=running'
| routes.running.next-hop-interface=='eth1'
| .......^`

		runTest(t, &testToRun)
	})
	sourceYAML = oldSourceYaml
}

func testFilterBadPath(t *testing.T) {
	t.Run("Filter list with non existing path", func(t *testing.T) {
		testToRun := withCaptureExpressions(t, `
base-iface-routes: routes.badfield.next-hop-interface==capture.default-gw.routes.running.0.next-hop-interface
`)
		testToRun.capturedStatesCache = `
default-gw:
  state: 
    routes:
      running:
      - destination: 0.0.0.0/0
        next-hop-address: 1.2.3.4
        next-hop-interface: eth1
        table-id: 254
`
		testToRun.err =
			`resolve error: eqfilter error: failed applying operation on the path: invalid path: cannot find key badfield in ` +
				`map[config:[map[destination:0.0.0.0/0 next-hop-address:192.168.100.1 next-hop-interface:eth1 table-id:254] ` +
				`map[destination:1.1.1.0/24 next-hop-address:192.168.100.1 next-hop-interface:eth1 table-id:254]] ` +
				`running:[map[destination:0.0.0.0/0 next-hop-address:192.168.100.1 next-hop-interface:eth1 table-id:254] ` +
				`map[destination:1.1.1.0/24 next-hop-address:192.168.100.1 next-hop-interface:eth1 table-id:254] ` +
				`map[destination:2.2.2.0/24 next-hop-address:192.168.200.1 next-hop-interface:eth2 table-id:254]]]
| routes.badfield.next-hop-interface==capture.default-gw.routes.running.0.next-hop-interface
| .......^`

		runTest(t, &testToRun)
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
        description: "1st ethernet interface"
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
		runTest(t, &testToRun)
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
		runTest(t, &testToRun)
	})
}

func testReplaceWithCaptureRef(t *testing.T) {
	t.Run("Replace list of structs field from capture reference with capture reference value", func(t *testing.T) {
		testToRun := test{
			capturedStatesCache: `
default-gw:
  state:
    routes:
      running:
      - destination: 0.0.0.0/0
        next-hop-address: 192.168.100.1
        next-hop-interface: br1
        table-id: 254

br1-bridge:
  state:
    interfaces:
    - name: br1
      type: linux-bridge
      bridge:
        port:
        - name: eth3
`,

			captureASTPool: `
default-gw-br1-first-port:
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
    path:
    - pos: 8
      identity: capture
    - pos: 9
      identity: br1-bridge
    - pos: 10
      identity: interfaces
    - pos: 11
      number: 0
    - pos: 12
      identity: bridge
    - pos: 13
      identity: port
    - pos: 14
      number: 0
    - pos: 15
      identity: name
`,

			expectedCapturedStates: `
default-gw:
  state:
    routes:
      running:
      - destination: 0.0.0.0/0
        next-hop-address: 192.168.100.1
        next-hop-interface: br1
        table-id: 254

br1-bridge:
  state:
    interfaces:
    - name: br1
      type: linux-bridge
      bridge:
        port:
        - name: eth3

default-gw-br1-first-port:
  state:
    routes:
      running:
      - destination: 0.0.0.0/0
        next-hop-address: 192.168.100.1
        next-hop-interface: eth3
        table-id: 254
`,
		}
		runTest(t, &testToRun)
	})
}
