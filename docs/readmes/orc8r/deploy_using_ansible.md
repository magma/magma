---
id: deploy_using_ansible
title: Deploy Orchestrator using Ansible (Beta)
hide_title: true
---

# Deploy Orchestrator using Ansible (Beta)

This how-to guide can be used to deploy Magma's Orchestrator on any cloud environment. 
It contains roles to set up a Kubernetes cluster and deploy Magma Orchestrator using helm charts.
For more information on Magma Deployer, please visit the project's
[magma-deployer](https://github.com/magma/magma-deployer).

> magma-deployer is in Beta and is not yet production ready or feature complete.

## Pre-requisites

- Ubuntu Jammy 22.04 VM / Baremetal machine 
- RAM: 8GB
- CPU: 4 cores
- Storage: 100GB


Quick Install:
```bash
sudo bash -c "$(curl -sL https://github.com/magma/magma-deployer/raw/main/deploy-orc8r.sh)"
```

Switch to `magma` user after deployment has finsished:
```bash
sudo su - magma
```

Once all pods are ready, setup NMS login:
```bash
cd ~/magma-deployer
ansible-playbook config-orc8r.yml
```

You can get your `rootCA.pem` file from the following location:
```bash
cat ~/magma-deployer/secrets/rootCA.pem
```
