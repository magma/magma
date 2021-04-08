## Build MME Ubuntu18 docker image

Depending on your installed `docker` version, you may have to do the following:

```bash
# Go to the root of the MAGMA repository in your workspace
cd $MAGMA
mv .dockerignore .mockerignore
```

The current `.dockerignore` file is targeted for the orchestrator images building
and is ignoring a lot of files we are needing.

Then you can build:

```bash
cd $MAGMA 
docker build --target magma-mme --tag magma-mme:latest --file lte/gateway/docker/mme/Dockerfile.ubuntu18.04  .
```

## Build MME RHEL8 podman image

We are using `podman3.0` or later versions.

First building a RHEL8 image, using a lot of developers packages, requires to pass the certificates to enable some YUM repositories.

**This implies that you are running your podman commands on an already certified RHEL system!**

```bash
mkdir -p tmp/ca tmp/entitlement
cp /etc/pki/entitlement/*pem tmp/entitlement
sudo cp /etc/rhsm/ca/redhat-uep.pem tmp/ca
```

Finally you can build:

`podman3` has a nice feature to change the location of the `.dockerignore` file.

```bash
cd $MAGMA
sudo podman build --target magma-mme --tag magma-mme:latest --ignorefile lte/gateway/docker/mme/.dockerignore --file lte/gateway/docker/mme/Dockerfile.rhel8 .
```

