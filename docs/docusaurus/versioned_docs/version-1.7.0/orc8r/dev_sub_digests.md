---
id: version-1.7.0-dev_sub_digests
title: Subscriber Digests
hide_title: true
original_id: dev_sub_digests
---

# Subscriber Digests

Orc8r manages subscriber configs through its *subscriberdb* service, syncing the data southbound to AGWs on a regular basis. To perform this efficiently, data digests are taken and compared to intelligently send only the data necessary for syncing with an AGW. This document provides an overview of this digest pattern.

## Motivation

Previously, Orc8r followed a basic data sync pattern that sent all subscriber configs down to all gateways once per configurable interval, regardless of whether the AGW actually needs the data.

This introduces a number of issues at scale, namely significant network pressure. For example, to support 20k subscribers over 500 active gateways, subscriber config sync alone would generate at least [**20TB of network pressure per month per network**](../proposals/p010_subscriber_scaling.md#context).

Therefore, starting with Magma v1.6, Orc8r uses deterministic protobuf digests to only sync subscriber configs when they're updated. This pattern utilizes the assumption `#reads >> #writes` to retain a [level-triggered architecture](https://haobaozhong.github.io/design/2019/08/28/level-edge-triggered.html), providing two major improvements

1. Instead of syncing at every interval, only sync when AGW is out of sync
2. Instead of syncing all data, only transmit updates

## Overview

The architectural additions implemented for the subscriber digests pattern are highlighted below.

![Subscriber Digests Architecture](assets/orc8r/subscriber_digests_architecture.png)

### Digest generation

This pattern relies on generating, storing, and communicating digests (aka hashes, or deterministically-encoded, constant-length "snapshots") of the data to be synced.

The subscriber digests pattern specifically generates and utilizes two types of digests

1. **Root digests**: digest of the entire set of subscriber configs in a network
2. **Leaf digests**: digest of a single subscriber config object. For a network, this list of per-subscriber digests is tracked and distributed en masse

NOTE: currently, a decoupling process between subscriber config objects and their gateway-specific APN resource configurations is conducted to ensure the generated digests are representative and also network-general. See [additional notes](#apn-resource-configs-handling) for more details.

### Subscriber digests cache

*subscriberdb_cache* is a single-pod service in charge of managing network digests. The service constantly generates new digests at a configurable interval for each network.

The generated digests are written to a cloud SQL store to which *subscriberdb* has read-only access. In this way, *subscriberdb* can directly read the most up-to-date digests from the SQL store, preserving the limited resources of the *subscriberdb_cache* singleton.

Additionally, *subscriberdb_cache* also caches the subscriber configs that the digests represent. This ensures consistency between the cached digests and the actual subscriber configs in a *subscriberdb* servicer response.

That is, every time an update of the digests are detected, *subscriberdb_cache* loads all subscriber configs from *configurator* and applies the batch update to store. In this way, *subscriberdb_cache* acts as the single source of truth from which the *subscriberdb* servicer loads all data required for the subscriber digests pattern.

*subscriberdb_cache* supports loading subscriber configs either directly by IDs or by pages. To that end, it generates its own page tokens using the same mechanism as in *configurator*.

### AGW digest store

Aside from its subscriber database, a gateway also manages SQL stores for its local root and leaf digests.

Every time a gateway synchronizes its subscriber configs with the cloud, it receives the latest cloud digests, which it writes to local SQL store. These digests represent the version of subscriber configs stored in the gateway, and they are used in the gateway's subsequent requests to the cloud.

### Cloud-AGW callpath

The current interactions between the *subscriberdb* cloud servicer and the AGW client involve 3 servicer endpoints

1. `CheckInSync`
    - AGW sends its local root digest
    - Cloud compares the AGW digest with the cloud digest. Depending on whether they are the same, the cloud returns a signal for whether the AGW is in-sync

2. If AGW is not in-sync: `Sync`
    - AGW sends its local leaf digests
    - Cloud compares the AGW digests with the cloud digests, and calculates the changeset between the two sets of leaf digests. That is, it locates the subscribers to renew and delete on the AGW side
        - If the changeset is of a reasonable (configurable) size, returns the data needed to update AGW locally, as well as the latest cloud digests
        - If not, returns a signal for AGW to conduct a full-on resync
3. If AGW receives resync signal: `ListSubscribers`
    - AGW sends desired page size, as well as the current page token (until all pages are fetched)
    - Cloud returns the corresponding cached page of data as well as the latest cloud digests

## Additional Notes

### APN resource configs handling

Currently, gateway-specific "APN resource" configurations are also stored in the subscriber config objects synced to AGWs. This is a legacy pattern that will be redressed in an upcoming refactor.

As an immediate workaround, to ensure the generated subscriber digests are the same for all gateways in a network, the APN resources configs are extracted from subscriber config objects beforehand, and captured in its own, separate digest. That is, the subscriber digests generated in *subscriberdb_cache* are concatenations of `apnResourceDigest` and `subscriberDigest`. In the generation of the `subscriberDigest`, the APN resource configurations are left blank.

Similarly, the cached subscriber config objects in *subscriberdb_cache* also don't contain APN resource configs. As a result, the *subscriberdb* servicers currently append APN resource configs to subscriber config objects loaded from cache, before transmitting the full configs down to AGWs.

### Orc8r-directed resync

For added reliability, and to remediate (trivially improbable) hash collisions, an Orc8r-directed resync is enforced on the AGWs at a configurable interval (e.g. once per day).

To this end, *subscriberdb* tracks the last resync times of each gateway, and enforces the resync for gateways which haven't undergone a resync in a while.
