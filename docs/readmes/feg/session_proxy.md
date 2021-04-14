# Session Proxy

## Overview of Session Proxy
session_proxy is a service from the FEG that interacts with PCRF and OCS.
It's main function is to translate GRPC messages from sessiond into Diameter protocol back and forth.

## Interfaces
1. Session Manager (sessiond)<br>
Magma PCEF is implemented by Session Manager (or sessiond). Session Proxy receive the GRPC messages
translates from sessiond and translates them into diameter messages for creation (CCR-I and CCA-I),
update (CCR-U and CCA-U) and termination (CCR-T and CCA-T).

2. PCRF (Gx)<br>
session_proxy implements Gx interface and will translate policy related calls into diameter AVPs to
be sent to PCRF. It supports some events like monitoring tracking, time based rules,
and PCRF initiated messages like reauthentication

3. OCS (Gy)<br>
session_proxy implements Gy interface and will translate charging related calls into diameter AVPs to
be sent to OCS. It supports charging reporting and OCS initiated messages like reauthentication

4. PoicyDb<br>
session_proxy will get static rules and omnipresent rules from policyDb to inject
them to the responses back to SessionD


##Session Proxy Configuration
Configuration of Session Proxy can be done through NMS or Swagger API. In both
cases the configuration of Session Proxy is through Gx and Gy labels/tabs.


## Magma PCRF-less configuration
Omnipresent rules can be used as a way to achieve a PCRF-less (Gx) configuration while
maintaining charging (Gy) capabilities. Omnipresent rules are configured on Orc8r and injected
by session proxy to any subscriber that creates a new session. So if those omnipresent rules define a
Rating Group (or charging key), that key will be used by to charge the user and will
be communicated through Gy. Note that **all** subscribers will need to have that Rating Group
configured on OCS so the reporting is tracked.

###Configure  Magma PCRF-less:
- Create omnipresent rules:<br>
On NMS, create a static rule. In case you want to use charging (Gy) for that rule add a Rating Group
Then check the box that says `omnipresent`. Once this is enable, any subscriber that attaches to the
network will get that rule installed. Remember to check that subscriber on OCS too.

- Disable Gx (optional): <br>
If there is not a PCRF to connect to, you can disable Gx so your session proxy doesn't try to connect to a
non-existing PCRF. To do so go to swagger API `/feg/{network_id}/gateways/{gateway_id}`
search for your Federated Gateway(using `feg network` and `feg gateway`) and modify `disableGx` under `Gx` key.
