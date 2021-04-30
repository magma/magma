---
id: version-1.5.0-deploy_orcl
title: Experimental Orc8r Deployer
hide_title: true
original_id: deploy_orcl
---

# Orcl

Orcl is a Orchestrator CLI. It is used for managing Orc8r deployment. It provides following subcommands.

- Configure
- Certs
- Install
- Upgrade
- Verify
- Cleanup
- Debug(perhaps in the future)

![Orcl Big Picture](assets/orc8r/orcl.png)
Orcl is packaged within orc8r_deployer. Orc8r deployer is a docker image which contains all the necessary prerequisites to deploy orc8r. The only requirements for running orc8r_deployer is that the the host machine must have [docker engine installed](https://docs.docker.com/get-docker/).

## Usage

```
./run_deployer runs the orc8r deployer container
orc8r deployer contains scripts which enable user to configure, run prechecks,
install, upgrade, verify and cleanup an orc8r deployment

Usage: run_deployer [-deploy-dir|-root-dir|-build|-h]
options:
-h           Print this Help
--deploy-dir  deployment dir containing configs and secrets (mandatory)
--root-dir    magma root directory
--build       build the deployer container
--test        'all' or any specific test function[run_unit_tests,check_helmcharts_insync, check_tfvars_insync ]
example: ./run_deployer.bash --deploy-dir ~/orc8r_15_deployment
Note that for the first time the script has to be run with -build option.
 ./run_deployer.bash --deploy-dir ~/orc8r_15_deployment --build
```

Deployment directory will be used to maintain the configs, secrets and also contain the deployment related files such as main.tf etc. It is important to use the same deployment directory during the upgrade so that the configuration variables can be reused.

In the following sections, we will discuss each of the orcl commands in detail.

## Configure Command

Every orc8r deployment relies on several configuration attributes. For example, cluster_name attribute is used to identify the orc8r  kubernetes cluster, orc8r_tag provides the image tag to be used during the deployment. Configure command enables user to easily configure the mandatory configs necessary for the deployment. It additionally provides the ability to also configure optional attributes through the ***set*** subcommand. Configure command also has subcommands like ***info*** to show all the possible configuration attributes and ***show*** to display the current configuration. Finally configure also contains a subcommand ***check*** which provides the ability to check if all mandatory configs needed by the deployment have been successfully configured.

```
# orcl configure --help
Usage: orcl configure [OPTIONS] COMMAND [ARGS]...

  Configure orc8r deployment variables

Options:
  -c, --component [infra|platform|service]
  --help                          Show this message and exit.

Commands:
  check  Check if all mandatory configs are present
  info   Display all possible config options in detail
  set    Set enables user to configure any configuration option.
  show   Display the current configuration
```

Note: Currently we don’t have the capability to sanitize the inputs. This might be a feature to be added later

Configure command relies mainly on a meta file `vars.yml`. This file contains all deployment variables with their name, type, description, default value, whether it is a required or an optional attribute and finally the apps which rely on this configuration attribute.
The configuration provided by the user is used to

- set the aws configuration(specifically access_key, secret_key, region),
- build terraform.tfvars.json(which is automatically loaded by terraform and provided to the root module)
- build main.tf and vars.tf(to ensure that we set the vars in the root module for it to be correspondingly consumed by the orc8r and orc8r-app modules)

For example,

```
module "orc8r" {
  source = "/root/magma/orc8r/cloud/deploy/terraform/orc8r-aws"
  cluster_name=var.cluster_name <-- generated from configured keys
  cluster_version=var.cluster_version
  orc8r_domain_name=var.orc8r_domain_name
  region=var.region
```

## Certs Command

Certs is a orcl command for managing certificates. Currently it has only one subcommand ***add*** which can be used to add application certs and self signed root certificates(optional).

```
# orcl certs --help
Usage: orcl certs [OPTIONS] COMMAND [ARGS]...
  Manage certs in orc8r deployment

Options:
  --help  Show this message and exit.

Commands:
  add  Add creates application and self signed(optional) certs
```


Certs relies on the orc8r configuration options being configured. Specifically it relies on `orc8r_domain_name` attribute to be configured.

## Install Command

Install command provides the ability to install an orc8r deployment based on provided configuration. Install command provides an optional ability to run prechecks prior to installation and also provides prechecks as a separate subcommand.

```
# orcl install --help
Usage: orcl install [OPTIONS] COMMAND [ARGS]...

  Deploy new instance of orc8r

Options:
  --help  Show this message and exit.

Commands:
  precheck  Performs various checks to ensure successful installation
```

Currently the installation process is dependent on terraform. Install command runs following subprocess commands

```
terraform init
terraform apply -target=module.orc8r -auto-approve
terraform apply -target=module.orc8r-app.null_resource.orc8r_seed_secrets -auto-approve
terraform apply

```

Following prechecks are performed prior to the installation

- Check if all mandatory configs are present
- Check if all secret files are present
- Check if we have atleast 3 availability zones in our deployment
- Check if secrets manager is already configured with the orc8r secrets id
- Check cloudwatch log group already exists
- In case of elastic search deployment, check if AWSServiceRoleForAmazonElasticsearchService IAM role already exists
- Check if orc8r deployment type is valid
- Check if image tag specified is present in image repository
- Check if helm chart specified is present in the helm repository

## Upgrade Command

Upgrade command provides the ability to install an orc8r deployment based on provided configuration. Upgrade command provides an optional ability to run prechecks prior to upgrade and also provides prechecks as a separate subcommand.
**Note: it is important to use the same deployment directory which was used during install. This is mainly to ensure that the terraform state is available and also reuse the configuration variables**

```
# orcl upgrade --help
Usage: orcl upgrade [OPTIONS] COMMAND [ARGS]...

  Upgrade existing orc8r deployment

Options:
  --help  Show this message and exit.

Commands:
  precheck  Precheck runs various checks to ensure successful upgrade
```

Upgrade command runs the following upgrade command

```
terraform apply
```

Following prechecks are performed prior to the upgrade

- Check if terraform state exists
- Check if all mandatory configs are present
- Check if we have atleast 3 availability zones in our deployment
- Check if deployed EKS cluster version is greater than what’s specified in the configuration
- In case of elastic search deployment, check if AWSServiceRoleForAmazonElasticsearchService IAM role already exists
- Check if deployed RDS instance version is greater than what’s specified in the configuration
- Check if image tag specified is present in image repository
- Check if helm chart specified is present in the helm repository

## Verify Command

Verify command provides the ability to run post-deployment checks on orc8r.

```
Usage: orcl verify [OPTIONS] COMMAND [ARGS]...

  Run post deployment checks on orc8r

Options:
  --help  Show this message and exit.

Commands:
  sanity
```

Currently the subcommand ***sanity*** runs the following checks

- Check if all kubernetes pods are healthy and dump the logs of unhealthy pods
- Invoke get all networks API for generic and LTE networks(magma/v1/networks, magma/v1/lte) to ensure if the orchestrator and lte pods are working as expected

## Cleanup Command

Cleanup command provides the ability to cleanup all resources deployed during orc8r deployment

```
# orcl cleanup --help
Usage: orcl cleanup [OPTIONS] COMMAND [ARGS]...

  Removes resources deployed for orc8r

Options:
  --help  Show this message and exit.

Commands:
  raw  Individually cleans up resources deployed for orc8r
```

Cleanup runs

```
terraform destroy
```

to perform the cleanup. Unfortunately terraform cleanup hasn’t been very reliable in the past. So this command provides ability to also directly cleanup the underlying resources when the terraform destroy command fails. It is though the subcommand ***raw***

```
Usage: orcl cleanup raw [OPTIONS]

  Individually cleans up resources deployed for orc8r

Options:
  --dryrun      Show resources to be cleaned up during raw cleanup
  --state TEXT  Provide state file containing resource information example,
                terraform.tfstate or terraform.tfstate.backup
  --override    Provide values to cleanup the orc8r deployment
  --help        Show this message and exit.
```

Currently the raw cleanup performs cleanup of the following

- Orc8r cloudwatch log group
- RDS instances
- Elastic search instance.
- EFS mount targets and volumes
- EKS cluster
- Autoscaling groups
- Elastic load balancers
- NAT gateways
- Internet gateways
- Subnets
- VPC
- Hosted records in hosted zone and hosted zone

Note: The raw cleanup attempts to identify the above mentioned resources through `terraform show --json`
Sometimes terraform destroy might cause state to be cleaned up and we have no way to identify the state to be cleaned up. So there is an optional knob to specify state. This can be used to provide terraform tfstate backup file to identify the resources to be cleaned up.
