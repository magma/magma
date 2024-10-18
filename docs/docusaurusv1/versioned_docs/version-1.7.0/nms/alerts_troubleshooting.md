---
id: version-1.7.0-alerts_troubleshooting
title: Alerts
hide_title: true
original_id: alerts_troubleshooting
---

# Alerts

## S1 Setup Failure

### Description

S1 setup connections are failing within network <network_id>, which means that the eNodeBs are not able to get provisioned with this network.

### Why is this important?

It is one of the key KPI and may impact network deployment/expansion/operations targets.

### Automated recommendation

1. Check the Orchestrator pods and make sure that they are ok. A basic check to ensure that there is no issue in Orchestrator, specifically in terms of logging.
2. Check the s1_set metric on NMS -> Metrics -> Grafana -> Networks -> S1 Setup (1h increase) to verify. You may further check any additional metrics in Grafana as well.
3. Navigate to NMS -> Dashboard, to check for any relevant alerts or events.
4. Navigate to NMS -> Network, to observe that all EPC and RAN parameters are ok as intended.
5. Navigate to NMS -> Equipment -> Gateways, to check if any Gateway is in ‘Bad’ health which is supposed to be ok.
6. Navigate to NMS -> Equipment -> eNodeB, to check if any eNodeB is in ‘Disconnected’ state. Determine which AGW it belongs to. It may be good indication of any faulty AGW.

### Troubleshooting steps

