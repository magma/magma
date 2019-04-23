---
id: quick_start_guide
title: Quick Start Guide
hide_title: true
---
# Quick Start Guide

With the [prereqs](prerequisites.md) installed, we can now set up a minimal 
end-to-end system on your development environment. In this guide, we'll start 
by running the LTE access gateway and orchestrator cloud, and then
register your local access gateway with your local cloud for management.

We will be spinning up a virtual machine and some docker containers for this 
full setup, so you'll probably want to do this on a system with at least 8GB 
of memory. Our development VM's are in the 192.168.80.0/24 address space, so
make sure that you don't have anything running which hijacks that (e.g. VPN).

In the following steps, note the prefix in terminal commands. `HOST` means to
run the indicated command on your host machine, and `MAGMA-VM` on the `magma`
vagrant machine under `lte/gateway`.

#### Provisioning the environment

Go ahead and open up 2 fresh terminal tabs. Start in

##### Terminal Tab 1:

We'll be setting up the LTE access gateway VM here.

```console
HOST [magma]$ cd lte/gateway
HOST [magma/lte/gateway]$ vagrant up
```

This will take a few minutes to spin up the VMs. Switch over to..

##### Terminal Tab 2:

Here, we'll be setting up the orchestrator docker containers.

```console
HOST [magma]$ cd orc8r/cloud/docker
HOST [magma/orc8r/cloud/docker]$ ./build.py
```

This will build the docker images for the orc8r.

#### Initial Run

Once all those jobs finish (should only be 5 minutes or so), we can build the
access gateway code and run it. Start in..

##### Terminal Tab 1:

```console
HOST [magma/lte/gateway]$ vagrant ssh
MAGMA-VM [/home/vagrant]$ cd magma/lte/gateway
MAGMA-VM [/home/vagrant/magma/lte/gateway]$ make run
```
This will take a while (we have a lot of CXX files to build).

##### Terminal Tab 2:

Here, we'll run all of our orchestrator cloud services. We'll also 
use this time to register the client certificate you'll need to access the API 
gateway for the controller running on your cloud VM.

```console
HOST [magma/orc8r/cloud/docker]$ docker-compose up -d
```

#### Connecting Your Local LTE Gateway to Your Local Cloud

At this point, you will have built all the code in the LTE access gateway and
the orchestrator cloud. All the services on the LTE access gateway and
orchestrator cloud are running, but your gateway VM isn't yet set up to
communicate with your local cloud.

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
HOST [magma/nms/fbcnms-projects/magmalte] $ docker-compose up -d
HOST [magma/nms/fbcnms-projects/magmalte] $ ./scripts/dev_setup.sh
```

After this, you will be able to access the UI by visiting [https://localhost](https://localhost), and using the email `admin@magma.test` and password `password1234`. If you see Gateway Error 502, don't worry, the NMS can take upto 60 seconds to finish starting up.
