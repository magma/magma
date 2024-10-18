---
id: version-1.7.0-deploy_install
title: Install AGW
hide_title: true
original_id: deploy_install
---

# Install Access Gateway on Ubuntu (Bare Metal)

> NOTE: Debian-based AGW deployments are no longer supported as of v1.6. If you want to install to Debian, refer to [v1.5 of the documentation](https://magma.github.io/magma/versions).

## Prerequisites

To set up a Magma Access Gateway, you will need a machine that
satisfies the following requirements:

- AGW_HOST: 64bit-X86 machine, baremetal strongly recommended
  (not virtualized). You will need two ethernet ports. We use
  `enp1s0` and `enp2s0` in this guide.
    - `enp1s0`: Will carry any traffic that is not S1. So data plane traffic (SGi)
    control plane traffic (Orc8r HTTP2), management (ssh)
    - `enp2s0`: S1 interface.

  *Note interface names might have different names on your hardware so just
  replace `enp1s0` and `enp2s0` with your current interfaces name
  in this guideline.

  *Note that the `agw_install_ubuntu.sh` script will rename the `enp1s0`
   interface to `eth0`.

  *Note if you don't want all internet traffic to go through `enp1s0`
  to separate control plane (Orc8r Http2 traffic) from user plane, you
  may want to add another interface and configure proper routing.

## Deployment

### 1. Create boot USB stick and install Ubuntu on your AGW host

- Download the Ubuntu Server 20.04 LTS .iso image from the Ubuntu website
- Create bootable usb using etcher [tutorial here](https://tutorials.ubuntu.com/tutorial/tutorial-create-a-usb-stick-on-macos#0)
- Boot your AGW host from USB
  (Press F11 to select boot sequence, :warning: This might be different for
  your machine). If you see 2 options to boot from USB, select the non-UEFI
  option.
- Install and configure you access gateway according to your network defaults.
    - Make sure to enable ssh server and utilities (untick every other)
- Connect your SGi interface to the internet and select this port during the
installation process to get an IP using DHCP.

### 2. Deploy magma on the  AGW_HOST

#### Run AGW installation

To install on server with DHCP configured SGi interface.

```bash
su
wget https://raw.githubusercontent.com/magma/magma/v1.6/lte/gateway/deploy/agw_install_ubuntu.sh
bash agw_install_ubuntu.sh
```

To Install on server with statically allocated SGi interface. Fow example:
SGi has 1.1.1.1/24 IP and upstream router IP is 1.1.1.200

```bash
su
wget https://raw.githubusercontent.com/magma/magma/v1.6/lte/gateway/deploy/agw_install_ubuntu.sh
bash agw_install_ubuntu.sh 1.1.1.1/24 1.1.1.200
```

The script will run a pre-check script that will prompt you what will change
on your machine. If you're okay with those changes reply `yes` and magma will
be installed. If `no` is replied It will stop the installation.

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

The machine will reboot but It's not finished yet, the script is still running in the background.
You can follow the output there

```bash
journalctl -fu agw_installation
```

When you see "AGW installation is done." It means that your AGW installation is done, you can make sure magma is running by executing:

```bash
service magma@* status
```

#### Post Install Check

Make sure you have `control_proxy.yml` file in directory /var/opt/magma/configs/
before running post install script.

```bash
bash /root/agw_post_install_ubuntu.sh
```
