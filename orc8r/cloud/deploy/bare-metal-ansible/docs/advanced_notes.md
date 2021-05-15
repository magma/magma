Advanced deployment with Ansible

deploy.sh script has optional deployment modes for skipping Kubespray
Kubernetes deployment. Here are the options:

* ./deploy.sh magma-infra
* ./deploy.sh magma

magma-infra.yml playbook deploys nfs-server-provisioner and metallb charts.
Optionally, it could deploy kubevirt if you want to set up a virtualized
(unsupported) AGW as a VM.

magma.yml deploys the logging stack, DB, and then orc8r itself. Individual
components can be enabled/disabled in ansible_vars.yaml in the Advanced options
section.

This project is still quite biased in terms of allowing users to bring their
own DMBS and ElasticSearch and/or Fluentd, but it should be made more flexible
in the nearest time.
