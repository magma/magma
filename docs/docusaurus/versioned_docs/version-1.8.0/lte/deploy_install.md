---
id: version-1.8.0-deploy_install
title: Install AGW
hide_title: true
original_id: deploy_install
---

# Install Access Gateway on Ubuntu (Bare Metal)

## Prerequisites

To set up a Magma Access Gateway, you will need a machine that
satisfies the following requirements:

- AGW_HOST: 64bit-X86 machine, baremetal strongly recommended
  (not virtualized). You will need two ethernet ports. In this guide,
  `enp1s0` and `enp2s0` are used:
    - `enp1s0`: Will carry any traffic that is not S1. That is, data plane traffic (SGi),
    control plane traffic (Orc8r HTTP2) and management (ssh).
    - `enp2s0`: S1 interface.

> NOTE:
>
> - Interface names might have different names on your hardware, so just
>   replace `enp1s0` and `enp2s0` with your current interface names
>   when following this guide.
>
> - The `agw_install_ubuntu.sh` script will rename the `enp1s0`
>   interface to `eth0`.
>
> - If you do not want all internet traffic to go through `enp1s0`
>  to separate control plane (Orc8r Http2 traffic) from user plane, you
>  may want to add another interface and configure proper routing.

## Deployment

### 1. Create boot USB stick and install Ubuntu on your AGW host

- Download the Ubuntu Server 20.04 LTS `.iso` image from the Ubuntu website. Verify its integrity by checking [the hash](https://releases.ubuntu.com/20.04/SHA256SUMS).
- Create a bootable USB using this [Etcher tutorial](https://tutorials.ubuntu.com/tutorial/tutorial-create-a-usb-stick-on-macos#0)
- Boot your AGW host from USB
    - Press F11 to select boot sequence. WARNING: This might be different for your machine.
    - If you see two options to boot from USB, select the non-UEFI option.
- Install and configure your Access Gateway according to your network defaults.
    - Make sure to enable ssh server and utilities (untick every other).
- Connect your SGi interface to the internet and select this port during the
installation process to get an IP using DHCP. (Consider enabling DHCP snooping to mitigate [Rogue DHCP](https://en.wikipedia.org/wiki/Rogue_DHCP))

### 2. Deploy magma on the AGW_HOST

#### Run AGW installation

To install on a server with a DHCP-configured SGi interface

```bash
su
wget https://raw.githubusercontent.com/magma/magma/v1.8/lte/gateway/deploy/agw_install_ubuntu.sh
bash agw_install_ubuntu.sh
```

To install on a server with statically-allocated SGi interface. For example,
if SGi has an IP of 1.1.1.1/24 and the upstream router IP is 1.1.1.200

```bash
su
wget https://raw.githubusercontent.com/magma/magma/v1.8/lte/gateway/deploy/agw_install_ubuntu.sh
bash agw_install_ubuntu.sh 1.1.1.1/24 1.1.1.200
```

The script will run a pre-check script that will tell you what will change on your machine
and prompt you for your approval. If you are okay with these changes, reply `yes` and magma will
be installed. If you respond with `no`, the installation will be stopped.

```bash
  - Check if Ubuntu is installed
  Ubuntu is installed
  - Check for magma user
  magma user is not Installed
  - Check if both interfaces are named eth0 and eth1
  Interfaces will be renamed to eth0 and eth1
  eth0 will be set to dhcp and eth1 10.0.2.1
  Do you accept those modifications and want to proceed with magma installation?(y/n)
  Please answer yes or no.
```

The machine will reboot but the installation is not finished yet; the script is still running in the background.
You can follow the output using

```bash
journalctl -fu agw_installation
```

When you see "AGW installation is done.", it means that your installation is complete. You can make sure magma is running by executing

```bash
service 'magma@*' status
```

#### Post-install check

Make sure you have the `control_proxy.yml` file in directory `/var/opt/magma/configs/`
before running the post-install script

```bash
bash /root/agw_post_install_ubuntu.sh
```
