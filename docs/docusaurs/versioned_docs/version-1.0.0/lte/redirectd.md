---
id: version-1.0.0-redirectd
title: Redirection
hide_title: true
original_id: redirectd
---
# Redirection
### Overview
When we enable redirection for the subscriber the next request(next packet) is
intercepted by the redirection server. The subscriber then completes the tcp
handshake with the redirection server. After the tcp handshake is complete the
subscriber sends an HTTP GET request and receives the HTTP 302 response from
the redirection server. Then it establishes a new tcp connection with the
redirection address provided in the 302 response and traffic to this address
is allowed to go through pipelined.

HTTPS is not supported as without ssl certificates this isn't possible.

### Server details
Redirectd runs a flask server which sends back HTTP 302(redirect) responses.
The server uses redis to lookup redirect information for incoming requests.
We save the information for subscribers in pipelined, using src_ip of request
to redirect_info lookup. When no such information is found return a 404,
this shouldn't happen, means redirect info wasn’t properly saved.

Redirectd is also a dynamic service, it is only launched when mconfig
dynamic_services array has a 'redirectd' entry.


### Redirection in pipelined
EnforcementController instantiates the required flows for forwarding subscriber
traffic to the redirection server. Pipelined also saves the redirect
information in redis using subcriber_ip as the key(mobilityd is used for
getting the subcriber ip from imsi).
All this is only done when subscriber PolicyRule has redirection enabled.
Example PolicyRule with enabled redirection:
```
policy = PolicyRule(
    id='redirect_test',
    priority=3,
    flow_list=flow_list,
    redirect=RedirectInformation(
        support=1,
        address_type=2,
        server_address="http://about.sha.ddih.org/"
    )
)
```

*Description of added flows:*
* Add flow to allow UDP traffic so DNS queries can go through
* Add flows with a higher priority that allow traffic to and from the
  redirection address provided in the redirect rule
  - if address_type is url submit a dns query and allow access to resolved IPs,
    the resolved IPs are stored in a ttl cache to decrease num of dns requests
  - if address_type is IPv4 allow access to that redirection IP address
  - ignore IPv6 address_type redirection as we don't support it
  - ignore SIP_URI address_type is not implemented
* Intercept tcp traffic from UE, send it to the redirection server, also
intercept tcp traffic from Redirection server and send it back to the UE.
This is done by adding an OVS flow with a learn action (when packets from UE
hit this flow a new flow will be instantiated to intercept traffic from
Redirection server)
  - Flow matching TCP traffic from the user with port 80, modify&send packets
  to the Redirection server. The learn action will 'save' the original dst_ip
  address by loading it into the instantiated flow
  - Flow instantiated from the learn action, matching TCP traffic from the
  server with the UE ip_addr/tcp_sport, modify&send packets back to UE
* Drop other traffic (default for all subscribers)

### Packet path breakdown
```
Packet protocol diagram (only includes fields that are being checked):
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|       Source IP Address       |     Destination IP Address    |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|        TCP Source Port        |      TCP Destination Port     |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

UE ip: 192.168.128.9
Remote url ip: 185.199.110.8
Redirect server ip: 192.168.128.1

For the tcp handshake, initial HTTP request the traffic flow looks like this:

UE sends a packet to remote url
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|         192.168.128.9         |         185.199.110.8         |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|             43040             |               80              |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

Packet gets modified in EnforcementController, sent to the Redirection server
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|         192.168.128.9         |        *192.168.128.1*        |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|             43040             |               80              |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

Redirection server responds, packet is sent back into OVS
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|         192.168.128.1         |         192.168.128.9         |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|               80              |             43040             |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

Packet gets modified in EnforcementController, sent back to the UE
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|        *185.199.110.8*        |         192.168.128.9         |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|               80              |             43040             |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

After getting a 302 response from the redirect server the traffic can go
straight to the redirected address without being changed in pipelined
```
