---
id: setup_deb
title: AGW Setup (Bare Metal)
sidebar_label: Setup (Bare Metal)
hide_title: true
---
# Access Gateway Setup (On Bare Metal)
### Prerequisites

To setup a Magma Access Gateway, you will need a machine that 
satisfies the following requirements:

###### Machine running Debian Stretch

Officially, we support kernel release 4.9.0.8, version 4.9.130-2 
compiled in October of 2018.
To check the kernel version running on your machine:

```
HOST:~ USER$ uname -a

Linux debian 4.9.0-8-amd64 #1 SMP Debian 4.9.130-2 (2018-10-27) x86_64 GNU/Linux
```

###### Two Network Adapters

Two network adapters are required for a Magma installation to work 
successfully.
Typically, an ethernet port will be used to connect to your eNodeB.
Magma expects the eNodeB to be connected to the 'eth1' interface
Your second network interface should connect to the internet and 
provide the machine with the SGi interface.


### Setup

###### 1. Access Keys to Magma package repository

First, you need access to the Magma package repository:

``HOST:~ USER$ apt-key adv --fetch-keys 
http://packages.magma.etagecom.io/pubkey.gpg``

###### 2. Add Magma package repository to apt

After getting access, then you'll want to add our package 
repository to Apt:

``HOST:~ USER$ echo 'deb http://packages.magma.etagecom.io 
stretch-stable main' > 
/target/etc/apt/sources.list.d/packages_magma_etagecom_io.list``

###### 3. Get and install Magma

And now to get and install Magma:

``
HOST:~ USER$ apt-get update
HOST:~ USER$ apt-get install magma
``

Magma AGW should now be installed on your machine.

###### 4. Start Magma services

Start the magmad service, which will kick off the rest of the 
services.

* Start the magmad service:

``HOST:~ USER$ sudo service magma@magmad start``

* To check status of Magma services

``HOST:~ USER$ sudo service magma@* status``

###### 5. Configure your access network

If all has gone well, then you have a working Magma Access Gateway.

If you haven't already, get the Orc8r up and running.
You can find the instructions to do so on the sidebar.

From here, you'll need to configure your controller and gateway for 
your desired setup.

Access the Swagger page, which should be running if your Orc8r 
instance is up and running properly.
You can find it at <https://localhost:9443/apidocs> 
or <https://192.169.99.99:9443>

For bare minimum configuration, you'll want to configure the 
following:

* /networks/{network_id}/gateways/{gateway_id}/configs
* /networks/{network_id}/gateways/{gateway_id}/configs/cellular
* /networks/{network_id}/configs/cellular
* /networks/{network_id}/configs/dns

For the most part, configurations can all be default.
You'll want to pay attention to ``/networks/{network_id}/configs/cellular`` as 
the configuration depends on the eNodeB you are using.

###### 6. Connect your eNodeB

Magma Access Gateway expects eNodeB to be connected to the 'eth1' 
interface.

For more information about support, connecting, and configuring 
your eNodeB, find the 'eNodeB Configuration' page.
