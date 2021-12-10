---
id: version-1.4.0-quick_start_guide
title: Quick Start Guide
hide_title: true
original_id: quick_start_guide
---
# Quick Start Guide

The quick start guide is for developing on Magma or just trying it out. Follow
the deployment guides under Orchestrator and Access Gateway if you are
installing Magma for a production deployment.

With the [prereqs](prerequisites.md) installed, we can now set up a minimal
end-to-end system on your development environment. In this guide, we'll start
by running the LTE access gateway and orchestrator cloud, and then
register your local access gateway with your local cloud for management.

We will be spinning up a virtual machine and some docker containers for this
full setup, so you'll probably want to do this on a system with at least 8GB
of memory. Our development VM's are in the 192.168.60.0/24, 192.168.128.0/24 and
192.168.129.0/24 address spaces, so make sure that you don't have anything
running which hijacks those (e.g. VPN).

In the following steps, note the prefix in terminal commands. `HOST` means to
run the indicated command on your host machine, and `MAGMA-VM` on the `magma`
vagrant machine under `lte/gateway`.

## Provisioning the environment

Go ahead and open up 2 fresh terminal tabs. Start in

### Terminal Tab 1: Provision the AGW VM

The development environment virtualizes the access gateway so you don't need
any production hardware on hand to test an end-to-end setup.
We'll be setting up the LTE AGW VM in this tab.

```bash
HOST [magma]$ cd lte/gateway
HOST [magma/lte/gateway]$ vagrant up magma
```

This will take a few minutes to spin up the VM. While that runs, switch over
to...

**Note**: If you are looking to test/develop the LTE features of AGW, without
cloud based network management, you can skip the rest of this guide and try the
[S1AP integration tests](../lte/s1ap_tests.md) now.

### Terminal Tab 2: Build Orchestrator

Here, we'll be building the Orchestrator docker containers.

```bash
HOST [magma]$ cd orc8r/cloud/docker
HOST [magma/orc8r/cloud/docker]$ ./build.py --all
```

This will build all the docker images for Orchestrator. The `vagrant up` from
the first tab should finish before the image building, so you should switch
to that tab and move on for now.

## Initial Run

Once `vagrant up` in the first tab finishes:

### Terminal Tab 1: Build AGW from Source

We will kick off the initial build of the AGW from source here.

```bash
HOST [magma/lte/gateway]$ vagrant ssh magma
MAGMA-VM [/home/vagrant]$ cd magma/lte/gateway
MAGMA-VM [/home/vagrant/magma/lte/gateway]$ make run
```

This will take a while (we have a lot of CXX files to build). With 2 extensive
build jobs running, now is a good time to grab a coffee or lunch. The first
build ever from source will take a while, but afterwards, a persistent ccache
and Docker's native layer caching will speed up subsequent builds
significantly.

You can monitor what happens in the other tab now:

### Terminal Tab 2: Start Orchestrator

Once the Orchestrator build finishes, we can start the development Orchestrator
cloud for the first time. We'll also use this time to register the local
client certificate you'll need to access the local API gateway for your
development stack.

To start Orchestrator (without metrics) is as simple as:

```bash
HOST [magma/orc8r/cloud/docker]$ ./run.py

Creating orc8r_postgres_1 ... done
Creating orc8r_test_1     ... done
Creating orc8r_maria_1    ... done
Creating elasticsearch    ... done
Creating fluentd          ... done
Creating orc8r_kibana_1   ... done
Creating orc8r_proxy_1      ... done
Creating orc8r_controller_1 ... done
```

If you want to run everything, including metrics, run:

```bash
HOST [magma/orc8r/cloud/docker]$ ./run.py --metrics

Creating orc8r_alertmanager_1     ... done
Creating orc8r_maria_1            ... done
Creating elasticsearch            ... done
Creating orc8r_postgres_1         ... done
Creating orc8r_config-manager_1   ... done
Creating orc8r_test_1             ... done
Creating orc8r_prometheus-cache_1 ... done
Creating orc8r_prometheus_1       ... done
Creating orc8r_kibana_1           ... done
Creating fluentd                  ... done
Creating orc8r_proxy_1            ... done
Creating orc8r_controller_1       ... done
```

