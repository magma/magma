Deploying Magma Orc8r on bare metal

NOTE: This deployment is only working for 1.3.x versions. 1.4 support will be
added soon.

The following file needs to be edited before deployment can start:
* ansible_vars.yaml

The values that need to be customized are the IP settings for your network and
the passwords which need to be generated. Additionally, you need to host a
docker repo and helm chart repo for orc8r. Further details can be found at
https://magma.github.io/magma/docs/orc8r/deploy_build#build-and-publish-helm-charts

The deployment assumes the host running Ansible and the target hosts are
running Ubuntu 18.04 or 20.04.

The following variables must be set before starting:
* orc8r_image_repo - Docker repo where controller, magmalte, and nginx containers are located
* orc8r_helm_repo - Helm chart repo where orc8r chart is located
* orc8r_domain - DNS domain where Magma services are published
* orc8r_nms_admin_email - Initial login user for NMS
* orc8r_chart_version - Helm chart version for orc8r
* orc8r_image_tag - orc8r controller Docker image tag
* orc8r_nms_image_tag - NMS magmalte Docker image tag
* orc8r_nginx_image_tag - orc8r nginx Docker image tag
* metallb_addresses - An IP range of at least 5 IP addresses for LoadBalancers for Magma

If your environment requires it, the following variables should be set:
* ansible_user - ssh user for your k8s hosts
* loadbalancer_apiserver - contains port and address where to set a loadbalancer
* vrrp_nic - network interface which keepalived should bind to on all hosts
* docker_insecure_registries - If your orc8r_image_repo is insecure, it should be configured here
* db_root_password - MariaDB root password (randomly generated if not specified)
* orc8r_db_pass - orc8r DB password (randomly generated if not specified)
* orc8r_nms_db_pass - nms DB password (randomly generated if not specified)
* orc8r_nms_admin_pass - password for NMS initial login user (randomly generated if not specified)

Once the config is set, just run the following command:

./deploy.sh

If the deployment succeeds, you will see information on how to log into Magma
web UI.

Note: external-dns configuration is not automatic here because it is intended
for on-premise. For DNS to work, you have three options:
* Configure external-dns service
* Update your DNS manually
* Locally update /etc/hosts
