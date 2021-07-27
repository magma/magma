---
id: ref_sync_rpc
title: Life of a Sync RPC
hide_title: true
---

# Life of a Sync RPC

- [Life of a Sync RPC](#life-of-a-sync-rpc)
	- [Introduction](#introduction)
		- [Purpose of a Sync RPC](#purpose-of-a-sync-rpc)
		- [Magma-Wide Architecture](#magma-wide-architecture)
	- [Details](#details)
		- [The SyncRPCService.SyncRPC Stream](#the-syncrpcservicesyncrpc-stream)
	- [Life of A AGW -> FEG RPC](#life-of-a-agw---feg-rpc)
	- [Related Magma Documentation or Proposals](#related-magma-documentation-or-proposals)

## Introduction

### Purpose of a Sync RPC

The `Sync RPC` service enables any-to-any RPC communications between the Access Gateways (AGW) and Foreign Egress Gateways (FEG) mediated by Magma cloud components.

The cloud mediation is required because gateways may be behind Network Address Translation (NAT). In these settings RPC channels cannot be established directly between Gateways (as no public destination port:ip is yet known / available).  `Sync RPC` represents one solution to this problem by supplying a single streaming bi-direction gRPC service between access gateway and its associated `Magma Orchestrator`, established by the gateway through any NAT (if present). RPCs between Gateways can then multiplex over this stream.

Services relayed over `Sync RPC` are the following (exhaustive list).

- FEG->AGW RPCs defined by the [CSFBGatewayService](../../../feg/protos/csfb.proto)
- AGW->FEG RPCs defined by [CSFBFedGWService](../../../feg/protos/csfb.proto)

The Sync RPC gRPC bi-directional streaming service is defined in [SyncRPCService.SyncRPC](../../../orc8r/protos/sync_rpc_service.proto).


### Magma-Wide Architecture

```
╔═════════════════════════════════════════════════════════╗
║                      Orchestrator                       ║
║ ┌───GatewayRequest                                      ║
║ │┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓  ║
║ │┃Dispatcher                          ┌─────────────┐┣━┓║
║ │┃                ┌───────────┐       │  feg_relay  │┃ ┃║
║ │┃            ┌──▶│Reply chans│       └─────────────┘┃ ┃║
║ │┃            │   └───────────┘ ┌───────────────────┐┃ ┃║
║ │┃  ┌─────────┴─────────┐       │ SyncRPCHttpServer │┃ ┃║
║ └╋─▶│ GatewayRPCBroker  │──────┐└───────────────────┘┃ ┃║
║  ┃  ├───────────────────┘      │                     ┃ ┃║
║  ┃  │SyncRPC                   │SyncRPC              ┃ ┃║
║  ┃  │                          │                     ┃ ┃║
║  ┗━┳╋━━━━━━━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━━━━━━━┛ ┃║
║    ┗╋━━━━━━━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━━━━━━━━━┛║
║ ┏━━━╋━━━━━━━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━┓               ║
║ ┃   │            ngingx        │        ┃               ║
║ ┗━━━╋━━━━━━━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━┛               ║
╚╦━━━━╋━━━━━━━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━╦══════════════╝
 ┃    │              ELB         │         ┃               
 ┗━━━━╋━━━━━━━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━┛               
      │                          │                         
      │                          │                         
╔═════╬═══════════════════╗ ╔════╬════════════════════╗    
║     │     AGW           ║ ║    │      FEG           ║    
║     │                   ║ ║    │                    ║    
║     │SyncRPC            ║ ║    │  SyncRPC           ║    
║  ┌──▼────────────────┐  ║ ║   ┌▼──────────────────┐ ║    
║  │      magmad       │  ║ ║   │      magmad       │ ║    
║  └────┬──────────────┘  ║ ║   └────┬──────────────┘ ║    
║       │gRPC In Custom   ║ ║        │gRPC in Custom  ║    
║       │ HTTP2 Stream    ║ ║        │ HTTP2 Stream   ║    
║       │                 ║ ║        │                ║    
║  ┌────▼──────────────┐  ║ ║   ┌────▼──────────────┐ ║    
║  │   control_proxy   │  ║ ║   │   control_proxy   │ ║    
║  │     (nghttpx)     │  ║ ║   │     (nghttpx)     │ ║    
║  └────┬──────────────┘  ║ ║   └─┬──┬──┬───────────┘ ║    
║       │                 ║ ║     │  │  │             ║    
║  ┌────▼──────────────┐  ║ ║     ▼  ▼  ▼             ║    
║  │        MME        │  ║ ║                         ║    
║  └───────────────────┘  ║ ║                         ║    
╚═════════════════════════╝ ╚═════════════════════════╝    
```

- `dispatcher`
  - Per Gateway {`AGW` or `FEG`} instance differentiated by `gwId`
  - Sends SyncRPCs towards gateways and awaits reply
- `nginx`
  - Load balancer front-end
  - Also responsible for authentication? Via SSL certs?
- `ELB`
  - Amazon's Elastic Load Balancer
  - Stand-in Ingress LB required for all Kubernetes clusters
- `nghttpx`
  - Gateway side termination of authentication?
  - Referenced as `control_proxy` in code and documentation
- `magmad`
  - Gateway side multiplexer between Sync RPC service (to cloud) and local daemons

## Details

### The SyncRPCService.SyncRPC Stream

Communication from the `Orchestrator` towards either `Access Gateways` or `Federated Gateways` is multiplexed over the [SyncRPCService.EstablishSyncRPCStream](../../../orc8r/protos/sync_rpc_service.proto) service.  The endpoints for this service are the `Dispatcher` kube in Orchestrator and the `magmad` instances in the gateways.  The service is defined as follows.

```protobuf
service SyncRPCService {
    // creates a bidirectional stream from gateway to cloud
    // so cloud can send in SyncRPCRequest, and wait for SyncRPCResponse.
    // This will be the underlying service for Synchronous RPC from the cloud.
    rpc EstablishSyncRPCStream (stream SyncRPCResponse) returns (stream SyncRPCRequest) {}
}
```

```protobuf
// SyncRPCRequest is sent down to gateway from cloud
message SyncRPCRequest {
    // unique request Id passed in from cloud
    uint32 reqId = 1;
    GatewayRequest reqBody = 2;
    // cloud will send a heartBeat every minute, if no other requests are sent
    // down to the gateway
    bool heartBeat = 3;
    // connClosed is set to true when the client closes the connection
    bool connClosed = 4;
}

// SyncRPCResponse is sent from gateway to cloud
message SyncRPCResponse {
    uint32 reqId = 1;
    GatewayResponse respBody = 2;
    // gateway will send a heartBeat if it hasn't received SyncRPCRequests from cloud for a while.
    // If it's a heartbeat, reqId and respBody will be ignored.
    bool heartBeat = 3;
}
```

Some notes.

- Each Cloud -> Gateway Request is specified by a `reqId` (though is this true even for `heartBeat` - should this be `OneOf`)?
- Each `SyncRPCRequest` Sent to an AGW results in a reply channel generated to await the `reqID` response within [GatewayRPCBroker](../../../orc8r/cloud/go/services/dispatcher/broker/broker.go).
- When and how is `connClosed` set? How does it ensure reception.

## Life of A AGW -> FEG RPC

Pre-Work:

1. FEG `feg_relay` establishes long-lived [SyncRPCService.SyncRPC](../../../orc8r/protos/sync_rpc_service.proto) stream to Orchestrator.Dispatcher
1. AGW's `magmad` establishes long-lived [SyncRPCService.SyncRPC](../../../orc8r/protos/sync_rpc_service.proto) stream to Orchestrator.Dispatcher

At Time of RPC:

1. `MME` wishes to signal a [CSFBFedGWService.PagingRej](../../../feg/protos/csfb.proto) to an assocaited `FEG` instance
1. `MME` uses [Service discovery](https://sourcegraph.com/github.com/magma/magma/-/blob/orc8r/gateway/python/magma/common/service_registry.py?L165:9) of `magmad` to look up remote endpoint on a [CSFBFedGWService](../../../feg/protos/csfb.proto) request (direct or through `control_proxy`)
2. `Service discovery` returns to the `MME` a `control_proxy` proxied path to `Orchestrator`
3. `MME` sends the [CSFBFedGWService.PagingRej](../../../feg/protos/csfb.proto) to `magmad` via ?? service
4. `magmad` relays the [CSFBFedGWService.PagingRej](../../../feg/protos/csfb.proto) through the `control_proxy` over the long running [SyncRPCService.SyncRPC](../../../orc8r/protos/sync_rpc_service.proto) service
   1. `Magmad` implements [SyncRPCClient](https://sourcegraph.com/github.com/magma/magma/-/blob/orc8r/gateway/python/magma/magmad/sync_rpc_client.py?L31:7&subtree=true#tab=references)
   2. The SyncRPCClient is responsible for certificate validation?
5. `control_proxy` (implemented by `nghttpx`) proxies the gRPC stream to Orchestrator's `ngingx` instance (this configuration is optional)
6. En route to Orchestrator's `ngingx` instance, passes through e.g. `ELB` as the ingress to the Kubernetes cluster
7. `Orchestrator`'s `ngingx` instance relays the SyncRPC request to a `Dispatcher` kube
8. The [GatewayRPCBroker](../../../orc8r/cloud/go/services/dispatcher/broker/broker.go) within the `Dispatcher` service looks up the associated FEG using <what key?> and relays the [CSFBFedGWService.PagingRej](../../../feg/protos/csfb.proto) Message to the associated `FEG` over the [SyncRPCService.SyncRPC](../../../orc8r/protos/sync_rpc_service.proto) service

## Related Magma Documentation or Proposals

- [Cloud-Native Orc8r Secrets Proposal](https://github.com/magma/magma/discussions/5728)
  - Breadcrumbs of AGW<->Cloud SSL auth story
- [Security Overview](../orc8r/architecture_security.md)
- [Security Debugging](../orc8r/dev_security.md)
- [AWS Stack](../orc8r/dev_aws_stack.md)
