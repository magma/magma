################################################################################
# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
################################################################################

output "orc8r_domain_name" {
  description = "Base domain name for Orchestrator application components."
  value       = var.orc8r_domain_name
}

output "eks_cluster_id" {
  description = "Cluster ID for the EKS cluster"
  value       = module.eks.cluster_id
}

output "kubeconfig" {
  description = "kubectl config file to access the EKS cluster"
  value       = module.eks.kubeconfig
}

output "eks_config_map_aws_auth" {
  description = "A k8s configmap to allow authentication to the EKS cluster."
  value       = module.eks.config_map_aws_auth
  sensitive   = true
}

output "efs_file_system_id" {
  description = "ID of the EFS file system created for k8s persistent volumes."
  value       = aws_efs_file_system.eks_pv.id
}

output "efs_provisioner_role_arn" {
  description = "ARN of the IAM role for the EFS provisioner."
  value       = aws_iam_role.efs_provisioner.arn
}

output "es_endpoint" {
  description = "Endpoint of the ES cluster if deployed."
  value       = join("", aws_elasticsearch_domain.es.*.endpoint)
}

output "secretsmanager_secret_name" {
  description = "Name of the secretsmanager secret for deployment secrets"
  value       = aws_secretsmanager_secret.orc8r_secrets.name
}

output "orc8r_db_host" {
  description = "Hostname of the Orchestrator RDS instance"
  value       = aws_db_instance.default.address
}

output "orc8r_db_name" {
  description = "Database name for Orchestrator RDS instance"
  value       = aws_db_instance.default.name
}

output "orc8r_db_port" {
  description = "Database connection port for Orchestrator RDS instance"
  value       = aws_db_instance.default.port
}

output "orc8r_db_user" {
  description = "Database username for Orchestrator RDS instance"
  value       = aws_db_instance.default.username
}

output "orc8r_db_pass" {
  description = "Orchestrator DB password"
  value       = aws_db_instance.default.password
  sensitive   = true
}

output "nms_db_host" {
  description = "Hostname of the NMS RDS instance"
  value       = aws_db_instance.nms.address
}

output "nms_db_name" {
  description = "Database name for NMS RDS instance"
  value       = aws_db_instance.nms.name
}

output "nms_db_user" {
  description = "Database username for NMS RDS instance"
  value       = aws_db_instance.nms.username
}

output "nms_db_pass" {
  description = "NMS DB password"
  value       = aws_db_instance.nms.password
  sensitive   = true
}

output "route53_zone_id" {
  description = "Route53 zone ID for Orchestrator domain name"
  value       = aws_route53_zone.orc8r.id
}

output "route53_nameservers" {
  description = "Route53 zone nameservers for external DNS configuration."
  value       = aws_route53_zone.orc8r.name_servers
}

output "external_dns_role_arn" {
  description = "IAM role ARN for external-dns"
  value       = aws_iam_role.external_dns.arn
}
