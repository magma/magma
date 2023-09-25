---
id: version-1.8.0-dev_minikube
title: Deploy on Minikube
hide_title: true
original_id: dev_minikube
---

# Deploy on Minikube

Deploying Orc8r on Minikube is the easiest way to test changes to the Helm charts. Most steps are similar to the main deployment guide, with a few differences.

## Prerequisites

### Spin up Minikube

Set up Minikube with one of the following commands, including sufficient resources and seeding the metrics config files:

- On macOS

  ```bash
  minikube start --cni=bridge --driver=hyperkit --memory=8gb --cpus=8 --kubernetes-version='v1.20.2' --mount --mount-string "${MAGMA_ROOT}/orc8r/cloud/docker/metrics-configs:/configs"
  ```

- On Linux

  ```bash
  minikube start --cni=bridge --driver=docker --memory=8gb --cpus=8 --kubernetes-version='v1.20.2' --mount --mount-string "${MAGMA_ROOT}/orc8r/cloud/docker/metrics-configs:/configs"
  ```

where the only difference between the two is the driver used.

> Note:
>
> - There are a lot of things that can go wrong with spinning up Minikube, so if you have problems check for documentation specific to your system.
> - Make sure you are not connected to a VPN when running this command.
> - For further information on the recommended drivers used above, see [Minikube configs](#minikube-configs).

Now install prerequisites for Orc8r and create the `orc8r` K8s namespace

```bash
helm repo add bitnami https://charts.bitnami.com/bitnami
helm upgrade --install \
    --create-namespace \
    --namespace orc8r \
    --set global.postgresql.auth.postgresPassword=postgres,global.postgresql.auth.database=magma,fullnameOverride=postgresql \
    postgresql \
    bitnami/postgresql
```

### Build and publish images

> NOTE: skip this step if you want to use the official container images at [linuxfoundation.jfrog.io](https://linuxfoundation.jfrog.io/).

There are 2 ways you can publish your own images: to a private registry, or to a localhost registry. Choose an option, then complete the relevant prerequisites:

1. Publish to private registry (specifically, we'll use [DockerHub](https://hub.docker.com/))
    - `docker login`
    - Use `registry.hub.docker.com/USERNAME` as your registry name
2. Publish to localhost registry (specifically, we'll use the [local Docker registry via Minikube](https://minikube.sigs.k8s.io/docs/handbook/registry/#docker-on-macos))
    - `minikube addons enable registry`
    - `docker run --rm -it --network=host alpine ash -c "apk add socat && socat TCP-LISTEN:5000,reuseaddr,fork TCP:$(minikube ip):5000"`
        - This should hang, open new tab for other commands
    - Use `localhost:5000` as your registry name

After completing the prerequisites listed above, follow the instructions at [Building Orchestrator](./dev_build.md#build-and-publish-container-images) to publish container images to the chosen registry.

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
export IMAGE_REGISTRY_URL=linuxfoundation.jfrog.io/magma-docker  # or replace with your registry
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

### Optional: fluentd secrets

To use fluentd on the minikube deployment, create additional fluentd secrets

```bash
cd ${CERTS_DIR}
openssl genrsa -out fluentd.key 2048
openssl req -new -key fluentd.key -out fluentd.csr -subj "/C=US/CN=fluentd.$domain"
openssl x509 -req -in fluentd.csr -CA certifier.pem -CAkey certifier.key -CAcreateserial -out fluentd.pem -days 3650 -sha256
```

Apply the secrets

```bash
cd ${MAGMA_ROOT}/orc8r/cloud/helm/orc8r
helm template orc8r charts/secrets \
  --namespace orc8r \
  --set-file 'secret.certs.files.fluentd\.pem'=${CERTS_DIR}/fluentd.pem \
  --set-file 'secret.certs.files.fluentd\.key'=${CERTS_DIR}/fluentd.key |
  kubectl apply -f -
```

### Create values file

A minimal values file is at `${MAGMA_ROOT}/orc8r/cloud/helm/orc8r/examples/minikube.values.yaml`

- Copy that file to `${MAGMA_ROOT}/orc8r/cloud/helm/orc8r.values.yaml`
- Replace `IMAGE_REGISTRY_URL` with your registry and `IMAGE_TAG` with your tag
- Make additional edits as desired

### Install charts

This section describes how to install based on local charts. However, you can also install charts from the official chart repositories

- Stable: <https://linuxfoundation.jfrog.io/artifactory/magma-helm-prod/>
- Test: <https://linuxfoundation.jfrog.io/artifactory/magma-helm-test/>

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

Optionally install the fluentd charts by editing the file at `${MAGMA_ROOT}/orc8r/cloud/helm/orc8r/charts/logging/values.yaml`

```yaml
fluentd_daemon:
  create: true

  image:
    repository: IMAGE_REGISTRY_URL
    tag: IMAGE_TAG
    pullPolicy: IfNotPresent

  env:
    elastic_host: "host.minikube.internal"
    elastic_port: "9200"
    elastic_scheme: "http"
```

```yaml
fluentd_forward:
  create: true

  # Domain-proxy output
  dp_output: true

  replicas: 1

  nodeSelector: {}
  tolerations: []
  affinity: {}

  image:
    repository: IMAGE_REPOSITORY
    tag: IMAGE_TAG
    pullPolicy: IfNotPresent

  env:
    elastic_host: "host.minikube.internal"
    elastic_port: "9200"
    elastic_scheme: "http"
    elastic_flush_interval: 5s
```

Replace `IMAGE_REGISTRY_URL` with your registry and `IMAGE_TAG` with your tag.
Then install the charts

```bash
helm upgrade --install --namespace orc8r --values ${MAGMA_ROOT}/orc8r/cloud/helm/orc8r.values.yaml orc8r .
```

It may take a couple minutes before all pods are finished being created, but if successful, you should get something like this

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

If using fluentd, there should be the pods `orc8r-fluentd-forward` and `orc8r-fluentd-daemon` as well

```bash
$ kubectl --namespace orc8r get pods

NAME                                             READY   STATUS    RESTARTS   AGE
mysql-57955549d5-n69pd                           1/1     Running   0          6m14s
nms-magmalte-7c84667c4c-pvtlj                    1/1     Running   0          2m58s
nms-nginx-proxy-5b86f479f7-lvjpn                 1/1     Running   0          2m58s
orc8r-alertmanager-57d5d6ccc4-ht4n4              1/1     Running   0          2m58s
orc8r-alertmanager-configurer-76cf8f8f57-rmwjf   1/1     Running   0          2m58s
orc8r-controller-fdf59f456-vvqr2                 1/1     Running   0          2m58s
orc8r-fluentd-forward-54794d4d86-599lv           1/1     Running   0          2m58s
orc8r-nginx-7d6c78647-n6d8z                      1/1     Running   0          2m58s
orc8r-prometheus-77dccb799b-w9z6z                1/1     Running   0          2m58s
orc8r-prometheus-cache-6d647df4d9-2wqkc          1/1     Running   0          2m58s
orc8r-prometheus-configurer-6d6d987c88-pgm87     1/1     Running   0          2m58s
orc8r-user-grafana-6498bb6959-rchx5              1/1     Running   0          2m58s
postgresql-0                                     1/1     Running   4          6d23h
```

```bash
$ kubectl --namespace kube-system get pods

NAME                               READY   STATUS    RESTARTS   AGE
coredns-74ff55c5b-xsz59            1/1     Running   0          3d
etcd-minikube                      1/1     Running   0          3d
kube-apiserver-minikube            1/1     Running   0          3d
kube-controller-manager-minikube   1/1     Running   0          3d
kube-proxy-cp6s5                   1/1     Running   0          3d
kube-scheduler-minikube            1/1     Running   0          3d
orc8r-fluentd-daemon-b7wpt         1/1     Running   0          4h
registry-46g6n                     1/1     Running   0          3d
registry-proxy-n6kml               1/1     Running   0          3d
storage-provisioner                1/1     Running   0          3d
```

Optionally, start the elasticsearch and kibana containers which handle the logs aggregated by fluentd

```bash
cd ${MAGMA_ROOT}/cloud/docker
./run.py
```

### Access logs through Kibana

The Orchestrator logs aggregated by fluentd can be accessed via Kibana in the web browser under `http://localhost:5601/`.
There should be two index patterns: `fluentd` (from fluentd-forward) and `logstash` (from fluentd-daemon).
The `logstash` index pattern collects all logs from the Orchestrator, and `fluentd` the logs from connected gateways.

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
# Tab 1 (should hang)
kubectl --namespace orc8r port-forward svc/orc8r-nginx-proxy 7443:8443 7444:8444 9443:443

# Tab 2
export CERTS_DIR=${MAGMA_ROOT}/.cache/test_certs  # mirrored from above
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

Log in to NMS at <https://magma-test.localhost:8081> using credentials: `admin@magma.test/password1234`

## Appendix

### Minikube configs

We use the [HyperKit driver](https://minikube.sigs.k8s.io/docs/drivers/hyperkit/) on macOS instead of the default [Docker driver](https://minikube.sigs.k8s.io/docs/drivers/docker/) because the former is more performant in supporting [hairpin mode](https://github.com/kubernetes/minikube/issues/1568) (via the [CNI network plugin](https://kubernetes.io/docs/concepts/extend-kubernetes/compute-storage-net/network-plugins/)). I.e., Orc8r has some services that make requests to themselves, which requires K8s support.

You could also accomplish this with `--cni=bridge --driver=docker`, but this causes a 2x increase in pod startup times. Note that the HyperKit driver is not supported on [Linux](https://minikube.sigs.k8s.io/docs/drivers/#linux).
