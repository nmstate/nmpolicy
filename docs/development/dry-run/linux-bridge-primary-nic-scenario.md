# NMPolicy stages with bridge at primary NIC scenario

The NMPolicy processing works like a pipeline where the chain like this:

1. Calculate captured state
   1. Lexer: convert capture expression to tokens
   2. Parser: convert capture tokens to AST
   3. Emitter: convert AST + current state into captured state 
2. Calculate desired state
   1. Lexer: convert place holder expression into tokens
   2. Parser: convert placeholder tokens into AST
   3. Emitter: Convert AST + captured state into idented yaml at the placeholder

To illustrate that this document uses the scenario where a linux bridge is 
created in top of the default gw NIC interface and represent then 
describe how the data is converted from stage to the next one

## Bridge on defaul Scenario NMPolicy, Captured State and Emitted Desired State

**Current State**
```yaml
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
```

**NMPolicy**
```yaml
capture:
  default-gw: routes.running.destination=="0.0.0.0/0"
  base-iface-routes: routes.running.next-hop-interface==capture.default-gw.routes.running[0].next-hop-interface
  base-iface: interfaces.name==capture.default-gw.routes.running[0].next-hop-interface
  bridge-routes: capture.base-iface-routes | routes.running.next-hop-interface:="br1"
  delete-base-iface-routes: capture.base-iface-route | routes.running.state:="absent"
  bridge-routes-takeover: capture.delete-base-iface-routes.routes.running + capture.bridge-routes.routes.running
desiredState:
  interfaces:
  - name: br1
    description: Linux bridge with base interface  as a port
    type: linux-bridge
    state: up
    ipv4: {{ capture.base-iface.interfaces[0].ipv4 }}
    bridge:
      options:
        stp:
          enabled: false
      port:
      - name: {{ capture.base-iface.interfaces[0].name }}
  routes:  
    config: {{ capture.bridge-routes-takeover.running }}

}}
```

**Captured State**
```yaml
default-gw:
  routes:
    running:
    - destination: 0.0.0.0/0
      next-hop-address: 192.168.100.1
      next-hop-interface: eth1
      table-id: 254
base-iface-routes:
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
base-iface:
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
bridge-routes:
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
bridge-routes-takeover:
  routes:
    running:
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
    - destination: 0.0.0.0/0
      next-hop-address: 192.168.100.1
      next-hop-interface: br1
      table-id: 254
    - destination: 1.1.1.0/24
      next-hop-address: 192.168.100.1
      next-hop-interface: br1
      table-id: 254
```

**Desied State**

```yaml
routes:  
  config:
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
  - destination: 0.0.0.0/0
    next-hop-address: 192.168.100.1
    next-hop-interface: br1
    table-id: 254
  - destination: 1.1.1.0/24
    next-hop-address: 192.168.100.1
    next-hop-interface: br1
    table-id: 254
interfaces:
- name: br1
  description: Linux bridge with eth1 as a port
  type: linux-bridge
  state: up
  ipv4:
    address:
    - ip: 10.244.0.1
      prefix-length: 24
    dhcp: false
    enabled: true
  bridge:
    options:
      stp:
        enabled: false
    port:
    - name: eth1
```

### Capture pipeline

Following is a golang snippet to convert the capture expressions into the
captured state:
```golang

func main() {                                                                   
    currentStateYAML := `                                                       
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
`                                                                               
    currentState := state.State{}                                               
                                                                                
    err := yaml.Unmarshal([]byte(currentStateYAML), currentState)               
    if err != nil {                                                             
        panic(err)                                                              
    }                                                                           
                                                                                
    capture := map[string]string{                                               
        "default-gw":               `routes.running.destination=="0.0.0.0/0"`,  
        "base-iface-routes":        `routes.running.next-hop-interface==capture.default-gw.routes.running.0.next-hop-interface`,                                                     
        "base-iface":               `interfaces.name==capture.default-gw.routes.running.0.next-hop-interface`,
        "bridge-routes":            `capture.base-iface-routes | routes.running.next-hop-interface="br1"`,
        "delete-base-iface-routes": `capture.base-iface-route | routes.running.state="absent"`,
        "bridge-routes-takeover":   `capture.delete-base-iface-routes.routes.running + capture.bridge-routes.routes.running`,
    }                                                                           
                                                                                
    commandByName := map[string]ast.Command{}                                   
    for ek, ev := range capture {                                               
        source := source.New(ev)                                                
        parser := parser.New(*source, lexer.NewLexer(source))                   
        command, err := parser.Parse()                                          
        if err != nil {                                                         
            panic(err)                                                          
        }                                                                       
                                                                                
        commandByName[ek] = *command                                            
    }                                                                           
                                                                                
    cache := capturer.CapturedStateByName{}                                     
    c := capturer.New(commandByName, cache)                                     
    capturedStateByName, err := c.Capture(currentState)                         
    if err != nil {                                                             
        panic(err)                                                              
    }

    fmt.Println(capturedStateByName)                                         
}
```

