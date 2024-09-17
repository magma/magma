## Sanity Tests

### TC001_SANITY_control_50UE_3rate.py

**Summary -**
Sanity test case that setups up the S1 connection and `attaches` a low number of users to the AGW.

**Setup -**

- Total Subs - 50
- Total eNBs - 12
- Attach Rate - 3 UEs/sec
- Iterations - 1
- Runtime - Less than 1 minute

**Validation Categories -**

- S1-AP
- UE_STATE_CHECK

### TC002_SANITY_data_30UE_3rate_active_idle.py

**Summary -**
Sanity test case that setups up the S1 connection and `attaches` a moderate number of users to the AGW with low data traffic. Further, the test will force devices to transition through active and idle states in a loop.

**Setup -**

- Total Subs - 30
- Total eNBs - 12
- Attach Rate - 3 UEs/sec
- Data Message Flow - HTTP_DL_500K_Per_UE
- Iterations - 4
- Runtime - Less than 10 minutes

**Validation Categories -**

- EMM
- ESM
- active_idle
- UE_STATE_CHECK

### TC003_SANITY_data_50UE_3rate_25M_attach_detach.py

**Summary -**
Sanity test case that setups up the S1 connection and `attaches` a low number of users to the AGW with low data traffic.

**Setup -**

- Total Subs - 50
- Total eNBs - 12
- Attach Rate - 3 UEs/sec
- Data Message Flow - HTTP_DL_500K_Per_UE
- Iterations - 1
- Runtime - Less than 3 minutes

**Validation Categories -**

- S1-AP
- EMM
- ESM
- Data Traffic
- UE_STATE_CHECK

### TC004_SANITY_control_200UE_3rate.py

**Summary -**
Sanity test case that setups up the S1 connection and `attaches` a moderate number of users to the AGW.

**Setup -**

- Total Subs - 200
- Total eNBs - 12
- Attach Rate - 3 UEs/sec
- Iterations - 1
- Runtime - Less than 2 minutes

**Validation Categories -**

- S1-AP
- EMM
- ESM
- UE_STATE_CHECK

### TC005_SANITY_data_200UE_3rate_100M_attach_detach.py

**Summary -**
Sanity test case that setups up the S1 connection and `attaches` a moderate number of users to the AGW with moderate data traffic.

**Setup -**

- Total Subs - 200
- Total eNBs - 12
- Attach Rate - 3 UEs/sec
- Data Message Flow - HTTP_DL_500K_Per_UE
- Iterations - 1
- Runtime - Less than 2 minutes

**Validation Categories -**

- S1-AP
- EMM
- ESM
- Data Traffic
- UE_STATE_CHECK

### TC006_SANITY_data_200UE_3rate_400M_attach_detach.py

**Summary -**
Sanity test case that setups up the S1 connection and `attaches` a moderate number of users to the AGW with high data traffic.

**Setup -**

- Total Subs - 200
- Total eNBs - 12
- Attach Rate - 3 UEs/sec
- Data Message Flow - HTTP_DL_2M_Per_UE
- Iterations - 1
- Runtime - Less than 2 minutes

**Validation Categories -**

- S1-AP
- EMM
- ESM
- Data Traffic
- UE_STATE_CHECK

### TC007_SANITY_control_400UE_3rate.py

**Summary -**
Sanity test case that setups up the S1 connection and `attaches` a high number of users to the AGW.

**Setup -**

- Total Subs - 400
- Total eNBs - 12
- Attach Rate - 3 UEs/sec
- Iterations - 1
- Runtime - Less than 3 minutes

**Validation Categories -**

- S1-AP
- EMM
- ESM
- UE_STATE_CHECK

### TC008_SANITY_data_400UE_3rate_200M_attach_detach.py

**Summary -**
Sanity test case that setups up the S1 connection and `attaches` a moderate number of users to the AGW with high data traffic.

**Setup -**

- Total Subs - 400
- Total eNBs - 12
- Attach Rate - 3 UEs/sec
- Data Message Flow - HTTP_DL_500K_Per_UE
- Iterations - 1
- Runtime - Less than 3 minutes

**Validation Categories -**

- S1-AP
- EMM
- ESM
- Data Traffic
- UE_STATE_CHECK

### TC009_SANITY_control_600UE_3rate.py

**Summary -**
Sanity test case that setups up the S1 connection and `attaches` a max number of users to the AGW.

**Setup -**

- Total Subs - 600
- Total eNBs - 12
- Attach Rate - 3 UEs/sec
- Iterations - 1
- Runtime - Less than 5 minutes

**Validation Categories -**

- S1-AP
- EMM
- ESM
- UE_STATE_CHECK

### TC010_SANITY_data_600UE_3rate_300M_attach_detach.py

