---
id: orc8r_stack_on_aws
title: Orchestrator Stack on AWS
hide_title: true
---

# Orchestrator Stack on AWS

![orc8r_aws_stack](assets/orc8r/orc8r_aws_stack.png)

## What does terraform deploy in my AWS account?

These following resources are the basic infrastructure necessary to deploy Magma services:

* **Kubernetes cluster:** Orc8r runs over a Kubernetes (K8s) cluster, Magma uses Elastic Kubernetes Service (EKS) to deploy the cluster. Other resources are created to support this cluster.
* **Elastic File System:** File system to hold persistent data for Kubernetes cluster.
* **Secrets:** Magma creates an instance in AWS Secret Manager to hold all the certificates and keys used in K8s cluster.
* **Network and Security:** In order to provide networking a Virtual Private Cloud is configured in AWS and to protect it's access a Security Group is configured for it.

* **ElasticSearch Domain:** This is used to store logs of the K8s pods and services
* **Database instances:** Magma needs two database instances one for Orc8r another for NMS, these instances are deployed using AWS Relational Databases Services (RDS)
* **DNS:** Magma uses AWS Route53 service to configure external DNS to access Orc8r's UI and API resources.


[Detailed list of configurations available for these resource](http://github.com/magma/magma/blob/master/orc8r/cloud/deploy/terraform/orc8r-aws/variables.tf). The values for these configurations can be defined/overwritten in your main.tf `orc8r` module. Some configurations do not have default values and are mandatory to be present in your main.tf. Other are already have predefined values that can be overwritten. Follows an example of some configurations:


|Name	|Type	|Description	|Default Value	|
|---	|---	|---	|---	|
|nms_db_password	|Mandatory	|Password for the NMS DB. Must be at least 8 characters.	|-	|
|orc8r_db_password	|Mandatory	|Password for the Orchestrator DB. Must be at least 8 characters.	|-	|
|orc8r_db_storage_gb	|Optional	|Capacity in GB to allocate for Orchestrator RDS instance.	|64	|
|nms_db_storage_gb	|Optional	|Capacity in GB to allocate for NMS RDS instance.	|16	|
|orc8r_db_engine_version	|Optional	|Postgres engine version for Orchestrator DB.	|9.6.15	|
|nms_db_engine_version	|Optional	|MySQL engine version for NMS DB.	|5.7	|
|cluster_version	|Optional	|Kubernetes version for the EKS cluster.	|1.17	|
|elasticsearch_version	|Optional	|ES version for ES domain.	|7.1	|

## What Magma deploys in my AWS-hosted Kubernetes?


Before starting to deploy in your Kubernetes environment, some AWS resource need to be linked to it, in order to have that Magma deploys:

* **Kubernetes Secrets:** Create secrets in Kubernetes to mirror the AWS secrets, also create secrets to hold the Orc8r and NMS certificates.
* **External DNS:** This is an Helm chart responsible for activating Route53 DNS on K8s
* **ElasticSearch Curator:** This is a Helm chart responsible for cleaning ElasticSearch from old logs.


Now that Kubernetes cluster is complete, is time to set up the applications inside of it. Magma performs that using Helm charts, follows a list of charts deployed in your cluster:

* **EFS Provisioner:** This is an Helm chart responsible for activating EFS (created previous) as a persistent volume in K8s
* **FluentD:** This is a Helm chart responsible for aggregating the logs from all AGWs sources.
* **OpenVPN:** This is an Helm chart responsible for creating a VPN server to access AGW securely.
* **Orc8r Application:** This is the Helm chart that deploys Oc8r and NMS applications in your Kubernetes infrastructure. This chart include multiple services and pods, including prometheus, grafana, NMS, orc8r-controller and other.

[Detailed list of configurations available for these resource](http://github.com/magma/magma/blob/master/orc8r/cloud/deploy/terraform/orc8r-helm-aws/variables.tf). The values of these configurations can be defined/overwritten in your main.tf `orc8r-app` module. Follows an example of some configurations:


|Name	|Type	|Description	|Default Value	|
|---	|---	|---	|---	|
|region	|Mandatory	|AWS region to deploy Orchestrator components into. The chosen region must provide EKS.	|-	|
|orc8r_domain_name	|Mandatory	|Base domain name for Orchestrator.	|-	|
|orc8r_deployment_type	|Mandatory	|Type of Orc8r deployment (fixed wireless access, federated fixed wireless access, or all modules)	|-	|
|orc8r_tag	|Mandatory	|Image tag for Orchestrator components.	|-	|
|deploy_openvpn	|Optional	|Flag to deploy OpenVPN server into cluster. This is useful if you want to remotely access AGWs.	|FALSE	|

## What does the Orc8r Helm chart contain?


Orc8r Helm chart is composed of 5 other charts:

1. **Logging:** deploys all the services related to collect and store logs, it includes:
    1. FluentD: for log aggregation
    2. Elastic Search Curator: to clean old logs
2. **Metrics:** deploys all the services related to collect, store and show metrics, it includes:
    1. Prometheus/Thanos for metric storage.
    2. Prometheus Cache
    3. Alert-Manager: For alert generation
    4. [More details](https://github.com/magma/magma/blob/master/orc8r/cloud/helm/orc8r/charts/metrics/README.md)
3. **NMS:** deploys the Orc8r web UI and API, it includes:
    1. Nginx HTTP proxy
    2. NMS service
    3. [More details](https://github.com/magma/magma/blob/master/orc8r/cloud/helm/orc8r/charts/nms/README.md)
4. **Orc8r:** deploys the Magma orchestration function
    1. Orc8r controller: Stream configuration to all AWGs in the network. Collect and export metrics from all the AGWs in the network.
    2. [More details](https://github.com/magma/magma/blob/master/orc8r/cloud/helm/README.md)
5. **Secrets:** is used to apply a set of secrets required by the Magma Orchestrator.
    1. [More details](https://github.com/magma/magma/blob/master/orc8r/cloud/helm/orc8r/charts/secrets/README.md)

## Which certificates Magma use?

Orc8r needs certificates to assure messages traveling over the internet are encrypted and secure. Certificates are created in pairs with a public certificate (.crt or .pem) and a private key (.key). The former needs to be safely stored and kept secret, while the latter can be distributed to clients.

**Attention:** All the certificates are created with validity time period, so make sure you know when your certificates expire and schedule a maintenance to update them.

### Certificate Authority (and rootCA)

Depending on the application you are running, it may require your certificate to be issued by an Authority, this means the certificate is signed by a third party that is trusted by both client and server application. You can create your own Authority certificate (rootCA) at the price of not being trusted by some applications (like browsers).

If you choose to self-sign your certificates, you will need to create a rootCA:

**rootCA.key, rootCA.pem:** certs for trusted root Certificate Authority, will be used to sign other certificates.

Follows the list of applications that need certificates in Magma:

* **Controller Certificate (controller.key, controller.crt):**
    * Certs for Orc8r deployment's public domain name,
    * Needs to be signed by a Certificate Authority
* **Admin Certificate** (**admin_operator.key.pem, admin_operator.pem):**
    * Admin certs for the initial admin operator (e.g. whoever's deploying Orc8r) to authenticate to the Orc8r proxy
        * are the files that NMS will use to authenticate itself with the Orchestrator API
    * Admin certs are used to access the Swagger API in the Orc8r and it is also used by NMS to connect to the Orc8r controller.
* **Fluentd Certificate (fluentd.key, fluentd.pem):**
    * Certs for fluentd endpoint, allowing gateways to securely send logs (fluentd is currently outside Orc8r proxy)
* **Certifier Certificate (certifier.key, certifier.pem):**
    * Certs for the controller's certifier service, providing more fine-grained access to controller services
* **Bootstrapper Key (bootstrapper.key):**
    * Private key for controller's bootstrapper service, used in gateway bootstrapping challenges