#### bridge-routes-takeover: capture.delete-base-iface-routes.routes.running + capture.bridge-routes.routes.running

**Tokens**
```yaml
- pos: 0
  type: IDENTITY(1)
  literal: capture
- pos: 7
  type: PATH(4)
  literal: PATH(4)
- pos: 8
  type: IDENTITY(1)
  literal: delete-base-iface-routes
- pos: 32
  type: PATH(4)
  literal: PATH(4)
- pos: 33
  type: IDENTITY(1)
  literal: routes
- pos: 39
  type: PATH(4)
  literal: PATH(4)
- pos: 40
  type: IDENTITY(1)
  literal: running
- pos: 48
  type: MERGE(11)
  literal: MERGE(11)
- pos: 50
  type: IDENTITY(1)
  literal: capture
- pos: 57
  type: PATH(4)
  literal: PATH(4)
- pos: 58
  type: IDENTITY(1)
  literal: bridge-routes
- pos: 71
  type: PATH(4)
  literal: PATH(4)
- pos: 72
  type: IDENTITY(1)
  literal: routes
- pos: 78
  type: PATH(4)
  literal: PATH(4)
- pos: 79
  type: IDENTITY(1)
  literal: running
- pos: 85
  type: EOF(0)
  literal: ""

```

**AST**
```yaml
pos: 0
merge:
  - pos: 0
    path:
      - pos: 0
        id: capture
      - pos: 8
        id: delete-base-iface-routes
      - pos: 33
        id: routes
      - pos: 40
        id: running
  - pos: 50
    path:
      - pos: 50
        id: capture
      - pos: 58
        id: bridge-routes
      - pos: 72
        id: routes
      - pos: 79
        id: running

```
#### default-gw: routes.running.destination=="0.0.0.0/0"

**Tokens**
```yaml
- pos: 0
  type: IDENTITY(1)
  literal: routes
- pos: 6
  type: PATH(4)
  literal: PATH(4)
- pos: 7
  type: IDENTITY(1)
  literal: running
- pos: 14
  type: PATH(4)
  literal: PATH(4)
- pos: 15
  type: IDENTITY(1)
  literal: destination
- pos: 26
  type: EQUALITY(10)
  literal: EQUALITY(10)
- pos: 28
  type: STRING(3)
  literal: 0.0.0.0/0
- pos: 38
  type: EOF(0)
  literal: ""

```

**AST**

![your-UML-diagram-name](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/qinqon/nmpolicy/main-scenario-dry-run/docs/development/dry-run/ast1.uml)

```golang
resolveCommand(currentState, command)
|- filterStateByEquality(currentState, arg{path: path{{id: "route"}, {id: "running"}, {id: "destination"}}, arg{string: "0.0.0.0/0"})
```

#### base-iface-routes: routes.running.next-hop-interface==capture.default-gw.routes.running.0.next-hop-interface

**Tokens**
```yaml
- pos: 0
  type: IDENTITY(1)
  literal: routes
- pos: 6
  type: PATH(4)
  literal: PATH(4)
- pos: 7
  type: IDENTITY(1)
  literal: running
- pos: 14
  type: PATH(4)
  literal: PATH(4)
- pos: 15
  type: IDENTITY(1)
  literal: next-hop-interface
- pos: 33
  type: EQUALITY(10)
  literal: EQUALITY(10)
- pos: 35
  type: IDENTITY(1)
  literal: capture
- pos: 42
  type: PATH(4)
  literal: PATH(4)
- pos: 43
  type: IDENTITY(1)
  literal: default-gw
- pos: 53
  type: PATH(4)
  literal: PATH(4)
- pos: 54
  type: IDENTITY(1)
  literal: routes
- pos: 60
  type: PATH(4)
  literal: PATH(4)
- pos: 61
  type: IDENTITY(1)
  literal: running
- pos: 68
  type: PATH(4)
  literal: PATH(4)
- pos: 69
  type: NUMBER(2)
  literal: "0"
- pos: 70
  type: PATH(4)
  literal: PATH(4)
- pos: 71
  type: IDENTITY(1)
  literal: next-hop-interface
- pos: 88
  type: EOF(0)
  literal: ""

```

![your-UML-diagram-name](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/qinqon/nmpolicy/main-scenario-dry-run/docs/development/dry-run/ast2.uml)

```golang
capture(currentState, "base-iface-routes", eqcmd)
|- resolveCommand(currentState, eqcmd)
|-- filterStateByEquality(currentState, arg{path{"route","running", "destination"}}, arg{path:{"capture", "default-gw", "routes", "running", 0, "next-hop-interface"}})
|--- walkPath(currentState, arg{path:{"capture", "default-gw", "routes", "running", 0, "next-hop-interface"}})
|---- capture("default-gw", defaultgwcmd)
```

