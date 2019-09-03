---
id: config_agw
title: AGW Configuration
sidebar_label: AGW Configuration
hide_title: true
---
# Access Gateway Configuration
### Prerequisites

Before beginning to configure your Magma Access Gateway, you will need to make
sure that it is running all services without crashing. You will also need a
working Orchestrator setup. Please follow the instructions in
"[Deploying Orchestrator](
https://facebookincubator.github.io/magma/docs/orc8r/deploying)" for a
successful Orchestrator installation.

### Configuration
1. Copy root CA, `rootCA.pem`, from the host running the Orchestrator to
`/var/opt/magma/tmp/certs/` on the AGW host.

2. Point your AGW to the right Orc8r instance. Create the file
`/var/opt/magma/configs/control_proxy.yml` with cloud_address and cloud_port
overrides:
```
cloud_address: controller.yourdomain.com
cloud_port: 443
bootstrap_address: bootstrapper-controller.yourdomain.com
bootstrap_port: 443
```
3. On AGW host, run `show_gateway_info.py` to get the AGW's hardware ID and
challenge key.

4. From the Swagger API page, `https://api.yourdomain.com/v1/`, add the
following configurations for your AGW. Use the provided sample for guidance or
check the available config model. Pay extra attention to the cellular
configurations and make sure they match the connected/planned eNodeB(s)

  - **/lte** Create a new LTE network

    Pick a unique `id` and a meaningful `name` and `description` for your
    network. Also, make sure only one of `fdd_config` or `tdd_config`
    configurations is included.

    ```
    {
      "cellular": {
        "epc": {
          "cloud_subscriberdb_enabled": false,
          "default_rule_id": "default_rule_1",
          "lte_auth_amf": "gAA=",
          "lte_auth_op": "EREREREREREREREREREREQ==",
          "mcc": "001",
          "mnc": "01",
          "mobility": {
            "ip_allocation_mode": "NAT",
            "nat": {
              "ip_blocks": [
                "192.168.0.0/16"
              ]
            },
            "reserved_addresses": [
              "192.168.0.1"
            ],
            "static": {
              "ip_blocks_by_tac": {
                "1": [
                  "192.168.0.0/16"
                ],
                "2": [
                  "172.10.0.0/16",
                  "172.20.0.0/16"
                ]
              }
            }
          },
          "network_services": [
            "metering",
            "dpi",
            "policy_enforcement"
          ],
          "relay_enabled": false,
          "sub_profiles": {
            "additionalProp1": {
              "max_dl_bit_rate": 20000000,
              "max_ul_bit_rate": 100000000
            },
            "additionalProp2": {
              "max_dl_bit_rate": 20000000,
              "max_ul_bit_rate": 100000000
            },
            "additionalProp3": {
              "max_dl_bit_rate": 20000000,
              "max_ul_bit_rate": 100000000
            }
          },
          "tac": 1
        },
        "feg_network_id": "",
        "ran": {
          "bandwidth_mhz": 20,
          "tdd_config": {
            "earfcndl": 39150,
            "special_subframe_pattern": 7,
            "subframe_assignment": 2
          }
        }
      },
      "description": "Network Description",
      "dns": {
        "enable_caching": false,
        "local_ttl": 0,
        "records": [
          {
            "a_record": [
              "192.88.99.142"
            ],
            "aaaa_record": [
              "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
            ],
            "cname_record": [
              "cname.example.com"
            ],
            "domain": "example.com"
          }
        ]
      },
      "features": {
        "features": {
          "networkType": "cellular"
        }
      },
      "id": "network_id",
      "name": "Network Name"
    }
    ```

  - **/networks/{network_id}/tiers** Register a tier

    ```
    {
      "gateways": [],
      "id": "default",
      "images": [],
      "name": "Default Tier",
      "version": "0.3.14-123456789-deadbeef"
    }
    ```

  - **/lte/{network_id}/gateways** Register a new LTE gateway

    Pick a unique `id` for your gateway. Also, use the `hardware_id` and `key`
    obtained from step 3 above.

    ```
    {
      "cellular": {
        "epc": {
          "ip_block": "192.168.128.0/24",
          "nat_enabled": true
        },
        "non_eps_service": {
          "arfcn_2g": [
            0
          ],
          "csfb_mcc": "001",
          "csfb_mnc": "01",
          "csfb_rat": 0,
          "lac": 1,
          "non_eps_service_control": 0
        },
        "ran": {
          "pci": 260,
          "transmit_enabled": true
        }
      },
      "connected_enodeb_serials": [],
      "description": "Sample Gateway description",
      "device": {
        "hardware_id": "==== HW ID HERE ====",
        "key": {
          "key": "==== HW KEY HERE ====",
          "key_type": "SOFTWARE_ECDSA_SHA256"
        }
      },
      "id": "gateway_id",
      "magmad": {
        "autoupgrade_enabled": true,
        "autoupgrade_poll_interval": 300,
        "checkin_interval": 60,
        "checkin_timeout": 10,
        "dynamic_services": []
      },
      "name": "Sample Gateway",
      "tier": "default"
    }
    ```

