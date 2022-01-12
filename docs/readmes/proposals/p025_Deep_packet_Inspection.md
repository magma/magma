---
id: 
title: Deep packet Inspection service to enforce policy rules on 5G SA Deployments
hide_title: true
---
# Deep packet inspection service to enforce policy rules on 5G SA Deployments
Status: Draft
Authors: prabinak 
Reviewers: Magma team
Last Updated: 2022-01-12

Discussion at
[https://github.com/magma/magma/issues/11134](https://github.com/magma/magma/issues/11134).

## Overview
DPI (Deep Packet Inspection) is a traffic recognition method that classifies IP traffic in real time. DPI identifies protocols, applications and application attributes. In addition to IP traffic classification, DPI engines extract protocol and application-based metadata,
providing insight into user behavior and application usage. Examples of metadata that can be
extracted from IP traffic include the following:

Metadata category                        Example metadata
Traffic volume 					Per user, per protocol, per application, per flow, per direction.
						
Service detection 				Differentiation between for example Skype audio and video calls

Quality of service              Jitter, throughput, latency, roundtrip time, ramp-up time, packet loss,retransmissions

## Purpose
TIP OCN Functional requirement REQ-OCN-18 for FWA deployment asks for application detection as a minimum DPI requirement. 
Adding DPI support on Magma will help enterprises that subscribe to FWA services get the following benefits
1)	DPI is required to deliver real-time intelligence about traffic to create the most effective solution.
2)	With DPI-enabled FWA Deployments, operators have the ability to tailor policies and adjust the traffic shaper based on time, package, or applications.
3)	DPI will be applied to improve network efficiency, but ultimately it will allow carriers to deliver tailored Edge / Application services that increase customer satisfaction, create differentiation, and provide revenue growth.

4) Policy and charging control: Provides policy control and charging software vendors with the capability to define bandwidth guarantees,
priorities and limits, offer fine-grained QoS for an additional fee and deliver real-time charging and billing support. It can offers a high detection rate (> 95 %) and accurate application identification for policy and billing purposes.

5) Enterprises can leverage DPI to block or throttle access to risky or unauthorized applications, block policy-violating usage patterns or prevent unauthorized data access within corporate-approved applications, stop data exfiltration attempts by external attackers or potential data leaks caused by both malicious and negligent insiders.

Wavelabs as a VAR (Value Added Reseller) for magma 5G SA core received RFPs from FWA service providers. DPI is a feature asked by all the service providers. 

## Gap (Design) Analysis
To provide DPI support on magma  OVS and Pipelined will need the following changes;
1.  Provide the new Opaque DPI Interface that integrates with the open source DPI plugin. (dynamically any 3rd party / external / customer DPI engine can be hooked with OVS seamlessly)
2.  Implement the sample plugin will demonstrate the DPI public API that was given for DPI-enabled-OVS. This sample plugin will simply write out the ethernet packet from OVS to a file.
3.  Implement the test-dpi plugin, to demonstrate the interfacing capablity of DPI with OVS. This can be extended to full-fledged DPI engine as we now have raw ethernet packet from OVS
4.  Implement the DPI engine
    Init the DPI engine
	Destroy the DPI engine
	Traffic processing 
	Logging support (which enables plugin developers to log their contents directly into OVS logging framework. A handy tool for debugging & maintenance)

5.  Kernel datapath modified to clone all the incoming packets to send it to userspace with special label for DPI.
6.  Create DPI controller for marking a flow with an App ID derived from DPI in Pipelined.
    Assigns the App ID to each new IP tuple.
	

    -------------------------------
    |          Table X            |
    |            DPI              |
    |- Assigns App ID to each new |
    |  IP tuple encountered       |
    |- Requires separate          |
    |  DPI engine                 |
    -------------------------------


## Delivery Approach
Feature will be delivered in a single milestone with the following 7 process gates
1. Design
2. Development & Unit Testing
3. code review 
4. Integration testing
5. resolve integration issues and regression issues
6. System test
7. User Acceptance test

The DPI will support for service based and application based traffic.
1.  Implement the test-dpi plugin to verify the DPI functionality.
2.  Add support for 5G rule policy for DPI services and applications.
	•	DPI service interface add with SMF
	•	DPId support for L4-L7 polices.
3.  Implement the Interface for 3rd party DPI engine can plugin.
4.  Integration with network Qos policy.  
5.  Implement the DPI engine.
6.  Implement DPI controller for marking a flow with an App ID derived from DPI in Pipelined.
7.  Implement the CLI for support dynamically load any DPI plugin.

## Test Plan
•	Verify the packet processing in DPI engine using logs
•	Traffic Testing UDP/TCP/ICMP with Video/Audio/Chat 
•	Verify the Logging API to provided to log onto the OVS logfile, directly from dpi-plugin using below Loglevels:
    DPIERR DPIINFO DPIWARN DPIDEBUG
•	Verify the DPI controller functionality in pipelined using test-DPI controller.
•	Introduce a command-line switch to dynamically load any DPI plugin. Here is the syntax
    sudo ovs-vswitchd --dpi-engine=<path to plugin>


## Roadmaps
This section should break out the development roadmap into a number of milestones.

MS1   		FUNCTIONAL AREA       DELIVERABLES        	                            COST
M1.0		Pipelined	      DPI controller 
                                      Classify the app/service_type
				      to a numeric identifier to export.
	------------------------------------------------------------------------				
		Sessiond	      Support QoS with DPI 			
			                      
	-----------------------------------------------------------------------           $68500
		DPI engine	      DPI enginee with OVS compatibility
			
	------------------------------------------------------------------------
		DPI-plugin            Provide the DPI public API for 
		                      any 3rd party / external / customer DPI engine
 		                      can be hooked with OVS seamlessly
	------------------------------------------------------------------------
		Integration Test	  Functionality verification
			                      
		
## Reference
https://www.fiercewireless.com/sponsored/deep-packet-inspection-getting-most-out-5g
https://www.thefastmode.com/expert-opinion/21162-dpi-supporting-sase-for-5g-security
https://www.thefastmode.com/expert-opinion/15934-how-dpi-drives-monetization-in-the-5g-era
https://github.com/kspviswa/dpi-enabled-ovs
https://kspviswa.github.io/dpi-enabled-ovs/
