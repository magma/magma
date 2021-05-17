---
id: p004_fua_restrict_feature
title: Service Restriction Feature
hide_title: true
---

# Final Unit Action - Service Restriction Feature

*Status: Accepted*\
*Author: @ymasmoudi*\
*Last Updated: 08/27*

## Overview 

Currently, upon exhaustion of subscriber quota, we support two final unit
actions : 
1/ session termination
2/ redirection to a predifined IP/URL address. 

A third option called service restriction would allow us to restrict traffic
to a list of IP filter rules with specified QoS. 

FUA-Restrict indicates to SessionD that the user's access MUST be restricted
according to IP filters definitions or IDs provided in AVPs. This is enabled
by sending CCA-I/U with a Final-Unit-Action AVP set to RESTRICT\_ACCESS.


## Propostion

FUA-Restrict can be enabled by including either the Restriction-Filter-Rule
AVP or the Filter-Id AVP in the Credit-Control-Answer message.

```
  [Final-Unit-Indication]
    {Final-Unit-Action}
    *[Restriction-Filter-Rule]
    *[Filter-Id]
    [Redirect-Server]
      {Redirect-Address-Type}
      {Redirect-Server-Address}
```

Restriction-Filter-Rule: of type IPFilterRule (FlowDescription). This will
need to be translated to a PolicyRule.
Filter-Id: Provides IP packet filter identifiers (In our case, It will be a
list of PolicyRule identifiers)

It is important to note that QoS can be enabled only when using Filter-Id AVP.
Restriction-Filter-Rule AVP would only provide allowed traffic filters.

This implementation will not provide Restriction-Filter-Rule as Filter-id
provides support of all the required use-cases. Nevertheless, 
Restriction-Filter-Rule could be translated to restriction dynamic rules with
lower priority and apply them similar to static rules. 
LocalEnforcer will ensure priority in this case.

As we would like to use any static rule when restriction is enabled, there will
no change to PolicyRules definition.

## How We Will Change Magma

**PolicyRule**

Policy Object definition remains the same as the current definition. At the
opposite of a policy rule used in the FUA-Redirect case, there is no need to
introduce any additional field for FUA-Restrict.
We will be able to apply any static rule when restricting access.


**Session Proxy**

We will limit this implementation to Filter-Id AVP only as It covers the desired
use-cases. The AVP will reference a PolicyRule Identifier.

```
  [Final-Unit-Indication]
    {Final-Unit-Action}
    *[Filter-Id]
    [Redirect-Server]
      {Redirect-Address-Type}
      {Redirect-Server-Address}
```

The support of Restriction-Filter-Rule can be added later if required. In this
case, ReceivedCredits will need to be updated to include a FlowDescription list.

We will need to translation IPFilterRule to FlowDescription as well.


**SessionD**

ServiceAction will need to be updated to support RESTRICT\_ACCESS action. Mostly,
we will need to be able to store a restricted list of rule ids received from
OCS/SessionProxy.

In addition, we will need to store the list of rules ID to be used when restriction
is enabled. A method get\_restricted\_rule\_ids will be added to retrive this list
prior to sending rules activation to pipelined.

We will re-use the same pipelined client method add\_gy\_final\_action\_flow as It
provides support for this case as well.

Static Rules will be sent using their original priority and will be applied as gy
final action flow.

The service state SERVICE\_RESTRICTED will need to be supported as well.


**Pipelined**

GYController needs to be updated to support restriction FUA. This translates to,
1/ activate all the flows as they are received from SessionD.
2/ Support QoS in GY APP (Include a Qos Manager)
3/ Updating statistics back to SessionD
4/ Support Gy rules enforcement stats in SessionD

Since some of these logic is common between enforcement and gy apps, we will need
to move it to  policy\_mixin.py


**Swagger API**

No change is planned


**NMS**

No change is planned


**Additional Changes**

Full support Restrict FUA-Restrict in OCS
Add integration tests for FUA-Restrict to cover Filter-Id with and without QoS.
