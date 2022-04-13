---
id: ipv6_agw
title: (Experimental) Configure IPv6 AGW
hide_title: true
---
# Configure IPv6 address for UE and SGi interfaces

## Overview

AGW has limited support for IPv6 address familty. It allows user to assign
IPv6 address to UE. IPv6 UE datapath is only supported in Non-NAT mode.

## Supported config

1. Non-NAT datapath with only router mode for IPv6.
2. Dual stack: IPv6 is supported in dual stack mode, IPv6 only AGW is not supported.
3. IPv6 Allocation Using IP-POOL
4. PDN type IPv6: IPv6 would assign single IPv6 address to UE.
5. PDN type IPv4v6: AGW would assign IPv4 and IPv6 address to UE.

## Steps to enable IPv6

1. Setup AGW in IPv4 and attached UE to AGW. Validate datapath using traffic test from UE over IPv4 network.
2. Once IPv4 datapath is validated, enable IPv6 config option in mme and spgw. `mme.yml:s1ap_ipv6_enabled` and `spgw.yml:s1_ipv6_enabled`.
3. Restart services on AGW using command: `config_stateless_agw.py flushall_redis`
4. Rerun the IPv4 traffic test from UE to validate IPv4 datapath
5. Once IPv4 datapath is validated with IPv6 enabled, you can open NMS to configure
IPv6 on orc8r.
6. Disable NAT from NMS AGW config and save the config. This would result in AGW services restart. Wait for 5 min and rerun the traffic test to validate Non-NAT datapath
7. Configure IPv6 networking for SGi interface: On NMS, in AGW config, Open ‘Edit Gateway’ dialog box assign IPv6 IP address along with subnet CIDR rage, also configure SGi IPv6 router IPv6 address. and save the configuration ![AGW SGi IPv6 configuration](assets/lte/SGi-IPv6.png?raw=true "AGW SGi IPv6 configuration")
8. SGi IPv6 configuration would result in AGW services restart. Wait for 5 min and rerun the traffic test to validate Non-NAT datapath
9. Validate IPv6 is working on AGW, run `ping6 fb.com` to make sure you have currect configuration on AGW.
10. Configura IP-POOL for UE IPv6 address: Open JSON config for AGW config and add key, value pair for IPv6 block. This step will get NMS support soon. ![AGW SGi IPv6 configuration](assets/lte/IPv6-block-config.png?raw=true "AGW SGi IPv6 configuration")
11. Configure PDN Type in APN config: On NMS, Open the APNs tab. There should be list of APN, update the PDN type of APN to IPv6 or IPv4v6 to enable IPv6 for the PDN. ![APN IPv6 configuration](assets/lte/APN-IPv6-config.png?raw=true "APN IPv6 configuration")
12. Validate IPv6 POOL config: on AGW run `mobility_cli.py list_ipv6_blocks`.

13. Validate APN type: dump subscriber data and check PDN type:

```bash
subscriber_cli.py get IMSI00101XXXXXXXXXX
sid {
  id: "00101XXXXXXXXXX"
}
lte {
  state: ACTIVE

}
network_id {
  id: "..."
}
state {
  ...
}
sub_profile: "default"
non_3gpp {
  ambr {
    max_bandwidth_ul: ...
    max_bandwidth_dl: ...
  }
  apn_config {
    service_selection: "oai.ipv6"
    qos_profile {
      class_id: 9
      priority_level: 15
    }
    ambr {
      max_bandwidth_ul: ...
      max_bandwidth_dl: ...
    }
    pdn: IPv6      # or IPV4V6
  }
}
sub_network {
}
```

12: Validate assigned IPv6 for UEs: Run `mobility_cli.py get_subscriber_table`.

13: If everything looks good, Run IPv6 data traffic test from UE.
