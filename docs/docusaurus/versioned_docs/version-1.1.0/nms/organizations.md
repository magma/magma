---
id: version-1.1.0-nms_organizations
title: Multi-Tenancy (Organizations)
hide_title: true
original_id: nms_organizations
---

# NMS Multi-Tenancy (Organizations)

We introduced multitenancy support in the Magma LTE NMS starting from v1.1.0.
Tentants in the Magma NMS are called "Organizations". Each organization owns
a subset of the networks provisioned on Orchestrator and the special `master`
organization administrates organizations in the system.

Users in organizations log into the NMS using a subdomain that matches their
organization name. For example, users of a `facebook` organization in the NMS
would access the NMS using `facebook.nms.yourdomain.com`.

## First-time Setup

When you deploy the nms for the first time, you’ll need to create a user that
has access to the master organization. Run the command

* Docker:
    * `docker-compose exec magmalte yarn setAdminPassword master <email> <password>`
* Kubernetes:
    * `kubectl exec -it \
          $(kubectl get pod -l app.kubernetes.io/component=magmalte -o jsonpath='{.items[0].metadata.name}') -- \
          yarn setAdminPassword master <admin user email> <admin user password>`
          
## DNS Resolution

We use `external-dns` to automatically set up an AWS Route53 DNS record that
points `*.nms.youdomain.com` to the NMS application. If you're managing your
domain name outside of Route53, you'll have to add an NS record `<org>.nms.`
for every new organization you add to the NMS. The list of nameservers to set
can be found in the AWS console for the Route53 zone or as the `nameservers`
output of the `orc8r-aws` Terraform module.

## Examples

### Single-tenant

Create one organization and give it access to all networks. This is essentially
the same as v1.0 when there was no tenancy support. The only difference is that
the NMS is accessible from the URL `magma-test.youdomain.com`

![Org with access to all networks](assets/nms/org_all_networks.png)

### Multiple Tenants

Create a second organization and give it access to specific networks:

![List of organizations](assets/nms/org_multiple_list.png)

Here, fb-test has access to all networks, while magma-test only has access
to the network `mpk_test`. Create a user in this organization to use it:

![Add user to org](assets/nms/org_add_user.png)

When you log in to `magma-test.youdomain.com` you will only be able to see the
network `mpk_test`. If you log into `fb-test.yourdomain.com`, you will
have access to all networks.
