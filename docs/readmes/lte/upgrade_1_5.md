---
id: upgrade_1_5
title: Upgrade to v1.5
hide_title: true
---

# Upgrade to v1.5

> **_NOTE:_** Please note that Fuji is the last release with support for Debian.

### Repo Change

In Magma Fuji (v1.5), Magma artifacts are now hosted on the new Magmacore repositories at 
[artifactory.magmacore.org](https://artifactory.magmacore.org/).
Gateways migrating from older Magma releases can run the migration script to update the sources accordingly.

The repository currently supports both Debian and Ubuntu OS flavors.

```bash
wget https://raw.githubusercontent.com/magma/magma/master/lte/gateway/release/upgrade_magma.sh
./upgrade_magma_sh
```

### Image Version

`1.5.0-1619628161-f023455f`
