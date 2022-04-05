---
id: version-1.7.0-debug_logs title: View Logs hide_title: true
title: debug_logs title: View Logs hide_title: true
original_id: debug_logs title: View Logs hide_title: true
---

# Debugging

Debugging Domain Proxy usually comes down to checking log content of individual pods. This document describes how to
access these logs.

## List pods

To list Domain Proxy pods running in production environment type:

```console
kubectl -n orc8r get pod -l app.kubernetes.io/name=domain-proxy
```

You will get an output similar to this:

```console
NAME                                                    READY   STATUS      RESTARTS   AGE
domain-proxy-active-mode-controller-7b984c6579-zmwrm    1/1     Running     0          13d
domain-proxy-configuration-controller-6d99c978f-b8h6b   1/1     Running     0          13d
domain-proxy-radio-controller-5c868696d9-s7vgg          1/1     Running     0          13d
```

## Check individual pods' logs

```console
# Last 1000 lines of logs on a specific pod
kubectl logs --tail=1000 domain-proxy-active-mode-controller-7b984c6579-zmwrm

# Last hour worth of logs
kubectl logs --since=1h domain-proxy-active-mode-controller-7b984c6579-zmwrm

# Live preview mode
kubectl logs -f domain-proxy-active-mode-controller-7b984c6579-zmwrm
```
