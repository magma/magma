---
id: deploy_ansible_install
title: Install Orchestrator on Bare Metal with Ansible
hide_title: true
---

# Install Orchestrator with Ansible

This page walks through a full, vanilla Orchestrator install using Ansible on bare metal or a virtual machine.

If you want to install a specific release version, see the notes in the
[deployment intro](https://docs.magmacore.org/docs/orc8r/deploy_intro).

## Advanced users

If you have an existing Kubernetes cluster or want to make changes to the
deployment, refer to the [advanced deployment notes](docs/advanced_notes.md).

## Prerequisites

We assume `MAGMA_ROOT` is set as described in the
[deployment intro](https://docs.magmacore.org/docs/orc8r/deploy_intro).

This walkthrough assumes you already have the following

- a registered domain name
- a physical or virtual server with at least 4 vCPU, 16 RAM, and 100Gb free storage

## Preparing for the install

The Ansible install of the Orchestrator is designed to automate as much of the deployment process as possible, but you'll need to do a little bit of preparation ahead of time.

### Preparing the deployment server

First let's prepare the server itself. We're assuming that you are running a stock install of Ubuntu 18.04 or 20.04, so you'll want to begin by installing some basic tools that we'll need.

1. Enable password-less sudo on the server. To do this, use 

   ```bash
   sudo visudo /etc/sudoers
   ```

   to edit the sudoers file. You want to add permission for passwordless sudo for your user to the end of the file. For example, we'd add this permission  for the ubuntu user using:

   ```bash
   ubuntu ALL=(ALL) NOPASSWD: ALL
   ```

   Obviously you'll want to substitute your own user for `ubuntu`. 

   To save and exit, press `<ctrl>-x` then `y` and `<enter>`.
1. Set up passwordless login for the deployment user.  If you don't have a private key, create that first.
   ```bash
   ssh-keygen -t rsa -b 4096 -C "your_email@domain.com"
   ```
   Then copy the key to `authorized_keys` for this server:
   ```bash
   ssh-copy-id yourusername@localhost
   ```
1. Clone the Orchestrator repository:
   ```bash
   git clone https://github.com/magma/magma.git
   ```
1. Set the `MAGMA_ROOT` variable:
   ```bash
   export MAGMA_ROOT=PATH_TO_YOUR_MAGMA_CLONE
   ```
1. If you're running a VirtualBox VM, you need to set the network as Bridged, and edit netplan. Open the appropriate file with
   ```bash
   sudo vi /etc/netplan/01-netcfg.yaml
   ```
   Then add the following content:
   ```bash
   network:
     version: 2
     renderer: networkd
     ethernets:
       enp0s4:
         dhcp4: yes
   ```
   and apply the changes with
   ```bash
   sudo netplan apply
   ```
Now we're ready to configure the deployment itself.
## Configure Orchestrator

While the deployment script takes care of all the work in doing the actual deployment, you do need to provide it with the various pieces of information it will need to do its job. You'll specify these values in the
```bash
$MAGMA_ROOT/orc8r/cloud/deploy/bare-metal-ansible/ansible_vars.yaml
```
file.

Values in this file you will need to check and/or set to make sure they match your own values include:

* `ansible_user`: The Linux user that will be running Ansible. This is the user your set up for passwordless sudo and passwordless ssh earlier.
* `orc8r_image_repo`: The repository that hosts the docker containers that make up your Orchestrator. 
* `orc8r_helm_repo`: The repository that hosts the helm charts that explain how to deploy your Orchestrator. This can be a public repo, or it can be a Github repo protected by a username and password.
* `orc8r_domain`: Your domain name. Can be `username.local` for test installations.
* `orc8r_nms_admin_email`: The email address for the Orchestrator administrator.
* `docker_insecure_registries`: Ansible is assuming your Docker registry has trusted SSL or https: support, so if that's not the case, add your image repository here. Should match your value for orc8r_image_repo
* `orc8r_chart_version`: Helm chart tag in the repo. In this exmple, we used v1.4.37.
* `orc8r_image_tag`: Orchestrator Docker image tag in the repo. In this example, we used v1.4.37.
* `orc8r_nms_image_tag`: ??? image chart tag in the repo. In this example, we used v1.3.7.
* `orc8r_nginx_image_tag`: Image tag for the Nginx container image in the repo. In this example, we used version v1.3.7.
* `metallb_addresses`: LoadBalancer settings for Magma services publicly exposed. These must be at least 4 sequential IP addresses on the same network as the ansible host. You should ping each of them to ensure there are no other hosts using them. 

There are a number of values that can usually be left with their default values:
```bash
# Change these to your private nameservers if needed
nameservers:
  - 8.8.8.8
  - 8.8.4.4
upstream_dns_servers:
  - 8.8.8.8
  - 8.8.4.4

## Kubernetes LoadBalancer
## Enabling a loadbalancer is optional, but provides HA for external access to Kubernetes API
## loadbalancer_apiserver and vrrp_nic need to be set if you want to enable this feature.

# Change to match the network of nodes where you want to run loadbalancer
# loadbalancer_apiserver:
#   address: 192.168.0.20
#   port: 8383
# Set to the NIC on which you want to run loadbalancer
# vrrp_nic: "eth0"

# Component passwords are randomly generated (and stored in credentials dir).
# Uncomment the following lines to explicitly set passwords
# db_root_password:
# orc8r_db_pass:
# orc8r_nms_db_pass:
# orc8r_nms_admin_pass:
```

Once you've configured the Ansible vars file, it's time to run the deployment.

## Deploy the Orchestrator

We're almost ready to deploy the orchestrator. The last step is to set the IP for the target host(s), as in:
```bash
export IPS=192.168.1.180
```
or if you have multiple hosts:
```bash
export IPS=192.168.1.180 192.168.181 192.168.1.182
```

Now you're ready to run the deployment.  From the `$MAGMA_ROOT/orc8r/cloud/deploy/bare-metal-ansible' directory, execute:
```bash
./deploy.sh
```

## Cleaning up

Once you've begun a deployment, if you find that you need to go back and start again, it's best to completely reset the environment with:
```bash
cd $MAGMA_ROOT/orc8r/cloud/deploy/bare-metal-ansible/
ansible-playbook -i inventory/cluster.local/hosts.yaml -b kubespray/reset.yml
```
