---
id: pipelined_packet_tracing
title: Pipelined packet-tracer
sidebar_label: Pipelined packet-tracer
hide_title: true
---
# Packet Tracer
Packet tracer generates traffic in pipelined and sends the packets through
magma OVS tables.

The packet tracer is implemented as a standalone Ryu app which sets the
test-packet register (reg5) to 1 so that other apps in pipelined "know" that the
packet is not a real one and act accordingly.

Packet tracer installs additional flows to the magma tables (if test-packet register is set):
* send-to-controller flow corresponding to every drop flow in pipelined
* send-to-controller whenever there is no match for the 
table (otherwise it will be dropped as the default behaviour for pipelined 
is to drop the packet if no matching flows were found)


# Usage
```bash
magtivate
packet_tracer_cli.py -h
packet_tracer_cli.py icmp --src_ip 0.0.0.0 --dst_ip 8.8.8.8
packet_tracer_cli.py arp
packet_tracer_cli.py http
```
