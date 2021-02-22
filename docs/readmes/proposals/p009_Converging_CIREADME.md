[WIP]# Proposal: Converging CI

Author(s): @tmdzk (more to be added)

Last updated:

Discussion at https://github.com/magma/magma/discussions/5063.

## Abstract

The current status of CI is very complex due to all the tool that we use. It is
very confusing for everyone to understand what run where and which tests are
mandatory or not.
The goal of that proposal is to find the best solution to converge all tools
into. That will help maintainers/contributors to take advantage of CI as much as
 they can.

## Background

What do we currently have:

- Github Action: that run simple jobs like DCO-checks or Netlify jobs (magma
  website)
- [Jenkins](https://jenkins.magmacore.org/) : Runs heavy task like integration test
on magma core, this jenkins has a couple of [packet](https://metal.equinix.com/)
(ex Packet), also build all the Vagrant boxes (dev) environment automatically
CircleCi https://app.circleci.com/pipelines/github/magma/magma: that runs every
single jobs: dockerized/on baremetal/in virtualized environment. We hacked
around to plug our own environments using a custom VPN server and a scheduling
api that allow to run remote jobs on Baremetal


## Proposals

Here are the proposed solution to the current problem, those solutions have to
be discussed during the DevOps meeting. No solution has been taken yet. The
proposal is just to get things started.

- Migrate everything to CircleCI

We would use CircleCI for every CI cases: docker environment, virtualized
environment, migrate all equinix nodes from jenins to CircleCI using their
custom runner feature

Pros:

Cons:

- Migrate everything to Jenkins

Pros:

Cons:

- Migrate everything to Github Action

Pros:

Cons:

- Migrate everything to a third Ci tool that we are not currently using
TraviCI, Harness.io

Pros:

Cons:

## Rationale

I'll continue the proposal with the `Migrate everything to Github Action` idea.


## Compatibility

Those changes are not backward compatible. It will require the DevOps Team to
work on that full time, however we can slowly canary, blue/green the new CI
tests to make sure that the new CI system we're not downgrading our current
services, and the new CI can take the load contributors inject to it.


## Implementation

Implementation roadmap and timeline:


- [1 month] Test Github Action for all magma use cases:


## Open issues (if applicable)

No open issue, just a discussion
