# Magma-Builder Docker Image

> :warning: **Bazel based builds are still experimental**

This Dockerfile is used to create a build space for all development with Bazel.

## Prerequisites

Throughout this guide we assume the `MAGMA_ROOT` environment variable is set to the local directory where you cloned the Magma repository

```bash
export MAGMA_ROOT=PATH_TO_YOUR_MAGMA_CLONE
```

## Build magma-builder Docker image

All docker-compose commands below are to be run from `$MAGMA_ROOT/experimental/bazel-base`.

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

## Generate go_repository via Gazelle

Gazelle is a tool that generates Bazel configurations from an existing Go project

Any time there is a dependency upgrade or a new Go dependency is added to the project, run the following

```bash
bazel run //:gazelle -- update-repos -from_file=src/go/go.mod -to_macro=go_repositories.bzl%go_repositories
```

This will output all `go_repository` configurations into `$MAGMA_ROOT/go_repositories.bzl`.
