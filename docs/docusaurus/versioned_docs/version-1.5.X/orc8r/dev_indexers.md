---
id: version-1.5.0-dev_indexers
title: State Indexers
hide_title: true
original_id: dev_indexers
---

# State Indexers

This document describes the state indexing pattern, which can be used e.g. to create reverse maps of reported gateway state.

## Motivation

Gateways use the state service as the default mechanism for reporting their state. The state service accepts type-annotated blobs, meaning gateways can report state in domain-specific formats.

Orc8r stores reported gateway state as blobs keyed by the primary key triplet `{network ID, state type, state key}`. We call this *primary state*. With this pattern, services can access reported state of a particular type with the state's key, which is often a well-known value like an IMSI or gateway hardware ID.

However, sometimes services need to access a piece of state and all they know is the value of part of the blob. For example, finding the IMSI that owns an IP address.

To support this functionality, Orc8r provides the state indexer pattern. State indexers receive reported state at the time of reporting, optionally transform it, then write it to their own storage under a new primary key of their choosing. We call this transformed state *derived state*.

## Overview

### State service

The state service stores gateway state as keyed blobs. Consider the following two examples, which we'll use as examples throughout this document

- `IMSI -> directory record blob`
- `HWID -> gateway status blob`

Since state values are stored as arbitrary serialized blobs, the state service has no semantic understanding of stored values. This means searching over stored values would otherwise require an O(n) operation. Examples include

- Find IMSI with given IP: must load all directory records into memory
- Find all gateways that haven't checked in recently: must load all gateway statuses into memory

### Derived state

The solution is to provide customizable, online mechanisms for generating derived state based on existing state. Existing, "primary" state is stored in the state service, and derived, "secondary" state is stored in whichever service owns the derived state. Examples include

- Reverse map of directory records
    - Primary state: `IMSI -> directory record`
    - Secondary state: `IP -> IMSI` (stored in e.g. directoryd)
- Reverse map of gateway checkin time
    - Primary state: `HWID -> gateway status`
    - Secondary state: `checkin time -> HWID` (stored in e.g. metricsd)
- List all gateways with multiple kernel versions installed
    - Primary state: `HWID -> gateway status`
    - Secondary state: list of gateways (stored in e.g. bootstrapper)

### State indexers

State indexers are Orchestrator services registering an `IndexerServer` under their gRPC endpoint. Any Orchestrator service can provide its own indexer servicer.

The state service discovers indexers using K8s labels. Any service with the label `orc8r.io/state_indexer` will be assumed to provide an indexer servicer.

Indexers provide two additional pieces of metadata: version and types

- *version* is a positive integer indicating when indexer requires reindexing
- *types* are a list of state types the indexer subscribes to

These metadata are indicated by K8s annotations

- `orc8r.io/state_indexer_version` is a positive integer
- `orc8r.io/state_indexer_types` is a comma-separated list of state types

### Reindexing

When an indexer's implementation changes, its derived state needs to be refreshed. This is accomplished by sending all existing state (of desired types) through the now-updated indexer.

An indexer indicates it needs to undergo a reindex by incrementing its version (exposed via the above-mentioned annotation). From there, the state service automatically handles the reindexing process.

Metrics and logging are available to track long-running reindex processes, as well an indexers CLI which reports desired and current indexer versions.

### Architecture overview

![State indexing codepath](assets/orc8r/state_indexing_codepath.png)

## Additional notes

### Implementing a custom indexer

To create a custom indexer, attach an `IndexerServer` to a new or existing Orchestrator service.

A service can only attach a single indexer. However, that indexer can choose to multiplex its functionality over any desired number of "logical" indexers.

See the orchestrator service for an example custom indexer.

### Automatic reindexing support

Automatic reindexing is *only supported with Postgres*. Deployments targeting Maria will need to use the indexer CLI to manually trigger reindex operations. Run `/var/opt/magma/bin/indexers` from an Orc8r application container to view the help text for this command.

### State can go stale

The state indexer pattern currently provides no mechanism for connecting primary and secondary state. This means secondary state can go stale. Where relevant, consumers of secondary state should take this into account, generally by checking the primary state to ensure it agrees with the secondary state.

### Existent race condition

There is a trivial but existent race condition during the reindex process. Since the index and reindex operations both use the Index gRPC method, and the index and reindex operations operate in parallel, it's possible for an indexer to receive an outdated piece of state from the reindexer. However, this requires

- Reindexer read old state
- New state reported, indexer read new state
- Indexer Index call completed
- Reindexer Index call completed

If this race condition is intolerable to the desired use case, the solution is to separate out the Index call into `Index` and `Reindex` methods. This is not currently implemented as we don't have a concrete use-case for it yet.
