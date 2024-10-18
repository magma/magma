---
id: version-1.7.0-sim_integration
title: SIM Management System Integration
hide_title: true
original_id: sim_integration
---

# Proposal: SIM Management System Integration

Author(s): andreilee

Last updated: August 27, 2021

## Context & scope

SIM platform integration is both a commonly requested integration capability, and also necessary feature.

The use case that will be focused on here is providing LTE access on a Magma network,
and allowing subscribers to have access when leaving the Magma network.
One way to provide this is outbound roaming,
and another way to provide this is with integration with a programmable SIM platform,
such as that provided by [Teal Communications](https://www.tealcom.io/).
Teal provides programmable SIMs that can provide connectivity globally
(or just locally if specified),
and have the home network be specified as the Magma network.

Another use for this is for CBRS **to be filled in**

### Goals

- eSIM platform integration with Teal Communications
- usage of eSIM platform should be optional for an operator using Magma
- all SIM management should be able to be done through the Magma NMS
- users of Magma should not need to log into 3rd party platforms like Teal's, there should be integration through NMS

## Proposal

### Overview

Magma will add optional integration with the eSIM platform provided by [Teal Communications](https://www.tealcom.io/).
Network operators running Magma will not need to

### NMS UI Changes

- SIM management as a tab under subscribers
    - as passthrough to underlying SIM management system being integrated with
-
- need a page to enable or disable SIM management

### API Changes

- Endpoint to upload file with auth opc and key to be correlated with UICCs
- Endpoint to update API keys for accessing SIM management system
- Subscribers resource changed
- Two types of subscribers:
    - Type 1: Magma managed SIM
    - Type 2: externally managed SIM (Teal Systems only for now)
- new UICC resource
- polymorphism used so subscribers endpoints can return both types of subscribers
- See [swagger spec](https://swagger.io/specification/v2/#schemaObject) - polymorphism
- `/magma/v1/lte/{network_id}/uicc`

### New Orc8r Service uicc for SIM management

- Orc8r-only service, in LTE directory
- optional service
- gRPC interface
    - upload file with auth opc and key to be correlated with SIMs
    - update API keys

### Orc8r Subscriberdb Changes

- Edit Subscriberdb -> subscribers are split to have SIM entries
    - subscribers can either have associated SIM EIDs from a separate managed system, or be Magma managed
    - orc8r subscriberdb should still stream the same updates to AGW

### AGW Subscriberdb Changes

- No changes, as the same data is still streamed from the Orc8r

## Alternatives considered

The main problem being addressed with this proposal is on how to provide
subscribers access when they leave the Magma network.
The primary alternative considered was implementation of outbound roaming,
which is more costly from an engineering perspective.

### Compatibility

[A discussion of the change with regard to backward / forward compatibility.]

- For any swagger definition changes, a database migration will be written
- No backwards compatibility issues with AGW, this should affect Orc8r/NMS only

### Observability and Debug

[A description, of how issues with this design would be observed and debugged
in various stages from development through production.]

Implementation will occur in several phases, so verification will occur with
each phase.

#### Phase 1 - Orc8r Subscriberdb Modification

Subscribers will be modified to either be associated with an externally managed
eSIM, or a Magma managed SIM.

#### Phase 2 - Orc8r Service uicc

Orc8r service `uicc` will be added, and API changes will be added.
Subscriberdb will be changed to interact with `uicc`.
In this phase, `uicc` will interact with Teal Communications eSIM management
platform.

#### Phase 3 - NMS Integration

All SIM management functionality will be provided on NMS.
This functionality should have already been present by the end of phase 2,
as the NMS should just use the Magma API provided by Orc8r.

### Security & privacy

No changes
