# Orc8r Deployer

Orc8r deployer is a docker image which contains all the necessary prerequisites to deploy orc8r. This docker image ships with orc8r deployer cli(orcl) which enables a user to configure, run installation prechecks, run upgrade prechecks, run post deployment checks and finally perform deployment cleanup if necessary.
(note: this tool is still in development)

## Requirements

The host machine must have docker engine installed.

## Usage

```
./run_deployer runs the orc8r deployer container
orc8r deployer contains scripts which enable user to configure, run prechecks
install, upgrade, verify and cleanup an orc8r deployment
Usage: run_deployer [-deploy-dir|-root-dir|-build|-h]
options:
-h           Print this Help
-deploy-dir  deployment dir containing configs and secrets (mandatory)
-root-dir    magma root directory
-build       build the deployer container
example: ./run_deployer -deploy-dir /tmp/orc8r_14_deployment
```

Deployment directory initially can either be empty or can contain config files. Once configured, the deployment directory will contain config and secrets directory.

Config directory will contain the terraform config files.

* infra.tfvars.json - configuration
* platform.tfvars.json
* service.tfvars.json

Secrets directory will contain the certificates generated through the cli.

### Orc8r Deployer CLI (orcl)

```
# orcl --help
Usage: orcl [OPTIONS] COMMAND [ARGS]...

  Orchestrator Deployment CLI

Options:
  --version  Show the version and exit.
  --help     Show this message and exit.

Commands:
  configure  Configure enables user to configure values needed for Orc8r...
  install    Install command enables user to run subcommands in context of...
```

**orcl configure**

```
root@de8ebf2c5e0e:~/project# orcl configure --help
Usage: orcl configure [OPTIONS] COMMAND [ARGS]...

  Configure option enables user to manage deployment related configuration.
  It can be used to configure all mandatory configuration values necessary
  for the deployment

Options:
  -c, --component [infra|platform|service]
  --help                          Show this message and exit.

Commands:
  check  Check option enables user to check if mandatory configuration has...
  info   Displays all possible configuration options along with its...
  set    Set enables user to configure any configuration option.
  show   Shows the current configuration
```


**orcl install**

```
root@de8ebf2c5e0e:~/project# orcl install --help
Usage: orcl install [OPTIONS] COMMAND [ARGS]...

  Install command enables user to run subcommands in context of Orc8r
  installation.

Options:
  --help  Show this message and exit.

Commands:
  addcerts  Addcerts which lets user create certificates for the...
  precheck  Precheck which runs various checks to ensure successful...
```

## Work done so far

* Docker image containing all pre reqs necessary for Orc8r deployment
* Configuration tool which enables
    * user to configure infra, platform and services
    * See the local configuration
    * View the detail information about all config knobs available
    * check if all required values are present
* Ability to add secrets (copied existing ansible script from cloudstrapper)
    * option to skip self signed certs
* scripts to perform installation prechecks
    * Infra
        * Check if all mandatory infra values have been configured
        * Check if we have atleast 3 availabiliity zones
        * Check if secrets manager secrets exists
    * Platform
        * Check if all mandatory platform values have been configured
        * Check if cloud watch related log groups are present
    * Service
        * Check if all mandatory service values have been configured
        * Check if Orc8r deployment type is valid
        * Check if helm repo is valid and it contains helm charts of right version and matches the deployment requirements
        * Check if docker repo is valid and if it contains images with expected tags

## Examples

### CLI configure examples:

