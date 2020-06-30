---
id: deploy_intro
title: Introduction
hide_title: true
---
# Deploying Orchestrator: Introduction

Before deployment, it may be useful to read through the
[prerequisites](../basics/prerequisites.md) and
[quick start guide](../basics/quick_start_guide.md) sections.

These pages will walk through the full process of spinning up a full
Orchestrator deployment, from building the various containers that you'll need
to deploying them onto Amazon Elastic Kubernetes Service (EKS). This
installation guide targets *production* environments - if you aren't ready for
this, the developer documentation will be up shortly.

Familiarity with the AWS console and the Kubernetes command line are expected.
The instructions in this section have been tested on MacOS and Linux. If you
are deploying from a Windows host, some shell commands may require adjustments.

If you want to get a head start on the development setup, you can build the
Orchestrator containers following this guide and use docker-compose at
`magma/orc8r/cloud/docker` to spin up the local version of Orchestrator.
