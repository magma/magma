---
id: version-1.4.0-deploy_install
title: Install Orchestrator
hide_title: true
original_id: deploy_install
---

# Install Orchestrator

This page walks through a full, vanilla Orchestrator install.

If you want to install a specific release version, see the notes in the
[deployment intro](./deploy_intro.md).

## Prerequisites

We assume `MAGMA_ROOT` is set as described in the
[deployment intro](./deploy_intro.md).

This walkthrough assumes you already have the following

- a registered domain name
- a blank AWS account
- an AWS credential with admin permissions

If your AWS account is not blank, this can cause errors while Terraforming.
If you know what you're doing, this is fine - otherwise, consider signing up
for a new account.

Finally, our install process assumes the chosen region contains at least 3
availability zones. This should be the case for all major regions.

## Assemble Certificates

Before Terraforming specific resources, we'll assemble the relevant
certificates.

First, create a local directory to hold the certificates you will use for
your Orchestrator deployment. These certificates will be uploaded to AWS
Secrets Manager and you can delete them locally afterwards.

```bash
mkdir -p ~/secrets/certs
cd ~/secrets/certs
```

You will need the following certificates and private keys placed in this
directory

1. The public SSL certificate for your Orchestrator domain,
with `CN=*.yourdomain.com`. This can be an SSL certificate chain, but it must be
in one file
2. The private key which corresponds to the above SSL certificate
3. The root CA certificate which verifies your SSL certificate

If you aren't worried about a browser warning, you can generate self-signed
versions of these certs

```bash
${MAGMA_ROOT}/orc8r/cloud/deploy/scripts/self_sign_certs.sh yourdomain.com
```

Alternatively, if you already have these certs, rename and move them as follows

1. Rename your public SSL certificate to `controller.crt`
2. Rename your SSL certificate's private key to `controller.key`
3. Rename your SSL certificate's root CA certificate to `rootCA.pem`
4. Put these three files under the directory you created above

Next, with the domain certs placed in the correct directory, generate the
application certs

```bash
${MAGMA_ROOT}/orc8r/cloud/deploy/scripts/create_application_certs.sh yourdomain.com
```

NOTE: `yourdomain.com` above should match the relevant Terraform variables in
subsequent sections. For example, if in `main.tf` the `orc8r_domain_name` is
`orc8r.yourdomain.com`, then that same domain should be used when requesting
or generating all the above certs.

Finally, create the `admin_operator.pfx` file, protected with a password of
your choosing

```bash
$ openssl pkcs12 -export -inkey admin_operator.key.pem -in admin_operator.pem -out admin_operator.pfx

Enter Export Password:
Verifying - Enter Export Password:
```

`admin_operator.pem` and `admin_operator.key.pem` are the files that NMS will
use to authenticate itself with the Orchestrator API. `admin_operator.pfx` is
for you to add to your keychain if you'd like to use the Orchestrator REST API
directly (on macOS, double-click the `admin_operator.pfx` file and add it to
your keychain, inputting the same password chosen above).

The certs directory should now look like this

```bash
$ ls -1 ~/secrets/certs/

admin_operator.key.pem
admin_operator.pem
admin_operator.pfx
bootstrapper.key
certifier.key
certifier.pem
controller.crt
controller.key
fluentd.key
fluentd.pem
rootCA.pem
rootCA.key
```

## Install Orchestrator

With the relevant certificates assembled, we can move on to Terraforming
the infrastructure and application.

### Initialize Terraform

Create a new root Terraform module in a location of your choice by creating a
new `main.tf` file. Follow the example Terraform root module at
`orc8r/cloud/deploy/terraform/orc8r-helm-aws/examples/basic` but make sure to
override the following parameters

- `nms_db_password` must be at least 8 characters
- `orc8r_db_password` must be at least 8 characters
- `orc8r_domain_name` your registered domain name
- `docker_registry` registry containing desired Orchestrator containers
- `docker_user`
- `docker_pass`
- `helm_repo` repo containing desired Helm charts
- `helm_user`
- `helm_pass`
- `seed_certs_dir`: local certs directory (e.g. `"~/secrets/certs"`)
- `orc8r_tag`: tag used when you published your Orchestrator containers
- `orc8r_deployment_type`: type of orc8r deployment (`fwa`, `federated_fwa`, `all`)

