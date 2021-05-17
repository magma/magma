---
id: deploy_dns
title: DNS Resolution
hide_title: true
original_id: deploy_dns
---
# DNS Resolution

In the following steps, replace `yourdomain.com` with the TLD or subdomain that
you've chosen to host Orchestrator on. It's important that you follow the
naming conventions for subdomains in order for your Access Gateways to
successfully communicate with Orchestrator.

First, grab the public-facing hostnames for the ELB instance fronting the
internet-facing Orchestrator components
(`orc8r-bootstrap-legacy`, `orc8r-clientcert-legacy`, `nginx-proxy`):

```bash
$ kubectl -n magma get svc -o='custom-columns=NAME:.metadata.name,HOSTNAME:.status.loadBalancer.ingress[0].hostname'

NAME                      HOSTNAME
magmalte                  &lt;none&gt;
nginx-proxy               ELB-ADDRESS1.elb.amazonaws.com
orc8r-bootstrap-legacy    ELB-ADDRESS2.elb.amazonaws.com
orc8r-clientcert-legacy   ELB-ADDRESS3.elb.amazonaws.com
orc8r-controller          &lt;none&gt;
orc8r-graphite            &lt;none&gt;
orc8r-metrics             &lt;none&gt;
orc8r-prometheus-cache    &lt;none&gt;
orc8r-proxy               ELB-ADDRESS4.elb.amazonaws.com
```

Set up the following CNAME records for your chosen domain:

| Subdomain | CNAME |
|-----------|-------|
| nms.yourdomain.com | nginx-proxy hostname |
| controller.yourdomain.com | orc8r-clientcert-legacy hostname |
| bootstrapper-controller.yourdomain.com | orc8r-bootstrap-legacy hostname |
| api.yourdomain.com | orc8r-clientcert-legacy hostname |

Wait for DNS records to propagate, then if you go to
`https://nms.yourdomain.com`, you should be able to log in as the admin user
that your created earlier and create your first network.
