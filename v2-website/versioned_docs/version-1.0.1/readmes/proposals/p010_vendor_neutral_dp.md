---
---
[TOC]

# Introduction

This document provides the Scope of Work definition for creating a
vendor neutral CBRS Domain Proxy (DP). A Domain Proxy provides proxy and
aggregation services for signaling between CBSD radios (eNBs or CPEs)
and a Spectrum Access System (SAS).

The interface between SAS and CBSD/DP is a standardized interface. The
interface between a DP and the CBSD is not standardized and can be based
on multiple protocols including the WInnForum defined SAS to CBSD/DP
REST API, TR-069, SNMP, or NETCONF. This SoW will include support for DP
to CBSD interface based on the standard WInnForum SAS to CBSD/DP
signaling protocol.

Support for additional DP to CBSD interface protocols is not in scope.

# References

|       |                                                                                                                                                                                                                                                                                                                                                                                            |
|-------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| \[1\] | SSC-Wireless Innovation Forum, “Signaling Protocols and Procedures for Citizens Broadband Radio Service (CBRS): Spectrum Access System (SAS) - Citizens Broadband Radio Service Device (CBSD) Interface Technical Specification”, WINNF-TS-0016-V1.x.x [&lt;u&gt;https://winnf.memberclicks.net/assets/CBRS/WINNF-TS-0016.pdf&lt;/u&gt;](https://winnf.memberclicks.net/assets/CBRS/WINNF-TS-0016.pdf) |
| \[2\] | SSC-Wireless Innovation Forum, “Test and Certification for Citizens Broadband Radio Service (CBRS); Conformance and Performance Test Technical Specification; CBSD/DP as Unit Under Test (UUT)”, WINNF-TS-0122-V1.x.x [&lt;u&gt;https://winnf.memberclicks.net/assets/CBRS/WINNF-TS-0122.pdf&lt;/u&gt;](https://winnf.memberclicks.net/assets/CBRS/WINNF-TS-0122.pdf)                                  |
| \[3\] | SSC-Wireless Innovation Forum, “CBRS Communications Security Technical Specification”, WINNF-TS-0065-V1.x.x [&lt;u&gt;https://winnf.memberclicks.net/assets/CBRS/WINNF-TS-0065.pdf&lt;/u&gt;](https://winnf.memberclicks.net/assets/CBRS/WINNF-TS-0065.pdf)                                                                                                                                            |
| \[4\] | [&lt;u&gt;https://github.com/Wireless-Innovation-Forum/Citizens-Broadband-Radio-Service-Device&lt;/u&gt;](https://github.com/Wireless-Innovation-Forum/Citizens-Broadband-Radio-Service-Device)                                                                                                                                                                                                        |
| \[5\] | [&lt;u&gt;https://github.com/Wireless-Innovation-Forum/Spectrum-Access-System&lt;/u&gt;](https://github.com/Wireless-Innovation-Forum/Spectrum-Access-System)                                                                                                                                                                                                                                          |

# Motivation

A vendor neutral domain proxy provides a number of benefits to a CBRS
deployment:

-   The aggregation and proxy function allows for a single connection
    between a customer network of CBSDs and the SAS. This provides a
    small attack surface by reducing the number of outbound connections
    traversing the customer firewall and edge devices.

-   A domain proxy can be deployed in a controlled stable environment
    allowing it to provide a consistent interface to the SAS on behalf
    of a CBSD which may be subject to intermittent loss of power or
    connectivity.

-   It provides a point of coordination/control for radio channel
    allocation in a multi-vendor deployment.

-   It allows for the application of additional grant requesting and
    maintenance logic which can correct errant behavior and/or provide
    an optimization layer over vendor SAS client algorithms.

# System Architecture

The DP will sit between the SAS and the CBSD or downstream DP.

# &lt;img src="media/image1.png" style="width:5.21875in;height:2.09375in" /&gt;

Neutral DP System Architecture

The vendor neutral DP will run as kubernetes managed containers. The
deployment of the vendor neutral DP can be in any generic kubernetes
environment such as an on premise compute node (or set of nodes) or in a
cloud environment such as AWS.

&lt;img src="media/image2.png" style="width:5.54167in;height:2.72917in" /&gt;

Microservices in the DP

# Specification 

The following sections provide functional requirements included in the
scope of work.

## SAS Client Interface Requirements (interface toward SAS)

&lt;table&gt;
&lt;tbody&gt;
&lt;tr class="odd"&gt;
&lt;td&gt;&lt;strong&gt;Requirement ID&lt;/strong&gt;&lt;/td&gt;
&lt;td&gt;&lt;strong&gt;Description&lt;/strong&gt;&lt;/td&gt;
&lt;td&gt;&lt;strong&gt;Priority&lt;/strong&gt;&lt;/td&gt;
&lt;/tr&gt;
&lt;tr class="even"&gt;
&lt;td&gt;SCI-01&lt;/td&gt;
&lt;td&gt;The SAS client interface shall expose a single SSL connection between the DP and the SAS used for all proxied connections.&lt;/td&gt;
&lt;td&gt;P1&lt;/td&gt;
&lt;/tr&gt;
&lt;tr class="odd"&gt;
&lt;td&gt;SCI-02&lt;/td&gt;
&lt;td&gt;&lt;p&gt;The SAS client interface shall be compliant with [1] when sending the following messages on behalf of proxied CBSDs:&lt;/p&gt;
&lt;ol type="1"&gt;
&lt;li&gt;&lt;p&gt;Registration Request&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;Spectrum Inquiry Request&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;Grant Request&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;Heartbeat Request&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;CBSD Measurement Report&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;Relinquishment Request&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;Deregistration Request&lt;/p&gt;&lt;/li&gt;
&lt;/ol&gt;&lt;/td&gt;
&lt;td&gt;P1&lt;/td&gt;
&lt;/tr&gt;
&lt;tr class="even"&gt;
&lt;td&gt;SCI-03&lt;/td&gt;
&lt;td&gt;The SAS client interface shall validate the SAS connection using a certificate that chains back to a Domain Proxy CA according to [3].&lt;/td&gt;
&lt;td&gt;P1&lt;/td&gt;
&lt;/tr&gt;
&lt;tr class="odd"&gt;
&lt;td&gt;SCI-04&lt;/td&gt;
&lt;td&gt;The SAS client interface shall comply with all interface test cases defined in [2].&lt;/td&gt;
&lt;td&gt;P1&lt;/td&gt;
&lt;/tr&gt;
&lt;tr class="even"&gt;
&lt;td&gt;SCI-05&lt;/td&gt;
&lt;td&gt;The SAS client interface shall be implemented with a TCP_NODELAY socket so that no extra queuing delay is added between the application and the wire.&lt;/td&gt;
&lt;td&gt;P1&lt;/td&gt;
&lt;/tr&gt;
&lt;/tbody&gt;
&lt;/table&gt;

## 

## CBSD Proxy Interface Requirements (interface toward CBSD)

&lt;table&gt;
&lt;tbody&gt;
&lt;tr class="odd"&gt;
&lt;td&gt;&lt;strong&gt;Requirement ID&lt;/strong&gt;&lt;/td&gt;
&lt;td&gt;&lt;strong&gt;Description&lt;/strong&gt;&lt;/td&gt;
&lt;td&gt;&lt;strong&gt;Priority&lt;/strong&gt;&lt;/td&gt;
&lt;/tr&gt;
&lt;tr class="even"&gt;
&lt;td&gt;CPI-01&lt;/td&gt;
&lt;td&gt;&lt;p&gt;The CBSD proxy interface shall accept incoming connections from the SAS client interfaces in CBSDs as defined in [1], receiving and processing request message types including:&lt;/p&gt;
&lt;ol type="1"&gt;
&lt;li&gt;&lt;p&gt;Registration Request&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;Spectrum Inquiry Request&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;Grant Request&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;Heartbeat Request&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;CBSD Measurement Report&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;Relinquishment Request&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;Deregistration Request&lt;/p&gt;&lt;/li&gt;
&lt;/ol&gt;&lt;/td&gt;
&lt;td&gt;P1&lt;/td&gt;
&lt;/tr&gt;
&lt;tr class="odd"&gt;
&lt;td&gt;CPI-02&lt;/td&gt;
&lt;td&gt;The CBSD proxy interface shall comply with security validation requirements for incoming proxy connections according to [3] with certificates that chain to a SAS CA.&lt;/td&gt;
&lt;td&gt;P1&lt;/td&gt;
&lt;/tr&gt;
&lt;tr class="even"&gt;
&lt;td&gt;CPI-03&lt;/td&gt;
&lt;td&gt;The CBSD proxy interface shall be implemented as a pluggable module to internal APIs such that other external interface protocol modules (TR-069, SNMP, NETCONF) can be supported in parallel to the SAS client interface defined in [1] on a per connection basis.&lt;/td&gt;
&lt;td&gt;P1&lt;/td&gt;
&lt;/tr&gt;
&lt;tr class="odd"&gt;
&lt;td&gt;CPI-04&lt;/td&gt;
&lt;td&gt;The CBSD proxy interface shall be implemented with a TCP_NODELAY socket so that no extra queuing delay is added between the application and the wire.&lt;/td&gt;
&lt;td&gt;P1&lt;/td&gt;
&lt;/tr&gt;
&lt;/tbody&gt;
&lt;/table&gt;

##  Aggregation and Proxy Management Requirements

&lt;table&gt;
&lt;tbody&gt;
&lt;tr class="odd"&gt;
&lt;td&gt;&lt;strong&gt;Requirement ID&lt;/strong&gt;&lt;/td&gt;
&lt;td&gt;&lt;strong&gt;Description&lt;/strong&gt;&lt;/td&gt;
&lt;td&gt;&lt;/td&gt;
&lt;/tr&gt;
&lt;tr class="even"&gt;
&lt;td&gt;APM-01&lt;/td&gt;
&lt;td&gt;The aggregation and proxy management process shall support two modes of operation including Active and Passive.&lt;/td&gt;
&lt;td&gt;P1&lt;/td&gt;
&lt;/tr&gt;
&lt;tr class="odd"&gt;
&lt;td&gt;APM-02&lt;/td&gt;
&lt;td&gt;The aggregation and proxy management process shall allow the selection of Active or Passive operation as a global setting for all connections on the CBSD Proxy Interface.&lt;/td&gt;
&lt;td&gt;P1&lt;/td&gt;
&lt;/tr&gt;
&lt;tr class="even"&gt;
&lt;td&gt;APM-03&lt;/td&gt;
&lt;td&gt;The aggregation and proxy management process shall allow the selection of Active or Passive operation on a per CBSD/DP basis on the CBSD Proxy Interface. This selection shall override the global configuration for the specified CBSD.&lt;/td&gt;
&lt;td&gt;P2&lt;/td&gt;
&lt;/tr&gt;
&lt;tr class="odd"&gt;
&lt;td&gt;APM-04&lt;/td&gt;
&lt;td&gt;The aggregation and proxy management process in Passive mode shall perform a proxy and forward function for all messages received on the CBSD Proxy Interface without modification to message body contents.&lt;/td&gt;
&lt;td&gt;P1&lt;/td&gt;
&lt;/tr&gt;
&lt;tr class="even"&gt;
&lt;td&gt;APM-05&lt;/td&gt;
&lt;td&gt;The aggregation and proxy management process in Active mode shall instantiate and maintain two state machines according to [1] on the SAS Client interface and separately on the CBSD Proxy Interface for each CBSD connection.&lt;/td&gt;
&lt;td&gt;P2&lt;/td&gt;
&lt;/tr&gt;
&lt;tr class="odd"&gt;
&lt;td&gt;APM-06&lt;/td&gt;
&lt;td&gt;The aggregation and proxy management process in Active mode shall have an extensible logic module that can create and/or modify message body contents independently for SAS Client Interface and CBSD Proxy Interface for a given CBSD.&lt;/td&gt;
&lt;td&gt;P2&lt;/td&gt;
&lt;/tr&gt;
&lt;tr class="even"&gt;
&lt;td&gt;APM-07&lt;/td&gt;
&lt;td&gt;The aggregation and proxy management process in Active mode shall allow logic modules to be specific to a CBSD class (e.g. FCC ID specific, list of FCC IDs).&lt;/td&gt;
&lt;td&gt;P2&lt;/td&gt;
&lt;/tr&gt;
&lt;tr class="odd"&gt;
&lt;td&gt;APM-08&lt;/td&gt;
&lt;td&gt;&lt;p&gt;The aggregation and proxy management process shall log state information for CBSDs and grants of Passive and Active mode connections in the CBSD &amp; Grant Database with metadata including at least:&lt;/p&gt;
&lt;ul&gt;
&lt;li&gt;&lt;p&gt;CBSD ID (FCC_ID + Serial #)&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;State&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;Message Type&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;Message Body&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;Response code (if applicable)&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;Timestamp&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;SAS Operator Name&lt;/p&gt;&lt;/li&gt;
&lt;/ul&gt;&lt;/td&gt;
&lt;td&gt;P1&lt;/td&gt;
&lt;/tr&gt;
&lt;tr class="even"&gt;
&lt;td&gt;APM-09&lt;/td&gt;
&lt;td&gt;The aggregation and proxy management process shall store state information for CBSDs and grants of Active mode connections in the CBSD &amp; Grant Database.&lt;/td&gt;
&lt;td&gt;P2&lt;/td&gt;
&lt;/tr&gt;
&lt;/tbody&gt;
&lt;/table&gt;

## User Interface Requirements

&lt;table&gt;
&lt;tbody&gt;
&lt;tr class="odd"&gt;
&lt;td&gt;&lt;strong&gt;Requirement ID&lt;/strong&gt;&lt;/td&gt;
&lt;td&gt;&lt;strong&gt;Description&lt;/strong&gt;&lt;/td&gt;
&lt;td&gt;&lt;strong&gt;Priority&lt;/strong&gt;&lt;/td&gt;
&lt;/tr&gt;
&lt;tr class="even"&gt;
&lt;td&gt;UI-01&lt;/td&gt;
&lt;td&gt;&lt;p&gt;The DP UI shall provide the ability to configure a SAS operator with SAS specific details:&lt;/p&gt;
&lt;ul&gt;
&lt;li&gt;&lt;p&gt;Provider Name&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;Host URL&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;User ID&lt;/p&gt;&lt;/li&gt;
&lt;/ul&gt;&lt;/td&gt;
&lt;td&gt;P1&lt;/td&gt;
&lt;/tr&gt;
&lt;tr class="odd"&gt;
&lt;td&gt;UI-02&lt;/td&gt;
&lt;td&gt;The DP UI shall provide the ability to configure a global default SAS operator from the configured SAS operator list.&lt;/td&gt;
&lt;td&gt;P1&lt;/td&gt;
&lt;/tr&gt;
&lt;tr class="even"&gt;
&lt;td&gt;UI-03&lt;/td&gt;
&lt;td&gt;The DP UI shall provide the ability to configure an override SAS operator for a selected set of CBSDs.&lt;/td&gt;
&lt;td&gt;P1&lt;/td&gt;
&lt;/tr&gt;
&lt;tr class="odd"&gt;
&lt;td&gt;UI-04&lt;/td&gt;
&lt;td&gt;The DP UI shall provide the ability to set the global default mode of operation to either Passive or Active&lt;/td&gt;
&lt;td&gt;P1&lt;/td&gt;
&lt;/tr&gt;
&lt;tr class="even"&gt;
&lt;td&gt;UI-05&lt;/td&gt;
&lt;td&gt;The DP UI shall provide the ability to configure override Passive or Active mode on a per CBSD.&lt;/td&gt;
&lt;td&gt;P1&lt;/td&gt;
&lt;/tr&gt;
&lt;tr class="odd"&gt;
&lt;td&gt;UI-06&lt;/td&gt;
&lt;td&gt;The DP UI shall display a table of current connections and their state.&lt;/td&gt;
&lt;td&gt;P1&lt;/td&gt;
&lt;/tr&gt;
&lt;tr class="even"&gt;
&lt;td&gt;UI-07&lt;/td&gt;
&lt;td&gt;&lt;p&gt;The DP UI shall provide the ability to filter historical logs based on:&lt;/p&gt;
&lt;ul&gt;
&lt;li&gt;&lt;p&gt;FCC_ID&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;CBSD Serial #&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;State&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;Message Type&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;Message Body&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;Response code (if applicable)&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;Timestamp&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;SAS Operator Name&lt;/p&gt;&lt;/li&gt;
&lt;/ul&gt;&lt;/td&gt;
&lt;td&gt;P1&lt;/td&gt;
&lt;/tr&gt;
&lt;tr class="odd"&gt;
&lt;td&gt;UI-08&lt;/td&gt;
&lt;td&gt;&lt;p&gt;The DP UI shall allow configuration of users with rights including at least:&lt;/p&gt;
&lt;ul&gt;
&lt;li&gt;&lt;p&gt;Read Only Users&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;Read/Write Users&lt;/p&gt;&lt;/li&gt;
&lt;/ul&gt;&lt;/td&gt;
&lt;td&gt;P1&lt;/td&gt;
&lt;/tr&gt;
&lt;tr class="even"&gt;
&lt;td&gt;UI-09&lt;/td&gt;
&lt;td&gt;The DP UI shall authenticate user credentials upon login of a user and enforce user permissions for the user session duration.&lt;/td&gt;
&lt;td&gt;P1&lt;/td&gt;
&lt;/tr&gt;
&lt;tr class="odd"&gt;
&lt;td&gt;UI-10&lt;/td&gt;
&lt;td&gt;The DP UI shall allow the export/download of log files in csv format that include all fields in requirement UI-07.&lt;/td&gt;
&lt;td&gt;P1&lt;/td&gt;
&lt;/tr&gt;
&lt;tr class="even"&gt;
&lt;td&gt;UI-11&lt;/td&gt;
&lt;td&gt;&lt;p&gt;The DP UI shall allow the user to select a connection/CBSD ID to view real-time updates on messaging including:&lt;/p&gt;
&lt;ul&gt;
&lt;li&gt;&lt;p&gt;CBSD ID (FCC_ID + Serial #)&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;State&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;Message Type&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;Message Body&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;Response code (if applicable)&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;Timestamp&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;SAS Operator Name&lt;/p&gt;&lt;/li&gt;
&lt;/ul&gt;&lt;/td&gt;
&lt;td&gt;P1&lt;/td&gt;
&lt;/tr&gt;
&lt;tr class="odd"&gt;
&lt;td&gt;UI-12&lt;/td&gt;
&lt;td&gt;&lt;p&gt;The DP UI shall provide a summary dashboard with graphical displays of:&lt;/p&gt;
&lt;ul&gt;
&lt;li&gt;&lt;p&gt;Session counts in each state per [1] including at least:&lt;/p&gt;
&lt;ul&gt;
&lt;li&gt;&lt;p&gt;Registered&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;Granted&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;Suspended&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;Authorized&lt;/p&gt;&lt;/li&gt;
&lt;/ul&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;A table of sessions that have changed state in the last 24 hours&lt;/p&gt;&lt;/li&gt;
&lt;/ul&gt;&lt;/td&gt;
&lt;td&gt;P1&lt;/td&gt;
&lt;/tr&gt;
&lt;tr class="even"&gt;
&lt;td&gt;UI-13&lt;/td&gt;
&lt;td&gt;The DP UI shall provide the ability to define a CBSD record based on FCC ID and Serial Number&lt;/td&gt;
&lt;td&gt;P2&lt;/td&gt;
&lt;/tr&gt;
&lt;tr class="odd"&gt;
&lt;td&gt;UI-14&lt;/td&gt;
&lt;td&gt;&lt;p&gt;The DP UI shall provide configuration interface to allow creation, modification, or deletion of CBSD records including the following installation parameters:&lt;/p&gt;
&lt;ul&gt;
&lt;li&gt;&lt;p&gt;CBSD Serial Number&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;FCC ID&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;User ID (SAS account user ID)&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;cbsdCategory&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;CBSD Type&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;antennaAzimuth&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;antennaBeamwidth&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;antennaDowntilt&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;antennaGain&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;height&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;heightType&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;indoorDeployment&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;latitude&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;longitude&lt;/p&gt;&lt;/li&gt;
&lt;/ul&gt;&lt;/td&gt;
&lt;td&gt;P2&lt;/td&gt;
&lt;/tr&gt;
&lt;tr class="even"&gt;
&lt;td&gt;UI-15&lt;/td&gt;
&lt;td&gt;The DP UI shall provide the ability to define a CBSD record as active or inactive allowing the DP to activate/deactivate CBSDs based on configuration.&lt;/td&gt;
&lt;td&gt;P2&lt;/td&gt;
&lt;/tr&gt;
&lt;/tbody&gt;
&lt;/table&gt;

## Deployment Architecture Requirements

&lt;table&gt;
&lt;tbody&gt;
&lt;tr class="odd"&gt;
&lt;td&gt;&lt;strong&gt;Requirement ID&lt;/strong&gt;&lt;/td&gt;
&lt;td&gt;&lt;strong&gt;Description&lt;/strong&gt;&lt;/td&gt;
&lt;td&gt;&lt;strong&gt;Priority&lt;/strong&gt;&lt;/td&gt;
&lt;/tr&gt;
&lt;tr class="even"&gt;
&lt;td&gt;DA-01&lt;/td&gt;
&lt;td&gt;The DP shall be deployable as a kubernetes managed container system.&lt;/td&gt;
&lt;td&gt;P1&lt;/td&gt;
&lt;/tr&gt;
&lt;tr class="odd"&gt;
&lt;td&gt;DA-02&lt;/td&gt;
&lt;td&gt;The DP shall be horizontally scalable based on number of connections/CBSDs.&lt;/td&gt;
&lt;td&gt;P1&lt;/td&gt;
&lt;/tr&gt;
&lt;tr class="even"&gt;
&lt;td&gt;DA-03&lt;/td&gt;
&lt;td&gt;&lt;p&gt;The DP shall be self contained such that it is deployable in any standard kubernetes environment including:&lt;/p&gt;
&lt;ul&gt;
&lt;li&gt;&lt;p&gt;On premises x86 compute platform&lt;/p&gt;
&lt;ul&gt;
&lt;li&gt;&lt;p&gt;Specs TBD&lt;/p&gt;&lt;/li&gt;
&lt;/ul&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;AWS EKS&lt;/p&gt;&lt;/li&gt;
&lt;/ul&gt;&lt;/td&gt;
&lt;td&gt;P1&lt;/td&gt;
&lt;/tr&gt;
&lt;tr class="odd"&gt;
&lt;td&gt;DA-04&lt;/td&gt;
&lt;td&gt;The DP delivery shall include all helm charts and any necessary installation scripts and procedures&lt;/td&gt;
&lt;td&gt;P1&lt;/td&gt;
&lt;/tr&gt;
&lt;tr class="even"&gt;
&lt;td&gt;DA-05&lt;/td&gt;
&lt;td&gt;&lt;p&gt;The DP shall support networking requirements on both SAS Client Interface and CBSD Proxy Interface including:&lt;/p&gt;
&lt;ul&gt;
&lt;li&gt;&lt;p&gt;1Gbps or 10Gbps NIC (copper or fiber for both)&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;DHCP or Static IP based on configuration&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;Multiple IP address support per interface&lt;/p&gt;&lt;/li&gt;
&lt;li&gt;&lt;p&gt;VLAN support based on configuration&lt;/p&gt;&lt;/li&gt;
&lt;/ul&gt;&lt;/td&gt;
&lt;td&gt;P1&lt;/td&gt;
&lt;/tr&gt;
&lt;/tbody&gt;
&lt;/table&gt;

# Schedule and Milestones

|                                                              |          |          |          |          |          |          |          |
| ------------------------------------------------------------ | -------- | -------- | -------- | -------- | -------- | -------- | -------- |
| **Milestones**                                               | **T0+1** | **T0+2** | **T0+3** | **T0+4** | **T0+5** | **T0+6** | **T0+7** |
| MS1: Setup development infrastructure (Airspan eNB + DP)     | x        |          |          |          |          |          |          |
| MS2: POC Implementation of SAS-facing signaling service based on \[1\] |          | x        | x        | x        | x        |          |          |
| MS3: POC CBSD Proxy Interface in Passive mode                |          |          |          | x        | x        | x        |          |
| MS4: GUI development                                         |          |          |          | x        | x        | x        |          |
| MS5: Lab testing and iterative bug-fixing for SAS-signaling (SAS discovery, authentication, CBSD registration, CBSD Heartbeat etc.) |          |          |          | x        | x        | x        |          |
| MS6: Testing and bug-fixing with a small pilot deployment (3-5 CBSDs) |          |          |          |          | x        | x        | x        |
| MS7: Documentation and beta deployment support gate          |          |          |          |          |          |          | x        |

# Q/A Feedback

Q: Can we run on VMware

A: To make it more future proof, we should build it on k8s (not VMs);
but vSphere should have the ability to deploy and manage k8s. so it can
be VMware k8s

Q: High availability: will these be active-active or active-standby?

A: The app itself running on k8s will be active-active. k8s controller
will provide for HA by managing the individual DP application
containers; with k8s managed HA is active-active.

Q: multi-vendor SAS capabilities

A: Yes, the plan is for multi-vendor SAS support. We have UI-2 in place
to allow for selection of a default. UI-3 requirement has been edited to
more specifically allow for batch override of SAS vendor for a selected
set of CBSDs e.g. a vendor group, or all radios at a site.

Q: Audit trail abilities

A: Full audit details will be displayable and exportable including all
configuration changes and SAS and eNB message logging.

Q: Easy eNB location mapping

A: All CBSD data the DP has will be exportable and we can have a
selection pane to allow selection of attributes to export e.g. if you
just want name and lat/long.

Q: What are Active/Passive modes:

A: The idea is to start with a simple Passive mode where there is not
message mangling or internal logic applied as messages are proxied from
the CBSD interface to the SAS interface. However, we want the
architecture to plan for and allow the insertion of logic modules in the
future that sit between the CBSD and SAS interface to allow for more
than a simply proxy. The motivation here is it allows for a DP to mask
undesirable behavior in CBSDs and/or optimize a multi-vendor deployment
by separately controlling SAS and CBSD messaging per CBSD.
