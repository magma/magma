## Federated Gateway (FeG)

The federated gateway provides remote procedure call (GRPC) based interfaces to standard 3GPP components, such as HSS (S6a, SWx), OCS (Gy), and PCRF (Gx). The exposed RPC interface provides versioning & backward compatibility, security (HTTP2 & TLS) as well as support for multiple programming languages. The Remote Procedures below provide simple, extensible, multi-language interfaces based on GRPC which allow developers to avoid dealing with the complexities of 3GPP protocols. Implementing these RPC interfaces allows networks running on Magma to integrate with traditional 3GPP core components


## TL;DR;
```bash
$ cd magma/feg/gateway/helm

$ cat vals.yaml
feg:
  image:
    docker_registry: docker.io/feg_
    tag: latest
  repo:
    url: https://github.com/facebookincubator/magma.git
    branch: master

image:
  repository: virtlet.cloud/cloud-images.ubuntu.com/xenial/current/xenial-server-cloudimg-amd64-disk1.img
  pullPolicy: IfNotPresent

$ helm upgrade --install feg --namespace magma ./feg --values=vals.yaml
```

## Overview

This chart installs the Magma Federated Gateway.

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

4. Build images for feg

```shell
   export COMPOSE_PROJECT_NAME=feg

   cd magma/feg/gateway/docker/

   source .env

   docker-compose -f docker-compose.yml build

   ./magma/orc8r/tools/docker/publish.sh -r REGISTRY -i gateway_go -u USERNAME -p passwordfile

   ./magma/orc8r/tools/docker/publish.sh -r REGISTRY -i gateway_python -u USERNAME -p passwordfile

```

## Configuration

The following table list the configurable parameters of the orchestrator chart and their default values.

| Parameter        | Description     | Default   |
| ---              | ---             | ---       |
| `manifests.configmap_env` | Enable feg configmap env. | `True` |
| `manifests.deployment` | Enable feg deployment. | `True` |
| `manifests.service` | Enable feg service. | `True` |
| `manifests.rbac` | Enable feg rbac. | `True` |
| `secrets.create` | Enable feg secret creation | `False` |
| `secret.gwinfo`   | Secret name containing feg gwinfo | `feg-secrets-gwinfo` |
| `feg.type` | Gateway type agrument. | `feg` |
| `feg.image.docker_registry` | FeG docker registry host. | `docker.io` |
| `feg.image.tag` | FeG docker images tag. | `latest` |
| `feg.image.username` | Docker registry username. | `` |
| `feg.image.password` | Docker registry password. | `` |
| `feg.proxy.local_port` | FeG proxy local port. | `8443` |
| `feg.proxy.cloud_address` | FeG proxy cloud address. | `orc8r-proxy` |
| `feg.proxy.cloud_port` | FeG proxy Cloud port. | `9443` |
| `feg.proxy.bootstrap_address` | FeG proxy bootstrap address. | `orc8r-bootstrap` |
| `feg.proxy.bootstrap_port` | FeG proxy bootstrap port. | `9444` |
| `feg.repo.url` | FeG magma repo url. | `https://github.com/facebookincubator/magma/` |
| `feg.repo.branch` | FeG magma repo branch. | `master` |
| `feg.bind.S6A_LOCAL_PORT` | FeG S6A local port. | `3868` |
| `feg.bind.S6A_HOST_PORT` | FeG S6A host port. | `3869` |
| `feg.bind.S6A_NETWORK` | FeG S6A network type. | `sctp` |
| `feg.bind.SWX_LOCAL_PORT` | FeG SWX local port. | `3868` |
| `feg.bind.SWX_HOST_PORT` | FeG SWX host port. | `3868` |
| `feg.bind.SWX_NETWORK` | FeG SWX network type. | `sctp` |
| `feg.bind.GX_LOCAL_PORT` | FeG GX local port. | `3907` |
| `feg.bind.GX_HOST_PORT` | FeG GX host port. | `0` |
| `feg.bind.GX_NETWORK` | FeG GX network type. | `tcp` |
| `feg.bind.GY_LOCAL_PORT` | FeG GY local port. | `3906` |
| `feg.bind.GY_HOST_PORT` | FeG GY host port. | `0` |
| `feg.bind.GY_NETWORK` | FeG GY network type. | `tcp` |
| `image.repository` | Virtlet image path | `virtlet.cloud/<image_path>` |
| `image.pullPolicy` | Virtlet Image pullpolicy | `IfNotPresent` |
| `labels.node_selector_key` | Target Node selector label Key. | `extraRuntime` |
| `labels.node_selector_value` | Target Node selector label value. | `virtlet` |
| `pod.replicas.server` | Number of instances to deploy for FeG server. | `1` |
| `pod.resources.enabled` | Enable resources requests and limits for Pods. | `False` |
| `pod.resources.server.requests` | Define resources requests and limits for FeG Pods. | `{}` |
| `nodeSelector` | Define which Nodes the Pods are scheduled on. | `{}` |
| `tolerations` | If specified, the pod's tolerations. | `[]` |
| `affinity` | Assign the FeG to run on specific nodes. | `{}` |
| `virtlet.vcpuCount` | Number of vcpu assigned to FeG VM. | `1` |
| `virtlet.rootVolumeSize` | Size of root volume assigned to FeG VM. | `10Gi` |
| `virtlet.diskDriver` | Virtlet disk driver type. | `virtio` |
| `virtlet.sshKeys` | default user ssh key for user ubuntu | `` |
| `virtlet.ssh_pwauth` | To Enable/Disable password auth. | `True` |
| `virtlet.user` | Add New user. | `` |