1. If any Orchestrator pods has some issue, please try to resolve that first.
2. Make sure that Gateways are successfully checking in. If they are not checked-in, means that their metrics won’t be reported. Please follow [these steps](https://magma.github.io/magma/docs/howtos/troubleshooting/agw_unable_to_checkin) to troubleshoot.
3. Perform relevant sync/changes if any update has been performed on NMS -> Network.
4. Make sure that configuration for [Gateway](https://magma.github.io/magma/docs/lte/deploy_config_agw), [eNodeB](https://magma.github.io/magma/docs/lte/deploy_config_enodebd) and [APN](https://magma.github.io/magma/docs/lte/deploy_config_apn)has been followed.
5. If any faulty Gateway has been identified, please consider rebooting _enodebd_ service. Please consider rebooting the device (this should be done in minimal traffic duration) if need be.
6. If still not resolved, then capture trace on eth1 interface on issue Gateway(s) to identify the case. Try to analyze and identify the cause. This will give an indication of any parameter inserted by eNodeB causing the issue.
7. If issue persists then get higher level support by providing relevant traces/logs and [additional files](https://magma.github.io/magma/docs/lte/debug_show_tech). Report issue with eNodeB (vendor, firmware etc) and Magma node details along with any issue found in step 6.

### Causes / Effects/ Solutions

| Cause                         |                   Possible Effects                   |                                                                                     Solutions |
| ----------------------------- | :--------------------------------------------------: | --------------------------------------------------------------------------------------------: |
| Metrics not reported.         |              Gateway check-in failure.               |                                                  Perform steps as mentioned above in point 1. |
| Orchestrator Pods have issue  |           Metrics not processed properly.            |                                                                            Troubleshoot pods. |
| Network level config changes. |               Not synced on equipment.               |                                                       Make relevant changes on all equipment. |
| Problematic eNodeB            | eNodeB sending a parameter not interpreted by Magma. |                                            Make sure to get relevant info from eNodeB vendor. |
| Problematic gateway           |         eNodeB not able to attach with them.         | Make sure to follow proper config and troubleshooting as mentioned in point 4,5, and 6 above. |

### How does this affect the SLA?

Major severity. It indicates that functionality of network has been impacted.

### what resource does this affect?

Network, eNodeB, Gateway, Subscribers.

## UE Attach Failure

### Description

S1 setup connections are failing within network <network_id>, which means that the eNodeBs are not able to get provisioned with this network.

### Why is this important?

Subscribers will not be able to access the services.

### Automated recommendation

1. Please check NMS -> Metrics -> Grafana (Network & Gateway) to verify and check additional metrics.
2. Navigate to NMS -> Dashboard, to check for any relevant alerts or events.
3. Navigate to NMS -> Network, to observe that all EPC and RAN parameters are ok as intended.
4. Navigate to NMS -> Equipment -> Gateways, to check if any Gateway is in ‘Bad’ health which is supposed to be ok and UE are coming under this Gateway or any other Gateway. This will narrow us down to issue Gateway.
5. Navigate to NMS -> Equipment -> eNodeB, to check if any eNodeB is in ‘Disconnected’ state.
6. Check all the services in corresponding Gateway(s) under which issue is been reported.

### Troubleshooting steps

1. Perform relevant sync/changes if any update has been performed on NMS -> Network.
2. Make sure that configuration for [Gateway](https://magma.github.io/magma/docs/lte/deploy_config_agw), [eNodeB](https://magma.github.io/magma/docs/lte/deploy_config_enodebd) and [APN](https://magma.github.io/magma/docs/lte/deploy_config_apn)has been followed.
3. Please follow troubleshooting steps from [here](https://magma.github.io/magma/docs/howtos/troubleshooting/user_unable_to_attach) for issue Gateway. Please note the error code.
4. Please check mme logs, verify if the service request failures are coming from a specific user/device/model/firmware.
5. You may use PromQL _ue_attach{networkID=&lt;NetworkID>,result="failure"}_ to isolate further.
6. If required, please consider rebooting the problematic device (this should be done in minimal traffic duration).
7. If still not resolved, then capture trace on eth1 interface on issue Gateway(s) to identify the case.
8. If issue persists then get higher level support by providing relevant traces/logs and [additional files](https://magma.github.io/magma/docs/lte/debug_show_tech). Report issue with eNodeB (vendor, firmware etc) and Magma node details along with any issue found in step 6.

### Causes / Effects/ Solutions

| Cause                                 |             Possible Effects              |                                              Solutions |
| ------------------------------------- | :---------------------------------------: | -----------------------------------------------------: |
| Network level config changes          |         Not synced on equipment.          |                Make relevant changes on all equipment. |
| Configuration issue.                  |           Unexpected parameters           |                   Make sure to follow step 2 as above. |
| Unknown Messages                      |    Magma not supporting those messages    | Check with vendor or file a feature request with Magma |
| APN misconfiguration                  | Not matched with that configured over NMS |                               Check step configuration |
| Cause related to subscription options |           Not matching with SIM           |                                Verify Auth Key and OPC |

### How does this affect the SLA?

Minor when below 90%.
Major when below 80%.

### what resource does this affect?

Network, eNodeB, Gateway, Subscribers.

## Gateway Checkin Failure

### Description

It was observed that AGW is not able to check-in to Orchestrator as described [here](https://magma.github.io/magma/docs/lte/deploy_config_agw).

### Why is this important?

Visibility of AGW will be lost. That means admins cannot perform AGW related configuration, maintenance, and monitoring from NMS.

### Automated recommendation

1. Navigate to NMS -> Dashboard, to check for any relevant alerts or events.
2. Navigate to NMS -> Equipment -> Gateways -> Click on issue AGW to check overview, events, logs etc.
3. Login to AGW and run `sudo checkin_cli.py.`
4. To checkout further logs, you may run `journalctl -u magma@magmad -f.`

### Troubleshooting steps

1. Please follow [this section](https://magma.github.io/magma/docs/howtos/troubleshooting/agw_unable_to_checkin) for troubleshooting steps.
2. If issue persists then get higher level support by providing relevant traces/logs and [additional files](https://magma.github.io/magma/docs/lte/debug_show_tech).

### Causes / Effects/ Solutions

| Cause                                                     |    Possible Effects     |                                                                                               Solutions |
| --------------------------------------------------------- | :---------------------: | ------------------------------------------------------------------------------------------------------: |
| Hostname and ports are changes in control_proxy.yml file. | Can break the check-in. |                                                                           Make sure it is properly set. |
| Location of rootCA.pem is changed.                        | Can break the check-in. | Make sure rootCA.pem is in the correct location defined in rootca_cert (specified in control_proxy.yml) |
| Certificates are expired.                                 | Can break the check-in. |                                                        Revive the certificates in AGW and Orchestrator. |
| Domain in-consistency                                     | Can break the check-in. |                                      Make sure the domain name is consistent across all configurations. |
| Make sure connection is ok.                               | Can break the check-in. |                                                Make sure the connection is ok and ports are ok as well. |

### How does this affect the SLA?

Critical as visibility of AGW is lost from Orchestrator.

### what resource does this affect?

AGW, Orchestrator, Network.

## eNodeB Failure To Connect

### Description

A new eNodeB is not able to connect to Access Gateway.

### Why is this important?

Impacts the new radio deployment.

### Automated recommendation

1. Navigate to NMS -> Dashboard, to check for any relevant alerts or events.
2. Navigate to NMS -> Equipment -> eNodeB-> Click on issue eNodeB to check overview and config.
3. Login to AGW:
   1. Check that enodebd service is running (for certified eNodeB).
   2. Check dnsd logs to observe if IP address has been assigned to eNodeB.
   3. Check basic status of eNodeB using command enodebd_cli.py get_all_status

### Troubleshooting steps

1. Please make sure that IP reachability is ok between eNodeB and AGW.
2. Please make sure to that all the steps mentioned [here](https://magma.github.io/magma/docs/lte/deploy_config_enodebd) are followed properly.
3. If issue still persists then get higher level support by providing relevant traces/logs and [additional files](https://magma.github.io/magma/docs/lte/debug_show_tech).

### Causes / Effects/ Solutions

| Cause           |       Possible Effects        |                                                                              Solutions |
| --------------- | :---------------------------: | -------------------------------------------------------------------------------------: |
| IP reachability |       Issue in routing        | Make sure that connection with AGW is ok (either direct or with any switch in between) |
| Configuration   | Steps have not been followed. |                         Make sure to follow all steps as mentioned in Step 2 of above. |

### How does this affect the SLA?

Major

### what resource does this affect?

AGW, Orchestrator, Network.

## No Metrics Available

### Description

There are no metrics obtained from gateway

### Why is this important?

Unable to monitor the gateway

### Automated recommendation

If AGW is unable to checkin(alert fired), recommend restarting the magmad service on AGW, verify backhaul, or power issues.

### Troubleshooting steps

#### Cause 1: AGW is not reachable

Make sure AGW is reachable, try to ssh the AGW and if not possible make sure there are no backhaul, physical connection or power issues onsite.

#### Cause 2: AGW checkin Orc8r problem

In AGW use `sudo checkin_cli.py` to test the Orc8r connection, verify the correct configuration in control_proxy.yml, validate certificates, verify Orc8r endpoints are reachable from AGW. Use the [troubleshooting steps](https://magma.github.io/magma/docs/howtos/troubleshooting/agw_unable_to_checkin) to investigate these steps further.

#### Cause 3: Unhealthy services in AGW

In AGW, use `sudo service magma@* status`  to verify the services are active in AGW. Services like `magmad, eventd` and `td-agent-bit` are important for reporting metrics, events and logs.

In AGW,  verify metrics are being generated using `service303_cli.py metrics &lt;service>.` For example,  `service303_cli.py metrics magmad`

#### Cause 4: Unhealthy services in Orc8r

Verify if the metrics are not being populated on prometheus or if there is an error querying data in prometheus.

Go to the Org8r swagger API and retrieve some data points from a Prometheus query. `GET /networks/{network_id}/prometheus/query`

In Orc8r, use `kubectl --namespace orc8r get pods,` to verify pods are properly running. Pods should have a status with “running”.

In Orc8r, we can debug the issues by dumping the logs on prometheus, prometheus-configmanager, prometheus-cache and metricsd

```bash
kubectl --namespace orc8r logs -l app.kubernetes.io/component=prometheus-configurer -c prometheus-configurer

kubectl --namespace orc8r logs -l app.kubernetes.io/component=prometheus -c prometheus

kubectl --namespace orc8r logs -l app.kubernetes.io/component=metricsd
helm --debug -n orc8r get values orc8r
```

#### Cause 6: Orc8r endpoints unreachable

Make sure you can reach Orc8r endpoints `api.yourdomain, controller.yourdomain` from NMS, you can ping these domains to verify the same. If you are unable to reach the domains try reaching external IPs of the pods. To get the external IPs use the command `kubectl --namespace orc8r get services`

### Causes / Effects/ Solutions

| Cause                       |                                    Possible Effects                                     |                                                                                                                               Solutions |
| --------------------------- | :-------------------------------------------------------------------------------------: | --------------------------------------------------------------------------------------------------------------------------------------: |
| AGW is not reachable        | Ranging from inability to manage and monitoring the network to complete loss of service |                                                                                     Operator to visit site to recover access to the AGW |
| AGW checkin Orc8r problem   |                    No management or monitoring available for one AGW                    | Multiple possible solutions in [troubleshooting steps](https://magma.github.io/magma/docs/howtos/troubleshooting/agw_unable_to_checkin) |
| Unhealthy services in AGW   |                    No management or monitoring available for one AGW                    |                                                                                Restart magmad service and debug further in service logs |
| Unhealthy services in Orc8r |                  No management or monitoring available for all network                  |                                                                  Debug further from prometheus, prometheus-configured and metricsd pods |
| Orc8r endpoints unreachable |                  No management or monitoring available for all network                  |                                                                                Make sure your domains resolve to the proper service IPs |

### How does this affect the SLA?

Major, as it blocks network monitoring.

### what resource does this affect?

AGW, NMS, Orchestrator.

## S6A Authentication Success Rate

### Description

Ratio of Success events of S6a Authentication, over total events of S6a Authentication

### Why is this important?

Measure unexpected increase of users not able to authenticate

### Automated recommendation

### Troubleshooting steps

#### Cause 1: DIAMETER_ERROR_USER_UNKNOWN (5001)

Verify if s6a_auth_failure has increased due to error code 5001.``This result code shall be sent by the HSS to indicate that the user identified by the IMSI is unknown

#### Cause 2: DIAMETER_ERROR_UNKNOWN_EPS_SUBSCRIPTION (5420)

Verify if s6a_auth_failure has increased due to error code 5420. This result code shall be sent by the HSS to indicate that no EPS subscription is associated with the IMSI.

#### Cause 3: DIAMETER_ERROR_RAT_NOT_ALLOWED (5421)

Verify if s6a_auth_failure has increased due to error code 5421. This result code shall be sent by the HSS to indicate the RAT type the UE is using is not allowed for the IMSI.

#### Cause 4: DIAMETER_ERROR_ROAMING_NOT_ALLOWED (5004)

Verify if s6a_auth_failure has increased due to error code 5004. This result code shall be sent by the HSS to indicate that the subscriber is not allowed to roam within the MME or SGSN area.

### Causes / Effects/ Solutions

| Cause                                          |           Possible Effects           |                                                                                                     Solutions |
| ---------------------------------------------- | :----------------------------------: | ------------------------------------------------------------------------------------------------------------: |
| DIAMETER_ERROR_USER_UNKNOWN (5001)             | User unable to attach to the network | Verify on logs which uses are being rejected. Check the subscription information for the user(s) in Orc8r/HSS |
| DIAMETER_ERROR_UNKNOWN_EPS_SUBSCRIPTION (5420) | User unable to attach to the network | Verify on logs which uses are being rejected. Check the subscription information for the user(s) in Orc8r/HSS |
| DIAMETER_ERROR_RAT_NOT_ALLOWED (5421)          | User unable to attach to the network | Verify on logs which uses are being rejected. Check the subscription information for the user(s) in Orc8r/HSS |
| DIAMETER_ERROR_ROAMING_NOT_ALLOWED (5004)      | User unable to attach to the network | Verify on logs which uses are being rejected. Check the subscription information for the user(s) in Orc8r/HSS |

### How does this affect the SLA?

Minor to Major, depending upon the success ratio. Say as below:
Minor when below 90%.
Major when below 80%.

### what resource does this affect?

AGW.

## Service Request Success Rate

### Description

Ratio of Success events of Service Request, over total events of Service Requests

### Why is this important?

An inactive UE in Idle state is unable to get activated to handle new traffic

### Automated recommendation

### Troubleshooting steps

#### Cause 1: Rejects coming from a new error code

Verify if the service request failures are coming from new error code

```promql
service_request{networkID=<NetworkID>,result="failure"}
```

#### Cause 2: Rejects coming from a single user

From mme logs, verify if the service request failures are coming from a single user or specific device brand/model/firmware. Sometimes, frequent attempts from a new user could degrade the metrics.

#### Cause 3: Rejects coming from a single AGW

Verify if service request failures have increased due to a specific AGW. Verify the configuration/version compared to other AGW. You can use the following PromQL to isolate the AGW \
`service_request{networkID=&lt;NetworkID>,result="failure"}`

### Causes / Effects/ Solutions

| Cause                                |                                Possible Effects                                |                                                                                                                                             Solutions |
| ------------------------------------ | :----------------------------------------------------------------------------: | ----------------------------------------------------------------------------------------------------------------------------------------------------: |
| Rejects coming from a new error code | An inactive UE in Idle state is unable to get activated to handle new traffic  |                                                                                           Get information about the new error code in 3GPP standard.) |
| Rejects coming from a single user    | An inactive UE in Idle state is unable to get activated to handle new traffic  | Identify which type of device the user is using. Verify if there was a change in the device(firmware upgrade) and try to test with previous versions. |
| Rejects coming from a single AGW     | An inactive UE in Idle state is unable to get activated to handle new traffic. |                                             Verify any recent configurations/upgrade that differ from other AGWs. If not, verify the same in the eNB. |

### How does this affect the SLA?

Minor to Major, depending upon the success ratio. Say as below:
Minor when below 90%.
Major when below 80%.

### what resource does this affect?

AGW.

## AGW Reboot

### Description

This alert would be tracking the restart of AGW intentionally triggered by user or triggered because of unknown reason which affects the subscribers served by the gateway.

### Why is this important?

Service cannot be provided to subscribers served by the gateway during the reboot process

### Automated recommendation

If other alerts related to power failure are firing, then the reboot might be related to power fluctuation and operator would be recommended to visit the site.

### Troubleshooting steps

- Check below metrics to confirm if traffic has been affected
    - Number of Connected eNodebs
    - Network of Connected UEs
    - Network of Registered UEs
    - Attach/ Reg attempts
    - Attach Success Rate
    - S6a Authentication Success Rate
    - Service Request Success Rate
    - Session Create Success Rate
    - Upload/Download Throughput
- Use last reboot to list the last logged in users and system last reboot time and date.
- Confirm the same information in /var/log/syslog, logs like kernel loading should indicate the AGW has been rebooted.
- Example of what the log could look like:
- magma kernel: [0.000000] Linux version 4.9.0-9-amd64 (debiankernel@lists.debian.org) (gcc version 6.3.0 20170516 (Debian 6.3.0-18+deb9u1)) #1 SMP Debian 4.9.168-1
- Use this timestamp and compare with the timestamp in the metrics degradation to confirm both events are related
- Verify the commands history matching the timestamp to confirm if AGW was intentionally restarted.
- If it wasn’t intentional then check the power connection and health of device running AGW

### Causes / Effects/ Solutions

| Cause                                                                                                  |                                      Possible Effects                                      |                             Solutions |
| ------------------------------------------------------------------------------------------------------ | :----------------------------------------------------------------------------------------: | ------------------------------------: |
| Command run intentionally/AGW reboot due to an unknown cause (power failure)                           | Gateway service(s)/ whole AGW restart causing disruption of services served by the gateway | If expected then it’s not applicable. |
| However if reboot was happened due to unknown reason, then root cause has to be identified accordingly |

### How does this affect the SLA?

Critical, as AGW will be unavailable.

### what resource does this affect?

It directly affects the gateway, which consequently affects all the subscribers connected to that gateway.

## Expected Service Restart

### Description

This alert would be tracking the restart of services intentionally triggered by user which affects the subscribers served by the gateway (mme, magmad pipelined, sessiond, mobilityd). This could be because of AGW reboot as well.

### Why is this important?

Service cannot be provided to subscribers served by the gateway during the restart process

### Automated recommendation

If AGW reboot alert has fired then service(s) restart would be an outcome of that event

### Troubleshooting steps

- Check below metrics to confirm if traffic has been affected
    - Number of Connected eNBs
    - Network of Connected UEs
    - Network of Registered UEs
    - Attach/ Reg attempts
    - Attach Success Rate
    - S6a Authentication Success Rate
    - Service Request Success Rate
    - Session Create Success Rate
    - Upload/Download Throughput
- Use last reboot to list the last logged in users and system last reboot time and date.
- Confirm the same information in /var/log/syslog, logs like kernel loading should indicate the AGW has been rebooted.
- Example of what the log could look like:
- magma kernel: [0.000000] Linux version 4.9.0-9-amd64 (debiankernel@lists.debian.org) (gcc version 6.3.0 20170516 (Debian 6.3.0-18+deb9u1)) #1 SMP Debian 4.9.168-1
- Use this timestamp and compare with the timestamp in the metrics degradation to confirm both events are related
- Verify the commands history in syslog matching the timestamp to confirm if the services were intentionally restarted.

### Causes / Effects/ Solutions

| Cause                                |                                      Possible Effects                                      |                             Solutions |
| ------------------------------------ | :----------------------------------------------------------------------------------------: | ------------------------------------: |
| Command run intentionally/AGW reboot | Gateway service(s)/ whole AGW restart causing disruption of services served by the gateway | If expected then it’s not applicable. |

### How does this affect the SLA?

Minor, as was done intentionally. It serves as more of a notification.

### what resource does this affect?

It directly affects the gateway, which consequently affects all the subscribers connected to that gateway.

## Unexpected Service Restart

### Description

This alert would be tracking the unexpected restart of services which affects the subscribers served by the gateway (mme, magmad pipelined, sessiond, mobilityd).

### Why is this important?

Service cannot be provided to subscribers served by the gateway.

### Automated recommendation

If AGW reboot alert has fired then service(s) restart would be an outcome of that event

### Troubleshooting steps

- Check below metrics to confirm if traffic has been affected
    - Number of Connected eNBs
    - Network of Connected UEs
    - Network of Registered UEs
    - Attach/ Reg attempts
    - Attach Success Rate
    - S6a Authentication Success Rate
    - Service Request Success Rate
    - Session Create Success Rate
    - Upload/Download Throughput
- Check for recent changes done before the issue was first observed
- If applicable, revert the recent changes and check if issue is still observed
- Capture the service crash syslogs and coredumps. Use the approximate time in metrics and look for the events in both syslogs and coredumps
- In syslogs located in /var/log, look for service terminating events and its previous logs. For example, in below mme service crash there is a segfault reported in mme service before its being terminated.
- Dec 5 22:25:55 magma kernel: [266759.489500] ITTI 3[13887]: segfault at 1d6d80 ip 000055b0080da0c2 sp 00007f529e6c0310 error 4 in mme[55b0077bd000+e79000]
- Dec 5 22:25:59 magma systemd[1]: magma@mme.service: Main process exited, code=killed, status=11/SEGV
- Service crashes with a segmentation fault will create coredumps in /tmp/ folder. Verify if coredumps have been created and obtain the coredump that matches the time of the outage/crash. Depending on the type of service crash the name of the coredump will vary. More detail in [https://magma.github.io/magma/docs/lte/dev_notes#analyzing-coredumps](https://magma.github.io/magma/docs/lte/dev_notes#analyzing-coredumps)
- Using latest mme binary, inspect coredumps via gdb and backtrace to find possible root cause
- **Obtain event that triggered the crash**. Every time a service restarts it will generate a log file (i.e. mme.log). Inside the coredump folder you will find the log (i.e. mme.log) that was generated just before the crash. To understand what was the event that triggered the crash, get the last event (Attach Request, Detach, timer expiring, etc.) in the log file.
- Note: If you can't find the timestamp of the crash in syslogs, you can use the last log generated in the log found in the coredump to get the exact timestamp of that crash.
- Investigate or seek for help. Use the collected information to investigate previous Github issues/bugs and confirm if this is a known issue or bug that has been fixed in a later version. Otherwise, open a new report with the information collected.

### Causes / Effects/ Solutions

| Cause                                          |      Possible Effects      |                                                                                            Solutions |
| ---------------------------------------------- | :------------------------: | ---------------------------------------------------------------------------------------------------: |
| Unsupported configuration in partner’s network | Gateway Crash/Restart Loop | If applicable revert the changes, check for known bugs for this version and get higher level support |

### How does this affect the SLA?

Major, as was done unintentionally. Investigation required to check why.

### what resource does this affect?

It directly affects the gateway, which consequently affects all the subscribers connected to that gateway.
