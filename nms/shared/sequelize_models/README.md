# sequelize-models

Uses the [Sequelize ORM](https://sequelize.org/)
and defines various models commonly used in different NMS products.

## Models

**Organization**

Access control to NMS products are separated based on organizations.

For example, in Magma, organizations can be assigned access to one or more
networks.

**User**

Users exist per-organization, and additional access control can be set for
individual users.

**Audit Log Entry**

A log of all actions taken by users on the NMS.

**Feature Flag**

## Yarn Commands

### dbDataMigration Usage

Used for migration of sequelize-models data from one DB to another

**Example: Manual Usage**

```
$ yarn dbDataMigrate

? Enter DB host: mariadb
? Enter DB port: 3306
? Enter DB database name: nms
? Enter DB username: root
? Enter DB password: [hidden]
? Enter DB SQL dialect: mariadb

DB Connection Config:
---------------------------
Host: mariadb:3306
Database: nms
Username: root
Dialect: mariadb

? Are you importing from the specified DB, or exporting to it?: import
? Would you like to run data migration with these settings?: Yes
Completed data migration, importing from specified DB
```

**Example: Automated Usage**

```
$ npm dbDataMigrate -- --username=nms --password=nms --database=nms --host=mariadb --port=3306 --dialect=mariadb --export --confirm

DB Connection Config:
---------------------------
Host: mariadb:3306
Database: nms
Username: nms
Dialect: mariadb

Completed data migration, exporting to specified DB
```
