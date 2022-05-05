---
id: architecture_overview
title: Overview
hide_title: true
---

# Overview

Domain Proxy is an application whose purpose is to communicate eNodeBs with SAS, send requests to SAS on behalf of them,
and maintain their desired state.

## Motivation

Vendor neutral Domain Proxy serves as a single point of contact for eNB-SAS communication for private networks using
Magma as the Operator Core.

Without a vendor neutral Domain Proxy an eNB has to request for SAS Spectrums / Grants on its own. This requires the eNB
to speak the same protocol as SAS (JSON over HTTPS), have the intelligence to send specific requests required by SAS,
and adhere to the sequence in which they need to be sent. Radios capable of doing this have their own (usually
internal) **proprietary domain proxy**.

One of the primary goals of Domain Proxy is to keep the validity of the grant obtained by a radio from SAS even in case
of intermittent network failures. This happens thanks to heartbeat requests being sent to SAS on behalf of the radio for
up to 4 hours (by default) of radio's inactivity.

## Implementation

### Backend

DP is a microservice application written in Python and Go. Internally components communicate with each other using **
gRPC** calls. In addition to the application, which is deployed as a separate entity, there is a `DP service` within the
Orchestrator, written in Go, acting as an API backend for the currently developed UI NMS portal.

### Frontend

NMS UI page for Domain Proxy management is currently under development.

## Big Picture

The picture below gives a very general look at the Domain Proxy as a proxy service between radios and SAS

![dp](assets/dp/dp_high_level_overview.png)

### Domain Proxy in Magma's Context

![dp](assets/dp/dp_in_context.png)

If we take a closer look, Domain Proxy as part of Magma is situated next to Orc8r,
(typically on the same Kubernetes cluster) and uses the same database as the Orc8r.

In addition there is a `DP service` within the Orc8r that serves as backend for NMS UI portal (under development)

Both `DP service` and Domain Proxy Application perform CRUD operations on the same tables in the database.

### Operation modes

Domain Proxy was designed to be able to work in two modes: `passive` and `active`.

#### Passive mode (currently deprecated)

The only purpose of `passive` mode is to capture incoming HTTPs requests from eNodeBs, bundle them together by type and
send over to SAS. This scenario however assumes that the radio is capable of sustaining an independent communication
with SAS. Since all the radios currently using Domain Proxy use `TR069` and need to be fully configured, we are now only
supporting `active` mode. The code for `passive` mode has been developed and tested as part of an
HTTPs [Protocol Controller](architecture_overview.md#protocol-controller) but is not a part of DP's deployment.

#### Active mode

Active Mode is currently the **default mode** of operation for Domain Proxy. It assumes that the radio needs to be fully
configured by external means (like eNodeB) and Domain Proxy is to carry out the entire communication with SAS for the
radio. This is achieved mainly by
[Active Mode Controller](architecture_overview.md#active-mode-controller)

### Domain Proxy components

In the picture we see three components that comprise Domain Proxy

![dp](assets/dp/dp_architecture.png)

#### Radio Controller

Written in Python; it is used to store individual serialized requests coming from eNodeBd in the database. A vital part
of `Radio Controller` is the `dp_service` - a GRPC server used by `Enodebd`'s `dp_client` for sending information about
radios and getting back information about their state and transmission parameters once they are authorized in SAS. It is
also used by another Domain Proxy component - `Active Mode Controller` to determine what type of requests should be
stored in the database based off a CBSD state (more on that in the next section).

#### Active Mode Controller

Written in Go; used to monitor and maintain desired state of radios.

Radios that do not have a proprietary domain proxy are unaware neither of the type of messages that need to be sent to
SAS, nor of their sequence. Therefore, a dedicated entity called `Active Mode Config` is stored in the database for each
radio as soon as EnodeBd contacts Domain Proxy with its information.

The `Active Mode Config` is looked at by the `Active Mode Controller` at an interval specified in the config and based
off the current state of the radio and its config, an appropriate request is sent to SAS on behalf of the radio.

E.g. if the radio is not registered yet, a `RegistrationRequest` will be sent. If the radio is marked for deletion,
a `RelinquishmentRequest`
will be issued. If the radio is authorized in SAS a `HeartbeatRequest` will be sent on its behalf at a specified
interval.

#### Configuration Controller

Written in Python; it is responsible for picking up pending requests from the database, bundling them together by type
and source and sending to SAS. Once a response from SAS is received, `Configuration Controller` will split it into
individual responses (if the response was bundled), connect them with corresponding requests and insert in the database.
The related request is also marked as processed.

#### *Protocol Controller

**NOTE**: `Protocol Controller` is a historical name for a component that is no longer an independent part of Domain Proxy, but still
needs to be mentioned. CBSD-SAS `Protocol Controller` was removed by [GitHub PR 12420](https://github.com/magma/magma/pull/12420).

`Protocol Controller` is used to handle incoming requests from an eNB or another domain proxy. It is meant to handle
messages sent using a specific protocol. Currently, the only eNBs connected to DP speak `TR069` protocol (Sercomm
Englewood and Baicells 430 Nova).

Instead, the component treated by Domain Proxy as the only `Protocol Controller` is currently
AGW's [Enodebd](lte/architecture_overview.md#enodebd)
being the configuration component for `TR069` based radios.

`Enodebd` communicates with Domain Proxy's `Radio Controller` as if it were just another DP component, over gRPC, using
a newly introduced component called `dp_client`.

In case if future CBSD models speak different protocols, a new `Protocol Controller` will need to be implemented.
