---
id: deploy_config_apn
title: Configure an APN
hide_title: true
---

# Configure Access Point Name

Before any of your subscribers can attach to your network, you will first have
to provision at least one APN.

Access point names allow devices to connect to different networks. With the APN
setting in place, the device is able to build a connection to the gateway.
UEs can successfully attach and connect to the Magma AGWs if they have a
valid APN configuration in their subscription profile on the network side.
Typically, UEs send APN information explicitly in their connection requests.
The AGW pulls APN information from subscription data synced down from
Orchestrator to verify that UEs have a valid subscription for the requested APN.
If APN information is missing from the connection request, the AGW picks the
first APN in the subscriber profile as the default APN and establishes a
connection session according to that default APN. Once APN information is
selected, Magma allocates IP address for the UE under that APN.

## Define APN Configurations using NMS

First, check that there is at least one APN provisioned for your network:

- Navigate to your NMS instance and on the sidebar click on "Traffic" button.

- In the newly opened page, on the top bar select "APNs".
If there are already APNs defined, it will show up on this page.

- You can edit or delete any of the existing APN configurations.
Note that the updates and deletions will be reflected automatically in
subscriber profiles and any new attaches as well as PDN connection requests
will be impacted by these changes.

![APN tab in NMS](assets/nms/APN/Configuration_By_NMS.png)


- You can also add a new APN configuration by clicking on the "Create New APN"
button and filling out the requested fields. After saving these changes, the
page should refresh with the new list of APNs and their configurations.

![Defining new APN](assets/nms/APN/define_new_APN.png)

  APN ID: Specifies the unique name used for this APN which could be configured
  on different terminals

  Class ID: Specifies QCI ranging from 1-9

  ARP Priority Level: Specifies the Priority Level of the specific ARP - An
  integer ranging 1-15. Smaller value demonstrates Higher priority.

  Maximum Required Bandwidth: APN Aggregated maximum bitrate specifies maximum
  bitrate which can be provided for all non-GBR bearers for this particular APN

    Upload(bps): Specifies the Uplink MBR
    Download(bps): Specifies the Downlink MBR

  ARP Pre-emption-Capability:* PCI value - Can be Enabled or Disabled, Enabled
  indicating that one service can preempt the resources of other services with lower
  Priority Level value.

  ARP Pre-emption-Vulnerability:* PVI value - Can be Enabled or Disabled, Enabled
  indicating that resources of a service can be preempted by a service with higher
  Priority Level value.



![Adding APN Configuration](assets/nms/APN/add_apnconfig.png)

## Add APN Configurations to Subscriber Profiles

The next step is to add one or more APN configurations to subscriber profiles
so that UEs can start consuming network services based on their APNs. From the left pane
choose the Sunscriber tab.

- For an existing subscriber, to update its subscription profile, click on the
edit field and perform a multi-select under the "Access Point Names".
- For a new subscriber, fill out the fields and from the Active APNs column choose desired APN(s)
(screenshot below shows the view after clicking on the "Add Subscriber"
button).
Once you save the updated or new subscriber information, the APNs
added to the subscriber profile will show be reflected on the page.

![Adding subscriber with APN](assets/nms/APN/Adding_APN_to_SubscriberPro.png)


## Define APN Configurations using API

