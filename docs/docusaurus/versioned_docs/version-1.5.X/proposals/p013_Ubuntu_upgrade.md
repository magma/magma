---
id: version-1.5.0-p013_Ubuntu_upgrade
title: AGW Ubuntu upgrade
hide_title: true
original_id: p013_Ubuntu_upgrade
---

# Overview

*Status: Accepted*\
*Author: @pshelar*\
*Last Updated: 03/15*\
*Targeted Release: 1.5*\
*Feedback requested from: @arunuke*

This document cover AGW ubuntu upgrade path

## Goal
Add support for Ubuntu based AGW.

## Introduction

There are multiple reasons to update base OS from current Debian to ubuntu distribution

1. Python 3.5 is EOL on September 13th, 2020. Current Debian does not support the latest python release. pip 21.0
   dropped support for Python 3.5 in January 2021.
2. Newer kernel is required for enabling next generation datapath performance.
3. Orc8r is already on Ubuntu. This would converge the magma platform on ubuntu distribution
4. With newer distribution we can relax AGW kernel dependency using *DKMS packages for OVS kernel module*

## Deployment scenarios for 1.5

### *Fresh install on Ubuntu:*

This is going to be default deployment OS for AGW from 1.5 onwards. Using ubuntu base OS for AGW 1.5 would also result
in smoother AGW updates for future releases.

### *Debian based AGW 1.4 to Debian based AGW 1.5:*

Current production deployment can use current Debian upgrade path for upgrading magma from 1.4 to 1.5. Existing workflow
remains same.

### *Debian based AGW 1.4 to Ubuntu based AGW 1.5:*

This deployment involves two steps first, upgrading underlying OS. OS upgrade can be performed either
inplace upgrade or re-install of Ubuntu OS. Second step installing AGW packages.
There is separate section below to discuss it in details.

## Debian depreciation:

Debian support will be deprecated on AGW 1.5. Default OS distribution for magma is going to be ubuntu. There would be
support for Debian based deployment.
This is subject to change depending on stability of ubuntu support for AGW 1.5.

## Drop Debian support:

Debian OS support would be dropped in AGW 1.6

## Extended AGW 1.5 release cycle:

To ease urgency of ubuntu upgrade, magma community could provide extra minor releases for 1.5. This way existing
Debian based production deployment can continue to use debian based AGW and receive bug fixes for 1.5 release.

## Debian to ubuntu upgrade solutions:

### *Zero downtime:*

1. Setup HA pair for AGW, either on-premises or on cloud, there is separate HA document for details.
2. Fail-over all UEs to the new standby AGW.
3. Follow steps in "Upgrade with single AGW node".
4. All UEs would fallback to updated AGW 1.5.
5. At this point the stand-by AGW can be removed.

### *Upgrade with single AGW node:*

This type of deployment results in downtime.

1. Bring down AGW services
2. Update base OS on AGW from Debian to ubuntu
3. Deploy AGW 1.5
4. Configure AGW
5. Start AGW services

## Major changes for Ubuntu upgrade.
1. Update Python based services to Python 3.8
2. Update OVS binaries from 2.8 to 2.14
3. Use OVS DKMS kernel module package for GTP tunneling to support range of kernels, as of now it supports 4.9.214,
   4.14.171, 4.19.110, 5.4.50 and 5.6.19
4. Add packaging scripts for Ubuntu
5. Update magma artifactory to host packages for Ubuntu
6. Add deployment script to deploy AGW on Ubuntu
7. Add CI pipeline for Ubuntu based AGW
