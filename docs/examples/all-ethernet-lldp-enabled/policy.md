{% raw %}
capture:
  ethernets: interfaces.type=="ethernet"
  ethernets-lldp: capture.ethernets | interfaces.lldp.enabled:=true

desiredState:
  interfaces: "{{ capture.ethernets-lldp.interfaces }}"
{% endraw %}
