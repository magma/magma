---
id: version-1.4.0-agw_140_upgrade
title: Upgrading from 1.3
hide_title: true
original_id: agw_140_upgrade
---
# Upgrading from 1.3

You can upgrade your access gateways doing SSH directly into them and
run an `apt-get upgrade`, with some required modification steps.

## Manual Upgrade

As of v1.3.2, the apt source needs to be updated in order to get the latest
tagged AGW build. Hence, it is required to modify
/etc/apt/sources.list.d/packages_magma_etagecom_io.list on the gateway.
Instead of “stretch stretch-1.3.3 main", please replace with
"stretch-1.4.0 main".

Upgrade from v1.3.x to 1.4.0 for the Access Gateway will require operators to
use a force option else the upgrade may fail due to unmet dependencies. To
mitigate this issue, it is recommended that the upgrade be done in the
following fashion:

```bash
sudo apt-get update
sudo apt-get upgrade magma -o Dpkg::Options::="--force-overwrite"
```
