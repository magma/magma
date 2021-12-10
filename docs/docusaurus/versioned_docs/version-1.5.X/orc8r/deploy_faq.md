---
id: version-1.5.0-deploy_faq
title: FAQs
hide_title: true
original_id: deploy_faq
---

# Deployment FAQs

## Terraform timed out when running `terraform apply`

https://magmacore.slack.com/archives/C018J8UMGMR/p1599228643121500

```
Error: rpc error: code = Unknown desc = release orc8r failed: timed out waiting for the condition on ../../main.tf line 22, in resource "helm_release" "orc8r":
 22: resource "helm_release" "orc8r" {
```

**Resolution steps**

- Check if there is an issue with Helm repo URL, password, or container image tag
- Ensure container image registry and Terraform values are correct
- Run following kubectl commands to get more details on the error encountered

```
kubectl --namespace orc8r get pods
kubectl --namespace orc8r describe pods
# pod status ImagePullBackOff, indicates that Image wasn't found
# pod can also be crash looping, pod status will provide this information,
# in this case, looking at the logs of individual container will help debug this further
```

- Get Helm's view of images

```
helm version (provides current version of helm)
helm -n orc8r list (provides the list of releases which are currently deployed under orc8r namespace)
helm -n orc8r get values orc8r (gets the values file for the orc8r release)
```

## NMS pod wasn't showing up properly

```
kubectl exec -it $(kubectl get pod -l app.kubernetes.io/component=magmalte -o jsonpath='{.items[0].metadata.name}') -- yarn setAdminPassword master xxxx@xxxx.com 1234
error: error executing jsonpath "{.items[0].metadata.name}": Error executing template: array index out of bounds: index 0, length 0. Printing more information for debugging the template:
        template was:
                {.items[0].metadata.name}
        object given to jsonpath engine was:
                map[string]interface {}{"apiVersion":"v1", "items":[]interface {}{}, "kind":"List", "metadata":map[string]interface {}{"resourceVersion":"", "selfLink":""}}
error: pod, type/name or --filename must be specified
```

**Resolution steps**

- The problem was that `deploy_nms` wasn't set to true in variables.tf. Set `deploy_nms` to true and rerun `terraform apply`

## No Resources found in default namespace after terraform apply

```
No resources found in default namespace after terraform apply
```

**Resolution steps**

- Terraform seems to reset the kubectl namespace — fix using `kubens` to select the orc8r namespace

## Errors creating Secrets Manager

```
error creating Secrets Manager Secret: InvalidRequestException:
You can't create this secret because a secret with this name is already
 scheduled for deletion
```

**Resolution steps**

- Change `secretsmanager_orc8r_secret` variable to new value
- AWS doesn't fully delete secrets for safety reasons, so deleted secrets need a name change before successful regeneration
- If you want to fully, "immediately" delete the secret (danger zone!), consider `aws secretsmanager delete-secret --secret-id your_secret_id --force-delete-without-recovery` [(ref)](https://docs.aws.amazon.com/secretsmanager/latest/userguide/manage_delete-restore-secret.html)

## Error: ValidationException: Domain is being deleted for aws_elasticsearch_domain_policy

```
`Error: ValidationException: Domain is being deleted` for `aws_elasticsearch_domain_policy`

```

**Resolution steps**

- Sometimes ES domains aren't deleted immediately. Change ES domain to new value and stop tracking old domain
- Ensure ES domain is actually marked for deletion, e.g. via AWS console
- Change `elasticsearch_domain_name` variable to new value
- `terraform state rm module.orc8r.aws_elasticsearch_domain.es[0]` stop tracking old ES domain

## Terraform apply failed due to some Python module error

```
Warning: Applied changes may be incompleteThe plan was created with the -target option in effect, so some changes
requested in the configuration may have been ignored and the output values may
not be fully updated. Run the following command to verify that no other
changes are pending:
 terraform planNote that the -target option is not suitable for routine use, and is provided
only for exceptional situations such as recovering from errors or mistakes, or
when Terraform specifically suggests to use it as part of an error message.Error: Error running command ' ../../scripts/create_orc8r_secrets.py \
 'orc8r-secrets' 'us-west-2' \
 "~/secrets/certs"
': exit status 1. Output: Traceback (most recent call last):
 File "../../scripts/create_orc8r_secrets.py", line 12, in <module>
 import boto3
ModuleNotFoundError: No module named 'boto3'
```

**Resolution steps**

- Identify and install the relevant Python module - `pip3 install boto3`
- Specifically ensure if deployment specific tooling has been installed

```
brew install aws-iam-authenticator kubectl helm terraform
python3 -m pip install awscli boto3
aws configure
```
