---
id: readme_agw
title: AGW Services/Sub-Components
sidebar_label: Services/Sub-Components
hide_title: true
---
# AGW Services/Sub-Components
## MME
 MME includes S1AP, NAS and MME_APP subcomponents. MME functions include:

1. S1AP external Interface with eNB
    1.  S1AP ASN.1 encode/decode
    2.  S1AP Procedures
2.  NAS external Interface with UE
    1. NAS message encode/decode
    2. NAS Procedures
    3. NAS state-machine for NAS EMM and NAS ESM protocols
3.  S11 like Interface with unified S-GW & P-GW
    1. Create and delete PDN Sessions
    2. Create/modify/delete default and dedicated bearers
4.  GRPC based S6a like interface towards FedGW
    1. To get authentication vector and subscriber profile to authenticate and authorize the subscriber
    2. To register the serving MME-id with HSS
    3. To receive the HSS initiated subscriber de-registration request
    4. To send purge request to HSS during UE de-registration
    5. To receive HSS reset indication
5. GRPC based SGs like interface towards FeGW
    1. To support NON-EPS services for the subscriber ( CS voice and CS-SMS)
6. Update serving GW-id for the subscriber to the FeGW
7. Statistics to track the number of eNodeBs connected, number of registered UEs, number of connected UEs and number of idle UEs.
8. MME APP maintains UE state machine and routes the message to appropriate modules based on UE state, context and received message.

## S-PGW Control Plane
S-PGW Control Plane functions include:

1. S11 like interface Interface with MME
    1. Create and delete PDN Sessions
    2. Create/modify/delete default and dedicate bearers
2. Interface with MobilityD to allocate and release IP address for the subscriber during PDN connection establishment and release, respectively
3. Interface with Sessiond/PCEF to trigger Gx and Gy session establishment for the subscriber during PDN connection establishment
4. Establish and release GTP tunnel during bearer setup and release

## Health Checker
Health checker reports 2 kinds of health status:
1. Access Gateway specific health which includes:
    * Number of allocated_ips
    * Number of core_dumps
    * Registration_success_rate
    * Subscriber table
2. Generic Health status which includes:
    * Gateway - Controller connectivity
    * Status for all the running services
    * Number of restarts per each service
    * Number of errors per each service
    * Internet and DNS status
    * Kernel version
    * Magma version

## Mobilityd
Mobilityd functions include:

1. Interface with orchestrator to receive IP address block during system bring-up.
2. Allocate and release  IP address for the  subscriber on the request from S-PGW Control Plane.

## Sessiond / PCEF
Sessiond implements the control plane for the PCEF functionality in Magma. Sessiond is responsible for the lifecycle management of the session state (credit and rules) associated with a user. It interacts with the PCEF datapath through pipelined for L2-L4 and DPId for L4-L7 policies.

## Pipelined
Pipelined is the control application that programs the OVS openflow rules. In implementation pipelined is a set of services that are chained together. These services can be chained and enabled/disabled through the REST API.  The README (https://github.com/facebookincubator/magma/blob/master/README.md) describes the contract in greater detail.

## PolicyDB
PolicyDB is the service that supports static PCRF rules. This service runs in both the AGW and the orchestrator. Rules managed through the rest API are streamed to the policydb instances on the AGW. Sessiond ensures these policies are implemented as specified.

## Subscriberdb
Subscriberdb is Magma's local version of HSS. Magma uses Subscriberdb to enable LTE data services through one network node like AGW for LTE subscribers.  It is deactivated for the deployments that make use of the MNO's HSS. It supports the following two S6a procedures:

1. S6a: Authentication Information Request and Answer (AIR/AIA)
2. S6a: Update Location Request and Answer (ULR/ULA)

Subscriberdb functions include:

1. Interface with Orchestrator to receive subscriber information such as IMSI, secret key (K) , OP, user-profile during system bring-up.
2. Generate Authentication vectors using* *Milenage Algorithm and share these with MME.
3. Share user profile with MME.

## OVS - Data path
OVS (http://www.openvswitch.org/) is used to implement basic PCEF functionality for user plane traffic. The control plane applications interacting with OVS are implemented in pipelined.

## Enodebd

Enodebd supports management of eNodeB devices that use TR-069 as management interface. This is used for both provisioning the eNodeB and collecting the performance metrics.It suppots followig data models:
1. Device Data  model : TR-181 and TR-098
2. Information Data model : TR-196


## Control Proxy
Control proxy manages the network transport between the gateways and the controller.

1. Control proxy abstract the service addressability, by providing a service registry which maps a user addressable name to its remote IP and port.
2. All traffic over HTTP/2, and are encrypted using TLS. The traffic is routed to individual services by encoding the service name in the HTTP/2 :authority: header.
3. Individual GRPC calls between a gateway and the controller are multiplexed over the same HTTP/2 connection, and this helps to avoid the connection setup time per RPC call.

## Command Line Interfaces

Several services listed above can be configured using CLIs, located under
magma/lte/gateway/python/scripts. These are:

1. Health Checker: agw_health_cli.py
2. Mobilityd: mobility_cli.py
3. Sessiond: session_manager_cli.py
4. Pipelined: pipelined_cli.py
5. PolicyDB: policydb_cli.py
6. Subscriberdb: subscriber_cli.py
7. Enodebd: enodebd_cli.py
8. State Tracing: state_cli.py

Each of these CLIs can be used in the gateway VM:

```bash
vagrant@magma-dev:~$ magtivate
(python) vagrant@magma-dev:~$ enodebd_cli.py -h

usage: enodebd_cli.py [-h]
                      {get_parameter,set_parameter,config_enodeb,reboot_enodeb,get_status}
                      ...

Management CLI for Enodebd

optional arguments:
  -h, --help            show this help message and exit

subcommands:
  {get_parameter,set_parameter,config_enodeb,reboot_enodeb,get_status}
    get_parameter       Send GetParameterValues message
    set_parameter       Send SetParameterValues message
    config_enodeb       Configure eNodeB
    reboot_enodeb       Reboot eNodeB
    get_status          Get eNodeB status
```
