---
id: p004_apn_correction
title: MME APN Correction 
hide_title: true
---

# Overview

*Status: Accepted*\
*Author: @ymasmoudi*\
*Last Updated: 09/22*


During session establishment, the phone may have an APN configured already.
In this case, this APN will be used as it has a higher precedence over the
one provided by the SIM card.

This behavior requires a feature to override this value at the MME level
based on IMSI prefixes prior to interrogating the HSS and creating the PDN
session.

This document describes the support of APN correction.


## Proposition

The PDN connectivity procedure is used by the UE to request the setup of a
default EPS bearer to a PDN by sending a PDN Connectivity Request message to
the network.

During this phase, the UE may communicate its APN following one of 3 modes:

1/ The UE may set the ESM information transfer flag in the PDN Connectivity
Request to indicate that it has ESM information that needs to be transferred
to the MME after the NAS signalling security has been activated.
In this case, the network initiates the ESM information request procedure in
which the UE can provide the MME with APN.

2/ The UE shall include the requested APN in the standalone PDN Connectivity
Request when requesting connectivity to an additional PDN.

3/ If no APN is included in the ESM Information Response or PDN Connectivity
Request message, then the UE is connected to the default APN.


The APN correction will be implemented upon reception of one of these two
messages, PDN Connectivity Request or ESM Information Response. It will
provide the possibility to filter IMSIs based on a prefix configured in the
mme.yml or mconfig. 
The APN is then overridden with the configured value.

```
cat /etc/magma/mme.yml
....
enable_apn_correction: false
apn_correction_map_list:
        - imsi_prefix: "00101"
          apn_override: "magma.ipv4"
```

The configuration will be limited to a maximum of 10 imsi prefix.


## How We Will Change Magma

**MME Change**

Two structures will be added to contain the mme apn correction configuration    

typedef struct apn_map_s {
  bstring imsi_prefix;
  bstring apn_override;
} apn_map_t;

typedef struct apn_map_config_s {
  int nb;
  apn_map_t apn_map[MAX_APN_CORRECTION_MAP_LIST];
} apn_map_config_t;


In addition, a new function will be introduced to filter and return the overridden
APN based on the MME NAS config.

This function will be called upon reception of the ESM Information Request or PDN
Connectivity Request.


**NMS Change**

Exposing APN configuration to NMS will be supported in a later phase.

