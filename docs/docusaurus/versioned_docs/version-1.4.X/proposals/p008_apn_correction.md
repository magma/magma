---
id: version-1.4.0-p008_apn_correction
title: MME APN Correction
hide_title: true
original_id: p008_apn_correction
---

# Overview

*Status: Accepted*\
*Author: @ymasmoudi*\
*Last Updated: 09/22*


During session establishment, the phone may have an APN configured already.
In this case, this APN will be used as it has a higher precedence over the
one provided by the SIM card.

This feature enables the overriding of UE requested APN with a network
specified APN via IMSI prefix based filtering. Up to 10 IMSI prefix filters
and corresponding APNs to overwrite with can be defined.

Note that the APN correction applies to a federated set up only where HSS
is used.

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

The APN correction, if enabled by setting enable_apn_correction to true, will
overwrite the APN requested by UE for UEs with matching imsi_prefix filter
specified as part of apn_correction_map_list with the value of key apn_override.


```
cat /etc/magma/mme.yml
....
enable_apn_correction: false
apn_correction_map_list:
        - imsi_prefix: "00101"
          apn_override: "magma.ipv4"
```

The configuration will be limited to a maximum of 10 imsi prefix filters.


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

Mconfig proto description will be changed to support APN correction configuration
pending the support of mme config in swagger and NMS.

**ORC8R Change**

Exposing MME configuration to swagger will be supported in a later phase.

**NMS Change**

Exposing APN correction configuration to NMS will be supported in a later phase.

