---
id: version-1.7.0-l3_transport
title: L3 transport for AGW
hide_title: true
original_id: l3_transport
---
# Configuration of SGi L3 transport

Magma AGW has configuration option to create and provision wireguard tunnel
for PGW traffic, This feature is targeted for inbound roaming deployment.

## Introduction

Today AGW operator need to have separate networking device to create the tunneling
infrastructure for pushing PGW traffic to aggregator host. This scheme also needs
separate public IP which is expensive to have these days. By providing support for
tunneling traffic from AGW to aggregator host we eliminate both requirement and
reduce operation cost of Magma AGW.

## Configure PipelineD

PipelineD manages all datapath tunnels, so naturally pipelineD needs configuration
for wire-guard tunnel.
the tunnel can be configured using pipelined.yml:

```text
sgi_tunnel:
 enabled: true
 type: wg
 enable_default_route: false
 tunnels:
 - wg_local_ip: 172.168.100.1/24
   peer_pub_key: 'VQT+tLY6/xF+k1WqrXeQzlfb8hWMVLcPdtCPvwIUNU0='
   peer_pub_ip: 1.2.3.4

```

- `sgi_tunnel`: option defines a L3 transport tunnels
- `enabled`: option to enable this functionality.
- `type`: Today pipelineD support wireguard tunnel type only so only supported type
is 'wg', There is plan to add support for IPSec tunnel.
- `enable_default_route`: Set this to true for routing all traffic via L3 transport
including SGi maganagement traffic
- `tunnels`: This section defines tunnel endpoint.
- `wg_local_ip`: IP address of the wireguard define on AGW
- `peer_pub_key`: public key of peer wireguard device.
- `peer_pub_ip`: public IP of remote wireguad tunnel host.

Once this configuration is committed to the pipelined.yml file, you need to
restart magma service. PipelineD generates private/public key pair if needed
and creates `magma_wg0` device for wireguard tunnel.

Validate the wireguard tunnel creation by checking `wg show` cmd.

```text

sudo wg show
interface: magma_wg0
  public key: UMlfoG9ZwCNhpm1ezfL1MfzT2qTBTzkLi30cSrgXTDc=
  private key: (hidden)
  listening port: 9333

peer: VQT+tLY6/xF+k1WqrXeQzlfb8hWMVLcPdtCPvwIUNU0=
  endpoint: 1.2.3.4:9333
  allowed ips: 172.168.100.0/24
  transfer: 0 B received, 740 B sent
  persistent keepalive: every 15 seconds
```

You need 'public key' from the AGW wiregaurd tunnel to configure the remote
endpoint.

## Configuration of SPGW

Enable following option to route PGW traffic over L3 tunnel.

```text
agw_l3_tunnel: true
```

The actual routing rules for PGW are created on demand on inbound roaming UE
is attached to the AGW.

### Validation of WireGuard tunnel

Validate that tunnel is created and traffic is routed for PGW tunnel using
`gw show`

For example in following scenario PGW IP is routed over wireguard device, you
can check `allowed ips` section of the `peer`.

```text
sudo wg show
interface: magma_wg0
  public key: UMlfoG9ZwCNhpm1ezfL1MfzT2qTBTzkLi30cSrgXTDc=
  private key: (hidden)
  listening port: 9333

peer: VQT+tLY6/xF+k1WqrXeQzlfb8hWMVLcPdtCPvwIUNU0=
  endpoint: 1.2.3.4:9333
  allowed ips: 192.168.60.141/32, 172.168.100.0/24
  transfer: 0 B received, 10.70 KiB sent
  persistent keepalive: every 15 seconds
```
