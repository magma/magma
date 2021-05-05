---
id: ref_pcap
title: PCAP Collection
hide_title: true
---

# PCAP collection

The following documents are a collection of `.pcap` files that can be opened
using [Wireshark](https://www.wireshark.org/).

The files below capture external Magma interfaces such as the S1AP, S6a, Gx, Gy,
and S8. Some PCAPs merge all of these interfaces into a single file in order to
make the entire call flow visible. If you wish to examine these traces
chronologically, you can sort the packets by time.

A few of the PCAP samples include user plane information. In these cases, you
may see the packets duplicated due to the capture of both the S1AP and SGi.

![Pcap Sample](assets/feg/pcap_sample.png?raw=true "PCAP Sample")

## Samples

[[pcap]](assets/feg/pcaps/gx_gy_combined_03_210505_132953-1.pcapng)
**Basic Session Cycle with Static and Dynamic Rules (Gx/Gy)**
One subscriber creates a session, receives multiple rules, and sends traffic
through the network.

[[pcap]](assets/feg/pcaps/gx_gy_combined_03_210505_132953-1.pcapng)
**Basic Session Cycle with Static and Dynamic Rules (Gx/Gy)**
32 subscribers create a session, receive multiple rules, and send traffic
through the network.

[[pcap]](assets/feg/pcaps/gx_gy_combined_05_210505_133111-1.pcapng)
**Gx Reporting and Quota Exhaustion**
One subscriber creates a session, receives a quota from OCS, uses up all the
quota available, and gets disconnected.

[[pcap]](assets/feg/pcaps/gx_gy_combined_06_210505_133132-1.pcapng)
**Gy Reporting and Quota Exhaustion**
One subscriber creates session, receives a monitor quota from the PCRF, uses up
all the quota available, and stops reporting. In this case, the session is not
terminated since this is Gx monitoring.

[[pcap]](assets/feg/pcaps/inbound_roaming_01_210505_130828-1.pcap)
**Inbound Roaming Basic Session Cycle**
One roaming subscriber begins a S8 session by sending 10 packets and receiving
10 packets.

[[pcap]](assets/feg/pcaps/inbound_roaming_05_210504_215514-1.pacp)
**Inbound Roaming IDLE to Connected UE initiated**
One roaming subscriber begins a S8 session, goes from "Idle" to "Connected"
followed by sending 10 packets and receiving 10 packets.
