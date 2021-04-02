
## Proposal: Tagging Infrastructure, Platform and Magma Services 

Author:  @arunuke

Requested Reviewers:  @karthiksubraveti, @hcgatewood, @mpgermano


## Abstract

Magma services (orchestrator) consume various infrastructure (compute, storage, networking) and platform (logging, database, monitoring) resources. These resources will be tagged (used interchangeably with ‘labeled’) for cost, monitoring, updating and automation related use-cases based on the target platform in Cloud, Hybrid and On-prem scenarios.


## Background

Magma’s deployment tools (ex: Cloudstrapper, Bare-metal) set up infrastructure and platform resources before deploying containerized Magma services using helm packages. Some resources are identified using the “Name” tag while many others are not. Additionally, simplified names used in tags conflict with other orchestrator deployments within the same organization. This causes issues when cleaning up an environment, identifying ownerships and which deployment a certain resource belongs to.


## Proposal

The purpose of this document is to define a framework allowing the deployer to tag every infrastructure, platform and Magma service which will help in uniquely identifying a resource that belongs to a given deployment cluster. For the purposes of this document, a deployment cluster relates to a complete deployment of Magma’s Orchestrator that includes infrastructure resources, shared services, Magma services and in certain cases, the associated gateways. The operational question here is to determine the resources that an administrator would want to manage as a single cohesive unit. 

The usage of tags is fairly common for Cloud-native applications. Cloud providers such as AWS/GCP and cluster managers such as Kubernetes support arbitrary tag names and related values as a key-value pair. Using a combination of these tools, all resources deployed in a Magma cluster can be tagged using various arbitrary values. This document proposes an initial set of key-value pairs that can be expanded to support future use-cases. 

Infrastructure and platform resources are deployed via Terraform in Cloud platforms with Ansible available for bare-metal platforms. Magma services are deployed as Kubernetes services. Each of these tools allows a native way to label resources. For example, AWS uses arbitrary tags with a “key” and a “value” that can be set for any deployed resource while Kubernetes allows users to define labels with a “key” and “value” as part of metadata. All resources deployed in a single Orchestrator development would share similar key-value pairs to indicate they are part of the same deployment and mandatory tags. Individual resources might choose to add other optional tags as needed.

Administrators will be able to request the system to auto-generate a UUID that will be used for all resources that are part of the deployment and returned to the user. To add new resources to an existing deployment, the administrator can provide the deployment’s UUID in a configuration file that will be used for all the new resources deployed.

To support migration and upgrade use-cases, where an AGW needs to be added to a newer version of Orchestrator, tags can be edited to reflect the value of the new orchestrator’s UUID without having any impact on services.


## Rationale

An alternate approach would be to consider using a platform-neutral way to tag individual resources. As an example, Terraform supports arbitrary tagging for resources deployed and it can be used as a consistent way across multiple Cloud platforms. However, not all resources in Magma are deployed via Terraform which eliminates the possibility of using a single tool. In this scenario, using a native feature keeps the number of tools required to the absolute minimum and ensures only essential tools are used.

An additional feature where additional clusters could use shared services across other clusters (ex: an Orchestrator cluster with UUID=x does not deploy an elasticsearch service but shares the elasticsearch service of Orchestrator cluster with UUID=y) was considered. However, it was dropped due to the disproportionate relationship of the complexity it introduces versus the benefits it provides. If an elasticsearch service is considered too large for a deployment, it can always be scaled back. This approach is preferred over introducing complexity of sharing services across Orchestrator clusters.


## Compatibility

This feature does not provide backward compatibility if the tag names of choice are used by existing deployments. This feature will not be backported and will be supported from the Magma release it lands, tentatively 1.5.


## Observability and Debug

All tags are retrievable using CLI options (aws, kubectl) and validated. The respective CLIs can also be used to filter resources based on tags.


## Implementation



-  Proposal review and inviting comments on Github discussions
-  Implementation plan for AWS and Kubernetes resources as the first deliverable to support currently supported deployment models (AWS or
On-Prem) for v1.5 targeted for April 2021.
-  Phase 1 will implement tagging for all AWS resources
-  Phase 2 will tag all other kubernetes resources in orchestrator


## Future Work


Most Cloud provider environments provide a way to tag resources. The current proposal will implement the findings on AWS and will extend the support to other Cloud provider platforms and on-prem options. However, as new services are containerized and managed via kubernetes, they can be tagged even if they are outside the scope of the orchestrator (ex: Cloud HA instances)


## Design



For tags, follow the current guidelines
  -   Minimal number of tags to start with and incrementally increasing as additional use-cases are identified
  	-   The primary use-case for now is to use labels and tags for A/B or Blue/Green or Canary based update of services
  	-   The secondary use-case for now is to associate all resources in a deployment together so that they can be searched for and cleaned up
  -   Identify mandatory, optional and resource-specific tags 
	  -   Mandatory: UUID, named magmaUUID
    -   For Cloud deployments, Terraform is currently used with a provider for the actual Cloud platform that hosts the services. Provider-native mechanisms will be used to tag the resource. An example block is given below where an AWS provider resource is tagged with the UUID.

```

	resource "aws_db_instance" "nms" { 
		identifier = var.nms_db_identifier 
		allocated_storage = var.nms_db_storage_gb 
		engine = "mysql" 
		engine_version  = var.nms_db_engine_version
		instance_class = var.nms_db_instance_class
  		Tags = { magmaUUID = "<128-bit Identifier>" } 
	} 
```


Any non-Cloud provider resources will be deployed via Kubernetes and can be tagged (labeled) as follows. These resources can be filtered via label selectors via the kubectl command. For instances starting up, use the kubernetes configuration file. To generate dynamic configuration manifests, extensions are available to create custom kubernetes operators in the future.

``` 
	template: 
          metadata:
            labels:
              app.kubernetes.io/component: orc8r
              app.kubernetes.io/instance: <Identifier>
```


For instances already running, use the kubectl label call to set it during runtime

```
            kubectl label pods <orc8r pods> app.kubernetes.io/magmaUUID=<128-bit Identifier>
```

The following tags will be used 
    -   UUID: 128-bit identifier to uniquely identify a deployment. The value would be generated randomly at the time of deployment and all resources deployed as part of this effort will be uniquely tagged. 
    -   Two alternatives under consideration for the UUID format
	-   [RFC 4122](https://tools.ietf.org/html/rfc4122) with a golang implementation available from [Go Packages.](https://github.com/google/uuid)
	-   [UUID v4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random))

             

##   Workflow

At deployment time, the administrator sets a flag for the deployer to automatically generate a UUID, set it to all resources and return it to
the user.
    -   Alternatively, the administrator sets the UUID value computed using the UUID generation tool in the configuration file that is used as a tag for all resources.

When new resources are being added to the deployment, they are added with the same tag to indicate they are part of the same deployment.
