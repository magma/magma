---
id: version-1.8.0-deploy_install_docker
title: Install Docker AGW
hide_title: true
original_id: deploy_install_docker
---

# Install Docker-based Access Gateway on Ubuntu

## Prerequisites

To set up a Magma Access Gateway, you will need a machine that
satisfies the following requirements:

- AGW_HOST: aarch64 or 64bit-X86 machine. You will need two ethernet ports. In this guide,
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
> - The `agw_install_docker.sh` script will rename the `enp1s0`
>   interface to `eth0`.
>
> - If you do not want all internet traffic to go through `enp1s0`
>  to separate control plane (Orc8r Http2 traffic) from user plane, you
>  may want to add another interface and configure proper routing.

## Deployment

### 1. Create boot USB stick and install Ubuntu on your AGW host

- Download the Ubuntu Server 20.04 LTS `.iso` image from the Ubuntu website
- Create a bootable USB using this [Etcher tutorial](https://tutorials.ubuntu.com/tutorial/tutorial-create-a-usb-stick-on-macos#0)
- Boot your AGW host from USB
    - Press F11 to select boot sequence. WARNING: This might be different for your machine.
    - If you see two options to boot from USB, select the non-UEFI option.
- Install and configure your Access Gateway according to your network defaults.
    - Make sure to enable ssh server and utilities (untick every other).
- Connect your SGi interface to the internet and select this port during the
installation process to get an IP using DHCP.

### 2. Deploy magma on the AGW_HOST

#### Do pre-installation steps

Become root user:

```bash
sudo -i
```

Copy your `rootCA.pem` file from orc8r to the following location:

```bash
mkdir -p /var/opt/magma/certs
vim /var/opt/magma/certs/rootCA.pem
```

#### Run AGW installation

Download AGW docker install script

```bash
wget https://github.com/magma/magma/raw/v1.8/lte/gateway/deploy/agw_install_docker.sh
bash agw_install_docker.sh
```

#### Configure AGW

Once you see the output `Reboot this machine to apply kernel settings`, reboot your AGW host.

Create `control_proxy.yml` file with your orc8r details:

```bash
cat << EOF | sudo tee /var/opt/magma/configs/control_proxy.yml
cloud_address: controller.orc8r.magmacore.link
cloud_port: 443
bootstrap_address: bootstrapper-controller.orc8r.magmacore.link
bootstrap_port: 443
fluentd_address: fluentd.orc8r.magmacore.link
fluentd_port: 24224

rootca_cert: /var/opt/magma/certs/rootCA.pem
EOF
```

Start your access gateway:

```bash
cd /var/opt/magma/docker
sudo docker-compose up -d
```

Now get Hardware ID and Challenge key and add AGW in your orc8r:

```bash
docker exec magmad show_gateway_info.py
```

Then restart your access gateway:

```bash
sudo docker-compose up -d --force-recreate
```
