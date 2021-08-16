## NMS

Install NMS which includes nginx and magmalte subcomponents.

## Install charts

```
$ helm upgrade --install nms ./nms --namespace=magma
```

## Overview

This chart installs the magma NMS. The NMS is the UI for managing, configuring, and monitoring networks.

## Prerequisites

1. we will first need the orc8r to be setup

2. API and Nginx certs present under secrets named `orc8r-secrets-certs`

3. MySql Database created for NMS

4. magmalte image ( build using Docker file https://github.com/magma/magma/blob/master/nms/packages/magmalte/Dockerfile )



## Configuration

The following table list the configurable parameters of the NMS chart and their default values.

| Parameter        | Description     | Default   |
| ---              | ---             | ---       |
| magmalte:        |
| `manifests.secrets` | Enable Magmalte secrets to store mysql info. | `true` |
| `manifests.deployment` | Enable Magmalte deployment | `true` |
| `manifests.service` | Enable Magmalte services. | `true` |
| `env.api_host` | orc8r proxy endpoint. | `[]` |
| `env.host` | Host to bind | `[]` |
| `env.port` | Magmalte service port bind. | `{}` |
| `env.mapbox_access_token` | Mapbox Access token. | `` |
| `env.mysql_host` | MySQL host IP/Name. | `` |
| `env.mysql_db` | NMS Database name. | `` |
| `env.mysql_user` | NMS Database user. | `` |
| `env.mysql_pass` | NMS Database password. | `` |
| `labels.node_selector_key` | Target Node selector label Key. | `` |
| `labels.node_selector_value` | Target Node selector label value. | `` |
| `image.repository` | Repository for NMS Magmalte image. | `nil` |
| `image.tag` | Tag for NMS Magmalte image. | `latest` |
| `image.pullPolicy` | Pull policy for NMS Magmalte image. | `IfNotPresent` |
| `service.type` | Service type for magmalte | `ClusterIP` |
| `service.http.port` | Service port number for magmalte | `8081` |
| `service.http.targetport` | Service targetport number for magmalte | `8081` |
| `service.http.nodePort` | Service nodePort number | `""` |
| Nginx:       |
| `manifests.configmap` | Enable Nginx configmap to store config files. | `true` |
| `manifests.deployment` | Enable Nginx deployment | `true` |
| `manifests.service` | Enable Nginx services. | `true` |
| `labels.node_selector_key` | Target Node selector label Key. | `` |
| `labels.node_selector_value` | Target Node selector label value. | `` |
| `image.repository` | Repository for NMS Nginx image. | `nil` |
| `image.tag` | Tag for NMS Nginx image. | `latest` |
| `image.pullPolicy` | Pull policy for NMS Nginx image. | `IfNotPresent` |
| `service.type` | Service type for Nginx | `ClusterIP` |
| `service.http.port` | Service port number for Nginx | `443` |
| `service.http.targetport` | Service targetport number for Nginx | `443` |
| `service.http.nodePort` | Service nodePort number | `""` |
| Global: |
| `pod.replicas.nginx.server` | Number of instances to deploy for Nginx server. | `1` |
| `pod.replicas.magmalte.server` | Number of instances to deploy for Magmalte server. | `1` |
| `pod.resources.enabled` | Enable resources requests and limits for Pods. | `False` |
| `pod.resources.nginx.requests` | Define resources requests and limits for Nginx Pods. | `{}` |
| `pod.resources.magmalte.requests` | Define resources requests and limits for Magmalte Pods. | `{}` |
| `proxy.nodeSelector` | Define which Nodes the Pods are scheduled on. | `{}` |
| `proxy.tolerations` | If specified, the pod's tolerations. | `[]` |
| `proxy.affinity` | Assign the orchestrator proxy to run on specific nodes. | `{}` |

```
- Install NMS chart:
```bash
$ cat values.yaml
magmalte:
  env:
    mapbox_access_token: ""
    mysql_host: mariadb.magma.svc.cluster.local
    mysql_db: magma
    mysql_user: magma
    mysql_pass: password

  image:
    repository: docker.io/magmalte

$ helm upgrade --install nms ./nms --namespace=magma

```
- Create Admin user:
```bash
kubectl exec -it -n magma $(kubectl get pod -n magma  \
-l app.kubernetes.io/component=magmalte -o jsonpath="{.items[0].metadata.name}") \
-- yarn run setAdminPassword admin@magma.test password1234
```

- NMS Dashboard should be reachable via https://<nginx_svc>

Get nginx_svc with following command 

```bash
kubectl get svc -n magma -l app.kubernetes.io/component=nginx,app.kubernetes.io/instance=nms \ 
-o jsonpath="{.items[0].spec.clusterIP}"
```
