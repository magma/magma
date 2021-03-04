---
id: intro
title: Introduction
hide_title: true
---
# Introduction

*This document serves as an introduction to Orchestrator functionality and
architecture. Operators may find this useful to better understand their
deployment.*

## Overview

Magma’s Orchestrator is a centralized controller for a set of networks.
Orchestrator handles the control plane for various types of gateways in Magma.
Orchestrator functionality is composed of two primary components:

* A standardized, vendor-agnostic **northbound REST API** which exposes
configuration and metrics for network devices
* **A southbound interface** which applies device configuration and reports
device status

![Orc8r Overview](assets/orc8r/orc8r_overview.png)

At a lower level, Orchestrator supports the following functionality:

* Network entity configuration (i.e. networks, gateway, subscribers, policies, etc.)
* Metrics querying via Prometheus and Grafana
* Event and log aggregation via Fluentd and Elasticsearch Kibana
* Config streaming for gateways, subscribers, policies, etc.
* Device state reporting (metrics + status)
* Request relaying between access gateways and federated gateways

## Architecture

Orchestrator follows a modular design. Core Orchestrator services
(located at `magma/orc8r/cloud/go`) provide domain agnostic implementations of
the functionality described above. Orchestrator modules, such as `lte`
(located at `magma/lte/cloud/go`) provide domain specific knowledge to the
core services. This modularity allows Orchestrator to remain flexible to new
use-cases while allowing operators to only deploy the services needed for their
specific use-case.

### Controller Supercontainer (v1.0-v1.3)

Prior to v1.4, the Orchestrator application was deployed in a single
“supercontainer” on Kubernetes.

![Orc8r Plugin Architecture](assets/orc8r/orc8r_plugin.png)
Modularity was achieved via the following plugin interface:

```
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
Plugins were built into the controller Docker image and loaded into memory for
each service at runtime.

### Service Mesh (v1.4+)

Starting in v1.4, Orchestrator is deployed as a service mesh. Every
Orchestrator service is now deployed with it’s own Kubernetes service and pod.

![Orc8r Service Mesh Architecture](assets/orc8r/orc8r_service_mesh.png)

The service mesh changes provide the following benefits:

* Major memory reduction from the removal of the Orchestrator plugins
* Better monitoring via open source tooling
* Ability to horizontally scale specific services

To achieve the same modularity described above, the in-memory plugin has been
decomposed into strictly RPC interfaces. Individual services in an Orchestrator
module register labels and annotations in its Kubernetes service definition to
declare what functionality the service provides. Core Orchestrator services
then query the Kubernetes API to discover service functionality.

### Configuration

The following diagram displays an example call flow for LTE gateway
configuration.

![Orc8r Configuration](assets/orc8r/orc8r_configuration.png)

### Metrics

The following diagram displays the call flow for the metrics pipeline in
Orchestrator.

![Orc8r Metrics](assets/orc8r/orc8r_metrics.png)

### Services

The following section outlines the functionality of each service in the
core Orchestrator, LTE, and FEG modules.

**Orchestrator**


* **Accessd**: Stores, manages and verifies operator identity objects and their
rights to access (read/write) Entities.
* **Analytics:** Periodically fetches and aggregates metrics for all deployed
Orchestrator modules, exporting the aggregations to Prometheus.
* **Bootstrapper:** Manages the certificate bootstrapping process for newly
registered gateways and gateways whose cert has expired.
* **Certifier:** Maintains and verifies signed client certificates and their
associated identities.
* **Configurator:** Maintains configurations and meta data for networks and
network entity structures.
* **Ctraced:** Handles gateway call traces, exposing this functionality via a
CRUD API.
* **Device:** Maintains configurations and meta data for devices in the network
(i.e gateways).
* **Directoryd:** Stores subscriber identity (i.e. IMSI, IP address,
MAC address) and location (i.e. gateway hardwareID).
* **Dispatcher:** Maintains SyncRPC connections (HTTP2 bidirectional streams)
with the gateways.
* **Metricsd:** Collects runtime metrics from the gateways and Orchestrator
services.
* **Obsidian**: Verifies API request access control and reverse proxies
requests to the Orchestrator service with the appropriate API handlers.
* **Orchestrator:** Provides:
    * Mconfig for Orchestrator related configuration (i.e. magmad, eventd, etc.)
    * Metrics exporting to Prometheus
    * Mconfig stream updates
    * CRUD API for core Orchestrator network entities (Networks, Gateways,
    Upgrade Tiers, Events, etc.)
* **Service Registry:** Provides service discovery for all services in the
Orchestrator by querying Kubernetes’ API server.
* **State:** Maintains reported state from devices in the network.
* **Streamer:** Fetches updates for various data streams (i.e. mconfig,
subscribers, etc.) from the appropriate Orchestrator service, returning these
to the gateways.
* **Tenants:** Provides CRUD interface for managing NMS tenants.

**LTE**

* **HA**: Provides interface for secondary gateways in an HA deployment to find
offload status for UEs.
* **LTE**: Provides:
    * Mconfig for LTE related configuration (i.e. MME, policydb, pipelined, etc.)
    * CRUD API for LTE network entities (LTE Networks, LTE Gateways, EnodeBs, etc.)
* **Policydb:** Manages subscriber policies via a northbound CRUD API and
a southbound policy stream.
* **SMSd:** Provides CRUD support for SMS messages to be fetched by LTE gateways.
* **Subscriberdb:** Manages subscribers via a northbound CRUD API and
a southbound subscriber stream.

**FEG**

* **FEG**: Provides
    * Mconfig for LTE related configuration (i.e. s6a_proxy, session_proxy, etc.)
    * CRUD API for LTE network entities (FEG Networks, Federated Gateways, etc.)
* **Feg Relay:** Relays requests between access gateways and federated gateways
* **Health:** Manages active/standby clusters of federated gateways

