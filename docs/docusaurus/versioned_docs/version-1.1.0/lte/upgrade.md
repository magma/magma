---
id: version-1.1.0-agw_110_upgrade
title: Upgrading from 1.0
sidebar_label: Upgrading from 1.0
hide_title: true
original_id: agw_110_upgrade
---
# Upgrading from 1.0

You can upgrade your access gateways remotely from the NMS or SSH directly
into them and run an `apt-get install`.

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

It is recommended that you conduct upgrades on a test prior to rolling it out to your entire network. 

Setting up a test tier is a two part process. First, create a test tier in the "Equipment" tab by clicking the "Upgrade" button. Next, click the "Add" button in the top right corner and enter in your test tier information and click the "Save" button. 

Once this is created, navigate back to the main "Equipment" tab page. Under the Gateways table select the "Actions" button on the row of the AGW you would like to add to the test tier.  Select the "Edit" option to take you into the configuration page and then select "Edit JSON" in the top right. At the bottom on the JSON object change the "tier" to match the name of your test tier. Click "Save" in the top right to complete the process.  

![JSON](assets/JSON_edit.png)

## Manual Upgrade

SSH into your target AGW then:

```bash
sudo apt-get update
sudo apt-get install magma=1.1.0-1589476391-5dbd6822
```
