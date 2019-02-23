# Magma

[![Build Status](https://travis-ci.com/facebookexternal/magma.svg?token=rUhsJxd4NdhJ6GBKWCki&branch=master)](https://travis-ci.com/facebookexternal/magma)

What is magma/why magma?

Architecture diagram

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

Go ahead and open up 3 fresh terminal tabs. Start in

##### Terminal Tab 1:

We'll be setting up the LTE access gateway VM here.

```console
HOST$ cd lte/gateway
HOST$ vagrant up magma
```

This will take a few minutes. Switch over to

##### Terminal Tab 2:

Here, we'll be setting up the datastore VM. This VM runs a PostgreSQL instance 
that backs the orchestrator cloud VM.

```console
HOST$ cd orc8r/cloud
HOST$ vagrant up datastore
```

This will only take a few seconds, but feel free to switch over to

##### Terminal Tab 3:

Here, we'll be setting up your orchestrator cloud VM. This simulates our
deployment environment for the cloud-based management plane.

```console
HOST$ cd ~/fbsource/fbcode/magma/orc8r/cloud
HOST$ vagrant up cloud
```

An Ansible error about a `.cache/test_certs` directory being missing is benign.
If you see an error early on about VirtualBox guest additions failing to 
install, you have a few extra steps to take after the provisioning fails:

```console
HOST [~/fbsource/fbcode/magma/orc8r/cloud]$ vagrant reload cloud
HOST [~/fbsource/fbcode/magma/orc8r/cloud]$ vagrant provision cloud
```

This will re-run the Ansible provisioning. It'll be a lot quicker this time 
around since most the steps will be skipped.

#### Initial Build

Once all those jobs finish (should only be 5 minutes or so), we can build the
code for the first time and run it. Start in

##### Terminal Tab 1:

```console
HOST [magma]$ vagrant ssh magma
MAGMA-VM [/home/vagrant]$ cd magma/lte/gateway
MAGMA-VM [/home/vagrant/magma]$ make run
```
This will take a while (we have a lot of CXX files to build). It's a good time
to switch over to

##### Terminal Tab 3:

Here, we'll build and run all of our orchestrator cloud services. We'll also 
use this time to register the client certificate you'll need to access the API 
gateway for the controller running on your cloud VM.

```console
HOST [magma/orc8r/cloud]$ vagrant ssh cloud
CLOUD-VM [/home/vagrant]$ cd magma/orc8r/cloud
CLOUD-VM [/home/vagrant/magma/orc8r/cloud]$ make run
```

Again, this will take a while. While both these builds are churning, it's a 
good time to grab lunch, or a coffee. When the make in this tab is done:

```console
HOST$ open magma/.cache/test_certs
```

This will open up a finder window. Double-click the `admin_operator.pfx` cert
in this directory, which will open up Keychain to import the cert. The
password for the cert is `mai`. If you use Chrome or Safari, this is all you
need to do. If you use Firefox, copy this file to your desktop, then go to
`Preferences -> PrivacyAndSecurity -> View Certificates -> Import` and select
it.

Linux/Windows users should replace the above steps with the system-appropriate
method to import a client cert.

If you're reading this because you recently nuked your VMs, you have one more
thing to do. This `magma/.cache` folder stays around even when you destroy
your VMs so you don't have to fiddle with Keychain Access every time you spin
up a new box. But if you're starting with a blank datastore VM, the cert that
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

#### Poking Around the Systems

You can access the orchestrator REST API at https://192.168.80.10:9443/apidocs.
The SSL cert is self-signed, so click through any security warnings your
browser gives you. You should be prompted for a client cert, at which point
you should select the `admin_operator` cert that you added to Keychain above.

If you want to see what the access gateway is doing, you can
`vagrant ssh magma` inside `lte/gateway`, then do a
`sudo tail -f /var/log/syslog`. If everything above went smoothly, you should
eventually (give it a minute or two) see something along the lines of:

```console
MAGMA-VM$ sudo service magma@* stop
MAGMA-VM$ sudo service magma@magmad restart
MAGMA-VM$ sudo tail -f /var/log/syslog
Sep 27 22:57:35 magma-dev magmad[6226]: [2018-09-27 22:57:35,550 INFO root] Checkin Successful!
Sep 27 22:57:35 magma-dev control_proxy[6418]: 2018-09-27T22:57:35.550Z [127.0.0.1 -> logger-controller.magma.test,8443] "POST /magma.LoggingService/Log HTTP/2" 200 5bytes 0.031s
Sep 27 22:57:36 magma-dev magmad[6226]: [2018-09-27 22:57:36,521 INFO root] [SyncRPC] Sending heartbeat
Sep 27 22:57:36 magma-dev control_proxy[6418]: 2018-09-27T22:57:36.859Z [127.0.0.1 -> metricsd-controller.magma.test,8443] "POST /magma.MetricsController/Collect HTTP/2" 200 5bytes 0.061s
Sep 27 22:57:36 magma-dev control_proxy[6418]: 2018-09-27T22:57:36.859Z [127.0.0.1 -> streamer-controller.magma.test,8443] "POST /magma.Streamer/GetUpdates HTTP/2" 200 7bytes 0.043s
Sep 27 22:57:36 magma-dev policydb[6292]: [2018-09-27 22:57:36,903 INFO root] Processing 0 policy updates (resync=True)
Sep 27 22:57:37 magma-dev control_proxy[6418]: 2018-09-27T22:57:37.326Z [127.0.0.1 -> streamer-controller.magma.test,8443] "POST /magma.Streamer/GetUpdates HTTP/2" 200 7bytes 0.008s
Sep 27 22:57:37 magma-dev subscriberdb[6295]: [2018-09-27 22:57:37,327 INFO root] Processing 0 subscriber updates (resync=True)
Sep 27 22:57:38 magma-dev magmad[6226]: [2018-09-27 22:57:38,525 INFO root] [SyncRPC] Sending heartbeat
Sep 27 22:57:38 magma-dev magmad[6226]: [2018-09-27 22:57:38,606 INFO root] [SyncRPC] Got heartBeat from cloud
Sep 27 22:57:55 magma-dev control_proxy[6418]: 2018-09-27T22:57:55.683Z [127.0.0.1 -> streamer-controller.magma.test,8443] "POST /magma.Streamer/GetUpdates HTTP/2" 200 40bytes 0.037s
Sep 27 22:57:55 magma-dev magmad[6226]: [2018-09-27 22:57:55,684 INFO root] Processing config update g1
Sep 27 22:57:55 magma-dev control_proxy[6418]: 2018-09-27T22:57:55.683Z [127.0.0.1 -> streamer-controller.magma.test,8443] "POST /magma.Streamer/GetUpdates HTTP/2" 200 7bytes 0.009s
```

More in-depth information about the orchestrator cloud and LTE access gateway
can be found in the additional documentation below.

#### Additional Documentation

This README just scratches the surface of the system. A full deployment
involves the orchestrator cloud, the LTE access gateway, the NMS web frontend
for the orchestrator cloud, and the federated gateway for integrating with
a legacy 3GPP core.

For more in-depth documentation into the different parts of the system, see
the following docs:

- Running the Orchestrator Cloud locally:
    - See [Running Orchestrator Cloud](docs/running_cloud.md)
- Running the LTE gateway locally:
    - See [Running LTE Gateway](docs/running_gateway.md)
- Running the Federated Gateway locally:
    - See [Running FeG](docs/running_feg.md)
    
## Join the Magma Community

- Mailing lists:
  - Join <magma-dev@googlegroups.com> for technical discussions
  - Join <magma-announce@googlegroups.com> for announcements
- Discord:
  - [magma\_dev](https://discord.gg/WDBpebF) channel

## License

Magma is BSD License licensed, as found in the LICENSE file.