**Summary -**
Sanity test case that setups up the S1 connection and `attaches` a high number of users to the AGW with high data traffic.

**Setup -**

- Total Subs - 600
- Total eNBs - 12
- Attach Rate - 3 UEs/sec
- Data Message Flow - HTTP_DL_500K_Per_UE
- Iterations - 1
- Runtime - Less than 6 minutes

**Validation Categories -**

- S1-AP
- EMM
- ESM
- Data Traffic
- UE_STATE_CHECK

### TC011_SANITY_data_3UE_600M_attach_detach.py

**Summary -**
Sanity test case that setups up the S1 connection and `attaches` a low number of users to the AGW with high data traffic. With APN AMBR configured at 100mbps. validation check is to make sure per user traffic is in accordance with apn ambr

**Setup -**

- Total Subs - 3
- Total eNBs - 12
- Attach Rate - 3 UEs/sec
- Data Message Flow - UDP_DL_300M_Per_UE
- Iterations - 1
- Runtime - Less than 6 minutes
- apn ambr - 100mbps

**Validation Categories -**

- S1-AP
- EMM
- ESM
- Data Traffic ~ 300mbps
- UE_STATE_CHECK

### TC012_SANITY_data_5rate_600UE_400M_attach_detach.py

**Summary -**
Sanity test case that setups up the S1 connection and `attaches` a high number of users to the AGW with high data traffic and high attach rate. This TC is purely to test congestion control on the AGW.

**Setup -**

- Total Subs - 600
- Total eNBs - 12
- Attach Rate - 5 UEs/sec
- Data Message Flow - HTTP_DL_500K_Per_UE
- Iterations - 1
- Runtime - Less than 6 minutes

**Validation Categories -**

- S1-AP
- EMM
- ESM
- Data Traffic ~ 300mbps
- UE_STATE_CHECK

### TC001_STATIC_multi_apn_100UE_500k.py

**Summary -**
Sanity test case that setups up the S1 connection and `attaches` a moderate number of users to the AGW with low data traffic. This test connects 2 APNs per device. One APN uses static IP allocation while the second APN uses IP Addressing over DHCP. These APNs run concurrently.

**Setup -**

- Total Subs - 100
- Total eNBs - 12
- Attach Rate - 3 UEs/sec
- Data Message Flow - HTTP_DL_250K_Per_UE (Double this as this as this is traffic **per** APN )
- Iterations - 1
- Runtime - Less than 10 minutes

**Validation Categories -**

- Data Traffic

### TC002_STATIC_IP_300UE_500k.py

**Summary -**
Sanity test case that setups up the S1 connection and `attaches` a moderate number of users to the AGW with low data traffic. This test connects 3 sets of devices with a different APNs per set. One APN uses static IP allocation, the second APN uses IP Addressing over DHCP, finally the third APN assigns IPs statically but router mode is used with this APN. These APNs run concurrently.

**Setup -**

- Total Subs - 100
- Total eNBs - 12
- Attach Rate - 3 UEs/sec
- Data Message Flow - HTTP_DL_250K_Per_UE (Double this as this as this is traffic **per** APN )
- Iterations - 1
- Runtime - Less than 10 minutes

**Validation Categories -**

- EMM
- ESM
- Data Traffic

## Feature Tests

### TC001_FEATURE_control_gtpu_echo.py

**Summary -**
Feature test case that setups up S1 connections and `attaches` a low number of users to the AGW with no data traffic. The UEs are then kept idle until the GTP-U echo message is sent from the ENBs. Echo messages are sent every minute as long and only sent when no other packets reset this timer.

**Setup -**

- Total Subs - 100
- Total eNBs - 55
- Attach Rate - 3 UEs/sec
- Iterations - 1
- Runtime - Less than 5 minutes

**Validation Categories -**

- S1-AP
- UE_STATE_CHECK
- ENB

### TC002_FEATURE_X2_HO_200UE_6enbs_500k_600sec.py

**Summary -**
Feature test case that setups up the S1 connection and `attaches` a moderate number of users to the AGW with high data traffic. This test case also facilitates X2 handover between eNBs.

**Setup -**

- Total Subs - 200
- Total eNBs - 8 (2 additional for Target)
- Attach Rate - 3 UEs/sec
- Data Message Flow - HTTP_DL_500K_Per_UE
- Iterations - 1
- Runtime - Less than 15 minutes

**Validation Categories -**

- EMM
- ESM
- Data Traffic
- Handover

### TC004_FEATURE_S1_HO_200UE_6enbs_500k_600sec.py

**Summary -**
Feature test case that setups up the S1 connection and `attaches` a moderate number of users to the AGW with high data traffic. This test case also facilitates S1 handover between eNBs.

