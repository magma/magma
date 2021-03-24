---
id: upgrade_1_5
title: Upgrade to v1.5
hide_title: true
---

# Upgrade to v1.5

This guide covers upgrading Orchestrator deployments from v1.4 to v1.5.

The v1.5 upgrade follows a standard procedure:

- Migrate NMS DB data to orc8r DB
- etc.

## NMS DB Data Migration

Run on a machine with access to the kubernetes cluster you are using to run Magma.

*Paste in shell prompt:*

`wget https://gist.githubusercontent.com/andreilee/7aa7d533e2e8f425222b1e6a016a6f5a/raw/286abea80fbec8d03df9949d74f6f2dd298d8222/pre-upgrade-migration.sh && chmod +x pre-upgrade-migration.sh && ./pre-upgrade-migration.sh`
