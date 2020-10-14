CWF Kubernetes Operator
===

The CWF Operator is a Kubernetes Operator used to deploy the CWF Helm
chart (located at `magma/cwf/gateway/helm`) in a highly available
configuration. The only configuration currently supported is 2 CWAGs
running in active/standby mode.

Note: The `cwf-virtlet` helm chart is deprecated. If deploying a new
gateway, use the chart located at `magma/cwf/gateway/helm/cwf-kubevirt`.

## Prerequisites

Before installing the CWF Operator, the `redis-ha` helm chart must be deployed:

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

### CWAG Helm Configuration

For each CWF helm release, there is a section of values that need to be
updated. In the `vals.yml` override for each cwf release, set `gateway_ha`
to `enabled`.

Then set the redis section to the port and bind address to match the redis-ha
deployment. `port` should be changed to the port in `vals-redis-ha.yml` an
d
`bind` should be changed to the `haproxy` svc name that was created.

```console
cwf:
  ...
  gateway_ha:
    enabled: false
  ...
  redis:
    port: 6380
    bind: redis-ha-haproxy
```

Once these values are set, deploy the CWF helm chart with these
values
 overrides specified.

## Gateway Configuration and Registration

Once the gateway pods have fully spun up, check to ensure that the installation
succeeded for both gateways. To do this:

```console
[~/] kubectl -n magma get VirtualMachineInstances
[~/] sudo virtctl -n magma console <cwfKubeVirtInstance>
[user@cwf01:~/] cd /var/opt/magma/docker
[user@cwf01:/var/opt/magma/docker] sudo docker ps
```

You should see all of the containers running:
```console
2bff32bb9787        docker.io/gateway_sessiond:tag    "/usr/local/bin/sess…"   1 hour ago        Up 1 hour (healthy)                       sessiond
912626ef88fa        docker.io/gateway_python:tag      "python3.5 -m magma.…"   1 hour ago        Up 1 hour                                 state
fc577508d4c4        docker.io/gateway_pipelined:tag   "python3.5 -m magma.…"   1 hour ago        Up 1 hour (healthy)                       policydb
441b92fec477        docker.io/gateway_python:tag      "python3.5 -m magma.…"   1 hour ago        Up 1 hour (healthy)                       directoryd
1c017b67b50f        docker.io/gateway_python:tag      "python3.5 -m magma.…"   1 hour ago        Up 1 hour                                 magmad
9d118131385e        docker.io/gateway_go:tag          "envdir /var/opt/mag…"   1 hour ago        Up 1 hour (healthy)                       eap_aka
0926c3c6613a        docker.io/gateway_python:tag      "/bin/bash -c '/usr/…"   1 hour ago        Up 1 hour                                 td-agent-bit
6f2c3c82b723        docker.io/gateway_go:tag          "envdir /var/opt/mag…"   1 hour ago        Up 1 hour (healthy)                       aaa_server
cd704bed6e75        docker.io/cwag_go:tag             "envdir /var/opt/mag…"   1 hour ago        Up 1 hour                                 health
f12eb89a7cbd        docker.io/gateway_python:tag      "sh -c '/usr/local/b…"   1 hour ago        Up 1 hour                                 control_proxy
3e806aa645ec        docker.io/gateway_python:tag      "python3.5 -m magma.…"   1 hour ago        Up 1 hour                                 eventd
f786fd8f9d6c        docker.io/gateway_go:tag          "envdir /var/opt/mag…"   1 hour ago        Up 1 hour                                 radiusd
beed2e197c13        docker.io/gateway_pipelined:tag   "sh -c 'set bridge c…"   1 hour ago        Up 1 hour (healthy)                       pipelined
e6f78db3d8ad        docker.io/gateway_go:tag          "/bin/bash -c 'envsu…"   1 hour ago        Up 1 hour (healthy)                       radius
9a097700e3a8        docker.io/gateway_python:tag      "/bin/bash -c '/usr/…"   1 hour ago        Up 1 hour                                 redis
```

Then register each gateway:
```console
[user@cwf01:~/] cd /var/opt/magma/docker
[user@cwf01:/var/opt/magma/docker] sudo docker-compose exec magmad /usr/local/bin/show_gateway_info.py
```

You should see output similar to:

```console
Hardware ID:
------------
d8ae2bd1-a517-4f48-a01d-648fdfcd28da

Challenge Key:
-----------
MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAEzQeuSAA1Rmb0UyCE87DuX1HiFJY3McKm2XqplyJIsiPkqg8HNQ7v1MZwnvVI28AF/JAmRiQXjJqWqIqL1n3IsvOa9XyqOvIruCSVPo8wUaVkbs5bx9nFkzUMthuBjVjY
```

Register the gateway by adding this information to the intended CWF network
in your NMS.

### Creating the HA Pair

Lastly, we need to configure an `HAPair` network entity to associate these
two gateways into an HA cluster. NMS support for this is not yet completed,
so this will have to be done through the Orchestrator API.

Navigate to the API, and find endpoint `/cwf/{network_id}/ha_pairs`. Use the
`POST` method to create the pair. An important part of this configuration is
the `transport_virtual_ip` field. This is the IP address that will be used by
WLC/APs to connect to the CWAG cluster. The IP should be an unused IP in the
SAME subnet as eth1 of the gateways.

As an example, below is the pair we would create if our two gateways had
gateway ID's `gw1` and `gw2`:

```console
{
  "config": {
    "transport_virtual_ip": "10.10.10.12/24"
  },
  "gateway_id_1": "gw1",
  "gateway_id_2": "gw2",
  "ha_pair_id": "cwf_pair1"
}
```

## Deploying CWF Operator

Lastly, we need to deploy the operator to manage this cluster. To do so,
copy the deploy directory at `magma/cwf/k8s/cwf_operator/deploy` to a
controller node.

Then run the following on the node:

```console
$ kubectl -n magma create -f deploy/crds/magma.cwf.k8s_haclusters_crd.yaml
$ kubectl -n magma create -f deploy/service_account.yaml
$ kubectl -n magma create -f deploy/role.yaml
$ kubectl -n magma create -f deploy/role_binding.yaml
```

Before creating the operator pod, the `operator.yaml` file will need to be
modified. The following fields should be updated:
 * `imagePullSecrets` - update to the correct secrets to pull the operator
 image
 * `image` - update to the correct image (e.g. `docker.io/operator:latest`)
 * `REDIS_ADDR` - update to the deployed redis addr for the redis-ha helm chart
 (e.g. `redis-ha-haproxy:6380`)

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
CWAGs. To do this, modify `gatewayResources` to the name of 2 CWAG resources
that will be configured in the active/standby cluster. For each resource,
define:
 * `gatewayID` - the gateway ID created in the NMS
 * `helmReleaseName` - the release name of the gateway helm deployment

If you are unsure what the helm release name should be, run `helm ls -n magma`
to check the list of releases.

After making this change, the yaml file should something like:
```console
apiVersion: magma.cwf.k8s/v1alpha1
kind: HACluster
metadata:
  name: example-hacluster
spec:
  gatewayResources:
    - gatewayID: "gw1"
      helmReleaseName: "cwf01"
    - gatewayID: "gw2"
      helmReleaseName: "cwf02"
  haPairID: "cwf_pair1"
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

**Note**: If creating more than 1 HA cluster, the operator need only be
deployed once. Just define a new HACluster resource
(a new `magma.cwf.k8s_v1alpha1_hacluster_cr.yaml` file) set to the
appropriate gateways and run `kubectl -n magma create -f <new_cluster_cr.yaml>`.