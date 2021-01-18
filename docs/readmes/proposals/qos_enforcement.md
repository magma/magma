# Magma QoS Policy and Federation API Cleanup

*Status: In-Review*\
*Author: @andreilee*\
*Last Updated: 08/13*

This document concerns:

1. The addition of configurable QoS for policies
2. The addition of more comprehensive associations between subscribers, QoS, 
policies, and APN
3. The enforcement of the above with an eventually consistent model
4. Details concerning implementing these changes in both 
federated/un-federated environments


## Use Cases and Capabilities to Be Enabled

We would like to enable QoS enforcement across the APN and bearer level.

**Private LTE**

As an example use case for QoS enforcement, a service provider may wish to 
segment users into tiers denoting their service quality, eg. bronze, silver, 
and gold. A service provider should then be able configure a QoS profiles for 
each of the three service levels.

Policy rules will each be assigned one of these QoS profiles. The same QoS
profile may be assigned to multiple rules.

Beyond segmenting subscribers into QoS levels, as our support for multi-APN
matures, we wish to enable QoS enforcement for multiple sessions 
of a single subscriber, one for each APN. Across these different sessions, each 
APN will have a different QoS for its default bearer. To enable this behavior, 
there should be the capability to assign a QoS profile to an APN for the 
default bearer.

It is also possible that a network operator would desire to have SGi
specified differently for each APN even when served from the same gateway. 

We will define the `ApnResource` as a configuration specifying how SGi 
interfaces are set up for an APN in a gateway. 

To enable the above, we have the following list of new API requirements for
Magma:

1. Configuration of policy QoS profiles
2. Configuration of APN entities
3. Configuration of ApnResource entities
4. Assignment of APN entities to subscribers (to denote which APNs a
subscriber has access to)
5. Assignment of ApnResource entities to gateways

**Federated LTE**

In federated use cases, we have slightly different requirements. APN entities 
may not need to be configured as they are provided by the federation network.
`ApnResource` entities will still need to be configured for the AGW.

We also wish to support cases where an HSS is not provided, and our 
`subscriberdb` service is used, but where the PCRF is federated, 
and so our `policydb` service is disabled.

Due to different federation requirements between operators who federate HSS,
Gx/Gy, or both, the following requirement is added:

6. Split our `relay_enabled` flag for enabling federation into two separate
flags, `hss_relay_enabled`, and `gx_gy_relay_enabled`

## How We Will Change Magma

### API Endpoints

To enable configuration of QoS profiles and APNs, we will support the 
following REST endpoints:

```
Policy QoS Profiles

GET     /lte/{network_id}/policy_qos_profiles
POST    /lte/{network_id}/policy_qos_profiles

GET     /lte/{network_id}/policy_qos_profiles/{qos_name}
PUT     /lte/{network_id}/policy_qos_profiles/{qos_name}
DELETE  /lte/{network_id}/policy_qos_profiles/{qos_name}

GET     /feg_lte/{network_id}/policy_qos_profiles
POST    /feg_lte/{network_id}/policy_qos_profiles

GET     /feg_lte/{network_id}/policy_qos_profiles/{qos_name}
PUT     /feg_lte/{network_id}/policy_qos_profiles/{qos_name}
DELETE  /feg_lte/{network_id}/policy_qos_profiles/{qos_name}



APNs 

GET     /lte/{network_id}/apns
POST    /lte/{network_id}/apns

GET     /lte/{network_id}/apns/{apn_name}
PUT     /lte/{network_id}/apns/{apn_name}
DELETE  /lte/{network_id}/apns/{apn_name}

GET     /feg_lte/{network_id}/apns
POST    /feg_lte/{network_id}/apns

GET     /feg_lte/{network_id}/apns/{apn_name}
PUT     /feg_lte/{network_id}/apns/{apn_name}
DELETE  /feg_lte/{network_id}/apns/{apn_name}



APN Resources

GET     /lte/{network_id}/gateways/{gateway_id}/resource_labels
POST    /lte/{network_id}/gateways/{gateway_id}/resource_labels

GET     /lte/{network_id}/gateways/{gateway_id}/resource_labels/{resource_label}
PUT     /lte/{network_id}/gateways/{gateway_id}/resource_labels/{resource_label}
DELETE  /lte/{network_id}/gateways/{gateway_id}/resource_labels/{resource_label}

GET     /feg_lte/{network_id}/gateways/{gateway_id}/resource_labels
POST    /feg_lte/{network_id}/gateways/{gateway_id}/resource_labels

GET     /feg_lte/{network_id}/gateways/{gateway_id}/resource_labels/{resource_label}
PUT     /feg_lte/{network_id}/gateways/{gateway_id}/resource_labels/{resource_label}
DELETE  /feg_lte/{network_id}/gateways/{gateway_id}/resource_labels/{resource_label}
```