The Orchestrator application containers will bootstrap certificates on startup
which are cached for future runs. Watch the directory `magma/.cache/test_certs`
for a file `admin_operator.pfx` to show up (this may take a minute or 2), then:

```bash
HOST [magma/orc8r/cloud/docker]$ ls ../../../.cache/test_certs

admin_operator.key.pem  bootstrapper.key        controller.crt          rootCA.key
admin_operator.pem      certifier.key           controller.csr          rootCA.pem
admin_operator.pfx      certifier.pem           controller.key          rootCA.srl

HOST [magma/orc8r/cloud/docker]$ open ../../../.cache/test_certs
```

In the Finder window that pops up, double-click `admin_operator.pfx` to add the
local client cert to your keychain. *The password for the cert is magma*.
In some cases, you may have to open up the Keychain app in MacOS and drag-drop
the file into the login keychain if double-clicking doesn't work.

If you use Firefox, you'll have to import this .pfx file into your browser's
installed client certificates. See [here](https://support.globalsign.com/customer/en/portal/articles/1211486-install-client-digital-certificate---firefox-for-windows)
for instructions. If you use Chrome or Safari, you may have to restart the
browser before the certificate can be used.

### Connecting Your Local LTE Gateway to Your Local Cloud

At this point, you will have built all the code in the LTE access gateway and
the Orchestrator cloud. All the services on the LTE access gateway and
orchestrator cloud are running, but your gateway VM isn't yet set up to
communicate with your local cloud.

We have a fabric command set up to do this:

```bash
HOST [magma]$ cd lte/gateway
HOST [magma/lte/gateway]$ fab -f dev_tools.py register_vm
```

This command will seed your gateway and network on Orchestrator with some
default LTE configuration values and set your gateway VM up to talk to your
local Orchestrator cloud. Wait a minute or 2 for the changes to propagate,
then you can verify that things are working:

```bash
HOST [magma/lte/gateway]$ vagrant ssh magma

MAGMA-VM$ sudo service magma@* stop
MAGMA-VM$ sudo service magma@magmad restart
MAGMA-VM$ sudo tail -f /var/log/syslog

# After a minute or 2 you should see these messages:
Sep 27 22:57:35 magma-dev magmad[6226]: [2018-09-27 22:57:35,550 INFO root] Checkin Successful!
Sep 27 22:57:55 magma-dev magmad[6226]: [2018-09-27 22:57:55,684 INFO root] Processing config update g1
Sep 27 22:57:55 magma-dev control_proxy[6418]: 2018-09-27T22:57:55.683Z [127.0.0.1 -> streamer-controller.magma.test,8443] "POST /magma.Streamer/GetUpdates HTTP/2" 200 7bytes 0.009s
```

## Using the NMS UI

Magma provides an UI for configuring and monitoring the networks. To set up
the NMS to talk to your local Orchestrator:

```bash
HOST [magma]$ cd nms/app/packages/magmalte
HOST [magma/nms/app/packages/magmalte] $ docker-compose build magmalte
HOST [magma/nms/app/packages/magmalte] $ docker-compose up -d
HOST [magma/nms/app/packages/magmalte] $ ./scripts/dev_setup.sh
```

After this, you will be able to access the UI by visiting
[https://magma-test.localhost](https://magma-test.localhost), and using the email `admin@magma.test`
and password `password1234`. We recommend Firefox or Chrome. If you see Gateway Error 502, don't worry, the
NMS can take upto 60 seconds to finish starting up.

You will probably want to enable this organization (magma-test) to access all networks,
so go to [master.localhost](https://master.localhost) and login with the same credentials.
Once there, you can click on the organization and then select "Enable all networks".

**Note**: If you want to test the access gateway VM with a physical eNB and UE,
refer to
the [Connecting a physical eNodeb and UE device to gateway
VM](../lte/dev_notes.md#connecting-a-physical-enodeb-and-ue-to-gateway-vm)
section.
