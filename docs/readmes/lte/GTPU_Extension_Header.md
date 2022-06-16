---
id: GTPU_Ext_Header
title: GTP Extensions
hide_title: true
---
# GTP Extensions

## Overview

Magma Gateway uses the Linux networking stack and OVS to program the packet pipeline on the gateway. The GTP-U protocol entity provides packet transmission and reception services to user plane entities in the RNC, eNodeB, SGW, PGW and TWAN. The GTP-U protocol entity receives traffic from a number of
GTP-U tunnel endpoints and transmits traffic to a number of GTP-U tunnel endpoints.

The GTP-U header is a variable length header whose minimum length is 8 bytes. There are three flags that are used to signal the presence of additional optional fields: the PN flag, the S flag and the E flag. The PN flag is used to signal the presence of N-PDU Numbers. The S flag is used to signal the presence of the GTP Sequence Number field. The E flag is used to signal the presence of the Extension Header field.

Extension Header flag (E): This flag indicates the presence of the Next Extension Header field. When it is set to '0', the Next Extension Header field is not present. When it is set to '1', the Next Extension Header field is present.

## GTP-U Extension Header

The Extension Header Length field specifies the length of the particular Extension header in 4 octet units. The Next Extension Header Type field specifies the type of any Extension Header that may follow a particular Extension Header. If no such Header follows, then the value of the
Next Extension Header Type shall be 0.

	 0                   1                   2                   3
     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
     +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
     | Ext-Hdr Length|                                               |
     +-+-+-+-+-+-+-+-+                                               |
     |                  Extension Header Content                     |
     .                                                               .
     .                                               +-+-+-+-+-+-+-+-+
     |                                               |  Next-Ext-Hdr |
     +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

In 5G SA, each flow is forwarded based on the appropriate QoS rules. QoS rules are configured by SMF as QoS profiles to UP components and these components perform QoS controls to PDUs based on rules. In downlink, a pipelineD pushes QFI into an extension header, and transmits the PDU to RAN.  In uplink, each UE obtains the QoS rule from SMF, and transmits PDUs with QFI containing the QoS rules to the RAN.

### Transfer of PDU Session Information for Uplink Data Traffic

The Transfer of PDU Session Information for uplink data packets involves transfer of control information elements related to the PDU Session from NG-RAN to Pipelined.

The UL PDU SESSION INFORMATION frame includes a QoS Flow Identifier (QFI) field associated with the transferred packet.

The below Information Elements present in the PDU Session Information frame:

PDU Type: The PDU Type indicates the structure of the PDU session UP frame. The field takes the value of the PDU Type it identifies: "0" for PDU Type 0. Value range: {0= DL PDU SESSION INFORMATION, 1=UL PDU SESSION INFORMATION, 2-15=reserved for future PDU type extensions}
Spare: The spare field is set to "0" by the sender and should not be interpreted by the receiver.
QoS Flow Identifier: When this IE is present, this parameter indicates the QoS Flow Identifier of the QoS flow to which the transferred packet belongs.
Padding: The padding is included at the end of the frame.

### Transfer of PDU Session Information for Downlink Data Traffic

The Transfer of PDU Session Information for downlink data packets involves transfer of control information elements related to the PDU Session from Pipelined to NG-RAN.

The DL PDU SESSION INFORMATION frame includes a QoS Flow Identifier (QFI) field associated with the transferred packet. The NG-RAN uses the received QFI to determine the QoS flow and QoS profile which are associated with the received packet.

The DL PDU session information frame includes the Reflective QoS Indicator (RQI) field to indicate whether the user plane reflective QoS is to be activated or not. This is only applicable if reflective QoS is activated.

## High Level Design

The following functionality is supported for QFI set in GTP Extension header:
QFI Configuration for uplink flow (Matching criteria)
QFI Configuration for downlink flow (Action criteria)
Uplink traffic match with uplink flow
Downlink traffic takes action with action criteria values.

### Pipelined

