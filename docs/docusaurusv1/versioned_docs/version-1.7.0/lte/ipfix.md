---
id: version-1.7.0-ipfix
title: Magma IPFIX Support
hide_title: true
original_id: ipfix
---
# Magma IPFIX Implementation

## Overview

Current AGW supports an OVS based Internet Protocol Flow Information Export(IPFIX) export. Magma does not provide the IPFIX collector, but there are plenty of free open source collectors to use.
Current architecture supports 2 main use cases of IPFIX exports:

- Conntrack mode (recommended)
- Default mode

We sample packets by installing a special OVS flow rule with a sample action, that inspects the packets and creates an IPFIX record for it(or append to an existing record if such exists) and is then periodically sent over to the collector.

## Conntrack mode

In Conntrack mode IPFIX sampling is done only two times for each unique flow.
This is done by leveraging the linux conntrack events via the ConntrackD service and generating an internal packet every time a conntrack event is generated.
The internal packet is recognized by the OVS flowtable as a special packet, gets forwarded to a specific IPFIX sampling table, hits the sampling rule(creating the IPFIX entry), and is then dropped.

**As we only sample 2 packets per flow the performance degredation is minimal, and every flow is captured.**

## Default mode

**This mode is not recommended as very low sampling is best for performance but will miss occasional short lived(low pkt count) flows.**

In default mode IPFIX sampling is done only on a specified percentage of the packets for each flow.

### Custom IPFIX fields

We support a list of custom fields in the IPFIX export:

- imsi
- msisdn
- apn_mac_addr
- apn_name
- app_name

Here is the nProbe config for the special custom fields
Table is as follows:

```text
# Name    STANDARD_ALIAS	PEN     FieldId     Len         Format
imsi            NONE       6888         899       8    dump_as_hex
msisdn          NONE       6888         900      16    dump_as_hex
apn_mac_addr    NONE       6888         901       6    dump_as_hex
apn_name        NONE       6888         902      24    dump_as_asci
app_name        NONE       6888         903       4    dump_as_hex
```

### Initialization

Configure IPFIX collector destination IP/Port from the NMS

Enable IPFIX in pipelined.yml

```yaml
# Add ipfix to static_services

# Sample configuration for IPFIX sampling
ipfix:
  enabled: true
  probability: 65
  collector_set_id: 1
  cache_timeout: 60
  obs_domain_id: 1
  obs_point_id: 1
```

If running with conntrackd:

Add `connectiond` to magma_services in magmad.yml

Add conntrackd configuration to pipelined:

```yaml
#Add conntrack to static_services

# Sample configuration for conntrackd setup
conntrackd:
    enabled: true
    zone: 897
```

## Debug and logs

Internal OVS Check

Check if the openvswitch flows are correctly created. With a user attached and having internet traffic check the ovs tables using the command.

`ovs-ofctl dump-flows gtp_br0`

If IPFix is properly working you should be able to see an entry in the table 203, the number of packets countered in this entry needs to increase as the user is consuming data, see the example below

```sh
ovs-ofctl dump-flows gtp_br0 # Running First Time
cookie=0x0, duration=3256.232s, table=203, n_packets=856, n_bytes=51370, priority=10 actions=sample(probability=65535,collector_set_id=1,obs_domain_id=1,obs_point_id=1,pdp_start_epoch=1,msisdn=default,apn_name=default,sampling_port=gtp0)

ovs-ofctl dump-flows gtp_br0 # Running Second Time
cookie=0x0, duration=3316.850s, table=203, n_packets=870, n_bytes=52210, priority=10 actions=sample(probability=65535,collector_set_id=1,obs_domain_id=1,obs_point_id=1,pdp_start_epoch=1,msisdn=default,apn_name=default8,sampling_port=gtp0)
```

### IPFix packets

Check if the IPFIx packets are leaving the AGW source. In this example the collector IP and port were configured to 1.1.1.1 and 65010 make sure to change those to reflect your pipelined.yml. Run the tshark analyser to check if the IPFix records are shown as expected:

```sh
tshark -i any -d udp.port=65010,cflow dst 1.1.1.1
1 0.000000000 172.17.3.151 → 1.1.1.1 CFLOW 439 IPFIX flow ( 395 bytes) Obs-Domain-ID=    1 [Data:281]
2 1.035991689 172.17.3.151 → 1.1.1.1 CFLOW 439 IPFIX flow ( 395 bytes) Obs-Domain-ID=    1 [Data:281]
```

### nProbe collection

After running nProbe check if it is collecting the flows from the AGW, in this scenario nProbe is running at port 65010, and the AGW uses ip 172.17.3.151, you should see the the following log from your nProbe command

```sh
05/Jul/2021 20:52:03 [collect.c:206] Flow collector listening on port 65010 (IPv4/v6)
05/Jul/2021 20:52:03 [export.c:545] Using TLV as serialization format
05/Jul/2021 20:52:03 [nprobe.c:10773] nProbe started successfully
05/Jul/2021 20:52:06 [collect.c:3038] Collecting flows from 172.17.3.151 [total: 1/4]
```