**Setup -**

- Total Subs - 200
- Total eNBs - 8 (2 additional for Target)
- Attach Rate - 3 UEs/sec
- Data Message Flow - HTTP_DL_500K_Per_UE
- Iterations - 1
- Runtime - Less than 15 minutes

**Validation Categories -**

- EMM
- ESM
- Data Traffic
- Handover

### TC005_FEATURE_header_enrichment_wo_ciphering.py

**Summary -**
Feature test case that setups up the S1 connection and `attaches` a moderate number of users to the AGW with high data traffic. In this test case, UEs make a HTTP get request to a local webserver and expect the HTTP GET header to be enriched with UE specific information. This TC does not encrypt the enriched header.

**Setup -**

- Total Subs - 100
- Total eNBs - 12
- Attach Rate - 5 UEs/sec
- Data Message Flow - HTTP_GET_HEADER_ENRICHMENT
- Iterations - 1
- Runtime - Less than 5 minutes

**Validation Categories -**

- S1-AP
- EMM
- ESM
- UE_STATE_CHECK
- HE

### TC006_FEATURE_APN_COR_50UE.py

**Summary -**
Feature test case that setups up the S1 connection and `attaches` a low number of users to the AGW. In this test case, UEs send an APN value in the ESM transfer message which is then overridden by the AGW (MME). Expectation is that all UE IMSIs configured for APN correction are successfully attached after the MME corrects the APN value requested by the UE.

**Setup -**

- Total Subs - 50 (10 with APN correction; 40 will get rejected)
- Total eNBs - 12
- Attach Rate - 3 UEs/sec
- Iterations - 1
- Runtime - Less than 2 minutes

**Validation Categories -**

- S1-AP
- EMM
- ESM
- UE STATE CHECK

### TC007_FEATURE_PLMN_Restric_50UE.py

**Summary -**
Feature test case that setups up the S1 connection and `attaches` a low number of users to the AGW. In this test case, two sets of UEs attempt to attach; UE SET 1 (PLMN 00101) is allowed to attach while UE SET 2 with (PLMN 00102) is rejected. Success in this test case is to note that only those with the correct PLMN attach while the rest are rejected.

**Setup -**

- Total Subs - 100 (50 per PLMN)
- Total eNBs - 12
- Attach Rate - 3 UEs/sec
- Iterations - 1
- Runtime - Less than 1 minute

**Validation Categories -**

- EMM
- ESM
- UE STATE CHECK

### TC008_FEATURE_QoS_flow_res_200UE_tcp.py

**Summary -**
Feature test case that setups up the S1 connection and `attaches` a low number of users to the AGW. In this test case, two sets of UEs attempt to attach; QOS Silver (Max 500Kbps) and QOS Gold (Max 1Mbps). Success in this test case is to note that the total aggregate tput observed is in line with these QOS profiles. Left side of the equation should yield a total of ~ 136Mbps.

**Setup -**

- Total Subs - 200 (100 per QOS profile)
- Total eNBs - 12
- Attach Rate - 3 UEs/sec
- Iterations - 1
- Runtime - Less than 12 minutes

**Validation Categories -**

- S1AP
- EMM
- ESM
- Data Traffic

### TC009_FEATURE_IMEI_Restric_200UE.py

**Summary -**
Feature test case that setups up the S1 connection and `attaches` a low number of users to the AGW. In this test case, we test the IMEI restriction feature. In UE SET 1 we restrict a specific IMEI while in UE SET 2 we block the entire TAC (Type Allocation Code). 99 successes are expected out of 200 Attach requests.

**Setup -**

- Total Subs - 200 (100 per UE Set)
- Total eNBs - 12
- Attach Rate - 3 UEs/sec
- Iterations - 1
- Runtime - Less than 2 minutes

**Validation Categories -**

- EMM
- ESM
- UE_STATE_CHECK

### TC010_FEATURE_ipfix.py

**Summary -**
Feature test case that setups up the S1 connection and `attaches` a moderate number of users to the AGW. In this test case the UEs all send a HTTP get request towards a server and the AGW generates IPFIX records for each 5 tuple created (per direction). TC is successful when IPFIX records for all 200 UEs are accounted for.

**Setup -**

- Total Subs - 200
- Total eNBs - 12
- Attach Rate - 3 UEs/sec
- Iterations - 1
- Runtime - Less than 12 minutes

**Validation Categories -**

- S1-AP
- EMM
- ESM
- UE_STATE_CHECK
- IPFIX

## Performance Tests

### TC001_PERFORMANCE_data_200UE_2M.py

**Summary -**
Performance test case that setups up the S1 connection and `attaches` a moderate number of users to the AGW with high traffic. This test case runs for ~8hours and is supposed to provide signal on platform stability.

