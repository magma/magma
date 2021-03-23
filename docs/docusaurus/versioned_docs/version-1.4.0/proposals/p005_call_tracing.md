---
id: version-1.4.0-p005_call_tracing
title: Call Tracing for Troubleshooting
hide_title: true
original_id: p005_call_tracing
---

# Tracing Support

*Status: In review*
*Feature owner: @andreilee*
*Feedback requested from: @amarpad, @xjtian, @karthiksubraveti*
*Last Updated: 10/21*

## Summary

As a general purpose debugging tool, a tracing tool will be added for
understanding control messaging flow through a Magma access gateway.

## Motivation

Example use cases include the following:

**Subscriber Tracing**

When issues are experienced by a subscriber on the network, network operators
would like to have a tool which allows them to see the messages between the UE,
eNodeB, and Magma components relating to the subscriber. This feature already
exists for other NMS tools. Ideally, this capture should be a packet capture to
fit with existing workflows using Wireshark. Wireshark has features to decode
and display relevant messages for Gx, Gy, S1AP, TR-069, etc., as well as
plugins to support decoding and viewing gRPC messages which are passed between
Magma gateway components.

**Protocol Tracing**

Operators would like to be able to capture traces for specific protocol stacks,
such as SCTP or DIAMETER. This could be used for troubleshooting specific
integrations.

## Goals

- Allow network operators to trace the control plane messaging for a subscriber
- Allow network operators to trace control plane messaging for a specific
protocol or interface
- Allow call trace captures to started and stopped from the NMS 
- Allow call trace captures to be filtered by subscriber/protocol/interface
control plane messaging

As an additional criteria which is not necessary, call trace captures should
be viewable in Wireshark.

## Implementation Phases

The implementation will be done in two phases.

**Phase 1 - MVP**

(October - November 27)

The first phase of implementation will focus on bringing an MVP. Call traces
will be able to be initiated from the NMS, and downloads of captures viewable
through Wireshark after capture will be available.

**Phase 2 - Call Trace Filtering**

(December - January)

In the second phase, features will be provided for filtering call traces by
subscriber, protocol, or interface.

## Phase 1

#### User Flow

For call tracing with Magma, we will provide the following user flow:
1. User accesses NMS
2. User selects either the gateway or subscriber that they would like to do
call tracing for
3. User starts the call trace
4. After a timeout, or by user control, the call trace stops
5. The call trace will be provided as a download to the user

### Design

#### Orchestrator - API

The API supports starting the call trace, ending the tracing session, and
downloading the resulting capture.

The API will be as follows:
```
GET      /networks/{network_id}/tracing/{trace_type}/{trace_id}
POST     /networks/{network_id}/tracing/{trace_type}/{trace_id}
PUT      /networks/{network_id}/tracing/{trace_type}/{trace_id}
DELETE   /networks/{network_id}/tracing/{trace_type}/{trace_id}

GET      /networks/{network_id}/tracing/{trace_type}/download/{trace_id}
```

To use call tracing via the REST API, the following steps should be followed.

To begin a call trace, the POST endpoint should be called to create a
`call_trace` resource. During this time, the type of call trace should be
specified, along with timeout, and the `trace_id`.

To end a call trace, either the user can wait for the timeout, or to manually
end the call trace, the PUT endpoint can be called, specifying `requested_end`
as `true`.

To download the call trace, the GET endpoint should return whether the call
trace has ended and is ready to download, and if it is, will return the
download URL.

```
definitions:
  call_trace:
   description: Mutable Call Trace
   type: object
   required:
     - trace_id
     - config
   properties:
     trace_id:
       type: string
       x-nullable: false
       example: "pair1"
     config:
       $ref: call_trace_config
     state:
       $ref: call_trace_state

  mutable_call_trace:
   description: Mutable Call Trace
   type: object
   required:
     - trace_id
     - config
   properties:
     trace_id:
       type: string
       x-nullable: false
       example: "pair1"
     requested_end:
       type: boolean
     config:
       $ref: call_trace_config

  call_trace_config:
    type: object
    description: Call Trace spec
    required:
      - type
    properties:
      trace_type:
        type: string
        x-nullable: false
        enum:
          - 'GATEWAY'
          - 'SUBSCRIBER'
          - 'PROTOCOL'
          - 'INTERFACE'
      imsi:
        type: string
      protocol:
        type: string
      interface:
        type: string
      timeout:
        type: integer
        format: uint32
        description: Timeout of call trace in seconds

  call_trace state:
    type: object
    description: Full state object of a call trace
    properties:
      call_trace_available:
        type: boolean
```

