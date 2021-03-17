# certs

Install cert-manager:
```bash
helm repo add jetstack https://charts.jetstack.io
helm repo update

helm install cert-manager jetstack/cert-manager \
  --create-namespace \
  --namespace cert-manager \
  --set installCRDs=true
```
---

Install certs:
```bash
helm install certs . \
  --set dnsDomain=magma.shubhamtatvamasi.com
```
