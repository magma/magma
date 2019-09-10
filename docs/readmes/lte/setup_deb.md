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


- Docker host where the container AGW_DEPLOY will be built
- AGW_HOST, 64bit-X86 machine
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

### 2. Prepare AGW_DEPLOY
- [AGW_DEPLOY] Build and run AGW_DEPLOY container

```bash
git clone https://github.com/facebookincubator/magma.git ~/magma  
cd ~/magma
docker build --build-arg CACHE_DATE="$(date)" -t agw_deploy -f lte/gateway/docker/deploy/Dockerfile .
docker run -it agw_deploy bash
```

### 3. Prepare AGW_HOST
- [AGW_HOST] Change interface names

```bash
su

sed -i 's/GRUB_CMDLINE_LINUX=""/GRUB_CMDLINE_LINUX="net.ifnames=0 biosdevname=0"/g' /etc/default/grub

grub-mkconfig -o /boot/grub/grub.cfg
sed -i 's/enp1s0/eth0/g' /etc/network/interfaces

echo "auto eth1
iface eth1 inet static
address 10.0.2.1
netmask 255.255.255.0" > /etc/network/interfaces.d/eth1

reboot
```

- [AGW_HOST] Make magma user a sudoer on AGW host

```bash
su
apt install -y sudo python-minimal aptitude
adduser magma sudo
echo "magma ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers
```

- [AGW_DEPLOY] Copy the public key

``` bash
# copy output
cat ~/.ssh/id_rsa.pub
```

- [AGW_HOST] paste the pub key to ~/.ssh/authorized_keys

```bash
su magma
mkdir -p ~/.ssh
# paste pub key inside
vi ~/.ssh/authorized_keys     
```

- [AGW_DEPLOY] Prepare Hostfile agw_hosts for deployment

``` bash
AGW_IP=10.0.2.1
USERNAME=magma
echo "[ovs_build]
$AGW_IP ansible_ssh_port=22 ansible_ssh_user=$USERNAME ansible_ssh_private_key_file=~/.ssh/id_rsa
[ovs_deploy]
$AGW_IP ansible_ssh_port=22 ansible_ssh_user=$USERNAME ansible_ssh_private_key_file=~/.ssh/id_rsa
" > ~/magma/lte/gateway/deploy/agw_hosts
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
