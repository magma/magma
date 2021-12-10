---
id: version-v1.0.1-deploy_install
title: Installing Orchestrator
hide_title: true
original_id: deploy_install
---
# Installing Orchestrator

## Creating Secrets

IMPORTANT: in all the below instructions, replace `yourdomain.com` with the
actual domain/subdomain which you've chosen to host Orchestrator on.

We recommend storing the following secrets on S3 in AWS, or any similar
object storage service on your preferred cloud provider.

Start first by creating a new directory somewhere to hold the secrets while
you create them:

```bash
mkdir -p ~/secrets/
```

### SSL Certificates

You will need the following certificates and private keys:

1. The public SSL certificate for your Orchestrator domain,
with CN=*.yourdomain.com. This can be an SSL certificate chain, but it must be
in one file
2. The private key which corresponds to the above SSL certificate
3. The rootCA certificate which verifies your SSL certificate.

If you already have these files, you can do the following:

1. Rename your public SSL certificate to `controller.crt`
2. Rename your SSL certificate's private key to `controller.key`
3. Rename your SSL certificate's root CA to `rootCA.pem`
4. Put these 3 files under a subdirectory `certs`

If you aren't worried about a browser warning, you can also self-sign these
certs. Change the values of the DN prompts as necessary, but pay *very* close
attention to the common names - these are very important to get right!

```bash
$ mkdir -p ~/secrets/certs
$ cd ~/secrets/certs
$ openssl genrsa -out rootCA.key 2048

Generating RSA private key, 2048 bit long modulus

$ openssl req -x509 -new -nodes -key rootCA.key -sha256 -days 3650 -out rootCA.pem

You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) []:US
State or Province Name (full name) []:CA
Locality Name (eg, city) []:Menlo Park
Organization Name (eg, company) []:Facebook
Organizational Unit Name (eg, section) []:Magma
Common Name (eg, fully qualified host name) []:rootca.yourdomain.com
Email Address []:admin@yourdomain.com

$ openssl genrsa -out controller.key 2048

Generating RSA private key, 2048 bit long modulus

$ openssl req -new -key controller.key -out controller.csr

You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) []:US
State or Province Name (full name) []:CA
Locality Name (eg, city) []:Menlo Park
Organization Name (eg, company) []:Facebook
Organizational Unit Name (eg, section) []:Magma
Common Name (eg, fully qualified host name) []:*.yourdomain.com
Email Address []:admin@yourdomain.com

Please enter the following 'extra' attributes
to be sent with your certificate request
A challenge password []:

$ openssl x509 -req -in controller.csr -CA rootCA.pem -CAkey rootCA.key -CAcreateserial -out controller.crt -days 365 -sha256

Signature ok
subject=/C=US/ST=CA/L=Menlo Park/O=Facebook/OU=Magma/CN=*.yourdomain.com/emailAddress=admin@yourdomain.com
Getting CA Private Key

$ rm controller.csr rootCA.srl
```

At this point, regardless of whether you self-signed your certs or acquired
them from a certificate provider, your `certs` subdirectory should look like
this:

```bash
$ ls
controller.crt  controller.key  rootCA.pem    rootCA.key
```

We *strongly* recommend moving `rootCA.key` to a more secure location at this
point. By default, the Helm chart for secrets below will upload it to EKS
as a Kubernetes secret and mount it to controller pods. If the private key
portion of the root CA is compromised, all TLS traffic to and from your cluster
will be compromised.

Keep rootCA.key in a place where you can access it again - you will need it to
renew the SSL certificate for the Orchestrator controller when it expires.

### Application Certificates and Keys

`certifier` is the Orchestrator service which signs client certificates. All
access to Orchestrator is authenticated by client SSL certificates, so you'll
need to create the verifying certificate for `certifier`.

Again, pay *very* close attention to the CN.

