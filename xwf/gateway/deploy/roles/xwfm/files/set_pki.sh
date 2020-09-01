#!/usr/bin/env bash
set -e

mkdir -p /etc/openvswitch

# Copying PKI
cat $ovs_rootca > /etc/openvswitch/cacert.pem
cat $ovs_key > /etc/openvswitch/sc-privkey.pem
crt $ovs_crt >  /etc/openvswitch/sc-cert.pem

# Changing the mode of ovs to be secure
ovs-vsctl set-ssl /etc/openvswitch/sc-privkey.pem \
    /etc/openvswitch/sc-cert.pem /etc/openvswitch/cacert.pem