---
id: version-1.0.0-setup_baremetal
title: AGW Setup (Bare Metal)
sidebar_label: Setup (Bare Metal)
hide_title: true
original_id: setup_baremetal
---
# AGW installation on baremetal


### HW requirements
- Ansible host
- AGW host, 64bit-X86 machine
  Two ethernet ports. We use enp1s0 and enp2s0 in this guide. They might have different names on your hardware so just replace enp1s0 and enp2s0 by your current interfaces name in this guidelind.
  One port is for the SGi interface (default: enp1s0) and one for the S1 interface (default: enp2s0)

### Install Ansible (on Ansible host)
- Mac Installation  
```bash
brew install Ansible
```

- Install ansible on a RHEL/CentOS Linux based system
```bash
sudo yum install ansible
```

- Install ansible on a Debian/Ubuntu Linux based system
```bash
sudo apt-get install software-properties-common
sudo apt-add-repository ppa:ansible/ansible
sudo apt-get update
sudo apt-get install ansible
```

- Install ansible using pip
```bash
pip install ansible
```

- Using source code
``` bash
cd ~
git clone git://github.com/ansible/ansible.git
cd ./ansible
source ./hacking/env-setup
```
When running ansible from a git checkout, one thing to remember is that you will need to setup your environment everytime you want to use it, or you can add it to your bash rc file:
```bash
echo "export ANSIBLE_HOSTS=~/ansible_hosts" >> ~/.bashrc
echo "source ~/ansible/hacking/env-setup" >> ~/.bashrc
```

### Install Debian Stretch (on AGW host)

1. Create boot USB stick

  - Download .iso image from [Debian mirror](http://cdimage.debian.org/mirror/cdimage/archive/9.9.0/amd64/iso-cd/debian-9.9.0-amd64-netinst.iso)
  - Create bootable usb using etcher [tutorial here](https://tutorials.ubuntu.com/tutorial/tutorial-create-a-usb-stick-on-macos#0)
  - Boot your AGW host from USB (Press F11 to select boot sequence, :warning: This might be different for your machine)
  - Select “Install” option.
  - Configuration.
    * Hostname: “magma”
    * Domain name : “”
    * Root password: “magma”
    * Username : “magma”
    * Password: “magma”
    * Partition disk: Use entire disk and put all files in one partition
    * Only install ssh server and utilities.
  - Connect you SGi interface to the internet and select this port during the installation process to get an IP using DHCP.


2. Prepare AGW
  - Change interface names

  ```bash
  su

  sed -i 's/GRUB_CMDLINE_LINUX=""/GRUB_CMDLINE_LINUX="net.ifnames=0 biosdevname=0"/g' /etc/default/grub

  grub-mkconfig -o /boot/grub/grub.cfg
  sed -i 's/enp1s0/eth0/g' /etc/network/interfaces

  echo “auto eth1
  iface eth1 inet static
  address 10.0.2.1
  netmask 255.255.255.0” > /etc/network/interfaces.d/eth1

  reboot
  ```

  - Make magma user a sudoer on AGW host

  ```bash
  su
  apt install -y sudo python-minimal
  adduser magma sudo
  echo "magma ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers
  ```

  - Generate key on your Ansible host

  ``` bash
  cd .ssh
  ssh-keygen
  chmod 600 id_rsa
  ssh-add id_rsa
  cat id_rsa.pub # copy text
  ```

  - On AGW host, add the pub key to authorized_keys

  ```bash
  su magma
  mkdir -p ~/.ssh
  vi ~/.ssh/authorized_keys     # paste pub key inside
  ```

  - Install Magma

  ``` bash
  #Clone magma repository on Ansible host
  git clone https://github.com/facebookincubator/magma.git ~/
  ```

  - Prepare for deployment

  ``` bash
  AGW_IP=10.0.2.1
  USERNAME=magma
  echo "[ovs_build]
  $AGW_IP ansible_ssh_port=22 ansible_ssh_user=$USERNAME ansible_ssh_private_key_file=~/.ssh/id_rsa
  [ovs_deploy]
  $AGW_IP ansible_ssh_port=22 ansible_ssh_user=$USERNAME ansible_ssh_private_key_file=~/.ssh/id_rsa
  " >> ~/magma/lte/gateway/deploy/agw_hosts
  ```

 - Build ovs

  ``` bash
  cd ~/magma/lte/gateway/deploy
  ansible-playbook -e "MAGMA_ROOT='~/magma' OUTPUT_DIR='/tmp'" -i agw_hosts ovs_gtp.yml
  ```

  - Deploy AGW

  ``` bash
  cd ~/magma/lte/gateway/deploy
  ansible-playbook -i agw_hosts -e "PACKAGE_LOCATION='/tmp'" ovs_deploy.yml
  ```
