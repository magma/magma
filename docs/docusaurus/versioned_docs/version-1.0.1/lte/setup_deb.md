---
id: version-1.0.1-setup_deb
title: AGW Setup (Bare Metal)
sidebar_label: Setup (Bare Metal)
hide_title: true
original_id: setup_deb
---
# Access Gateway Setup (On Bare Metal)
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

- Download .iso image from [Debian mirror](http://cdimage.debian.org/mirror/cdimage/archive/9.9.0/amd64/iso-cd/debian-9.9.0-amd64-netinst.iso)
- Create bootable usb using etcher [tutorial here](https://tutorials.ubuntu.com/tutorial/tutorial-create-a-usb-stick-on-macos#0)
- Connect AGW host to a working internet connection and select this port during the installation process to get an IP using DHCP
- Boot your AGW host from USB
  (Press F11 to select boot sequence, :warning: This might be different for
  your machine). If you see 2 options to boot from USB, select the non-UEFI
  option.
- Select “Install” option
- Select a language : "English"
- Select your location 
- Configure the keyboard : "American English"
- Primary network interface "enp1s0"
- Configuration
  * Hostname: “magma”
  * Domain name : “”
  * Root password: “magma”
  * Username : “magma”
  * Password: “magma”
  * Partition disk: "Use entire disk"
  * Select disk to partition: "sda"
  * Partitioning scheme: "All files in one partition"
  * Configure the package manager and HTTP proxy information: “”
  * Software selection : Only tick ssh server and utilities (untick every other)
  * Install the GRUB boot loader on a hard disk : "Yes" and select "/dev/sda"
  * Finish the installation : "Continue" 

### 2. Deploy magma on the  AGW_HOST

```bash
su
wget https://raw.githubusercontent.com/magma/magma/v1.0/lte/gateway/deploy/agw_install.sh
bash agw_install.sh
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
