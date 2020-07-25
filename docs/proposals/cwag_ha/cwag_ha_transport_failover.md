# CWAG HA Transport Failover - Design Proposal

*Status: In-Review*\
*Author: @mpgermano*\
*Reviewers: @cwf-team*\
*Last Updated: 07/24*

# Overview

In H1 of 2020, Magma implemented CWAG high availability through an active/standby design. 
The initial design relied on the external WLC/AP to properly route RADIUS and GRE traffic 
to the IP of the active gateway. This document outlines a proposal to handle transport 
failover as a part of the CWAG HA failover mechanism.

## Transport Failover

The HA design currently relies on the WLC/AP to support redundant GRE tunnels and RADIUS servers. 
It works in the following way:

1. When the active becomes unhealthy, an ICMP drop rule is added to eth1 of the failed active. 
   The RADIUS server is also stopped for this gateway.
2. The WLC detects that the primary transport (both GRE and RADIUS) is down.
3. After a configurable number of seconds, the WLC will begin using the secondary server 
   and tunnel configured. 

This design is problematic for a couple of reasons:

* Not all WLC/APs support transport redundancy (e.g. Cambium)
* The failover determination is separate for the RADIUS and GRE transports. This means that
  control-plane traffic could be sent to the standby while data-plane traffic is routed to 
  the active (or vice versa).

To fix this, we want to have a single IP that can properly route traffic to whichever gateway
is currently designated as active. 

### Using a Virtual IP

We can configure a virtual IP address in the NMS as a network level config that will be 
streamed to each CWAG in the cluster. This VIP will use device eth1 of the CWAG. For example:

```
eth1: flags=4163<UP,BROADCAST,RUNNING,MULTICAST> mtu 1500
 inet 10.22.20.4 netmask 255.255.255.0 broadcast 10.22.20.255
 ether e6:a1:54:03:c0:2f txqueuelen 1000 (Ethernet)
 
eth1:0: flags=4163<UP,BROADCAST,RUNNING,MULTICAST> mtu 1500
 inet **10.22.20.5** netmask 255.255.255.0 broadcast 10.255.255.255
 ether e6:a1:54:03:c0:2f txqueuelen 1000 (Ethernet)
```

The CWF Operator, which manages the active/standby status and associated initialization, will be 
responsible for triggering the transfer of the IPs on failover. When a failover occurs, the standby 
is initialized first. If that succeeds, the active is then initialized. The happy path for a 
failover will consist of the following steps:

Standby:
* Delete the virtual IP 

Active:
* Add the virtual IP
* Send gratuitous ARP for the VIP
* Update default route to use VIP as source IP
* Restart sessiond

By configuring the RADIUS authentication server IP, RADIUS accounting server IP, and GRE tunnel 
endpoint to use the VIP, our steady state behavior should see that traffic is always routed to 
the active.

### VIP Transfer

The interesting cases to consider are when the gRPC call from the operator to the demoted CWAG cannot 
be made successfully. Consider the following cases:

1. Standby gateway no longer exists due to a pod eviction or node failure. In this case, the crash 
   will have caused the VIPs to be released.
2. Network partition between demoted CWAG and the rest of the cluster. This is a scenario in which 
   the VIP is still held by the standby, despite the Kubernetes cluster seeing the pod as 
   unreachable. Given that Carrier Wifi uses an on-premise deployment, this case would occur due to
   hardware failure. The assumption made here is that such cases of hardware failure will result in
   the node being inoperable, rather than partitioning the network. Therefore, this case will not be
   considered below.
3. Connection failure from Operator to CWAG’s health service due to transient connection issues or 
   an unavailable health service.

To handle cases #1 and #3, the operator can utilize Kubernetes’ API server when trying to down the 
VIP on the standby. The operator will:

* Query if the demoted gateway service and pod exists. 
    * If the resource doesn’t exist or is an initialization state, it is safe to transfer the VIPs. 
* If the resource does exist then the operator sends a gRPC call to the CWAG’s health service to disable the VIP.
    * If this succeeds, then proceed to promote the active.
* If the RPC call to down the VIP fails, there are two cases that could be occurring. 
    * There is a connectivity issue between the operator and the demoted gateway.
    * The pod was just recreated but is yet to be fully initialized. This can occur as pod recreation 
      takes 15-20s while full gateway initialization takes ~1min (using a virtlet+CWAG pre-installed 
      image). 
* Thus if the RPC call fails, the operator will query the Kubernetes API to delete, and thus recreate,
  the demoted gateway. This is done to aggressively ensure that both gateways don’t hold the VIP at 
  the same time.

Many times the deletion of the standby pod will be overkill, occurring after the pod was just recreated. 
The downside to this is that the demoted gateway will be unavailable for a small duration longer. 
The upside is that we’re guarding against potential dual ownership of the VIP. 
This tradeoff seems reasonable.
