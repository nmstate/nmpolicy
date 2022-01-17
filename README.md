---
---

# Introduction

An expressions driven declarative API for dynamic network configuration

# Motivation

When networking configuration for a cluster is needed and all the details 
are common between the nodes in the cluster a NMState yaml configuration is 
enough.

Problems arise when some of the network configuration details 
are different between nodes and depend on the current node network state.

For that a different NMState yaml configuration needs to be generated 
per node and that's not convenient for big clusters and 
also at scale up scenarios.

The NMPolicy goal is to solve this problem.
Given a node network state and a network configuration policy 
(common to the cluster), 
the NMPolicy tool will generate a node specific desired network state.

Previously without the help from NMPolicy a cluster user needed to apply the 
following configurations per node at a three nodes cluster to create a linux-bridge on 
top of an interface and clone the mac, also it has to hardcode the name of the 
interface, that can be different between nodes on some clusters.

node01:
```yaml
desiredState:
  interfaces:
  - name: br1
    type: linux-bridge
    state: up
    mac-address: 00:00:5E:00:00:01
    ipv4:
      dhcp: true
      enabled: true
    bridge:
      options:
        stp:
          enabled: false
        port:
        - name: eth1
```

node02:
```yaml
desiredState:
  interfaces:
  - name: br1
    type: linux-bridge
    state: up
    mac-address: 00:00:5E:00:00:02
    ipv4:
      dhcp: true
      enabled: true
    bridge:
      options:
        stp:
          enabled: false
        port:
        - name: eth1

```

node03
```yaml
desiredState:
  interfaces:
  - name: br1
    type: linux-bridge
    state: up
    mac-address: 00:00:5E:00:00:03
    ipv4:
      dhcp: true
      enabled: true
    bridge:
      options:
        stp:
          enabled: false
        port:
        - name: eth1
```

The example at this page show how to do that without harcoding the nic name
and the mac addresses.

# How it works

It's implemented on top of [nmstate](https://nmstate.io/), nmpolicy generates a 
nmstate desired state as output, given an input of a 
[policy spec](https://nmstate.io/nmpolicy/user-guide/102-policy-syntax.html) and a 
nmstate current state.

This is a simple nmpolicy example to connect a nic that is referenced by a 
default gateway to a bridge:

<!--
{% raw %}
-->
```yaml
capture:
  default-gw: routes.running.destination=="0.0.0.0/0"
  base-iface: interfaces.name==capture.default-gw.routes.running.0.next-hop-interface
desiredState:
  interfaces:
  - name: br1
    description: DHCP aware Linux bridge to connect a nic that is referenced by a default gateway
    type: linux-bridge
    state: up
    mac-address: "{{ capture.base-iface.interfaces.0.mac-address }}"
    ipv4:
      dhcp: true
      enabled: true
    bridge:
        port:
        - name: "{{ capture.base-iface.interfaces.0.name }}"
```
<!--
{% endraw %}
-->

# Use

To start using nmpolicy you can go directly to one of the following 
[documentation](https://nmstate.io/nmpolicy) chapters:
- [library usage](https://nmstate.io/nmpolicy/user-guide/101-library.html)
- [policy syntax](https://nmstate.io/nmpolicy/user-guide/102-policy-syntax.html)
- [examples](https://nmstate.io/nmpolicy/examples.html)
