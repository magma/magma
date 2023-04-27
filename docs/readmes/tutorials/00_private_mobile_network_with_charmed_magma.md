---
id: 00_overview
title: Overview
hide_title: true
---

# Operate your own private mobile network with Charmed Magma

In this tutorial, we will use Juju to deploy and run Magma's 4G core network on AWS.
We will also deploy a radio and cellphone simulator from the [srsRAN](https://www.srslte.com/)
project to simulate usage of this network. You will need:

- An AWS account[^1]
- A public domain
- A computer[^2] with the following software installed:
    - [juju 2.9](https://juju.is/docs/olm/install-juju)
    - [kubectl](https://kubernetes.io/docs/tasks/tools/)
    - [aws cli](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html)
    - [eksctl](https://docs.aws.amazon.com/eks/latest/userguide/eksctl.html)

## Table of contents

1. [Getting Started](01_getting_started.md)
2. [Deploying Magma Orchestrator](02_deploying_magma_orchestrator.md)
3. [Deploying Magma Access Gateway](03_deploying_magma_access_gateway.md)
4. [Integrating Magma Access Gateway with Magma Orchestrator](04_integrating_magma_access_gateway_with_magma_orchestrator.md)
5. [Deploying the radio simulator](05_deploying_the_radio_simulator.md)
6. [Simulating user traffic](06_simulating_user_traffic.md)
7. [Destroying the environment](07_destroying_the_environment.md)

[^1]: This tutorial uses AWS as the cloud provider. You can use any cloud provider
that Juju supports. See [Juju Clouds](https://juju.is/docs/olm/juju-supported-clouds)
for more information.
[^2]: All the commands were tested from a Ubuntu 22.04 LTS machine.