```bash
$ openssl genrsa -out certifier.key 2048

Generating RSA private key, 2048 bit long modulus

$ openssl req -x509 -new -nodes -key certifier.key -sha256 -days 3650 -out certifier.pem

...
Common Name (eg, fully qualified host name) []:certifier.yourdomain.com

$ openssl genrsa -out bootstrapper.key 2048

Generating RSA private key, 2048 bit long modulus

$ ls
bootstrapper.key  certifier.key     certifier.pem     controller.crt    controller.key    rootCA.pem    rootCA.key
```

The last command created a private key for the `bootstrapper` service, which
is the mechanism by which Access Gateways acquire their client certificates
from `certifier`.

### Environment Secrets

Go into the AWS management console, choose "RDS", and find the hostname of your
orc8r RDS instance (make sure not to choose the NMS RDS instance). Note this
down, then continue:

```bash
mkdir -p ~/secrets/envdir
cd ~/secrets/envdir
echo "STREAMER,SUBSCRIBERDB,METRICSD,CERTIFIER,BOOTSTRAPPER,METERINGD_RECORDS,ACCESSD,OBSIDIAN,DISPATCHER,DIRECTORYD" > CONTROLLER_SERVICES
echo "dbname=orc8r user=orc8r password=<YOUR ORC8R DB PASSWORD> host=<YOUR ORC8R RDS ENDPOINT>" > DATABASE_SOURCE
```

### Static Configuration Files

Orchestrator microservices can be configured with static YAML files. In this
deployment, the only one you'll have to create will be for `metricsd`:

```bash
mkdir -p ~/secrets/configs/orc8r
cd ~/secrets/configs/orc8r
touch metricsd.yml
```

Put the following contents into `metricsd.yml`:

```
profile: "prometheus"

prometheusQueryAddress: "http://orc8r-prometheus:9090"
prometheusPushAddresses:
  - "http://orc8r-prometheus-cache:9091/metrics"

alertmanagerApiURL: "http://orc8r-alertmanager:9093/api/v2/alerts"
prometheusConfigServiceURL: "http://orc8r-config-manager:9100"
alertmanagerConfigServiceURL: "http://orc8r-config-manager:9101"
```

## Initial Helm Deploy

Copy your secrets into the Helm subchart where you cloned Magma:

```bash
cp -r ~/secrets magma/orc8r/cloud/helm/orc8r/charts/secrets/.secrets
```

We need to set up the EKS cluster before we can `helm deploy` to it:

```bash
cd magma/orc8r/cloud/helm/orc8r
kubectl apply -f tiller-rbac-config.yaml
helm init --service-account tiller --history-max 200
# Wait for tiller to become 'Running'
kubectl get pods -n kube-system | grep tiller

kubectl create namespace magma
```

Next, create a `vals.yml` somewhere in a source controlled directory that you
own (e.g. adjacent to your terraform scripts). Fill in the values in caps
with the correct values for your docker registry and Orchestrator hostname:

