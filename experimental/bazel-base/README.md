# Magma-Builder Docker Image

> :warning: **Bazel based builds are still experimental**

This Dockerfile is used to create a build space for all development with Bazel.

All docker-compose commands below are to be run from `$MAGMA_ROOT/experimental/bazel-base`.

## Build magma-builder Docker image

To build magma-builder base image, run the following.

```bash
docker-compose build magma-builder
```

## Run bazel commands

To run bazel commands, exec into a magma-builder container,

```bash
docker-compose run magma-builder bash
```

Once insider the container, bazel can be run like this,

```bash
# To build all targets
bazel build --config=docker ...
# To build a specific target (Ex: session_manager.proto)
bazel build --config=docker lte/protos:session_manager_cpp_proto
# To run all tests
bazel test --config=docker ...
```

## Format bazel files

To format all bazel related files, exec into a bazel container and run the following

```bash
bazel run //:buildifier
```
