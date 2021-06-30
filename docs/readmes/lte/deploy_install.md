---
id: deploy_install
title: Install AGW on Debian Stretch.
hide_title: true
---

For Installing AGW on Ubuntu refer to : [_deploy_on_ubuntu_](https://magma.github.io/magma/docs/lte/deploy_install_ubuntu)
# Install Access Gateway on Debian Stretch (Bare Metal)

### AGW deployment on Debian distribution is DEPRECATED! Release 1.5 is last release supporting debian based AGW.

## Prerequisites

To setup a Magma Access Gateway, you will need a machine that
satisfies the following requirements:

- AGW_HOST: 64bit-X86 machine, hardware strongly recommended (not virtualized).
  You will need two ethernet ports. We use enp1s0 and enp2s0 in this guide.
  They might have different names on your hardware so just replace enp1s0 and
  enp2s0 with your current interfaces name in this guideline.
  One port is for the SGi interface (default: enp1s0) and one for the S1
  interface (default: enp2s0). Note that the `agw_install.sh` script will
  rename the `enp1s0` interface to `eth0`.

## Deployment

### 1. Create boot USB stick and install Debian on your AGW host

- Download .iso image from [Debian mirror](https://cdimage.debian.org/cdimage/archive/9.13.0/amd64/iso-cd/debian-9.13.0-amd64-netinst.iso)
- Create bootable usb using etcher [tutorial here](https://tutorials.ubuntu.com/tutorial/tutorial-create-a-usb-stick-on-macos#0)
- Boot your AGW host from USB
  (Press F11 to select boot sequence, :warning: This might be different for
  your machine). If you see 2 options to boot from USB, select the non-UEFI
  option.
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
- Connect your SGi interface to the internet and select this port during the
installation process to get an IP using DHCP.

### 2. Deploy magma on the  AGW_HOST

*Run AGW installation*

```bash
su
wget https://raw.githubusercontent.com/facebookincubator/magma/v1.5.2/lte/gateway/deploy/agw_install.sh
bash agw_install.sh
```

The script will run a pre-check script that will prompt you what will change
on your machine. If you're okay with those changes reply `yes` and magma will
be installed. If `no` is replied It will stop the installation.

```bash
  - Check if Debian is installed
  Debian is installed
  - Check for correct Linux Headers
  New Linux Headers will be Installed
  - Check for magma user
  magma user is not Installed
  - Check if both interfaces are named eth0 and eth1
  Interfaces will be renamed to eth0 and eth1
  eth0 will be set to dhcp and eth1 10.0.2.1
  Do you accept those modifications and want to proceed with magma installation?(y/n)
  Please answer yes or no.
  ```

The machine will reboot but It's not finished yet, the script is still running in the background.
You can follow the output there

```bash
journalctl -fu agw_installation
```

When you see "AGW installation is done." It means that your AGW installation is done, you can make sure magma is running by executing:

```bash
service magma@* status
```

*Post Install Check*

``` bash
cd ~/magma/lte/gateway/deploy
./agw_post_install.sh
```
