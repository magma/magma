---
id: version-1.7.0-agw_unable_to_checkin
title: Access Gateway Unable to Check-in to Orchestrator
hide_title: true
original_id: agw_unable_to_checkin
---

# Access Gateway Unable to Check-in to Orchestrator

## About

**Description:** After deploying AGW and Orchestrator, it is time to make AGW accessible from Orchestrator. After following github Magma AGW configuration [guide](../../lte/deploy_config_agw.md), it was observed that AGW is not able to check-in to Orchestrator.

**Environment:** AGW and Orc8r deployed.

**Affected components:** AGW, Orchestrator

## Resolution

1. Diagnose AGW and Orchestrator setup with script checkin_cli.py. If the test is not successful, the script would provide potential root cause for a problem. A successful script will look like below

    ```text
    AGW$ sudo checkin_cli.py

    1. -- Testing TCP connection to controller-staging.magma.etagecom.io:443 --
    2. -- Testing Certificate --
    3. -- Testing SSL --
    4. -- Creating direct cloud checkin --
    5. -- Creating proxy cloud checkin --
    Success!
    ```

    If the output is not successful, the script will recommend some steps to resolve the problem. After following the steps the problem has not been resolved, follow below steps.

2. Make sure that the hostnames and ports specified in control_proxy.yml file in AGW are properly set.
Sample control_proxy.yml file

    ```yaml
    cloud_address: controller.yourdomain.com
    cloud_port: 443
    bootstrap_address: bootstrapper-controller.yourdomain.com
    bootstrap_port: 443

    rootca_cert: /var/opt/magma/tmp/certs/rootCA.pem
    ```

3. Verify the certificate rootCA.pem is in the correct location defined in rootca_cert (specified in control_proxy.yml)

4. Make sure the certificates have not expired.
    Note: To obtain certificate information you can use `openSSL x509 -in certificate -noout -text`
    - In AGW: rootCA.pem
    - In Orc8r: rootCA.pem, controller.cert

5. Verify the domain is consistent across AGW and Orc8r and the CN matches with the domain
    - CN in rootCA.pem AGW
    - CN in Orc8r for root and controller certificates.
    - The domain in `main.tf`

6. Verify connectivity between AGW and Orc8r.  Choose the port and domain obtained in `control_proxy.yml`. You can use telnet, example below:
 `telnet bootstrapper-controller.yourdomain.com 443`

7. Verify the DNS resolution of the bootstrap and controller domain.
    - In AGW: You can ping or telnet to your bootstrap and controller domain from AGW to verify which AWS address is being resolved.
    - In Orc8r: Verify which external-IP your cluster is assigned. You can use the command: `kubectl get services`

    The address resolved in AGW should be the same defined in Orc8r. If not,  verify your DNS resolution.

8. Verify that there are no errors in AGW magmad service.

    `AGW$ sudo tail -f /var/log/syslog | grep -i "magmad"`

9. From Orchestrator, get all pods and verify attempts from AGW are reaching Orc8r in nginx and look for any bootstrapping erors in the bootstrapper.

    - First, you can use below command to get all pods from orc8r

        ```bash
        kubectl -n orc8r get pods
        ```

        For example, boostrapper and nginx Orc8r pods should look something like below:

        ```text
        orc8r-bootstrapper-775b5b8f6d-89spq              1/1     Running   0          37d
        orc8r-bootstrapper-775b5b8f6d-gfmrp              1/1     Running   0          37d
        orc8r-nginx-5f599dd8d5-rz4gm                     1/1     Running   0          37d
        orc8r-nginx-5f599dd8d5-sxpzf                     1/1     Running   0          37d
        ```

    - Next, using the pod name get the logs from the pod with below command. Check if there is any problematic log for the related pod

        ```bash
        kubectl -n orc8r logs -f <nginx podname>
        kubectl -n orc8r logs -f <bootstrapper podname>
        ```

        For example: `kubectl -n orc8r get logs orc8r-bootstrapper-775b5b8f6d-89spq`

10. Try restarting magmad services.

    ```text
    AGW$ sudo service magma@magmad restart
    ```

11. If issue still persists, please  file github issues or ask in our support channels <https://magmacore.org/join-the-open-source-community/>
