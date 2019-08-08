## Carrier Wireless Gateway (cwf)

This Chart helps in deploying CWF on Kubernetes as pods with help of Virtlet ( https://docs.virtlet.cloud/user-guide/real-cluster/ )

## TL;DR;
```bash
$ cd magma/cwf/gateway/helm

$ cat vals.yaml
cwf:
  image:
    docker_registry: docker.io/cwf_
    tag: latest
  repo:
    url: https://github.com/facebookincubator/magma.git
    branch: master
image:
  repository: virtlet.cloud/cloud-images.ubuntu.com/bionic/current/bionic-server-cloudimg-amd64.img
  pullPolicy: IfNotPresent

$ helm upgrade --install cwf --namespace magma ./cwf --values=vals.yaml
```

## Overview

This chart installs the Magma CWF Gateway.

## Prerequisites
1. We will first need the orc8r to be setup

2. Check if rootCA.pem is present under secrets named `orc8r-secrets-certs`

```bash

   kubectl -n magma describe  secrets orc8r-secrets-certs
   
   Data
   ====
   rootCA.key:              1675 bytes
```
3. Install Virtlet to run Virtual Machine in POD

	https://docs.virtlet.cloud/user-guide/real-cluster/ 

4. Build images for cwf

```shell
   export COMPOSE_PROJECT_NAME=cwf

   cd magma/cwf/gateway/docker/

   source .prod_env

   docker-compose build

   ./magma/orc8r/tools/docker/publish.sh -r REGISTRY -i gateway_c -u USERNAME -p passwordfile

   ./magma/orc8r/tools/docker/publish.sh -r REGISTRY -i gateway_python -u USERNAME -p passwordfile

```

## Configuration

The following table list the configurable parameters of the orchestrator chart and their default values.

| Parameter        | Description     | Default   |
| ---              | ---             | ---       |
| `manifests.configmap_bin` | Enable cwf configmap bin. | `True` |
| `manifests.configmap_env` | Enable cwf configmap env. | `True` |
| `manifests.secrets` | Enable cwf secrets. | `True` |
| `manifests.deployment` | Enable cwf deployment. | `True` |
| `manifests.service` | Enable cwf service. | `True` |
| `manifests.rbac` | Enable cwf rbac. | `True` |
| `cwf.type` | Gateway type agrument. | `cwf` |
| `cwf.image.docker_registry` | CWF docker registry host. | `docker.io` |
| `cwf.image.tag` | CWF docker images tag. | `latest` |
| `cwf.image.username` | Docker registry username. | `` |
| `cwf.image.password` | Docker registry password. | `` |
| `cwf.proxy.local_port` | CWF proxy local port. | `8443` |
| `cwf.proxy.cloud_address` | CWF proxy cloud address. | `orc8r-proxy` |
| `cwf.proxy.cloud_port` | CWF proxy Cloud port. | `9443` |
| `cwf.proxy.bootstrap_address` | CWF proxy bootstrap address. | `orc8r-bootstrap` |
| `cwf.proxy.bootstrap_port` | CWF proxy bootstrap port. | `9444` |
| `cwf.repo.url` | CWF magma repo url. | `https://github.com/facebookincubator/magma/` |
| `cwf.repo.branch` | CWF magma repo branch. | `master` |
| `image.repository` | Virtlet image path | `virtlet.cloud/<image_path>` |
| `image.pullPolicy` | Virtlet Image pullpolicy | `IfNotPresent` |
| `labels.node_selector_key` | Target Node selector label Key. | `extraRuntime` |
| `labels.node_selector_value` | Target Node selector label value. | `virtlet` |
| `pod.replicas.server` | Number of instances to deploy for CWF server. | `1` |
| `pod.resources.enabled` | Enable resources requests and limits for Pods. | `False` |
| `pod.resources.server.requests` | Define resources requests and limits for CWF Pods. | `{}` |
| `nodeSelector` | Define which Nodes the Pods are scheduled on. | `{}` |
| `tolerations` | If specified, the pod's tolerations. | `[]` |
| `affinity` | Assign the CWF to run on specific nodes. | `{}` |
| `virtlet.vcpuCount` | Number of vcpu assigned to CWF VM. | `1` |
| `virtlet.rootVolumeSize` | Size of root volume assigned to CWF VM. | `10Gi` |
| `virtlet.diskDriver` | Virtlet disk driver type. | `virtio` |
| `virtlet.sshKeys` | default user ssh key for user ubuntu | `` |
| `virtlet.ssh_pwauth` | To Enable/Disable password auth. | `True` |
| `virtlet.user` | Add New user. | `` |


## Installation steps

1. Install CWF

	helm upgrade --install cwf --namespace magma ./cwf --values=vals.yaml

2. Register the gateway with the orchestrator

login to cwf VM and execute below command

```shell
   
   a. Get IP of CWF VM:
   
         cwf_ip=$(kubectl -n magma get svc -l app.kubernetes.io/component=gateway,app.kubernetes.io/instance=cwf \ 
         -o jsonpath="{.items[0].spec.clusterIP}")
    
   b. SSH into VM
   
         ssh -t testuser@${cwf_ip} "cd /var/opt/magma/docker ; sudo docker-compose exec magmad /usr/local/bin/show_gateway_info.py"
      
   c. Note down the H/w Id and Challenge Key: 
  
  
      Hardware ID:
      ------------
      5cc30126-a218-4492-b654-485fc7bdac6f

      Challenge Key:
      -----------
      MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAE2lAV8Dj1ZQEeQlJ/M9/iYXmiVLC7l5QU7IvrNe+lLsu2MuGz4hjNwFPLmG      /x055Zqzh++8LsXQSKJ0mgV9AUB87xyFGt1wGjvaUa8Jea1ZMRMd1lJ+IsKA606HeaQfVq

```

3. Login to NMS Dahsboard and register New Gateway


