---
id: readme
title: Federated Gateway (FeG)
sidebar_label: Overview
hide_title: true
---
# Federated Gateway (FeG)
The federated gateway provides remote procedure call (GRPC) based interfaces to standard 3GPP components, such as 
HSS (S6a, SWx), OCS (Gy), and PCRF (Gx). The exposed RPC interface provides versioning & backward compatibility, 
security (HTTP2 & TLS) as well as support for multiple programming languages. The Remote Procedures below provide 
simple, extensible, multi-language interfaces based on GRPC which allow developers to avoid dealing with the 
complexities of 3GPP protocols. Implementing these RPC interfaces allows networks running on Magma to integrate 
with traditional 3GPP core components.

![Federated Gateway architecture diagram](https://github.com/facebookincubator/magma/blob/master/docs/readmes/assets/federated_gateway_diagram.png?raw=true "FeG Architecture")

The Federated Gateway supports the following features and functionalities:

1. Hosting centralized control plane interface towards HSS, PCRF, OCS and MSC/VLR on behalf of distributed AGW/EPCs.
2. Establishing diameter connection with HSS, PCRF and OCS directly as 1:1 or via DRA. 
3. Establishing SCTP/IP connection with MSC/VLR.
4. Interfacing with AGW over GPRC interface by responding to remote calls from EPC (MME and Sessiond/PCEF) components,
    converting these remote calls to 3GPP compliant messages and then sending these messages to the appropriate core network 
    components such as HSS, PCRF, OCS and MSC.  Similarly the FeG receives 3GPP compliant messages from HSS, PCRF, OCS and MSC 
    and converts these to the appropriate GPRC messages before sending them to the AGW. 



Please see the **[Magma Product Spec](https://github.com/facebookincubator/magma/blob/master/docs/Magma_Specs_V1.1.pdf)** for more detailed information.

## Running the System

First, ensure you have the necessary developer prerequisites and have built the orc8r as instructed [here](https://github.com/facebookincubator/magma#developer-prereqs).

While a production system would include the access gateway, orc8r and federated gateway, for feg development
only the orc8r *datastore*, *orc8r cloud* and *feg* VM's are necessary.

Now, we will spin up the federated gateway VM. Our development VM's are in the
192.168.80.0/24 address space, so make sure that you don't have anything running 
already on that IP (e.g. VPN).
 
#### Provisioning the Virtual Machine

In the following steps, note the prefix in terminal commands. `HOST` means to
run the indicated command on your host machine, `CLOUD-VM` on the `cloud`
vagrant machine under `orc8r/cloud`, and `FEG-VM` on the `feg` vagrant
machine under `feg/gateway`.

Open two new terminal instances. Start in

##### Terminal Tab 1:

```console
HOST [magma]$ cd orc8r/cloud
HOST [magma/orc8r/cloud]$ vagrant up datastore
HOST [magma/orc8r/cloud]$ vagrant up cloud
HOST [magma/orc8r/cloud]$ vagrant ssh cloud
CLOUD-VM [/home/vagrant]$ cd magma/orc8r/cloud
CLOUD-VM [/home/vagrant/magma/orc8r/cloud]$ make run
```

##### Terminal Tab 2:

We'll now provision the federated gateway dev VM:

```console
HOST [magma]$ cd feg/gateway
HOST [magma/feg/gateway]$ vagrant up feg
```

This will take some time. Once finished:

```console
HOST [magma/feg/gateway]$ vagrant ssh feg
FEG-VM [/home/vagrant]$ cd magma/feg/gateway
FEG-VM [/home/vagrant]$ make run
```

#### Connecting Your Local Federated Gateway to Your Local Cloud

At this point, you will have built all the code in the federated gateway and
the orchestrator cloud. All the services on the federated gateway and
orchestrator cloud are running, but your feg VM isn't yet set up to
communicate with your local cloud VM.

We have a fabric command set up to do this:

```console
HOST [magma]$ cd feg/gateway
HOST [magma/feg/gateway]$ fab register_feg_vm
```

At this point, your federated gateway VM is streaming configurations from your
cloud VM and sending status, health and metrics back to your cloud VM.

If you want to see what the federated gateway is doing, you can run
`vagrant ssh feg` inside `feg/gateway`. Then run `sudo tail -f /var/log/syslog`. 
If everything above went smoothly, you should eventually (give it a minute or two) see 
something along the lines of:

```console
FEG-VM$ sudo service magma@* stop
FEG-VM$ sudo service magma@magmad start
FEG-VM$ sudo tail -f /var/log/syslog
Sep 27 22:57:35 magma-feg-dev magmad[6226]: [2018-09-27 22:57:35,550 INFO root] Checkin Successful!
Sep 27 22:57:55 magma-feg-dev magmad[6226]: [2018-09-27 22:57:55,684 INFO root] Processing config update g1
Sep 27 22:57:55 magma-feg-dev control_proxy[6418]: 2018-09-27T22:57:55.683Z [127.0.0.1 -> streamer-controller.magma.test,8443] "POST /magma.orc8r.Streamer/GetUpdates HTTP/2" 200 7bytes 0.009s
```

## Federated Gateway Services & Tools

The following services run on the federated gateway:
 - `s6a_proxy` - translates calls from GRPC to S6a protocol between AGW and HSS 
 - `session_proxy` - translates calls from GRPC to gx/gy protocol between AGW and PCRF/OCS
 - `csfb` - translates calls from GRPC interface to csfb protocol between AGW and VLR
 - `swx_proxy` - translates GRPC interface to SWx protocol between AGW and HSS
 - `gateway_health` - provides health updates to the orc8r to be used for 
 achieving highly available federated gateway clusters (see **[Magma Product Spec](https://github.com/facebookincubator/magma/blob/master/docs/Magma_Specs_V1.1.pdf)**
 for more details)

Associated tools for sending requests and debugging issues can be found
at `magma/feg/gateway/tools`. 
