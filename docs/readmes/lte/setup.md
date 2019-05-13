---
id: setup
title: AGW Setup (With Vagrant)
sidebar_label: Setup (With Vagrant)
hide_title: true
---
# Access Gateway Setup (With Vagrant)
### Prerequisites
Make sure that you have Virtualbox and Vagrant installed.
We use this to manage our development VMs.

### Overview
To bring up the AGw VM with Vagrant, simply run the following:
```
HOST:magma/lte/gateway USER$ vagrant up magma
```
Vagrant will bring up the AGw VM, and then ansible will provision it.

Once the VM is brought up and provisioned successfully, you can ssh into
the machine, and build and run the code
```
HOST:magma/lte/gateway USER$ vagrant ssh magma
AGW:~ USER$ cd magma/lte/gateway
AGW:~/magma/lte/gateway USER$ make run
```

If it has all run successfully, then you are ready to attach an eNodeB.
