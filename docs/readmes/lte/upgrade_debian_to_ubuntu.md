---
id: upgrade_debian_to_ubuntu
title: Migrate AGWs to Ubuntu
hide_title: true
---

# Migrate AGWs from Debian to Ubuntu

> NOTE: Only Ubuntu-based Access Gateways are supported starting from Magma release 1.6, Debian support is deprecated.

For a smooth transition from a Debian to Ubuntu access gateway, back up local configs including

- Gateway's hardware ID
- Gateway's long-term challenge key (`gw_challenge.key`)
- Gateway's control_proxy config (`control_proxy.yml`)
- Root CA certificate (`rootCA.pem`)
- Any static IPs or routes provisioned directly on the existing gateway

## Before migration

Copy the contents from the following file locations to a local backup storage:

- Gateway's hardware ID: `/etc/snowflake`
- Challenge key: `/var/opt/magma/certs/gw_challenge.key`
- Control proxy config: `/var/opt/magma/configs/control_proxy.yml`
- Root CA certificate : `/var/opt/magma/tmp/certs/rootCA.pem`

After copying these, gather the output of `show_gateway_info.py` as reference and follow the [AGW Install Guide](./deploy_install.md).

## After migration

- Once Magma is installed on the Access Gateway, copy the contents from local backup to configure the new access gateway
- Restart `magmad` services to enable the new configs: `service magma@magmad restart`
- The output of `show_gateway_info.py` should match the one captured before migration
- Verify access gateway is checked-in to the existing Orchestrator and configs are being streamed successfully: `journalctl -fu magma@magmad`
