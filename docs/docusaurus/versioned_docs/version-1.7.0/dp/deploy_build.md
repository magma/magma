---
id: version-1.7.0-deploy_build
title: Build and deployment
hide_title: true
original_id: deploy_build
---

# Domain Proxy Installation

This page describes the installation of orc8r domain-proxy component.

## Prerequsites

In order to work correctly Domain Proxy requires Orchestrator to be installed as it is tightly coupled with it. We
assume `MAGMA_ROOT` is set as described in the [deployment intro](orc8r/deploy_intro.md).

This walkthrough assumes you already have the following:

- Properly installed Orc8r using this procedure

## Certificates

Before updating terraform variables, proper set of certificates must be placed in `seed_certs_dir` directory.

There are three components required to communicate with production grade SAS, namely:

- Certificate file (e.g `tls.crt`): Public key signed by a CA recognized by the SAS operator.
- Key file (e.g `tls.key`): Private key used to verify your public key.
- Certificate authority chain file (e.g `ca.crt`): Chain of certificates that verifies the CA.

## Install Domain Proxy

Change your directory to your terraform root module.

Adjust the following terraform parameters:

- `dp_enabled`: `true`
- `dp_sas_crt`: certificate filename (defaults to `tls.crt`)
- `dp_sas_key`: private key filename (defaults to `tls.key`)
- `dp_sas_ca`: ca chain filename (defaults to `ca.crt`)
- `dp_sas_endpoint_url`: must be set to point to your SAS which will accept provided certificate.

Finally, apply update to your terraform module.

```console
terraform apply
```

The following new services should be available on your cluster.

```console
kubectl --namespace orc8r get pod -l app.kubernetes.io/name=domain-proxy
NAME                                                    READY   STATUS      RESTARTS   AGE
domain-proxy-active-mode-controller-7b984c6579-zmwrm    1/1     Running     0          12d
domain-proxy-configuration-controller-6d99c978f-b8h6b   1/1     Running     0          12d
domain-proxy-radio-controller-5c868696d9-s7vgg          1/1     Running     0          12d
```

## Upgrade the Deployment

You can upgrade the deployment by changing one or both of the following variables in your root Terraform module, before
running terraform apply.

- `orc8r_tag` container image version.
- `dp_orc8r_chart_version` Domain proxy helm chart version.
