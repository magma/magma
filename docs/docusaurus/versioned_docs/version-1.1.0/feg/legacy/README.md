---
id: version-1.1.0-readme
title: Federated Gateway (FeG)
sidebar_label: Overview
hide_title: true
original_id: readme
---
# Federated Gateway (FeG)
The federated gateway provides remote procedure call (GRPC) based interfaces to standard 3GPP components, such as 
HSS (S6a, SWx), OCS (Gy), and PCRF (Gx). The exposed RPC interface provides versioning & backward compatibility, 
security (HTTP2 & TLS) as well as support for multiple programming languages. The Remote Procedures below provide 
simple, extensible, multi-language interfaces based on GRPC which allow developers to avoid dealing with the 
complexities of 3GPP protocols. Implementing these RPC interfaces allows networks running on Magma to integrate 
with traditional 3GPP core components.

![Federated Gateway architecture diagram](https://github.com/facebookincubator/magma/blob/master/docs/readmes/assets/federated_gateway_diagram.png?raw=true "FeG Architecture")

The Federated Gateway supports the following features and functionalities:

1. Hosting centralized control plane interface towards HSS, PCRF, OCS and MSC/VLR on behalf of distributed AGW/EPCs.
2. Establishing diameter connection with HSS, PCRF and OCS directly as 1:1 or via DRA. 
3. Establishing SCTP/IP connection with MSC/VLR.
4. Interfacing with AGW over GRPC interface by responding to remote calls from EPC (MME and Sessiond/PCEF) components,
    converting these remote calls to 3GPP compliant messages and then sending these messages to the appropriate core network 
    components such as HSS, PCRF, OCS and MSC.  Similarly the FeG receives 3GPP compliant messages from HSS, PCRF, OCS and MSC 
    and converts these to the appropriate GRPC messages before sending them to the AGW. 



Please see the **[Magma Product Spec](https://github.com/facebookincubator/magma/blob/master/docs/Magma_Specs_V1.1.pdf)** for more detailed information.

## Federated Gateway Services & Tools

The following services run on the federated gateway:
 - `s6a_proxy` - translates calls from GRPC to S6a protocol between AGW and HSS 
 - `session_proxy` - translates calls from GRPC to gx/gy protocol between AGW and PCRF/OCS
 - `csfb` - translates calls from GRPC interface to csfb protocol between AGW and VLR
 - `swx_proxy` - translates GRPC interface to SWx protocol between AGW and HSS
 - `gateway_health` - provides health updates to the orc8r to be used for 
 achieving highly available federated gateway clusters (see **[Magma Product Spec](https://github.com/facebookincubator/magma/blob/master/docs/Magma_Specs_V1.1.pdf)**
 for more details)
 - `radiusd` - fetches metrics from the running radius server and exports them

Associated tools for sending requests and debugging issues can be found
at `magma/feg/gateway/tools`. 
