# Outbound Roaming

[Link to discussed detailed doc](https://docs.google.com/document/d/17OYrYw6PMxFEWi_Ubx8eD1Qt3UT2VbchR_FIT5ivxI0)

## 1 Objectives

The objective of this work is to add support in Magma for deployments where outbound roaming is required. Specifically, adding support to expose the server sides of s6a and s8 to roaming partner core networks. This will enable outbound roaming to be supported in a 3GPP compliant and standard way which operators are used to.
Software built to accomplish this will be open source under BSD-3 license and will be committed to the Magma software repository under the governance of the Linux foundation, such that it can be effectively maintained in the future releases.

## 2 Background

### 2.1: Terminology

**S6a** refers to the Authentication, Location, and Service definition interface of an HSS in a mobile core.
**S8** refers to the control and user plane interfaces used when user data is home routed during an outbound roaming session.
**PCEF** refers to Policy and Charging Enforcement Function.
**HPLMN** refers to Home Public Land Mobile Network.
**VPLMN** refers to Visited Public Land Mobile Network.

### 2.2 Outbound roaming

Outbound roaming is when a subscriber gains access to services from a visited operator (VPLMN) other than their home network operator (HPLMN). In this case the HPLMN is a Magma deployment, and the VPLMN is another network operator with which this HPLMN has a roaming relationship.

When outbound roaming, the standard interfaces are s6a and s8 which allow the VPLMN to authenticate the SIM credentials with the home network and route user traffic back to the home network. See diagram below from 3GPP TS 23.401 V17.3.0 (2021-12).

![arch](https://user-images.githubusercontent.com/93994458/154709583-14f9bb52-486b-4908-b250-e521627b6f04.png)

Currently, Magma supports only inbound roaming where magma is the VPLMN and another operator/core is the HPLMN. Magma supports the client side of both S6a and S8 interfaces.

By implementing the server side of S6a and S8, this project will add support for outbound roaming to Magma.

## 3 Implementation

### 3.1 Diagrams

The Roaming Gateway diagram shows the services which will run in the Roaming gateway. These services will support the user plane data transport for both inbound roaming GTP aggregation and outbound roaming when traffic is home routed.

![roaming gateway](https://user-images.githubusercontent.com/93994458/154709967-7cd453ba-2070-4733-9e3f-dcc6f38624d2.png)

Overall, when combined with the orchestrator and federation gateway, the proposed architecture is depicted below showing all key interfaces for inbound and outbound roaming.

![arch](https://user-images.githubusercontent.com/93994458/154710530-f535e6c1-8221-4de4-9e7b-5fc974f2ee3d.png)

The proposed architecture gives the Magma deploying operator the ability to locate the user plane carrying Roaming Gateway outside of AWS where data transport rates are more competitive and/or in a location favorable for connection to the GRX/IPX network.

### 3.2 Scope of Change

#### 3.2.1 Roaming Gateway Services

The following existing services will be ported to run in the Roaming Gateway:

- Health
- S8_proxy
- PolicyDB
- SessionD
- PipelineD
- MobilityD
- SubscriberDB
- Control Proxy
- Magmad
- Eventd
- td-agent-bit

Similar to FeG and AGW, Roaming Gateway will be a full Magma gateway with control proxy connection to orchestrator and health service.

##### 3.2.1.1 Health (modified)

In addition to health metrics added in the GTP Gateway proposal for inbound roaming, the following health metrics will also be added to health:

- S8-U Inbound Traffic Rate (Mbps)
- S8-U Outbound Traffic Rate (Mbps)
- S8-U Inbound Bytes
- S8-U Outbound Bytes
- Active Inbound Sessions
- Active Outbound Sessions

##### 3.2.1.2 PolicyDB + SessionD + PipelineD + MobilityD(modified)

The combination of policyDB, sessiond, and pipelined provide the PCEF functionality to the Roaming gateway for home routed data sessions. This ensures that policy rules, QoS profiles, Rating Groups, and charging can be applied to subscriber sessions while getting services from a VPLMN.
Once the GTP-C session is created on FeG, Roaming Gateway sessiond process will be notified via gRPC and OVS flow will be installed to enable GTP-U tunnel and start sending data to the SGi interface of the Roaming Gateway associated with the session APN.  Having sessiond manage the tunnel lifecycle mirrors how it is done in the case of 5G in AGW and should reuse code where possible.

##### 3.2.1.3 MobilityD

UE IP address allocation will be provided by MobilityD. The UE IP request will be moved from MME function to sessiond to provide consistency across Roaming Gateway and AGW. IP allocation via NAT will be supported with a configurable NAT pool for IPs. Other IP address allocation schemes are out of scope.

Currently, default bearer and general internet access is provided to home routed sessions. Dedicated bearer and other services are out of scope.

##### 3.2.1.4 SubscriberDB (modified)

SubscriberDB will be added to Roaming Gateway to support S6a_proxy originated message handling of AIR, ULR and PUR responses.

##### 3.2.1.5 S8_Proxy (new)

S8_proxy service will implement the GTP-U protocol for establishing sessions between the Visitor SGW and the
Magma GTP Gateway, the S8-U interface.

##### 3.2.1.6 Control Proxy (modified)

Control Proxy will run on the Roaming Gateway to provide a configuration and message passing interface to the orchestrator. Control proxy and service registry will be updated to support the new GTP_proxy service endpoint.

##### 3.2.1.7 Magmad (modified)

Magmad will run on GTP Gateway to support gateway checkin and configuration management processes. Modifications will be limited to anything needed to handle new mconfig format for new services.

##### 3.2.1.8 Eventd / td-agent-bit (modified)

Eventd and td-agent-bit will run on GTP gateway to support log and metric reporting. Modifications will be to handle new event types and new metrics specific to roaming gateways.

#### 3.2.2 Federation Gateway Services

##### 3.2.2.1 S6a_proxy (modified)

S6a_proxy will be extended on FeG to support AIR, ULR, and PUR messages. These messages will be proxied to the Roaming Gateway where subscriberDB will be running and used to generate response. Content for AIA, ULA, and PUA will be proxied back to FeG where S6a_proxy will translate back to diameter. The main work will be adding support for converting new S6a messages to and from Diameter protocol to gRPC and adding routing to get the gRPC messages from FeG to Roaming Gateway.

New messages to be supported by s6a_proxy:

- Incoming Authentication Information Request (AIR)
- Incoming User Location Update Request (ULR)
- Incoming Purge UE Request (PUR)

##### 3.2.2.2 S8_proxy (modified)

S8_proxy service will be extended to support the GTP-C protocol for establishing sessions between the Visitor SGW and the home PGW, the S8-C interface. The northbound interface will add support 3GPP compliant S8-C signaling for the following messages:

- Incoming Create Session Request
- Incoming Delete Session Request
- Incoming GTP Echo Request

S8_proxy will be extended to convert the Diameter messages above to gRPC for communication with Roaming Gateway service sessiond.

### 3.3 Call flow

In the call flow below, red number messages indicate new functionality supported by a Magma element.

![call flow](https://user-images.githubusercontent.com/93994458/154710857-f8f23a91-415a-48b9-ab08-cf58450a21c6.png)

## 4 Schedule & Roadmap

Schedule and roadmap provided in grant issue.
