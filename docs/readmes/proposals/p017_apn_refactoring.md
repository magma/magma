---
id: p017_apn_refactoring
title: APN Refactoring Proposal
hide_title: true
---

# Proposal: APN Refactoring

Author(s): @ksubraveti

Last updated: 05/27/2021

Discussion at
[https://github.com/magma/magma/issues/7189](https://github.com/magma/magma/issues/7189).

## **Abstract**

Current APN models in magma have issues around poor abstraction and tight coupling
with other entities. This causes inefficiencies in the product and also makes it
harder to maintain and extend it in future.
This proposal is to cleanup existing data models and adding new services
in orc8r and gateways to handle the new data models.

## **Summary**

APN is an access point name. APN is used to identify the packet data network(PDN), the UE wants to be connected to.
From Magma’s perspective, APN configuration consists of two main entities.

* APN - network level entity which defines APN information and Qos characteristics.
* APN resource - gateway level entity which provides transport information for APN

Currently two main issues exists in the way APNs are handled

* Poor abstraction - APN resource entity contains low level network information including vlan_id, gateway_ip, gateway_mac information. It is undesirable to directly associate the low level network information with an APN resource entity and abstracting these details into a separate transport entity is important for enabling the transport to evolve independent of APN resource mapping. For e.g: Immediate use case is for adding vxlan as a new transport mechanism for APN. In future, we could possibly extend this to other forms of transport including GRE tunnels, IP in IP etc.
* APN resource configuration is too tightly coupled with Subscriberdb data models. This makes any further changes to APN config hard and unnecessarily affects subscriber data models.

## **Existing Implementation**
![existing](assets/proposals/p017/apn_old.png)

### **Orc8r**

* Subscriberdb service provides ListSubscribers API which provides a list of all subscribers associated with their complete denormalized APN information. [_https://github.com/magma/magma/blob/master/lte/protos/subscriberdb.proto#L174_](https://github.com/magma/magma/blob/master/lte/protos/subscriberdb.proto#L174)
* PolicyDB service exposes per subscriber APN to rule mapping through streamer interface.

### **Gateway**

* Subscriberdb holds the APN config for each subscriber as it gets streamed down from the orc8r.
* During subscriber session set up, we have two distinct workflows

1. Configure APN transport configuration into OVS((MME <-> Mobilityd <-> Subscriberdb)

* Mobilityd gets APN resource config from subscriberdb and passes transport information to MME.
* MME controller pushes this vlan tagging information to OVS(table 0)

1. Configure APN QOS information into OVS (PolicyDB, MME <-> Sessiond <-> Pipelined)
    1. MME pushes APN AMBR information to sessiond
    2. PolicyDB pushes per subscriber APN specific rules to sessiond
    3. Sessiond combines both this information and builds the flow request to pipelined
    4. Pipelined install flows with apn ambr and qos information into OVS.

## **Proposal**

* Add new transport level entity which encapsulates the data network information specific to each gateway and refactor "apn_resource" and rename "apn_resource" → "gateway_apn_config"
* Add APN DB servicer on the Orc8r to provide APN specific configs. This will include both network level and gateway level configs.
* Add APN DB client on ConfigDB(rename subscriberDB -> ConfigDB) and periodically poll apn configs from Orc8r and store locally on sqlite db.
* Refactor mme and enhance pipelined code to get transport level configs from subscriberdb and directly program APN transport information through pipelined.

## **API changes**

/lte/{network_id}/gateways/{gateway_id}/apn_configs (Add new apn config to the gateway)
gateway_apn_config

* apn_id
* Permit/Block (APN filtering)
* DNS configs (dns_primary, dns_secondary)
* Data Network Config
    * enable_static_ip_assignment
    * Address allocation mode (IP_POOL/DHCP)
    * Ip_block
    * ip6_block
    * gateway_addr [ip/ip6]
    * gateway_mac
    * transport_profile
        * egress interface or logical tunnels - (vlan/vxlan/gre etc)


## Implementation

![new](assets/proposals/p017/apn_new.png)

### Orc8r

* Current subscriberdb in Orc8r should stream the APNs similar to how it streams policy information. It just be a list of APNs associated with the subscriber instead of pushing down the APN configuration for every subscriber.
* New ApnDB service needs to be added on Orc8r for handling APN related configuration. ApnDB can be a grpc servicer colocated within existing subscriberdb(since both share the same fault domain).

### Gateway

* Add an APN db python client on the gateway which can pull the updates periodically from Orc8r and store it locally in SQLite db(more details in the **state management** section below). This can be co located with existing subscriber db service running on the gateway. Additionally, rename the service on the gateway to generic Config DB service.
* Add APN DB rpc servicer to provide APN network and transport level configuration. Given the limited set of APNs on any gateway or network(single digit APNs), we should completely cache this information in memory and should never reach sqlite unless upon a restart.
* Existing Subscriberdb rpc servicer can be slightly modified to pull the additional APN network configuration, IP allocation information to avoid modification on mobilityd side for it to perform IP allocation as it does currently.
* Modify Pipelined to pull the corresponding APN transport configuration and configure OVS independent of subscriber information.
* Given that APN changes are infrequent and involve considerable network reconfiguration, It would be cleaner to restart all the services to ensure that the configs and cached state across all services is consistent and uses the latest APN configuration.

### State Management

In Magma, we store config state in a persistent store and status information in a memory cache status information is periodically pushed to Orc8r. This pattern ensures that Orc8r acts as a management plane and gateway can be fully functional in a headless environment. This is the rationale behind storing the APN information in persistent store such as SQLite.

[Medium Term]

We have some [_kludge with policydb_](https://fb.quip.com/kWZ3AIOc3lyb) persisting its config in redis. This has to be refactored and moved to the newly proposed Config DB. This will ensure consistency in the way we store our configs and will improve redis performance since we can eliminate it persisting the state.
[Longer Term]
We remove references to redis and sqlite in the code and abstract it behind our services such as ConfigDB service and state service.


## **Implementation Plan**

1. NMS/Orc8r (3 weeks)
    1. Defining a new data model for above mentioned APN changes in orc8r and add obsidian handlers to support REST API changes and add corresponding unit tests
    2. Adding APN servicer to stream the APN state southbound and cleaning up APN state in corresponding subscriber state
    3. DB migration script to migrate existing APN information contained within subscribers into a separate APN entity
    4. NMS changes to support new APN API and configure gateway APN config


1. Gateway ConfigDB(1 week)
    1. Adding AGW client to handle streamed APN information and store it in SQLite
    2. Add APN servicer to handle new RPC requests from Pipelined
    3. Add cache to store APN information in memory.
2. Gateway (Mobility/Pipelined)

We have two options to implement this, we will choose it depending on time in 1.6 release.

1. Option 1: this implementation would take less time and handle Partner (Blink) ask about per APN IP pools
    1. Call APN API from mobilityD and cache the APN transport info in mobilityD at start up.
    2. On any IP allocation RCP call, mobilityD can return the required information as today. Rest of the services remain the same.
2. Option 2: This is clean design
    1. mobilityD:
        1. mobilityD only allocates IP address, so it queries IP pool or ip-allocation config from APN config
    2. PipelineD:
        1. SessionD would pass APN name to pipelineD in activateFlows RPC
        2. PipelineD queries APN transport parameters from APN service using the APN name. PipelinedD can cache this info to avoid RPN calls on every session activate flow request.
        3. PipelineD builds flows as per the transport parameters



[Medium Term]

1. Refactor policyDB and move corresponding config state management to ConfigDB

[Long Term]

1. Abstract out AGW references to redis and sqlite from all services and enable plug and play approach to persistent store and memory cache.

## **References**

* https://gist.github.com/karthiksubraveti/0daee7f5446cc72460497e247d427ee6
* _https://raw.githubusercontent.com/magma/magma/d691319dd40e0a7d2822a993f819beca7490e65d/docs/readmes/lte/Attach_call_flow_in_Magma.txt_

