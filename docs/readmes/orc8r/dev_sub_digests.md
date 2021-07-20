---
id: dev_sub_digests
title: Subscriber Digests
hide_title: true
---

# Subscriber Digests

Orc8r manages subscriber data through its `subscriberdb` service, and streams the data southbound to AGWs on a regular basis. In this process, a digest-cache pattern is used to intelligently recognize and send down only the data necessary for syncing AGWs. This document overviews the said pattern.

## Motivation

Previously, orc8r follows a vanilla data transmission pattern that sends all subscriber data down to all gateways at every configurable interval.

This introduces various issues at scale, a major one being the significant amount of network pressure this pattern necessitates. For example, to achieve our [GA scalability targets](https://docs.google.com/document/d/1P306DmuC1CFi7bqz4_VVi7v4Ewr1y-dD47C-s5XzdJ8/edit), we need to support **~20TB of network load per month per network**, which is definitively unsustainable.

Therefore, starting from v1.6, as part of our [subscriber scaling project](../proposals/p010_subscriber_scaling.md), orc8r implements the subscriber digests pattern to minimize network pressure and optimize the current mode of subscriber data transmission. This pattern utilizes the assumption of reads >> writes, and focuses on two major improvements
1. Instead of syncing at every interval, only sync when AGW is out-of-sync
2. Instead of syncing all data, sync the cloud-AGW data changeset when reasonable

## Overview

The architectural additions implemented for the subscriber digests pattern are highlighted below.

![Subscriber Digests Architecture](assets/orc8r/subscriber_digests_architecture.png)

### Digest Generation

This pattern relies on taking deterministically encoded data snapshots, a.k.a. "digests", that are used to represent the local data version. To this end, orc8r library code `mproto` is created using the golang third-party library `proto` to conduct consistently deterministic data serializations.

The subscriber digests pattern specifically generates and utilizes two types of digests

1. **Flat digests**: a digest of the entire set of subscriber data of a network
2. **Per-subscriber digests**: a digest of a single subscriber data object; for a network, the list of per-subscriber digests is tracked and distributed en masse

NOTE: Currently, a decoupling process between subscriber data objects and their gateway-specific APN resource configurations is conducted to make sure the generated digests is representative and also network-general. See [additional notes](./dev_sub_digests.md#apn-resource-handling) for more details.

### Subscriber Digests Cache
`subscriberdb_cache` is a single-pod service in sole charge of managing network digests. The service consists mainly of worker code that constantly generates new flat digests and per-subscriber digests at a configurable interval for each network.

The generated digests are written to a cloud SQL store to which `subscriberdb` service has read-only access. In this way, `subscriberdb` can directly read the most recently cached digests from the SQL store, instead of making extraneous gRPC calls over the network and having multiple servicer workers compete for the `subscriberdb_cache` endpoints.

Additionally,  `subscriberdb_cache` also caches the subscriber data that the digests represent. This is to ensure consistency between the cached digests and the actual subscriber data objects in a `subscriberdb` servicer response.

That is, every time an update of the digests are detected, `subscriberdb_cache` loads all subscriber data from `configurator` and applies the batch update to store. In this way, `subscriberdb_cache` acts as the single source of truth from which the `subscriberdb` servicer loads all data required for the subscriber digests pattern.

`subscriberdb_cache` supports loading subscriber data either directly by IDs or by pages. To that end, it generates its own page tokens using the same mechanism as in  `configurator` for simplicity.

### AGW Digest Store

Aside from its subscriber database, a gateway also manages SQL stores for its local flat and per-subscriber digests.

Every time a gateway engages in a sync session with the cloud, it receives the latest cloud digests, with which it updates its local SQL store. These digests represent the version of subscriber data stored in the gateway, and they are used in the gateway's subsequent requests to the cloud.

### Cloud-AGW Callpath

The current interactions between the `subscriberdb` cloud servicer and the AGW client involve 3 servicer endpoints

1. `CheckSubscribersInSync`
    - AGW sends its local flat digest
    - Cloud compares the AGW digest with the cloud digest. Depending on whether they are the same, the cloud returns a signal for whether the AGW is in-sync

2. If AGW is not in-sync: `SyncSubscribers`
    - AGW sends its local per-subscriber digests
    - Cloud compares the AGW digests with the cloud digests, and calculates the changeset between the two sets of per-subscriber digests. That is, it locates the subscribers to renew or to delete on the AGW side
        - If the changeset is of a reasonable size (configurable), returns the data needed to update AGW locally, as well as the latest cloud digests
        - If not, returns a signal for AGW to conduct a full-on resync
3. If AGW receives resync signal: `ListSubscribers`
    - AGW sends desired pageSize, as well as the current page token (until all pages are fetched)
    - Cloud returns the corresponding page of data loaded from `subscriberdb_cache`, as well as the latest cloud digests

## Additional Notes

### APN Resource Configs Handling

Currently, gateway-specific APN resources configurations are also stored in subscriber data objects streamed down to AGWs. This is an issue to be addressed by ongoing refactoring projects.

Until a fix is landed, to ensure that the generated subscriber digests are the same for all gateways in a network, the APN resources configs are extracted from subscriber data objects beforehand, and captured in its own digest instead. That is, the subscriber digests generated in `subscriberdb_cache` are concatenations of `apnResourceDigest` and `subscriberDigest`; in the generation of the `subscriberDigest`, the APN resource configurations are left blank.

Similarly, the cached subscriber data objects in `subscriberdb_cache` also don't contain APN resource configs. As a result, the `subscriberdb` servicers would currently append APN resource configs to subscriber data objects loaded from cache, before streaming it down to AGWs.

### Orc8r-directed Resync

For reliability, and also to avoid hash collision-induced edge cases in the subscriber digests pattern, an orc8r-directed resync is enforced on the AGWs at a configurable interval (e.g. once per day).

To this end, `subscriberdb` tracks the last resync times of each gateway, and enforces the resync for a gateway hasn't undergone a resync in a while.

NOTE: Since the cloud servicer doesn't track the success status of the requests on the AGW end, it takes an AGW request for the last page of subscriber data as an approximate indication that the AGW has come to the end of a resync.
