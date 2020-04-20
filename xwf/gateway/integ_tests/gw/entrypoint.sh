#!/bin/bash

[[ -z "${CTRL_IP}" ]] && CtrlIP="$(getent hosts ofproxy | awk '{ print $1 }')" || CtrlIP="${CTRL_IP}"

# Copy files to /etc/magma it must be here and not in dockerfile because the volume
# are shared and may be taint on the local host
cp cwf/gateway/configs/* /etc/magma/
cp xwf/gateway/configs/* /etc/magma/
cp orc8r/gateway/configs/templates/* /etc/magma/

# run the other part of the cwag install ansible here
ANSIBLE_CONFIG=cwf/gateway/ansible.cfg  ansible-playbook cwf/gateway/deploy/cwag.yml -i "localhost," -c local -v

# Create the XWF-M lan gateway, we use it as gw, dhcp server and DNS server
ifconfig gw0 10.100.0.1 netmask 255.255.0.0 up

# remove IPv6 DHCP server
mv xwf/gateway/deploy/roles/dhcpd/files/isc-dhcp-server xwf/gateway/deploy/roles/dhcpd/files/isc-dhcp-server.back
sed '/^INTERFACESv6/s/^/#/' xwf/gateway/deploy/roles/dhcpd/files/isc-dhcp-server.back > xwf/gateway/deploy/roles/dhcpd/files/isc-dhcp-server

# run XWF install ansible here
ANSIBLE_CONFIG=cwf/gateway/ansible.cfg ansible-playbook -e xwf_ctrl_ip="${CtrlIP}" xwf/gateway/deploy/xwf.yml -i "localhost," -c local -v

# run DNS server
dnsmasq

# loop forever
tail -f /dev/null