#### base-iface: interfaces.name==capture.default-gw.routes.running.0.next-hop-interface

**Tokens**
```yaml
- pos: 0
  type: IDENTITY(1)
  literal: interfaces
- pos: 10
  type: PATH(4)
  literal: PATH(4)
- pos: 11
  type: IDENTITY(1)
  literal: name
- pos: 15
  type: EQUALITY(10)
  literal: EQUALITY(10)
- pos: 17
  type: IDENTITY(1)
  literal: capture
- pos: 24
  type: PATH(4)
  literal: PATH(4)
- pos: 25
  type: IDENTITY(1)
  literal: default-gw
- pos: 35
  type: PATH(4)
  literal: PATH(4)
- pos: 36
  type: IDENTITY(1)
  literal: routes
- pos: 42
  type: PATH(4)
  literal: PATH(4)
- pos: 43
  type: IDENTITY(1)
  literal: running
- pos: 50
  type: PATH(4)
  literal: PATH(4)
- pos: 51
  type: NUMBER(2)
  literal: "0"
- pos: 52
  type: PATH(4)
  literal: PATH(4)
- pos: 53
  type: IDENTITY(1)
  literal: next-hop-interface
- pos: 70
  type: EOF(0)
  literal: ""

```

**AST**
```yaml
pos: 0
equal:
  - pos: 0
    path:
      - pos: 0
        id: interfaces
      - pos: 11
        id: name
  - pos: 17
    path:
      - pos: 17
        id: capture
      - pos: 25
        id: default-gw
      - pos: 36
        id: routes
      - pos: 43
        id: running
      - pos: 51
        idx: 0
      - pos: 53
        id: next-hop-interface

```
#### bridge-routes: capture.base-iface-routes | routes.running.next-hop-interface="br1"

**Tokens**
```yaml
- pos: 0
  type: IDENTITY(1)
  literal: capture
- pos: 7
  type: PATH(4)
  literal: PATH(4)
- pos: 8
  type: IDENTITY(1)
  literal: base-iface-routes
- pos: 26
  type: PIPE(5)
  literal: PIPE(5)
- pos: 28
  type: IDENTITY(1)
  literal: routes
- pos: 34
  type: PATH(4)
  literal: PATH(4)
- pos: 35
  type: IDENTITY(1)
  literal: running
- pos: 42
  type: PATH(4)
  literal: PATH(4)
- pos: 43
  type: IDENTITY(1)
  literal: next-hop-interface
- pos: 61
  type: ASSIGN(9)
  literal: ASSIGN(9)
- pos: 62
  type: STRING(3)
  literal: br1
- pos: 66
  type: EOF(0)
  literal: ""
- pos: 66
  type: EOF(0)
  literal: ""

```

**AST**
```yaml
pos: 0
path:
  - pos: 0
    id: capture
  - pos: 8
    id: base-iface-routes
pipe:
    pos: 28
    assign:
      - pos: 28
        path:
          - pos: 28
            id: routes
          - pos: 35
            id: running
          - pos: 43
            id: next-hop-interface
      - pos: 62
        string: br1

```
#### delete-base-iface-routes: capture.base-iface-route | routes.running.state="absent"

**Tokens**
```yaml
- pos: 0
  type: IDENTITY(1)
  literal: capture
- pos: 7
  type: PATH(4)
  literal: PATH(4)
- pos: 8
  type: IDENTITY(1)
  literal: base-iface-route
- pos: 25
  type: PIPE(5)
  literal: PIPE(5)
- pos: 27
  type: IDENTITY(1)
  literal: routes
- pos: 33
  type: PATH(4)
  literal: PATH(4)
- pos: 34
  type: IDENTITY(1)
  literal: running
- pos: 41
  type: PATH(4)
  literal: PATH(4)
- pos: 42
  type: IDENTITY(1)
  literal: state
- pos: 47
  type: ASSIGN(9)
  literal: ASSIGN(9)
- pos: 48
  type: STRING(3)
  literal: absent
- pos: 55
  type: EOF(0)
  literal: ""
- pos: 55
  type: EOF(0)
  literal: ""

```

**AST**
```yaml
pos: 0
path:
  - pos: 0
    id: capture
  - pos: 8
    id: base-iface-route
pipe:
    pos: 27
    assign:
      - pos: 27
        path:
          - pos: 27
            id: routes
          - pos: 34
            id: running
          - pos: 42
            id: state
      - pos: 48
        string: absent

```

***Captured state***

The `capturer.CapturedStateByName` is a `map[string]interface{}`, like the one
generated by `yaml.Unmarshal`.


