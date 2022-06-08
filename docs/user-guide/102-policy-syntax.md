---
title: "Policy Syntax"
tag: "user-guide"
---

The policy contains the following fields:
- `capture`: Contains some simple expressions
to evaluate the current network configuration and store the result at 
variables.
- `desiredState`: Contains network state yaml with references to `capture` entries.

__Example:__
```yaml
{% include_absolute 'examples/bridge-on-default-gw-dhcp/policy.md' %}
```

## Capture syntax

The capture is a map of capture entries with an identifier as the key that can 
be referenced with ```capture.[id]```

The map value is a capture entry expression that evaluates the current 
network state and stores the result.

The following is a semi-formal definition of the capture entry expression:

```
<letter> ::= "A" | "B" | "C" | "D" | "E" | "F" | "G"
       | "H" | "I" | "J" | "K" | "L" | "M" | "N"
       | "O" | "P" | "Q" | "R" | "S" | "T" | "U"
       | "V" | "W" | "X" | "Y" | "Z" | "a" | "b"
       | "c" | "d" | "e" | "f" | "g" | "h" | "i"
       | "j" | "k" | "l" | "m" | "n" | "o" | "p"
       | "q" | "r" | "s" | "t" | "u" | "v" | "w"
       | "x" | "y" | "z"
<digit> ::= [0-9]
<number> ::= <digit>+
<identity> ::= <letter> ( <digit> | "-" | <letter> )*
<dot> ::= "."
<path> ::= <identity> ( <dot> ( <identity> | <number> ))*
<string> ::= \" (<all characters>)* \"

<captureid> ::= <identity>
<capturepath> ::= "capture" <dot> <captureid> <path>
<eqoperator> ::= "=="
<eqexpression> ::= <path> <eqoperator> (<string> | <number> | <capturepath>)
<replaceoperator> ::= ":="
<replaceexpression> ::= <path> <replaceoperator> (<string> | <number> | <capturepath>)
<pathexpression> ::= <path>
<expression> ::= <pathexpression> | <eqexpression> | <replaceexpression>
<pipe> ::= "|"
<pipedexpression> ::= <capturepath> <pipe> <expression>
```

### Path ```<path>```
The path expression contains different "steps" separated by dots, each "step" 
can be a key from a map or the index starting with 0 from a list.
```
interfaces.name
routes.running.0.next-hop-interface
```

To reference a capture entry from a path the reserved word `capture` has to be
used followed by a dot and the capture entry name:
```
capture.default-gw.routes.running.0
capture.primary-nic.interfaces.0.name
```

### Equality filter ```<eqexpression>```
Filter the current state based on specific state values. 
The filter follows a simple syntax, similar to jsonpath. 
Values may be explicit or appear as references to other expressions. 
The resulting output is a full NMState state containing only the filtered values.
```
interfaces.name == "eth1"
routes.destination == "0.0.0.0/0"
dns.server == "192.168.1.1"
interfaces.name == capture.default-gw.interfaces.0.name
```

### Path filter ```<pathexpression>```
Filter out current state to include only the data matching the ```<path>```

Following is a path filter to filter out everything except the running DNS
configuration:
```
dns-resolver.running
```

This will create a capture entry with the the following Nmstate
```yaml
dns-resolver:
  running:
    search:
    - redhat.com
    server:
    - 8.8.8.8
```

This can be referened later on with:

```
{% raw %}
"{{ capture.[capture name].dns-resolver.running }}"
{% endraw %}
```

### Replace ```<replaceexpression>```
These commands can replace values from the specified fields 
at the input NMState and they can reference other capture entries.
```
routes.running.next-hop-interface := "br1"
```

### Pipe ```<pipexpression>```
When expressions are piped the output from the left expression is passed 
to the input of the right command.
```
capture.base-iface-routes | routes.running.next-hop-interface := "br1"
```

## Desired state syntax

The state follows [NMState](https://nmstate.io/examples.html) syntax and will include optionally references 
to capture entries so they can be expanded in-place. 
capture references have to be enclosed between {% raw %}```"{{``` and ```}}"```{% endraw %} expressions, the
`desiredState` field can be expressed using JSON or YAML.

The only supported expressions are capture entry reference path like the following
```
capture.base-iface.interfaces.0.mac-address
capture.base-iface.interfaces.0.name
```

For example to override the routes config from a capture ```new-routes``` 
the following can be specified

{% raw %}
```yaml
routes:
  config: "{{ capture.new-routes.running }}"
```
{% endraw %}

Or clone the mac-address from a capture entry ```primary-nic```
{% raw %}
```yaml
interfaces:
- name: br1
  type: linux-bridge
  state: up 
  mac-address: "{{ capture.primary-nic.interfaces.0.mac-address }}"
```
{% endraw %}