```
imagePullSecrets:
  - name: orc8r-secrets-registry

secrets:
  create: true
  docker:
    registry: YOUR-DOCKER-REGISTRY
    username: YOUR-DOCKER-USERNAME
    password: YOUR-DOCKER-PASSWORD
  

proxy:
  image:
    repository: YOUR-DOCKER-REGISTRY/proxy
    tag: YOUR-CONTAINER-TAG

  replicas: 2

  service:
    name: orc8r-bootstrap-legacy
    type: LoadBalancer

  spec:
    hostname: controller.YOURDOMAIN.COM

  nodeSelector:
    worker-type: controller

controller:
  image:
    repository: YOUR-DOCKER-REGISTRY/controller
    tag: YOUR-CONTAINER-TAG

  replicas: 2

  migration:
    new_handlers: 1
    new_mconfigs: 1

  nodeSelector:
    worker-type: controller

metrics:
  imagePullSecrets:
    - name: orc8r-secrets-registry

  metrics:
    volumes:
      prometheusData:
        volumeSpec:
          hostPath:
            path: /prometheusData
            type: DirectoryOrCreate
      prometheusConfig:
        volumeSpec:
          hostPath:
            path: /configs/prometheus
            type: DirectoryOrCreate

  prometheus:
    create: true
    nodeSelector:
      worker-type: metrics

  configmanager:
    create: true
    image:
      repository: YOUR-DOCKER-REGISTRY/config-manager
      tag: YOUR-CONTAINER-TAG
    nodeSelector:
      worker-type: metrics

  alertmanager: 
    create: true
    nodeSelector:
      worker-type: metrics

  prometheusCache:
    create: true
    image:
      repository: YOUR-DOCKER-REGISTRY/prometheus-cache
      tag: YOUR-CONTAINER-TAG
    limit: 500000
    nodeSelector:
      worker-type: metrics

  grafana:
    create: true
    image:
      repository: YOUR-DOCKER-REGISTRY/grafana
      tag: YOUR-CONTAINER-TAG
    nodeSelector:
      worker-type: metrics

nms:
  imagePullSecrets:
    - name: orc8r-secrets-registry

  magmalte:
    manifests:
      secrets: false
      deployment: false
      service: false
      rbac: false

    image:
      repository: YOUR-DOCKER-REGISTRY/magmalte
      tag: YOUR-CONTAINER-TAG

    env:
      api_host: controller.YOURDOMAIN.COM
      mysql_host: YOUR RDS MYSQL HOST
      mysql_user: magma
      mysql_pass: YOUR RDS MYSQL PASSWORD
  nginx:
    manifests:
      configmap: false
      secrets: false
      deployment: false
      service: false
      rbac: false

    service:
      type: LoadBalancer

    deployment:
      spec:
        ssl_cert_name: controller.crt
        ssl_cert_key_name: controller.key
```

NMS won't work without a client certificate, so we've turned off those
deployments for now. We'll create an admin cert and upgrade the deployment
with NMS once the core Orchestrator components are up.

At this point, if your `vals.yml` is good, you can do your first helm deploy:

```bash
cd magma/orc8r/cloud/helm/orc8r
helm install --name orc8r --namespace magma . --values=PATH_TO_VALS/vals.yml
```

## Creating an Admin User

First, find a `orc8r-controller-` pod in k8s:

```bash
$ kubectl -n magma get pods

NAME                                      READY   STATUS    RESTARTS   AGE
orc8r-configmanager-896d784bc-chqr7       1/1     Running   0          X
orc8r-controller-7757567bf5-cm4wn         1/1     Running   0          X
orc8r-controller-7757567bf5-jshpv         1/1     Running   0          X
orc8r-alertmanager-c8dc7cdb5-crzpl        1/1     Running   0          X
orc8r-grafana-6446b97885-ck6g8            1/1     Running   0          X
orc8r-prometheus-6c67bcc9d8-6lx22         1/1     Running   0          X
orc8r-prometheus-cache-6bf7648446-9t9hx   1/1     Running   0          X
orc8r-proxy-57cf989fcc-cg54z              1/1     Running   0          X
orc8r-proxy-57cf989fcc-xn2cw              1/1     Running   0          X
```

Then:

```bash
export CNTLR_POD=$(kubectl -n magma get pod -l app.kubernetes.io/component=controller -o jsonpath='{.items[0].metadata.name}')
kubectl exec -it ${CNTLR_POD} bash

# The following commands are to be run inside the pod
(pod)$ cd /var/opt/magma/bin
(pod)$ envdir /var/opt/magma/envdir ./accessc add-admin -cert admin_operator admin_operator
(pod)$ openssl pkcs12 -export -out admin_operator.pfx -inkey admin_operator.key.pem -in admin_operator.pem

Enter Export Password:
Verifying - Enter Export Password:

(pod)$ exit
```

