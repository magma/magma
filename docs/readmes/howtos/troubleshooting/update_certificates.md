---
id: update_certificates
title: Update rootCA and controller SSL certificates
hide_title: true
---
# Update rootCA and controller SSL certificates

**Description:** This document describes the steps to update certificates `controller.crt`, `controller.key`, `rootCA.pem` and `rootCA.key` on Orchestrator/Access Gateway. These steps should follow on an Orc8r already deployed where is required to extend the expiration date of the certificates. This shouldn't be use for change of CN.

**Environment:** Orchestrator in Kubernetes/AWS

**Affected components:** AGW, Orchestrator

**Configuration steps on Orchestrator:**


1. Create a new rootCA.pem, rootCA.key, controller.crt and controller.key.

    - The public SSL certificate for your Orchestrator domain,
    with `CN=*.yourdomain.com`. This can be an SSL certificate chain, but it must be
    in one file

    - The private key which corresponds to the above SSL certificate

    - The root CA certificate which verifies your SSL certificate

    If you aren't worried about a browser warning, you can generate self-signed
    versions of these certs

    ```bash
    ${MAGMA_ROOT}/orc8r/cloud/deploy/scripts/self_sign_certs.sh yourdomain.com
    ```

    Alternatively, if you already have these certs, rename and move them as follows

    - Rename your public SSL certificate to `controller.crt`

    - Rename your SSL certificate's private key to `controller.key`

    - Rename your SSL certificate's root CA certificate to `rootCA.pem`


2. Move these certificates into `~/secrets/certs`.

3. Other certificates `bootstrapper.key`, `certifier.key`, `certifier.pem`, `fluentd.key`, `fluentd.pem`, `admin_operator.key.pem,` `admin_operator.pem` and `admin_operator.pfx`  **don't** need to change and you can copy into the folder `~/secrets/certs`. The certs directory should now look like this

```bash
$ ls -1 ~/secrets/certs/

admin_operator.key.pem
admin_operator.pem
admin_operator.pfx
bootstrapper.key
certifier.key
certifier.pem
controller.crt
controller.key
fluentd.key
fluentd.pem
rootCA.pem
rootCA.key
```


4. Opt to upgrade modules and plugins as part of their respective installation steps.

`terraform init –upgrade`

5. Run the following commands:

`$ terraform taint module.orc8r-app.null_resource.orc8r_seed_secrets`

`$ terraform apply -target=module.orc8rapp.null_resource.orc8r_seed_secrets`

`$ terraform apply`


NOTE: `terraform apply` command outputs the “plan” of what it intends to add,destroy,modify. Please scrutinize this output before typing “yes” on the confirm prompt. If there are any changes that are not consistent with your expectations, please cancel the run. You can specifically target the secrets portion by doing `terraform apply -target=module.<module information>`

6. You can remove the secrets folder from your local disk. But make sure you have a copy of the `rootCA.key`

7. Kill Controllers and Proxy pods one by one.

`kubectl delete pods <pod>`

**Configuration steps on Access Gateway:**

1. Update the new `rootCA.pem` in `/var/opt/magma/tmp/certs/rootCA.pem`

2. Restart magmad service

`AGW$ sudo service magma@magmad restart`

3. You can validate the connection between your AGW and Orchestrator:

```bash
AGW$ journalctl -u magma@magmad -f
# Look for the following logs
# INFO:root:Checkin Successful!
# INFO:root:[SyncRPC] Got heartBeat from cloud
# INFO:root:Processing config update gateway_id

AGW$ sudo checkin_cli.py

1. -- Testing TCP connection to controller-staging.magma.etagecom.io:443 --
2. -- Testing Certificate --
3. -- Testing SSL --
4. -- Creating direct cloud checkin --
5. -- Creating proxy cloud checkin --
```
