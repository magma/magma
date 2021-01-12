create database if not exists ${nms_db_user};
create user if not exists '${nms_db_user}'@'localhost' identified by '${nms_db_pass}';
create user if not exists '${nms_db_user}'@'%' identified by '${nms_db_pass}';
GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, INDEX, DROP, ALTER, CREATE TEMPORARY TABLES, LOCK TABLES ON ${nms_db_user}.* TO '${nms_db_user}'@'localhost';
GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, INDEX, DROP, ALTER, CREATE TEMPORARY TABLES, LOCK TABLES ON ${nms_db_user}.* TO '${nms_db_user}'@'%';

create database if not exists ${orc8r_db_user};
create user if not exists '${orc8r_db_user}'@'localhost' identified by '${orc8r_db_pass}';
create user if not exists '${orc8r_db_user}'@'%' identified by '${orc8r_db_pass}';
GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, INDEX, DROP, ALTER, CREATE TEMPORARY TABLES, LOCK TABLES ON ${orc8r_db_user}.* TO '${orc8r_db_user}'@'localhost';
GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, INDEX, DROP, ALTER, CREATE TEMPORARY TABLES, LOCK TABLES ON ${orc8r_db_user}.* TO '${orc8r_db_user}'@'%';
