#!/bin/bash
while inotifywait -e close_write /etc/dhcp/dhcpd.conf; do
    supervisorctl restart dhcpd
done
