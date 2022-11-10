---
id: setup
title: AGW Setup (With Vagrant)
sidebar_label: Setup (With Vagrant)
hide_title: true
---

# Access Gateway Setup (With Vagrant)

## Prerequisites

To develop and manage a Magma VM, you must have the following applications installed locally:

- Virtualbox
- Vagrant
- Ansible

## Steps

To bring up an Access Gateway (AGW) VM using Vagrant:

- Run the following command:

``HOST:magma/lte/gateway USER$ vagrant up magma``

Vagrant will bring up the VM, then Ansible will provision the VM.

- Once the Access Gateway VM is up and provisioned, run the following command to ssh into the VM:

``HOST:magma/lte/gateway USER$ vagrant ssh magma``

- Next follow the steps below to get the magma services running:

1. Create links for cli scripts: `cd $MAGMA_ROOT && bazel/scripts/link_scripts_for_bazel_integ_tests.sh`
2. Use bazel systemd unit files: `sudo cp $MAGMA_ROOT/lte/gateway/deploy/roles/magma/files/systemd_bazel/* /etc/systemd/system/ && sudo systemctl daemon-reload`
3. Build the services: `cd $MAGMA_ROOT && magma-build-agw` (Note: this will take some time for the initial build, but will be fast for follow-up builds.)
4. Start the access gateway: `magma-restart`

Once the Access Gateway VM is running successfully, proceed to attaching the eNodeB.