**Setup -**

- Total Subs - 200
- Total eNBs - 12
- Attach Rate - 3 UEs/sec
- Data Message Flow - HTTP_DL_2M_Per_UE
- Iterations - 1
- Runtime - 8 Hours

**Validation Categories -**

- S1AP
- EMM
- ESM
- Data Traffic

### TC002_PERFORMANCE_data_600UE_500K.py

**Summary -**
Performance test case that setups up the S1 connection and `attaches` a high number of users to the AGW with moderate traffic. This test case runs for ~8hours and is supposed to provide signal on platform stability.

**Setup -**

- Total Subs - 600
- Total eNBs - 12
- Attach Rate - 3 UEs/sec
- Data Message Flow - HTTP_DL_500K_Per_UE
- Iterations - 1
- Runtime - 8 Hours

**Validation Categories -**

- S1AP
- EMM
- ESM
- Data Traffic

### TC003_PERFORMANCE_data_300UE_12enbs_3rate_active_idle_12h.py

Performance test case that setups up the S1 connection and `attaches` a moderate number of users to the AGW with moderate data traffic. Further, the test will force devices to transition through active and idle states in a loop.

**Setup -**

- Total Subs - 300
- Total eNBs - 12
- Attach Rate - 3 UEs/sec
- Data Message Flow - HTTP_DL_500K_Per_UE
- Iterations - 1
- Runtime - 12-13 Hours

**Validation Categories -**

- S1AP
- EMM
- ESM
- UE_STATE_CHECK
- active_idle

## Availability Tests

### TC001_AVAILABILITY_MIX.py

Availability test case that setups up the S1 connection and `attaches` a moderate number of users to the AGW with moderate data traffic. Further during the availability testing period, probing traffic will be sent to ensure that the data plane and control plane are available.

**Setup -**

- Total Subs - 240
- Total eNBs - 6
- Attach Rate - 3 UEs/sec
- Data Message Flow - HTTP_DL_500K_Per_UE
- Iterations - 1
- Runtime - 3 Hours

**Validation Categories -**

- Availability

## Validation Categories

- S1AP
    - S1 Setup Requests == # of eNBs
    - S1 Setup Responses == # of eNBs
    - S1 Release Timeouts <= 0.1% of S1 Release Requests
    - Init Context Setup Fail <= 0.1% of Init Context Setup Requests
- ESM
    - PDN Connectivity Requests >= Total Subs
    - PDN Connectivity Success >= 99th pctl of Total Subs
- EMM
    - Attach Requests >= Total subs
    - Attach Accepts >= 99th pctl of Total subs
    - Service Rejects <= 0.1% of Service Requests
    - Service Request Timeout <= 0.1% of Service Requests
- Data Traffic
    - Sanity
        - Total_Subs *DMF_Per_UE_Tput* 0.95 <= third_quartile(l3_server_bitrate)
        - Total_Subs *DMF_Per_UE_Tput* 0.95 <= third_quartile(l3_client_bitrate)
    - All others
        - Total_Subs *DMF_Per_UE_Tput* 0.95 <= median(l3_server_bitrate)
        - Total_Subs *DMF_Per_UE_Tput* 0.95 <= median(l3_client_bitrate)
- UE_STATE_CHECK
    - MME State = 0
    - SPGW State = 0
    - S1AP State = 0
    - table0 Flows = 3
    - table12 Flows = 2
    - table13 FLows = 2
    - mobility_cli.py get_subscriber_table = 0
- HANDOVER
    - S1HO Setup Requests Sent / 2 == # of eNBs
    - S1HO Setup Responses Recvd / 2 == # of eNBs
    - Handoff Successes == Handoff Attempts
- HE
    - (Num of GET requests with header enriched / Total Num of GET requests) * 100 >= 95%
- active_idle
    - (L3 Client # of Pings sent / L3 Client # of ping responses revd) * 100 >= 95%
    - Pings are sent only once per active period. So in essence very similar to availability measurement
- Availability
    - Total_Unavailable_Time_Control = Total_Test_Duration - Total_Sessions_Connects (Metrics provided by Spirent)
    - Total_Unavailable_Time_Data = Total_Sessions_Connects - Total_Data_Traffic_Verified (Metrics provided is provided by Spirent)
    - Total_Unavailable_Time = Total_Unavailable_Time_Control + Total_Unavailable_Time_Data
    - Avail = (Total_Test_Duration - Total_Unavailable_Time) / Total_Test_Duration
- IPFIX
    - Total # of cflow.srcaddr (unique) == Total # of subs
- ENB
    - Total # of "Echo Requests Sent" > 0
    - Total # of "Echo Requests Sent" == Total # of "Echo Responses Received"
