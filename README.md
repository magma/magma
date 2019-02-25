# Magma

[![Build Status](https://travis-ci.com/facebookincubator/magma.svg)](https://travis-ci.com/facebookincubator/magma.svg)

Magma is an open-source software platform that gives network operators an open, flexible and extendable mobile core network solution. Magma enables better connectivity by:

* Allowing operators to offer cellular service without vendor lock-in with a modern, open source core network
* Enabling operators to manage their networks more efficiently with more automation, less downtime, better predictability, and more agility to add new services and applications
* Enabling federation between existing MNOs and new infrastructure providers for expanding rural infrastructure
* Allowing operators who are constrained with licensed spectrum to add capacity and reach by using Wi-Fi and CBRS


## Magma Architecture

The figure below shows the high-level Magma architecture. Magma is 3GPP generation (2G, 3G, 4G or upcoming 5G networks) and access network agnostic (cellular or WiFi). It can flexibly support a radio access network with minimal development and deployment effort.

Magma has three major components:

* **Access Gateway:** The Access Gateway (AGW) provides network services and policy enforcement. In an LTE network, the AGW implements an evolved packet core (EPC), and a combination of an AAA and a PGW. It works with existing, unmodified commercial radio hardware.

* **Orchestrator:** Orchestrator is a cloud service that provides a simple and consistent way to configure and monitor the wireless network securely. The Orchestrator can be hosted on a public/private cloud. The metrics acquired through the platform allows you to see the analytics and traffic flows of the wireless users through the Magma web UI.

* **Federation Gateway:** The Federation Gateway integrates the MNO core network with Magma by using standard 3GPP interfaces to existing MNO components.  It acts as a proxy between the Magma AGW and the operator's network and facilitates core functions, such as authentication, data plans, policy enforcement, and charging to stay uniform between an existing MNO network and the expanded network with Magma.

![Magma architecture diagram](docs/images/magma_overview.png?raw=true "Magma Architecture")

## Developer Prereqs

We develop all components of the system on virtual machines managed by
[Vagrant](https://www.vagrantup.com/). This helps us ensure that every
developer has a consistent development for cloud and gateway development. We 
support macOS and Linux host operating systems, and developing on Windows
should be possible but has not been tested.

First, install [Virtualbox](https://www.virtualbox.org/wiki/Downloads) and 
[Vagrant](http://www.vagrantup.com/downloads.html) so you can set up your
development VMs.

Then, install some additional prereqs (replace `brew` with your OS-appropriate
package manager as necessary):

```console
$ brew install python3
$ pip3 install ansible fabric3 requests
$ vagrant plugin install vagrant-vbguest
```

## Running the System

With the prereqs out of the way, we can now set up a minimal end-to-end system
on your development environment. In this README, we'll start by running the
LTE access gateway and orchestrator cloud inside their development VM's, and
registering your local access gateway with your local cloud for management.

We will be spinning up 3 virtual machines for this full setup, so you'll 
probably want to do this on a system with at least 8GB of memory. Our 
development VM's are in the 192.168.80.0/24 address space, so make sure that
you don't have anything running which hijacks that (e.g. VPN).

In the following steps, note the prefix in terminal commands. `HOST` means to
run the indicated command on your host machine, `CLOUD-VM` on the `cloud`
vagrant machine under `orc8r/cloud`, and `MAGMA-VM` on the `magma` vagrant
machine under `lte/gateway`.

#### Provisioning the Virtual Machines

Go ahead and open up 2 fresh terminal tabs. Start in

##### Terminal Tab 1:

We'll be setting up the LTE access gateway VM here.

```console
HOST [magma]$ cd lte/gateway
HOST [magma/lte/gateway]$ vagrant up
```

This will take a few minutes. Switch over to

##### Terminal Tab 2:

Here, we'll be setting up the orchestrator cloud and datastore VMs. The datastore VM runs a PostgreSQL instance.

```console
HOST [magma]$ cd orc8r/cloud
HOST [magma/orc8r/cloud]$ vagrant up
```

An Ansible error about a `.cache/test_certs` directory being missing is benign.
If you see an error early on about VirtualBox guest additions failing to 
install, you have a few extra steps to take after the provisioning fails:

```console
HOST [magma/orc8r/cloud]$ vagrant reload cloud
HOST [magma/orc8r/cloud]$ vagrant provision cloud
```

This will re-run the Ansible provisioning. It'll be a lot quicker this time 
around since most of the steps will be skipped.

#### Initial Build

Once all those jobs finish (should only be 5 minutes or so), we can build the
code for the first time and run it. Start in

##### Terminal Tab 1:

```console
HOST [magma/lte/gateway]$ vagrant ssh
MAGMA-VM [/home/vagrant]$ cd magma/lte/gateway
MAGMA-VM [/home/vagrant/magma]$ make run
```
This will take a while (we have a lot of CXX files to build). It's a good time
to switch over to

##### Terminal Tab 2:

Here, we'll build and run all of our orchestrator cloud services. We'll also 
use this time to register the client certificate you'll need to access the API 
gateway for the controller running on your cloud VM.

```console
HOST [magma/orc8r/cloud]$ vagrant ssh
CLOUD-VM [/home/vagrant]$ cd magma/orc8r/cloud
CLOUD-VM [/home/vagrant/magma/orc8r/cloud]$ make run
```

Again, this will take a while. While both these builds are churning, it's a 
good time to grab lunch, or a coffee.

If you're reading this because you recently nuked your VMs, you have one more
thing to do. This `magma/.cache` folder stays around even when you destroy
your VMs. But if you're starting with a blank datastore VM, the cert that
your browser has isn't registered with the system anymore
(you deleted those tables). We have a Make target that will re-register the
cached test cert with your datastore VM:

```console
CLOUD-VM [/home/vagrant]$ cd magma/orc8r/cloud
CLOUD-VM [/home/vagrant/magma/orc8r/cloud]$ make restore_admin_operator
```

#### Connecting Your Local LTE Gateway to Your Local Cloud

At this point, you will have built all the code in the LTE access gateway and
the orchestrator cloud. All the services on the LTE access gateway and
orchestrator cloud are running, but your gateway VM isn't yet set up to
communicate with your local cloud VM.

We have a fabric command set up to do this:

```console
HOST [magma]$ cd lte/gateway
HOST [magma/lte/gateway]$ fab register_vm
```

At this point, your access gateway VM is streaming configuration from your
cloud VM and sending status and metrics back to your cloud VM.

If you want to see what the access gateway is doing, you can
`vagrant ssh magma` inside `lte/gateway`, then do a
`sudo tail -f /var/log/syslog`. If everything above went smoothly, you should
eventually (give it a minute or two) see something along the lines of:

```console
MAGMA-VM$ sudo service magma@* stop
MAGMA-VM$ sudo service magma@magmad restart
MAGMA-VM$ sudo tail -f /var/log/syslog
Sep 27 22:57:35 magma-dev magmad[6226]: [2018-09-27 22:57:35,550 INFO root] Checkin Successful!
Sep 27 22:57:55 magma-dev magmad[6226]: [2018-09-27 22:57:55,684 INFO root] Processing config update g1
Sep 27 22:57:55 magma-dev control_proxy[6418]: 2018-09-27T22:57:55.683Z [127.0.0.1 -> streamer-controller.magma.test,8443] "POST /magma.Streamer/GetUpdates HTTP/2" 200 7bytes 0.009s
```

#### Using the NMS UI

Magma provides an UI for configuring and monitoring the networks. To run the UI, first install [Docker](https://www.docker.com/) in your host. Then:

```console
HOST [magma]$ cd nms/fbcnms-projects/magmalte
HOST [magma/nms/fbcnms-projects/magmalte] $
docker-compose up -d
HOST [magma/nms/fbcnms-projects/magmalte] $
docker-compose run magmalte yarn run setAdminPassword admin@magma.test password1234
```

After this, you will be able to access the UI by visiting [https://localhost](https://localhost), and using the email `admin@magma.test` and password `password1234`.

## Join the Magma Community

- Mailing lists:
  - Join [magma-dev](https://groups.google.com/forum/#!forum/magma-dev) for technical discussions
  - Join [magma-announce](https://groups.google.com/forum/#!forum/magma-announce) for announcements
- Discord:
  - [magma\_dev](https://discord.gg/WDBpebF) channel

See the [CONTRIBUTING](CONTRIBUTING.md) file for how to help out.

## License

Magma is BSD License licensed, as found in the LICENSE file.
The EPC is OAI is offered under the OAI Apache 2.0 license, as found in the LICENSE file in the OAI directory.
