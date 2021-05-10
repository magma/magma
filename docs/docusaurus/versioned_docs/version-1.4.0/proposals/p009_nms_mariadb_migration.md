---
id: version-1.4.0-p009_nms_mariadb_migration
title: NMS MariaDB Migration
hide_title: true
original_id: p009_nms_mariadb_migration
---

# Overview

*Status: Accepted*\
*Author: @andreilee*\
*Last Updated: 03/14*\
*Targeted Release: 1.5*

This document concerns:

1. Migration of NMS db from MariaDB to to Postgres to share with orc8r
2. Migration Process

## Goals

Currently the Magma NMS and orc8r are running on separate database
technologies. The orc8r uses postgres for model storage, whereas the NMS uses
MariaDB for storage of users and organizations.

There are no particular benefits for NMS to use MariaDB over Postgres
(or vice-versa), and the continued usage of MariaDB incurs additional costs
for network operators.
Particularly, for network operators using Magma on AWS, there are operating
costs to running both MariaDB and Postgres, over using just one or the other.
The goal of this migration is to have both NMS and orc8r use the same DB 
technology to reduce operating expenses of running a network on Magma.

While orc8r does also use other technologies such as Prometheus and
Elasticsearch, this document does not concern costs associated to such
technologies.

## Success Criteria

Several success criteria 

- Migration script only required to run once after upgrade to 1.5
- Networks running on 1.5 do not have expenses related to running MariaDB
- Migration process works in the three scenarios:
    - with bare Docker only
    - with Terraform and Kubernetes
    
## Current Setup

#### NMS Storage

The NMS for the most part only consumes the orc8r REST API, and few data
objects are interacted with without using the orc8r API.
This data includes only the following:
- Organizations
- NMS Users
- NMS Feature Flags
- NMS Audit Log Entries

To store this NMS-native data in MariaDB, the [Sequelize](https://sequelize.org/master/)
ORM is used. Sequelize is configured through environment variables to use
MariaDB, though Sequelize is DB agnostic, and can be configured to use
Postgres.

#### Terraform

Terraform specifies a single AWS DB instance for Postgres, and another AWS
DB instance for MariaDB on NMS.

## Proposed End Goal Setup

#### NMS Storage

For NMS-native data, Sequelize will still be used, but will be configured to
use Postgres. A separate logical Postgres DB will be used for namespacing
purposes between orc8r data and NMS data.

#### Terraform

A single AWS DB instance will be specified for Postgres, with two logical DBs.

References:
- [Terraform: AWS DB Resource Specification](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_instance)
- [AWS FAQ: How many databases or schemas can I run within a DB instance?](https://aws.amazon.com/rds/faqs/)
- [StackOverflow: Using Terraform for multiple logical DBs](https://stackoverflow.com/questions/52542244/provision-multiple-logical-databases-with-terraform-on-aws-rds-cluster-instance)


## Migration Process

This section is split up between the process that operators must follow to
successfully migrate NMS off MariaDB, and the technical details to achieve
this migration.

#### Developer Process - Docker 

Migration will not be supported for developers.

Rebuilding and restarting docker containers for NMS will simply start up
Postgres instead of MariaDB.

#### Operator Process - Kubernetes Without Terraform

Unsupported.

#### Operator Process - Kubernetes With Terraform

1. Upgrade orc8r and NMS version to 1.5
2. Mark terraform meta-argument `run_mariadb` to `true`
3. Terraform
4. Run migration script to migrate MariaDB data to Postgres
5. Mark terraform meta-argument `run_mariadb` to `false1

#### Technical Details

**Migration Script**

NMS interacts with the DB through an ORM called Sequelize. To migrate data
from MariaDB to Postgres, two instances of the Sequelize object will be
created. One instance will read all data from MariaDB, and the second instance
will write all data to Postgres.

References:
- [StackOverflow: Sequelize Using Multiple Databases](https://stackoverflow.com/questions/37078970/sequelize-using-multiple-databases)

**Update to FBC sequelize-models**

NMS uses the sequelize-models dependency for its [Sequelize models](https://sequelize.org/master/manual/model-basics.html).
This dependency is only capable of opening one DB connection at a time, with
the connection specified through environment variables. A pull request will
be upstreamed to this dependency to remove this limitation so that two
Sequelize connections can be opened at a time for migration off of MariaDB.

**Terraform Changes**

To allow the migration off of MariaDB to proceed in 1.5, but not incur
significant operational costs from running MariaDB, the `aws_db_instance` for
running MariaDB will be optional in Terraform.

To achieve this with Terraform, meta-arguments will be used.
An `aws_db_instance` resource will be specified in Terraform for MariaDB which
will be optional, based on the specified meta-argument.

Additionally, to eliminate namespace collisions between orc8r table names and
NMS-specific table names in Postgres, the same `db_instance` will be used, but
two logical databases will be used. 

MariaDB is required to run for the migration to proceed, but afterwards can be
turned down. To achieve this with Terraform, meta-arguments will be used.
An `aws_db_instance` resource will be specified in Terraform for MariaDB which
will be optional, based on the specified meta-argument.

References:
- [Conditional Resources in Terraform](https://dev.to/tbetous/how-to-make-conditionnal-resources-in-terraform-440n)
- [Terraform: count Meta-Argument](https://www.terraform.io/docs/language/meta-arguments/count.html)
- [Terraform: db_instance Resource](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_instance)
