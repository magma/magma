---
id: version-1.4.0-p006_subscriber_state_view
title: Displaying Run-time Subscriber State in NMS
hide_title: true
original_id: p006_subscriber_state_view
---
*Feature owner: @themarwhal*

# Subscriber State Reporting From GW → Orc8r

## Summary of currently available run-time subscriber state

Currently, our view of subscribers is fairly fragmented over states reported by various services and event logs.

### DirectoryD (IMSI → SessionID/MSISDN)

* On session creation, IMSI to SessionID/MSISDN is added as an entry by SessionD. However, this is not ideal as SessionID is only unique per session not IMSI. This means that in a multi APN scenario, the first entry gets overwritten when the second session is initialized.
* For CWF only, we also have an additional entry in DirectoryD that maps from IMSI→UE IPv4 added by PipelineD.

### MobilityD (IMSI+APN → IPv4)

* MobilityD also reports internal state to the orc8r from which we parse out IMSI+APN ↔ IPv4 mapping



## Summary of currently available subscriber/policy config

For non-Federated deployments, all subscriber and subscriber->policy configurations are made on the Orc8r. So the Orc8r at least has a sense of the desired state.
For federated deployments where the subscriber configurations are federated to HSS or the policy configurations are federated to PCRF, the Orc8r does not have any visibility into the configurations received.


## General Goal

The goal here is to export consolidated session state owned by SessionD to the Orc8r so that we can 

1. View all existing session states, uniquely identified by IMSI+APN 
2. View runtime state of active policies + QoS



## Proposed State from SessionD → GW State service

Here is the state Protobuf type:

```
message State {
    // Type determines how the value is deserialized and validated on the cloud service side
    string type = 1; // 'subscriber_state'
    string deviceID = 2; // IMSI'
    // Value contains the operational state json-serialized.
    bytes value = 3; // JSON serialized Subscriber State type below
    uint64 version = 4; // ?
}
```

Here is the blob value:

```
  subscriber_state:
    type: object
    description: An object that maps IMSI->APN->SubscriberState
    properties:
      imsi:
        type: string
      sessions:
        type: array
        items:
          $ref: '#/definitions/policy_rule'
  session_state:
    type: object
    description: Used to describe a session's state per APN
    additionalProperties:
      type: object
      properties:
        apn:
          type: string
        msisdn:
          type: string
        session_id:
          type: string
        ipv4:
          type: string
        session_start_time:
          type: uint32
        lifecycle_state: 
          type: string
        active_duration_sec:
          type: uint32
        active_policy_rules:
          type: array
          items:
            $ref:# A JSON-serialized version of protobuf PolicyRule
```

### Notes on the subscriber_state fields

* IMSI+APN uniquely identifies each active session
    * In SessionD there could be situations where we multiple sessions for the 
    same IMSI+APN at the same time, but there is only one ACTIVE session at a time. 
    Additionally, we do not keep around non-ACTIVE sessions for longer than 5 seconds typically. 
* MSISDN value is propagated from MME to SessionD on session creation. The value is fetched from HSS/SubscriberDB. 
* This is the first round of fields suggested, but the following fields could be useful for debugging as well
    * Attach time & session duration in seconds (Update: Included now)

**[P1]** As with other gateway states reported by the state service, SessionD’s state will also be sent as an untyped JSON blob. These subscriber states will be reported from the Access network. (LTE/FEG_LTE/CWF etc.)

## Proposed Orc8r API Endpoints

**[P1]** Define read-only `lte` subscriber state endpoints that dumps the most recently reported subscriber state. 

**[P2]** We can expose a similar endpoint for other access network types, `cwf`. 

```
/lte/{network_id}/subscriber_state/:
    get:
      summary: A list of subscriber states in this network
      tags:
        - Federated LTE Networks
      parameters:
        - $ref: './orc8r-swagger-common.yml#/parameters/network_id'
      responses:
        '200':
          description: Subscriber State
          schema:
            type: array
            items:
              $ref: '#/definitions/subscriber_state'
        default:
          $ref: './orc8r-swagger-common.yml#/responses/UnexpectedError'

/lte/{network_id}/subscriber_state/{imsi}:
    get:
        summary: A list of subscriber states in this network keyed by IMSI
        tags:
        - Federated LTE Networks
        parameters:
        - $ref: './orc8r-swagger-common.yml#/parameters/network_id'
        - $ref:  './lte-policydb-swagger.yml#/parameters/imsi'
        responses:
        '200':
            description: Subscriber State
            schema:
            type: array
            items:
                $ref: '#/definitions/subscriber_state'
        default:
            $ref: './orc8r-swagger-common.yml#/responses/UnexpectedError'
```

## Proposed NMS Changes

**[P1]** For this limited scope, subscriber_state data can be displayed as part of the FederatedLTE network page. The NMS will use the `lte` endpoint defined above to display this information.

## Timelines (Includes Testing Time)

### Gateway Items (Owned by Marie Bremner)

* **[P1]** Define ReportedStates response handler for SessionD and dump state in the response. (~3 days)

### Orc8r Items (Owned by Hunter Gatewood)

* **[P1]** Define API handlers for LTE network to display JSON blob  (~1 week)
* **[P2]** Work with Marie to identify any fields that will need to be used for purposes other than displaying

### NMS Items (Owned by Karthik Shyam Krishnan Subraveti)

* **[P1]** Add a viewing page for subscriber state in FeGLTE network (~1-2weeks)
