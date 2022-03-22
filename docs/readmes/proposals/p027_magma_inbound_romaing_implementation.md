# Magma extensions for inbound roaming design response
Author(s): [@just-now]
Last updated: 03/20/2022

## 1 Background and Objective
The objective of this proposal is to react on p026 proposal w.r.t. to
code changes in magma. This document presents a high level design
describing enablement of additional features related to inbound
roaming, such as traffic monitoring, quota management, etc. with the
help of existing S5 based session management code.
The main purposes of this document are: to be inspected by Magma
developers, architects and peer designers to ascertain that high level
design is aligned with Magma architecture and other designs, and
contains no defects; to serve as a design reference document.

## 2 Implementation Scope

### 2.1 Assumptions
 - Currently, for S8 sessions user-plane tunnels are created and
   controlled on MME side without engagement of SessionD and
   PipelineD.
 - It's possible to create fake S5 session for each S8 session to cope
   with traffic quota management.
 - It's possbile to disable PipelineD-related activities in fake S5
   session except traffic monitoring like `AsyncPipelinedClient::poll_stats()`.
 - It's possbile to associate dedicated `SessionReporter` with fake S5
   session to send CCR-requests to remote OCS through Session Proxy.
 - It's possible to reuse S5 session termination flow on quota exhaust
   or any other known event to trigger termination of associated S8
   session.
 - Design doesn't take into account failures and availability of
   related remove and local services like remote OCS availability
   accessed through Session Proxy; PipelineD or Redis crash restart.

### 2.2 Associated S5 and S8 sessions flows
 - Introduce fake S5 session type used along with S8 sessions.
 - Create S8 session flow:
   - Create S5 fake session if it failed do S5 session cleanup;
   - Create associated S8 session if it failed do S5+S8 session cleanup.
 - Termination S8 session flow:
   - Terminate S5 session;
   - Terminate S8 session.
 - Terminate S5+S8 sessions on quota exhaust:
   - `SpgwServiceImpl::DeleteBearer` is being sent to SPGW when S5
     session is being terminated due to quota exhaust event;
   - Reuse `SpgwServiceImpl::DeleteBearer` or nested handlers to
     terminate S8 session in case if S5 session is a fake session.
 - Credit control flow:
   - For each Home PLMN create `SessionReporter` and corresponding thread.
   - Associate a specific reported with a fake S5 session (put a
     pointer into session, for example).
   - For fake sessions, use associated reporters.
 - Related PipelineD flows:
   - As S8 user plane tunnels are managed by MME, pipelined tunnel
   management for fake sessions needs to be switched off:
```
[14:44:13] just-now-i9@session_manager git:(master)$ gg 'pipelined_client_->'
LocalEnforcer.cpp:160:    pipelined_client_->setup_cwf(session_infos, quota_updates, ue_mac_addrs,
LocalEnforcer.cpp:164:    pipelined_client_->setup_lte(session_infos, epoch, callback);
LocalEnforcer.cpp:287:  pipelined_client_->poll_stats(
LocalEnforcer.cpp:457:    pipelined_client_->deactivate_flows_for_rules_for_termination(
LocalEnforcer.cpp:534:  pipelined_client_->activate_flows_for_rules(
LocalEnforcer.cpp:564:    pipelined_client_->delete_ue_mac_flow(
LocalEnforcer.cpp:615:  pipelined_client_->deactivate_flows_for_rules_for_termination(
LocalEnforcer.cpp:682:  pipelined_client_->add_gy_final_action_flow(imsi, ip_addr, ipv6_addr, teids,
LocalEnforcer.cpp:823:        pipelined_client_->activate_flows_for_rules(
LocalEnforcer.cpp:882:        pipelined_client_->activate_flows_for_rules(
LocalEnforcer.cpp:935:        pipelined_client_->deactivate_flows_for_rules(
LocalEnforcer.cpp:976:          pipelined_client_->deactivate_flows_for_rules(
LocalEnforcer.cpp:1261:  pipelined_client_->update_subscriber_quota_state(
LocalEnforcer.cpp:1489:        pipelined_client_->deactivate_flows_for_rules(
LocalEnforcer.cpp:1827:    pipelined_client_->deactivate_flows_for_rules(imsi, ip_addr, ipv6_addr,
LocalEnforcer.cpp:1834:    pipelined_client_->activate_flows_for_rules(
LocalEnforcer.cpp:2000:        pipelined_client_->add_ue_mac_flow(
LocalEnforcer.cpp:2073:  pipelined_client_->update_ipfix_flow(sid, ue_mac_addr,
SessionStateEnforcer.cpp:498:  pipelined_client_->set_upf_session(sess_info, pending_activation,
```
