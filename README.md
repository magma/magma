# Magma

[![magma](https://circleci.com/gh/magma/magma/tree/v1.4.svg?style=shield)](https://circleci.com/gh/magma/magma/tree/v1.4)

Magma is an open-source software platform that gives network operators an open, flexible and extendable mobile core network solution. Magma enables better connectivity by:

* Allowing operators to offer cellular service without vendor lock-in with a modern, open source core network
* Enabling operators to manage their networks more efficiently with more automation, less downtime, better predictability, and more agility to add new services and applications
* Enabling federation between existing MNOs and new infrastructure providers for expanding rural infrastructure
* Allowing operators who are constrained with licensed spectrum to add capacity and reach by using Wi-Fi and CBRS


## Magma Architecture

The figure below shows the high-level Magma architecture. Magma is 3GPP generation (2G, 3G, 4G or upcoming 5G networks) and access network agnostic (cellular or WiFi). It can flexibly support a radio access network with minimal development and deployment effort.

Magma has three major components:

* **Access Gateway:** The Access Gateway (AGW) provides network services and policy enforcement. In an LTE network, the AGW implements an evolved packet core (EPC), and a combination of an AAA and a PGW. It works with existing, unmodified commercial radio hardware.

* **Orchestrator:** Orchestrator is a cloud service that provides a simple and consistent way to configure and monitor the wireless network securely. The Orchestrator can be hosted on a public/private cloud. The metrics acquired through the platform allows you to see the analytics and traffic flows of the wireless users through the Magma web UI.

* **Federation Gateway:** The Federation Gateway integrates the MNO core network with Magma by using standard 3GPP interfaces to existing MNO components.  It acts as a proxy between the Magma AGW and the operator's network and facilitates core functions, such as authentication, data plans, policy enforcement, and charging to stay uniform between an existing MNO network and the expanded network with Magma.

![Magma architecture diagram](docs/readmes/assets/magma_overview.png?raw=true "Magma Architecture")

## Usage Docs
The documentation for developing and using Magma is available at: [https://docs.magmacore.org/docs/basics/introduction.html](https://docs.magmacore.org)

## Join the Magma Community

- Mailing lists:
  - Join [magma-dev](https://groups.google.com/forum/#!forum/magma-dev) for technical discussions
  - Join [magma-announce](https://groups.google.com/forum/#!forum/magma-announce) for announcements
- Slack:
  - [magma\_dev](https://join.slack.com/t/magmacore/shared_invite/zt-g76zkofr-g6~jYiS3KRzC9qhAISUC2A) channel

See the [CONTRIBUTING](CONTRIBUTING.md) file for how to help out.

## License

Magma is BSD License licensed, as found in the LICENSE file.

The EPC originates from OAI (OpenAirInterface Software Alliance) and is offered under the same BSD-3-Clause License.
