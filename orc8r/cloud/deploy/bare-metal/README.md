Deploying Magma on bare metal

The following files need to be edited before deployment can start:
* deploy_ansible_vars.yaml
* orc8r_settings (shell env file)

The values that need to be customized are the IP settings for your network and
the passwords which need to be generated.

Once the config is set, just run the following commands:

./deploy.sh
./deploy_charts.sh

If the deployment succeeds, you will see information on how to log into Magma
web UI.

Note: external-dns configuration is not automatic here because it is intended
for on-premise. You are expected to update DNS or /etc/hosts on your own.
