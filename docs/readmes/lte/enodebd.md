---
id: enodebd
title: eNodeB Configuration
hide_title: true
---
# eNodeB Configuration
### Overview
To handle automatic configuration of eNodeB devices on your network, Magma 
uses the enodebd service. The enodebd service is responsible for handling
the O&M interface between Magma and any connected eNodeB. The enodebd service
can be disabled if you configure your eNodeB devices manually.

*Magma officially supports auto-configuration of the following devices:*
* Baicells Nova-243 Outdoor FDD/TDD eNodeB
  - Firmware Version: BaiBS_RTS_3.1.6
* Baicells mBS1100 LTE-TDD Base Station
  - Firmware Version: BaiStation_V100R001C00B110SPC003
  
*Magma supports the following management protocols:*
* TR-069 (CWMP)

*Magma supports configuration of the following data models:*
* TR-196 data model
* TR-181 data model

Magma is able to configure a single eNodeB connected to the gateway box.
The Magma team plans to add support for auto-configuration of multiple eNodeB
devices connected to a single gateway. The Magma team also plans to add
support for more devices and management protocols.

### User Guide
A few simple steps are required by the user for eNodeB to interface with the
enodeb auto-configuration service. If you have followed our Magma Setup Guide
then these steps should be redundant.

*1. Connect eNodeB to eth1 interface of Magma gateway*

Magma enodebd can only interface with an eNodeB through the eth1 interface.
After connecting your eNodeB to the gateway through the eth1 interface, you
may want to check a few things if auto-configuration is not working.
Check that the gateway's eth1 interface is not managed through DHCP.
Interface eth1 should have a statically configured IPv4 address of
`192.168.60.142`.
Magma will be running a DHCP server to assign an IP address to your connected
eNodeB.

*2. Set eNodeB management server URL to `baiomc.cloudapp.net:48080`*

Magma uses DNS hijacking to point the eNodeB to the configuration server
being run by enodebd. `baiomc.cloudapp.net:48080` will point to
`192.88.99.142`, the IP address that the configuration server is being hosted
on.

*3. Make sure your Magma gateway is registered on the NMS*

Network configuration settings are propagated to the eNodeB starting the NMS
interface. Default configuration settings that will be used if gate

### Basic Debugging
Use the `enodebd_cli.py` tool to check basic status of eNodeB. It also allows
for querying the value of parameters, setting them, and sending reboot requests
to the eNodeB.

Outside of `enodebd_cli.py`, check `/var/log/syslog`.
