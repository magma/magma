# [WIP] Scaled Prometheus Pipeline

**Author: Scott Smith**

**Feature Owner: Scott Smith**

**Requires Feedback From: Jacky Tian**

Details:
* I use the term pushgateway to reference [prometheus-edge-hub](https://github.com/facebookincubator/prometheus-edge-hub)
* "Prod" and "Staging" refer to FB-hosted deployments of Orchestrator which are currently our largest deployments and our source of performance data.

### Goal: Improve metrics query speed on large deployments

The current magma metrics deployment involves only a single Prometheus server, with no option for scaling. While this is sufficient for writing, the query side is already seeing some performance issues. For example, loading grafana on prod takes more than 5 seconds. This is due to the queries grafana has to make to populate the list of networkIDs and gatewayIDs being slow. While the production deployment is currently the largest that exists, we don’t want to have partners run into serious metrics performance issues as they build large networks. We need to scale prometheus to improve query performance.

## Data

### Prometheus

On prod we currently have ~347 active gateways, and ~105k metric series

Query times for PromQL query `count({__name__=~".+"})` 

|Time Range	|Query Time	|
|---	|---	|
|1m	|3141ms	|
|5m	|7801ms	|
|15m	|6966ms	|
|30m	|6025ms	|
|1h	|6153ms	|
|2h	|8186ms	|
|6h	|11957ms	|
|12h	|17983ms	|
|1d	|TIMEOUT	|

Even the smallest query range takes over 3 seconds, and over any significant period of time the query times out. This is definitely not “snappy” and will not provide a good user experience with a large network. This data fits in with prometheus recommendations that queries should not have more than ~100k series.

### Grafana

The main problem with long query times is less responsive dashboards.

Time to load grafana dashboard on prod NMS: **>15s **from click until all graphs are populated. Most of this time is spent looking at all series in order to determine the set of available `networkID`s and `gatewayID`s. 

## Options

Several scaling solutions exist for prometheus, as this is a well-known limitation and prometheus itself does not have any built in way of horizontal scaling.

### Thanos

[Thanos](https://improbable.io/blog/thanos-prometheus-at-scale) is a very popular project which allows for easy and customizable scaling of prometheus monitoring pipelines. From the start, I believe this will be the easiest and most powerful solution to the problem. Thanos consists of several components, all of which can be used independently. The most relevant to us is the `Querier` which allows for the querying of data across multiple prometheus servers. A simple architecture diagram from Thanos shows how this works in a typical deployment:
![image.png](./image.png)Here we see multiple prometheus servers with the Thanos `sidecar` attached. This allows for the rest of the thanos components to work together. Then, the `Querier` components are able to accept PromQL queries and retrieve data from any set of the prometheus servers. Here the queriers are replicated, but that is only for HA purposes and does not improve scalability as far as I know. Additionally there are components to enable long term storage, but we won’t need that for this task. A simpler diagram for our use case would look like:

![scaledprometheus1.png](./scaledprometheus1.png)

With this setup, we only need to deploy the Thanos `sidecar` and the `Querier`, along with multiple prometheus servers to achieve horizontal scalability. I’ll discuss the details, specifically how to push metrics to multiple prometheus servers later.

### Cortex

Another option for scaling prometheus is [Cortex](https://cortexmetrics.io/). Similar to Thanos, but more complicated. In a Cortex deployment, we would deploy multiple prometheus servers and configure them to write to Cortex, and deploy the Cortex querier. 


![cortex](./cortex.png)
One major downside of Cortex is that it [requires](https://cortexmetrics.io/docs/production/running-in-production/) a long-term storage option. This makes it unlikely to be a good candidate for this since it would increase the cost of deploying and require us to add LTS (which is not a requirement at the moment) instead of just focusing on improving query performance.

## A note on multi-tenancy

Both of these projects support some form of multi-tenancy, however I don’t believe they are useful for our situation since they tie the tenancy to a specific prometheus server. This means we would have to have a 1:1 mapping from servers to tenants, which may not match up with our scaling needs. For example there could be a deployment with one massive network, and several smaller ones. Queries on the big network would not be improved this way.

## Configmanager and multiple prometheus servers

One problem that arises with multiple prometheus servers is handling alerting configurations. Alerting rules are used by the prometheus server to evaluate rules against incoming metrics, but in a model with multiple servers this requires a lot of careful planning and tradeoffs due to metrics being split across servers, but rules only applying to a single server at a time.

Thanos provides a solution to this problem in the [Ruler](https://thanos.io/components/rule.md/) component. It accepts rules files in the same format as prometheus, and evaluates the rules using a Thanos Querier which has access to all servers. Feature parity not documented (POST request to reload configuration files) validated by investigating [source code](https://github.com/thanos-io/thanos/blob/master/cmd/thanos/rule.go#L567).

To keep our alerting system all we have to do is configure the Thanos Ruler to use the alerting rules files managed by prometheus-configmanager, and remove the rules files from the individual prometheus servers.

Note: Alertmanager is not affected by the switch to Thanos, and there will still only be a single alertmanager instance.

## Deploying Thanos in Orchestrator

### Scaling the Pushgateway

All metrics go through the controller (which is already scalable and placed behind a load balancer), where they are pushed to the `prometheus-edge-hub` (or pushgateway). The single prometheus server then scrapes from this pushgateway. 
Current metrics pipeline diagram:

![currentMetricsPipeline.png](./currentMetricsPipeline.png)
Since the pushgateway simply works as a temporary holder of metrics for prometheus to scrape, it does not currently work with a multiple-server model. If two prometheus servers scrape a pushgateway, there is no way to send only a subset of metrics to each server. This could be done of course, but it would add complexity to the pushgateway which I want to keep as simple as possible. We will need to deploy multiple pushgateways (one per prometheus server) and push a subset of metrics to each pushgateway. 

![newMetricsPipeline.png](./newMetricsPipeline.png)

To avoid making performance worse, we of course need to ensure that for each series, it is always sent to the same pushgateway and prometheus server. To do that, we need to implement some form of hashing. Preferably, this would all be done on the controllers, to avoid having to introduce a new “distributor” component. 

### Distributing Metrics

On prod, out of the ~105k metric series, only about 1.5k do not belong to a gateway. The average number of metric series per gateway is only ~350. This makes `gatewayID` a useful key to hash on. We simply map each `gatewayID` to a specific pushgateway, and all non-gateway metrics also to a single specific pushgateway. Even a very simple implementation with this method will ensure that no server has more than about 1000 series more than average. Importantly, all of this work can be done on the controller in the `Exporter`. 

**Discussion:**
While hashing on the controller provides an easy solution which is already scalable, it may be more architecturally sound to instead shift this responsibility to a load balancer in front of the pushgateways. Unfortunately, we cannot use a naive load-balancer since the metrics must be directed to the correct gateway. It would be possible to modify the pushgateway as a load balancer with this functionality. Simply configure it to hash the metrics as they come in and send them to the correct pushgateway to be scraped. However we then are left with another bottleneck to deal with. For this reason, I believe hashing on the controllers is the best option.

### Dealing with Rescaling

When a user decides to rescale their deployment (either up or down), series will inevitably be mapped to different pushgateways. This will mean that series will be scraped by different prometheus servers as well. Note that this will not cause any issues when querying/alerting as Thanos can handle the same series existing across multiple servers. The only effect this will have is some temporary non-optimal distribution of series since some will exist in two locations. Once the local server retention time has passed, everything will be back to normal, but even during that period we expect no noticeable impact to performance or user experience.

We will use consistent hashing on the exporters to ensure that as few series as possible are moved.

### Pushgateway Redundancy

We can leverage the addition of multiple gateways to provide failure protection. As discussed above, it's not a significant problem if some series get shifted to a different prometheus server. Because of this, we can protect metrics against a pushgateway failure by simply pushing metrics to the next gateway in the hash circle if the first request fails. While I do not expect this to happen much if at all (as far as I know a pushgateway has never failed while in prod/staging), this provides extra protection at no cost.

## Configuration

We should make this as easy to configure as possible. Some goals include:

* Make this entire setup optional. The single prometheus server works well for small to medium deployments.
* Configuring thanos to run in Orc8r should require no more than a few values in the helm chart. (e.g. `create_thanos` and `prometheus_servers_count)`

* Since we are moving the prometheus-configmanager to act on the Thanos Ruler component, we can manage prometheus server configuration with Kubernetes and ConfigMaps. This will be used to correctly set scrape targets for the servers with respect to the multiple pushgateways.

## Optional Object Storage

Thanos supports long-term storage and compaction via the `Store Gateway` and `Compactor` components. While the system can work just fine by using only local storage with the prometheus servers, the [object storage](https://thanos.io/storage.md/) makes it possible to store metrics for much longer without using much more storage space. Thanos supports object storage with S3, which is the obvious path for Orchestrator since production deployments are done on AWS. Additionally, for on-prem deployments OpenStack Swift is a supported option.

Adding object storage as an optional feature can be done after the first iteration of deploying Orchestrator with Thanos is done. 

## Development Plan

|Step	|Est. Time	|
|---	|---	|
|Deploy Thanos locally and experiment with loads to validate query time improvements	|2 wk	|
|Deploy Thanos with a single prometheus server in Orchestrator to test that the Querier and sidecar don’t cause any unexpected side effects.	|1 wk	|
|Integrate configmanager with Thanos and validate alerting/configuration	|1 wk	|
|Implement hashing in the prometheus-edge-hub exporter	|2 wk	|
|Deploy locally with multiple pushgateways/servers	|2 wk	|
|Configure helm chart (ideally so scaling is as simple as setting a single value)	|2 wk	|
|Deploy to staging/prod	|1 wk	|
|Add Object storage option	|2 wk	|


