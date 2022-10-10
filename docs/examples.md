---
title: Examples
toc: true
toc_label: "NMPolicy Examples"
toc_icon: "network-wired"
toc_sticky: true
---


## Linux bridge on top of default gw NIC with DHCP

{% include_relative examples/example.md example="bridge-on-default-gw-dhcp" %}

## Linux bridge on top of default gw NIC without DHCP

{% include_relative examples/example.md example="bridge-on-default-gw-no-dhcp" %}

## OVS SLB bond between primary and secondary nics

It uses the `description` field to filter between primary and secondary NIC.

{% include_relative examples/example.md example="ovs-slb-bond-primary-secondary" %}

## Set all linux bridges down

{% include_relative examples/example.md example="all-linux-bridges-down" %}

## Convert DHCP aware interface to static addressing

{% include_relative examples/example.md example="convert-dhcp-to-static" %}

## Turn LLDP to true at interfaces active interfaces

{% include_relative examples/example.md example="all-ethernet-up" %}

## Create a linux-bridge with all the interfaces matching description

{% include_relative examples/example.md example="bridge-interfaces-by-description" %}
