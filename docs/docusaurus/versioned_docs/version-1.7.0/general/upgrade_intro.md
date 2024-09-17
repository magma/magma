---
id: version-1.7.0-upgrade_intro
title: Introduction
hide_title: true
original_id: upgrade_intro
---

# Upgrade Introduction

This section walks through upgrading different components of the Magma topology. This specific page provides a high-level overview of the upgrade process.

## Orc8r-gateway compatibility

Orc8r and gateways both follow a SemVer-like versioning of `MAJOR.MINOR.PATCH`, with the following constraints

- Gateway version must be <= Orc8r version
- Gateway can diverge at most 1 minor version from Orc8r

For example

- ✅ Gateway 1.4.0, Orc8r 1.4.0
- ✅ Gateway 1.3.0, Orc8r 1.4.0
- ✅ Gateway 1.3.3, Orc8r 1.4.0
- ✅ Gateway 1.3.3, Orc8r 1.4.10
- 🚨 Gateway 1.2.0, Orc8r 1.4.0 (more than 1 minor)
- 🚨 Gateway 1.4.1, Orc8r 1.4.0 (gateway > Orc8r)

## Orc8r-gateway upgrade flow

Based on these [compatibilities](#orc8r-gateway-compatibility), the following upgrade flow is prescribed

1. Upgrade all gateways to Orc8r's current version
2. Upgrade Orc8r: 1 minor version and/or any number of patch versions
3. (Optional) Upgrade all gateways to Orc8r's current version
