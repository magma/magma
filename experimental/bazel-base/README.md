> :warning: **Bazel based builds are still experimental**

# Bazel-Build Docker Image
This Dockerfile is used to create a build space for all development with Bazel.

## Build bazel-build Docker image

To build bazel-build base image, run the following.

```bash
# MAGMA_ROOT should be set to repo root
export PATH_TO_DOCKERFILE=$MAGMA_ROOT/experimental/bazel-base/Dockerfile
docker build -t magma/bazel-build -f $PATH_TO_DOCKERFILE $MAGMA_ROOT
```

## Run bazel commands

To run bazel commands, exec into a bazel-build container,

```bash
docker run -v $MAGMA_ROOT:/magma -v $MAGMA_ROOT/lte/gateway/configs:/etc/magma -i -t magma/bazel-build:latest /bin/bash
```

Once insider the container, bazel can be run like this,

```bash
# To build all targets
bazel build ...
# To build a specific target (Ex: session_manager.proto)
bazel build lte/protos:session_manager_cpp_proto
# To run all tests
bazel test ...
```

## Format bazel files

To format all bazel related files, exec into a bazel container and run the following

```bash
bazel run //:buildifier
```
