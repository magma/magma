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

## Deploy Orchestrator

Quick Install:
```
sudo bash -c "$(curl -sL https://github.com/magma/magma-deployer/raw/main/deploy-orc8r.sh)"
```

Switch to `magma` user after deployment has finsished:
```
sudo su - magma
```

Check if all pods are ready
```
kubectl get pods
```

setup NMS login:
```
cd ~/magma-deployer
ansible-playbook config-orc8r.yml
```

Update `/etc/hosts` file with the following entries:
```
10.86.113.153 api.magma.local
10.86.113.153 magma-test.nms.magma.local
10.86.113.153 fluentd.magma.local
10.86.113.153 controller.magma.local
10.86.113.153 bootstrapper-controller.magma.local
```
> Replace the IP address with the one you see in `kubectl get nodes -o wide`


You can access NMS at the following URL:

https://magma-test.nms.magma.local


You can get your `rootCA.pem` file from the following location:
```
cat ~/magma-deployer/secrets/rootCA.pem
```
