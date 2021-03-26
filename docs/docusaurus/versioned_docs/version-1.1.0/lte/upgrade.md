---
id: version-1.1.0-agw_110_upgrade
title: Upgrade to  v1.1
sidebar_label: Upgrade to v1.1
hide_title: true
original_id: agw_110_upgrade
---
# Upgrade to v1.1

You can upgrade your access gateways remotely from the NMS or SSH directly
into them and run an `apt-get install`.

The Access Gateway version needs to be equal to or less than the version
 of your Orc8r. We recommend you update your Orc8r first. 

## NMS Autoupgrade

Navigate to the "Configure" tab of the NMS and select the tab "Upgrade". Find
the tier that your target AGW is a member of (probably `default`), then set
the desired software version for that tier to `1.1.0-1589476391-5dbd6822`.
The AGWs in this tier will pull this configuration and upgrade automatically.
By default, we check for an upgrade every 5 minutes.

![1.1.0 upgrade](assets/agw_110_upgrade.png)

When the gateway completes its upgrade, you should see that its reported
software version in this tab has changed. If it hasn't changed, something
probably went wrong with the autoupgrade. You will probably have to SSH into
the AGW to troubleshoot the installation (see the next section).

If you want to roll out an upgrade slowly to your fleet, you can segment your
AGWs into different tiers and upgrade tiers one-by-one. Use the NMS to create
new tiers and assign your AGWs to them.

## Manual Upgrade

SSH into your target AGW then:

```bash
sudo apt-get update
sudo apt-get install magma=1.1.0-1589476391-5dbd6822
```
