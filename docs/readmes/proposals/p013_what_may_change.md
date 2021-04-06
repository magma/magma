---
id: p013_what_may_change
title: What May Change
hide_title: true
---

[This is a template for Magma's change proposal process, documented
[here](README.md).]

# Proposal: What May Change

Author(s): Scott Moeller, Marie Bremner

Last updated: 04/06/2021

Discussion at
[https://github.com/magma/magma/issues/4888](https://github.com/magma/magma/issues/4888).

## Abstract

Proposals can be made more robust if they are put through thought experiments of potential change.  By considering numerous low-but-nonzero probability system pressures or changes, a design might avoid some degree of breakage or limit the re-work blast radius should these events come to pass.  The list of brainstormed non-zero possibility events in this document are intended to represent larger structural changes and supplement standard system failure design analysis (security, availability, reliability, scalability, observability, maintainability...).

## How to use the List

Take the proposed solution, iterate down the list of possible stressors and determine the degree of impact and what subsystems of the proposal are impacted.  Watch out for collections of stressors that would trigger rework of multiple components, then explore whether a differing proposal could better isolate the subsystems impacted to reduce blast radius of these stressors.

## Possible Stressors

People have differing estimates of what probability terms mean. Here we will define the ~ranges for clarity.  See [If You Say Something Is “Likely,” How Likely Do People Think It Is?](https://hbr.org/2018/07/if-you-say-something-is-likely-how-likely-do-people-think-it-is).  Further, having too many levels of probability makes the exercise less useful - so here we stick with just three {`Unlikely`, `Possible`, and `Likely`}.  Note that even `Unlikely` sresstors should be considered for impact, as in aggregate they may cover decent outcome probability - and our assessment of their `Unlikly`ness is likely faulty (see Black Swan estimation!).

| Probability | ~Range |
|-|-|
| Unlikely | <10% |
| Possible | 10-80% |
| Likely | >80% |

### General Stressor Table

| Stressor | Probability |
| :--- | ---:|
| ARM platform support | Likely |
| Bare-metal AGW performance testing | Likely |
| AWS Marketplace CI testing | Likely |
| Support for Google Cloud, Microsoft Azure | Likely |
| Migration to a new CI provider | Likely |
| Migration of most CI tests to Containers | Likely |
| Increased coverage and number of integration test architectures | Likely |
| Upgrade / downgrade integration tests across wide swaths of versions | Likely |
| Continuous deployment for Magma development, Magma canary | Likely |
| Integration with new third-party CI test labs | Possible |
| Migration of most AGW software to the cloud | Possible |
| Forced removal of library dependency | Possible |
| Changes in prevalance of various programming languages | Possible |
| Changes to build automation / testing (Makefiles, CMake, etc) | Possible |
| Off-host build automation | Possible |
| Changes to dynamic vs static linking of c/cpp binaries | Possible |
| Changes to AGW package distribution | Possible |
| Changes to AGW linux distribution | Possible |
| Migration to multi-repo | Unlikely |
| Migration away from OVS datapath | Unlikely |
| Project resourcing drastic reduction | Unlikely |

### 3GPP Stressor Table

| Stressor | Probability |
| :--- | ---:|
| Converged core | Likely |

