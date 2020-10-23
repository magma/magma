---
id: p006_hook_teravm_to_ci
title: Hooking up TeraVM infra to CircleCi
hide_title: true
---

# Hooking up TeraVM infra to CircleCi

*Status: In review*
*Feature owner: @tmdzk*
*Feedback requested from: @uri200 and tiplab team*
*Last Updated: 10/23*

## Summary

TeraVM has proven to be an efficient tool in terms of catching magma issue and
should be part of our entire CI process. The goal of that proposal is to find
the best way to include them to our CI workflow.

#### What is TeraVM?

TeraVM is an application emulation and security performance solution, delivering
comprehensive test coverage for application services, wired and wireless
networks.
[TeraVM commercial link](https://www.viavisolutions.com/en-us/products/teravm)

## Motivation

TipLab and magma engineers have been working hard to get that environnement up
and running and they already catch a lot of bugs in lte/feg and orc8r code.
However TerraVM is only run manually by those engineers.


## Goals

Be able to run TeraVM on master at least on a daily basis.

## Implementation Phases

**Phase 1 - Doc design all CI implementations**

(October 2020)

There are several ways to implement TeraVM in the CI and we need to list all
solutions and make the one we pick is the safest and the most efficient

**Phase 2 - Select one CI implementation and execute it**

(November 2020 - December 2020)

In the second phase, we'll meet and decide which CI implementation will be used
we need to divide the work inbetween magma team/tiplab team/FreedomFi

**Phase 3 - Test it on master on a daily basis**

(December 2020 - January 2021)

When it's done we need to monitor all tests and make sure they are stable and
every single developer knows how to take advantage of those

**Phase 4 - Make it release blocker**

(January 2021)

TeraVm tests have to be release blocker.

**Phase 5 - Improve reporting**

(January 2021 - ???)

Native TeraVM reporting is definitely not the best, keeping track of recent
tests, history is not ideal and we need a find a way to store test results and
history in a better way.

### Design of the implementation
---


#### Design 1

+----------------------------------------------------+
|                                                    |
| CircleCI Cloud                                     |
|                                                    |
|  +-------------------+     +-------------------+   |
|  |                   | 1.spin                  |   |
|  |                   +---->|   CircleCI node   |   |
|  |                   |     |                   |   |
|  |     Circle CI     |     ++-------+ +-------++   |
|  |                   |     || Open  | |  Fab  ||   |
|  |                   |     ||  VPN  | |  file ||   |
|  |                   |     |+-------+ +-------+|   |
|  +-------------------+     +-------------------+   |
|                            | |              |      |
+----------------------------------------------------+
             |---------------| |2. update     | 2.update
+----------------------------------------------------+
|            |                 |              |      |
| Tiplab     |3. Trigger test  |              |      |
| Infra      |                 |              |      |
|   +--------v----+  +---------v--+  +--------v----+ |
|   |             |  |            |  |             | |
|   |             |  |            |  |             | |
|   |   NG40      |  |  AGW       |  | FEG         | |
|   |             |  |            |  |             | |
|   +-------------+  +------------+  +-------------+ |
|                                                    |
+----------------------------------------------------+

##### Components

CircleCI Cloud:

- CircleCI main component
- CircleCI nodes: spinned up on the fly for every job

Tiplab Infra:

- AGW
- FEG
- NG40: TeraVM "controller"

##### Workflow

1. CircleCi spins up a new CircleCI node in order to run the job where it will
automatically run an OpenVPN service that will connect to the TipLab VPN and
be able to access all the TipLab Infra

2. CircleCI node can access the TipLab Infra and will update both agw and feg
with the latest available tag using the FabFile (in the future those packages
can be build directly on the CI node and transferred to the components directly
or uploaded to a special TeraVM registry)

3. CircleCI node can access the TipLab Infra and will use the fabfile to Trigger
ng40 tests.

##### Pros

- Easy to implement

##### Cons

- CircleCI will move to FreedomFi which means we'll have to share vpn secrets
with them.
- OpenVPN server is running on an infra that we don't control



#### Design 2

+-----------------------------------------------------+
|  CircleCI Cloud                                     |
|                                                     |
|                                                     |
|  +------------+               +-------------------+ |
|  |            |               |                   | |
|  |            |1. Spin node   |                   | |
|  |CircleCI    +-------------->+CircleCI node      | |
|  |            |  with secret  |                   | |
|  |            |               |                   | |
|  +------------+         +-----+-------------------+ |
|                         |                           |
|                         |                           |
+-----------------------------------------------------+
                          |
+-----------------------------------------------------+
|                         |                           |
|  AWS/Packet             | 2. Call API service       |
|                         | with token                |
|        +----------------v--------------------+      |
|        | TipLab Bastion                      |      |
|        | +----------+   +---------+  +-----+ |      |
|        | | Api      |   |OpenVPN  |  | logs| |      |
|        | | SERVICE  |   |SERVICE  |  |     | |      |
|        | +----------+   +---------+  +-----+ |      |
|        +-------------------------------------+      |
+-----------------------------------------------------+
               |            |                |
+-----------------------------------------------------+
|              |3.update    | 3.update       | 4.Trigger
|  TipLab Infra|            |                | tests  |
|              |            |                |        |
|  +-----------v+  +--------v-----+   +------v------+ |
|  |            |  |              |   |             | |
|  |            |  |              |   |             | |
|  |  AGW       |  |   FEG        |   |   NG40      | |
|  |            |  |              |   |             | |
|  |            |  |              |   |             | |
|  +------------+  +--------------+   +-------------+ |
|                                                     |
+-----------------------------------------------------+


##### Components

CircleCI Cloud:

- CircleCI main component
- CircleCI nodes : spinned up on the fly for every job

AWS/Packet:

- Tiplab Bastion: Access point to the tiplab infra   

Tiplab Infra:

- AGW
- FEG
- NG40: TeraVM "controller"

##### Workflow

1. CircleCi spins up a new CircleCI node that will have a token in order to
call the api hosted on the TipLab Bastion

2. CircleCI node call the Tiplab bastion with the token and trigger a test
execution

3. Tiplab bastion can access the TipLab Infra and will update both agw and feg
with the latest available tag using the FabFile (in the future those packages
can be build directly on the CI node and transferred to the components directly
or uploaded to a special TeraVM registry)

3. Tiplab bastion can access the TipLab Infra and will use the fabfile to
Trigger ng40 tests.

##### Pros

- The Safest option
- We control the infra that run the VPN service
- Another layer of authentification using tokens
- We can dig in the option of whitelisting the CI nodes (Not sure if it's
possible)
- TipLab Bastion is reusable for other specific infra
- More Control on logs/access

##### Cons

- New components are thrown in the picture (Tiplab Bastion)
- More dev time
- More maintenance


#### Design 3

Using TeraVM cloud (AWS). No design for this option as I don't know the
feasibility of it.
