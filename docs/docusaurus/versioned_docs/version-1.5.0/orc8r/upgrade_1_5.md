---
id: version-1.5.0-upgrade_1_5
title: Upgrade to v1.5
hide_title: true
original_id: upgrade_1_5
---

# Upgrade to v1.5

This guide covers upgrading Orchestrator deployments from v1.4 to v1.5.

The v1.5 upgrade follows a standard procedure:

- Migrate NMS DB data to orc8r DB
- etc.

## NMS DB Data Migration

Run on a machine with access to the kubernetes cluster you are using to run Magma.

*Paste in shell prompt:*

`wget https://raw.githubusercontent.com/magma/magma/master/nms/app/packages/magmalte/scripts/fuji-upgrade/pre-upgrade-migration.sh && chmod +x pre-upgrade-migration.sh && ./pre-upgrade-migration.sh`
