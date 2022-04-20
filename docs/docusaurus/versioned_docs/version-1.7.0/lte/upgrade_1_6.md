---
id: version-1.7.0-upgrade_1_6
title: Upgrade to v1.6
hide_title: true
original_id: upgrade_1_6
---

# Upgrade to v1.6

## Fresh Install Notes

Instructions on fresh gateway installs are provided in the [Install AGW section](deploy_install.md)

## Upgrading from previous releases

Beginning with release v1.5 (Fuji), Magma artifacts are now hosted on the new Magmacore artifactory. Gateways migrating from older Magma releases should run the migration script to update the sources accordingly.

To upgrade an existing AGW, please run the following upgrade script

```bash
wget https://raw.githubusercontent.com/magma/magma/master/lte/gateway/release/upgrade_magma.sh
chmod +x upgrade_magma.sh
./upgrade_magma.sh
```

## Image Version

Ubuntu: `1.6.0-1625592603-1f26ba81`