```
root@f6664a7e4de7:~/project# orcl configure -c platform
Configuring platform deployment variables

nms_db_password[password1234]:
orc8r_db_password[password1234]:

root@1233a8960a8e:~/project# orcl configure set
component to configure (infra, platform, service): infra
name of the variable: cluster_name
value of the variable: foobar

root@f6664a7e4de7:~/project# orcl configure check -c platform
Checking platform deployment variables

root@f6664a7e4de7:~/project# orcl configure show -c platform
platform Configuration
+-------------------+---------------+
| Name              | Configuration |
+-------------------+---------------+
| nms_db_password   | password1234  |
| orc8r_db_password | password1234  |
+-------------------+---------------+

root@f6664a7e4de7:~/project# orcl configure info -c platform
platform Configuration Options
+------------------------------------------+-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------+--------+----------+-----------+
| Name                                     | Description                                                                                                                                                                 | Type   | Required | Component |
+------------------------------------------+-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------+--------+----------+-----------+
| deploy_elasticsearch                     | Flag to deploy AWS Elasticsearch service as the datasink for aggregated logs.                                                                                               | bool   | False    | tf        |
| deploy_elasticsearch_service_linked_role |  Flag to deploy AWS Elasticsearch service linked role with cluster. If you've already created an ES service linked role for another cluster, you should set this to false.  | bool   | False    | tf        |
| efs_project_name                         | Project name for EFS file system                                                                                                                                            | string | False    | tf        |
| elasticsearch_az_count                   | AZ count for ES.                                                                                                                                                            | number | False    | tf        |
| elasticsearch_dedicated_master_count     | Number of dedicated ES master nodes.                                                                                                                                        | number | False    | tf        |
| elasticsearch_dedicated_master_enabled   | Enable/disable dedicated master nodes for ES.                                                                                                                               | bool   | False    | tf        |
| elasticsearch_dedicated_master_type      | Instance type for ES dedicated master nodes.                                                                                                                                | string | False    | tf        |
| elasticsearch_domain_name                | Name for the ES domain.                                                                                                                                                     | string | False    | tf        |
| elasticsearch_domain_tags                | Extra tags for the ES domain.                                                                                                                                               | map    | False    | tf        |
| elasticsearch_ebs_enabled                | Use EBS for ES storage.                                                                                                                                                     | bool   | False    | tf        |
| elasticsearch_ebs_iops                   | IOPS for ES EBS volumes.                                                                                                                                                    | number | False    | tf        |
| elasticsearch_ebs_volume_size            | Size in GB to allocate for ES EBS data volumes.                                                                                                                             | number | False    | tf        |
| elasticsearch_ebs_volume_type            | EBS volume type for ES data volumes.                                                                                                                                        | string | False    | tf        |
| elasticsearch_instance_count             | Number of instances to allocate for ES domain.                                                                                                                              | number | False    | tf        |
| elasticsearch_instance_type              | AWS instance type for ES domain.                                                                                                                                            | string | False    | tf        |
| elasticsearch_version                    | ES version for ES domain.                                                                                                                                                   | string | False    | tf        |
| global_tags                              | n/a                                                                                                                                                                         | map    | False    | tf        |
| nms_db_engine_version                    | MySQL engine version for NMS DB.                                                                                                                                            | string | False    | tf        |
| nms_db_identifier                        | Identifier for the RDS instance for NMS.                                                                                                                                    | string | False    | tf        |
| nms_db_instance_class                    | RDS instance type for NMS DB.                                                                                                                                               | string | False    | tf        |
| nms_db_name                              | DB name for NMS RDS instance.                                                                                                                                               | string | False    | tf        |
| nms_db_password                          | Password for the NMS DB.                                                                                                                                                    | string | True     | tf        |
| nms_db_storage_gb                        | Capacity in GB to allocate for NMS RDS instance.                                                                                                                            | number | False    | tf        |
| nms_db_username                          | Username for default DB user for NMS DB.                                                                                                                                    | string | False    | tf        |
| orc8r_db_engine_version                  | Postgres engine version for Orchestrator DB.                                                                                                                                | string | False    | tf        |
| orc8r_db_identifier                      | Identifier for the RDS instance for Orchestrator.                                                                                                                           | string | False    | tf        |
| orc8r_db_instance_class                  | RDS instance type for Orchestrator DB.                                                                                                                                      | string | False    | tf        |
| orc8r_db_name                            | DB name for Orchestrator RDS instance.                                                                                                                                      | string | False    | tf        |
| orc8r_db_password                        | Password for the Orchestrator DB.                                                                                                                                           | string | True     | tf        |
| orc8r_db_storage_gb                      | Capacity in GB to allocate for Orchestrator RDS instance.                                                                                                                   | number | False    | tf        |
| orc8r_db_username                        | Username for default DB user for Orchestrator DB.                                                                                                                           | string | False    | tf        |
+------------------------------------------+-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------+--------+----------+-----------+


```



### CLI install example:

```
root@f6664a7e4de7:~/project# orcl install precheck
[WARNING]: No inventory was parsed, only implicit localhost is available
[WARNING]: provided hosts list is empty, only localhost is available. Note that
the implicit localhost does not match 'all'
['No config file found; using defaults',
 '',
 'PLAY [localhost] '
 '***************************************************************',
 '',
 'TASK [infra : check if required infra variables are set] '
 '***********************',
 'changed: [localhost] => {"changed": true, "cmd": ["orcl", "configure", '
 '"check", "-c", "infra"], "delta": "0:00:00.155063", "end": "2021-03-26 '
 '02:30:14.708979", "failed_when_result": false, "rc": 0, "start": "2021-03-26 '
 '02:30:14.553916", "stderr": "", "stderr_lines": [], "stdout": "Checking '
 'infra deployment variables", "stdout_lines": ["Checking infra deployment '
 'variables"]}',
 '',
 'TASK [infra : Open configuration file] '
 '*****************************************',
 'changed: [localhost] => {"changed": true, "cmd": "cat '
 '/root/project/configs/infra.tfvars.json", "delta": "0:00:00.004723", "end": '
 '"2021-03-26 02:30:14.903529", "rc": 0, "start": "2021-03-26 '

 root@f6664a7e4de7:~/project# orcl install addcerts
[WARNING]: No inventory was parsed, only implicit localhost is available
[WARNING]: provided hosts list is empty, only localhost is available. Note that
the implicit localhost does not match 'all'
['No config file found; using defaults',
```

## Terraform Installation

Currently terraform installation is still manual. All configuration values are pushed into ‘terraform.tfvars.json’ which is automatically loaded by terraform.

```
cd /root/project

# initialize terraform
terraform init

# bring up infra and platform resources
`terraform apply ``-``target``=``module``.``orc8r`
`export`` KUBECONFIG``=``$``(``realpath kubeconfig_orc8r``)`

`# seed the certificates`
`terraform apply ``-``target``=``module``.``orc8r``-``app``.``null_resource``.``orc8r_seed_secrets `

# bring up the orc8r services
`terraform apply`

```

## Tasks Remaining

* Adding upgrade prechecks
* Adding post deployment checks
* Adding deployment cleanup
* Automating installation

## Note:

* vars.yml - currently a static file which holds all the configuration variables. It contains variables to be configured for AWS and terraform. It could be generated when we are able to parse the tf files to JSON. That’s probably the next step once the tool changes are complete.
* main.tf.j2 and vars.tf.j2 are the jinja templates for main terraform and variable terraform file. It is mainly templatized to include the configurable variables added.