If you don't know what values to put for the `docker_*` and `helm_*` variables,
go through the [building Orchestrator](./deploy_build.md) section first.

Make sure that the `source` variables for the module definitions point to
`github.com/magma/magma//orc8r/cloud/deploy/terraform/<module>?ref=v1.4`.
Adjust any other parameters as you see fit - check the READMEs for the
relevant Terraform modules to see additional variables that can be set.

Finally, initialize Terraform

```bash
$ terraform init

Initializing modules...

Initializing the backend...

Initializing provider plugins...

Terraform has been successfully initialized!
```

### Terraform Infrastructure

The two Terraform modules are organized so that `orc8r-aws` contains all the
resource definitions for the cloud infrastructure that you'll need to run
Orchestrator and `orc8r-helm-aws` contains all the application components
behind Orchestrator. On the very first installation, you'll have to
`terraform apply` the infrastructure before the application. On later changes
to your Terraform root module, you can make all changes at once with a single
`terraform apply`.

With your root module set up, run

```bash
$ terraform apply -target=module.orc8r

# NOTE: actual resource count will depend on your root module variables
Apply complete! Resources: 70 added, 0 changed, 0 destroyed.
```

This `terraform apply` will create a
[kubeconfig](https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/)
file in the same directory as your root Terraform module. To access the
K8s cluster, either set your KUBECONFIG environment variable to point to this
file or pull this file into your default kubeconfig file at `~/.kube/config`.

