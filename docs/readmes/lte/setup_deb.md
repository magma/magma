---
id: setup_deb
title: AGW Setup (Bare Metal)
sidebar_label: Setup (Bare Metal)
hide_title: true
---
# Access Gateway Setup (On Bare Metal)
## Prerequisites

To setup a Magma Access Gateway, you will need a machine that
satisfies the following requirements:


- Docker host where the container AGW_DEPLOY will be built. This container can be run directly on your machine or a remote host, as long as It can reach your Gateway.
- AGW_HOST, 64bit-X86 machine, we strongly recommend hardware.
  Two ethernet ports. We use enp1s0 and enp2s0 in this guide. They might have different names on your hardware so just replace enp1s0 and enp2s0 with your current interfaces name in this guideline.
  One port is for the SGi interface (default: enp1s0) and one for the S1 interface (default: enp2s0)

## Deployment
### 1. Create boot USB stick and install Debian on your AGW host

- Download .iso image from [Debian mirror](http://cdimage.debian.org/mirror/cdimage/archive/9.9.0/amd64/iso-cd/debian-9.9.0-amd64-netinst.iso)
- Create bootable usb using etcher [tutorial here](https://tutorials.ubuntu.com/tutorial/tutorial-create-a-usb-stick-on-macos#0)
- Boot your AGW host from USB (Press F11 to select boot sequence, :warning: This might be different for your machine)
- Select “Install” option.
- Network missing firmeware "No"
- Primary network interface "enp1s0"
- Configuration.
  * Hostname: “magma”
  * Domain name : “”
  * Root password: “magma”
  * Username : “magma”
  * Password: “magma”
  * Partition disk: "Use entire disk"
  * Select disk to partition: "sda"
  * Partitioning scheme: "All files in one partition"
  * Only tick ssh server and utilities (untick every other)
- Connect your SGi interface to the internet and select this port during the installation process to get an IP using DHCP.

### 2. Prepare AGW_HOST
- [AGW_HOST] Prepare AGW_HOST

```bash
su
wget https://raw.githubusercontent.com/facebookincubator/magma/master/lte/gateway/deploy/agw_prepare.sh
sh agw_prepare.sh
```

A prompt will pop up to as you if you want to stop removing linux-image-4.9.0-11-amd64 please hit: No

### 3. Prepare AGW_DEPLOY
- [AGW_DEPLOY] Build and run AGW_DEPLOY container

```bash
git clone https://github.com/facebookincubator/magma.git ~/magma  
cd ~/magma
docker build --build-arg CACHE_DATE="$(date)" -t agw_deploy -f lte/gateway/docker/deploy/Dockerfile .
docker run -it agw_deploy bash
scp ~/.ssh/id_rsa.pub magma@10.0.2.1:/home/magma/.ssh/authorized_keys
```

### 4. Build openvswitch on Access Gateway
- [AGW_DEPLOY] Run build playbook

``` bash
cd ~/magma/lte/gateway/deploy
ansible-playbook -e "MAGMA_ROOT='~/magma' OUTPUT_DIR='/tmp'" -i agw_hosts ovs_build.yml
```

### 5. Deploy Access Gateway
- [AGW_DEPLOY] Run deploy playbook

``` bash
cd ~/magma/lte/gateway/deploy
ansible-playbook -i agw_hosts -e "PACKAGE_LOCATION='/tmp'" ovs_deploy.yml
```
