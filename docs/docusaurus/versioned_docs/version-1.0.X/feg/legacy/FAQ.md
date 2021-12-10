---
id: version-1.0.1-faq
title: FAQ
hide_title: true
original_id: faq
---
# FAQ

1. Do I need to run the federated gateway as an individual developer?

   - It is highly unlikely you'll need this component. Only those who plan
   to integrate with a Mobile Network Operator will need the federated gateway.

2. I'm seeing 500's in `/var/log/syslog`. How do I fix this?

    - Ensure your cloud VM is up and services are running
    - Ensure that you've run `register_feg_vm` at `magma/feg/gateway` on your host machine

3. I'm seeing 200's, but streamed configs at `/var/opt/magma/configs` aren't being updated?

    - Ensure the directory at `/var/opt/magma/configs` exists
    - Ensure the gateway configs in NMS are created (see [link](https://github.com/facebookincubator/magma/blob/master/docs/Magma_Network_Management_System.pdf) for more instructions)
    - Ensure one of the following configs exist:
        - [Federated Gateway Network Configs](https://127.0.0.1:9443/apidocs#/Networks/post_networks__network_id__configs_federation)
        - [Federated Gateway Configs](https://127.0.0.1:9443/apidocs#/Gateways/post_networks__network_id__gateways__gateway_id__configs_federation)
