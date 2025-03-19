#!/bin/bash

if usr/share/openvswitch/scripts/ovs-ctl status | grep ovsdb-server | grep -q running && 
    /usr/share/openvswitch/scripts/ovs-ctl status | grep ovs-vswitchd | grep -q running; then
    exit
else
    exit 1
fi