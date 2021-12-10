---
id: version-1.4.0-dev_minikube_deployment
title: Deploy on Minikube
hide_title: true
original_id: dev_minikube_deployment
---

# Deploy on Minikube

Deploying orchestrator on minikube is the easiest way to test changes to the helm
charts. Most steps are similar to the main deployment guide, but there are some
significant differences and many things in there you don't need to worry about.

> Note: All terminal commands are run from HOST unless otherwise specified

> Helm commands are written for Helm 3

## Steps

### Build and publish images

Follow the instructions at [Building Orchestrator](./deploy_build.md#build-and-publish-container-images).

In the end you should have your container images published to a registry.

### Setup Minikube and Helm

Setup minikube with the following command to give it enough resources and seed the
metrics config files:

```
$ minikube start --memory=8192 --cpus=8 --kubernetes-version=v1.18.0 --mount --mount-string "${MAGMA_ROOT}/orc8r/cloud/docker/metrics-configs:/configs"
```

> Note: This has been tested on MacOS. There are a lot of things that can go wrong
> with spinning up minikube, so if you have problems check for documentation specific
> to your system. Also make sure you are not connected to a VPN when running this command.

Now you can install preqrequisites for orc8r and create the magma namespace in kubernetes:

```
$ helm install \
    postgresql \
    --create-namespace \
    --namespace magma \
    --set postgresqlPassword=postgres,postgresqlDatabase=magma,fullnameOverride=postgresql \
    bitnami/postgresql
```

Mysql is a requirement to run the NMS (you can skip this step if you don't want the NMS)

```
$ helm install mysql \
  --namespace magma \
  --set mysqlRootPassword=password,mysqlUser=magma,mysqlPassword=password,mysqlDatabase=magma \
    stable/mysql
```

> Note: You may need to run `helm repo add bitnami https://charts.bitnami.com/bitnami` if the chart is not found

### Generate Secrets

First you'll need to create a few certs and the kubernetes secrets that orchestrator
uses.

Generate the NMS certificate

```
$ cd ${MAGMA_ROOT}/.cache/test_certs
$ openssl req -nodes -new -x509 -batch -keyout nms_nginx.key -out nms_nginx.pem -subj "/CN=*.localhost"
```

Move the certs to the charts directory

```
$ cd ${MAGMA_ROOT}/orc8r/cloud/helm/orc8r
$ mkdir -p charts/secrets/.secrets/certs
$ cp -r ../../../../.cache/test_certs/* charts/secrets/.secrets/certs/.
```

Create the kubernetes secrets

```
helm template orc8r charts/secrets \
    --namespace magma \
    --set-string secret.certs.enabled=true \
    --set-file secret.certs.files."rootCA\.pem"=charts/secrets/.secrets/certs/rootCA.pem \
    --set-file secret.certs.files."bootstrapper\.key"=charts/secrets/.secrets/certs/bootstrapper.key \
    --set-file secret.certs.files."controller\.crt"=charts/secrets/.secrets/certs/controller.crt \
    --set-file secret.certs.files."controller\.key"=charts/secrets/.secrets/certs/controller.key \
    --set-file secret.certs.files."admin_operator\.pem"=charts/secrets/.secrets/certs/admin_operator.pem \
    --set-file secret.certs.files."admin_operator\.key\.pem"=charts/secrets/.secrets/certs/admin_operator.key.pem \
    --set-file secret.certs.files."certifier\.pem"=charts/secrets/.secrets/certs/certifier.pem \
    --set-file secret.certs.files."certifier\.key"=charts/secrets/.secrets/certs/certifier.key \
    --set-file secret.certs.files."nms_nginx\.pem"=charts/secrets/.secrets/certs/nms_nginx.pem \
    --set-file secret.certs.files."nms_nginx\.key\.pem"=charts/secrets/.secrets/certs/nms_nginx.key \
    --set=docker.registry=$DOCKER_REGISTRY \
    --set=docker.username=$DOCKER_USERNAME \
    --set=docker.password=$DOCKER_PASSWORD |
    kubectl apply -f -
```

### Create Values File

A minimum values file that can be used to deploy orc8r is at `${MAGMA_ROOT}/orc8r/cloud/helm/orc8r/examples/minikube_values.yml`.
Use that file and make sure to replace `<DOCKER_REGISTRY>` with your registry and `<TAG>` with your tag.

Save this file wherever you want.

### Install with Helm

```
$ cd ${MAGMA_ROOT}/orc8r/cloud/helm/orc8r
$ helm dep update
$ helm install orc8r --namespace magma . --values=<path-to-values-file>
```

If everything worked you should get something like this:

```
$ kubectl -n magma get pods

NAME                                             READY   STATUS    RESTARTS   AGE
mysql-57955549d5-n69pd                           1/1     Running   0          6m14s
nms-magmalte-7c84667c4c-pvtlj                    1/1     Running   0          2m58s
nms-nginx-proxy-5b86f479f7-lvjpn                 1/1     Running   0          2m58s
orc8r-alertmanager-57d5d6ccc4-ht4n4              1/1     Running   0          2m58s
orc8r-alertmanager-configurer-76cf8f8f57-rmwjf   1/1     Running   0          2m58s
orc8r-controller-fdf59f456-vvqr2                 1/1     Running   0          2m58s
orc8r-nginx-7d6c78647-n6d8z                      1/1     Running   0          2m58s
orc8r-prometheus-77dccb799b-w9z6z                1/1     Running   0          2m58s
orc8r-prometheus-cache-6d647df4d9-2wqkc          1/1     Running   0          2m58s
orc8r-prometheus-configurer-6d6d987c88-pgm87     1/1     Running   0          2m58s
orc8r-user-grafana-6498bb6959-rchx5              1/1     Running   0          2m58s
postgresql-0                                     1/1     Running   4          6d23h
```

### Access API

Create an admin user:

```
kubectl exec -it -n magma \
    $(kubectl get pod -n magma -l app.kubernetes.io/component=controller -o jsonpath="{.items[0].metadata.name}") -- \
    /var/opt/magma/bin/accessc add-existing -admin -cert /var/opt/magma/certs/admin_operator.pem admin_operator
```

Now make sure that the API and your certs are working:

```
# Tab 1
kubectl -n magma port-forward $(kubectl -n magma get pods -l app.kubernetes.io/component=nginx-proxy -o jsonpath="{.items[0].metadata.name}") 9443:9443
```

```
# Tab 2
$ curl -k \
  --cert MAGMA_ROOT/orc8r/cloud/helm/orc8r/charts/secrets/.secrets/certs/admin_operator.pem \
  --key MAGMA_ROOT/orc8r/cloud/helm/orc8r/charts/secrets/.secrets/certs/admin_operator.key.pem \
https://localhost:9443
"hello"
```

### Access NMS

Follow the instructions to [Create an admin user](./deploy_install.md#create-an-nms-admin-user)

Port-forward nginx:

```
kubectl -n magma port-forward svc/nginx-proxy  8443:443
```

Login to NMS at https://magma-test.localhost:8443 using credentials: admin@magma.test/password1234
