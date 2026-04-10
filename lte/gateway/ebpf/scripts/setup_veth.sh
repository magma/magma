#!/bin/bash
# setup_veth.sh
# Create veth pairs for testing eBPF integration with OVS

VETH1=${1:-veth0}
VETH2=${2:-veth1}

# Delete if they already exist
ip link del $VETH1 2>/dev/null
ip link del $VETH2 2>/dev/null

# Create veth pair
ip link add $VETH1 type veth peer name $VETH2

# Bring interfaces up
ip link set $VETH1 up
ip link set $VETH2 up

echo "Veth pair $VETH1 <--> $VETH2 created and up"
