---
id: architecture_overview
title: Overview
hide_title: true
---

# Architecture Overview

This document serves as an introduction to Orchestrator functionality and
architecture. Operators may find this useful to better understand their
deployment.

## Functionality

Magma’s Orchestrator is a centralized controller for a set of networks.
Orchestrator handles the control plane for various types of gateways in Magma.
Orchestrator functionality is composed of two primary components

- A standardized, vendor-agnostic [**northbound REST API**](https://app.swaggerhub.com/apis/MagmaCore/Magma/1.0.0) which exposes
configuration and metrics for network devices
- **A southbound interface** which applies device configuration and reports
device status

![Orc8r Overview](assets/orc8r/orc8r_overview.png)

At a lower level, Orchestrator supports the following functionality

- Network entity configuration (networks, gateway, subscribers, policies, etc.)
- Metrics querying via Prometheus and Grafana
- Event and log aggregation via Fluentd and Elasticsearch Kibana
- Config streaming for gateways, subscribers, policies, etc.
- Device state reporting (metrics and status)
- Request relaying between access gateways and federated gateways

## Architecture

Orchestrator follows a modular design. Core Orchestrator services
(located at `orc8r/cloud/go`) provide domain-agnostic implementations of
the functionality described above. Orchestrator modules, such as `lte`
(located at `lte/cloud/go`), provide domain-specific knowledge to the
core services. This modularity allows Orchestrator to remain flexible to new
use-cases while allowing operators to only deploy the services needed for their
specific use-case.

### Controller Supercontainer (pre-v1.3)

Prior to v1.4, the Orchestrator application was deployed in a single
“supercontainer” on Kubernetes.

![Orc8r Plugin Architecture](assets/orc8r/orc8r_plugin.png)

Modularity was achieved via the following plugin interface

```go
type OrchestratorPlugin interface {
  GetName()
  GetServices()
  GetSerdes()
  GetMconfigBuilders()
  GetMetricsProfiles()
  GetObsidianHandlers()
  GetStreamerProviders()
  GetStateIndexers()
}
```

Each Orchestrator module (e.g. `lte`, `feg`, `cwf`) implemented this interface.
Plugins were built into the controller container image and loaded into memory for
each service at runtime.

### Service Mesh (v1.4+)

Starting in v1.4, Orchestrator is deployed as a service mesh. Every
Orchestrator service is now deployed with its own Kubernetes service and pod.

![Orc8r Service Mesh Architecture](assets/orc8r/orc8r_service_mesh.png)

The service mesh changes provide the following benefits

- Major memory reduction from the removal of the Orchestrator plugins
- Better monitoring via open source tooling
- Ability to horizontally scale specific services

To achieve the same modularity described above, the in-memory plugin has been
decomposed into strictly RPC interfaces. Individual services in an Orchestrator
module register labels and annotations in its Kubernetes service definition to
declare what functionality the service provides. Core Orchestrator services
then query the Kubernetes API to discover service functionality.

### Configuration

The following diagram displays an example call flow for LTE gateway
configuration

![Orc8r Configuration](assets/orc8r/orc8r_configuration.png)

### Metrics

The following diagram displays the call flow for the metrics pipeline

![Orc8r Metrics](assets/orc8r/orc8r_metrics.png)

### Services

This section outlines the functionality of each service in the core
Orchestrator, LTE, and FeG modules.

**Orchestrator**

- *accessd* stores, manages and verifies operator identity objects and their
rights to access (read/write) entities
- *analytics* periodically fetches and aggregates metrics for all deployed
Orchestrator modules, exporting the aggregations to Prometheus
- *bootstrapper* manages the certificate bootstrapping process for newly
registered gateways and gateways whose cert has expired
- *certifier* maintains and verifies signed client certificates and their
associated identities
- *configurator* maintains configurations and metadata for networks and
network entity structures
- *ctraced* handles gateway call traces, exposing this functionality via a
CRUD API
- *device* maintains configurations and metadata for devices in the network
(e.g. gateways)
- *directoryd* stores subscriber identity (e.g. IMSI, IP address,
MAC address) and location (gateway hardware ID)
- *dispatcher* maintains SyncRPC connections (HTTP2 bidirectional streams)
with gateways
- *metricsd* collects runtime metrics from gateways and Orchestrator
services
- *obsidian* verifies API request access control and reverse proxies
requests to Orchestrator services with the appropriate API handlers
- *orchestrator* provides
    - Mconfigs for configuration of core gateway service configurations
     (e.g. magmad, eventd, state)
    - Metrics exporting to Prometheus
    - CRUD API for core Orchestrator network entities (networks, gateways,
    upgrade tiers, events, etc.)
- *service_registry* provides service discovery for all services in the
Orchestrator by querying Kubernetes's API server
- *state* maintains reported state from devices in the network
- *streamer* fetches updates for various data streams (e.g. mconfig,
subscribers, etc.) from the appropriate Orchestrator service, returning these
to the gateways
- *tenants* provides CRUD interface for managing NMS tenants

**LTE**

- *ha* provides interface for secondary gateways in an HA deployment to find
offload status for UEs
- *lte* provides
    - Mconfigs for configuration of LTE-related gateway service configurations
      (e.g. mme, pipelined, policydb)
    - CRUD API for LTE network entities (LTE networks, LTE gateways, eNodeBs, etc.)
- *policydb* manages subscriber policies via a northbound CRUD API and
a southbound policy stream
- *smsd* provides CRUD support for SMS messages to be fetched by LTE gateways
- *subscriberdb* manages subscribers via a northbound CRUD API and
a southbound subscriber stream

**FeG**

- *feg* provides
    - Mconfigs for configuration of FeG-related gateway service configurations
      (e.g. s6a_proxy, session_proxy)
    - CRUD API for LTE network entities (FeG networks, federated gateways, etc.)
- *feg relay* relays requests between access gateways and federated gateways
- *health* manages active/standby clusters of federated gateways
