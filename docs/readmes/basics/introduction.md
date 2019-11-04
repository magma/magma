---
id: introduction
title: Introduction
hide_title: true
---
# Introduction

Magma is an open-source software platform that gives network operators an open, flexible and extendable mobile core network solution. Magma enables better connectivity by:

* Allowing operators to offer cellular service without vendor lock-in with a modern, open source core network
* Enabling operators to manage their networks more efficiently with more automation, less downtime, better predictability, and more agility to add new services and applications
* Enabling federation between existing MNOs and new infrastructure providers for expanding rural infrastructure
* Allowing operators who are constrained with licensed spectrum to add capacity and reach by using Wi-Fi and CBRS

## Magma Architecture

The figure below shows the high-level Magma architecture. Magma is designed to be 3GPP generation and access network (cellular or WiFi) agnostic. It can flexibly support a radio access network with minimal development and deployment effort.

Magma has three major components:

* **Access Gateway:** The Access Gateway (AGW) provides network services and policy enforcement. In an LTE network, the AGW implements an evolved packet core (EPC), and a combination of an AAA and a PGW. It works with existing, unmodified commercial radio hardware.

* **Orchestrator:** Orchestrator is a cloud service that provides a simple and consistent way to configure and monitor the wireless network securely. The Orchestrator can be hosted on a public/private cloud. The metrics acquired through the platform allows you to see the analytics and traffic flows of the wireless users through the Magma web UI.

* **Federation Gateway:** The Federation Gateway integrates the MNO core network with Magma by using standard 3GPP interfaces to existing MNO components.  It acts as a proxy between the Magma AGW and the operator's network and facilitates core functions, such as authentication, data plans, policy enforcement, and charging to stay uniform between an existing MNO network and the expanded network with Magma.

![Magma architecture diagram](assets/magma_overview.png?raw=true "Magma Architecture")
