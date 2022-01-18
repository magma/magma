# Access Gateway

Update **control_proxy.yml** file as per your domain:
```bash
sudo mkdir -p /var/opt/magma/configs/
sudo vim /var/opt/magma/configs/control_proxy.yml
```

```yaml
cloud_address: controller.magma.local
cloud_port: 443
bootstrap_address: bootstrapper-controller.magma.local
bootstrap_port: 443
fluentd_address: fluentd.magma.local
fluentd_port: 443

rootca_cert: /var/opt/magma/tmp/certs/rootCA.pem
```
> Note: if you are using Private IP then you will have to update `/etc/hosts` file as per your domain.

Update **rootCA.pem** file from as per your Orchestrator:
```bash
sudo mkdir -p /var/opt/magma/tmp/certs/
sudo vim /var/opt/magma/tmp/certs/rootCA.pem
```

Get Hardware ID and Challenge key from the Access Gateway:
```bash
show_gateway_info.py
```

Stop and restart magma services:
```bash
sudo systemctl stop magma@*
sudo systemctl restart magma@magmad
```

Follow the `magmad` service logs for getting updates on Bootstrap and Checkin:
```bash
sudo journalctl -fu magma@magmad
```

### Cleanup

Remove old gateway keys and network config files from the Access Gateway:
```bash
sudo rm /var/opt/magma/certs/gateway.crt
sudo rm /var/opt/magma/certs/gateway.key
sudo rm /var/opt/magma/configs/gateway.mconfig
```
