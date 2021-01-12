---
id: version-1.2.0-config_agw_bridged
title: AGW Bridged Mode
hide_title: true
original_id: config_agw_bridged
---
# Magma AGW Bridged mode

Bridged mode on the Access Gateway refers to any IP allocation strategy that
does not NAT assigned UE IP addresses at the AGW.

There are 3 components to reconfigure to enable bridged mode:

- Bridged (Non-NAT) config option for AGW
- UE IP allocation mode
- AGW management plane traffic

This document will explain how to configure these options using the
Orchestrator API. You can also use the JSON editor components on the NMS to
update these values if you don't have access to the API. Navigate to the
corresponding component on the NMS (e.g. gateway, network, subscriber) and
use the JSON editor to edit its configuration.

## Enable Bridged Mode

To configure an AGW in bridged mode, you need to disable NAT on the AGW.
Note that this is a gateway-level configuration, so you will need to change
this for every AGW in the network that you want to run in bridged mode.

```
API: /lte/{network_id}/gateways/{gateway_id}/cellular
{
  "epc": {
    ...
    "nat_enabled": false
  },
```

This mode will turn on the `uplink_br0` interface in the AGW to bridge UE
traffic to the SGi interface. You will need to configure an internet gateway
node that handles NAT and routing for all user plane traffic. The AGW will
detect the (default) internet gateway from its networking stack and start ARP
probes to get the gateway MAC address. This gateway MAC address is used as the
target host for all bridged user plane traffic. If the AGW is not able complete
ARP resolution, it floods user plane packets in SGi L2 segment. To avoid
flooding SGi network, you can statically set the fallback MAC address of the
internet gateway in `/var/opt/magma/configs/pipelined.yml` with the key
`uplink_gw_mac`.

When switching between NAT modes, we recommend rebooting the AGW. You can do
this from the NMS or manually if you have direct access to your AGW.

## UE IP Address Allocation

How you choose to allocate IP addresses to your UEs will heavily depend on SGi
network design and integration with the rest of your LTE system. On UE attach,
the PGW component of the AGW allocates an IP address for the UE. We currently
support the following IP allocation options:

- IP Pool
- DHCP
- Static Assigned IP address
- Multi APN: Separate address space per APN

The first two modes are basic IP allocation algorithms, and the next two are
add-ons. The most straightforward (and default) configuration is to assign
an IP pool for the AGW to assign UE IPs from.

The first step in IP address allocation is to configure mobility parameters on
Orchestrator. This is a network-wide configuration, so any mobility options
you configure will be passed down to all AGWs in that network.

### IP Pool

This is the default configuration for new AGWs. We call this "NAT" mode on the
API.

```
API: /lte/{network_id}/cellular
{
  "epc": {
    ....
    "mobility": {
      ....
      "ip_allocation_mode": "NAT",
      "nat": {
        "ip_blocks": [
          "192.168.0.0/16"
        ]
      },
      ....
```

In this case, the AGW allocates IP addresses from the configured IP block on a
first-come-first-served basis.

### DHCP

To enable this feature, set `ip_allocation_mode` to `DHCP_BROADCAST`.

```
API: /lte/{network_id}/cellular
{
  "epc": {
    ....
    "mobility": {
      ....
      "ip_allocation_mode": "DHCP_BROADCAST",
      ....
```

In this mode, the AGW will generate a DHCP request packet on each UE attach
event. The MAC address for the DHCP request will be deterministically generated
from the UE's IMSI. This DHCP request packet is sent over the SGi network.
After receiving the corresponding DHCP lease, `mobilityd` service on the AGW
will parse the packet and complete the attach procedure.

### Static IP Assignments

You can also assign a specific IP address to any UE for an APN session on the
subscriber entity. If static IP assignments are enabled but a UE does not have
an assigned IP, the AGW will fall back to the network-wide IP allocation
strategy.

