---
id: version-1.6.0-RDS_upgrade
title: Orchestrator DB Upgrade
hide_title: true
original_id: RDS_upgrade
---

# Summary 

AWS will deprecate all RDS instances that are running on 9.6.X version as of January 2022. In preparation for that, the following instructions will allow a network administrator to upgrade the RDS version to Version 12.8. 

> **_NOTE:_** This upgrade procedure is mandatory for all RDS versions running Postgres 9.6.X version. Versions 10 and above are not mandatory at this time.*

## Pre-Requisites 

1. Access to AWS Dashboard (read/write). 

2. Permissions to modify the orc8r deployment via Terraform

3. Access to k8’s cluster and the kubeconfig file. 

4. A [_new snapshot_](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_CreateSnapshot.html) created in case a database recovery is needed.

## Downtime:

This operation is going to cause noticeable downtime. During this time the NMS will not be reachable and new subscribers cannot be added to gateways. If a public cloud orc8r instance is being used for Federated setups, attach and policy management will be impacted. 

Depending on the starting version, up to two upgrades might be necessary. Based on our testing, this can take up to an hour (including validation). It is recommended to schedule this work in a maintenance window.


## Execution:

Three minor modifications are necessary to the terraform files. Please note that these modifications were committed to the Magma Repo in [_PR 10602_](https://github.com/magma/magma/pull/10602). 

Please note that these changes need to be made to the local copy of the modules downloaded.

**Change 1:** 
In `.terraform/modules/orc8r/orc8r/cloud/deploy/terraform/orc8r-aws/db.tf` add the following snippet: 

`apply_immediately = var.orc8r_db_apply_immediately`

**Change 2:**
In `.terraform/modules/orc8r/orc8r/cloud/deploy/terraform/orc8r-aws/variables.tf` add the following snippet:

```
variable "orc8r_db_apply_immediately"{
  description = "Flag to immediately upgrade RDS without waiting for a maintenance window"
  type = bool
  default = false
}
```

**Change 3:**
In `main.tf` set the following variables in `module "orc8r" {...}` stanza:

`orc8r_db_engine_version = "12.8"`
`orc8r_db_apply_immediately = "true"`


> **_NOTE:_** If the current version of the `orc8r_db_engine_version` is _not_ `9.6.23` (i.e. a minor version *older* than 9.6.23), the RDS has to be first upgraded to this version (using the same process) and then a second time to `12.8`



> **_NOTE_**: `terraform init --upgrade` should be run once at the start of this activity. This will ensure that the latest changes from Github are downloaded to the local version on your workstation. `terraform init --upgrade` should NOT be run between upgrades, else all the local changes will be overwritten. 


After confirming the notes above, run the following commands:

* `terraform plan`
* `terraform apply`

## Logs and Validation 

* Each time the RDS instance is upgraded, please check the `terraform plan` output. Each time, it should explicitly specify that the instance will be *updated in place*. 
* Once satisfied that the RDS instance will be updated in place, you can run the `terraform apply` command.
* Logs from an RDS upgrade (9.6.22 -> 9.6.23 -> 12.8) have been captured [_here_](https://gist.github.com/sudhikan/8e42985bc0db13512c9cd602d8acab3a) for reference.
* Between each upgrade step, please verify that you are able to login to the NMS, access the subscriber page, modify values, and gateways are successfully checking-in. 
* It is also recommended to take snapshots before each upgrade. Refer to the link in the prerequisites section for more information.

## Cleanup 

* Modify main.tf as follows:

`orc8r_db_apply_immediately = "false"`

* Run `terraform plan` and `terraform apply`. These are not expected to have any impact other than to disable immediate major upgrades of the RDS instance.

```
An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # module.orc8r.aws_db_instance.default will be updated in-place
  ~ resource "aws_db_instance" "default" {
      ~ apply_immediately                     = true -> false
        id                                    = "orc8rdb"
        name                                  = "orc8r"
        tags                                  = {}
        # (48 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.
```

* After terraform has finished running, run `terraform show` to validate the DB settings.

```
# module.orc8r.aws_db_instance.default:
resource "aws_db_instance" "default" {
    address                               = "orc8rdb.abcdy9gi9jk5o.eu-west-2.rds.amazonaws.com"
    allocated_storage                     = 64
    allow_major_version_upgrade           = true
    apply_immediately                     = false
    arn                                   = "arn:aws:rds:eu-west-2:************:db:orc8rdb"
    auto_minor_version_upgrade            = true
    availability_zone                     = "eu-west-2a"
    backup_retention_period               = 7
    backup_window                         = "01:00-01:30"
    ca_cert_identifier                    = "rds-ca-2019"
    copy_tags_to_snapshot                 = false
    customer_owned_ip_enabled             = false
    db_subnet_group_name                  = "orc8r"
    delete_automated_backups              = true
    deletion_protection                   = false
    enabled_cloudwatch_logs_exports       = []
    endpoint                              = "orc8rdb.abcdy9gi9jk5o.eu-west-2.rds.amazonaws.com:5432"
    engine                                = "postgres"
    engine_version                        = "12.8"
    engine_version_actual                 = "12.8"
```

## Backout

* If for some reason the DB upgrades fail, please [_restore_](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_RestoreFromSnapshot.html) your database from the snapshot created in the prerequisites section.

## Screenshots

![nms_rds_upgrade](assets/orc8r/rds_upgrade_nms.png)
> **_NOTE:_** Expected output when accessing the NMS during the upgrade. Errors are expected on screen.

![nms_upgrade_aws](assets/orc8r/rds_upgrade_aws_9_6_22.png)
> **_NOTE:_** Upgrade from 9.6.22 to 9.6.23. Per AWS’s RDS upgrade guide, jumping to 12.8 can only be done from 9.6.23.

![rds_upgrade_complete](assets/orc8r/rds_upgrade_aws_complete.png)
> **_NOTE:_** This screenshot displays a successful upgrade. 
