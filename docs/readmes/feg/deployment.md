---
id: deployment
title: Packaging, Deployment, and Upgrades
hide_title: true
---
# Packaging, Deployment, and Upgrades

All necessary federated gateway components are packaged using a fabric
command located at `magma/feg/gateway/fabfile.py`. To run this command:

```console
HOST [magma]$ cd magma/feg/gateway
HOST [magma/feg/gateway]$ fab package
```

This command will create a zip called `magma_feg_<hash>.zip` that is
pushed to S3 on AWS. It can then be copied from S3 and installed on any host.

## Installation

To install this zip, run:

```console
INSTALL-HOST [/home]$ mkdir -p /tmp/images/
INSTALL-HOST [/home]$ cp magma_feg_<hash>.zip /tmp/images
INSTALL-HOST [/home]$ cd /tmp/images
INSTALL-HOST [/tmp/images]$ sudo unzip -o magma_feg_<hash>.zip
INSTALL-HOST [/tmp/images]$ sudo ./install.sh
```

After this completes, you should see: `Installed Succesfully!!`

## Upgrades

If running in an Active/Standby configuration, the standard procedure for
upgrades is as follows:

1. Find which gateway is currently standby
2. Stop the services on standby gateway
3. Wait 30 seconds
4. Upgrade standby gateway
5. Stop services on active gateway
6. Wait 30 seconds (standby will get promoted to active)
7. Upgrade (former) active gateway

Please note that this sequence will lead to an outage for 30-40 seconds.
