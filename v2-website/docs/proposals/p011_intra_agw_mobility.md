---
id: p011_intra_agw_mobility
title: Intra-AGW Mobility (S1 Mobility)
hide_title: true
---

# Overview

*Status: In-Review*\
*Author: @shaddi*\
*Feedback requested from: @ulaskozat, @pshelar*\
*Last Updated: Mar 15, 2021*\
*Targeted Release: 1.5*

This document covers Intra-AGW mobility (aka S1 mobility support).

## Goal

Add S1 mobility support for intra-AGW mobility in situations where X2 handover is not available.

## Introduction

Currently, mobility in Magma is limited to RAN a single AGW where RAN elements support the X2 interface. This isn't always feasible; S1 mobility is the "fallback" method for mobility in LTE networks when X2 isn't avaialble. The scope of this proposal is enhancing intra-AGW mobility (that is, mobility that does not cross an AGW boundary) using S1 handover. Although *inter-AGW* mobility is not within scope of this proposal, the signalling mechanisms necessary for S1 handover will be useful for future mobility proposals.

The intended use case for this proposal is single cell tower (or small cluster of cell sites) connected to a single AGW, potentially using a mix of eNB equipment, where the network operator wishes to support seamless mobility anchored at the AGW. UEs should maintain 

## Proposal

The traditional S1 mobility call flow is our guide. A good overview of this particular flow is available [here](https://www.eventhelix.com/lte/handover/s1/).

### S1AP Signalling

We will add support to the MME for the following S1AP messages (see TS36.314 for additional details).

1. HandoverRequired
1. HandoverRequest
1. HandoverRequestAcknowledge
1. HandoverCommand
1. HandoverNotify
1. HandoverCancel
1. HandoverCancelAcknowledge

This will require modifications to the S1AP and MME_APP tasks for the mme.

### Data forwarding

S1 handover is used when there is no direct forwarding path between the two eNodeBs. So, we need to create an indirect forwarding path on the S11 interface and then modify bearers to the UE after handover has completed.

The following new proceedures will be added to the SGW:

1. CreateIndirectDataForwardingTunnelRequest
1. CreateIndirectDataForwardingTunnelResponse

### Out of scope
- RAN configuration: we assume that the RAN will be configured appropriately out-of-band to support mobility between cells (i.e., the neighbor table will be configured properly on all eNodeBs
- Mobility across AGWs
- Support for a ``large number'' of eNodeBs attached to a single AGW: The distributed nature of Magma means that any given AGW is responsible for a relatively small number of eNodeBs. This proposal does not change this assumption.
- Inter-RAT handover: Only intra-LTE handover will be supported at this time.

## Testing
Ideally, we'd have `s1aptester` tests for this feature; however, `s1aptester` does not currently support S1 handover. In the meantime, we'll be testing with the following setups:

- TM500 testbed set up @ FB. Test cases:
  - Test S1-Handover with 2 eNodeB and with Single Subscribers
  - Test S1-Handover with Multiple eNodeBs
  - Test S1-Handover with multi UE
  - Scale Test S1-Handover with 10 subscribers
  - (Negative Test Cases) - eNB rejects the HO

- Local test bed: 2x Baicells + Pixel 3 as UE

