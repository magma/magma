---
id: version-1.3.0-agw_130_upgrade
title: Upgrade to v1.3
hide_title: true
original_id: agw_130_upgrade
---
# Upgrade to v1.3

You can upgrade your access gateways remotely from the NMS or SSH directly
into them and run an `apt-get install`.

The Access Gateway version needs to be equal to or less than the version
 of your Orc8r. We recommend you update your Orc8r first.
 
## NMS Autoupgrade

If you've set up your Access Gateways in upgrade tiers already, you can upgrade
them automatically from the NMS.

Navigate to the "Equipment" tab of the NMS and select "Upgrade" on the upper
right. In the next modal window, find the tier that your target AGW is a
member of (probably `default`), then set the desired software version for that
tier to `1.3.0-1602477016-723feee0`.
The AGWs in this tier will pull this configuration and upgrade automatically.
By default, we check for an upgrade every 5 minutes.

![1.2.0 upgrade](assets/agw_120_1.png)

![1.2.0 upgrade](assets/agw_130_2.png)

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
sudo apt-get install magma=1.3.0-1602477016-723feee0
```
