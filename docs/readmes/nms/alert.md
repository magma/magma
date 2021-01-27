---
id: alerts
title: Alerts
hide_title: true
---
# Alerts

Alerts are an important part of our NMS. We highly recommend the operators to run their networks with alerts always turned on. Without alerts, it is impossible to debug any potential issues happening on the network in a timely fashion.
In this guide we will discuss

* Viewing Alerts
* Alert receiver configuration
* Alert rules configuration
    * Predefined alerts
    * Custom Alerts
* Troubleshooting

## Viewing alerts

### Top Level Network Dashboard

Alert dashboard displays the current firing alerts in a table tabbed by severity. In each of the columns we additionally display the labels passed along with the alert.
![viewing_alerts_1](assets/nms/userguide/alerts/viewing_alerts1.png)
### Alarm component’s Alert Tab

Alerts can also be viewed from the alert tab in the Alarm table.
![viewing_alerts_2](assets/nms/userguide/alerts/viewing_alerts2.png)
![viewing_alerts_3](assets/nms/userguide/alerts/viewing_alerts3.png)

## Alert Receivers

An alert Receiver is created to push the alert notification in real time so that the operator is notified when the alert is fired. Following example details the steps involved in creating a slack alert receiver and configuring it in NMS.

### Example: Adding Slack Channel as Alert Receiver

**Generate Slack Webhook URL:**

* Create an App: Go to https://api.slack.com/apps?new_app=1 and click on “Create New App”. Enter the App Name and the Slack Workspace.
* Click on “Incoming Webhooks” and change “Active Incoming Webhooks” to On

![alert_recv1](assets/nms/userguide/alerts/alert_recv1.png)
* Scroll down and create a new Webhook by clicking on “Add New Webhook to Workspace”. Select the Slack Channel name.

![alert_recv2](assets/nms/userguide/alerts/alert_recv2.png)
* Copy the “Webhook URL” once it is generated.

![alert_recv3](assets/nms/userguide/alerts/alert_recv3.png)
**Create a new Alert Receiver in NMS:**
![alert_recv4](assets/nms/userguide/alerts/alert_recv4.png)

**Testing the newly added alert receiver**

Add a dummy alert to verify if the alert receiver is indeed working. A dummy alert expression can be constructed
with a PromQL advanced expresssion of `vector(1)` as shown below.

![alert_recv5](assets/nms/userguide/alerts/alert_recv5.png)

Look for the notification on the slack channel
![alert_recv6](assets/nms/userguide/alerts/alert_recv6.png)

## Configuring Alert rules

Alert rule configuration lets us define rules for triggering alerts.

### Predefined Alerts

Magma NMS comes loaded with some default set of alert rules. These default rules aren’t configured automatically. If an operator chooses to use these default rules, they can do it by clicking “Sync Predefined Alerts” button in Alert rule configuration tab. As shown below
![alerts1](assets/nms/userguide/alerts/alerts1.png)
![alerts2](assets/nms/userguide/alerts/alerts2.png)

Currently predefined alerts configure following alerts
* CPU percent on the gateway is running > 75% in last 5 minutes
* Unsuccessful S1 setup in last 5 minutes
* S6A authorization failures in last 5 minutes
* Upon exceptions when bootstrapping a gateway in last 5 minutes
* When services were restarted unexpectedly in last 5 minutes
* When a UE attach resulted in a failure in last 5 minutes
* When there were any service restarts in last 5 minutes

**Note:** Operator will have to go and additionally specify the receiver in each of these alerts when they want to be notified of the alerts as follows,
![alerts3](assets/nms/userguide/alerts/alerts3.png)

### Custom Alert Rules

Operators can create custom alert rules by creating an expression based on metrics.

**Metrics Overview**
Magma gateways collect various metrics at a regular intervals and push them into Orchestrator. Orchestrator stores these metrics in a Prometheus instance. The Prometheus instance along with AlertManager provides us support in querying various metrics on the system and setting alerts based on that.

