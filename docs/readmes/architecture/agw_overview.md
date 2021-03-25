---
id: agw_overview
title: AGW Overview
hide_title: true
---

# Access Gateway Overview

The main service within the access gateway (AGW) is Magmad, which brings up all the services hosted within the AGW. The other major services and components hosted within the AGW are: the MME, S-PGW, Health Checker, Mobilityd, Sessiond/PCEF, Pipelined, PolicyDB, Subscriberdb, OVS-data path, Enodebd, and the Control Proxy. Together, these components help to facilitate and manage data both to and from the user. 

![AGW diagram](assets/agw_services.png)

## Access Gateway Services and  Components

1. **Sctpd** - High availability SCTP interface to radio front end. Sctpd forms the endpoint of the S1-C connection. 

2. **Mme** - Implements S1AP, NAS and MME subcomponents for LTE control plane Also implements SGW and PGW control plane. If the mme is restarted, S1 connection will be restarted and users service will be affected. Mobilityd, pipelined, and sessiond are restarted in this process as well. 

3. **Enodebd** - Enodebd supports management of eNodeB devices that use TR-069 as a management interface. This is used for both provisioning the eNodeB and collecting the performance metrics. It also acts as a statistics reporter for externally managed eNodeBs. It supports following data models:
* Device Data model : TR-181 and TR-098
* Information Data model : TR-196

4. **Magmad** - Parent service to start all Magma services, owns the collection and reporting of metrics of services, and also acts as the bootstrapping client with Orchestrator.

5. **Dnsd** - local DNS and DHCP server for the eNB. 

6. **Subscriberdb** - Magma uses Subscriberdb to enable LTE data services through one network node like AGW for LTE subscribers. It is deactivated for the deployments that make use of the MNO's HSS. It supports the following two S6a procedures: S6a: Authentication Information Request and Answer (AIR/AIA) and S6a: Update Location Request and Answer (ULR/ULA). Subscriberdb also supports these additional functions:
   * Interface with Orchestrator to receive subscriber information such as IMSI, secret key (K), OP, user-profile during system bring-up.
   * Generate Authentication vectors using* *Milenage Algorithm and share these with MME.
   * Share user profile with MME.

7. **Mobilityd** - IP Address Management Service. It primarily functions as an interface with the orchestrator to receive an IP address block during system bring-up. It can also allocate and release IP addresses for the subscriber on the request from S-PGW Control Plane.

8. **Directoryd** - Lookup service where you are able to push different keys and attribute pairs for each key

9. **Sessiond** - Sessiond implements the control plane for the PCEF functionality in Magma. Sessiond is responsible for the lifecycle management of the session state (credit and rules) associated with a user. It interacts with the PCEF datapath through pipelined for L2-L4 and DPId for L4-L7 policies.

10. **Policydb** - PolicyDB is the service that supports static PCRF rules. This service runs in both the AGW and the orchestrator. Rules managed through the rest API are streamed to the policydb instances on the AGW. Sessiond ensures these policies are implemented as specified.

11. **Dpid** - Deep packet inspection service to enforce policy rules.

12. **Pipelined** - Pipelined is the control application that programs the OVS openflow rules. Pipelined is a set of services that are chained together. These services can be chained and enabled/disabled through the REST API. If pipelined is restarted, users service will be affected.

13. **Eventd** - Service that acts like an intermediary for different magma services, using the service303 interface, it will receive and push the generated registered events to the td-agent-bit service on the gateway, so these can be then later sent to Orchestrator. These events will be sent to ElasticSearch where they can be queried. 

14. **Monitord ** - a dynamic service that monitors the CPEs connected to the AGW. It will send ICMP pings to the CPEs connected to the gateway and report if they are active. 

15. **SMSd ** - service that functions as the AGW interface that will sync the SMS information with Orc8r. 

16. **Td-agent-bit** - consists of Fluentbit and is run as a dynamic service on the AGW. To use td-agent-bit the user must modify the gateway magmad configuration. It is used for log aggregation and event logging where it takes input from syslog and the events service and forwards the output to the Orc8r. It is received on the Orc8r by Fluentd where it is stored in Elasticsearch. 

17. **Ctraced** - Service used for managing call tracing on the AGW. The Tshark tool is used for packet capture and filtering. Packet captures are sent back to the orc8r and viewable on the NMS. Preferred usage for call tracing is through the NMS.

18.  **Health** - Health checker service that verifies the state on MME, mobilityd, sessiond and pipelined and cleans corrupt state if necessary.

19. **Control Proxy** - Control proxy manages the network transport between the gateways and the controller. It additionally provides the following functionality:
   * Control proxy abstract the service addressability, by providing a service registry which maps a user addressable name to its remote IP and port.
   * All traffic over HTTP/2, and are encrypted using TLS. The traffic is routed to individual services by encoding the service name in the HTTP/2 :authority: header.
  * Individual GRPC calls between a gateway and the controller are multiplexed over the same HTTP/2 connection, and this helps to avoid the connection setup time per RPC call.

