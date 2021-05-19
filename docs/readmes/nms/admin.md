---
id: admin
title: Admin
hide_title: true
---
# Administration

NMS provides an admin page to manage the NMS itself. Admin page currently has three main functionalities

* Manage Users
* Manage Networks
* View Audit Logs

## Users

New users to NMS can be added or edited through the admin page.
![admin](assets/nms/userguide/admin/admin.png)
![users](assets/nms/userguide/admin/users.png)

Users can be configured as either a superuser/user/readonly user.
Super User will have capability to view all the metrics in the organization through Grafana.
User will have the capability to view only the metrics for assigned network in the organization.
Read Only user can only view the various components in NMS. They cannot make any edits to the underlying configuration.

## Audit Logs

All configuration change made through NMS can be viewed through the audit log component. Audit log component comes with a convenient interface to search, sort and view the entire JSON payload of the configuration change.
![audit_log](assets/nms/userguide/admin/audit_log1.png)
![audit_log](assets/nms/userguide/admin/audit_log2.png)
![audit_log](assets/nms/userguide/admin/audit_log3.png)


## Networks

Admin page also contains a network page to add any network to the given organization.
![network](assets/nms/userguide/admin/network.png)