We currently support following metrics on our Access gateways.

|Metric Identifier	|Description	|Category	|Service	|
|---	|---	|---	|---	|
|	|	|	|	|
|s1_connection	|eNodeB S1 connection status	|eNodeB	|MME	|
|user_plane_bytes_ul	|User plane uplink bytes	|eNodeB	|	|
|user_plane_bytes_dl	|User plane downlink bytes	|eNodeB	|	|
|enodeb_mgmt_connected	|eNodeB management plane connected	|eNodeB	|enodebd	|
|enodeb_mgmt_configured	|eNodeB is in configured state	|eNodeB	|enodebd	|
|enodeb_rf_tx_enabled	|eNodeB RF transmitter enabled	|eNodeB	|enodebd	|
|enodeb_rf_tx_desired	|eNodeB RF transmitter desired state	|eNodeB	|enodebd	|
|enodeb_gps_connected	|eNodeB GPS synchronized	|eNodeB	|enodebd	|
|enodeb_ptp_connected	|eNodeB PTP/1588 synchronized	|eNodeB	|enodebd	|
|enodeb_opstate_enabled	|eNodeB operationally enabled	|eNodeB	|enodebd	|
|enodeb_reboot_timer_active	|Is timer for eNodeB reboot active	|eNodeB	|enodebd	|
|enodeb_reboots	|eNodeb reboots counter	|eNodeB	|enodebd	|
|rcc_estab_attempts	|RRC establishment attempts	|eNodeB	|enodebd	|
|rrc_estab_successes	|RRC establishment successes	|eNodeB	|enodebd	|
|rrc_reestab_attempts	|RRC re-establishment attempts	|eNodeB	|enodebd	|
|rrc_reestab_attempts_reconf_fail	|RRC re-establishment attempts due to reconfiguration failure	|eNodeB	|enodebd	|
|rrc_reestab_attempts_ho_fail	|RRC re-establishment attempts due to handover failure	|eNodeB	|enodebd	|
|rrc_reestab_attempts_other	|RRC re-establishment attempts due to other cause	|eNodeB	|enodebd	|
|rrc_reestab_successes	|RRC re-establishment successes	|eNodeB	|enodebd	|
|erab_estab_attempts	|ERAB establishment attempts	|eNodeB	|enodebd	|
|erab_estab_failures	|ERAB establishment failures	|eNodeB	|enodebd	|
|erab_estab_successes	|ERAB establishment successes	|eNodeB	|enodebd	|
|erab_release_requests	|ERAB release requests	|eNodeB	|enodebd	|
|erab_release_requests_user_inactivity	|ERAB release requests due to user inactivity	|eNodeB	|enodebd	|
|erab_release_requests_normal	|ERAB release requests due to normal cause	|eNodeB	|enodebd	|
|erab_release_requests_radio_resources_not_available	|ERAB release requests due to radio resources not available	|eNodeB	|enodebd	|
|erab_release_requests_reduce_load	|ERAB release requests due to reducing load in serving cell	|eNodeB	|enodebd	|
|erab_release_requests_fail_in_radio_proc	|ERAB release requests due to failure in the radio interface procedure	|eNodeB	|enodebd	|
|erab_release_requests_eutran_reas	|ERAB release requests due to EUTRAN generated reasons	|eNodeB	|enodebd	|
|erab_release_requests_radio_conn_lost	|ERAB release requests due to radio connection with UE lost	|eNodeB	|enodebd	|
|erab_release_requests_oam_intervention	|ERAB release requests due to OAM intervetion	|eNodeB	|enodebd	|
|pdcp_user_plane_bytes_ul	|User plane uplink bytes at PDCP	|eNodeB	|enodebd	|
|pdcp_user_plane_bytes_dl	|User plane downlink bytes at PDCP	|eNodeB	|enodebd	|
|ip_address_allocated	|Total IP addresses allocated	|AGW	|mobilityd	|
|ip_address_released	|Total IP addresses released	|AGW	|mobilityd	|
|s6a_auth_success	|Total successful S6a auth requests	|AGW	|subscriberdb	|
|s6a_auth_failure	|Total failed S6a auth requests	|AGW	|subscriberdb	|
|s6a_location_update	|Total S6a lcoation update requests	|AGW	|subscriberdb	|
|diameter_capabilities_exchange	|Total Diameter capabilities exchange requests	|AGW	|subscriberdb	|
|diameter_watchdog	|Total Diameter watchdog requests	|AGW	|subscriberdb	|
|diameter_disconnect	|Total Diameter disconnect requests	|AGW	|subscriberdb	|
|dp_send_msg_error	|Total datapath message send errors	|AGW	|pipelined	|
|arp_default_gw_mac_error	|Total errors with default gateway MAC resolution	|AGW	|pipelined	|
|openflow_error_msg	|Total openflow error messages received by code and type	|AGW	|pipelined	|
|unknown_pkt_direction	|Counts number of times a packet is missing its flow direction	|AGW	|pipelined	|
|enforcement_rule_install_fail	|Counts number of times rule install failed in enforcement app	|AGW	|pipelined	|
|enforcement_stats_rule_install_fail	|Counts number of times rule install failed in enforcement stats app	|AGW	|pipelined	|
|network_iface_status	|Status of a network interface required for data pipeline	|AGW	|pipelined	|
|	|	|	|	|
|subscriber_icmp_latency_ms	|Reported latency for subscriber in milliseconds	|AGW	|monitord	|
|	|	|	|	|
|magmad_ping_rtt_ms	|Gateway ping metrics in milliseconds	|AGW	|magmad	|
|cpu_percent	|System-wide CPU utilization as percentage over 1 second	|AGW	|magmad	|
|swap_memory_percent	|Percent of memory that can be assigned to processes	|AGW	|magmad	|
|virtual_memory_percent	|Percent of memoery that can be assigned to processes without the system going to swap	|AGW	|magmad	|
|mem_total	|Total memory	|AGW	|magmad	|
|mem_available	|Available memory	|AGW	|magmad	|
|mem_used	|Used memory	|AGW	|magmad	|
|mem_free	|Free memory	|AGW	|magmad	|
|disk_percent	|Percent of disk space used for the volume mounted at root	|AGW	|magmad	|
|bytes_sent	|System-wide network I/O bytes sent	|AGW	|magmad	|
|bytes_received	|System-wide network I/O bytes received	|AGW	|magmad	|
|temperature	|Temperature readings from system sensors	|AGW	|magmad	|
|checkin_status	|1 for checkin success, and 0 for failure	|AGW	|magmad	|
|bootstrap_exception	|Count for exceptions raised by bootstrapper	|AGW	|magmad	|
|unexpected_service_restarts	|Count of unexpected service restarts	|AGW	|magmad	|
|unattended_upgrade_status	|Unattended Upgrade status	|AGW	|magmad	|
|service_restart_status	|Count of service restarts	|AGW	|magmad	|
|enb_connected	|Number of eNodeb connected to MME	|AGW	|MME	|
|ue_registered	|Number of UE registered succesfully	|AGW	|MME	|
|ue_connected	|Number of UE connected	|AGW	|MME	|
|ue_attach	|Number of UE attach success	|AGW	|MME	|
|ue_detach	|Number of UE detach	|AGW	|MME	|
|s1_setup	|Counter for S1 setup success	|AGW	|MME	|
|mme_sgs_eps_detach_indication_sent	|SGS EPS detach indication sent	|AGW	|MME	|
|sgs_eps_detach_timer_expired	|SGS EPS Detach Timer expired	|AGW	|MME	|
|sgs_eps_implicit_detach_timer_expired	|SGS EPS Implicit detach Timer expired	|AGW	|MME	|
|mme_sgs_imsi_detach_indication_sent	|SGS IMSI detach indication sent	|AGW	|MME	|
|sgs_imsi_detach_timer_expired	|SGS IMSI detach timer expired	|AGW	|MME	|
|sgs_imsi_implicit_detach_timer_expired	|SGS IMSI implicit detach timer expired	|AGW	|MME	|
|mme_spgw_delete_session_rsp	|SPGW delete session response	|AGW	|MME	|
|initial_context_setup_failure_received	|Initial context setup failure received	|AGW	|MME	|
|nas_service_reject	|NAS Service Reject	|AGW	|MME	|
|mme_s6a_update_location_ans	|S6a Update location	|AGW	|MME	|
|sgsap_paging_reject	|SGS Paging Reject	|AGW	|MME	|
|duplicate_attach_request	|Duplicate attach request	|AGW	|MME	|
|authentication_failure	|Auth failure	|AGW	|MME	|
|nas_auth_rsp_timer_expired	|NAS Auth response timer expired	|AGW	|MME	|
|emm_status_rcvd	|EMM status received	|AGW	|MME	|
|emm_status_sent	|EMM status sent	|AGW	|MME	|
|nas_security_mode_command_timer_expired	|NAS security mode command timer expired	|AGW	|MME	|
|extended_service_request	|Extended service request	|AGW	|MME	|
|tracking_area_update_req	|Tracking area update request	|AGW	|MME	|
|tracking_area_update	|Tracking area update success	|AGW	|MME	|
|service_request	|Service request	|AGW	|MME	|
|security_mode_reject_received	|Security mode reject received	|AGW	|MME	|
|mme_new_association	|New SCTP association	|AGW	|MME	|
|ue_context_release_command_timer_expired	|UE context release command timer expired	|AGW	|MME	|
|enb_sctp_shutdown_ue_clean_up_timer_expired	|SCTP shutdown UE clean up timer expired	|AGW	|MME	|
|ue_context_release_req	|UE context release request	|AGW	|MME	|
|s1ap_error_ind_rcvd	|S1AP error indication received	|AGW	|MME	|
|s1_reset_from_enb	|S1 reset from eNB	|AGW	|MME	|
|nas_non_delivery_indication_received	|NAS non delivery indication received	|AGW	|MME	|
|spgw_create_session	|SPGW create session success	|AGW	|MME	|
|ue_pdn_connection	|UE PDN connection	|AGW	|MME	|
|ue_pdn_connectivity_req	|UE PDN connectivity request	|AGW	|MME	|
|ue_reported_usage / up	|Reported TX traffic for subscriber / session in bytes	|AGW	|sessiond	|
|ue_reported_usage / down	|Reported RX traffic for subscriber / session in bytes	|AGW	|sessiond	|
|ue_dropped_usage / up	|Reported dropped TX traffic for subscriber / session in bytes	|AGW	|sessiond	|
|ue_dropped_usage / down	|Reported dropped RX traffic for subscriber / session in bytes	|AGW	|sessiond	|
More up to date information might be available from the “Metrics Explorer” in the metrics component.
![alerts4](assets/nms/userguide/alerts/alerts4.png)


