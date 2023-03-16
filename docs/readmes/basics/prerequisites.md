---
id: prerequisites
title: Prerequisites
hide_title: true
---

# Prerequisites

These are the prerequisites to setting up a full private LTE Magma deployment.
Additional prerequisites for developers can be found in the [contributor guide on Github](https://github.com/magma/magma/wiki/Contributor-Guide).

## Operating System

Currently, the main development operating system (OS) is **macOS**. Documentation is mainly focused on that operating system.
To develop on a **Linux OS**, the package manager (brew for macOS) will need to be replaced by the appropriate package manager for the respective Linux distribution (e.g. apt, yum, etc.).
**Windows OS** is currently _not_ supported as developing environment, due to some dependencies on Linux-only tools during setup, such as Ansible or `fcntl`. You can try to use a [DevContainer setup](https://github.com/magma/magma/wiki/Contributing-Code-with-VSCode#using-devcontainer-for-development) though.

## Development Tools

Development can occur from multiple OS's, where **macOS** and **Ubuntu** are **explicitly supported**, with additional polish for macOS.

**Note:** If you still want to contribute from a different OS, you will need to figure out some workarounds to install the tooling. You might want to follow one of the guides, either macOS or Ubuntu, and replicate the steps in your preferred OS.

### macOS

1. Install the following tools

   1. [Docker and Docker Compose](https://docs.docker.com/desktop/install/mac-install/)
   2. [Homebrew](https://brew.sh/)
   3. [VirtualBox](https://www.virtualbox.org/)
   4. [Vagrant](https://vagrantup.com)

   ```bash
   brew install go@1.20 pyenv
   # NOTE: this assumes you're using zsh.
   # See the above pyenv install instructions if using alternative shells.
   echo 'export PATH="/usr/local/opt/go@1.20/bin:$PATH"' >> ~/.zshrc
   echo 'eval "$(pyenv init --path)"' >> ~/.zprofile
   echo 'eval "$(pyenv init -)"' >> ~/.zshrc
   exec $SHELL
   # IMPORTANT: close your terminal tab and open a new one before continuing
   pyenv install 3.8.10
   pyenv global 3.8.10
   pip3 install ansible fabric jsonpickle requests PyYAML
   vagrant plugin install vagrant-vbguest vagrant-disksize vagrant-reload
   ```

   **Note**: In the case where installation of `fabric` through pip was unsuccessful,
   try switching to other package installers. Try running `brew install fabric`.

   You should start Docker Desktop and increase the memory
   allocation for the Docker engine to at least 4GB (Preferences -> Resources ->
   Advanced). If you are running into build/test failures with Go that report
   "signal killed", you likely need to increase Docker's allocated resources.

   ![Increasing docker engine resources](assets/docker-config.png)

### Ubuntu

1. Install the following tools
   1. [Docker and Docker Compose](https://docs.docker.com/engine/install/ubuntu/)
   2. [VirtualBox](https://www.virtualbox.org/wiki/Linux_Downloads)
   3. [Vagrant](https://www.vagrantup.com/downloads)
2. Install golang version 1.20.1.

   1. Download the tar file.

      ```bash
      wget https://linuxfoundation.jfrog.io/artifactory/magma-blob/go1.20.1.linux-amd64.tar.gz
      ```

   2. Extract the archive you downloaded into `/usr/local`, creating a Go tree in `/usr/local/go`.

      ```bash
      sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.20.1.linux-amd64.tar.gz
      ```

   3. Add `/usr/local/go/bin` to the PATH environment variable.

      ```bash
      export PATH=$PATH:/usr/local/go/bin
      ```

   4. Verify that you've installed Go by opening a command prompt and typing the following command

      ```bash
      go version
      ```

      You should expect something like this

      ```bash
      go version go1.20.1 linux/amd64
      ```

3. Install `pyenv`.

   1. Update system packages.

      ```bash
      sudo apt update -y
      ```

   2. Install some necessary dependencies. **If you are using `zsh` instead of `bash`, replace** `.bashrc` **for** `.zshrc`.

      ```bash
      apt install -y make build-essential libssl-dev zlib1g-dev libbz2-dev    libreadline-dev libsqlite3-dev wget curl llvm libncurses5-dev  libncursesw5-dev xz-utils tk-dev libffi-dev liblzma-dev python-openssl git
      ```

      **Note**: For Ubuntu 22.04, use `python3-openssl` instead of `python-openssl`.

   3. Clone `pyenv` repository.

      ```bash
      git clone https://github.com/pyenv/pyenv.git ~/.pyenv
      ```

   4. Configure `pyenv`.

      ```bash
      echo 'export PYENV_ROOT="$HOME/.pyenv"' >> ~/.bashrc
      echo 'export PATH="$PYENV_ROOT/bin:$PATH"' >> ~/.bashrc
      echo -e 'if command -v pyenv 1>/dev/null 2>&1; then\n eval "$(pyenv init -) "\nfi' >> ~/.bashrc
      exec "$SHELL"
      ```

   5. Create python virtual environment version 3.8.10.

      ```bash
      pyenv install 3.8.10
      pyenv global 3.8.10
      ```

        **Note**: The `pyenv` installation [might fail with a segmentation fault](https://github.com/pyenv/pyenv/issues/2046). Try using `CFLAGS="-O2" pyenv install 3.8.10` in that case.

4. Install `pip3` and its dependencies.

   1. Install `pip3`.

      ```bash
      sudo apt install python3-pip
      ```

   2. Install the following dependencies

      ```bash
      pip3 install ansible fabric jsonpickle requests PyYAML
      ```

5. Install `vagrant` necessary plugin.

   ```bash
   vagrant plugin install vagrant-vbguest vagrant-disksize vagrant-reload
   ```

    Make sure `virtualbox` is the default provider for `vagrant` by adding the following line to your `.bashrc` (or equivalent) and restart your shell: `export VAGRANT_DEFAULT_PROVIDER="virtualbox"`.

## Downloading Magma

You can find Magma code on [Github](https://github.com/magma/magma).

To download Magma current version, or a specific release do the following

```bash
git clone https://github.com/magma/magma.git
cd magma

# in case you want to use a specific version of Magma (for example v1.8)
git checkout v1.8

# to list all available releases
git tag -l
```

## Deployment Tooling

First, follow the previous section on [developer tools](#development-tools). Then, install some
additional prerequisite tools.

### macOS

Install necessary dependencies and configure the aws cli

```bash
brew install aws-iam-authenticator kubectl helm terraform
python3 -m pip install awscli boto3
aws configure
```

### Ubuntu

Install the following

1. [aws-iam-authenticator for Linux](https://docs.aws.amazon.com/eks/latest/userguide/install-aws-iam-authenticator.html).
2. [kubectl for Linux](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/#install-using-native-package-management).
3. [Helm for Linux](https://helm.sh/docs/intro/install/).
4. [Terraform for Linux](https://learn.hashicorp.com/tutorials/terraform/install-cli).
5. awscli

   ```bash
   sudo apt install awscli
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
requirements for the AGW production hardware are

1. 2+ physical Ethernet interfaces
2. AMD64 dual-core processor around 2GHz clock speed or faster
3. 4GB RAM
4. 32GB or greater SSD storage

In addition, in order to build the AGW, you should have on hand

1. A USB stick with 2GB+ capacity to load a Debian Stretch ISO
2. Peripherals (keyboard, screen) for your production AGW box for use during
   provisioning

### RAN Equipment

We currently have tested with the following EnodeB's

1. Baicells Nova 233 TDD Outdoor
2. Baicells Nova 243 TDD Outdoor
3. Assorted Baicells indoor units (for lab deployments)

Support for other RAN hardware can be implemented inside the `enodebd` service
on the AGW, but we recommend starting with one of these EnodeBs.
