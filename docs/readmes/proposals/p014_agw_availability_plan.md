---
id: p014_agw_availability_plan
title: AGW Availability Improvement Plan
hide_title: true
---

# Table of Contents

Generated with [github-markdown-toc](https://github.com/ekalinin/github-markdown-toc).

- [Table of Contents](#table-of-contents)
- [Document Objectives](#document-objectives)
  - [Regions of Document that are In Progress](#regions-of-document-that-are-in-progress)
- [Immediate Strategy and Priorities](#immediate-strategy-and-priorities)
  - [Summary](#summary)
  - [Per-Action Rationale and GH Issues](#per-action-rationale-and-gh-issues)
    - [Turn up GH Pull Review linters and static analysis tooling](#turn-up-gh-pull-review-linters-and-static-analysis-tooling)
    - [Gather Daemon Crash Stack Traces](#gather-daemon-crash-stack-traces)
    - [Audit MME AssertFatal calls](#audit-mme-assertfatal-calls)
    - [Turn up gRPC Telemetry](#turn-up-grpc-telemetry)
    - [Add latency telemetry](#add-latency-telemetry)
    - [Form a plan for OpenTracing in AGW](#form-a-plan-for-opentracing-in-agw)
    - [Document expected system state invariants and add to telemetry](#document-expected-system-state-invariants-and-add-to-telemetry)
    - [Author AGW stability chaos test](#author-agw-stability-chaos-test)
    - [Author AGW state chaos test](#author-agw-state-chaos-test)
    - [Lab test Clang build with ThreadSan](#lab-test-clang-build-with-threadsan)
    - [Pilot unit tests for uncovered code domains](#pilot-unit-tests-for-uncovered-code-domains)
    - [Address entangled c/c++ tech debt](#address-entangled-cc-tech-debt)
    - [Improve contributor velocity](#improve-contributor-velocity)
- [Measuring Availabiility](#measuring-availabiility)
  - [Using Exising Phone-Home Telemetry](#using-exising-phone-home-telemetry)
  - [Ground Truth](#ground-truth)
  - [Taxonomy of Software-Driven Outage Causes and Observability](#taxonomy-of-software-driven-outage-causes-and-observability)
  - [Taxonomy Details](#taxonomy-details)
    - [Functional behavior escape](#functional-behavior-escape)
    - [c and c++ manual memory managmement errors](#c-and-c-manual-memory-managmement-errors)
    - [Inability to recover daemons and serve traffic](#inability-to-recover-daemons-and-serve-traffic)
    - [gRPC Failures](#grpc-failures)
    - [Inconsistent state across AGW daemons](#inconsistent-state-across-agw-daemons)
    - [Mismatched IP data plane configuration](#mismatched-ip-data-plane-configuration)
    - [Overloaded IP data plane](#overloaded-ip-data-plane)
    - [Memory resource starvation](#memory-resource-starvation)
    - [Congestive collapse of 3GPP control plane processing](#congestive-collapse-of-3gpp-control-plane-processing)
- [Comprehensive Action Space](#comprehensive-action-space)
  - [Mind Map of Potential Actionns](#mind-map-of-potential-actionns)
  - [Per-Action GH Issues and Status](#per-action-gh-issues-and-status)

# Document Objectives

This document represents inputs from many Magma CODEOWNERS, Technical Leaders, and Managers. It is intended to be viewed as a community generated document.

Note that the following has a focus on `C` and `C++` codebase and binaries, though there are areas of large overlap with e.g. Python codebase of the AGW.

Note that some fraction of the below Improvement Plan has already been executed in the first quarter of 2021.  The mind map diagram depicts the space of possible availbiltiy enhancing actions for the system as it existed prior to the start of the quarter.

## Regions of Document that are In Progress

The following portions of this document are under active authorship. The document is being released early for feedback and visibility.

|Topic|Section|What is incomplete|
|-|-|-|
|Availability|[Ref](#using-exising-phone-home-telemetry)|Form a plan for availability estimation using existing phone-home telemetry
|Failure Taxonomy|[Ref](#taxonomy-of-software-driven-outage-causes-and-observability)|Columns `FB`, `Partner` and `Alert` need per-row references and planning
|Strategy|[Ref](#immediate-strategy-and-priorities)|Needs per-priority rationale updates / editing
|Mindmap|[Ref](#mind-map-of-potential-actions)|Needs GH Issue links for many possible tasks (some GH Issues don't yet exist either)

# Immediate Strategy and Priorities

## Summary

- Without tooling support, new code landing rate exceeds team ability to review
  - Immediately [Turn up GH Pull Review linters and static analysis tooling](#turn-up-gh-pull-review-linters-and-static-analysis-tooling)
- Gain automated insight into crash behavior of Magma Daemons in production
  - [Gather Daemon Crash Stack Teraces](#gather-daemon-crash-stack-traces) with present focus on sentry.io stack traces
- Some inputs cause MME to AssertFatal
  - [Audit MME AssertFatal calls and migrate to error handling where reasonable #fuzz](#audit-mme-assertfatal-calls)
- Stateless AGW pushes system further into distributed computing, but lacks tooling for test and debug
  - Improve performance engineering distributed systems observability
    - [Turn up gRPC telemetry](#turn-up-grpc-telemetry)
    - [Add latency telemetry](#add-latency-telemetry)
    - [Form a plan for OpenTracing in AGW](#form-a-plan-for-opentracing-in-agw)
  - Implement First Chaos Tests
    - [Documenting expected system state invariants and add to telemetry](#document-expected-system-state-invariants-and-add-to-telemetry)
    - [Author AGW stability chaos test](#author-agw-stability-chaos-test)
    - [Author AGW state chaos test](#author-agw-state-chaos-test)
- Sanitize against concurrency bugs
  - Based on V1.4 release - there may be some hidden concurrency that is not intended
  - [Run lab integration tests with c/c++ binaries built using Clang's tsan](#)
- Some important code regions lack examples or appraoch for test coverage
  - [Pilot unit tests for uncovered code domains](#pilot-unit-tests-for-uncovered-code-domains)
- Our entwined c/c++ exposes us to bugs, prevents health and safety improvements
  - [Address entangled c/c++ tech debt](#address-entangled-cc-tech-debt)
- Continuous improvement of developer velocity
  - [Improve contributor velocity](#improve-contributor-velocity)

## Per-Action Rationale and GH Issues
### Turn up GH Pull Review linters and static analysis tooling

**Immediate gains and pays compounding interest.**

While automated testing in all stripes is the most effective and thorough prevention for regressions, broadly it is also the most work to move the needle on. From where we are today, test authorship is a long haul fix and one we need to enable and support. We will work on that immediately, but linting and static code analysis can get us some percentage of wins across the entire code base (while accelerating review velocity) - so that has been a priority this quarter.  Additionally, code uniformity, standards, and linting - increase the code health and set community expectations on code quality.

It is also the case that some forms of non-adherance to Coding Standards makes test authorship harder.

High priority tasks in this category include clang-tidy Github Pull Request annotations, and application of Google's cpplint.  Cpplint ~requires we migrate our codebase formatting to Google Style, which has other benefits and is discussed in the associated GH Issue below.

- [x] [#4865](https://github.com/magma/magma/issues/4865) - Achieve additional clang build support for MME, track warnings in CI
- [x] [#4865](https://github.com/magma/magma/issues/4865) - [CI] Clang build support for MME, track warning counts through CI
- [x] [#4866](https://github.com/magma/magma/issues/4866) - [CI] Annotate GH Pull Requests with cranked up GCC / warning flags
- [x] [#5117](https://github.com/magma/magma/issues/5117) - Build Docker image capable of MME C test execution
- [x] [#5628](https://github.com/magma/magma/pull/5628) - [CI] Turn up Dockerfile Linting in GH PR
- [x] [#5850](https://github.com/magma/magma/pull/5850) - [CI] Turn up Github Action for Shellcheck
- [x] [#5962](https://github.com/magma/magma/pull/5962) - [AGW][Python] Add reviewdog annotations for Python diffs
- [x] [#5965](https://github.com/magma/magma/pull/5965) - [DevExp] Add ReviewDog misspell for PR spellcheck
- [x] [#6004](https://github.com/magma/magma/pull/6004) - [DevExp] Add reviewdog yamllint linter for PR annotations
- [ ] [#4867](https://github.com/magma/magma/issues/4867) - Enable clang-tidy reporting for AGW, annotate in GH Action
- [ ] [#6187](https://github.com/magma/magma/issues/6187) - [Proposal] migrate clang-format formatting to Google Style
- [ ] [#6163](https://github.com/magma/magma/pull/6163) - [DevExp] Add PR Linting from cpplint - Google Style linter

### Gather Daemon Crash Stack Traces

Similar to many other software defects in this document, with Magma's new Statess AGW migration, terminatino of daemons does not necessarily directly result in existing user interruption (though it may interfere with in-progress 3GPP signaling for e.g. Attach). But these binary crashes stress other types of bugs in the AGW, such as [Distributed computing bugs](#inconsistent-state-across-agw-daemons).  These subsequent bugs directly lead to avilability outages and service requests.

In order to prioritize Magma core developer efforts, automated collection and de-duplicatino of stack traces is essential.  themarwhal@gh has landed the capability to turn on Sentry.io stack trace capture, through configuration by the operator, in V1.5.  Facebook will be using this in all laboratory testing, and will be working with interested partners to establish their stack trace pipeliens and report findings to the Magma Core team.  There is ongoing improvement work still necessary - being led by themarwhal@.

- [x] [#5048](https://github.com/magma/magma/issues/5048) - [AGW] Automate process of crash stack trace reporting to DEV team

### Audit MME AssertFatal calls

More than one SEV (customer outage which required Magma Core software fixes) were root caused to `AssertFatal` calls within the Magma MME.  The control plane signal processing appears vulnerable to both unexpected inputs (perhaps triggered by atypical UEs or UE behavior) or signaling sequencing.  It appears that many of these `AssertFatal` calls need not be fatal, and would be better suited to error logging, telemetry accumulation, and discard of the offending traffic or messages.

- [ ] [#4879](https://github.com/magma/magma/issues/4879) - Audit MME AssertFatal calls and migrate to error handling where reasonable #fuzz

### Turn up gRPC Telemetry

When attempting to diagnose e.g. Magma AGW congestive collapse, there is a risk of chasing the appearance of CPU loading problems inside one daemon, while the driver could be systemic (e.g. feedback loops in control behavior).  Exploring the root cause of inter-daemon behavior is most efficiently started - in part - by observing the inter-daemon RPC telemetry.  If necessary, this couples well with [RPC Tracing](#form-a-plan-for-opentracing-in-agw).

- [ ] [#6394](https://github.com/magma/magma/issues/6394) - [AGW] Turn up per-service gRPC server/client telemetry

### Add latency telemetry

Observability of latency in critical interfaces of e.g. the Magma MME today (among others) is done exclusively through logs and human parsing / computation.  This creates barriers to ease of access, and limits observability into larger aggregate statistics of behavior.  Aside from [Turning up gRPC Telemetry](#turn-up-grpc-telemetry), we would benefit from dropping telemetry latency estimates into any critical processing. As immediate examples, of many, we would want to gather latency distributions for **state conversion per type** and **write latency to Redis**.

TODO: needs GH Issue(s)

### Form a plan for OpenTracing in AGW

In bisecting misbhavior of complex distributed systems, tracing is a superb tool for insight into inter-daemon interactions.  Magma has attempted some degree of manual tracing support in logs, but a community-supported open source tracing altenrative would be much more thorough / batteries-included.  This effort, combined with [Turning up gRPC Telemetry](#turn-up-grpc-telemetry) - are the first step in asessing certain classes of loading misbehavior, which Magma is actively exploring in order to improve scalability.

Tracing tooling has the auxiliary benefit of forming automated documentation of RPC interactions between daemons.

hcgatewood@gh is leading an investigation into OpenTracing unified solutions for Magma cloud and edge (AGW).  This will provide greatly enhanced observability and debug for RPC chains that cross multiple AGW daemons.

- [ ] [#6358](https://github.com/magma/magma/issues/6358) - Explore and describe general distributed tracing solution

### Document expected system state invariants and add to telemetry

Several very high value classes of tests cannot know the precise expected system state as the test progresses or completes.  This includes some system integration tests, and chaos tests (#6340), as two examples.  Further, production environments have very dynamic state but would benefit from invariant property checks, logging, telemetry and alerts - where feasible.

With this motivation, it makes sense to enumerate the system invariants that are anticipated in the AGW.  Some or many of these may already be tested in e.g. the LTE `integ_test` for example.  But a thorough accounting (and filling any holes) would provide strong real-world benefit if followed up with new tests, property checks, and telemetry.

- [ ] [#6386](https://github.com/magma/magma/issues/6386) - [AGW] Document expected system state invariants

### Author AGW stability chaos test

Historically there have been complex interactions between Magma daemons and their startup / restart sequencing.  This easily leads to health / deadlock issues, and can be effectively guarded against by a stability chaos test.  See the GH Issue for a description - but ultimately this is looking to ensure that when a random Magma daemon crashes, the system returns to function after some threshold duration (also covered is initial startup / rooting out any races).

- [ ] [#6339](https://github.com/magma/magma/issues/6339) - [AGW] Author daemon-restart-stability chaos test

### Author AGW state chaos test

As Magma Access Gateways are now very much subject to distributed computing bugs (more so than ever with stateless migration) - there are classes of issues which will best be discovered by [Jepsen style](https://jepsen.io/) chaos testing. Specifically, subjecting the system to a known set of system configuration transactions (to the extent possible), then measuring system invariants while subjecting the system to chaos behaviors (restart of daemons, killing sockets relied upon by gRPC sessions, etc).

The advantage to this style of test is that it is extremely capable of rooting out misbehavior of complex inter-dependent systems, often around transactional guarantees or other race conditions.

The downside is that the findings are harder to debug, as subsystem root cause isolation and reproduction is not so simple as other testing styles..

- [ ] [#6340](https://github.com/magma/magma/issues/6340) - [AGW] Author daemon-restart-state-health chaos test

### Lab test Clang build with ThreadSan

Now that we can build the MME with Clang, we should at least run a quick laboratory test with a tsan-enabled MME build (if not other C++ binaries) in order to see whether any concurrency failures are detected.  Thread sanitizer has a very low false positive rate, so its findings can reliably be assumed to be legitimate.

- [ ] [#6384](https://github.com/magma/magma/issues/6384) - [AGW] Run unit tests, lab tests with ThreadSan turned on

### Pilot unit tests for uncovered code domains

Some domains within Magma (e.g. MME) lack unit test coverage, and require various library or software architecture updates in order to enable testing.  Discovering what sorts of changes are necessary, enacting them, and leaving a sample set of example tests for other contributors - is an important task to get started on.  This will enable us to ask contributors to author unit tests on more code contributions, and enable gradual code coverager improvement.

TODO: Needs Github Issue(s)

### Address entangled c/c++ tech debt

Pulled from rationale in [#6385](https://github.com/magma/magma/issues/6385).

>The mixing of c and c++ in the present Magma codebase presents several very real challenges to reliability, test, and general code health.  These include but are not limited to the following.
>
>- Mixed logging strategies
>- Mixed testing frameworks
>- Mixed best-practices to support mocking / test injection
>- Extern ambiguity and lack of clear best practice
>- Subtle inter-language incompatibilities (certain initializations, empty struct sizes, ...)
>- Memory management tooling barrier
>- Prometheus telemetry client incompatibilities

While a rewrite of some or most c to c++ is in discussion, we might find that actually separating our c and c++ logic is a faster path to removal of these technical debt blockers (subject to performance evaluation). Related, how is the c to c++ migratino is giong to deal with the large volume of new 5G code landing?

- [ ] [#####](link) - TODO find c->c++ migration GH Issue
- [ ] [#6385](https://github.com/magma/magma/issues/6385) - [AGW] Consider breaking up mixed c/c++ binaries

### Improve contributor velocity

Contributor velocity is directly attributable to availability in multiple ways. Here are just a few reasons.

- Velocity clearly increases the rate at which availbility efforts can be accomplished broadly
- Velocity and developer experience often lead to higher quality product at the time of first PR
  - E.g. in-editor linting tools, clang-tidy, gcc warnings, etc
- Better test automation leads reviewers to feel free to ask for rebase / CI retest
  - Which can detect bad PRs
- Velcocity aided by better build tooling
  - Enables reviewers to ask for test authorship
  - Enables reviewers to feel free to question flaky tests

- [x] [#5710](https://github.com/magma/magma/pull/5710) - [DevEx] VS Code remote-containers setup

# Measuring Availabiility

Magma's waterfall release process and customer deployment strategies mean that new software to measure availability takes upwards of a 6-9 months to develop and deploy to the field broadly.  There may be opportunity for tightly integrated select partner outreach and coordination to deploy new measurement approaches on a smaller scale earlier, Magma will be exploring that.  The following table attempts to enumerate many but not all classes of availability outage. It is useful to have this collection of failure modes in mind, when evaluating metrics for observability of availability.

## Using Exising Phone-Home Telemetry

[In Progress](#in-progress)

With release V1.4, Magma has enabled turn-up of [phone-home metrics](https://docs.magmacore.org/docs/next/nms/metrics#list-of-metrics-which-are-currently-available). One goal of these metrics is to support availabilty estimation.

In this section, we will propose processing of these metrics to estimate availability, and if possible isolate root causes from the [taxonomy of software-driven outage causes](#taxonomy-of-software-driven-outage-causes-and-observability) below

## Ground Truth

Measuring ground truth availability data would require deployed UEs with continuous probing for healthy internet access.  Anything short of this is going to require omission of certain portions of outage root cause.  Unforutnately, this type of data gathering is extremely high effort to establish, and has real ongoing maintenance costs.  **For now, we will assume this is not a feasible strategy and will have to do our best with partial knowledge**.

## Taxonomy of Software-Driven Outage Causes and Observability

[1] Represent a linked plan to estimate through automated test or telemetry,  failures of this type, by end of 1H 2021.

[2] Represents operator's ability to establish alerting on failures of this type, by end of 1H 2021.

| Class | Description | Details | FB [1] | Partner [1] | Alert [2] |
|-|-|-|-|-|-|
| AGW Daemon Health | Functional behavior escape | [Ref](#functional-behavior-escape) | Some | Some | Some
| AGW Daemon Health | c/c++ manual memory management errors | [Ref](#c-and-c-manual-memory-managmement-errors) | Yes | Yes | Yes
| AGW Daemon Health | Inability to recover daemons and serve traffic | [Ref](#inability-to-recover-daemons-and-serve-traffic) | Yes | Yes | Yes
| AGW Daemon Health | gRPC failures | [Ref](#grpc-failures) | No | No | No
| AGW Daemon Health | Memory resource starvation | [Ref](#memory-resource-starvation) | [Yes](#todo-link-to-priorities)* | [Yes](#todo-link-to-priorities)* | Yes?
| AGW State | Inconsistent state across AGW daemons | [Ref](#inconsistent-state-across-agw-daemons) | No | No | No
| 3GPP Controle Plane Loading | Congestive collapse of 3GPP control plane processing | [Ref](#congestive-collapse-of-3gpp-control-plane-processing) | [Yes](#todo-link-to-priorities)* | [Yes](#todo-link-to-priorities)* | Yes?
| IP Dataplane | Mismatched IP data plane configuration | [Ref](#mismatched-IP-data-plane-configuration) | No | No | No
| IP Dataplane | Overloaded IP data plane | [Ref](#overloaded-ip-data-plane) | No | No | No?

## Taxonomy Details

### Functional behavior escape

This class includes availability outages due to things like unaticipated message processing or sequencing that sometimes trigger **AssertFatal**s, or regression in a feature or previously supported mode or UE.

### c and c++ manual memory managmement errors

These types of errors come in the flavor of segfaults or state corruption. With the inclusion of the AGW Stateless mode - the restart of an AGW daemon may not impact the availability of users directly (aside from perhaps any outstanding Attach or other state updates that were in progress) - but these restarts do stress and trigger other bugs that can cause larger scale outage (e.g. [Distributed computing bugs](#inconsistent-state-across-agw-daemons)).

### Inability to recover daemons and serve traffic

There are complex interdependencies between Magma AGW daemons, such that restart of some daemons necessitates restart of others (orchestrated by healthd and systemd).  It is possible to reach a state in which daemons repeatedly are restarted and do not come up in a healthy way, providing continuous outage of the 3GPP control plane.

### gRPC Failures

Magma makes liberal use of the gRPC library for inter-process communications both inside the AGW, and between the AGW and cloud elements of Magma.  Failures of gRPC communications channels are intended to be self-recovering, but in some language implementations and some modes of use - they may not always be.  Further, limitations on e.g. gRPC's maximum receive proto size have in the past caused outages for Magma customers.

### Inconsistent state across AGW daemons

This represents the classic distributed computing challenge.  Though it existed previously as well, the new Stateless Magma AGW mode exposes more potential for distributed computing problems across AGW daemons and applied state, including but not limited to root causes such as.

1. Failure to roll-back applied state (e.g. pipelined timeouts, ...)
2. Non-atomically applied state (e.g. race conditions)
3. Partially applied state (e.g. some daemons fail to make progress or drop updates)
4. Durability of state (e.g. unflushed writes, corruption, etc)

### Mismatched IP data plane configuration

In this scenario, aspects of the network configuration of the AGW (e.g. OVS or linux kernel) which are intended to be applied **by pipelined** for a particular UE, are not actually applied as intended.  Possible causes include.

- Caching errors in state application
  - E.g. pipelined or other **think** that configuration is applied, but it is not
- Unrecoverable errors in applying configuration

### Overloaded IP data plane

The AGW IP data plane will at some loading naturally begin to drop packets.  When this behavior begins depends upon hardware selection, hardware offloading capability, and possibly any bugs in the associated configuration.

### Memory resource starvation

This class of failure occurs most often when resource leakage is occuring.  Presently, with lack of per-daemon memory sandboxing, the AGW OOM-kills the first process to fail a memory allocation (possibly not the agent of the leak).

### Congestive collapse of 3GPP control plane processing

Here, the 3GPP control plane processing needed to support EU operations (attach, detach, roaming, etc) are overwhelmed by the rate of events.  The impacts of this can be of two degrees of severity.

- Temporary availability outage until loading decreases
- Long-running availability outage due to feedback loops

# Comprehensive Action Space

## Mind Map of Potential Actionns

This Mind Map is intended to be a ~comprehensive list of possible activities that would each incrementally improve the availability of the Magma product.  Explanation of each mind-map element is provided below the figure with a GH Issue reference (see the #### suffix on the mind map elements for GH issue numbers).  A prioritization and explanation of sequencing is additionally below that.

Bold lines represent work package trees with every leaf being a GH Issue.

Thin lines represent (some but far from all) of the dependency linkage. These help motivate sequencing.

```                                                                                
   LLVM Build                                                                                                           
 ━━━Support━━━━━━━━━━━━━━━━━Clang build (4865)━                                                                         
                                            ▲                                                                           
                                            └──────────────────────────────────────────────────────────────────────┬───┐
                                                                                                                   │   │
                                                                              ┏━━Infrastructure (5872)━            │   │
                                                                              ┃                                    │   │
           Code Health        Finding Count By Subdir                         ┣━━━━━━━TODOs (6308)                 │   │
       ━━━━━Monitoring━━━━┳━━━━━━per Master Commit━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫                                    │   │
                          ┃                                                   ┣━AssertFatals (4879)━               │   │
                          ┃                                                   ┃                                    │   │
                          ┃                                                   ┃      ┏━━━━━━━━━━Clang-Tidy (4867)━─┤   │
                          ┃                                                   ┃      ┃                             │   │
                          ┃                   ┏━Coverage filtering (6348)━    ┃      ┣━━━━━━Clang Warnings (6309)━─┘   │
                          ┗━━━━━Code Coverage━┃                               ┃      ┃                                 │
                                        ▲     ┗━━━see Integ test coverage━    ┃      ┣━━━━━━━━GCC Warnings (4866)━     │
                                        │                                     ┃      ┃                                 │
  ┌─────────────────────────────────────┘                                     ┃      ┣━━━━━━━━━━━━Hadolint (5628)━     │
  │                                       ┏━VS Code remote-containers (5710)━ ┗━━━━━━┫                                 │
  │                                       ┃                              ▲           ┣━━━━━━━━━━Shellcheck (5850)━     │
  │                               Dev Env ┣━━━━VS Code extensions (6346)━└──────────▶┃                                 │
  │                              ┏━━━━━━━━┫                                          ┣━━━━━━━━━Tflint (Terraform)━     │
  │                              ┃        ┗━━Git precommit hooks (6347)━━───────────▶┃                                 │
  │               Code Health    ┃                                                   ┣━━━━━━━━━━━━━cpplint (6163)━     │
  │           ━━━━━Automation━━━━┫                                                   ┃                           ▲     │
  │                     ▲        ┃        ┏━━━reviewdog support (6266)━              ┣━━━━━━━━clang-format (5498)┛     │
  │                     │        ┃        ┃                       ▲                  ┃                                 │
  │                     │        ┃  CI    ┣━━━━━━━━━━━━━CI Linters┛─────────────────▶┣━━━━━Git banned func (5385)━     │
  │                     │        ┗━━━━━━━━┫                                          ┃                                 │
  │                     │                 ┣━━━━linter tests (6263)━                  ┣━━━━━━━━━━━━misspell (5965)━     │
  │                     └──┐              ┃                                          ┃                                 │
  │    Code Review         │              ┗━━━━Mutation Tests━                       ┣━━━wemake-python-styleguide━     │
  │ ━━━Velocity and ━━━┳───┘                              ▲                 ┌──────▶ ┃                     (5962)      │
  │      Quality       ┗GitHub team reviews (6218)━       │                 │        ┃                                 │
  │                                                       │                 │        ┗━━━━━━━━━━━━yamllint (6004)━     │
  │                                                       │                 │                                          │
  │                                                       │                 │                                          │
  │                                                       │                 │                                          │
  │                                                       │                 └────────────────────────────────────────┐ │
  ├───────────────────────────────────────────────────────┘                                                          │ │
  │                                                                                                                  │ │
  │                                                                                                                  │ │
  │                           Source Dependency      ┏━Automated Version Re-Pin Process (6355)━                      │ │
  │                        ┏━━━━━━Strategy━━━━━━━━━━━┫                                                               │ │
  │                        ┃                         ┗━Static & Dynamic Linking Re-Eval━━━━━━━━Folly (6350)━         │ │
  │                        ┃   Build Speed,                                                                          │ │
  │                        ┃   Reliability,                                                                          │ │
  │                        ┣━Reproducibility,━━━━━━━━bazel build (4114)━                     Multiply redefined      │ │
  │                        ┃  Multi-language                                 ┏━━━━━━build/link flags, cstd, etc━     │ │
  │                        ┃               ▲                                 ┃                                       │ │
  │                        ┃               │              Monorepo of        ┃    Many CMake Projects tied with      │ │
  │                        ┃               │         ┏━━Separate Builds━━━━━━╋━━━━━━━━━━━━━━━━━━━━━━━━Makefiles━     │ │
  │                        ┃               │         ┃                     ▲ ┃                                       │ │
  │                        ┃               │         ┃                     │ ┗━━━━━Protobuf Copies and Rebuilds━     │ │
  │                        ┃               │         ┃                     │                                         │ │
  │     Build Tooling      ┃               │         ┃                     │                                         │ │
  │  ━━━━━━━━━━━━━━━━━━━━━━┫               │         ┃                     │                                         │ │
  │                        ┃               │         ┃                     │                                         │ │
  │                        ┃               │         ┃                     │                                         │ │
  │                        ┃               │         ┃                     │                                         │ │
  │                        ┃               │Or...    ┃   Missing CMake  ┏──┘                                         │ │
  │                        ┃                         ┣━━Best Practices━━┫                                            │ │
  │                        ┃   CMake Tech Debt       ┃                  ┗━Fix cmake (5222)━                          │ │
  │                        ┗━━━━━━━━━━━━━━━━━━━━━━━━━┫                                                               │ │
  │                                        ▲         ┣━━C_STANDARD (6239)━                                           │ │
  │                                        │         ┃                                                               │ │
  │                                        │         ┗━━━━━━━cmake_format━                                           │ │
  │                                        │                                                                         │ │
  │                                        └──────────────────┐                                                      │ │
  │                                                           │                                                      │ │
  │                    ┏━━━━━━━━━━━━Speed of test execution━──┤                                                      │ │
  │                    ┃                                      │                                                      │ │
  │                    ┣━━━━━━━━━━━━East of test authorship━──┤                                                      │ │
  │     Test Infra     ┃                                      │            ┏━timeline insights (5847)━               │ │
  │  ━━━━━━━━━━━━━━━━━━╋━━━━━━━━━Test execution selectivity━──┤      CI    ┃                                         │ │
  │     ▲  ▲           ┃                                      │  ━━━━━━━━━━┃                                         │ │
  │     │  │           ┗━━━━━━━━━━━━━Test de-flake workflow━──┘            ┃      AGW                                │ │
  │     │  │                                                               ┗Containerization━━Containerize (3410)━   │ │
  │     │  └──────────────────────────┐                                                                              │ │
  │     │                         ┏───┘                                                                              │ │
  │     │                         ┣━━━Performance testing (5012)━━────────────────────────────────────────────       │ │
  │     │                         ┃                                                                                  │ │
  │     │                         ┣━Integ test coverage (6343)━━━━────────────────────────────────────────────       │ │
  │     │                         ┃                                                                                  │ │
  │     │     Integration Tests   ┃    integ_test    ┏━━━━━━━Augment tests━                                          │ │
  │     │  ━━━━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━━━━┫                                                               │ │
  │     │               ▲         ┃                  ┗━━━━━Break up (6342)━──────────────────────────────────┐       │ │
  │     │               │         ┃                                                                          │       │ │
┌─┼─────┼───────────────┘         ┃   Basic Health                                                           │       │ │
│ │     │                         ┣━━━━━━━━━━━━━━━━━━━Healthy startup (6284)━                                │       │ │
│ │     │                         ┃                                                                          │       │ │
│ │     │                         ┃                                                                          │       │ │
│ │     │                         ┃   Distributed    ┏━━━━━━━━chaos stability (6339)━────────────────────────┤       │ │
│ │     │                         ┗━━━━Computing━━━━━┫                                                       │       │ │
│ │     │                                            ┗━━━━━chaos state health (6340)━────────────────────────┤       │ │
│ │     │                                                                                                    │       │ │
│ │     │                                                                                                    │       │ │
│ │     └─────────────────────────────────────────────────────────────────────────────────────────────┐      │       │ │
│ │                                                                                                   │      │       │ │
│ └────────────────────────────────────────────────────────────────────────────────────────────────┐  │      │       │ │
│                                                                                                  │  │      │       │ │
│                                                                                                  │  │      │       │ │
│                        c/c++ mixing   ┏━━━━━━━━━━━━━━━━replace c━                                │  │      │       │ │
│                    ┏━━━━━━━━━━━━━━━━━━┫                                                          │  │      │       │ │
│                    ┃        ▲         ┃ OR       separate c/c++                                  │  │      │       │ │
│                    ┃        │         ┗━━━━━━━━processes (6385)━                                 │  │      │       │ │
│                    ┃        └────────────────────────────────────────────────────────┐           │  │      │       │ │
│                    ┃                                                                 │           │  │      │       │ │
│                    ┃  Missing Critical ┏━━━absl::Status & absl::StatusOr (6151)━     │           │  │      │       │ │
│                    ┣━━━━━Libraries━━━━━┫                                             │           │  │      │       │ │
│                    ┃            ▲      ┗━━━━━━━━absl::Time━                          │           │  │      │       │ │
│                    ┃            │                                                    │           │  │      │       │ │
│                    ┃            └─────────────┐                                      │           │  │      │       │ │
│                    ┃                        ┏─┘                                      │           │  │      │       │ │
│                    ┃                        ┣────────────────────────────────────────┘           │  │      │       │ │
│                    ┃                        ┣━━━Interface prevalence━                            │  │      │       │ │
│                    ┃    Design for Test     ┃                                                    │  │   (5117)     │ │
│       C/C++ Tech   ┣━━━━━━━━━━━━━━━━━━━━━━━━╋━━━━━━Swappable effects━                            │  │      │       │ │
│     ━━━━━Debt━━━━━━┫                 ▲      ┃                                                    │  │      │       │ │
│                    ┃                 │      ┣━━━━━━━━Simulated clock━                            │  │      │       │ │
│                    ┃                 │      ┃                                                    │  │      │       │ │
│                    ┃                 │      ┗━━━━━━gRPC mocks (5531)━                            │  │      │       │ │
│                    ┃                 │                                                           │  │      │       │ │
│                    ┃                 └───────────────────────────────────────────────┐           │  │      │       │ │
│                    ┃                                             Author Unit Tests ┏─┘           │  │      │       │ │
│                    ┃                                         ┏━━━━━━━━━━━━━━━━━━━━━╋─────────────┘  │      │       │ │
│              ┌─────╋──────────────────────────────┐          ┃                     ┣────────────────┘      │       │ │
│              │     ┃                              ▼          ┃                     ┗───────────────────────┘       │ │
│              │     ┃                Improve Localized Test   ┃   Demonstrate New    ┏━━━━━━━Fuzz Testing (4873)━   │ │
│              │     ┃             ┏━━━━━━━━━Coverage━━━━━━━━━━╋━━━━Test Paradigms━━━━┫                              │ │
│              │     ┃             ┃                           ┃                      ┗━━━Property Testing (5042)━   │ │
│              │     ┃             ┃                           ┃                                                     │ │
│              │     ┃ Code Health ┃                           ┗━MME state machine chaos test (4934)━                │ │
│              │     ┣━━━━━━━━━━━━━┃                                                                                 │ │
│              │     ┃             ┣━━━━━━━━━━━━Apply best practices (5357,5490,6187..)━━━───────────────────────────┘ │
│              │     ┃             ┃                                                                                   │
│              │     ┃             ┃      Missing &                                                                    │
│              │     ┃             ┣━━Extraneous Header━━━━━━IncludeWhatYouUse (4868)━─────────────────────────────────┤
│              │     ┃             ┃      includes                                                                     │
│              │     ┃             ┃                                                                                   │
│              │     ┃             ┣━━━━━━━━━━━Unrooted include paths━                                                 │
│              │     ┃             ┃                                                                                   │
│              │     ┃             ┃    Code Structure     ┏━━Header-only include public APIs━                         │
│              │     ┃             ┗━━━━━━━━━━━━━━━━━━━━━━━┫                                                           │
│              │     ┃                                     ┗━━━━━Cyclical dependencies (4869)━                         │
│              │     ┃                                                                                                 │
│              │     ┃                                                                                                 │
│              │     ┃                           Old / Mutated /     ┏━━━━━━━━━━━━━━bstr━                              │
│              │     ┃                       ┏━━Untested Vendored━━━━┫                                                 │
│              │     ┃    Dependencies       ┃      Libraries        ┣━━━━━━━━━hashtable━                              │
│              │     ┣━━━━━━━━━━━━━━━━━━━━━━━┫                       ┃                                                 │
│              │     ┃                       ┃                       ┣━━━━━━━━━━━━━━ITTI━                              │
│              │     ┃                       ┃                       ┃                                                 │
│              │     ┃                       ┃                       ┗━━━━━━━━━━━━━━3gpp━                              │
│              │     ┃                       ┃                                                                         │
│              │     ┃                       ┃                                                                         │
│              │     ┃                       ┃    Unmaintained       ┏━━━━━━━━━━libfluid━                              │
│              │     ┃                       ┗━━━Dynamic Linked ━━━━━┫                                                 │
│              │     ┃                            Dependencies       ┣━━━━━━━━━cpp_redis━                              │
│              │     ┃                                               ┃                                                 │
│              │     ┃                                               ┗━━━━━━━━━━━━━asn1c━                              │
│              │     ┃                                                                                                 │
│              │     ┃                       ┏━━━━━━━━telemetry augmentation━                                          │
│              │     ┃  Observability        ┃                                                                         │
│              │     ┗━━━━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━logging library━                                          │
│              │                             ┃                                                                         │
│              └────────────┐                ┗━━━━━━━━━━━━━logging verbosity━                                          │
│                           │                                                                                          │
└────────────────────────┐  │                                                                                          │
                      ┏──┘  │                                                                                          │
                      ┃     │                                                                                          │
                      ┣─────┘                                                                                          │
                      ┃                                                                                                │
                      ┣━━━━━━━━━━━━━━msan━─────────────────────────────────────────────────────────────────────────────┤
                      ┃                                                                                                │
                      ┣━━━━━━━━━━━━hwasan━─────────────────────────────────────────────────────────────────────────────┤
        Sanitizers    ┃                                                                                                │
     ━━━━━━━━━━━━━━━━━╋━━threadsan (6384)━─────────────────────────────────────────────────────────────────────────────┘
                      ┃                                                                                                 
                      ┣━━━━━━━━━━━━━ubsan━                                                                              
                      ┃                                                                                                 
                      ┗━━━━━━━━━asan/lsan━                                                                              
                                                                                                                        
                                                                                                                        
                         Health    ┏━━━━━━━━━Magma health━                                                              
                    ┏━━Monitoring━━┃                                                                                    
                    ┃              ┗━━━━━━━━━━3GPP health━                                                              
                    ┃                                                                                                   
                    ┃    Crash                                                                                          
                    ┣━━Reporting━━━━━━━━sentry.io (5048)━                                                               
                    ┃                                                                                                   
                    ┃                           ┏━━━━━━━━telemetry support━                                             
                    ┃                  Local    ┃                                                                       
                    ┃             ┏━━━━━━━━━━━━━╋━━━━gRPC telemetry (6394)━                                             
       Failure      ┃  Telemetry  ┃             ┃                                                                       
 ━━━Observability━━━╋━━━━━━━━━━━━━┫             ┣━system invariants (6386)━                                             
        ▲           ┃             ┃             ┃                                                                       
        │           ┃             ┃             ┣━━━━━━━━━━Datapath health━                                             
        │           ┃             ┃             ┃                                                                       
        │           ┃             ┃             ┗━━━━━━━━━Critical latency━                                             
        │           ┃             ┃  Phone Home                                                                         
        │           ┃             ┗━━━━━━━━━━━━━━━━availability estimator━                                              
        │           ┃                                                                                                   
        │           ┃   Tracing                                                                                         
        │           ┣━━━━━━━━━━━━━━━━━OpenTracing (6358)━                                                               
        │           ┃                                                                                                   
        │           ┃                                                                                                   
        │           ┃                            ┏━━━━━━━━heap profiler━                                                
        │           ┃  Profiling      gperftools ┃                                                                      
        │           ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━leak detector━                                                
        │                 ▲                      ┃                                                                      
        │                 │                      ┗━━━━━━━━━cpu profiler━                                                
        └─────────────┐   │                                                                                             
                    ┏─┘   │                                                                                             
                    ┣─────┘                                                                                             
                    ┃      SEV                                                                                          
                    ┣━━━━━━━━━━━━━━━───┐                      to Improve Localized Test                                 
                    ┃                  │                       ┌──────────────────────▶                                 
     Issue Escape   ┃   Laboratory     ▼                       │   to Integration Tests                                 
 ━━━━━Review and ━━━╋━━Regressions━━─▶ ━━━━test augmentation━──┼──────────────────────▶                                 
     Forward Fix    ┃                  ▲                       │          to CI Linting                                 
                    ┃  Integration     │                       └──────────────────────▶                                 
                    ┗━━━━━━Test ━━━━───┘                                                                                
                       Regressions                                                                                      
```

## Per-Action GH Issues and Status

- [x] [#4865](https://github.com/magma/magma/issues/4865) - Achieve additional clang build support for MME, track warnings in CI
- [x] [#4865](https://github.com/magma/magma/issues/4865) - [CI] Clang build support for MME, track warning counts through CI
- [x] [#4866](https://github.com/magma/magma/issues/4866) - [CI] Annotate GH Pull Requests with cranked up GCC / warning flags
- [x] [#5048](https://github.com/magma/magma/issues/5048) - [AGW] Automate process of crash stack trace reporting to DEV team
- [x] [#5117](https://github.com/magma/magma/issues/5117) - Build Docker image capable of MME C test execution
- [x] [#5628](https://github.com/magma/magma/pull/5628) - [CI] Turn up Dockerfile Linting in GH PR
- [x] [#5710](https://github.com/magma/magma/pull/5710) - [DevEx] VS Code remote-containers setup
- [x] [#5850](https://github.com/magma/magma/pull/5850) - [CI] Turn up Github Action for Shellcheck
- [x] [#5962](https://github.com/magma/magma/pull/5962) - [AGW][Python] Add reviewdog annotations for Python diffs
- [x] [#5965](https://github.com/magma/magma/pull/5965) - [DevExp] Add ReviewDog misspell for PR spellcheck
- [x] [#6004](https://github.com/magma/magma/pull/6004) - [DevExp] Add reviewdog yamllint linter for PR annotations
- [ ] [#3410](https://github.com/magma/magma/issues/3410) - AGW Containerization
- [ ] [#4114](https://github.com/magma/magma/issues/4114) - Explore Bazel as build system replacement
- [ ] [#4867](https://github.com/magma/magma/issues/4867) - Enable clang-tidy reporting for AGW, annotate in GH Action
- [ ] [#4879](https://github.com/magma/magma/issues/4879) - Audit MME AssertFatal calls and migrate to error handling where reasonable #fuzz
- [ ] [#5012](https://github.com/magma/magma/issues/5012) - [Reliability] Benchmarking performance between stateful and stateless modes
- [ ] [#5222](https://github.com/magma/magma/issues/5222) - [AGW] Our Cmake files need standardization / cleanup / linting
- [ ] [#5331](https://github.com/magma/magma/issues/5531) - [SessionD] Migrate unit tests to use protoc generate client mocks instead of a manually created server mock
- [ ] [#5357](https://github.com/magma/magma/issues/5357) - [MME] Address any pertinent clang compiler warnings in MME
- [ ] [#5385](https://github.com/magma/magma/issues/5385) - [AGW] Audit for use of Git's list of banned C functions
- [ ] [#5498](https://github.com/magma/magma/pull/5498) - [CI] Add Clang-Format linting as a GH Action on Pull Requests
- [ ] [#5490](https://github.com/magma/magma/issues/5490) - [AGW] Silence extraneous GCC -Wunused-parameter warnings
- [ ] [#5847](https://github.com/magma/magma/issues/5847) - [CI] Improve visibility into test duration regressions
- [ ] [#5872](https://github.com/magma/magma/issues/5872) - [Proposal] CI metrics pipeline to track master branch behavior
- [ ] [#6163](https://github.com/magma/magma/pull/6163) - [DevExp] Add PR Linting from cpplint - Google Style linter
- [ ] [#6218](https://github.com/magma/magma/discussions/6218) - Proposed approvers-* teams for PR reviews
- [ ] [#6151](https://github.com/magma/magma/issues/6151) - [AGW] Magma c++ needs adoption of Status and Result types
- [ ] [#6239](https://github.com/magma/magma/issues/6239) - [AGW] Migrate CMake to use of C_STANDARD
- [ ] [#6263](https://github.com/magma/magma/issues/6263) - [CI] Investigate Github Actions testing framework
- [ ] [#6266](https://github.com/magma/magma/issues/6266) - [CI] Workaround to support Reviewdog pr-gh-review comments
- [ ] [#6284](https://github.com/magma/magma/issues/6284) - [CI][PrecommitCheck] Add a sanity check to bring up all FeG, CWAG containers to see everything builds and doesn't crash
- [ ] [#6308](https://github.com/magma/magma/issues/6308) - [CI] Track TODO counts in Magma codebase by major branch
- [ ] [#6339](https://github.com/magma/magma/issues/6339) - [AGW] Author daemon-restart-stability chaos test
- [ ] [#6340](https://github.com/magma/magma/issues/6340) - [AGW] Author daemon-restart-state-health chaos test
- [ ] [#6342](https://github.com/magma/magma/issues/6342) - [AGW] Break up LTE integ_tests
- [ ] [#6343](https://github.com/magma/magma/issues/6343) - [AGW] Integration test code coverage
- [ ] [#6346](https://github.com/magma/magma/issues/6346) - [DevExp] Add available linting tools as extensions in VS Code .devcontainer
- [ ] [#6347](https://github.com/magma/magma/issues/6347) - [DevExp] Explore Git pre-commit hooks for linting
- [ ] [#6350](https://github.com/magma/magma/issues/6350) - [AGW] Move Folly dependency to static linking
- [ ] [#6355](https://github.com/magma/magma/issues/6355) - [AGW] Form strategy for dependency pinning and upgrade automation
- [ ] [#6358](https://github.com/magma/magma/issues/6358) - Explore and describe general distributed tracing solution
- [ ] [#6384](https://github.com/magma/magma/issues/6384) - [AGW] Run unit tests, lab tests with ThreadSan turned on
- [ ] [#6385](https://github.com/magma/magma/issues/6385) - [AGW] Consider breaking up mixed c/c++ binaries
- [ ] [#6386](https://github.com/magma/magma/issues/6386) - [AGW] Document expected system state invariants
- [ ] [#6394](https://github.com/magma/magma/issues/6394) - [AGW] Turn up per-service gRPC server/client telemetry
