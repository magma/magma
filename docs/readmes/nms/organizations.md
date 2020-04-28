---
id: nms_organizations
title: NMS Organizations
hide_title: true
---

# NMS Organizations

In version 1.1 of Magma, the NMS includes the concept of “Organizations”, or segmentations of networks which controls what networks an organization can access. Organizations are controlled via the `master` organization, which is a special org accessible from `master.<nms-hostname>` once you have created a user for that org (see: first-time setup)

### First-time Setup

When you deploy the nms for the first time, you’ll need to create a user that has access to the master organization, run the command

* Docker:
    * `docker-compose exec magmalte yarn setAdminPassword master <email> <password>`
* Kubernetes:
    * `kubectl exec <magmalte-container> -- yarn setAdminPassword master <email> <password>`

### Examples


Single-tenant: Create one organization and give it access to all networks

![Org with access to all networks](assets/nms/org_all_networks.png)

* This is essentially the same as before. The only difference is that the NMS is accessible from the URL `<organization-name>.<hostname>`

Multiple Tenants
* Create a second organization and give it access to specific networks

![List of organizations](assets/nms/org_multiple_list.png)

* Here, fb-test has access to all networks, while magma-test only has access to the network `mpk_test`
* Create a user in this organization to use it

![Add user to org](assets/nms/org_add_user.png)

* When you log in to `magma-test.<hostname>` you will only be able to see the network `mpk_test`, however if you log into `fb-test.<hostname>`, you will have access to all networks



### Migration Details

Organizations only carry information about which networks are accessible to which user, so upgrading your deployment should require minimal migration work. If you want to remain with a single-tenant deployment simply create an organization that has access to all networks, add users to this org, and now everything works just as before. All configurations for the pre-existing networks will still be there.