### API Entity Definitions

Entity definitions for QoS, APN, and subscribers will change. 
The final state is summarized here:

```
policy_qos_profile:
  id: string
  class_id: enum[0,255]
  max_req_bw_dl: uint32
  max_req_bw_ul: uint32
  # Guaranteed bit rate
  # Not enforced by Magma yet
  gbr:
    gbr_ul: uint32
    gbr_dl: uint32
  # Allocation and retention priority
  # Not enforced by Magma yet
  arp:
    arp_priority_level: uint32
    arp_preemption_capability: bool
    arp_preemption_vulnerability: bool
```

This is an addition, as there is an existing `qos_profile` definition, which
will continue to be used only for the `apn_configuration` entity.

This `policy_qos_profile` will be marshalled to the `FlowQos` protobuf
message.

```
flow_qos:
  ...
```

The `policy_qos_profile` will replace the `flow_qos` swagger entity
in funtctionality.
`flow_qos` should be marked as deprecated and should be discouraged from further
use in the API.

```
apn_configuration:
  ...
  # Unchanged
  ...
```

As before, a single QoS profile will be associated to an APN object at creation.
This will define the QoS for the default bearer, in combination with the `ambr`
field, and is unchanged.

```
apn_resource:
  id: string
  # Identifies that the SGi interface is set up for the specified APN
  apn_name: string
  # Specify either VLAN or the interface name for SGi
  # VLAN to use when communicating over SGi interface
  vlan_id: string
```

The `apn_resource` defines how the APN is serviced by each access gateway.
Options are provided for IP allocation, and the uplink. It is possible that
two APNs are serviced by separate uplink interfaces on the same access gateway.

```
mutable_lte_gateway:
  id: string
  name: string
  description: string
  ...
  apn_resources: array[string]
  ...
```

We are adding the `apn_resources` field to `lte_gateway` and
`mutable_lte_gateway` and they are specified by ID.

```
network_epc_configs:
  ...
  hss_relay_enabled: bool
  gx_gy_relay_enabled: bool
  ...
```

The `relay_enabled` flag will be split to enable federation under different
circumstances, with finer control. The combination of `hss_relay_enabled` and
`gx_gy_relay_enabled` together have the same functionality as the current
`relay_enabled`.

```
subscriber:
  id: string
  name: string
  ...
  # The subscriber has access to the default bearer for all APNs listed here
  active_apns: array[string]
  active_base_names: array[base_name]
  # Policy IDs listed in active_policies define which policies will be installed
  # for every APN of the subscriber
  active_policies: array[string]
  # Keyed by APN IDs, and values list out the policy IDs
  # policy IDs define what policies will be installed for the subscriber for 
  # each APN beyond default bearers
  active_policies_by_apn: dict[string, array[string]]
  ...
```

To now specify APN-specific policies, the field `active_policies_by_apn` will
be added. If an APN key is not present in `active_policies_by_apn`, then the
subscriber does not have access to that APN.

To ensure that edges are maintained correctly in `configurator`, a backend
entity `apn_policy_profile` should be maintained (and not visible through the
API). This `apn_policy_profile` should have no fields, but only edges. The
subscriber will have edges to its own `apn_policy_profile`s, and each
`apn_policy_profile` will have edges to policies to identify which policies
are installed for each (IMSI, APN) combination. Each APN will also have edges
to its corresponding `apn_policy_profile`s. When policies are deleted, then
the subscriber will no longer be assigned those policies for its APNs. When APNs
are deleted, each subscriber should no longer have assignments to policies for
those APNs.

```
policy_rule:
  id: string
  flow_list: array[flow_description]
  priority: uint32
  rating_group: uint32
  monitoring_key: string
  ...
  # Set to be deprecated and replaced by qos_profile
  qos: flow_qos
  # Refers to the ID of the qos_profile
  qos_profile: string
```

One QoS object can be associated to each policy object.
The same QoS object can be associated to multiple different policies.

