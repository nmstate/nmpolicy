{% raw %}
capture:
  ethernets: interfaces.type=="ethernet"
  ethernets-ipv4-dhcp-disabled: capture.ethernets | interfaces.ipv4.dhcp==false

desiredState:
  interfaces: "{{ capture.ethernets-ipv4-dhcp-disabled.interfaces }}"
{% endraw %}
