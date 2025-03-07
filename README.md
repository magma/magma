<h1 align="center">
    <a href="https://www.magmacore.org/"><img src="https://raw.githubusercontent.com/magma/magma/master/docs/docusaurus/static/img/magma-logo-purple.svg" alt="Magma" width="550"></a>
</h1>

<h3 align="center">Connecting the Next Billion People</h3>

<p align="center">
    <a href="https://github.com/magma/magma/blob/master/LICENSE"><img src="https://img.shields.io/badge/license-BSD3clause-blue.svg" alt="License"></a>
    <a href="https://github.com/magma/magma/releases"><img src="https://img.shields.io/github/release/magma/magma" alt="GitHub Release"></a>
    <a href="https://github.com/magma/magma/wiki/Contributing-Code"><img src="https://img.shields.io/badge/PRs-welcome-brightgreen.svg" alt="PR's Welcome"></a>
    <a href="https://github.com/magma/magma/graphs/contributors"><img src="https://img.shields.io/github/contributors/magma/magma" alt="GitHub contributors"></a>
    <a href="https://github.com/magma/magma/commits/master"><img src="https://img.shields.io/github/last-commit/magma/magma" alt="GitHub last commit"></a>
    <a href="https://github.com/magma/magma/commits/master"><img src="https://img.shields.io/github/commit-activity/y/magma/magma" alt="GitHub commit activity the past week"></a>
    <a href="https://codecov.io/gh/magma/magma"><img src="https://codecov.io/gh/magma/magma/branch/master/graph/badge.svg" alt="CodeCov"></a>
</p>

Magma is an open-source software platform that gives network operators an open, flexible and extendable mobile core network solution. Magma enables better connectivity by:

- Allowing operators to offer cellular service without vendor lock-in with a modern, open source core network
- Enabling operators to manage their networks more efficiently with more automation, less downtime, better predictability, and more agility to add new services and applications
- Enabling federation between existing MNOs and new infrastructure providers for expanding rural infrastructure
- Allowing operators who are constrained with licensed spectrum to add capacity and reach by using Wi-Fi and CBRS

## Documentation

- [Magma Website](https://magmacore.org/): Project landing page
- [Docs](https://magma.github.io/magma/docs/basics/introduction.html): Deployment, configuration and usage information
- [Code](https://github.com/magma): Source code
- [Contributing](https://github.com/magma/magma/wiki/Contributor-Guide): Contributor Guide
- [Wiki](https://wiki.magmacore.org/): Meeting notes and project team resources
- [Rewards Program](REWARDS_PROGRAM.md): How to participate in the Magma rewards program

## Join the Magma community

Direct detailed technical questions to [the GitHub Discussions page](https://github.com/magma/magma/discussions). Your question might already have an answer, and if not we want the answer you get to be available in the future.

Contact us in real-time using Slack: [magmacore workspace](https://slack.magmacore.org/).

## Magma Architecture

The figure below shows the high-level Magma architecture. Magma is 3GPP generation (2G, 3G, 4G or upcoming 5G networks) and access network agnostic (cellular or WiFi). It can flexibly support a radio access network with minimal development and deployment effort.

Magma has three major components

- **Access Gateway.** The Access Gateway (AGW) provides network services and policy enforcement. In an LTE network, the AGW implements an evolved packet core (EPC), and a combination of an AAA and a PGW. It works with existing, unmodified commercial radio hardware.

- **Orchestrator.** Orchestrator is a cloud service that provides a simple and consistent way to configure and monitor the wireless network securely. The Orchestrator can be hosted on a public/private cloud. The metrics acquired through the platform allows you to see the analytics and traffic flows of the wireless users through the Magma web UI.

- **Federation Gateway.** The Federation Gateway integrates the MNO core network with Magma by using standard 3GPP interfaces to existing MNO components.  It acts as a proxy between the Magma AGW and the operator's network and facilitates core functions, such as authentication, data plans, policy enforcement, and charging to stay uniform between an existing MNO network and the expanded network with Magma.

![Magma architecture diagram](docs/readmes/assets/magma_overview.png?raw=true "Magma Architecture")

## License

Magma is BSD License licensed, as found in the LICENSE file.

The EPC originates from OAI (OpenAirInterface Software Alliance) and is offered under the same BSD-3-Clause License.

## Risks

The Magma materials are provided in accordance with the licenses made available in the LICENSE file. Prior to using the materials, it is highly recommended that you test and verify that the materials meet your specific requirements, including, without limitation, any and all security and performance requirements.

## Security

Responsible disclosures from independent researchers are gratefully accepted. See the [Security Policy](/magma/magma/security/policy) for submission details and [Security Overview for Contributors](https://github.com/magma/magma/wiki/Security-Overview-for-Contributors) to learn about other ways of contributing.

We appreciate very much the valuable disclosures (and often fixes) from the following security researchers:

- Nathaniel Bennett (University of Florida). The contributions by Mr. Bennett required deep knowledge of code, networking, security, and Magma. We're grateful.
- Yilu Dong, Tianchang Yang, Abdullah Al Ishtiaq, Syed Md Mukit Rashid, Ali Ranjbar, Kai Tu, Tianwei Wu, Md Sultan Mahmud, Syed Rafiul Hussain (The Pennsylvania State University)
- Guarang Maheta
- Phi Trần
