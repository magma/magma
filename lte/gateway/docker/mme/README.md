## MME Docker Container

### Ubuntu 20.04 with GCC 9.3 and support for mme build + test

```
cd lte/gateway/docker/mme/
docker build -t magma-mme-build -f Dockerfile.ubuntu20.04 ../../../../
```

### Ubuntu 18.04 with GCC 6.x and 7.x and support for mme build

```
cd lte/gateway/docker/mme/
docker build -t magma_mme -f Dockerfile.ubuntu18.04 ../../../../
```