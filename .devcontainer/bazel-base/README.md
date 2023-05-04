# Magma-Builder Docker Image

> :warning: **Bazel based builds are still experimental**

This Dockerfile is used to create a build space for all development with Bazel.

## Prerequisites

Throughout this guide we assume the `MAGMA_ROOT` environment variable is set to the local directory where you cloned the Magma repository

```bash
export MAGMA_ROOT=PATH_TO_YOUR_MAGMA_CLONE
```

## Build magma-builder Docker image

All `docker compose` commands below are to be run from `$MAGMA_ROOT/.devcontainer/bazel-base`.

To build magma-builder base image, run the following.

```bash
docker compose --compatibility build magma-builder
```

## Run bazel commands

To run bazel commands, exec into a magma-builder container,

```bash
docker compose --compatibility run magma-builder bash
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

To format all bazel related files run the following

```bash
cd $MAGMA_ROOT
./bazel/scripts/run_buildifier.sh format
```