To enable this feature, first set the appropriate feature flag on the network:

```
API: /lte/{network_id}/cellular
{
  "epc": {
    ....
    "mobility": {
      "enable_static_ip_assignments": true,
      ....
```

After this, you can assign a static IP for a specific APN for any registered
subscriber.

````
API: /lte/{network_id}/subscribers
"static_ips": {
    "active_apn_name": "192.168.100.1",
  }
````

Note you can only assign a static IP to an active APN (i.e. provision the
corresponding APN first). The Static IP cannot overlap with the configured
IP block for pool (NAT) allocation or the DHCP server subnet.
If you need to change IP that a UE is currently using, first detach the UE
(you can set the subscriber's status to INACTIVE), update the IP configuration,
then re-attach the subscriber (set status back to ACTIVE).

### Experimental: Multi-APN IP Allocation

NOTE: This is currently an experimental feature.

Like static IP assignments, this is gated behind a network-wide feature flag:

```
API: /lte/{network_id}/cellular
{
  "epc": {
    ....
    "mobility": {
      "enable_multi_apn_ip_allocation": true,
      ....
```

Multi APN IP address assignment today is limited to DHCP based IP assignment.
You will needs set up a different VLAN network for each APN. SGi port needs to
be a trunk port so that it can access the VLAN nework associated with each APN.
Along with VLAN, you can also define Gateway IP address and Gateway MAC address
for UE traffic to be programmed by the AGW in the UE's datapath.

If the internet gateway IP address is not specified, it is parsed from the DHCP
response for the UE IP address allocation. If the internet gateway MAC address
is not specified, the AGW will ARP the Gateway IP address in the respective
VLAN. Your internet gateway needs to respond to this broadcast ARP
(i.e. source IP is `0.0.0.0`) from the AGW. In the event of a missing internet
gateway IP address or failure to resolve ARP, the AGW will flood all UE packets
on the respective L2 segment using the `ff:ff:ff:ff:ff:ff` MAC address.

These configuration options are specified at the AGW level per APN:

````
API: /lte/{network_id}/gateways/{gateway_id}
  "apn_resources": {
    "inet": {
      "apn_name": "inet",
      "gateway_ip": "internet_gateway_host_ip",
      "gateway_mac": "internet_gateway_host_mac",
      "id": "id_of_this_apn_resource",
      "vlan_id": 0
    },

````

There is an option to statically configure the internet gateway MAC address for
an AGW if you're using an identical uplink network for all APNs.
In `/var/opt/magma/configs/pipelined.yml`, set the key `uplink_gw_mac`
to the MAC address of your internet gateway. This value will be used as a
fallback if all other methods of determining the MAC of the uplink internet
gateway have failed.

## AGW Management Traffic

Magma allows to SGi inteface to act as management interface.
By default SGi interface would start DHCP client on the port.
In NATed mode this sets this IP address to SGi interface.
This sets IP address on SGi eth port in NATed mode. In bridged mode
this IP is set to Bridge interface.

There is orc8r API to assign static IP for each AGW.

```
API: /lte/{network_id}/gateways/{gateway_id}/cellular
{
  "epc": {
    ...
    "sgi_management_iface_static_ip": "1.1.1.1/24"
    ...
  },
```

You can also set VLAN for SGi management port using following orc8r.
This only works in Bridged mode.

```
API: /lte/{network_id}/gateways/{gateway_id}/cellular
{
  "epc": {
    ...
    "sgi_management_iface_vlan": 100
    ...
  },
```

If you set VLAN in default mode, AGW could send DHCP requests in this VLAN.

When NAT mode is changed, SGi interface is added to AGW bridge and IP address
from SGi port is removed. After this, AGW (without
`sgi_management_iface_static_ip` config) requests DHCP ip address for the
bridge interface. This could result in IP address change when AGW NAT mode is
changed. You can avoid it by statically assigning IP via API or DHCP.
