---
id: version-1.0.1-setup
title: CWAG Setup (With Vagrant)
sidebar_label: Setup (With Vagrant)
hide_title: true
original_id: setup
---
# CWF Access Gateway Setup (With Vagrant)
### Prerequisites
To develop and manage a Magma VM, you must have the following applications installed locally:

* Virtualbox
*  Vagrant
* Ansible

### Steps

To bring up a Wifi Access Gateway (CWAG) VM using Vagrant:

* Run the following command:

``HOST:magma/cwf/gateway USER$ vagrant up cwag``

Vagrant will bring up the VM, then Ansible will provision the VM.


* Once the CWAG VM is up and provisioned, run the following commands:

``HOST:magma/cwf/gateway USER$ vagrant ssh cwag``<br>
``AGW:~ USER$ cd magma/cwf/gateway/docker``<br>
``AGW:~/magma/cwf/gateway/docker USER$ docker-compose build --parallel``
``AGW:~/magma/cwf/gateway/docker USER$ docker-compose up -d``

After this, all the CWAG docker containers should have been brought up 
successfully.
