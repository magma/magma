---
id: 
title:CG-NAT support on magma to scale up concurrent sessions
hide_title: true
---
# CG-NAT support on magma to scale up concurrent sessions
Status: Draft
Authors: prabinak
Reviewers: Magma team
Last Updated: 2022-01-12

Discussion at
[https://github.com/magma/magma/issues/11136](https://github.com/magma/magma/issues/11136).

## Overview
Proliferation of wireless and Internet-enabled devices drove the creation of IPv6 as IPv4 addresses were rapidly depleted. All of the RIRs (regional Internet registries) have exhausted off their IPv4 allocations. IPv6 adoption has finally taken off due to wide support from technology vendors and service providers. Given that IPv4-addressed infrastructure will be around for a long time, it is up to service providers to make IP address translation transparent to users. FWA Service providers need a solution that will help them seamlessly optimize network operations that have both IPv4 and IPv6 addressed traffic.

3GPP TR 123.975, CG-NAT recommends service providers deploy native network address translation solutions such as NAT44 and NAT64. It provides carrier-grade scalability by offering a very high number of IP address translations, very fast NAT translation setup rates, high throughput, and high-speed logging. TIP FWA requirement REQ-OCN-14 asks for CG-NAT support on OCN.

Wavelabs as a VAR (Value Added Reseller) for magma 5G SA core received RFPs from FWA service providers. CG-NAT is a feature asked by all the service providers. 

## Purpose
CG-NAT feature on magma will enable FWA service providers offer IPv4 and IPv6 connectivity supporting concurrent sessions on the core network.
Following features will be enabled on magma with CG-NAT support;
Translate between IPv6 and IPv4 addresses
Gives service providers with IPv6 endpoints access to IPv4 content and destinations.
Port allocation for a session is performed dynamically out of assigned blocks.
Provides access to IPv4 services for mobile and wireline IPv6-only networks without encapsulation.
Stateless mapping of private IPv4 addresses to public addresses. So, Stateless implementation improves scalability.

## Use cases
CGNAT must be deployed to enable key capabilities such as:

1.  Enablement of IP address expansion by relying on the CGNAT to overcome the IPv4 address exhaustion, with the support of NAT64/DNS64 and NAT46 seamless IPv4/v6 connectivity
2.  Enhanced threat prevention by hiding subscribers’ and infrastructures’ IP addresses from the Internet.
3.  High scalability to support the rapid growth in the number of subscribers and devices to substantially increase revenue.
4.  Carrier Grade NAT as a Lifecycle Strategy: Service providers need to implement a network address translation strategy that includes both a short-term plan to address the preservation of their existing IPv4 address allocation and a long-term plan to seamlessly migrate to an IPv6 infrastructure. This requires a solution that provides a robust set of Carrier Grade Network Address Translation capabilities and addresses the entire lifecycle of the transition from IPv4 to IPv6.

## Design
This Proposal intends to implement python based CGNAT and integrate with pipelined. It generates CGNAT rules using netmap and handled by pipelined on UPF.

The Implementation is a "py-cgnat" Python library and CLI program for generating firewall rules to deploy Carrier-Grade NAT, besides translating a given IP and port to its private address and vice versa. The methodology consists in building netmap rules at 1:32 public-private ratio, mapping a range of 2000 ports for each client. Works for any netmask, since that follow the 1:32 ratio:

The following tasks are need to implement for achieve this functionality:
1.  Create a new optional controller based on configuration in pipelined.yml.
2.  Define the services in PipelineD for enable or disable.
3.  For generating the rules using pycgnat CLI For translating a private IP to its public one, use the direct option:
	pycgnat 100.64.0.0/20 203.0.113.0/25 trans --direct 100.64.2.15
    pycgnat 100.64.0.0/20 203.0.113.0/25 trans -d 100.64.2.15

4.  For translatig a public IP and port to its private IP correspondent, use the reverse option:
    pycgnat 100.64.0.0/20 203.0.113.0/25 trans --reverse 203.0.113.20:13578
    pycgnat 100.64.0.0/20 203.0.113.0/25 trans -r 203.0.113.20:13578
	
5.  To use these functionalities directly in new controller using import.
6.  Multiple sessions can configure in a iptables single term.
7.  Below features are available in controller:
    Generate CGNAT rules for AWS based on Private address pool and Public adddress pool target from netmap.
    Calculate the public IP and port range from private IP given.
	Maintain Dict containing the public_ip and port_range for the query.
	Calculate the private IP and port range from public IP given.
	Maintain Dict containing the private_ip and port_range for the query.

## Delivery Approach
Feature will be delivered in one milestone with the following 6 process gates
1. Design
2. Development & Unit Testing
3. code review 
4. Integration testing
5. resolve integration issues and regression issues
6. System test


### Milestone1 - Support for CG-NAT and Integrate with PipelineD

Tasks to be handled on new PipelineD controller (CG-NAT App)
1.  Generate CGNAT rules for Magma VM on Private address pool and Public adddress pool target from netmap.
2.  Calculate the public IP and port range from private IP given.
3.	Maintain Dict containing the public_ip and port_range for the query.
4.	Calculate the private IP and port range from public IP given.
5.	Maintain Dict containing the private_ip and port_range for the query.
6.  Make parser as IP:port format into tuple.
7.  Split an IPv4 network into a list of subnets according to given netmask.
8.  CLI Stub for functionality verification.

Tasks to be handled on PipelineD
1.  Add PipelineD config file for enable/disable APP.
2.  Add Service manager of pipelined for start the APP.
3.  Provide the GRPC method for configure the IP address range.
 
Unit tests will be added for all new functions introduced. 

## Test Plan
 
Following is the set of tests or scenarios to verify dual stack Support.
### Integration Testing using CLI Stub.
1.	Execute the PDU Session establishment, Release with IPv4 type.
2.	Execute the PDU Session establishment, Release with IPv6 type.
3.	Traffic Testing for IPv4 and IPv6


## Feature Roadmap

Feature will be delivered in One Milestones. Each milestone duration is 45 calendar days. 

MS   		FUNCTIONAL AREA       DELIVERABLES        	                               COST
M1.0		New Controller	      Configure, generating rules,
                                      Parsing.
		------------------------------------------------------------					
		PipelineD	      Configure, Start Service,			
			              GRPC message support                                     $51500
		--------------------------------------------------------
		CLI STUB	      Add CLI stub for functionality verification
				      UT for all respective modules.
		------------------------------------------------------------------
		Integration Testing   Integration Test cases are passed on the designated
		                      simulation environment.
-------------------------------------------------------------------------------------------


## Reference

https://www.etsi.org/deliver/etsi_tr/123900_123999/123975/13.00.00_60/tr_123975v130000p.pdf
https://github.com/williamabreu/py-cgnat/blob/master/pycgnat
https://www.fortinet.com/solutions/mobile-carrier/4g-5g/carrier-grade-nat
https://www.f5.com/pdf/products/big-ip-cgnat-datasheet.pdf
https://www.cisco.com/c/en/us/td/docs/routers/crs/software/crs-r6-6/cgnat/configuration/guide/b-cgnat-cg-crs-66x/m-implementing-cgn-crs.html#concept_FC12D656CA794B1899D86C2E6E1EF883

