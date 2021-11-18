---
id: prerequisites
title: Prerequisites
hide_title: true
---
# Prerequisites

These are the prerequisites to setting up a full private LTE Magma deployment.
Additional prerequisites for developers can be found in the [developer's guide](../contributing/contribute_onboarding.md).

## Operating System

Currently, the main development operating system (OS) is **macOS**. Documentation is mainly focused on that operating system.
To develop on a **Linux OS**, the package manager (brew for macOS) will need to be replaced by the appropriate package manager for the respective Linux distribution (e.g. apt, yum, etc.).
**Windows OS** is currently _not_ supported as developing environment, due to some dependencies on Linux-only tools during setup, such as Ansible or `fcntl`. You can try to use a [DevContainer setup](../contributing/contribute_vscode.md#open-a-devcontainer-workspace-with-github-codespaces) though.

## Development Tools

Install the following tools:

1. [Docker](https://www.docker.com) and Docker Compose
2. [Homebrew](https://brew.sh/) *only* for macOS users
3. [VirtualBox](https://www.virtualbox.org/)
4. [Vagrant](https://vagrantup.com)

Replace `brew` with your OS-appropriate package manager as necessary, or see
the [pyenv installation instructions](https://github.com/pyenv/pyenv#installation).

```bash
brew install go@1.13 pyenv

# NOTE: this assumes you're using zsh.
# See the above pyenv install instructions if using alternative shells.
echo 'export PATH="/usr/local/opt/go@1.13/bin:$PATH"' >> ~/.zshrc
echo 'eval "$(pyenv init --path)"' >> ~/.zprofile
echo 'eval "$(pyenv init -)"' >> ~/.zshrc
exec $SHELL

# IMPORTANT: close your terminal tab and open a new one before continuing

pyenv install 3.8.10
pyenv global 3.8.10

pip3 install ansible fabric3 jsonpickle requests PyYAML
vagrant plugin install vagrant-vbguest
```

**Note**: In the case where installation of `fabric3` through pip was unsuccessful,
try switching to other package installers. For example, for macOS users, try
running `brew install fabric`.

If you are on macOS, you should start Docker for Mac and increase the memory
allocation for the Docker engine to at least 4GB (Preferences -> Resources ->
Advanced). If you are running into build/test failures with Go that report
"signal killed", you likely need to increase Docker's allocated resources.

![Increasing docker engine resources](assets/docker-config.png)

## Downloading Magma

You can find Magma code on [Github](https://github.com/magma/magma).

To download Magma current version, or a specific release do the following:

```bash
git clone https://github.com/magma/magma.git
cd magma

# in case you want to use a specific version of Magma (for example v1.6)
git checkout v1.6

# to list all available releases
git tag -l
```

## Build/Deploy Tooling

We support building the AGW and Orchestrator on macOS and Linux host operating
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
2. Registered domain for Orchestrator endpoints

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
which can support a Debian or Ubuntu 20.04 Linux installation. The basic system
requirements for the AGW production hardware are:

1. 2+ physical Ethernet interfaces
2. AMD64 dual-core processor around 2GHz clock speed or faster
3. 4GB RAM
4. 32GB or greater SSD storage

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
