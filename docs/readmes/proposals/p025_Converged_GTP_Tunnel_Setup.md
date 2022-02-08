---
id: p025_Converged_GTP_Tunnel_Setup
title: Converged GTP Tunnel Setup for 4G
hide_title: true
---
# Converged GTP Tunnel Setup for 4G

*Status: Draft*\
*Authors: Prabinak*\
*Reviewers: Magma team*\
*Last Updated: 2022-02-08*

## **Objective**

Add support for to make GTP tunnel operations by Sessiond for 4G Call as similar to 5G Call Flows.

## **Motivation**

Today Magma AGW is making the GTP tunnel operations using below two ways for 4G call:

1. Using openflow controller or libgtpnl from SPGW.
2. Or Making gRPC call from MME to Pipelined.

For option-2 is doing extra gRPC call to configure the GTP tunnel in OVS but
not using the existing gRPC call from MME to Sessiond to Pipelined (Activate_flow and deactivate_flow).

Intent of the proposal is to make common gRPC call to configure GTP flows,
Enforcement flows and enforcement stats flows as similar to 5G flows.
For this proposal, we reduce the operation cost of Magma AGW.

## Design

To extend the existing gRPC methods of MME-Sessiond and Sessiond-Pipelined to configure GTP tunnel flows.
Following scenarios or process are need to implement using gRPC:

### GTP tunnel Creation

1. MME-Sessiond:
    rpc CreateSession(LocalCreateSessionRequest) returns (LocalCreateSessionResponse) {}
    rpc UpdateTunnelIds(UpdateTunnelIdsRequest) returns (UpdateTunnelIdsResponse) {}
2. Sessiond-Pipelined:
    Activate flows for a subscriber based on predefined flow templates
    rpc ActivateFlows (ActivateFlowsRequest) returns (ActivateFlowsResult) {}

### GTP tunnel Deletion

1. MME-Sessiond:
    rpc EndSession(LocalEndSessionRequest) returns (LocalEndSessionResponse) {}
2. Sessiond-Pipelined:
    Deactivate flows for a subscriber
    rpc DeactivateFlows (DeactivateFlowsRequest) returns (DeactivateFlowsResult) {}

![Create and Delete GTP tunnel block diagram](assets/Create_Delete_tbl0.png)

### GTP tunnel flows forward and suspend notifications

1. MME-Sessiond:
    rpc Suspend_notification_request(LocalNotificationRequest) returns (LocalNotificationResponse) {}
    rpc Resume_Data_notification_request(LocalNotificationRequest) returns (LocalNotificationResponse) {}
2. Sessiond-Pipelined:
    Suspend flows for a subscriber
    rpc SuspendGTPFlows(SuspendFlowsRequest) returns (SuspendFlowsResult) {}
    rpc ResumeGTPFlows(ResumeFlowsRequest) returns (ResumeFlowsResult) {}

![Notifications of Suspend and Resume GTP tunnel block diagram](assets/Notifications_Suspend_Resume_tbl0.png)

### GTP tunnel for IdleSession

1. MME-Sessiond:
    rpc IdleSession(LocalIdleSessionRequest) returns (LocalIdleSessionResponse) {}

2. Sessiond-Pipelined:
    rpc PagingFlows(PagingFlowsRequest) returns (PagingFlowsResult) {}

### GTP tunnel for PagingNotification

1. Sessiond-MME:
    rpc PagingNotification(LocalPagingInfo) returns (void) {}

2. Pipelined-Sessiond:
    rpc SendPagingRequest(UPFPagingInfo) returns (SmContextVoid) {}

![Idle Mode and Paging GTP tunnel block diagram](assets/Idle_mode_Paging_tbl0.png)

### **Configuration**

AGW need following configuration:

1. pipelined_managed_tbl0: false in pipelined.yml and spgw.yml

## **Implementation plan**

1. Create new gRPC methods for MME-Sessiond in session_manager.proto.
2. Create new gRPC methods for Sessiond-Pipelined in pipelined.proto.
3. Add gRPC handler (service endpoints) for new methods in Sessiond,MME and Pipelined.
4. Add support Paging notifications for 4G call in sessiond.

## **Testing Plan**

The following tests scenario are proposed to verify to cover this functionality:

### Integration Testing

1. Verify the functionality using S1ap test cases.
2. Write UT test cases for sessiond and pipelined.
3. Write the python stub for verify the functionality.  

### Regression Testing

Following tests will be verified on magma master branch to make sure there is no breakage introduced through this proposal.
• Verify the 5G and 4G basic call flows and paging using UERANSIM, Teravm etc.
• Need to verify the traffic test.
