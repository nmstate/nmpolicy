---
title: "Library"
tag: "user-guide"
---

The NMPolicy project can be integrated as a golang library calling the public 
function [GenerateState](https://pkg.go.dev/github.com/nmstate/nmpolicy/nmpolicy#GenerateState).

## Example calling GenerateState 

The following golang code will generate a linux bridge
on top of the default gateway nic using harcoded yamls, following the
syntax at [policy spec](102-policy-syntax.html), also note that the
`desiredState` field can follow YAML or JSON syntax, in this case is 
YAML to make it more human readable:

[Run it](https://go.dev/play/p/eVH2Aa_Ma-I)
{% raw %}
```golang
package main

import (
	"fmt"

	"github.com/nmstate/nmpolicy/nmpolicy"
	"github.com/nmstate/nmpolicy/nmpolicy/types"
)

func main() {
	policySpec := types.PolicySpec{
		Capture: map[string]string{
			"default-gw": `routes.running.destination=="0.0.0.0/0"`,
		},
		DesiredState: []byte(`                                                  
interfaces:                                                                     
- name: br1                                                                     
  type: linux-bridge                                                            
  state: up                                                                     
  ipv4:                                                                         
    dhcp: true                                                                  
    enabled: true                                                               
  bridge:                                                                       
    port:                                                                       
    - name: "{{ capture.default-gw.routes.running.0.next-hop-interface }}"
`),
	}
	currentState := []byte(`                                                    
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
  mac-address: 00:00:5E:00:00:01                                                
  ipv4:                                                                         
    address:                                                                    
    - ip: 10.244.0.1                                                            
      prefix-length: 24                                                         
    - ip: 169.254.1.0                                                           
      prefix-length: 16                                                         
    dhcp: true                                                                  
    enabled: true                                                               
`)
	generatedState, err := nmpolicy.GenerateState(policySpec, currentState, types.NoCache())
	if err != nil {
		panic(err)
	}

	fmt.Println(string(generatedState.DesiredState))
}
```
{% endraw %}

It should output the following network state 
```yaml
interfaces:
- bridge:
    port:
    - name: eth1
  ipv4:
    dhcp: true
    enabled: true
  name: br1
  state: up
  type: linux-bridge
```

Also the result from the ```capture``` evaluation will be returned at the
```generateState.Cache.Capture``` [field](https://pkg.go.dev/github.com/nmstate/nmpolicy@v0.2.0/nmpolicy/types#CachedState)
```yaml
default-gw:
  metaInfo:
     time: "2021-12-15T13:45:40Z"
  state:
    routes:
      running:
      - destination: 0.0.0.0/0
        next-hop-address: 192.168.100.1
        next-hop-interface: eth1
        table-id: 254
```

## Captured states cache

The capture entries are calculated and resolved based on the state of the 
network at the point the capture entry is created.

Sometimes it is necessary to re-generate the NMState output from a NMPolicy 
but not rely on the current network state.

For example when updating the desired state, when an error at desiredState 
needs to be fixed or when generated NMState is lost and needs to be 
re-applied. Usually for those scenarios you need the old network state 
so NMPolicy is applied to the same state it was designed for.

For those scenarios the tool has an extra output that contains the results 
from capture expression evaluation that can be used as input to the tool 
instead of the normal currentState, allowing to fix desiredState fields at a 
NMPolicy. The cache contains the captured state referenced by desiredState 
or another capture, not all the captured states.

When the cached captured state is used the capture expressions evaluation from 
NMPolicy is ignored and the cache is used instead, 
so if changes are done at the capture they will be ignored.

The format will be a map with the capture entry
key the value result of the capture entry expression evaluation.

Example:

```yaml
{% include_absolute 'examples/bridge-on-default-gw-dhcp/captured.yaml' %}
```