## Add a new APN to the network
/lte/{network_id}/apns (https://api-staging.magma.etagecom.io/swagger/v1/ui/#/APNs/post_lte__network_id__apns)
In this section a new APN could to be created while mentioning the network ID in order to bind them.

Related Parameters:

- AMBR: APN Aggregated maximum bitrate specifies maximum bitrate which can be provided for all non-GBR bearers for this particular APN

  - max_bandwidth_dl: Specifies the Downlink MBR in bps

  - max_bandwidth_ul: Specifies the Uplink MBR in bps

- QoS Profile: Specifies QoS parameters for this APN

  - Class_id: Specifies Non-GBR QCI ranging from 5-9

  - Preemption_capability: PCI value - Can be set to true or false, true indicating that one service can preempt the resources of other services with lower Priority Level value.
  - Preemption_vulnerability: PVI value - Can be set to true or false, true indicating that resources of a service can be preempted by a service with higher Priority Level
  - Priority_level: Specifies the Priority Level of the specific ARP - An integer ranging 1-15. Smaller value demonstrates Higher priority

- APN Name: Specifies the unique name used for this APN which could be configured on different terminals
![Adding new APN by API](assets/nms/APN/Add_new_APN_API.png)

## Adding a specific APN to Subscriber’s Profile
### Update an Existing Subscriber profile
/lte/{network_id}/subscribers/{subscriber_id} (https://localhost:9443/apidocs/v1/#/Subscribers/put_lte__network_id__subscribers__subscriber_id_)
APN names could be specified while modifying an existing subscriber profile
![Update existing subscriber by API](assets/nms/APN/Update_Existing_Subscribe_API.png)

### Add a new subscriber to the network
/lte/{network_id}/subscribers (https://api-staging.magma.etagecom.io/swagger/v1/ui/#/Subscribers/post_lte__network_id__subscribers)
APN names could be specified while defining a new subscriber while adding the new IMSI
![Adding new subscriber by API](assets/nms/APN/Add_new_Subscriber_API.png)

### Notes

- The first APN listed under "Active APNs" for each subscriber becomes the
default APN that will be used if the UE omits the APN information in its
connection requests.

- The subscriber data is streamed down to AGWs periodically and the new configs
should be reflected on AGW with some lag (typically, 1-2 minutes).

- To check if an AGW is already updated, on the AGW, you can run the following
command to retrieve local subscriber data:

`subscriber_cli.py get IMSI<15 digit IMSI>`

An example output for a hypothetical user with IMSI 001010000000001 and APNs
`internet`, `ims` is shown below:

```text
sid {
  id: "001010000000001"
}
lte {
  state: ACTIVE
  auth_key: "..." # not shown in this example
}
network_id {
  id: "my_network"
}
state {
}
sub_profile: "default"
non_3gpp {
  apn_config {
    service_selection: "ims"
    qos_profile {
      class_id: 5
      priority_level: 9
    }
    ambr {
      max_bandwidth_ul: 100000
      max_bandwidth_dl: 100000
    }
  }
  apn_config {
    service_selection: "internet"
    qos_profile {
      class_id: 9
      priority_level: 15
    }
    ambr {
      max_bandwidth_ul: 100000000
      max_bandwidth_dl: 200000000
    }
  }
}
```

## APN Override Config in MME

To override the UE requested APN with a network specified APN, you can use the
`enable_apn_correction` and `apn_correction_map_list` in `/etc/magma/mme.yml`.

```yaml
enable_apn_correction: false
apn_correction_map_list:
        - imsi_prefix: "00101"
          apn_override: "magma.ipv4"
```

If `enable_apn_correction` is set to `true`, MME will override the original APN
based on the specified IMSI-prefix filtering.
We support up to 10 IMSI prefix filters and corresponding APNs.
See under Proposals for a more detailed design doc.

## Configuring Multi APNs

It is possible to set up a multi APN architecture on the Magma Access Gateways with
similar/different IP allocation schemes per APN.

### Features covered
- Multi-APN
- DHCP Based IP Allocation
- Static IP Allocation
- Router Mode
- SGi Management Static IP Allocation

### Network Diagram

![Network Diagram4](assets/nms/APN/Netwotk_diagram.png)

### Network Planning

For Multi-APN setup, the AGW uplink port has to be configured as a trunk port with
separate VLANs configured for the various functions (e.g. Management, Static, DHCP etc.).


  - For the DHCP based IP allocation, this setup assumes that there is an L2 path between
  the AGW and the DHCP Server / Gateway.
  - To showcase the Router Mode configuration, static IP range will be routed through the
  AGW. The upstream L3 network will be configured with a static route to forward all return
  traffic towards the SGi management interface IP of the AGW.
  - Finally, to begin configuring this setup, the gateway is expected to be set up in a
  traditional NAT configuration and connected to the orc8r first, else, the gateway cannot
  be configured for multi-apn through the cloud).


### IP & VLAN Mapping

![IP and vLAN Mappinassets/nms/APN/IP_vLAN_Mapping.png)


### NMS Network Configuration

Create two APNs in the NMS:
![Creating two APNS](assets/nms/APN/Create_two_APNs.png)
Configure Network to support DHCP-BROADCAST, Static IP assignment, and Multi-apn support.

![support DHCP-BROADCAST Config](assets/nms/APN/support_DHCP-BROADCAST.png)

### NMS Gateway Configuration

  #### EPC configs
    Disable NAT
    Provision an SGi gateway IP
    Provision an SGi IP (cidr notation)
    (optional) In this setup the management VLAN is set up as the native vlan on the switch.
    If a native VLAN is not desired, add a VLAN ID in this field. Note that the switch
    configurations will be different from what is shown in subsequent sections.

  ![EPC Configs](assets/nms/APN/EPC_configs.png)

  #### APN Resources
    - Add two APNs here; these should match the APNs specified during network configuration.
    - Static IP APN will not have a VLAN designation as this will be routed mode traffic.
      Please note that static IP allocation does support VLAN tagging, however, for the purposes of
      using Routed Mode configuration, Static IP is not being configured with a VLAN tag.
      See the “closing” section for more information.
    - DHCP IP allocation APN (“dhcp”) will be configured with VLAN 168

![APN Resources](assets/nms/APN/APN_Resources.png)

#### Subscriber Configuration

  - Configure Active APN for the subscriber; please note that only one example is shown here.
Configuration for the second APN will also be necessary for the appropriate subscribers.
  - Configure static IP for “static” APN.


![Multi-APN Subscriber Configuration](assets/nms/APN/MultiAPN_Subscriber_Configuration.png)
![APN Static IP Configuration](assets/nms/APN/StaticIP_config.png)

#### Switch Port Configuration

```
interface  0/22
description 'multi-apn magma gw'
switchport mode trunk
switchport trunk native vlan 270
switchport trunk allowed vlan 168,270
exit
```
NOTE: Native VLAN here is being used for Management only and should be excluded if management
traffic should also be tagged.

#### L3 Router Configuration

```
ip route 10.22.128.0/19 10.22.70.97
```

#### Closing

In this setup, traffic from 10.22.168.0/22 network will be handled by the NAT/DHCP server first and
then forwarded to the L3 router. UE traffic from 10.22.128.0/19 will be sent untagged (hence over VLAN 270)
over the SGi interface. The upstream network will need a static route in order to properly route the return traffic.


If routed mode is not desired, the static IP APN (“static”) can also be set up with a VLAN tag.
However, as of v1.6, the default route for the static IP APN needs to be configured locally on the AGW.
Else, traffic will egress over the learned default route of the AGW, which might not be the desired behavior.
Furthermore, the VLAN ID for the Static IP APN will need to be aded to the Edge L2 Switch configuration.
![Closing](magma/docs/readmes/assets/nms/APN/Closing.png)