## Installation steps

1. Create persistent gateway info (optional)

    If you want your gateway pod to have the same gwinfo on pod 
    recreation, first follow the steps below.
    
    #### Creating Gateway Info
    If creating a gateway for the first time, you'll need to create a snowflake
    and challenge key before installing the gateway. To do so:

    ```
    $ docker login <DOCKER REGISTRY>
    $ docker pull <DOCKER REGISTRY>/gateway_python:<IMAGE_VERSION>
    $ docker run -d <DOCKER_REGISTRY>/gateway_python:<IMAGE_VERSION> python3.5 -m magma.magmad.main

    This will output a container ID such as
    f3bc383a95db16f2e448fdf67cac133a5f9019375720b59477aebc96bacd05a9

    Run the following, substituting your container ID:
    $ docker cp <container ID>:/etc/snowflake charts/secrets/.secrets
    $ docker cp <container ID>:/var/opt/magma/certs/gw_challenge.key /charts/secrets/.secrets 
    ```
   
    If you're instead upgrading your gateway to have persistent gwinfo,
    copy the `etc/snowflake` and `/var/opt/magma/certs/gw_challenge.key` from 
    your gateway to `charts/secrets/.secrets` of where this chart is stored.

    Ensure that `secrets.create` is set to true in your vals.yaml override

2. Install FeG 

	helm upgrade --install feg --namespace magma orc8r --values=vals.yaml

3. Register the gateway with the orchestrator

    Login to the Feg VM and execute below command:

    ```shell
   
   a. Get IP of FeG VM:
   
         feg_ip=$(kubectl -n magma get svc -l app.kubernetes.io/component=gateway,app.kubernetes.io/instance=feg \ 
         -o jsonpath="{.items[0].spec.clusterIP}")
    
   b. SSH into VM
   
         ssh testuser@${feg_ip} "/var/opt/magma/docker/docker-compose exec magmad /usr/local/bin/show_gateway_info.py"
      
   c. Note down the H/w Id and Challenge Key: 
  
  
      Hardware ID:
      ------------
      5cc30126-a218-4492-b654-485fc7bdac6f

      Challenge Key:
      -----------
      MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAE2lAV8Dj1ZQEeQlJ/M9/iYXmiVLC7l5QU7IvrNe+lLsu2MuGz4hjNwFPLmG      /x055Zqzh++8LsXQSKJ0mgV9AUB87xyFGt1wGjvaUa8Jea1ZMRMd1lJ+IsKA606HeaQfVq

    ```

4. Login to NMS Dahsboard and register New Gateway


