################################################################################
# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
################################################################################

module orc8r {
  source = "../../../orc8r-aws"

  region = "us-west-2"

  nms_db_password             = "mypassword"
  orc8r_db_password           = "mypassword"
  secretsmanager_orc8r_secret = "orc8r-secrets"
  orc8r_domain_name           = "orc8r.example.com"

  vpc_name     = "orc8r"
  cluster_name = "orc8r"

  deploy_elasticsearch          = true
  elasticsearch_domain_name     = "orc8r-es"
  elasticsearch_version         = "7.1"
  elasticsearch_instance_type   = "t2.medium.elasticsearch"
  elasticsearch_instance_count  = 2
  elasticsearch_az_count        = 2
  elasticsearch_ebs_enabled     = true
  elasticsearch_ebs_volume_size = 32
  elasticsearch_ebs_volume_type = "gp2"
}

module orc8r-app {
  source = "../.."

  region = "us-west-2"

  orc8r_domain_name     = module.orc8r.orc8r_domain_name
  orc8r_route53_zone_id = module.orc8r.route53_zone_id
  external_dns_role_arn = module.orc8r.external_dns_role_arn

  secretsmanager_orc8r_name = module.orc8r.secretsmanager_secret_name
  seed_certs_dir            = "~/orc8r.test.secrets/certs"

  orc8r_db_host = module.orc8r.orc8r_db_host
  orc8r_db_name = module.orc8r.orc8r_db_name
  orc8r_db_user = module.orc8r.orc8r_db_user
  orc8r_db_pass = module.orc8r.nms_db_pass

  nms_db_host = module.orc8r.nms_db_host
  nms_db_name = module.orc8r.nms_db_name
  nms_db_user = module.orc8r.nms_db_user
  nms_db_pass = module.orc8r.orc8r_db_pass

  docker_registry = "registry.hub.docker.com/foobar"
  docker_user     = "foobar"
  docker_pass     = "mypassword"

  helm_repo = "example.jfrog.io"
  helm_user = "foobar"
  helm_pass = "mypassword"

  eks_cluster_id = module.orc8r.eks_cluster_id

  efs_file_system_id       = module.orc8r.efs_file_system_id
  efs_provisioner_role_arn = module.orc8r.efs_provisioner_role_arn

  elasticsearch_endpoint = module.orc8r.es_endpoint

  orc8r_chart_version = "1.4.7"
  orc8r_tag           = "1.0.1"
  deploy_nms          = true
}

output "nameservers" {
  value = module.orc8r.route53_nameservers
}
