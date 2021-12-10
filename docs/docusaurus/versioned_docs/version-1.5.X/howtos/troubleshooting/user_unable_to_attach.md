---
id: version-1.5.0-user_unable_to_attach
title: User is unable to attach to Magma
hide_title: true
original_id: user_unable_to_attach
---
# User is unable to attach to Magma

**Description:** This document describe the steps triage issues when a user is unable to attach to the network, AGW is rejecting the attach request.

**Environment:** AGW and Orc8r deployed.

**Affected components:** AGW, Orchestrator

**Prerequisites**:

- APN configuration has been done in Orc8r
- AGW has been successfully  registered in Orc8r
- User has been registered in Orc8r with the correct parameters( Authentication and APN)
- eNB has been successfully integrated with AGW


**Triaging steps**:

1. Run the command sudo mobility_cli.py get_subscriber_table, the output will show the list of subscribers(by IMSI) that  are currently  attached to the network. If the user is showing in the list, it means the attach request was successfully completed. If the IMSI is not showing there, proceed with the next step
2. Verify in the pcap trace or NMS logs the reason of attach reject.
    a. To obtain a pcap trace in AGW you can run the following command `sudo tcpdump -i any sctp -w <destination path>`
    b. To verify the logs in NMS you can go in NMS to Equipment->Select the AGW-> Logs->Search logs


The attach request flow should follow below signal flow, where sctpd, mme, subscriberdb, mobilityd and sessiond are used services in the AGW:

![Attach flow](assets/lte/attach_flow.png)


Below you will find possible causes of why the Attach Request could be rejected by AGW.


**Causes related to invalid messages or Unknown messages**

UE may include in the Attach Request a parameter Magma doesn’t support. Magma will reject the request with the following cause:

```
SEMANTICALLY_INCORRECT 95
INVALID_MANDATORY_INFO 96
MESSAGE_TYPE_NOT_IMPLEMENTED 97
MESSAGE_TYPE_NOT_COMPATIBLE 98
IE_NOT_IMPLEMENTED 99
CONDITIONAL_IE_ERROR 100
MESSAGE_NOT_COMPATIBLE 101
PROTOCOL_ERROR 111
```

Suggested Solution:

1. Identify which parameter the UE is sending in the Attach Request and compare with the ones Magma support. You can find which parameters magma is currently supported from the code. https://github.com/magma/magma/blob/master/lte/gateway/c/oai/tasks/nas/emm/msg/AttachRequest.c

    Note: Make sure to select the branch you are currently have in your network(v1.1, v1.2, v1.3 , etc)

2. Once you identify the parameter, verify with the UE/CPE vendor if you can disable it. If you can't disable, file a feature request.


**Misconfiguration of APN**

Using the pcap trace, verify the APN configured in NMS/Orc8r matches the APN sent by UE.
- You will find the APN sent by UE in the “PDN connectivity request”( Attach Request)  or in the “ESM information transfer”
- Make sure the APN matches the one configured in Orc8r for this user



**Cause related to subscription options**
If the user is not correctly provisioned  in the NMS/Orc8r,  the attach request will be  rejected  with the following causes:

```
EMM_CAUSE_IMEI_NOT_ACCEPTED 5
EMM_CAUSE_EPS_NOT_ALLOWED 7
EMM_CAUSE_BOTH_NOT_ALLOWED 8
EMM_CAUSE_PLMN_NOT_ALLOWED 11
EMM_CAUSE_TA_NOT_ALLOWED 12
EMM_CAUSE_ROAMING_NOT_ALLOWED 13
EMM_CAUSE_EPS_NOT_ALLOWED_IN_PLMN 14
EMM_CAUSE_NO_SUITABLE_CELLS 15
EMM_CAUSE_CSG_NOT_AUTHORIZED 25
EMM_CAUSE_NOT_AUTHORIZED_IN_PLMN 35
EMM_CAUSE_NO_EPS_BEARER_CTX_ACTIVE 40
```

Suggested Solution:

1. Double check with the Auth Key and OPC are correctly set for this user
2. Verify the Auth Key and OPC values with the SIM vendor
3. If issue still persists, please  file github issues or ask in our support channels https://www.magmacore.org/community/
