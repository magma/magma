---
id: version-1.2.0-prerequisites
title: Prerequisites
hide_title: true
original_id: prerequisites
---
# Prerequisites

These are the prerequisites to setting up a full private LTE Magma deployment.
Additional prerequisites for developers can be found in the developer's guide.

## Development Tools

Install the following tools:

1. [Docker](https://www.docker.com) and Docker Compose
2. [Homebrew](https://brew.sh/) *only* for MacOS users
3. [VirtualBox](https://www.virtualbox.org/)
3. [Vagrant](https://vagrantup.com)

Replace `brew` with your OS-appropriate package manager as necessary, or see
the [pyenv installation instructions](https://github.com/pyenv/pyenv#installation).

```bash
brew install pyenv

# Replace .zshrc with your appropriate shell RC file
# IMPORTANT: Use .bash_profile, not .bashrc for bash
echo -e 'if command -v pyenv 1>/dev/null 2>&1; then\n  eval "$(pyenv init -)"\nfi' >> ~/.zshrc
exec $SHELL
pyenv install 3.7.3
pyenv global 3.7.3

pip3 install ansible fabric3 jsonpickle requests PyYAML
vagrant plugin install vagrant-vbguest
```

If you are on MacOS, you should start Docker for Mac and increase the memory
allocation for the Docker engine to at least 4GB (Preferences -> Advanced).

![Increasing docker engine resources](assets/docker-config.png)

## Build/Deploy Tooling

We support building the AGW and Orchestrator on MacOS and Linux host operating
systems. Doing so on a Windows environment should be possible but has not been
tested. You may prefer to use a Linux virtual machine if you are on a Windows
host.

First, follow the previous section on developer tools. Then, install some
additional prerequisite tools (replace `brew` with your OS-appropriate package
manager as necessary):

```bash
brew install aws-iam-authenticator kubectl helm terraform
python3 -m pip install awscli boto3
aws configure
```

### Orchestrator and NMS

Orchestrator deployment depends on the following components

1. AWS account
2. Docker image repository (e.g. Docker Hub, JFrog)
3. Helm chart repository (e.g. JFrog, Github)*
4. Registered domain for Orchestrator endpoints

\* We describe setting up a private GitHub repository as a Helm repository in
the [building Orchestrator](../orc8r/deploy_build.md) section.

We recommend deploying the Orchestrator cloud component of Magma into AWS.
Our open-source Terraform scripts target an AWS deployment environment, but if
you are familiar with devops and are willing to roll your own, Orchestrator can
run on any public/private cloud with a Kubernetes cluster available to use.
The deployment documentation will assume an AWS deployment environment - if
this is your first time using or deploying Orchestrator, we recommend that you
follow this guide before attempting to deploy it elsewhere.

Provide the access key ID and secret key for an administrator user in AWS
(don't use the root user) when prompted by `aws configure`. Skip this step if
you will use something else for managing AWS credentials.

## Production Hardware

### Access Gateways

Access gateways (AGWs) can be deployed on to any AMD64 architecture machine
which can support a Debian Linux installation. The basic system requirements
for the AGW production hardware are:

1. 2+ physical ethernet interfaces
2. AMD64 dual-core processor around 2GHz clock speed or faster
3. 2GB RAM
4. 128GB-256GB SSD storage

In addition, in order to build the AGW, you should have on hand:

1. A USB stick with 2GB+ capacity to load a Debian Stretch ISO
2. Peripherals (keyboard, screen) for your production AGW box for use during
provisioning

### RAN Equipment

We currently have tested with the following EnodeB's:

1. Baicells Nova 233 TDD Outdoor
2. Baicells Nova 243 TDD Outdoor
3. Assorted Baicells indoor units (for lab deployments)

Support for other RAN hardware can be implemented inside the `enodebd` service
on the AGW, but we recommend starting with one of these EnodeBs.
