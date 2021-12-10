---
id: version-1.5.0-thanos
title: Thanos
hide_title: true
original_id: thanos
---

# Scaled Metrics with Thanos

[Thanos](https://thanos.io/) is an open-sourced scaled prometheus solution used
to provide features like global query view, unlimited retention, and downsampling
and compaction. Orchestrator provides large deployments the option to deploy
with Thanos in order to have a more robust metrics pipeline.

## Deploying Orc8r with Thanos

The terraform module makes deploying Thanos very easy. In your `main.tf` you
just need to set the following values:

```
module orc8r {
    thanos_enabled = true
}

module orc8r-app {
  thanos_enabled = true
  thanos_object_store_bucket_name = "<globally-unique-bucket-name>"
}
```

Choose a value for `thanos_object_store_bucket_name` that will be globally
(across all of AWS) unique, but other than that the value doesn't matter.

That's all you need to do to deploy with Thanos! All interacting components
will be adjusted accordingly, so the NMS/Grafana will work the same as before.
If you don't care about the internals you can stop reading here.

If you run `kubectl --namespace orc8r get pods` it should now look like this:
```
NAME                                             READY   STATUS    RESTARTS   AGE
fluentd-6fb9f57dff-ljmfw                         1/1     Running   0          24h
fluentd-6fb9f57dff-p54p9                         1/1     Running   0          24h
nms-magmalte-f4bbf4cfb-tqblm                     1/1     Running   0          24h
nms-nginx-proxy-57b8585d6-4ml6s                  1/1     Running   0          24h
orc8r-alertmanager-84d79f774b-4svrs              1/1     Running   0          24h
orc8r-alertmanager-configurer-68d6c55c9c-6q9xg   1/1     Running   0          24h
orc8r-controller-7494c96646-4w4jp                1/1     Running   0          20h
orc8r-controller-7494c96646-7fcg5                1/1     Running   0          20h
orc8r-nginx-5f9d7f4bcc-cz5ld                     1/1     Running   0          20h
orc8r-nginx-5f9d7f4bcc-rszpr                     1/1     Running   0          20h
orc8r-prometheus-5bdd644fd8-mm8gb                2/2     Running   0          24h
orc8r-prometheus-cache-f84884575-7vw8d           1/1     Running   0          24h
orc8r-prometheus-configurer-69df67988-w9dc6      2/2     Running   0          20h
orc8r-thanos-compact-66dd4d974b-jwzjk            1/1     Running   0          20h
orc8r-thanos-query-5d5cb888bd-vm9t8              1/1     Running   0          114m
orc8r-thanos-store-0-7479bf59f6-97wbp            1/1     Running   0          114m
orc8r-user-grafana-bc644b4fc-28nmf               1/1     Running   0          24h
```

Notice that the prometheus pod now has another container running, this is the
[thanos sidecar](https://thanos.io/v0.17/components/sidecar.md/). There is another
sidecar that runs with `prometheus-configurer`, and then three more components
that run independently: compact, query, and store.

## Advanced configuration options

The default infrastructure setup deploys an additional node for Thanos, since
there is one component that requires significant on-node ephemeral storage.
However, you may want to deploy more nodes if you want to make sure thanos
components run on different nodes than the rest of orc8r. To do that you can
override the default value for `thanos_worker_groups` in the `orc8r` module.
The default value is:
```
[
    {
      name                 = "thanos-1"
      instance_type        = "m5d.xlarge"
      asg_desired_capacity = 1
      asg_min_size         = 1
      asg_max_size         = 1
      autoscaling_enabled  = false
      kubelet_extra_args = "--node-labels=compute-type=thanos"
    },
  ]
```

To add more workers, either adjust the `asg_...` values in that object, or add
another entry to that array of worker groups. To specify thanos components to
run on specific nodes, just set the following variables in the `orc8r-app` module:
```
thanos_query_node_selector = "thanos"
thanos_store_node_selector = "thanos"
```
> Note: set the value to the same value you used for `--node-labels=compute-type=<value>`
> in order to run on that worker group

These are advanced configuration options, and we don't expect them to be necessary,
but are available to give more fine-grained control over your deployment.
