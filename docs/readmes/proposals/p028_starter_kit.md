---
id: p028_starter_kit
title: "Magma Starter Kit: Community Onboarding Package"
hide_title: true
---

# P028: Magma Starter Kit

Author(s): [@shahrahman-fb]

Last updated: 03/30/2026

## 1 Problem

Magma's barrier to entry is too high.
Today, a new community member who wants to deploy a private LTE
network must piece together information from 10+ documentation files,
independently source compatible hardware, provision SIM cards with
no vendor guidance, and choose from four different orchestrator
deployment methods with no clear recommendation for beginners.

The result is that only organizations with existing telecom expertise
can successfully deploy Magma.
This limits community growth, contribution, and real-world validation
of the platform.

Specific gaps:

- **No hardware bill of materials (BOM).**
  The docs list supported eNodeBs (Baicells Nova-243, mBS1100) but
  provide no purchasing guidance, no tested combinations, and no
  price expectations.
- **No SIM procurement or programming guide.**
  Docs reference "programmable SIMs" but do not explain where to buy
  them, how to program them, or what tools are needed.
- **No consolidated first-deployment guide.**
  Setup instructions are split across `deploy_install.md`,
  `deploy_config_enodebd.md`, `deploy_intro.md`, subscriber
  provisioning docs, and a 7-part Juju tutorial.
- **No software-only option for evaluation.**
  Users who want to test Magma without purchasing radio hardware have
  no clearly documented path (srsRAN integration exists but is buried
  in tutorial step 5 of 7).
- **No cost or resource estimates.**
  New users cannot budget for a deployment without deep research.

## 2 Solution

Define a **Magma Starter Kit** — a documented, tested, end-to-end
package that a new community member can follow to deploy a working
private LTE network from scratch.

The starter kit is not a product to ship.
It is a **specification plus guide** — a bill of materials, a
step-by-step deployment guide, and a set of validation tests that
confirm the deployment works.

### 2.1 Kit Tiers

Three tiers address different use cases and budgets:

#### Tier 1: Software-Only Evaluation Kit

**Cost:** $0 (existing hardware)
**Time to deploy:** 2-4 hours
**Audience:** Developers, evaluators, students

| Component | Specification |
|-----------|--------------|
| Host machine | x86-64, 16GB RAM, 100GB SSD, Ubuntu 22.04 |
| AGW | Docker-based deployment on host |
| Orchestrator | Local Docker Compose (dev mode) |
| Radio | srsRAN software radio simulator |
| UE | srsUE software UE simulator |
| SIM | Virtual — built-in test credentials |

**What you get:** A fully functional Magma EPC with simulated radio
and UE, capable of end-to-end data sessions.
No physical hardware or SIM cards required.

#### Tier 2: Minimal Hardware Lab Kit

**Estimated cost:** $1,500 - $3,000 USD
**Time to deploy:** 1-2 days
**Audience:** Small operators, university labs, community networks

| Component | Specification | Estimated Cost |
|-----------|--------------|----------------|
| AGW server | Mini PC (e.g., Intel NUC or equivalent), dual-core 2GHz+, 8GB RAM, 128GB SSD, 2x GbE | $300-500 |
| eNodeB | Baicells Nova-243 (FDD) or Neutrino-244 (TDD), CBRS Band 48 | $800-1,500 |
| SIM cards | 10x programmable USIM cards + USB SIM reader/writer | $50-100 |
| Orchestrator | AWS t3.xlarge or equivalent (can share AGW host for lab) | $0-100/mo |
| Antennas | Included with eNodeB or low-cost omni antenna | $0-100 |
| Cabling | Cat6 Ethernet, SMA/N-type RF cables as needed | $50 |
| UE | Any unlocked LTE phone or USB dongle supporting kit band | $50-200 |

**What you get:** A real private LTE cell serving real devices with
real SIM cards.
Suitable for indoor lab testing, small campus deployments, or
community network pilots.

#### Tier 3: Field Deployment Kit

**Estimated cost:** $5,000 - $15,000 USD
**Time to deploy:** 1-2 weeks
**Audience:** Operators, ISPs, enterprises

| Component | Specification |
|-----------|--------------|
| AGW server | Rack-mount or ruggedized industrial PC, 16GB RAM, 256GB SSD, 2x GbE, out-of-band management |
| eNodeB | Baicells Nova-243 outdoor unit, CBRS Band 48, PoE powered |
| SIM cards | 50-100x programmable USIM cards |
| Orchestrator | Production AWS/GCP deployment via Terraform or Ansible |
| Spectrum | CBRS SAS registration (US) or appropriate local license |
| Backhaul | Fiber or fixed wireless backhaul to AGW |
| UE devices | 5-10x unlocked LTE phones/CPE devices |

**What you get:** A production-grade outdoor private LTE cell
suitable for campus, rural, or enterprise deployment.

### 2.2 Starter Kit Guide Structure

A new top-level document, `docs/readmes/basics/starter_kit.md`,
would serve as the consolidated entry point.

**Outline:**

