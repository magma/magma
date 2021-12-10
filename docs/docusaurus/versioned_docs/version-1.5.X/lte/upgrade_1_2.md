---
id: version-1.5.0-upgrade_1_2
title: Upgrade to v1.2
hide_title: true
original_id: upgrade_1_2
---

# Upgrade to v1.2

You can upgrade your access gateways remotely from the NMS or SSH directly
into them and run an `apt-get install`.

## NMS Autoupgrade

If you've set up your Access Gateways in upgrade tiers already, you can upgrade
them automatically from the NMS.

Navigate to the "Equipment" tab of the NMS and select "Upgrade" on the upper
right. In the next modal window, find the tier that your target AGW is a
member of (probably `default`), then set the desired software version for that
tier to `1.2.0-1600052642-41609608`.
The AGWs in this tier will pull this configuration and upgrade automatically.
By default, we check for an upgrade every 5 minutes.

![1.2.0 upgrade](assets/agw_120_1.png)

![1.2.0 upgrade](assets/agw_120_2.png)

When the gateway completes its upgrade, you should see that its reported
software version in. If it hasn't changed, something probably went wrong with
the autoupgrade. You will probably have to SSH into the AGW to troubleshoot
the installation (see the next section).

If you want to roll out an upgrade slowly to your fleet, you can segment your
AGWs into different tiers and upgrade tiers one-by-one. Use the NMS to create
new tiers and assign your AGWs to them.

## Manual Upgrade

SSH into your target AGW then:

```bash
sudo apt-get update
sudo apt-get install magma=1.2.0-1600052642-41609608
```
