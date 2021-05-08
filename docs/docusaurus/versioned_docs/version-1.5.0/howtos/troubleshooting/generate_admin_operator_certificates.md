---
id: version-1.5.0-generate_admin_operator_certificates
title: Generate and update admin_operator certificates
hide_title: true
original_id: generate_admin_operator_certificates
---
# Generate and update admin_operator certificates

**Description:**  NMS is unable to communicate with the controller or API access is no longer functional due to certificate used for the TLS handshake has expired. Below guide provide the steps to generate new admin_operator certificates and upload it to the NMS k8s pod for NMS access and to the browser for the API access.

**Environment:** Orchestrator, NMS in Kubernetes/AWS

**Affected components:** Orchestrator, NMS

**Configuration steps:**


1. Log on to a host that has kubectl access to your orc8r cluster.
2. Access the shell on your controller pod.

`export CNTLR_POD=$(kubectl -n orc8r get pod -l app.kubernetes.io/component=controller -o jsonpath='{.items[0].metadata.name}')`

`kubectl exec -it ${CNTLR_POD} bash`

3. Run the following commands inside the k8s controller pod to generate a new `admin_operator` cert. The `admin_operator.pfx` file is only for API access for the user to get to the browser.

```
# The following commands are to be run inside the pod
(pod)$ cd /var/opt/magma/bin
(pod)$ envdir /var/opt/magma/envdir ./accessc add-admin -cert admin_operator admin_operator
(pod)$ openssl pkcs12 -export -out admin_operator.pfx -inkey admin_operator.key.pem -in admin_operator.pem

Enter Export Password:
Verifying - Enter Export Password:
```

4. Copy the k8s certs from the controller k8s pod to the local directoy where all your secrets are held.

NOTE: Make sure the location where these certificates are being copied to is the same location which is referenced in the main.tf file for the seed_certs_dir variable.
```
cd ~/secrets/certs
for certfile in admin_operator.pem admin_operator.key.pem admin_operator.pfx
do
    kubectl cp orc8r/${CNTLR_POD}:/var/opt/magma/bin/${certfile} ./${certfile}
done
```
5. At this point the new certs have been generated, you can use the `admin_operator.pfx` to validate that you can reach the API endpoint. However, they still need to be applied to the k8s secrets manager so that they can be mounted in the nms pod. To do this, initialize your terraform first using terraform init -upgrade.

6. To replace certs, we have to first taint the secrets in terraform so that terraform knows to destroy those secrets first and then re-apply them. Taint them using terraform `taint module.orc8r-app.null_resource.orc8r_seed_secrets`

7. Once the secrets have been tainted, you can then go ahead and apply the new secrets by running `terraform apply -target=module.orc8rapp.null_resource.orc8r_seed_secrets` followed by `terraform apply`

NOTE: `terraform apply` command outputs the “plan” of what it intends to add,destroy,modify. Please scrutinize this output before typing “yes” on the confirm prompt. If there are any changes that are not consistent with your expectations, please cancel the run. You can specifically target the secrets portion by doing `terraform apply -target=module.<module information>`

8. Once the terraform apply succeeds, the NMS will not automatically get the new certificates until we destroy the existing pod and force the replication controller to instantiate another instance of the pod. To do this, run the following commands:
```
export NMS_POD=$(kubectl -n orc8r get pod -l  app.kubernetes.io/component=magmalte -o jsonpath='{.items[0].metadata.name}')
kubectl -n orc8r delete pod NMS_POD
```
9. Once the pod reinitializes, it should have the latest admin_operator certs and be able to re-establish communication with the controller.
