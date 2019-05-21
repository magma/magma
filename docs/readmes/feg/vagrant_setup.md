---
id: vagrant_setup
title: Federated Gateway (Vagrant Setup)
sidebar_label: FeG Vagrant Setup
hide_title: true
---
## Prerequisites

To develop and manage a Magma VM, you must have the following applications installed locally:

* Virtualbox
* Vagrant
* Ansible

To setup a Federated Gateway (FeG), you must also have Orchestrator (orc8r) installed. To install and build Orchestrator (orc8r), follow the instructions at facebookincubator (https://github.com/facebookincubator)/magma (https://github.com/facebookincubator/magma).

A production Magma system requires an access gateway. However, for FeG development, you do not need an access gateway VM. You only need the following VMs:

* orc8r datastore
* orc8r cloud
* feg

The Magma VMs are stored in the 192.168.80.0/24 address space. So, make sure you do not have anything running in that address space (such as a VPN).

## Steps

These steps describe how to:

* Provision the Cloud VM
* Provision the FeG VM
* Connect the FeG to the Cloud
* Monitor the Federated Gateway

**Note:** In these steps, the following terminal command prefixes are used:

* HOST - means run the command at the local host machine
* CLOUD-VM - means run the command on the cloud vagrant machine under orc8r/cloud
* FEG-VM - means run the command on the feg vagrant machine under feg/gateway


To provision the Cloud VM:

1. In terminal, open two tabs.
2. In a the first terminal tab, run the following commands:

``HOST [magma]$ cd orc8r/cloud``<br>
``HOST [magma/orc8r/cloud]$ vagrant up datastore``<br>
``HOST [magma/orc8r/cloud]$ vagrant up cloud``<br>
``HOST [magma/orc8r/cloud]$ vagrant ssh cloud``<br>
``CLOUD-VM [/home/vagrant]$ cd magma/orc8r/cloud``<br>
``CLOUD-VM [/home/vagrant/magma/orc8r/cloud]$ make run``<br>

To provision the Feg VM:

* In the second terminal tab, run the following commands:
    *Note:* The *vagrant up feg* command may take several minutes to run.

``HOST [magma]$ cd feg/gateway``<br>
``HOST [magma/feg/gateway]$ vagrant up feg``<br>
``HOST [magma/feg/gateway]$ vagrant ssh feg``<br>
``FEG-VM [/home/vagrant]$ cd magma/feg/gateway``<br>
``FEG-VM [/home/vagrant]$ make run``<br>

At this point, all services in the Federated Gateway and Orchestrator should be running.

To connect the FeG to the Cloud:

* Run the following fabric commands:

``HOST [magma]$ cd feg/gateway``<br>
``HOST [magma/feg/gateway]$ fab register_feg_vm``

At this point, the FeG VM should be receiving configuration settings from the
Cloud VM, and sending status, health, and metrics back to the Cloud VM.

To Monitor the Federated Gateway:

* Run the following commands inside feg/gateway:

``vagrant ssh feg``<br>
``sudo tail -f /var/log/syslog. ``

After a few minutes, it should display syslog information similar to this:

``FEG-VM$ sudo service magma@* stop``<br>
``FEG-VM$ sudo service magma@magmad start``<br>
``FEG-VM$ sudo tail -f /var/log/syslog``<br>
``Sep 27 22:57:35 magma-feg-dev magmad[6226]: [2018-09-27 22:57:35,550 INFO root] Checkin Successful!``<br>
``Sep 27 22:57:55 magma-feg-dev magmad[6226]: [2018-09-27 22:57:55,684 INFO root] Processing config update g1``<br>
``Sep 27 22:57:55 magma-feg-dev control_proxy[6418]: 2018-09-27T22:57:55.683Z [127.0.0.1 -> streamer-controller.magma.test,8443] "POST /magma.orc8r.Streamer/GetUpdates HTTP/2" 200 7bytes 0.009s``



```