`policy_qos_profile` replaces the `flow_qos` in functionality, but because of
the change from 1:1 mapping used previously for `flow_qos`, a new field is
introduced instead of reusing the old field.

A migration is required to create a `qos_profile` for each `flow_qos`, and
assign it to the corresponding `policy_rule`.


### Pipelined Changes for Enforcement

At a high level, changes in `pipelined` are required to support the following:
* (P1) Enforce QoS via AMBR across the bearer level, and APN level (not
subscriber level, as UE_AMBR is enforced on RAN side)
* (P1) Enforce VLAN specification in `apn_resource` definition supplied to the
gateway for each APN
* (P2) Load/KPI reporting (to be used for admission control and session 
lifecycle management by sessiond)
* (P3) GBR enforcement


### Policydb Changes with Streaming

The `policydb` service across the orc8r and gateway requires changes to support 
QoS enforcement:
* (P1) Rename `relay_enabled` flag to `gx_gy_relay_enabled`
* (P1) Orc8r `policydb` service should be updated to stream to the gateway 
all the mappings between basenames, policies, QoS profiles, APNs, data plans 
and subscribers
* (P1) Gateway `policydb` service should interface with `sessiond` to trigger
rule activations/deactivations when there are updates to subscribers'
active policies. This will occur through a new gRPC service implemented by
`sessiond`


### Session Manager Changes

The session manager will be changed to support the following:
* (P1) Provide gRPC methods for `policydb` to set the policies installed for
each subscriber
* (P1) Accept the default QoS profile with APN that is received from MME
* (P1) Dedicated bearer Life Cycle Management (LCM): Associations of policies 
(with QoS fields) with the subscribers will define the lifecycle of 
dedicated bearers
* (P2) Admission control and pre-emption. This is enforcing the AMBR fields 
specified in the current `apn_configuration`


## Policydb and Session Manager gRPC Interface

As specified in the previous two sections, `sessiond` will provide a gRPC
service for `policydb` to set installed policies for attached subscribers.

An outline of the gRPC service is as follows:

```
message SessionPolicySet {
  bool apply_subscriber_wide = 1;
  string apn = 2;
  repeated StaticRuleInstall static_rules = 3;
  repeated DynamicRuleInstall dynamic_rules = 4;
}

// Send only for subscribers for which there were updates
message SubscriberPolicySet {
  string imsi = 1;
  repeated SessionPolicySet = 2;
}

message SubscriberPolicyUpdates {
  repeated SubscriberPolicySet policies_subscriber = 1;
}

service SessionUpdateResponder {
  // Update the specified subscribers with the currently active rules
  // 
  rpc UpdateActiveRules (SubscriberPolicyUpdates) returns (magma.orc8r.Void) {}
}
```

It is expected that `policydb` will receive streamer updates of active rules
from the orc8r. For all subscribers that have updated active rules, it will
call `UpdateActiveRules` to update `sessiond` on all subscribers for which
there were changes. For the subset of which there are subscribers connected
on the gateway, `sessiond` will enforce rule changes.

To ensure that rule activations/deactivations are not missed, a P2 goal is to
ensure that `policydb` will periodically update `sessiond` on subscribers for
which no streamed updates were received.


### Subscriberdb Changes

* (P1) Rename `relay_enabled` to `hss_relay_enabled`
* (P1) Stream down to the gateway the default QoS policy, and `apn_resource`
as part of the APN Configuration
* (P1) Send to `mobilityd` the updated `ApnConfiguration` including default
QoS policy and `apn_resource`


### Mobilityd Changes

* (P1) Accept the updated `ApnConfiguration` protobuf with included QoS policy,
and `apn_resource` 


### NMS Changes

Our front-end interface needs to change to allow these configurations to be
made easily. These are second priority. Some API changes cause
backwards-compatibility issues and will need to be updated soon.

* (P2) ApnResource configuration on creation/edit of LTE gateways
* (P2) Apn Profile creation/edit modals
* (P2) Policy QoS Profile creation/edit modals
* (P2) Policy creation/edit includes fields to specify QoS
* (P2) LTE network creation/edit includes fields to specify supported APNs
* (P2) Subscriber creation/edit includes fields to specify policy assignments
per APN
* (P2) Creation of APN entities disabled when HSS is federated


### FeG Changes

* (P1) Split support of `relay_enabled` into `hss_relay_enabled` and
`gx_gy_relay_enabled`

