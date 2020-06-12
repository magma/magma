output "helm_vals" {
  description = "Helm values for the orc8r deployment"
  value       = helm_release.orc8r.values
  sensitive   = true
}