Printing the captured state with `fmt.Println(capturedStateByName)` will output
the following: 

```
map[base-iface:map[interfaces:[map[ipv4:map[address:[map[ip:10.244.0.1 prefix-length:24]] dhcp:false enabled:true] name:eth1 state:up type:ethernet]]] base-iface-routes:map[routes:map[running:[map[destination:0.0.0.0/0 next-hop-address:192.168.100.1 next-hop-interface:eth1 table-id:254] map[destination:1.1.1.0/24 next-hop-address:192.168.100.1 next-hop-interface:eth1 table-id:254]]]] bridge-routes:map[routes:map[running:[map[destination:0.0.0.0/0 next-hop-address:192.168.100.1 next-hop-interface:br1 table-id:254] map[destination:1.1.1.0/24 next-hop-address:192.168.100.1 next-hop-interface:br1 table-id:254]]]] bridge-routes-takeover:map[routes:map[running:[map[destination:0.0.0.0/0 next-hop-address:192.168.100.1 next-hop-interface:eth1 state:absent table-id:254] map[destination:1.1.1.0/24 next-hop-address:192.168.100.1 next-hop-interface:eth1 state:absent table-id:254] map[destination:0.0.0.0/0 next-hop-address:192.168.100.1 next-hop-interface:br1 table-id:254] map[destination:1.1.1.0/24 next-hop-address:192.168.100.1 next-hop-interface:br1 table-id:254]]]] default-gw:map[routes:map[running:[map[destination:0.0.0.0/0 next-hop-address:192.168.100.1 next-hop-interface:eth1 table-id:254]]]]]
```


### Desired State pipeline

The desired state has some placeholder elements using "{{" as prefix and "}}" 
as suffix, they contain a valid `ast.Command` so the same lexer/parser from 
capture can be used here. The desiredState emitter will fail if something 
else aprt from the `ast.Command.Path` has being set. 


To identify the placeholders the golang `strings.Index` can be use like the 
following example:

```golang
func main() {                                                                   
    desiredStateTemplate := `                                                   
interfaces:                                                                     
  - name: br1                                                                   
    description: Linux bridge with base interface  as a port                    
    type: linux-bridge                                                          
    state: up                                                                   
    ipv4: {{ capture.base-iface.interfaces[0].ipv4 }}                           
    bridge:                                                                     
      options:                                                                  
        stp:                                                                    
          enabled: false                                                        
      port:                                                                     
      - name: {{ capture.base-iface.interfaces[0].name }}                       
  routes:                                                                       
    config: {{ capture.bridge-routes-takeover.running }}                        
    `                                                                           
    bIdx := 0                                                                   
    snippet := desiredStateTemplate                                             
    for {                                                                       
        snippet = snippet[bIdx:]                                                
        bIdx = strings.Index(snippet, "{{")                                     
        if bIdx == -1 {                                                         
            break                                                               
        }                                                                       
        snippet = snippet[bIdx:]                                                
        eIdx := strings.Index(snippet, "}}")                                    
        if bIdx == -1 {                                                         
            panic("missing }}")                                                 
        }                                                                       
        expression := snippet[2:eIdx]                                           
        fmt.Println(expression)                                                 
        bIdx = eIdx + 2                                                         
    }                                                                           
}                             
```

The `expression` there can be converted into a `ast.Command` and the 
the `capturer` golang package is already capable of resolving references
with the function `Walk`, so using `Walk` with converted `expression` will
return the the captured state referenced.

```golang
source := source.New(expression)                                        
parser := parser.New(*source, lexer.NewLexer(source))                   
command, err := parser.Parse()                                          
if err != nil {                                                         
    panic(err)                                                          
}                                                                       
resolved, err := c.WalkPath(currentState, command.Path)                 
if err != nil {                                                         
    panic(err)                                                          
}                                                                       
fmt.Println(resolved)               
```

One of the trickiest part is discovering the level of indenting to render the
yaml snippets from captured state, this can be calculated taking into account
the number of elements in the path.

For example for `capture.base-iface.interfaces[0].name` since the final value 
is "scalar" there is no identation to do.

For `capture.base-iface.interfaces[0].ipv4` we have two levels `interfaces` and 
`ipv4` so we indent by 2

Then for `capture.bridge-routes-takeover.running` we just use one identation 
level.

To concatenate the different snippets a golang strings.Builder can be used.

An alternative to using `strings.Index` can be using the yaml.Node from 
https://pkg.go.dev/gopkg.in/yaml.v3#Node, we can traverse the tree and look for
the placeholder at value and in that case convert to capturedState and do
the `Node.Encode`

```golang
TODO
```

![your-UML-diagram-name](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/qinqon/nmpolicy/main-scenario-dry-run/docs/development/dry-run/figure1.iuml)
