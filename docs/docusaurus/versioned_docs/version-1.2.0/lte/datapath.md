---
id: version-1.2.0-openvswitch
title: Magma Datapath
hide_title: true
original_id: openvswitch
---
# Magma Datapath
## Overview
Magma gateway is a software appliance. Magma Gateway uses linux networking stack and OVS to program packet pipeline on the gateway. OVS gives us tremendous programmability to classify and process packets on the gateway

![datapath components](https://github.com/facebookincubator/magma/blob/master/docs/readmes/assets/AGW-OVS.png?raw=true)

OVS configuration has two major component.
1. Ports: AGW mostly uses static ports. OVSDB configuration is done at deployment time. All endpoint are dynamically learned in openflow pipeline.
2. Flows: AGW configures flows primarily via pipelineD. First table is configured by MME.

## Life of packet on AGW:
OVS flow table is performing L2 / L3 stateless classification. Currently Magma does not use any stateful service from OVS. Lets take example of packet from UE toward internet.
1. Packet arrives on downlink port. This is UDP packet with GTP header.
2. Packet is received on the GTP Tunnel device.
3. Tunnel device decap packet and send it to OVS datapath.
4. Datapath extracts flow key from the packet.
5. Datapath lookup flow table in kernel.
6. For the very first packet the flow would be missing in the table
7. Kernel does upcall to userspace and sends packet to ovs userspace ovs-vswitchd daemon.
8. Userspace processes the packet forwading actions and pushes down kernel flow to the OVS module.
9. Subsequent packets from the same UE would be handled in the kernel without involving userspace. OVS kernel module does hash lookup and executes actions associated with the flow. Most of the time final action is to egress the packet on OVS local port.

### After OVS:
1. NAT is done by the Linux kernel NetFilter component.
2. After egressing packet from OVS LOCAL port, it flows through the linux routing table towards forwarding chain.
3. PGW configures NAT table for UE traffic.
4. After NAT Packet is sent to uplink port.
5. QoS: QoS policies are configured Enforcement app, This app create TC classes and then sets packet mark in OVS flow table. This packet mark is used to select TC queue for a packet. Magma today uses TC for rate limiting.
6. After applying QoS policy, the packet is sent ovs to actual physical port.

## Major feature that impact datapath.
### ARP:
Magma configures special table for repsonding to ARP requests of UE ip. Magma AGW has L3 tunnel connection to UE, so the AGW needs to handle all incoming ARP requests.

### UE monitoring.
Magma has special port call 'mtr0' to probe UE, Magma has monitoring service that send periodic icmp packets to probe.

## Useful utilities for debugging datapath issues
1. first check subscriber table using `subscriber_cli.py`
2. Use `enodebd_cli.py` if subscriber table does not look right.
3. `ovs-ofctl` can be used to examine flow tables configured by AGW controllers.
4. `ovs-dpctl` dumps current state of kernel flow table. This table reflects active connections in the system. So it needs running traffic between UE and internet server.
5. Check pipelineD and OVS log files for any issues. OVS debug logging can be dynamically enabled by `sudo ovs-appctl vlog/set dbg`
6. `tcpdump` There are three intermediate net devices that can be checked for UE traffic.
    - downlink device: This would show UE traffic with GTP header.
    - gtp_br0: This would be decapsulated packet for UE.
    - uplink device: This would be NATed packet egressing on uplink device.
7. `sodu dmesg` to check any error messages from kernel.
8. `conntrack -L` to examine NAT table.
