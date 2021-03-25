---
id: debug_logs
title: View Logs
hide_title: true
---

# View Logs

This document describes how to view Orc8r logs in production deployments.

## Set log verbosity

Prerequisite: choose target Orc8r pod

```sh
export ORC_POD=$(kubectl --namespace orc8r get pod -l app.kubernetes.io/component=orchestrator -o jsonpath='{.items[0].metadata.name}')
```

Set pod [log verbosity](https://pkg.go.dev/github.com/golang/glog#V)

```sh
# Full verbosity
kubectl --namespace orc8r exec -it ${ORC_POD} -- /var/opt/magma/bin/service303_cli log_verbosity 10 obsidian

# Default verbosity
kubectl --namespace orc8r exec -it ${ORC_POD} -- /var/opt/magma/bin/service303_cli log_verbosity 0 obsidian
```

## Tail logs

Prerequisite: install [`stern`](https://github.com/wercker/stern)

```sh
brew install stern  # or alternative installer for your host machine
```

Tail logs from target pods

```sh
# All pods in a deployment
stern --namespace orc8r orc8r-state --since 2m

# All pods in multiple deployments
stern --namespace orc8r 'orc8r-(state|streamer|certifier|obsidian)' --since 2m

# All Orc8r application pods
stern --namespace orc8r --selector app.kubernetes.io/name=orc8r --exclude-container nginx --since 2m

# All NMS application pods
stern --namespace orc8r --selector app.kubernetes.io/name=nms --exclude-container nginx --since 2m
```

## Kibana

Orc8r can also forward its own logs, and aggregated gateway logs, to a configured [Elasticsearch](https://www.elastic.co/what-is/elasticsearch) endpoint. The logs can then be made available for consumption via [Kibana](https://www.elastic.co/what-is/kibana).

