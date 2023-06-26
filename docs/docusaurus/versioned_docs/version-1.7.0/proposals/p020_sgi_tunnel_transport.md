---
id: version-1.7.0-p020_sgi_tunnel_transport
title: L3 transport for AGW
hide_title: true
original_id: p020_sgi_tunnel_transport
---
# L3 transport for SGi

*Status: Draft*\
*Authors: Pravin Shelar*\
*Reviewers: Magma team*\
*Last Updated: 2021-8-15*

## **Objective**

Add support for tunnel from AGW (SGi interface) to aggregator host.

## **Motivation**

Today AGW operator need to have separate networking device to create the tunneling
infrastructure for pushing SGi traffic to aggregator host. This scheme also needs
separate public IP which is expensive to have these days. By providing support for
tunneling traffic from AGW to aggregator host we eliminate both requirement and
reduce operation cost of Magma AGW.

### Other use cases

1. L3 tunnel can be potentially used for routing traffic to PGW aggregator host.
2. L3 tunnel can be used to connect to AGW from the aggregator host.

## Design

Initially IPSec tunnel was considered for tunneling traffic to the aggregator
host. But after exploring other tunneling techniques it became clear that we
need to look at `wireguard` first. There are following advantages using `wireguard`
over IPSec tunnels

1. `wireguard` performs better than IPsec or OpenVPN. There are a bunch of benchmarking
   results on internet[1].
2. It avoids dependency on OVS IPSec support.
3. AGW can apply QoS (rate limiting) policies before egressing traffic on the tunnel.
4. `wireguard` Allows interface based L3 networking device, IPSec has this
   functionality, but it is limited and available on newer kernel. This allows
   simpler integration for the other use cases mentioned above.

So the plan is to implement `wireguard` based tunnel first and then implement
IPSec. This doc covers both tunnel implementation.

## **Wireguard Design**

This feature has two important aspects

1. Configuration of WireGuard tunnel
2. Datapath programming

### **Configuration**

AGW need following configuration for WireGuard tunnel

1. Install required binaries for WireGuard Tunnel
2. Enable or disable SGi Tunnel
3. Type of SGi Tunnel, `wireguard`, IPSec, etc
4. Tunnel configuration
   1. Peer Public IP address
   2. Peer Public Key.
   3. AGW `wireguard` interface IP

### **Datapath Programming for WireGuard Tunnel**

![WireGuard datapath in NonNAT setup](assets/agw_wg_tun1.png)

Magma data path programming for WireGuard Tunnel will involve the following steps

1. Setup the `wireguard` tunnel
    1. Create the tunnel device
    2. Apply the configuration from magma config
        1. If the `wireguard` key is not provided generate the GW key
        2. export the key over directoryD for operator visibility.
    3. Tunnel local IP address, this would allow remote ssh/connectivity directly
       to the AGW from aggregator.
    4. Store the configuration on persistent device.
    5. Link the configuration to systemd to have the device configured on boot up.
    6. Setup the tunnel for NAT traversal.
2. IP Address allocations for UE: `wireguard` tunnel transport would support
   existing IP allocation schemes
    1. IP-POOL or Static IP based UE IP allocation would work the best, since
       there would not be any need to configure uplink-bridge
    2. In case DHCP based allocation is required we need to setup uplink-br0
       for managing DHCP traffic.
3. Setup the flows to drive the traffic to/from tunnel port.

## **Implementation plan**

1. This implementation would focus on introducing `wireguard` tunnel configuration
   via pipelined yml file. The `wireguard` tunnel would be implemented using
   `wireguard` commandline.
2. PipelineD would create WG.
3. Use key configuration from magma configuration and map it to `wireguard`
   tunnel configuration.
4. Egress packet handling: In case `wireguard` is enabled, push all uplink packets
   to the `wireguard` tunnel. This includes SGi traffic, traffic to PGW as well as
   management traffic.
5. Ingress packet: Linux router would take care of forwarding the ingress traffic
   to respective processing next element. either OVS or to AGW management application.

## **Testing Plan**

I am planning on writing integration tests to verify this in a vagrant
development environment. This would done by creating `wireguard` tunnel between
AGW and trfserver.

## IPSec based tunnel

### Advantage over `wireguard`

1. IPSec provides more options for encryption algorithms to choose from.
2. There are more third party devices that support IPSec than Wireguard.
3. With certain tuning/configuration we can achieve comparable performance to
   `wireguard`.

### **Datapath Programming**

Magma data path programming with IPSec will involve the following steps

1. Setup IPSec tunnel
    1. Create the tunnel port in OVSDB
    2. Apply the configuration from magma config
    3. Setup the tunnel for NAT traversal.
2. IP Address allocations for UE: IPSec tunnel transport would support existing
   IP allocation schemes
    1. IP-POOL or Static IP based UE IP allocation would work the best, since
       there would not be any need to configure uplink-bridge
    2. In case DHCP based allocation is required we need to setup uplink-br0
       for managing the DHCP traffic.
3. Setup the flows to drive the traffic to/from the tunnel port.

![IPSec datapath in NonNAT setup](assets/agw_ipsec_tun1.png)

## **Scope of IPSec Tunnel**

This is the scope for version 1.

1. IPSec tunnel would support ‘self-signed certificate’ for Authentication
2. There would be single tunnel port for all APNs
3. VLAN configuration for APN would be ignored.
4. Tunnel device would be L3 tunnel
5. Uplink QoS would not work due to use of skb-mark, IPSec also utilizes
   skb-mark.
6. IPSec would be only supported on ubuntu and kernel (> 4.14).

[1] <https://www.wireguard.com/performance/>
    <https://core.ac.uk/download/pdf/322886318.pdf>