#### Orchestrator - ctraced

On the gateway, the ctraced service will manage call tracing for the network.
This service will be responsible for initiating call tracing sessions on
gateways, and servicing downloads of call trace captures.

Call traces will be stored in postgresql, and periodically, ctraced should
delete old traces. Traces will expire after a week, and so a cronjob will be
responsible for deleting old traces. This setting will be configurable.

#### Orchestrator - Trace storage

Storage of call trace captures will be done in postgresql, stored as binary
data. A new table will be created for storage of call traces.

A migration script is required to ensure that existing deployments will get
this new table for storing call traces.

#### Gateway - ctraced

The ctraced service on the gateway will be a new service, responsible for
managing call tracing. Options will be provided for gateway wide
packet-capture, and later on, options will be provided for protocol, interface,
or subscriber specific call tracing.

```
service CallTraceService {
  rpc StartCallTrace (StartTraceRequest) returns (StartTraceResponse) {}

  rpc EndCallTrace (EndTraceRequest) returns (EndTraceResponse) {}
}


message StartTraceRequest {
  enum TraceType {
    ALL = 0;
    SUBSCRIBER = 1;
    PROTOCOL = 2;
    INTERFACE = 3;
  }

  TraceType trace_type = 1;
  // IMSI specified only if trace_type is SUBSCRIBER
  string imsi = 2; // Include prefix 'IMSI'

  enum ProtocolName {
    SCTP = 0;
    DIAMETER = 1;
  }
  // Protocol name specified only if trace_type is PROTOCOL
  ProtocolName protocol = 3;
  enum InterfaceName {
    S1AP = 0;
    GX = 1;
    GT = 2;
  }
  // Interface name specified only if trace_type is INTERFACE
  InterfaceName interface = 4;
}

message StartTraceResponse {
  bool success = 1; // May fail due to an existing tracing session
}

message EndTraceRequest {}

message EndTraceResponse {
  bool success = 1; // May fail due to no existing tracing session
  bytes trace_content = 1; // Max size of 4MB
}
```

To perform call tracing, ctraced will use `tshark`, a terminal-based wireshark
tool to do packet capture on the gateway. This will be installed as a
provisioning step on the gateway.

Call tracing will be initiated from the orc8r, and the packet-capture will be
transferred back to orc8r from AGW via gRPC. Currently there is a limit for
gRPC message sizes to 4MB, which limits the size of packet capture. Given that
the packet capture should filter out data plane traffic, this should be
sufficient for short tracing sessions, for an MVP.

#### NMS - Interface

A new section will be added to the NMS to initiate call traces,
see past call traces, and download them.

---

## Phase 2

#### User Flow

Phase 2 user flow will be the same from the user interface. The differentiation
will be in the call trace provided to the operator, which will have its tracing
information filtered by the specified subscriber.

1. User accesses NMS
2. User options of the trace: protocol, interface, and subscriber/gateway
3. User starts the call trace
4. After a timeout, or by user control, the call trace stops
5. The call trace will be provided as a download to the user

#### Gateway - ctraced

The ctraced service will need to be modified to filter messages. For protocols 
including S1-AP or interfaces like Gx, tshark will be used to filter for the
relevant subscriber, protocol, or interface.

For gRPC messages, tshark does not provide features to filter.
Separate gRPC tracing will need to be implemented, and from this separate gRPC
trace, ctraced will filter for relevant messages. Together with the packet
capture acquired through tshark, these two traces will be provided for download
to the user.

#### Tracing Options

Tracing will be able to be filtered for the following interfaces:
- S6a
- Gx
- Gy

Or for the following protocols:
- S1-AP
- Diameter Credit-Control Application
- TR-069




