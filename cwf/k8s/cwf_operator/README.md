CWF Kubernetes Operator
===

The CWF Operator is a Kubernetes Operator used to deploy the CWF Helm
chart (located at `magma/cwf/gateway/helm`) in a highly available
configuration. The only configuration currently supported is 2 CWAGs
running in active/standby mode.

**Note**: This feature only works AP/WLCs that support redundant
RADIUS servers and GRE tunnels. Magma is in the process of removing this
requirement.

## Prerequisites

Before installing the CWF Operator, the following prereqs must be done to
ensure that HA will work properly:
* Redis-ha helm chart deployed
* CWAGs configured to use deployed redis server
* Stateless sessiond enabled
* RADIUS server configured to use redis

### Deploying HA Redis

We can do this by deploying the `redis-ha` helm chart. Below are the
recommended value overrides to deploy with:

```console
redis:
  port: 6380

persistentVolume:
  enabled: false
hostPath:
  path: "/data/{{ .Release.Name }}"
nodeSelector:
  node:mysql

## Enables a HA Proxy for better LoadBalancing / Sentinel Master support.
## Automatically proxies to Redis master.
## Recommend for externally exposed Redis clusters.
## ref: https://cbonte.github.io/haproxy-dconv/1.9/intro.html
haproxy:
  enabled: true
  replicas: 3
```

Given these values, the helm chart can be deployed with:
```console
helm upgrade --install redis-ha stable/redis-ha --namespace magma --values=vals-redis-ha.yml
```

### CWAG Configuration

For each CWF helm release, there is a section of values that need to be updated
to use the redis server deployed above. In the `vals.yml` for each cwf release,
`port` should be changed to the port `vals-redis-ha.yml` and `bind` should be
changed to the `haproxy` svc name that was created.

```console
cwf:
  ...
  ...
  redis:
    port: 6380
    bind: redis-ha-haproxy
```

To enable stateless sessiond, on each gateway, modify the file at
`/etc/magma/sessiond.yml`. The field `support_stateless: false` should
be set to `true`.

To enable stateless RADIUS server, on each gateway, the file at
`/var/opt/magma/docker/docker-compose.yml` should be updated to properly
set the redis server:

```console
radius:
    ...
    ...
    environment:
      - STORAGE_TYPE=redis
      - REDIS_ADDR=redis-ha-haproxy:6380
```

**Note**: The sessiond and radius configuration changes will be incorporated
into the cwf helm chart in the future.

## Deploying CWF Operator

To deploy the operator, copy the deploy directory at
`magma/cwf/k8s/cwf_operator/deploy` to a controller node.
Then run the following on the node:

```console
$ kubectl -n magma create -f deploy/crds/magma.cwf.k8s_haclusters_crd.yaml
$ kubectl -n magma create -f deploy/service_account.yaml
$ kubectl -n magma create -f deploy/role.yaml
$ kubectl -n magma create -f deploy/role_binding.yaml
```

Before creating the operator pod, the `operator.yaml` file will need to be
modified with the appropriate `imagePullSecrets` and `image` fields. These
fields should correspond to the Docker registry that contains the
`cwf_operator` image (e.g. `image: docker.io/cwf_operator:latest`).

Now, create the operator pod:
```console
$ kubectl -n magma create -f deploy/operator.yaml
```

After this, `kubectl -n magma get pods` should display the operator running:


```console
cwf-operator-6b568c447d-hvgg8                 1/1     Running            0          2m45s
cwf01-8597b79ff8-h4csd                        1/1     Running            0          19d
cwf02-747b5bf75b-4njp9                        1/1     Running            0          19d
```

Lastly, we need to modify
`deploy/crds/magma.cwf.k8s_v1alpha1_hacluster_cr.yaml` to track our intended
CWAGs. To do this, modify `gatewayResourceNames` to the name of the 2 CWAGs
that will be configured in the active/standby cluster. The names should be
the **helm release name** of the deployed cwf helm charts. If you are unsure
what these should be, run `helm ls` to check the list of releases.

After making this change, the yaml file should something like:
```console
apiVersion: magma.cwf.k8s/v1alpha1
kind: HACluster
metadata:
  name: cwf-hacluster
spec:
  gatewayResourceNames:
   - "cwf01"
   - "cwf02"
```

Create this custom resource by running:
```console
$ kubectl -n magma create -f deploy/crds/magma.cwf.k8s_v1alpha1_hacluster_cr.yaml
```

