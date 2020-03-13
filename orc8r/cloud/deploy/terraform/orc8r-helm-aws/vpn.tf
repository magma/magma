resource "helm_release" "ovpn" {
  count = var.deploy_openvpn ? 1 : 0

  chart      = "openvpn"
  name       = "openvpn"
  namespace  = var.orc8r_kubernetes_namespace
  repository = data.helm_repository.stable.id

  # TCP ovpn because ELB does not support UDP
  values = [<<EOT
  openvpn:
    OVPN_K8S_POD_NETWORK: null
    OVPN_K8S_POD_SUBNET: null
    OVPN_PROTO: tcp
    redirectGateway: false
  service:
    annotations:
      external-dns.alpha.kubernetes.io/hostname: vpn.${var.orc8r_domain_name}
  persistence:
    existingClaim: ${kubernetes_persistent_volume_claim.storage["openvpn"].metadata.0.name}
  EOT
  ]

  # Cert creation can take some time
  timeout = 900
}