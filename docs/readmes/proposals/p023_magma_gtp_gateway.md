# Proposal: Magma GTP Gateway

Author(s): [@arsenii-oganov]

Last updated: 12/16/2021

## 1. Objectives

The objective of this work is to add a GTP aggregation gateway to the Magma system. This Gateway will aid in scale deployments where S8 inbound roaming is supported by reducing the number of integration points for roaming partners.

Software built to accomplish this will be open source under BSD-3-Clause license and will be committed to the Magma software repository under the governance of the Linux foundation, such that it can be effectively maintained in the future releases.

## 2. Background

In traditional telecom deployments using a centralized (or non-distributed) core, roaming interfaces are limited in number. These interfaces often connect to an IPX/GRX network provider to aggregate roaming traffic into a small set of interfaces, regardless of the number of roaming partners an operator may have.

In the case of S8 roaming interconnections, the connection supports GTP-C and GTP-U traffic between the SGW in the visited network and the PGW in the home network. The SGW IP for the GTP-U endpoint must be routable from the home network. When using an IPX/GRX provider, this means having a globally routable IP address dedicated to the IPX/GRX network interface of an SGW.

Magma has a distributed core architecture. Each cell site, where the Access Gateway (AGW) element usually resides, has an SGW. The requirement that each site have a globally routable IP dedicated to the IPX/GRX connection (i.e. not a general traffic WAN port ISP provided address) is a challenge for operators using Magma core for roaming. Furthermore, IPX/GRX vendors are reluctant to support large numbers of VPN terminations coming from each AGW.

## 3. Implementation

### 3.1 Block Diagram

![Block Diagram](https://user-images.githubusercontent.com/93994458/146327309-6ab4ce9b-3bee-4cae-91c8-699d35c7ae9a.png)

### 3.2 GTP Gateway Scope of Change

Below are system requirements for the GTP Gateway. These requirements assume the presence of Magma support for:

#### 3.2.1 Wireguard/IPsec support in AGW for S8 GTP traffic

Random TEID assignment for S8 tunnels across AGWs
Wireguard / IPsec Connection Termination
Pipelined supports configuration of Wireguard tunnels which connect automatically when the service is started. These tunnels will be used for connecting AGWs to the GTP Gateway.

The GTP Gateway will need to support termination of 10k wireguard connections per instance. Instances will be horizontally scalable. This will produce a reduction in IP address usage of 1/10,000.

#### 3.2.2  GTP Traffic Flow Learning

GTP traffic needs to be routed to the correct tunnel in the downlink direction. The GTP Gateway will use OVS rules (or potentially another mapping / CGNAT solution) to keep track of GTP tunnel to Wireguard tunnel mapping.

The initial implementation will use OVS to learn flows based on UE APN IP Address. Learning will occur by seeing uplink traffic originating from a tunnel, learning the UE source IP address, and creating a flow mapping for return traffic.

Alternate solutions may use @mark attribute of the sk_buff struct or other native linux networking stack functions/attributes for flow mapping. Such solutions will be explored for scale performance and evaluated against the OVS solution.

#### 3.2.3 Orc8r Integration

The GTP Gateway will be integrated in the Magma ecosystem. Integration will have the following components:

- GTP Gateway will be deployable via a script prepared for bare metal host deployment similar to AGW.
- GTP Gateway will have a connection to orc8r for configuration, health and performance metrics similar to FeG.
- GTP Gateway will include control proxy to allow for gRPC between GTP Gateway and other services in the Magma deployment.
- GTP Gateway specific API endpoints will be created to support configuration and health monitoring:

    - List all GTP Gateways
    - Register a new GTP Gateway
    - Delete a GTP Gateway
    - Get a specific GTP Gateway

#### 3.2.4 Gateway_Health Service

Gateway health service for GTP Gateway will report on the following metrics:

- GTP Gateway specific health which includes:

    - Number of active VPN connections
    - Number of active GTP connections
    - Status of PGW VPN (if present)
    - Total Ingress throughput per NIC
    - Total Egress throughput per NIC

- Generic Health status which includes:

    - Gateway - Controller connectivity
    - Status for all the running services
    - Number of restarts per each service
    - Number of errors per each service
    - Internet and DNS status
    - Kernel version
    - Magma version

#### 3.2.5 HA Active Active

GTP Gateway will be deployed in HA Active Active configuration. GTP Gateway instances do not need to be colocated. Session loading will be balanced between the available instances.

#### 3.2.6 Infrastructure

GTP Gateway will be deployed on Equinix Metal hosts, or other bare metal hosts. Equinix metal C3.Small will be used to start. These hosts support:

- PROC1 x Xeon E 2278G 8-Core Processor @ 3.4Ghz
- RAM 32GB
- 2 x 480GB SSD
- NIC2 x 10Gbps Bonded Ports

This specification will be considered the base supported specification until further performance data is collected.

NOTE: Equinix Metal chosen for it’s access to IPX/GRX network connections.

#### 3.2.7 NMS GTP Gateway

NMS Elements will be created for display of key configuration and health elements including:

- View configuration
- SW Version
- Gateway - Controller connectivity
- Interface IP addresses
- View health metrics
- Number of active VPN connections
- Number of active GTP connections
- Status of PGW VPN (if present)
- Total Ingress throughput per NIC
- Total Egress throughput per NIC

## 4. Roadmap & Schedule

**MS1:** PoC Implementation via Equinix Dallas DC connected to GigSky PGW using Test eSIM

**MS2:** Integration with Magma Orc8r including health metrics.

**MS3:** HA Active Active Support Complete

**MS4:** NMS Integration Complete

Schedule: TBD
