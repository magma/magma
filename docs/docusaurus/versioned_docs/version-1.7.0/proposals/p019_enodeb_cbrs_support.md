---
id: version-1.7.0-p019_enodeb_cbrs_support
title: Enodebd CBRS Support
hide_title: true
original_id: p019_enodeb_cbrs_support
---

# Proposal: [Enodebd CBRS Support]

Author(s): [@amarpad]

Last updated: [07/20/2021]

Discussion at <https://github.com/magma/magma/issues/8196>.

## Context & scope

CBRS radios require specific radio parameters to be set based on grants that
can be obtained from third party spectrum databases called SAS. In Magma this
requires communication between the enodebd service that manages the radio
configuration and the domain proxy that runs in the orchestrator that acts as
a single interface point towards the SAS on behalf of all radios on the network

There are many advantages to having communication go through the proxy as
opposed to directly from the enodeb including security. Further, having the
desired state of the network including the radio parameters managed in the
orchestrator has been a primary goal of Magma.

## Goal

- Allow enodebs to obtain spectrum leases from SAS on bring up  
- Allow enodebs to heartbeat to the SAS and if needed the SAS can revoke the
  grant  
- Deregister a cbrs eNB from the SAS when it is deleted from magma  
- Relinquish a grant if the eNB stops radiating  

## Non-goals / simplifications

- For simplicity we collapse the need for a SAS grant with the transmit enabled
  status of the radio. i.e. if a CBRS radio is transmitting we assume it needs
  a grant.

## Proposal

The top level goal of this proposal is to have the domain proxy on the
orchestrator to serve as the truth and act as a desired state store for the
enodeb and rely on a control channel between enodebd->domain proxy and
enodeb->enodebd to keep the actual state in sync with the desired state.

### User facing

- NMS: NMS will explicitly support specific CBRS certified radios from vendors.
  The knowledge of these radios being CBRS certified is obtained out of band.
  i.e. for practical purposes the onboarding of these radios are going to be
  similar to what we have today.

### API changes

- Swagger: Swagger will introduce a new vendor name for CBRS certified radios

### Mconfig changes

- mconfig will be extended to introduce a new boolean field indicating if the
  radio is a CBRS radio.

### Enodebd changes (only applies to CBRS enodebs)

There are two subcases for enodeb bring up, eNB provisioned for first time/after
factory reset. eNB reconnecting to enodebd after reboot/disconnect etc.

#### New enodebd provisioning

enodebd is in the new enodeb provisioning flow if there is no SAS related
state for the eNB in Redis.

#### Towards eNB

Precondition: SAS state in Redis and in-memory is none.

- Parse Inform message as today, if serial number maps to enodeb with CBRS
  enabled and gps location absent, transition to enabling GPS state (new state),
  else complete session state (same complete session as today).
- Enable GPS using model specific tr-69 node, disable transmit and transition
  to complete session state
- Read GPS location, serial number as part of next inform, complete session.
- Radio status here is GPS is enabled, but radio is not transmitting

#### Towards cloud

Precondition: GPS state, serial number available.

- Make GRPC call to domain proxy for registration request with the following
  fields populated.

  ```json
  "cbsdSerialNumber": "abcd1234",
  "installationParam": {
    "latitude": 37.425056,
    "longitude": -122.084113,
  }
  ```

- The domain proxy will persist this state in the database and will interact
  with the SAS on behalf of the enodeb using these parameters.
- On response from the SAS the domain proxy will persist the grant in the
  database.
- Enodebd will periodically poll the domain proxy to retrive the desired
  state of the radio, once the SAS grant has been recvd the domain proxy
  will respond with a payload like

  ```json
  "Tx enabled": yes
  "Tx Power" : xxxx
  "Frequency" : yyyy
  ```

  this poll frequency is configurable.
- Version 1: Restart enodebd, when enodebd comes back up it will be in the
  reboot/reconnect enodeb flow.

#### Reboot/Reconnect endoebd provisioning

- Populate in-memory tx power, frequency and other SAS parameter states that
  specific vendors require on startup.
- On Inform message if serial number maps to CBRS enabled enodeb, set cbrs
  params and transmit + power like today. Enable transmit.
- enodebd will continue to poll the domain proxy as above for the transmit
  payload.

#### Hearbeat and spectrum revocation

The domain proxy will heartbeat with the SAS on behalf of all enodebs
provisioned at a frequency specified by the SAS.

- On spectrum revocation the domain proxy will update the database state
  associated with the eNB to disable transmit.
- enodebd's periodic poll request will return the following

  ```json
  "Tx enabled": false
  ```

- enodebd will remove the sas redis state associated with the radio and restart.
- enodebd will move into new enodebd provisioning flow.

#### Handling radio relocation

The poll request from the enodebd can contain the lat/long associated with the
enodeb, if the location changes, the domain proxy can reset the Tx enabled flag
which will cause the enodeb to stop transmitting and restart provisioning based
on the new lat long. The Domain proxy inturn will need to explicitly
relinquish the grant from the SAS using the "Relenquish" method.

## Alternatives considered

- Leverage sync RPC or bi-directional GRPC for communication from the DP to
  the enodebd instead of a poll loop: We deprioritized this option to use
  poll for the following reasons:
    - The domain proxy instance that receives the heartbeat revoking the
    spectrum grant might still need to persist the state in case the enodebd
    is unreachable. Further modeling the event as state is more Magmaesque.
    Poll is a simple way of state synchronization using a set interface
    - Sync RPC and GRPC based streaming are more complicated.

- Not restarting the enodebd service when SAS gets configured: This is just
  a phase-1 implementation, further as the interface towards the enodebd is
  idempotent restarting the enodebd service should have no effect on traffic.

## Cross-cutting concerns

- SAS interactions are message oriented, by persisting the desired enodeb
  state we model this more as a level trigger. This comes at the expense of
  responsiveness (capped by poll frequency), but this is an explicit tradeoff.

## Compatibility

N/A

## Observability and Debug

enodebd_cli should be enhanced to show SAS related information

## Security & privacy

Reuse existing communication channels

## Open issues (if applicable)

- Optimize away the need to restart enodebd on SAS param changes
