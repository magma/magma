---
id: version-1.7.0-setup
title: AGW Setup (With Vagrant)
sidebar_label: Setup (With Vagrant)
hide_title: true
original_id: setup
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

- Once the Access Gateway VM is up and provisioned, run the following commands:

```text
HOST:magma/lte/gateway USER$ vagrant ssh magma
AGW:~ USER$ cd magma/lte/gateway
AGW:~/magma/lte/gateway USER$ make run
```

Once the Access Gateway VM is running successfully, proceed to attaching the eNodeB.