**Custom Alert Configuration**

An alert configuration consists of

* Name/Description of the alert
* Alert Definition
* Alert receiver to be notified when the alert is fired (optional)
* Additional labels which can be added to provide more information about the alert. (optional)

Alert definition consists of a metric expression (a [Prometheus PromQL expression](https://prometheus.io/docs/prometheus/latest/querying/basics/)) and a duration attribute which specifies the time for the expression to be true, following which the alert is fired.
![alerts5](assets/nms/userguide/alerts/alerts5.png)
![alerts6](assets/nms/userguide/alerts/alerts6.png)

We can create a custom alert either from a simple expression or an advanced expression.
In case of a simple expression, we can choose a metric from the dropdown and construct an expression based on that
as shown below
For e.g.
![simple_alert](assets/nms/userguide/alerts/simple_alert.png)

In case of an advanced expression([PromQL cheatsheet](https://promlabs.com/promql-cheat-sheet/])) which might involve applying different functions on metric, we can
type the advanced promQL expression directly in the textbox
![advanced_alert](assets/nms/userguide/alerts/advanced_alert.png)


The following examples show how we can create custom alerts on the above mentioned metrics.

**eNB Down Alert:**
This alert will fire if eNB Rf Tx is down in any of the gateways in your network.

Expression:

```
sum by(gatewayID) (enodeb_rf_tx_enabled{networkID="<your network ID>"} < 1)
```

Duration:
10 minutes

**No Connected eNB Alert:**
This alert will fire if the connected eNB count falls to ‘0’ for any of the gateways in your network.

Simple: Select “enb_connected” metric from the dropdown and construct the if statement as **“if enb_connected < 1”**

Advanced Expression:

```
enb_connected{networkID="<your network ID>"} < 1
```

Duration:
5 minutes

**Free Memory is < 10 Alert:**
This alert will fire if the free memory of any of the gateways in your network is less than 10%

Advanced Expression:

```
((1 - avg_over_time(mem_available{networkID="mpk_dogfooding"}[5m]) / avg_over_time(mem_total{networkID="mpk_dogfooding"}[5m])) * 100) > 90
```

Duration:
15 minutes

**High Disk Usage Alert:**
This alert will fire if the disk usage of any of the gateways in your network is more than 80%
Simple Way: Select “disk*_percent*” metric from the dropdown and construct the if statement as **“if disk*_percent* > 80”**
Advanced Expression:

```
(disk_percent{networkID="<your network ID>"}) > 80
```

Duration:
15 minutes

**Attach Success Rate Alert:**
This alert will fire if the attach success rate of any of the gateways in your network is less than 50% for a 3h window.
Expression:

```
(sum by(gatewayID) (increase(ue_attach{action="attach_accept_sent",networkID="<your network ID>"}[3h]))) * 100 / (sum by(gatewayID) (increase(ue_attach{action=~"attach_accept_sent|attach_reject_sent|attach_abort",networkID="<your network ID>"}[3h]))) < 50
```

Duration:
15 minutes

Brief Explanation:
ue_attach metric is tagged with action, networkID labels. Action labels can contain be either "attach_accept_sent", "attach_reject_sent", "attach_abort".
Here we are computing the percentage of the increase in ue_attach counter for a successful attach against all ue_attach actions including rejected and aborted
actions and triggering an alert if the success rate for the ue_attach action is less than 50%.


**Dip in User Plane Throughput Alert:**
This alert will fire if for any of the gateway in your network, the User Plane throughput dips by over 70% when compared day-over-day.

Expression:

```
(sum by(gatewayID) (((rate(pdcp_user_plane_bytes_dl{networkID="<your network ID>"}[1h])) - (rate(pdcp_user_plane_bytes_dl{networkID="<your network ID>"}[1h] offset 1d))) / (rate(pdcp_user_plane_bytes_dl{networkID="<your network ID>"}[1h] offset 1d))) < -0.7)
```

Duration:
15 minutes

Dip in Connected UEs Alert:
This alert will fire if for any of the gateway in your network, connected UEs dip by over 50% when compared day-over-day.

Expression:

```
(ue_connected{networkID="<your network ID>"} - ue_connected{networkID="<your network ID>"} offset 1d) / (ue_connected{networkID="<your network ID>"}) < -0.5
```

Duration:
15 minutes

## Troubleshooting

In case we are having issues with alerts. Logs from the following services will give more information on debugging this further.

```
kubectl logs -n orc8r -l [app.kubernetes.io/component=alertmanager](http://app.kubernetes.io/component=alertmanager)
kubectl logs -n orc8r -l [app.kubernetes.io/component=alertmanager-configurer](http://app.kubernetes.io/component=alertmanager-configurer)
kubectl logs -n orc8r -l app.kubernetes.io/component=prometheus-configurer -c prometheus-configurer
kubectl logs -n orc8r -l app.kubernetes.io/component=prometheus -c prometheus
kubectl logs -n orc8r -l app.kubernetes.io/component=metricsd
```


