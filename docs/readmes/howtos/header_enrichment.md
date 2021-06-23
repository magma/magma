---
id: he_api
title: Header Enrichment
hide_title: true
---

# Header Enrichment

This feature would allow operators to enable header enrichment for UE HTTP traffic. This way AGW could add subscriber
information to HTTP requests. There could be privacy implication of this feature, so operator should check local
laws before using this feature.
Today AGW would add following parameters to HTTP request:
1. IMSI
2. MSISDN

Following are steps to configure Header enrichment

## 1. Enable header enrichment Feature for AGW
There are two option to enable header enrichment:
1. Enable it in pipelineD config file
   You would need to ssh on the AGW and add following line to pipelineD config file '/etc/magma/pipelined.yml'

   ```
   he_enabled: true
   ```
2. Enable via '/LTE/{Network-id}/{Gateway-id}/' API
   You would need to define following parameters for Header enrichment under 'cellular' parameter.
   ```
   "he_config": {
      "enable_header_enrichment": true,

      "enable_encryption": false,

        "he_hash_function": "MD5",
          "hmac_key": "Xs21Ncas87"

        "he_encryption_algorithm": "RC4",
           "encryption_key": "C14r0315v0x",

        "he_encoding_type": "BASE64",
    },
    ```
* enable_header_enrichment: This would add plain text to HTTP requests
* enable_encryption: This would encrypt the data passed in HTTP header. First the data is hashed then encrypted and encoded as per configuration.
* he_hash_function: User Data would be hashed using specified algorithm and key from 'hmac_key'.
  Supported algorithms: MD5, HEX, SHA256
* he_encryption_algorithm: Hashed data is encrypted using specified algorithm and key from 'encryption_key'.
  Supported algorithms: "RC4", "AES256_CBC_HMAC_MD5", "AES256_ECB_HMAC_MD5", "GZIPPED_AES256_ECB_SHA1"
* he_encoding_type: Encode the UE data using specified algorithm.
  Supported algorithms: "BASE64", "HEX2BIN"


## 2. Define Policy Rule
You need to define policy that defines flow along with list of URL that needs
header enrichment.

This can be done via NMS. Use 'traffic/policy' Page define egress policy.
Following example show required parameters for Header enrichment rule.

```
{
  "flow_list": [
    {
      "action": "PERMIT",
      "match": {
        "direction": "UPLINK",
        "ip_dst": {
          "address": "192.168.0.0/16",
          "version": "IPv4"
        },
        "ip_proto": "IPPROTO_TCP",
        "ip_src": {
          "version": "IPv4"
        },
        "tcp_dst": 80
      }
    }
  ],
  "header_enrichment_targets": [
    "abc.com"
  ],
  "id": "test123",
  "priority": 1,
  "rating_group": 1000,
  "tracking_type": "NO_TRACKING"
}
```

### Please note this rule needs to have single flow with following parameters
1. Mandatory field: `action` equals  `permit`.
2. Mandatory field: `Direction` needs to be `uplink`.
3. Mandatory field: `tcp_dst` port needs to be `80 (http)`.
4. Mandatory field: `ip_dst` is IP address range of the HTTP servers.
5. Mandatory field: `header_enrichment_targets` which is list of URLs. Do not
   specify http protocol in URL name, for example `http://abc.com` is not supported,
   but `abc.com` is a valid target.
6. Define Policy tracking according to your requirements.

## 3. Apply the policy to subscribers
Now apply the rule to subscribers. This can be using by
'/LTE/{network-ID}/Subscribers/{subscriber-id}/' API

```
"active_policies": [
    "test123"
],

```
