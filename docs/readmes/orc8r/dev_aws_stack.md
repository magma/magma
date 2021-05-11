---
id: dev_aws_stack
title: AWS Stack
hide_title: true
---

# AWS Stack

This document describes the Orchestrator stack when it's deployed to AWS.

![orc8r_aws_stack](assets/orc8r/orc8r_aws_stack.png)

## What does Terraform deploy in my AWS account?

The following resources are the basic infrastructure necessary to deploy Magma services

* **Kubernetes cluster (EKS).** Orc8r runs over a Kubernetes (K8s) cluster, Magma uses Elastic Kubernetes Service to deploy the cluster. Other resources are created to support this cluster.
* **Elastic File System (EFS).** File system to hold persistent data for Kubernetes cluster.
* **Secrets.** Magma creates an instance in AWS Secret Manager to hold all the certificates and keys used in K8s cluster.
* **Network and Security.** In order to provide networking a Virtual Private Cloud (VPC) is configured in AWS and to protect it's access a Security Group (SG) is configured for it.

* **ElasticSearch Domain.** This is used to store logs of the K8s pods and services
* **Database instances.** Magma needs a database instance for Orc8r and NMS, this instance is deployed using AWS Relational Databases Services (RDS)
* **DNS.** Magma uses AWS Route53 service to configure external DNS to access Orc8r's UI and API resources.

Magma provides a set of configuration parameters to control these resources, follow a short list with some configurations.

|Name	|Type	|Description	|Default Value	|
|---	|---	|---	|---	|
|orc8r_db_password	|Mandatory	|Password for the Orchestrator DB. Must be at least 8 characters.	|-	|
|orc8r_db_storage_gb	|Optional	|Capacity in GB to allocate for Orchestrator RDS instance.	|64	|
|orc8r_db_engine_version	|Optional	|Postgres engine version for Orchestrator DB.	|9.6.15	|
|cluster_version	|Optional	|Kubernetes version for the EKS cluster.	|1.17	|
|elasticsearch_version	|Optional	|ES version for ES domain.	|7.1	|

The values for these configurations can be defined/overwritten in your `main.tf` `orc8r` module.
Some configurations do not have default values and are mandatory to be present in your `main.tf`.
Other already have predefined values and can be overwritten.
See the [complete list of available configurations](http://github.com/magma/magma/blob/master/orc8r/cloud/deploy/terraform/orc8r-aws/variables.tf).

## What Magma deploys in my AWS-hosted Kubernetes?


Before starting to deploy in your Kubernetes environment, some AWS resource need to be linked to it, in order to have that Magma deploys.

* **Kubernetes Secrets.** Create secrets in Kubernetes to mirror the AWS secrets, also create secrets to hold the Orc8r and NMS certificates.
* **External DNS.** This is an Helm chart responsible for activating Route53 DNS on K8s
* **ElasticSearch Curator.** This is a Helm chart responsible for cleaning ElasticSearch from old logs.


Now that Kubernetes cluster is complete, is time to set up the applications inside of it. Magma performs that using Helm charts, follows a list of charts deployed in your cluster:

* **EFS Provisioner.** This is an Helm chart responsible for activating EFS (created previous) as a persistent volume in K8s
* **FluentD.** This is a Helm chart responsible for aggregating the logs from all AGWs sources.
* **OpenVPN.** This is an Helm chart responsible for creating a VPN server to access AGW securely.
* **Orc8r Application.** This is the Helm chart that deploys Oc8r and NMS applications in your Kubernetes infrastructure. This chart include multiple services and pods, including prometheus, grafana, NMS, orc8r-controller and other.

Magma provides a set of configuration parameters to control these resources, follow a short list with some configurations:

|Name	|Type	|Description	|Default Value	|
|---	|---	|---	|---	|
|region	|Mandatory	|AWS region to deploy Orchestrator components into. The chosen region must provide EKS.	|-	|
|orc8r_domain_name	|Mandatory	|Base domain name for Orchestrator.	|-	|
|orc8r_deployment_type	|Mandatory	|Type of Orc8r deployment (fixed wireless access, federated fixed wireless access, or all modules)	|-	|
|orc8r_tag	|Mandatory	|Image tag for Orchestrator components.	|-	|
|deploy_openvpn	|Optional	|Flag to deploy OpenVPN server into cluster. This is useful if you want to remotely access AGWs.	|FALSE	|

The values of these configurations can be defined or overwritten in your `main.tf` files's `orc8r-app` module.
See the [complete list of configurations available for these resource](http://github.com/magma/magma/blob/master/orc8r/cloud/deploy/terraform/orc8r-helm-aws/variables.tf).

## What does the Orc8r Helm chart contain?


Orc8r Helm chart is composed of 5 other charts:

1. **Logging.** deploys all the services related to collect and store logs, it includes:
    1. FluentD: for log aggregation
    2. Elastic Search Curator: to clean old logs
2. **Metrics.** deploys all the services related to collect, store and show metrics, it includes:
    1. Prometheus for metric storage.
    2. Prometheus Cache
    3. Alert-Manager: For alert generation
    4. [More details](https://github.com/magma/magma/blob/master/orc8r/cloud/helm/orc8r/charts/metrics/README.md)
3. **NMS.** deploys the Orc8r web UI and API, it includes:
    1. Nginx HTTP proxy
    2. NMS service
    3. [More details](https://github.com/magma/magma/blob/master/orc8r/cloud/helm/orc8r/charts/nms/README.md)
4. **Orc8r.** deploys the Magma orchestration function
    1. Orc8r controller: Stream configuration to all AWGs in the network. Collect and export metrics from all the AGWs in the network.
    2. [More details](https://github.com/magma/magma/blob/master/orc8r/cloud/helm/README.md)
5. **Secrets.** is used to apply a set of secrets required by the Magma Orchestrator.
    1. [More details](https://github.com/magma/magma/blob/master/orc8r/cloud/helm/orc8r/charts/secrets/README.md)

## Which certificates does Magma use?

Orc8r needs certificates to assure messages traveling over the internet are encrypted and secure. Certificates are created in pairs with a public certificate (`.crt` or `.pem`) and a private key (`.key`). The former needs to be safely stored and kept secret, while the latter can be distributed to clients. Please read the [certificates section on Orchestrator architecture](https://magma.github.io/magma/docs/next/orc8r/dev_security#certificates).

> ***Note.*** All the certificates are created with validity time period, so make sure you know when your certificates expire and schedule a maintenance to update them.


