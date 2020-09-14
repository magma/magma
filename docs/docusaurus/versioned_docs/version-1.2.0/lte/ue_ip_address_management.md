---
id: version-1.2.0-config_agw_bridged
title: Magma Bridge mode setup.
hide_title: true
original_id: config_agw_bridged
---
# Magma AGW Bridged mode
Following are important aspects for Bridged mode configuration:
- Bridged (Non-NAT) config option for AGW
- UE IP allocation mode
- AGW management plane traffic

This doc would explain configuration option related to these features using orc8r APIs. Some of these configurations could be achieved using NMS UI and json editor.

# Gateway Bridged mode setting
To configure Gateway in Bridged mode, you need to disable NAT on gateway. Note that this is per AGW configuration, so this needs to be changed for every AGW in the network that needs bridged mode.
```
API: /lte/{network_id}/gateways/{gateway_id}/cellular
{
  "epc": {
    ...
    "nat_enabled": false
  },
```
This mode will turn-on ``uplink_br0`` in the AGW to bridge UE traffic to SGi interface. Operator needs to configure intenet gateway node that handle NAT and routing for all user plane traffic.
AGW detects the (default) internet gateway from AGW networking stack and start ARP probes to get gateway mac address. This gateway MAC address is used as target host for all bridged user plane traffic. In case of AGW is not able complete ARP resolution of internet gateway node, it floods user plane packets in SGi L2 segment.
To avoid flooding SGi network, operator can change pipeline.yml file and change fallback mac address to internet GW mac-address ``uplink_gw_mac``.


When switching between NAT modes, it is recommended to reboot AGW host.

# UE IP address allocation
Next important configuration is UE IP allocation. This heavily depends on SGI network design and integration with rest of LTE system.
On UE attach PGW component of AGW allocates IP address for UE. Magma provides multiple options to allocate IP address to UE.
- IP Pool
- DHCP
- Static Assigned IP address
- Multi APN: Separate address space per APN

First two are basic IP allocation algorithms, next two are add-ons.Most straightforward configuration is assigning IP Pool for IP allocation.

UEs can successfully attach and get connected to the Magma AGWs if they have a valid APN configuration in their subscription profiles on the network side. Typically, UEs send APN information explicitly in their connection requests. Magma AGW pulls APN information from the subscription data to verify that UEs have indeed subscription for the requested APN. If APN information is missing from the connection request, AGW picks the first APN in the subscriber profile as the default APN and establishes a connection session according to that default APN. Once APN information is selected, Magma allocates IP address for the UE under that APN.

The first step in IP address allocation is to make sure that mobility related configuration is currectly set in orc8r:
As listed above there are multuple Options to select IP address
for the UE. DHCP, statip IP and multi APN otions make most sense in Bridged (Non-NAT) modes.
These are network wide options, so if you change it all AGW in this network would be using new ip allocation configuration.

## IP Pool (ip-block) based IP allocation
This is default configuration for AGW, to keep backward compability, it is called NAT.
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
In this case, AGW allocates IP address from the ip-block on first come first serve basis.

## DHCP based IP allocation
To enable this feature, you need set ip_allocation_mode to DHCP broadcast.

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

In this mode AGW generates DHCP request packet on UE attach event. IMSI is used for generating MAC address. This packet is sent over SGi network, after receiveing DHCP lease, MobilityD parsed the packet; then AGW complete UE attach operation.

## Static IP assignments
Operator can also assign specific IP address to subset of UE for an APN session. This information is pulled from subscriberDB for requested APN. If AGW could not find static IP for the subscriber, it uses underlying IP allocator (IP-POOL or DHCP) to finish IP allocation.

To enable static Ip allocation feature, set the feature flag:
```
API: /lte/{network_id}/cellular
{
  "epc": {
    ....
    "mobility": {
      "enable_static_ip_assignments": true,
      ....
```

Static IPs can be assigned to specific UE using subscriber API, You need to set the IP for a APN.
````
API: /lte/{network_id}/subscribers
"static_ips": {
    "active_apn_name": "192.168.100.1",
  }
````
Please note you can only assign static IP to active APN.
Static IP can not overlap with configured IP-POOL ip-block or DHCP server subnet.
If operator needs to change Current UE IP, first dettach the UE, then change static IP configuration, wait for 2 min and then attcha again.

## Multi APN IP allocation
First this feature needs to be enabled via orc8r API.
```
API: /lte/{network_id}/cellular
{
  "epc": {
    ....
    "mobility": {
      "enable_multi_apn_ip_allocation": true,
      ....
```

Multi APN IP address assignment today is limited to DHCP based IP assignment. Operator needs to setup different vlan network for each APN. SGi port needs to a trunk port so that it can access vlan nework associated with different APNs.
Along with VLAN you can also define Gateway IP address and Gateway MAC address for UE traffic. This allows AGW to program this information in UE datapath.
If Gateway IP address is not provided it is parsed from DHCP responce for UE IP address allocation.
If gatway Mac address is not provided AGW would ARP GW IP address in the respective VLAN. Internet gateway needs to respond to this ARP, this ARP is broadcast ARP (source IP is '0.0.0.0'). Internet gayteway needs to responds to such ARP request from AGW.
In event of Missing Gateway IP address or failure to resolve ARP, AGW would flood all UE packet on respective L2 segment using ``ff:ff:ff:ff:ff:ff`` mac address.

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

### Internet gateway mac address
If you have defined identical network for all APNs and the gateway mac address is same; there is option to set the mac address
in ``pipeline.yml: uplink_gw_mac``. Once this is set the operator does not need to provide gateway IP for any of APN. No need allow gateway ARP resolution from AGW. This value is used as destination eth address when all methodes of discovering gateway mac address have failed. Downside if it manual configuration needed on each AGW.


# AGW management plane traffic
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

You can also set vlan for SGi management port using following orc8r. This only works in Bridged mode.
```
API: /lte/{network_id}/gateways/{gateway_id}/cellular
{
  "epc": {
    ...
    "sgi_management_iface_vlan": 100
    ...
  },
```
If you set vlan in default mode, AGW could send DHCP requests in this vlan.

### SGi interface IP address
When NAT mode is changed, SGi interface is added to AGW bridge and IP address from SGi port is removed. After this, AGW (without `sgi_management_iface_static_ip` config) requests DHCP ip address for the bridge interface. This could result in IP address change when AGW NAT mode is changed. You can avoid it by statically assigning IP via API or DHCP.
