---
id: dev_minikube
title: Deploy on Minikube
hide_title: true
---

# Deploy on Minikube

Deploying Orc8r on Minikube is the easiest way to test changes to the Helm charts. Most steps are similar to the main deployment guide, with a few differences.

## Prerequisites

### Build and publish images

> NOTE: you can skip this step if you want to use the official container images at [artifactory.magmacore.org](https://artifactory.magmacore.org/).

Follow the instructions at [Building Orchestrator](./dev_build.md#build-and-publish-container-images).

In the end, your container images should be published to a registry.

### Spin up Minikube

Set up Minikube with the following command, including sufficient resources and seeding the metrics config files

```bash
minikube start --cni=bridge --driver=hyperkit --memory=8gb --cpus=8 --mount --mount-string "${MAGMA_ROOT}/orc8r/cloud/docker/metrics-configs:/configs"
```

> Note: This has been tested on MacOS. There are a lot of things that can go wrong with spinning up Minikube, so if you have problems check for documentation specific to your system. Also make sure you are not connected to a VPN when running this command.

Now install prerequisites for Orc8r and create the `orc8r` K8s namespace

```bash
helm repo add bitnami https://charts.bitnami.com/bitnami
helm upgrade --install \
    --create-namespace \
    --namespace orc8r \
    --set postgresqlPassword=postgres,postgresqlDatabase=magma,fullnameOverride=postgresql \
    postgresql \
    bitnami/postgresql
```

## Install

### Generate secrets

If you haven't already, the easiest way to generate these secrets is temporarily spinning up a local, Docker-based deployment

```bash
export CERTS_DIR=${MAGMA_ROOT}/.cache/test_certs
cd ${MAGMA_ROOT}/orc8r/cloud/docker && ./build.py && ./run.py && sleep 30 && docker-compose down && ls -l ${CERTS_DIR} && cd -
```

### Apply secrets

Create the K8s secrets

```bash
export IMAGE_REGISTRY_URL=docker.artifactory.magmacore.org  # or replace with your registry
export IMAGE_REGISTRY_USERNAME=''
export IMAGE_REGISTRY_PASSWORD=''

export CERTS_DIR=${MAGMA_ROOT}/.cache/test_certs  # mirrored from above

cd ${MAGMA_ROOT}/orc8r/cloud/helm/orc8r
helm template orc8r charts/secrets \
  --namespace orc8r \
  --set-string secret.certs.enabled=true \
  --set-file 'secret.certs.files.rootCA\.pem'=${CERTS_DIR}/rootCA.pem \
  --set-file 'secret.certs.files.bootstrapper\.key'=${CERTS_DIR}/bootstrapper.key \
  --set-file 'secret.certs.files.controller\.crt'=${CERTS_DIR}/controller.crt \
  --set-file 'secret.certs.files.controller\.key'=${CERTS_DIR}/controller.key \
  --set-file 'secret.certs.files.admin_operator\.pem'=${CERTS_DIR}/admin_operator.pem \
  --set-file 'secret.certs.files.admin_operator\.key\.pem'=${CERTS_DIR}/admin_operator.key.pem \
  --set-file 'secret.certs.files.certifier\.pem'=${CERTS_DIR}/certifier.pem \
  --set-file 'secret.certs.files.certifier\.key'=${CERTS_DIR}/certifier.key \
  --set-file 'secret.certs.files.nms_nginx\.pem'=${CERTS_DIR}/controller.crt \
  --set-file 'secret.certs.files.nms_nginx\.key\.pem'=${CERTS_DIR}/controller.key \
  --set=docker.registry=${IMAGE_REGISTRY_URL} \
  --set=docker.username=${IMAGE_REGISTRY_USERNAME} \
  --set=docker.password=${IMAGE_REGISTRY_PASSWORD} |
  kubectl apply -f -
```

### Create values file

A minimal values file is at `${MAGMA_ROOT}/orc8r/cloud/helm/orc8r/examples/minikube_values.yml`

- Copy that file to `${MAGMA_ROOT}/orc8r/cloud/helm/orc8r.values.yaml`
- Replace `IMAGE_REGISTRY_URL` with your registry and `IMAGE_TAG` with your tag
- Make additional edits as desired


### Install charts

This section describes how to install based on local charts. However, you can also install charts from the official chart repositories
- Stable: https://artifactory.magmacore.org/artifactory/helm/
- Test: https://artifactory.magmacore.org/artifactory/helm-test/

Install base `orc8r` chart

```bash
cd ${MAGMA_ROOT}/orc8r/cloud/helm/orc8r
helm dep update
helm upgrade --install --namespace orc8r --values ${MAGMA_ROOT}/orc8r/cloud/helm/orc8r.values.yaml orc8r .
```

Optionally install other charts, like the `lte` charts, similarly

```bash
cd ${MAGMA_ROOT}/lte/cloud/helm/lte-orc8r
helm dep update
helm upgrade --install --namespace orc8r --values ${MAGMA_ROOT}/orc8r/cloud/helm/lte.values.yaml lte .
```

If successful, you should get something like this

```bash
$ kubectl --namespace orc8r get pods

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

## Configure

### Access Orc8r

Create an Orc8r admin user

```bash
kubectl exec -it --namespace orc8r deploy/orc8r-orchestrator -- \
  /var/opt/magma/bin/accessc \
  add-existing -admin -cert /var/opt/magma/certs/admin_operator.pem \
  admin_operator
```

Now ensure the API and your certs are working

```bash
# Tab 1
kubectl --namespace orc8r port-forward svc/orc8r-nginx-proxy 7443:8443 7444:8444 9443:443

# Tab 2
# Assumes CERTS_DIR is set
curl \
  --insecure \
  --cert ${CERTS_DIR}/admin_operator.pem \
  --key ${CERTS_DIR}/admin_operator.key.pem \
  https://localhost:9443

# Should output: "hello"
```

### Access NMS

Follow the instructions to [create an NMS admin user](./deploy_install.md#create-an-nms-admin-user)

Port-forward Nginx

```bash
kubectl --namespace orc8r port-forward svc/nginx-proxy 8081:443
```

Log in to NMS at https://magma-test.localhost:8081 using credentials: `admin@magma.test/password1234`

## Appendix

### Minikube configs

We use the [HyperKit driver](https://minikube.sigs.k8s.io/docs/drivers/hyperkit/) instead of the default [Docker driver](https://minikube.sigs.k8s.io/docs/drivers/docker/) because the former is more performant in supporting [hairpin mode](https://github.com/kubernetes/minikube/issues/1568) (via the [CNI network plugin](https://kubernetes.io/docs/concepts/extend-kubernetes/compute-storage-net/network-plugins/)). I.e., Orc8r has some services that make requests to themselves, which requires K8s support.

You could also accomplish this with `--cni=bridge --driver=docker`, but this causes a 2x increase in pod startup times.
