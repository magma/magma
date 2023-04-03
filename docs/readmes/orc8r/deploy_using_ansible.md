---
id: deploy_using_ansible
title: Deploy Orchestrator using Ansible (Beta)
hide_title: true
---

# Deploy Orchestrator using Ansible (Beta)

This how-to guide can be used to deploy Magma's Orchestrator on any cloud environment. 
It contains roles to set up a Kubernetes cluster and deploy Magma Orchestrator using helm charts.
For more information on Magma Deployer, please visit the project's
[magma-deployer](https://github.com/magma/magma-deployer).

> magma-deployer is in Beta and is not yet production ready or feature complete.

## Pre-requisites

- Ubuntu Jammy 22.04 VM / Baremetal machine 
- RAM: 8GB
- CPU: 4 cores
- Storage: 100GB

## Deploy Orchestrator

Quick Install:

```
sudo bash -c "$(curl -sL https://github.com/magma/magma-deployer/raw/main/deploy-orc8r.sh)"
```

Following roles will be installed:

```
Sunday 02 April 2023  10:22:34 +0530 (0:00:00.044)       0:08:41.557 ********** 
=============================================================================== 
kubernetes ------------------------------------------------------------ 197.79s
orc8r ----------------------------------------------------------------- 141.46s
prerequisites ---------------------------------------------------------- 99.03s
docker ----------------------------------------------------------------- 41.85s
secrets ---------------------------------------------------------------- 11.68s
openebs ----------------------------------------------------------------- 8.10s
fluentd ----------------------------------------------------------------- 4.20s
postgresql -------------------------------------------------------------- 3.82s
metallb ----------------------------------------------------------------- 3.70s
haproxy ----------------------------------------------------------------- 2.99s
prometheus_cache_cleanup ------------------------------------------------ 2.61s
elasticsearch ----------------------------------------------------------- 2.57s
gather_facts ------------------------------------------------------------ 1.66s
dns --------------------------------------------------------------------- 0.04s
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ 
total ----------------------------------------------------------------- 521.51s
```

Switch to `magma` user after deployment has finished:

```
sudo su - magma
```

Check if all pods are ready:

```
kubectl get pods
```

```
NAME                                            READY   STATUS    RESTARTS   AGE
elasticsearch-master-0                          1/1     Running   0          10m
fluentd-7b5ffff8f8-5rr8c                        1/1     Running   0          10m
haproxy-f9c95678b-dndcj                         1/1     Running   0          10m
nms-magmalte-56495b6ff4-668kn                   1/1     Running   0          10m
nms-nginx-proxy-65bc67cd44-9pbrh                1/1     Running   0          10m
orc8r-accessd-7fcd7dc9b7-gg9tz                  1/1     Running   0          10m
orc8r-alertmanager-67c59fc8fd-bh9hx             1/1     Running   0          10m
orc8r-alertmanager-configurer-b47b95b69-m7vkc   1/1     Running   0          10m
orc8r-analytics-848557ccdf-jwr4f                1/1     Running   0          10m
orc8r-base-acct-7f7f6c5577-9667d                1/1     Running   0          10m
orc8r-bootstrapper-6d44bb55b-cpntw              1/1     Running   0          10m
orc8r-certifier-65fdd4776b-qxpq9                1/1     Running   0          10m
orc8r-configurator-76df9f4b9b-bxqp9             1/1     Running   0          10m
orc8r-ctraced-5c7f5496cc-5p2l7                  1/1     Running   0          10m
orc8r-cwf-6bc574bdbd-dqqx8                      1/1     Running   0          10m
orc8r-device-6cf589fc5d-7vrcz                   1/1     Running   0          10m
orc8r-directoryd-54d6975897-9j28t               1/1     Running   0          10m
orc8r-dispatcher-6b794c95b9-fccxn               1/1     Running   0          10m
orc8r-eventd-58b6dd8c5-8zjgh                    1/1     Running   0          10m
orc8r-feg-7dffff9cb5-zx7c7                      1/1     Running   0          10m
orc8r-feg-relay-65c8c68f6-9974l                 1/1     Running   0          10m
orc8r-ha-5d8d565b6f-g5m9v                       1/1     Running   0          10m
orc8r-health-54f67d778c-mfzfn                   1/1     Running   0          10m
orc8r-lte-658b4ff8fc-4npc5                      1/1     Running   0          10m
orc8r-metricsd-7f56b47d98-llwcr                 1/1     Running   0          10m
orc8r-nginx-7fb49b4489-h545t                    1/1     Running   0          10m
orc8r-nprobe-5d86b4f99f-m9424                   1/1     Running   0          10m
orc8r-obsidian-7bfd89d4fb-ndq69                 1/1     Running   0          10m
orc8r-orc8r-worker-545d669cc4-6k72z             1/1     Running   0          10m
orc8r-orchestrator-686cf7bc6f-8k4q5             1/1     Running   0          10m
orc8r-policydb-5d584fd576-mr8kp                 1/1     Running   0          10m
orc8r-prometheus-6d77968679-5pf9t               1/1     Running   0          10m
orc8r-prometheus-cache-5d679b8847-gnqz5         1/1     Running   0          10m
orc8r-prometheus-configurer-5d54d6556-sdfsr     1/1     Running   0          10m
orc8r-service-registry-6755f7f8f6-hznwq         1/1     Running   0          10m
orc8r-smsd-659f9d4c4d-7qz9t                     1/1     Running   0          10m
orc8r-state-97bb66f49-mjhhz                     1/1     Running   0          10m
orc8r-streamer-7744c46486-pbw5l                 1/1     Running   0          10m
orc8r-subscriberdb-5cdc7c599d-qv9sp             1/1     Running   0          10m
orc8r-subscriberdb-cache-7d7d5fff78-65fpl       1/1     Running   0          10m
orc8r-tenants-6cd4888466-hsh5r                  1/1     Running   0          10m
orc8r-user-grafana-597dfff79-cnzdx              1/1     Running   0          10m
postgresql-0                                    1/1     Running   0          12m
```

Now setup NMS login:

```
cd ~/magma-deployer
ansible-playbook config-orc8r.yml
```

## DNS Setup

Get the External IP address of the `haproxy` service:

```
kubectl get svc haproxy
```

```
NAME      TYPE           CLUSTER-IP      EXTERNAL-IP     PORT(S)                                     AGE
haproxy   LoadBalancer   10.43.177.100   10.86.113.153   80:31192/TCP,443:32665/TCP,1024:30097/TCP   10m
```

Update `/etc/hosts` file with the following entries:

```
10.86.113.153 api.magma.local
10.86.113.153 magma-test.nms.magma.local
10.86.113.153 fluentd.magma.local
10.86.113.153 controller.magma.local
10.86.113.153 bootstrapper-controller.magma.local
```

> Replace the External IP with the one you got from the previous step.

You can access NMS dashboard at the following URL:

https://magma-test.nms.magma.local


## Access Gateway Setup

You can get your `rootCA.pem` file from the following location for connecting your Access Gateway:

```
cat ~/magma-deployer/secrets/rootCA.pem
```

Update `/var/opt/magma/configs/control_proxy.yml` file with the following content:

```
cloud_address: controller.magma.local
cloud_port: 443
bootstrap_address: bootstrapper-controller.magma.local
bootstrap_port: 443
fluentd_address: fluentd.magma.local
fluentd_port: 443

rootca_cert: /var/opt/magma/tmp/certs/rootCA.pem
```

Update `/etc/hosts` file in Access Gateway with the following entries:

```
10.86.113.153 fluentd.magma.local
10.86.113.153 controller.magma.local
10.86.113.153 bootstrapper-controller.magma.local
```
