# Containerized AGW

This folder contains container image definitions for AGW services.

The containers need to run on a host that has a patched Open vSwitch installation.
The script [agw_install_docker.sh](../deploy/agw_install_docker.sh) can configure an Ubuntu machine to act as a
host system for the containerized AGW.

There are Ansible playbooks that create an EC2 instance and prepare it to act as
a host for the containerized AGW. Currently the preparation is not complete and
still requires running `agw_install_docker.sh` at the end, see [Deploying the
containerized AGW on AWS](#deploying-the-containerized-agw-on-aws).

## Building the images

The images can be built with `cd $MAGMA_ROOT/lte/gateway/docker && docker compose --compatibility build`.
On an Arm architecture with the 5.4 kernel, the images can be built with `cd $MAGMA_ROOT/lte/gateway/docker && docker compose --compatibility build --build-arg CPU_ARCH=aarch64 --build-arg DEB_PORT=arm64`.

## Deploying the containerized AGW on AWS

- Build the docker images and push them to a registry of your choice
- Look at [the cloudstrapper readme](../../../experimental/cloudstrapper/README.md) for an AWS setup
    - If you are only interested in running the AGW, you only need steps 1 and 5.1
        - Step 1 will produce required AWS resources like security groups, S3 buckets, and a keypair, and download the keypair to a local directory that you specify
        - Step 5.1 will start and customize an Ubuntu machine and create an AMI snapshot of that machine. The machine stays running after the snapshot is taken.
        - You don't need to configure Github and Docker credentials in the `secrets.yaml`
    - After running step 5.1, you should be able to SSH to the created EC2 instance
- On the EC2 instance
    - Optional: Copy custom `rootCA.pem` to  `/var/opt/magma/certs` with permissions 400
    - Copy `agw_install_docker.sh` to the instance and run it to finish the preparation of the host
    - Adapt the docker registry in `/var/opt/magma/docker/.env`
        - Example for the registry setting: `DOCKER_REGISTRY=registry.hub.docker.com/arunuke/`
    - Make changes to config files and restart services by running `/var/opt/magma/docker/agw_upgrade.sh` or by running the `agw_install_docker.sh` script

## Running the containerized AGW locally on the magma VM

The magma VM defined in [../Vagrantfile](../Vagrantfile) can be used to run the
containerized AGW by either using a fabric command on the host machine:

```bash
cd ${MAGMA_ROOT}/lte/gateway
fab start-gateway-containerized
```

or by running the following steps inside the VM:

```bash
sudo rm -rf /etc/snowflake && sudo touch /etc/snowflake
cd $MAGMA_ROOT && bazel/scripts/link_scripts_for_bazel_integ_tests.sh
bazel build `bazel query "attr(tags, util_script, kind(.*_binary, //orc8r/... union //lte/... union //feg/... except //lte/gateway/c/core:mme_oai))"`
for component in redis nghttpx td-agent-bit; do cp "${MAGMA_ROOT}"/{orc8r,lte}/gateway/configs/templates/${component}.conf.template; done
sed -i 's/init_system: systemd/init_system: docker/' "${MAGMA_ROOT}"/lte/gateway/configs/magmad.yml
sudo systemctl start magma_dp@envoy

# Optional: If you want to connect to an orc8r, copy the `rootCA.pem` from the orc8r
# to `/var/opt/magma/certs/rootCA.pem`. For example, in a typical magma-dev VM:
cp ${MAGMA_ROOT}/.cache/test_certs/rootCA.pem /var/opt/magma/certs/

cd $MAGMA_ROOT/lte/gateway/docker
docker compose --compatibility build
docker compose --compatibility up
```

Note that with the containerized AGW we ultimately want to get rid of the dependency
on a VM. However, we are not there yet as the containerized AGW currently depends
on a patched Open vSwitch installation on the host machine. The magma VM happens
to have the right packages installed, and thus can currently be used as a quick
and dirty way to run the containers locally.

### Running the S1AP integration tests against the containerized AGW

To run the tests, first start the docker containers:

```
cd $MAGMA_ROOT/lte/gateway/docker
docker compose down # If containers are already running
docker compose --compatibility up
```

The test VM can then be set up and the tests executed by following
[these instructions](https://magma.github.io/magma/docs/next/lte/s1ap_tests#test-vm-setup).
