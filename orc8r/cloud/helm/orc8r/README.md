## Orchestrator

Installs orchestrator which includes proxy / controller subcomponents.

## TL;DR;
```bash
$ cat vals.yaml
imagePullSecrets:
  - name: orc8r-secrets-registry
proxy:
  image:
    repository: docker.io/proxy
  spec:
    hostname: controller.magma.test
controller:
  image:
    repository: docker.io/controller

$ helm install --name orc8r --namespace magma orc8r --values=vals.yaml
```

## Overview

This chart installs the magma orchestrator. The chart expects a set of secrets
to exists during installation, See charts/secrets for more info.

## Configuration

The following table list the configurable parameters of the orchestrator chart and their default values.

| Parameter        | Description     | Default   |
| ---              | ---             | ---       |
| `imagePullSecrets` | Reference to one or more secrets to be used when pulling images. | `[]` |
| `secrets.create` | Create orchestrator secrets. See charts/secrets subchart. | `false` |
| `secret.certs` | Secret name containing orchestrator certs. | `orc8r-secrets-certs` |
| `secret.configs` | Secret name containing orchestrator configs. | `orc8r-secrets-configs` |
| `secret.envdir` | Secret name containing orchestrator envdir. | `orc8r-secrets-envdir` |
| `proxy.service.annotations` | Annotations to be added to the proxy service. | `{}` |
| `proxy.service.labels` | Proxy service labels. | `{}` |
| `proxy.service.type` | Proxy service type. | `ClusterIP` |
| `proxy.service.port.clientcert.port` | Proxy client certificate service external port. | `9443` |
| `proxy.service.port.clientcert.targetPort` | Proxy client certificate service internal port. | `9443` |
| `proxy.service.port.clientcert.nodePort` | Proxy client certificate service node port. | `nil` |
| `proxy.service.port.open.port` | Proxy open service external port. | `9444` |
| `proxy.service.port.open.targetPort` | Proxy open service internal port. | `9444` |
| `proxy.service.port.open.nodePort` | Proxy open service node port. | `nil` |
| `proxy.image.repository` | Repository for orchestrator proxy image. | `nil` |
| `proxy.image.tag` | Tag for orchestrator proxy image. | `latest` |
| `proxy.image.pullPolicy` | Pull policy for orchestrator proxy image. | `IfNotPresent` |
| `proxy.spec.hostname` | Magma controller domain name. | `""` |
| `proxy.replicas` | Number of instances to deploy for orchestrator proxy. | `1` |
| `proxy.resources` | Define resources requests and limits for Pods. | `{}` |
| `proxy.nodeSelector` | Define which Nodes the Pods are scheduled on. | `{}` |
| `proxy.tolerations` | If specified, the pod's tolerations. | `[]` |
| `proxy.affinity` | Assign the orchestrator proxy to run on specific nodes. | `{}` |
| `controller.service.annotations` | Annotations to be added to the controller service. | `{}` |
| `controller.service.labels` | Controller service labels. | `{}` |
| `controller.service.type` | Controller service type. | `ClusterIP` |
| `controller.service.port` | Controller web service external port. | `8080` |
| `controller.service.targetPort` | Controller web service internal port. | `8080` |
| `controller.service.portStart` | Controller service port range start. | `9079` |
| `controller.service.portEnd` | Controller service inclusive port range end. | `9108` |
| `controller.image.repository` | Repository for orchestrator controller image. | `nil` |
| `controller.image.tag` | Tag for orchestrator controller image. | `latest` |
| `controller.image.pullPolicy` | Pull policy for orchestrator controller image. | `IfNotPresent` |
| `controller.spec.postgres.db` | Postgres database name. | `magma` |
| `controller.spec.postgres.host` | Postgres database host. | `postgresql` |
| `controller.spec.postgres.port` | Postgres database port. | `5432` |
| `controller.spec.postgres.user` | Postgres username. | `postgres` |
| `controller.spec.postgres.pass` | Postgres password. | `postgres` |
| `controller.replicas` | Number of instances to deploy for orchestrator controller. | `1` |
| `controller.resources` | Define resources requests and limits for Pods. | `{}` |
| `controller.nodeSelector` | Define which Nodes the Pods are scheduled on. | `{}` |
| `controller.tolerations` | If specified, the pod's tolerations. | `[]` |
| `controller.affinity` | Assign the orchestrator proxy to run on specific nodes. | `{}` |

## Running in Minikube
- Start Minikube with 8192 MB of memory and 4 CPUs. This example uses Kuberenetes version 1.14.1 and uses [Minikube Hypervisor Driver](https://kubernetes.io/docs/tasks/tools/install-minikube/#install-a-hypervisor):
```bash
$ minikube start --memory=8192 --cpus=4 --kubernetes-version=v1.14.1
```
- Install Helm Tiller:
```bash
$ kubectl apply -f tiller-rbac-config.yaml
$ helm init --service-account tiller --history-max 200
# Wait for tiller to become 'Running'
$ kubectl get pods -n kube-system | grep tiller
```
- Create a namespace for orchestrator components:
```bash
$ kubectl create namespace magma
```
- Install Postgres Helm chart:
```bash
$ helm install \
    --name postgresql \
    --namespace magma \
    --set postgresqlPassword=postgres,postgresqlDatabase=magma,fullnameOverride=postgresql \
    stable/postgresql
```
- Copy orchestrator secrets:
```bash
cd magma/orc8r/cloud/helm/orc8r
mkdir -p charts/secrets/.secrets/certs
# You need to add the following files to the certs directory:
#   bootstrapper.key certifier.key certifier.pem vpn_ca.crt vpn_ca.key
#   controller.crt controller.key rootCA.pem
# The controller.crt, controller.key and rootCA.pem are the certificate info
# for your public domain name.
# For local testing, you can do the following after running Orc8r using docker:
cp -r ../../../../.cache/test_certs/* charts/secrets/.secrets/certs/.
```
- Install orchestrator secrets:
```bash
export DOCKER_REGISTRY=<registry>
export DOCKER_USERNAME=<username>
export DOCKER_PASSWORD=<password>
helm template charts/secrets \
    --name orc8r-secrets \
    --namespace magma \
    --set=docker.registry=${DOCKER_REGISTRY} \
    --set=docker.username=${DOCKER_USERNAME} \
    --set=docker.password=${DOCKER_PASSWORD} \
     | kubectl apply -f -
```
- Install orchestrator chart:
```bash
$ cat vals.yaml
imagePullSecrets:
  - name: orc8r-secrets-registry
proxy:
  image:
    repository: docker.io/proxy
  spec:
    hostname: controller.magma.test
controller:
  image:
    repository: docker.io/controller

$ helm install --name orc8r --namespace magma . --values=vals.yaml

# In the future, if you want to upgrade the helm chart, run:
$ helm upgrade orc8r . -f --values=vals.yaml
```
- Add the admin in the datastore:
```bash
kubectl exec -it -n magma \
    $(kubectl get pod -n magma -l app.kubernetes.io/component=controller -o jsonpath="{.items[0].metadata.name}") -- \
    /var/opt/magma/bin/accessc add-existing -admin -cert /var/opt/magma/certs/admin_operator.pem admin_operator
```
- Port forward traffic to orchestrator proxy:
```bash
kubectl port-forward -n magma svc/orc8r-proxy 9443:9443

# If using minikube, run:
minikube service orc8r-proxy -n magma --https
```
- Orchestrator proxy should be reachable via https://localhost:9443 and
requires magma client certificate to be installed on browser.