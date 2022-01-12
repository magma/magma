---
id: 
title: Dual stack support for user data connections in 5G
hide_title: true
---
# Proposal: Dual stack Support for user data connections in 5G
Status: Draft
Authors: prabinak
Reviewers: Magma team
Last Updated: 2022-01-11

Discussion at
[https://github.com/magma/magma/issues/11129](https://github.com/magma/magma/issues/11129).

## Overview
The main driver for IPv6 is a possible IPv4 (RFC1918) address shortage and the need for address harmonization in 5G, across all networks from day one. Taking into account 5G with network slicing, slice specific networks will be placed, depending on the use case, either close to the cell site (for example, for URLLC use cases) or somewhere centrally to provide a centralized service for all network segments.  A harmonized 5G address concept offers the freedom to gradually evolve 5G networks and services, without the need of usual IP hacks like dedicated backbone VPNs per network segment.

Due to the mixed nature of current transport networks and the long time it takes to transition from IPv4 to IPv6, Dual-Stack support becomes essential.  All Dual-Stack scenarios that occur during the early deployment phase of magma 5G SA and later must be supported by the 5G nodes. The main idea is to use specific GUA (Global Unicast Addresses) IPv6 addresses on the same network for LTE and 5G infrastructure.


## Purpose
Magma needs to support dual stack to enable operators manage both IPv4 & IPv6 network domains for FWA (Fixed Wireless Access) use cases. For use case details refer to Section 8.2.2 of 3GPP specification TS 129.561. TIP OCN has indicated IPv6 support for FWA as a highly desirable requirement. Refer to REQ-OCN-04 of TIP OCN FWA requirements document. 
 It gives network operators an open, flexible and extendable mobile core network solution. Deployment of Large Scale NAT (LSN) also known as Carrier Grade Nat (CGN) in IPv6 network domain is easier.
Magma 5G SA FWA deployments will likely have the following scenarios 
•	Dual-stack connectivity with Limited Public IPv4 Address Pools
•	UEs with IPv6-only connection and Enterprise Application end pionts using IPv6 addresses
•	IPv4 applications running on a Dual-stack Enterprise node with an assigned IPv6 prefix and a shared IPv4 address and having to access IPv4 services.
•	Enterprise Data Centers are accessed by 5G devices that are IPv6 compliant.
•	Design of IPv4 subnets for massive IOT use cases is becoming increasingly difficult


Magma 5G SA today has support for IPv6 address allocation to UEs, but does not have complete IPv6 support in user plane for 5G SA deployments. 
Wavelabs as a VAR (Value Added Reseller) for magma 5G SA core received RFPs from FWA service providers. Ipv4v6 suppport is a feature asked by all the service providers. 
This proposal intends to add dual stack support in user plane to enable IPv6, IPv4v6 based 5G SA FWA deployments and massive IOT deployments.

## Gap analysis 
To support dual stack 5G SA deployments the following issues need to be addressed on magma;
•	Required to support basic PDU Sessions for IPv6 & IPv4v6.
•	Need to support N3 and N6 traffic flows with both IPv4 and IPv6 network addresses.
•	Need to support Usage report functionality for IPv4v6 and IPv6 PDU Sessions. 
•	Need to support dual mode gnodeb interface. 
•	Need to support IPv6 Neighbour Solicitation requests based on local cache information. 
•	Need to support IP Packet Filter Set:
		Source/destination IP address or IPv6 prefix
		Type of Service (TOS) (IPv4) / Traffic class (IPv6) and Mask.
•	Need to support IPV6 Flow Label in the data path

### Example GTP Packet with IPv6 support
The below figure shows an example GTP-U packet over IPv6 that carries an extension header for QFI and an IPv6 PDU.
Outer IPv6 Header's DSCP value(EF) is marked at sender tunnel endpoint node based on QFI value which is contained in GTP-U Extension Header (PDU Session Container).

0                   1                   2                   3
0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
Outer IPv6 Header
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|Version|     DSCP=EF   |               Flow Label              |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|          Payload Length       | NxtHdr=17(UDP)|    Hop Limit  |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                                                               |
+                                                               +
|                                                               |
+                   Source IPv6 Address                         +
|                        2001:db8:1:1::1                        |
+                                                               +
|                                                               |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                                                               |
+                                                               +
|                                                               |
+              Destination IPv6 Address                         +
|                        2001:db8:1:2::1                        |
+                                                               +
|                                                               |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

Outer UDP Header
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|     Source Port = xxxx        |         Dest Port = 2152      |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|         UDP Length            |    UDP Checksum (Non-zero)    |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

GTP-U header
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
| 0x1 |1|0|1|0|0|     0xff      |           Length              |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                           TEID = 1654                         |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|    Sequence Number = 0        |N-PDU Number=0 |NextExtHdr=0x85|
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

GTP-U Extension Header (PDU Session Container)
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|  ExtHdrLen=2  |Type=0 |0|0|   |0|0|   QFI     | PPI |  Spare  |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                    Padding                    |NextExtHdr=0x0 |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

Inner IPv6 Header
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|Version|     DSCP=0    |               Flow Label              |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|          Payload Length       |    NexttHdr   |    Hop Limit  |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                                                               |
+                                                               +
|                                                               |
+                   Source IPv6 Address                         +
|                        2001:db8:2:1::1                        |
+                                                               +
|                                                               |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                                                               |
+                                                               +
|                                                               |
+              Destination IPv6 Address                         +
|                        2001:db8:3:1::1                        |
+                                                               +
|                                                               |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

Payload
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                                                               |
|                                                               |
.                        TCP/UDP/etc., Data                     .
.                                                               .
|                                                               |
|                                                               |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+



## Delivery Approach
Feature will be delivered in 2 milestones. Each milestone will have the following 5 process gates
1. Design
2. Development & Unit Testing
3. code review 
4. Integration testing
5. resolve integration issues and regression issues
Before finishing the 2nd milestone, the feature shall pass the following process gates 
6. System test
7. User Acceptance test

### Milestone1 - Support for IPv6 and IPv4IPv6 sessions

Tasks to be handled on mobilityd
1.	Configure IPv6 block in config file.
2.	Allocate the IP address for correspoding IMSI and PDU Type.
3.	store the IPv6 address and IMSI mapping in local store.
4.	Look up the IPv6 address for each IMSI and vice-versa

Tasks to be handled on Subscriberdb
1.	Need to store the IPv6 address as a UE IP address format.

Tasks to be handled on AMF
1.	Based on PDU session Type, AMF will support mobilityd interface, Subscriberdb interface.
2.	Dual-mode support for PDU session.
3.	Sessiond interface support for UE IPv6 address.

Tasks to be handled on Sessiond
1.	Sessiond will support UE IPv6 address for PDU session establishment & Release.
 
Tasks to be handled on UPF
1.	GTP tunnel creation for UE Ipv6 address.
2.	QoS flow creation for UE IPv6 address.

Unit tests will be added for all new functions introduced. 

### Milestone2 

Support gnb endpoint, enodeb_iface and NAT with IPv6
1.	Dual mode IP configuration support on enodeb_iface and NAT interface in config file.
2.	Validate Uplink & Downlink traffic flows based on PDU Type
3.	Support multi tunnel with IPv6.

Paging Support 
1.	Sessiond needs to invoke mobilityd servicer to retrieve IPv6 address for the given IMSI. 
2.	UPF needs to add paging flows in OVS for the UE IPv6.
 
Support Usage report, charging and credit for IPv6 traffic
1.	UPF needs to send the usage report for IPv6 PDU Sessions.
2.	Sessiond needs to have IPv6 support for charging, monitoring of PDU Sessions.

Add Unit tests for all modules that are changed. 
Make sure the regression tests are passing for all impacted modules. 


## Test Plan
 
Following is the set of tests or scenarios to verify dual stack Support.
### Integration Testing using UERANSIM or equivalent simulator
1.	Execute the PDU Session establishment, Release, Idle mode, Paging sceanrios with IPv4 type.
2.	Execute the PDU Session establishment, Release, Idle mode, Paging sceanrios with IPv6 type.
3.	Execute the PDU Session establishment, Release, Idle mode, Paging sceanrios with IPv4IPv6 type.
4.	Traffic Testing for IPv4 and IPv6


## Feature Roadmap

Feature will be delivered in 2 Milestones. Each milestone duration is 45 calendar days. 

MS   		FUNCTIONAL AREA       DELIVERABLES        	                             COST
M1.0		Mobilityd	  Configure, Allocate IP address,
                                  storing the IP address.
                ----------------------------------------------------					
		Subscriberdb	  Store IPv6 address as a UE IP 			
			          address format
		----------------------------------------------------
		AMF, SMF, UPF	  IPv4v6, IPv6 support as UE IP 		              $51500
			          address for PDU session establishment, 
				  Release, Paging, Idle mode.
				  UT for all respective modules.
		------------------------------------------------------------
		Integration Testing	integration Test cases are passed on the designated
		                      simulation environment.
-------------------------------------------------------------------------------------------
								
M2.0		Dual stack Interface  Dual mode IP configuration support 
                                      in enodeb_iface and nat interface
 				      in config file.
		--------------------------------------------------------
		SMF, UPF	  Usage Report, monitoring			              $51500
			          and charging for IPv6 PDU session.
		---------------------------------------------------------
		UPF, OVS	  Traffic flows for N3 and N6 interface
				  with dual mode configuartion.
                ---------------------------------------------------------
                SIT & UAT         Feature is passing all the tests in Real
			          equipment environment. 
				  Feature user / customer runs the confirmance 
				  tests and accepts the feature. 
----------------------------------------------------------------------------------------------
 
## References

1) https://github.com/magma/magma/issues/2999
2) https://docs.magmacore.org/docs/faq/faq_magma
3) https://www.ietf.org/archive/id/draft-tang-iiot-architecture-00.html
4) https://5ghub.us/ipv6-for-5g-transport/
5) https://www.osti.gov/servlets/purl/1642737
6) https://www.etsi.org/deliver/etsi_ts/129500_129599/129561/16.04.00_60/ts_129561v160400p.pdf

  