Now on your host, copy down the client certificates for the admin operator we
just created into the secrets directory:

```bash
cd ~/secrets/certs
for certfile in admin_operator.pem admin_operator.key.pem admin_operator.pfx
do
    kubectl -n magma cp ${CNTLR_POD}:/var/opt/magma/bin/${certfile} ./${certfile}
done
```

`admin_operator.pem` and `admin_operator.key.pem` are the files that NMS will
use to authenticate itself with the Orchestrator API. `admin_operator.pfx` is
for you to add to your keychain if you'd like to use the Orchestrator REST API
directly (on MacOS, double-click this file and add it to your keychain).

## Deploying NMS

Now that we've got an admin operator cert, we can deploy NMS. Edit the
`vals.yml` from above, and change the `nms` section to the following:

```
nms:
  imagePullSecrets:
    - name: orc8r-secrets-registry

  magmalte:
    manifests:
      secrets: true
      deployment: true
      service: true
      rbac: false

    image:
      repository: YOUR-DOCKER-REGISTRY/magmalte
      tag: YOUR-CONTAINER-TAG

    env:
      api_host: controller.YOURDOMAIN.COM
      mysql_host: YOUR RDS MYSQL HOST
      mysql_user: magma
      mysql_pass: YOUR RDS MYSQL PASSWORD
  nginx:
    manifests:
      configmap: true
      secrets: true
      deployment: true
      service: true
      rbac: false

    service:
      type: LoadBalancer

    deployment:
      spec:
        ssl_cert_name: controller.crt
        ssl_cert_key_name: controller.key
```

You'll just flip all the `manifests` keys to `true` except `rbac`.

Next, copy your `secrets` directory back to the chart (to pick up the admin
certificate), and upload to to S3 (this step is optional, but you should have
some story for where you're storing these).

```bash
rm -r magma/orc8r/cloud/helm/orc8r/charts/secrets/.secrets
cp -r ~/secrets magma/orc8r/cloud/helm/orc8r/charts/secrets/.secrets
aws s3 cp magma/orc8r/helm/orc8r/charts/secrets/.secrets s3://your-bucket --recursive
# Delete the local secrets after you've uploaded them
rm -r ~/secrets
```

We can upgrade the Helm deployment to include NMS components now:

```bash
cd magma/orc8r/cloud/helm/orc8r
helm upgrade orc8r . --values=PATH_TO_VALS/vals.yml
kubectl -n magma get pods
```

Wait for the NMS pods (`nms-magmalte`, `nms-nginx-proxy`) to transition into
`Running` state, then create a user on the NMS:

```bash
kubectl exec -it -n magma \
  $(kubectl -n magma get pod -l app.kubernetes.io/component=magmalte -o jsonpath='{.items[0].metadata.name}') -- \
  yarn setAdminPassword <admin user email> <admin user password>
```

## Upgrading the Deployment

We recommend an upgrade procedure along these lines:

1. `git checkout` the tag of the most recent release on Github
2. Rebuild all the images and push them
3. Update the image tags in vals.yml
4. `aws s3 cp` the secrets bucket in S3 into `.secrets` under the secrets
subchart in Magma
5. Upgrade helm deployment with `helm upgrade`
6. Delete the `.secrets` folder

We've automated steps 4-6 with a fabfile under
`magma/orc8r/cloud/helm/orc8r/fabfile.py`. You can upgrade your deployment
using this fabfile like this:

```bash
fab deploy:PATH_TO_VALS_YML,PATH_TO_TERRAFORM_KUBECONFIG,S3_BUCKET_NAME
```

where `PATH_TO_VALS_YML` is the full path to `vals.yml` on your machine,
`PATH_TO_TERRAFORM_KUBECONFIG` is the full path to the `kubeconfig_orc8r` file
produced by Terraform, and `S3_BUCKET_NAME` is the name of the S3 bucket where
you've uploaded your secrets.
