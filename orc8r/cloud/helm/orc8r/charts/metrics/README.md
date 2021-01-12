# Thanos Helm Deployment

Thanos provides scaling and data retention improvements to the prometheus metrics deployment.
It can be deployed by following the instructions below.

## Set up AWS infrastructure

Thanos is configured to use an S3 bucket to store long-term data, while only keeping a small amount of data in the Prometheus server. To support this, you must have an accessible S3 bucket
created.

### Create S3 Bucket and IAM User

- Create a bucket here: https://s3.console.aws.amazon.com/s3/home

  - Make sure it is in the same region as your orchestrator deployment, and disable public access.

- Create an IAM user that has full s3 access permissions
  - Save the access code and secret key

> NOTE: You can also reuse and existing IAM user as long as they have s3 access.

## Setup Helm

You will need to set some values and regenerate secret configs to enable thanos.

### Download Thanos Chart

```
$ cd MAGMA_ROOT/orc8r/cloud/helm/orc8r/charts/metrics
$ helm dep update
```

### Set values file

Add the following section to your values file:

```
metrics:
  thanos:
    enabled: true
    objstore:
      config:
        bucket: <S3-Bucket-Name>
        endpoint: s3.<region-code>.amazonaws.com
        region: <region-code>
        access_key: <ACCESS_KEY>
        secret_key: <SECRET_KEY>
```

### Generate and configure secrets

Similar to the instructions at TODO Get the right URL HERE!!! https://magma.github.io/magma/docs/orc8r/deploy_intro, configure the necessary secrets but this time add a value to enable Thanos configurations:

```
helm template orc8r-secrets charts/secrets \
    --namespace magma \
    --set-string secret.certs.enabled=true \
    --set-file secret.certs.files."rootCA\.pem"=charts/secrets/.secrets/certs/rootCA.pem \
    --set-file secret.certs.files."controller\.crt"=charts/secrets/.secrets/certs/controller.crt \
    --set-file secret.certs.files."controller\.key"=charts/secrets/.secrets/certs/controller.key \
    --set-file secret.certs.files."admin_operator\.pem"=charts/secrets/.secrets/certs/admin_operator.pem \
    --set-file secret.certs.files."admin_operator\.key\.pem"=charts/secrets/.secrets/certs/admin_operator.key.pem \
    --set-file secret.certs.files."certifier\.pem"=charts/secrets/.secrets/certs/certifier.pem \
    --set-file secret.certs.files."certifier\.key"=charts/secrets/.secrets/certs/certifier.key \
    --set-file secret.certs.files."nms_nginx\.pem"=charts/secrets/.secrets/certs/nms_nginx.pem \
    --set-file secret.certs.files."nms_nginx\.key\.pem"=charts/secrets/.secrets/certs/nms_nginx.key \
    --set-thanos_enabled=true \
    --set=docker.registry=$DOCKER_REGISTRY \
    --set=docker.username=$DOCKER_USERNAME \
    --set=docker.password=$DOCKER_PASSWORD |
    kubectl apply -f -
```

### Redeploy with the thanos-enabled values file

Follow instructions at <INSERT URL> and make sure to use the correct values file. If this succeeds your deployment should look like

```
NAME                                             READY   STATUS    RESTARTS   AGE
mysql-57955549d5-jvs4t                           1/1     Running   0          6m33s
nms-magmalte-6cdb5dc7f-qtlt8                     1/1     Running   0          8m45s
nms-nginx-proxy-5b86f479f7-8v59t                 1/1     Running   0          8m45s
orc8r-alertmanager-57d5d6ccc4-nkp8q              1/1     Running   0          8m45s
orc8r-alertmanager-configurer-76cf8f8f57-mgr7h   1/1     Running   0          8m45s
orc8r-controller-76948bc84-s67nd                 1/1     Running   0          8m44s
orc8r-nginx-7d6c78647-98zb7                      1/1     Running   0          8m45s
orc8r-prometheus-649ffd7ccb-hhjj2                2/2     Running   0          8m45s
orc8r-prometheus-cache-6d647df4d9-wqnnk          1/1     Running   0          8m45s
orc8r-prometheus-configurer-d474d69cc-rwrvz      2/2     Running   0          8m44s
orc8r-thanos-compact-cd995b9fd-wjcb9             1/1     Running   0          8m45s
orc8r-thanos-query-5c85d57886-dfpxb              1/1     Running   0          8m44s
orc8r-thanos-store-0-7479bf59f6-54mg5            1/1     Running   0          8m45s
orc8r-user-grafana-6498bb6959-lhmpj              1/1     Running   0          8m45s
postgresql-0                                     1/1     Running   0          22m
```

No additional work is required and monitoring and alerting will work as usual, but now Thanos
controls querying and long-term data storage in S3. Since the prometheus server retention time is now set to 6h, you can
be certain that thanos querying is working once you see data that is older than 6 hours from when you deployed. Additionally,
you can check the contents of your s3 bucket to see if prometheus data is making it there (will happen approximately every 2 hours). The bucket will look like this eventually:

```
aws s3 ls s3://<bucket-name> | head -n 8
                           PRE 01EN12K3EAXRDN1E3JDYXGF8FN/
                           PRE 01EN12N4W5GM8GPHJ3QCCJ3BKD/
                           PRE 01EN75YA9D51XHMKG5CM27PBFQ/
                           PRE 01EN9Z8JZZGW8GSK8SCB2NSHAA/
                           PRE 01ENAHV7V9M3AVD2XBY6SQ8RT3/
                           PRE 01ENAKTEP36W3QZSM7YX3MW3NS/
                           PRE 01ENARPZ7VHD106BCGDE994HAN/
                           PRE 01ENAZJPMZJAXXKH14B14TS0AX/
```