5. Validate connection between AGW and Orchestrator. This can be done by
verifying the proper propagation of AGW configurations or by monitoring syslogs.
  - Check `magmad` logs for successful checkin with Orchestrator
  ```
  journalctl -u magma@magmad -f
  # Look for the following logs
  # INFO:root:Checkin Successful!
  # INFO:root:[SyncRPC] Got heartBeat from cloud
  # INFO:root:Processing config update gateway_id
  ```

  - Verify the proper propagation of AGW configurations in the file
  `/var/opt/magma/configs/gateway.mconfig`. The following provides an example
  for the content of this file.

    ```
      {
        "configsByKey": {
          "control_proxy": {
            "@type": "type.googleapis.com/magma.mconfig.ControlProxy",
            "logLevel": "INFO"
          },
          "dnsd": {
            "@type": "type.googleapis.com/magma.mconfig.DnsD",
            "logLevel": "INFO",
            "enableCaching": false,
            "localTTL": 0,
            "records": [
              {
                "aRecord": [
                  "192.88.99.142"
                ],
                "aaaaRecord": [
                  "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
                ],
                "cnameRecord": [
                  "cname.example.com"
                ],
                "domain": "example.com"
              }
            ]
          },
          "enodebd": {
            "@type": "type.googleapis.com/magma.mconfig.EnodebD",
            "logLevel": "INFO",
            "pci": 260,
            "earfcndl": 0,
            "bandwidthMhz": 20,
            "plmnidList": "00101",
            "subframeAssignment": 0,
            "specialSubframePattern": 0,
            "allowEnodebTransmit": true,
            "tac": 1,
            "csfbRat": "CSFBRAT_2G",
            "arfcn2g": [
              0
            ],
            "tddConfig": {
              "earfcndl": 39150,
              "subframeAssignment": 2,
              "specialSubframePattern": 7
            },
            "fddConfig": null,
            "enbConfigsBySerial": {}
          },
          "magmad": {
            "@type": "type.googleapis.com/magma.mconfig.MagmaD",
            "logLevel": "INFO",
            "checkinInterval": 60,
            "checkinTimeout": 10,
            "autoupgradeEnabled": true,
            "autoupgradePollInterval": 300,
            "packageVersion": "0.3.14-123456789-deadbeef",
            "images": [],
            "tierId": "",
            "featureFlags": {},
            "dynamicServices": []
          },
          "metricsd": {
            "@type": "type.googleapis.com/magma.mconfig.MetricsD",
            "logLevel": "INFO"
          },
          "mme": {
            "@type": "type.googleapis.com/magma.mconfig.MME",
            "logLevel": "INFO",
            "mcc": "001",
            "mnc": "01",
            "tac": 1,
            "mmeGid": 1,
            "mmeCode": 1,
            "enableDnsCaching": false,
            "relayEnabled": false,
            "nonEpsServiceControl": "NON_EPS_SERVICE_CONTROL_OFF",
            "csfbMcc": "001",
            "csfbMnc": "01",
            "lac": 1,
            "cloudSubscriberdbEnabled": false,
            "attachedEnodebTacs": []
          },
          "mobilityd": {
            "@type": "type.googleapis.com/magma.mconfig.MobilityD",
            "logLevel": "INFO",
            "ipBlock": "192.168.128.0/24"
          },
          "pipelined": {
            "@type": "type.googleapis.com/magma.mconfig.PipelineD",
            "logLevel": "INFO",
            "ueIpBlock": "192.168.128.0/24",
            "natEnabled": true,
            "defaultRuleId": "default_rule_1",
            "relayEnabled": false,
            "services": [
              "METERING",
              "DPI",
              "ENFORCEMENT"
            ]
          },
          "policydb": {
            "@type": "type.googleapis.com/magma.mconfig.PolicyDB",
            "logLevel": "INFO"
          },
          "sessiond": {
            "@type": "type.googleapis.com/magma.mconfig.SessionD",
            "logLevel": "INFO",
            "relayEnabled": false
          },
          "subscriberdb": {
            "@type": "type.googleapis.com/magma.mconfig.SubscriberDB",
            "logLevel": "INFO",
            "lteAuthOp": "EREREREREREREREREREREQ==",
            "lteAuthAmf": "gAA=",
            "subProfiles": {
              "additionalProp1": {
                "maxUlBitRate": "100000000",
                "maxDlBitRate": "20000000"
              },
              "additionalProp2": {
                "maxUlBitRate": "100000000",
                "maxDlBitRate": "20000000"
              },
              "additionalProp3": {
                "maxUlBitRate": "100000000",
                "maxDlBitRate": "20000000"
              }
            },
            "relayEnabled": false
          }
        },
        "metadata": {
          "createdAt": "1567027497"
        }
      }
    ```
