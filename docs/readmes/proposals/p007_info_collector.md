---
id: p007_info_collector
title: Information collector for debugging
hide_title: true
---

# Information collector for debugging

*Status: In review*
*Feature owner: @tmdzk*
*Feedback requested from: @apad @karthiksubraveti @andreilee*
*Last Updated: 11/06/2020*

## Summary

As magma is open source it's deployed by developers/companies all across the
globe and on different type of hardware/machine. Core maintainers help as much
as they can in order to help people using/deploying magma. However sometimes
some misconfiguration, bugs happen. In order to speed up troubleshooting we
would like to make a simple tool that the user can run on every single Component
of magma in order to share with the community the current status of their
infrastructure/deployment/software.

## Motivation

When a contributor is seeking for help on the community mailing list or slack,
they usually copy paste a bunch of logs or just throw the error message. Most of
the time this message is not enough in order to help/troubleshoot with them.
We want tocreate a tool that will gather all information needed in order to give
the context/logs/system information


## Goals

Create a simple "issue" template like we use on
[magma's github](https://github.com/magma/magma/issues/new?assignees=&labels=bug&template=bug_report.md&title=)
that would be generated directly on the misconfigured/buggy component. On top
of filling that template, the tool will also gather essential information that
will be useful for debugging/troubleshooting.
_Useful information :_

AGW:
- System information (machine type/distribution/space left/`/proc`)
- Necessary packages installed and their version:
As we're piling up commits and new components/upgrading packages we have to know
which version of which component is running.
- system logs:
kern logs/cron/auth
- Service status:
output of system magma@* show
- pcap (see @andreii @karthik)
- Service logs of main components (journalctl output?)
- Execute the python scripts that catch data flows and save output
- openssl output to check certificates validity
- gather magma specific configs

cwag/feg:
- System information (machine type/distribution/space left/`/proc`)
- docker status (running container/age/version..)
- docker logs of containers


orc8r:
- TBD not enough knowledges on our orc8r

## Implementation Phases

#### Phase 1: Finish design of the tool and specs

This is a simple design with for now only a few information collected (at least
for other component than `agw`) It has to be enriched by the team and discussed.

#### Phase 2: Work on the information collector

Implement all steps of the information collector listed above

#### Phase 3: Test it on different infra/Components

This tool has to be tested on different component in order to make sure it's
not crashing on specific infra.

### Design of the implementation

The design is pretty straight forward. It will be a bash script that will
gather all information above store them in a temp dir and tar them.
Different flag can be passed
- `--skip-<action>` For example to skip heavy logs (like `--skip-pcap`)
- `--add-file` to add more files/informations
- `--issue` to fill the [magma's github](https://github.com/magma/magma/issues/new?assignees=&labels=bug&template=bug_report.md&title=) and add it to the
tarfile