Pipelined will extract qos structure from uplink PDR which is coming from SMF. QFI value gets from QoS structure and sets it in uplink flow as a match criteria.
                       match = MagmaMatch(tunnel_id=i_teid, qfi=<qfi_value>, in_port=gtp_portno)
Pipelined will extract qos structure from downlink PDR which is coming from SMF. QFI value gets from QoS structure and sets it in downlink flow as an action criteria.
                       actions.append(parser.OFPActionSetField(qfi = <qfi_value>)

### RYU Controller

Ryu is a Python library that provides an API wrapper for programming OVS.

There are two parts to taking care of QFI functionality support:
Nicira Extension Actions Structures (NXAction)
Nicira Extended Match Structures (OFPMatch)

### OVS Openflow

Using the OpenFlow protocol, the controller can add, update, and delete flow entries in flow tables, both reactively (in response to packets) and proactively. Each flow table in the switch contains a set of flow entries; each flow entry consists of match fields, counters, and a set of instructions to apply to matching packets.

### OVS Kernel

For Egress flow, one needs to build the extension header and set the QFI value. It needs to be Passed down to kernel datapath for pushing it to GTP-U packets.
For Ingress flow, extract the QFI value from skb of GTP header and set to the tunnel key.

## High Level Call Flow

Call flow for QFI set in OVS.

```mermaid
sequenceDiagram
    CPE->>MME(AMF): PDU Session Establishment Request
    MME(AMF)->>SESSIOND: SetSMFContext
    SESSIOND->>PIPELINED: SetSMFSession(GTP)
    PIPELINED->>RYU_CONTROLLER: Set QFI value in match and action criteria
    RYU_CONTROLLER->>OPENFLOW: NXM_NX_QFI
    OPENFLOW->>OVS_DATAPATH: MFF_QFI and set tunnel QFI
    CPE->>OVS_DATAPATH: GTP Traffic for uplink/downlink with QFI value
```

## Common Issues and Troubleshooting

1. First check GTP table entries using `sudo ovs-ofctl -O OpenFlow13 dump-flows gtp_br0 table=0`
2. Check whether GTP port `2152` is listening or not by using `sudo netstat -a|grep 2152`
   If it is not listening, check the `OVS` status using `sudo ovs-vsctl show`
3. `ovs-dpctl` dumps the current state of kernel flow table. This table reflects active connections in the system, so it needs running traffic between UE andinternet server.
4. If the below error is observed by running `sudo ovs-vsctl show`

   Port gtp0
   Interface gtp0
   type: gtp
   options: {key=flow, remote_ip=flow}
   error: "could not add network device gtp0 to ofproto (Address family not supported by protocol)"

   Then GTP tunnel type is changed to `gtpu` in `/etc/network/interfaces.d/gtp` file.
5. GTP kernel module is included as part of OVS module. So no need to insert gtp.ko

6. To setup DEV environment, run the following command on the magma-dev VM
   `sudo bash ~/magma/third_party/gtp_ovs/ovs-gtp-patches/2.15/dev.sh setup`

7. To Run OVS GTP tests on OVS kernel datapath:
   `sudo bash ~/magma/third_party/gtp_ovs/ovs-gtp-patches/2.15/dev.sh build_test`

8. OVS debug logging can be dynamically enabled by `sudo ovs-appctl vlog/set dbg`
   For a specific module,
   `sudo ovs-appctl vlog/set netdev dbg`
   `sudo ovs-appctl vlog/set ofproto dbg`
   `sudo ovs-appctl vlog/set vswitchd dbg`
   `sudo ovs-appctl vlog/set dpif dbg`

9. Debug the traffic issues in fastpath
   Enable the OVS debug logging
   Check the logs using `sudo dmesg`

10. Stop and start the `OVS Service` using below commands:
    `sudo /usr/share/openvswitch/scripts/ovs-ctl stop`
    `sudo /usr/share/openvswitch/scripts/ovs-ctl start`
