--- 
id: deploy_agw_kubernetes
title: Deploy AGW in Kubernetes
hide_title: true
---
# Deploy AGW in Kubernetes
## Prerequisites
Make sure you setup the master and the worker nodes. They must run Ubuntu Server 20.04 LTS, and the master node must have ssh key access to the worker nodes with the user `magma`.
Optionally you can setup an ansible controller where you will run the deploy script or run it directly from the master node.


The worker nodes must satisfy the following requirements:


- aarch64 or 64bit-X86 machine.
- Two ethernet ports.


The interfaces are renamed to `eth0` and `eth1` during the deploy process. Interface `eth0` will carry any traffic that is not S1. That is,data plane traffic (SGi), control plane traffic (Orc8r HTTP2) and management (ssh). The interface `eth1` carries S1 traffic.


> NOTE:
> - Interface names might have different names on your hardware, so just
> replace `enp1s0` and `enp2s0` with your current interface names
> when following this guide.
> - The `agw_install_docker.sh` script will rename the `enp1s0`
> interface to `eth0`.
> - If you do not want all internet traffic to go through `enp1s0`
> to separate control plane (Orc8r Http2 traffic) from user plane, you
> may want to add another interface and configure proper routing.

## Deployment

### 1. Prepare the cluster

- Configure the user `magma` with SUDO privileges to the master and worker nodes.

- Configure ssh key access from the worker to the nodes for the `magma` user.

- Add your master and worker nodes to the ansible inventory file at `$MAGMA_ROOT/lte/gateway/deploy/agw_k8s_hosts.yml`. Make sure that the `ansible_user` is set to `magma` and the `ansible_sudo_pass` password is right. Set local ansible connection for the master node if you are deploying from it. 

- Copy the `rootCA.pem` file from the orc8r to the following location in the master node:

```bash

mkdir -p /var/opt/magma/certs

vim /var/opt/magma/certs/rootCA.pem

```

### 2. Deploy AGW Kubernetes cluster

- cd to the `$MAGMA_ROOT/lte/gateway/deploy` directory.

- Run the deployment script:

```bash

sudo ./agw_install_k8s.sh

```

- You shall see the joined nodes with the command `kubectl get nodes`.

#### Deploying in single node

If you are running a single node cluster, make sure that the master node is available to the scheduler by untainting it.

```bash

kubectl taint nodes <master-node-name> node-role.kubernetes.io/control-plane:NoSchedule-

```

### 3. Install AGW

Follow the AGW Helm Deployment guide in `$MAGMA_ROOT/lte/gateway/deploy/agwc-helm-charts` directory.

If you are deploying AGW in multiple workers, ensure the Node Selection is configured properly in the Helm charts deployment files.