```
1. Choose Your Kit Tier
   - Decision matrix (budget, use case, timeline)

2. Bill of Materials
   - Per-tier component list with purchase links / vendor options
   - SIM card vendor recommendations
   - Tested hardware combinations

3. Network Planning
   - IP addressing plan (SGi, S1, management)
   - Spectrum selection (CBRS, test bands)
   - Firewall / port requirements
   - Estimated AWS costs for Orc8r (if cloud-hosted)

4. Step-by-Step Deployment
   a. Deploy Orchestrator (simplified — one recommended path)
   b. Install AGW (bare metal or Docker)
   c. Connect AGW to Orchestrator
   d. Configure eNodeB (or start srsRAN simulator)
   e. Provision SIM cards
   f. Add subscribers via NMS
   g. Attach UE and verify data session

5. Validation Checklist
   - [ ] AGW registered in Orc8r
   - [ ] eNodeB connected and transmitting
   - [ ] UE attaches and gets IP
   - [ ] UE can ping external host
   - [ ] NMS shows subscriber activity

6. Troubleshooting
   - Common failure modes and fixes
   - Log locations and diagnostic commands
   - How to collect a support bundle

7. Next Steps
   - Adding more subscribers
   - Scaling to multiple eNodeBs
   - Enabling QoS policies
   - Contributing back to Magma
```

### 2.3 SIM Card Guide

A dedicated sub-guide covering:

- **Recommended vendors:** sysmocom (sysmoUSIM-SJS1),
  Gemalto/Thales, or equivalent programmable USIM providers
- **Required tools:** PC/SC-compatible smart card reader
  (e.g., Omnikey 3121), pySim or similar open-source SIM
  programming tool
- **Programming steps:** Write IMSI, Ki, OPc, and PLMN to blank
  SIM using pySim
- **Matching EPC config:** Ensure AGW subscriber provisioning
  matches SIM credentials
- **Test SIM profile:** A reference profile for lab use with
  documented test credentials

### 2.4 Tested Hardware Combinations

Publish a **compatibility matrix** of tested AGW + eNodeB + UE
combinations:

| AGW Hardware | eNodeB | Band | UE Tested | Status |
|-------------|--------|------|-----------|--------|
| Intel NUC i5 | Baicells Nova-243 FDD | CBRS B48 | Pixel 6, Galaxy S21 | Verified |
| Mini PC (generic x86) | Baicells Neutrino-244 | CBRS B48 | Quectel RM500Q | Verified |
| Raspberry Pi 4 (Docker) | srsRAN simulated | Band 7 | srsUE | Verified |

*Initial matrix populated through community testing; expanded over
time via contribution.*

## 3 Implementation Plan

### Phase A: Documentation (2-4 weeks)

1. Write the consolidated starter kit guide
   (`docs/readmes/basics/starter_kit.md`)
2. Write the SIM card programming guide
   (`docs/readmes/basics/sim_programming.md`)
3. Create the hardware compatibility matrix
   (`docs/readmes/basics/hardware_compatibility.md`)
4. Update `docs/readmes/basics/quick_start_guide.md` to reference
   the starter kit as the recommended starting point

### Phase B: Software-Only Kit Validation (2-4 weeks)

1. Create a Docker Compose file or script that brings up
   AGW + Orc8r + srsRAN in a single command
2. Pre-configure test subscriber credentials
3. Validate end-to-end data flow on Ubuntu 22.04
4. Document in starter kit guide

### Phase C: Hardware Kit Validation (4-8 weeks)

1. Acquire Tier 2 hardware (or solicit community lab access)
2. Walk through the guide end-to-end on fresh hardware
3. Document failure modes and troubleshooting steps
4. Publish compatibility matrix results

### Phase D: Community Feedback (ongoing)

1. Announce starter kit on GitHub Discussions and Slack
2. Solicit community-tested hardware combinations
3. Iterate on guide based on feedback
4. Consider partnering with hardware vendors for discounted
   community bundles

## 4 Non-Goals

- **Shipping physical kits.**
  This proposal defines a specification and guide, not a product.
  Hardware procurement is the user's responsibility.
- **Commercial vendor partnerships.**
  While vendor recommendations are included for convenience,
  this is not an endorsement or commercial arrangement.
- **5G NR support.**
  The starter kit targets LTE (4G).
  5G SA support can be added as a future extension once Magma's
  5G capabilities mature.
- **Production deployment guide.**
  The starter kit is for lab, pilot, and small-scale deployments.
  Production deployments at scale require additional planning
  (HA, monitoring, capacity) beyond the kit's scope.
- **Replacing existing documentation.**
  The starter kit consolidates and links to existing docs rather
  than duplicating them.
  Detailed component documentation remains in its current
  locations.

## 5 Success Metrics

- A new user with no prior Magma experience can complete a
  Tier 1 (software-only) deployment in under 4 hours
- A new user can complete a Tier 2 (hardware lab) deployment
  in under 2 days, including hardware procurement lead time
- The starter kit guide is the #1 recommended entry point for
  new community members
- At least 5 community members report successful deployments
  using the guide within 6 months of publication
- Hardware compatibility matrix has at least 3 community-verified
  combinations within 6 months

## 6 Open Questions

1. **CBRS vs. international spectrum:**
  The Tier 2/3 kits assume CBRS (US).
  Should the guide include parallel paths for other regulatory
  environments (EU, India, Africa)?

2. **Orc8r hosting recommendation:**
  For Tier 1, Docker Compose is clear.
  For Tier 2, should we recommend local Orc8r (Ansible on same
  host) or cloud-hosted (AWS)?
  Cost vs. simplicity tradeoff.

3. **srsRAN version and compatibility:**
  srsRAN has evolved significantly.
  Which version is tested and supported with current Magma?

4. **All-in-one script feasibility:**
  How close can we get to a single `./deploy-starter-kit.sh`
  command for Tier 1?
  What are the blockers?
