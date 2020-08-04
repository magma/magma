# Session Management in Magma

## Overview of SessionD
SessionD is the main service responsible for managing and enforcing session 
configurations. 

### Interfaces
SessionD interfaces with three distinct external components. 
1. Access (MME, AAA) <br>
 This component receives all connection and session related information from 
 the UE. It also handles mobility management. But from SessionD's point of view,
 it notifies SessionD whenever a UE is attached or detached. This notification 
 will include any context that SessionD needs for enforcement.
2. Policy (PolicyDB, FeG+PCRF/OCS) <br>
 This component is responsible for propagating any session or policy 
 configuration that SessionD must enforce. Whenever a new UE is introduced, 
 SessionD will reach out to this component to fetch any configuration. Depending
 on the configuration, SessionD maybe reach out to this component for updates. 
3. User Plane (PipelineD) <br>
 This component is responsible for policy and QoS enforcement. It receives 
 any relevant policy and QoS configuration from SessionD and periodically 
 reports usage accordingly.
 
### Internal Architecture 
![Internal Architecture Diagram](SessionD_Architecture.png)

## Configurations (Not Exhaustive)
- FeG Relay <br>
  A federation gateway, FeG, is an entity responsible for hosting a centralized 
  control plane interface towards an operator core. In terms of session 
  management, the SessionProxy service on a FeG will relay all session related 
  interactions. <br>
  This feature can be enabled on the Orc8r with a network wide `relay_enabled` 
  flag. When it is disabled, a local PolicyDB service is used in place. This
  service will directly fetch session configurations from the Orc8r. <br> 
  The feature is disabled by default.
- Omnipresent Rules <br>
  *This feature is currently only relevant for the federated case.* <br>
  Omnipresent rules are policies that are applied to all subscribers in a 
  network. The list of such policies are configured on the Orc8r as a network 
  configuration and are streamed down to all gateways in the network.
  These policies are added onto the list of policies configured by the PCRF. 
- Zero Wallet Detections <br>
  *This feature is currently only relevant for the CWF case.* <br>
  This configures a way for SessionD to detect when a subscriber is out of valid
  wallet. When an empty wallet is detected, the empty status is propagated to 
  PipelineD, which hosts a flask server that indicates the state. After a 
  configured number of seconds, the session will be terminated and the user will
  be kicked out.
  The configuration is done in Orc8r, but it is not currently exposed by an API.
  This timeout is configured in `sessiond.yml` with 
  `cwf_quota_exhaustion_termination_on_init_ms`.
  <br>
  The detection methods are described below: <br>
  - GxTrackedRules: Wallet is empty if there are no active PCRF tracked policy.
    PCRF tracked policies are policies with tracking type `PCRF_ONLY` or 
    `PCRF_AND_OCS`. 
- `magma/lte/gateway/configs/sessiond.yml` has more configurations that are 
   managed on each gateway.
   
## Stateless Operation
A stateless SessionD service can be enabled by the `support_stateless` flag in 
`sessiond.yml`. This feature allows SessionD to be restarted without requiring 
previously existing sessions to be re-authenticated.

With this feature, SessionD stores all session state updates into Redis for 
persistent storage. When a session states needs to be modified, the state is
loaded from storage and modified in memory. When the modification is complete,
the updates are written back into Redis. When the feature is disabled, all 
states are kept in memory.

`SessionStore` is the interface to storage that makes the decision of where the 
states are kept based on the `support_stateless` flag.
Updates to stored session state are done with the `UpdateCriteria` structures 
passed into `SessionStore`.

To ensure that SessionD can always be restarted properly, session state 
should always be consistent with services with which it interacts. 
On SessionD restart, session state will be synced to MME service, which 
it will use as the source of truth for active sessions. 
Based on the active sessions, PipelineD will be instructed to 
allow/disallow traffic.
