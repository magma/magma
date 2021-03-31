---
id: deploy_build
title: Build FeG
hide_title: true
---

# Build Federation Gateway Components

If you cloned Magma using git, make sure you are checked out on the release you
intend to build.

In case you need to change the version you can:

```bash
# to list all releases
git tag -l

# to switch to a different release (for example v1.3.3)
git checkout v1.3.3

# to switch to master (developement version)
git checkout master
```

Once you are on a proper version of Magma, make sure your Docker daemon is running.
Then go run those commands to build FeG.

```bash
cd magma/feg/gateway/docker
docker-compose build --parallel
# if build fails try with sudo and without parallelization
sudo docker-compose build
```

Note that you are building FeG from your local repository. There is no need to
change content `.env`

If this is your first time building the FeG, this may take a while.

When this job finishes, you will have built FeG on your local machine. You can
check the images using docker. You should `gateway_python` and `gateway_go`
among others images that were used during the build process.
```bash
docker images
```

In case you want to host FeG on your image registry do the following to upload these
images:

```bash
../../../orc8r/tools/docker/publish.sh -r <REGISTRY> -i gateway_python
../../../orc8r/tools/docker/publish.sh -r <REGISTRY> -i gateway_go
```

In case you built Magma CWF (Carrier Wi-FI), you also need to upload `gateway_radius`.
```bash
../../../orc8r/tools/docker/publish.sh -r <REGISTRY> -i gateway_radius
```