For example, with the [`realpath`](https://linux.die.net/man/1/realpath) utility
installed, you can set the kubeconfig with

```bash
export KUBECONFIG=$(realpath kubeconfig_orc8r)
```

### Terraform Secrets

From your same root Terraform module, seed the certificates and secrets we
generated earlier by running

```bash
$ terraform apply -target=module.orc8r-app.null_resource.orc8r_seed_secrets

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```

The secrets should now be successfully uploaded to AWS Secrets Manager.

NOTE: if this isn't your first time applying the `orc8r_seed_secrets` resource,
you'll need to first
`terraform taint module.orc8r-app.null_resource.orc8r_seed_secrets`.

### Terraform Application

With the underlying infrastructure and secrets in place, we can now install the
Orchestrator application.

From your same root Terraform module, install the Orchestrator application
by running

```bash
$ terraform apply

Apply complete! Resources: 16 added, 0 changed, 0 destroyed.
```

### Create an Orchestrator Admin User

The NMS requires some basic certificate-based authentication when making
calls to the Orchestrator API. To support this, we need to add the relevant
certificate as an admin user to the controller.

NOTE: in the below `kubectl` commands, use the `-n` flag, or
[`kubens`](https://github.com/ahmetb/kubectx), to select the appropriate K8s
namespace (by default this is `orc8r`). Also, assumes kubeconfig is set
correctly from above.

Create the Orchestrator admin user with the `admin_operator` certificate
created earlier

```bash
export ORC_POD=$(kubectl get pod -n orc8r -l app.kubernetes.io/component=orchestrator -o jsonpath='{.items[0].metadata.name}')
kubectl -n orc8r exec ${ORC_POD} -- envdir /var/opt/magma/envdir /var/opt/magma/bin/accessc add-existing -admin -cert /var/opt/magma/certs/admin_operator.pem admin_operator
```

If you want to verify the admin user was successfully created, inspect the
output from

```bash
$ kubectl -n orc8r exec ${ORC_POD} -- envdir /var/opt/magma/envdir /var/opt/magma/bin/accessc list-certs

# NOTE: actual values will differ
Serial Number: 83550F07322CEDCD; Identity: Id_Operator_admin_operator; Not Before: 2020-06-26 22:39:55 +0000 UTC; Not After: 2030-06-24 22:39:55 +0000 UTC
```

At this point, you can `rm -rf ~/secrets` to remove the certificates from your
local disk (we recommend this for security). If you ever need to update your
certificates, you can create this local directory again and `terraform taint`
the `null_resource` to re-upload local certificates to Secrets Manager. You'll
also need to add a new admin user with the updated `admin_operator` cert.

### Create an NMS Admin User

Create an admin user for the `master` organization on the NMS

```bash
export NMS_POD=$(kubectl -n orc8r get pod -l  app.kubernetes.io/component=magmalte -o jsonpath='{.items[0].metadata.name}')
kubectl -n orc8r exec -it ${NMS_POD} -- yarn setAdminPassword master ADMIN_USER_EMAIL ADMIN_USER_PASSWORD
```

## DNS Resolution

EKS has been set up with
[ExternalDNS](https://github.com/kubernetes-sigs/external-dns), so AWS Route53
will already have the appropriate CNAME records for the relevant subdomains of
Orchestrator at this point. You will need to configure your DNS records on
your managed domain name to use the Route53 nameservers in order to resolve
these subdomains.

The example Terraform root module has an output `nameservers` which will list
the Route53 nameservers for the hosted zone for Orchestrator. Access these
via `terraform output` (you have probably already noticed identical output
from every `terraform apply`). Output should be of the form:

```
Outputs:

nameservers = [
  "ns-xxxx.awsdns-yy.org",
  "ns-xxxx.awsdns-yy.co.uk",
  "ns-xxxx.awsdns-yy.com",
  "ns-xxxx.awsdns-yy.net",
]
```

If you chose a subdomain prefix for your Orchestrator domain name in your
root Terraform module, you only need to provide a single NS record to your
domain registrar, mapping the subdomain to the above name servers.
For example, for the subdomain `orc8r`, this record would notionally take
the form `{ orc8r -> [ns-xxxx.awsdns-yy.org, ...] }`.

If you didn't choose a subdomain prefix, then you can still point the whole
domain to AWS via the single NS record. Alternatively, if this is undesirable,
provide NS records for each of the following subdomains

1. nms
2. controller
3. bootstrapper-controller
4. api

For example, for the domain `mydomain`, these records would notionally take
the form
`{ nms -> [ns-xxxx.awsdns-yy.org, ...], controller -> [ns-xxxx.awsdns-yy.org, ...], ... }`.

## Verify the Deployment

After a few minutes the NS records should propagate. Confirm successful
deployment by visiting the master NMS organization at e.g.
`https://master.nms.yoursubdomain.yourdomain.com` and logging in with the
`ADMIN_USER_EMAIL` and `ADMIN_USER_PASSWORD` provided above.

NOTE: the `https://` is required. If you self-signed certs above, the browser
will rightfully complain. Either ignore the browser warnings at your own risk
(some versions of Chrome won't allow this at all), or e.g.
[import the root CA from above on a per-browser basis
](https://stackoverflow.com/questions/7580508/getting-chrome-to-accept-self-signed-localhost-certificate).

For interacting with the Orchestrator REST API, a good starting point is the
Swagger UI available at `https://api.yoursubdomain.yourdomain.com/apidocs/v1/`.

If desired, you can also visit the AWS endpoints directly. The relevant
services are `nginx-proxy` for NMS and `orc8r-nginx-proxy` for Orchestrator
API. Remember to include `https://`, as well as the port number for
non-standard TLS ports.

```bash
$ kubectl get services

# NOTE: values will differ, e.g. the EXTERNAL-IP column
NAME                            TYPE           CLUSTER-IP       EXTERNAL-IP                       PORT(S)                                                     AGE
fluentd                         LoadBalancer   172.20.213.111   aaa.us-west-2.elb.amazonaws.com   24224:31621/TCP                                             3h13m
magmalte                        ClusterIP      172.20.197.108   <none>                            8081/TCP                                                    3h13m
nginx-proxy                     LoadBalancer   172.20.1.201     www.us-west-2.elb.amazonaws.com   443:32422/TCP                                               3h13m
orc8r-accessd                   ClusterIP      172.20.128.137   <none>                            9180/TCP                                                    3h13m
orc8r-alertmanager              ClusterIP      172.20.165.206   <none>                            9093/TCP                                                    3h13m
orc8r-alertmanager-configurer   ClusterIP      172.20.92.62     <none>                            9101/TCP                                                    3h13m
orc8r-analytics                 ClusterIP      172.20.152.243   <none>                            9180/TCP                                                    3h13m
orc8r-bootstrap-nginx           LoadBalancer   172.20.232.199   xxx.us-west-2.elb.amazonaws.com   80:31116/TCP,443:31302/TCP,8444:31093/TCP                   3h13m
orc8r-bootstrapper              ClusterIP      172.20.65.124    <none>                            9180/TCP                                                    3h13m
orc8r-certifier                 ClusterIP      172.20.89.150    <none>                            9180/TCP                                                    3h13m
orc8r-clientcert-nginx          LoadBalancer   172.20.143.232   yyy.us-west-2.elb.amazonaws.com   80:30546/TCP,443:31400/TCP,8443:30781/TCP                   3h13m
orc8r-configurator              ClusterIP      172.20.56.203    <none>                            9180/TCP                                                    3h13m
orc8r-ctraced                   ClusterIP      172.20.134.117   <none>                            9180/TCP                                                    3h13m
orc8r-device                    ClusterIP      172.20.103.126   <none>                            9180/TCP                                                    3h13m
orc8r-directoryd                ClusterIP      172.20.4.31      <none>                            9180/TCP                                                    3h13m
orc8r-dispatcher                ClusterIP      172.20.124.178   <none>                            9180/TCP                                                    3h13m
orc8r-ha                        ClusterIP      172.20.201.112   <none>                            9180/TCP                                                    3h13m
orc8r-lte                       ClusterIP      172.20.225.103   <none>                            9180/TCP,8080/TCP                                           3h13m
orc8r-metricsd                  ClusterIP      172.20.159.39    <none>                            9180/TCP,8080/TCP                                           3h13m
orc8r-nginx-proxy               LoadBalancer   172.20.52.234    zzz.us-west-2.elb.amazonaws.com   80:30034/TCP,8443:31884/TCP,8444:31829/TCP,443:30124/TCP    3h13m
orc8r-obsidian                  ClusterIP      172.20.41.215    <none>                            9180/TCP,8080/TCP                                           3h13m
orc8r-orchestrator              ClusterIP      172.20.172.120   <none>                            9180/TCP,8080/TCP                                           3h13m
orc8r-policydb                  ClusterIP      172.20.95.10     <none>                            9180/TCP,8080/TCP                                           3h13m
orc8r-prometheus                ClusterIP      172.20.65.141    <none>                            9090/TCP                                                    3h13m
orc8r-prometheus-cache          ClusterIP      172.20.111.91    <none>                            9091/TCP,9092/TCP                                           3h13m
orc8r-prometheus-configurer     ClusterIP      172.20.106.4     <none>                            9100/TCP                                                    3h13m
orc8r-service-registry          ClusterIP      172.20.146.78    <none>                            9180/TCP                                                    3h13m
orc8r-smsd                      ClusterIP      172.20.63.198    <none>                            9180/TCP,8080/TCP                                           3h13m
orc8r-state                     ClusterIP      172.20.185.245   <none>                            9180/TCP                                                    3h13m
orc8r-streamer                  ClusterIP      172.20.57.35     <none>                            9180/TCP                                                    3h13m
orc8r-subscriberdb              ClusterIP      172.20.238.111   <none>                            9180/TCP,8080/TCP                                           3h13m
orc8r-tenants                   ClusterIP      172.20.173.57    <none>                            9180/TCP,8080/TCP                                           3h13m
orc8r-user-grafana              ClusterIP      172.20.149.141   <none>                            3000/TCP                                                    3h13m
```

## Upgrade the Deployment

You can upgrade the deployment by changing one or both of the following
variables in your root Terraform module, before running `terraform apply`

- `orc8r_tag` container image version
- `orc8r_chart_version` Helm chart version

Changes to the Terraform modules between releases may require some updates to
your root Terraform module - these will be communicated in release notes.