To verify that the operator is configured properly,
`kubectl -n magma logs -f <cwf_operator_pod>` should show:
```console
I0708 08:21:48.704915       1 main.go:65] cmd "level"=0 "msg"="Operator Version: 0.0.1"
I0708 08:21:48.705194       1 main.go:66] cmd "level"=0 "msg"="Go Version: go1.13.4"
I0708 08:21:48.705206       1 main.go:67] cmd "level"=0 "msg"="Go OS/Arch: linux/amd64"
I0708 08:21:48.705215       1 main.go:68] cmd "level"=0 "msg"="Version of operator-sdk: v0.16.0"
I0708 08:21:48.705418       1 leader.go:46] leader "level"=0 "msg"="Trying to become the leader."
I0708 08:21:51.852198       1 leader.go:88] leader "level"=0 "msg"="No pre-existing lock was found."
I0708 08:21:51.858843       1 leader.go:108] leader "level"=0 "msg"="Became the leader."
I0708 08:21:53.861893       1 listener.go:44] controller-runtime/metrics "level"=0 "msg"="metrics server is starting to listen"  "addr"="0.0.0.0:8383"
I0708 08:21:53.862047       1 main.go:114] cmd "level"=0 "msg"="Registering Components."
I0708 08:21:53.862233       1 controller.go:92] helm.controller "level"=0 "msg"="Watching resource"  "apiVersion"={"Group":"charts.helm.k8s.io","Version":"v1alpha1"} "kind"="Cwf" "namespace"="" "reconcilePeriod"="10s"
I0708 08:22:01.090487       1 metrics.go:97] metrics "level"=0 "msg"="Metrics Service object created"  "Service.Name"="cwf-operator-metrics" "Service.Namespace"="magma"
I0708 08:22:07.190899       1 main.go:131] cmd "level"=0 "msg"="Starting the Cmd."
I0708 08:22:07.191324       1 controller.go:164] controller-runtime/controller "level"=0 "msg"="Starting EventSource"  "controller"="hacluster-controller" "source"={"Type":{"metadata":{"creationTimestamp":null},"spec":{"gatewayResourceNames":null},"status":{"active":"","activeInitState":"","standbyInitState":""}}}
I0708 08:22:07.191328       1 controller.go:164] controller-runtime/controller "level"=0 "msg"="Starting EventSource"  "controller"="cwf-controller" "source"={"Type":{"apiVersion":"charts.helm.k8s.io/v1alpha1","kind":"Cwf"}}
I0708 08:22:07.291680       1 controller.go:171] controller-runtime/controller "level"=0 "msg"="Starting Controller"  "controller"="hacluster-controller"
I0708 08:22:07.291717       1 controller.go:190] controller-runtime/controller "level"=0 "msg"="Starting workers"  "controller"="hacluster-controller" "worker count"=1
I0708 08:22:07.291684       1 controller.go:171] controller-runtime/controller "level"=0 "msg"="Starting Controller"  "controller"="cwf-controller"
I0708 08:22:07.291737       1 controller.go:190] controller-runtime/controller "level"=0 "msg"="Starting workers"  "controller"="cwf-controller" "worker count"=1
I0708 08:46:05.673708       1 hacluster_controller.go:112] controller_hacluster "level"=0 "msg"="Reconciling Cluster" "Request.Name"="cwf-hacluster" "Request.Namespace"="magma"
I0708 08:46:05.673747       1 hacluster_controller.go:131] controller_hacluster "level"=0 "msg"="No active is currently set. Setting active" "Request.Name"="cwf-hacluster" "Request.Namespace"="magma" "gateway"="cwf02"
I0708 08:46:05.699199       1 hacluster_controller.go:112] controller_hacluster "level"=0 "msg"="Reconciling Cluster" "Request.Name"="cwf-hacluster" "Request.Namespace"="magma"
I0708 08:46:05.902238       1 hacluster_controller.go:146] controller_hacluster "level"=0 "msg"="Fetched active health status" "Request.Name"="cwf-hacluster" "Request.Namespace"="magma" "health"="HEALTHY" "message"="gateway status appears healthy"
I0708 08:46:05.904689       1 hacluster_controller.go:152] controller_hacluster "level"=0 "msg"="Fetched standby health status" "Request.Name"="cwf-hacluster" "Request.Namespace"="magma" "health"="HEALTHY" "message"="gateway status appears healthy"
```