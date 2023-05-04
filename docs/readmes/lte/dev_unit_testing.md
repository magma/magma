---
id: dev_unit_testing
title: Test AGW
hide_title: true
---

# Test Access Gateway

This guide covers tips for quickly validating AGW changes.

## Run all unit tests

Unit testing for AGW can be done either inside the magma-dev VM, or inside the devcontainer or bazel-base Docker containers.

To SSH into the magma-dev VM, run

```bash
[HOST] cd $MAGMA_ROOT/lte/gateway
[HOST] vagrant up magma
[HOST] vagrant ssh magma
```

To start the devcontainer, run

```bash
[HOST] cd $MAGMA_ROOT
[HOST] docker run -v ${MAGMA_ROOT}:/workspaces/magma/ -v ${MAGMA_ROOT}/lte/gateway/configs:/etc/magma/ -it ghcr.io/magma/magma/devcontainer:latest /bin/bash
```

To start the bazel-base container, run

```bash
[HOST] cd $MAGMA_ROOT
[HOST] docker run -v ${MAGMA_ROOT}:/workspaces/magma/ -v ${MAGMA_ROOT}/lte/gateway/configs:/etc/magma/ -it ghcr.io/magma/magma/bazel-base:latest /bin/bash
```

To run all existing unit tests, run the following command from inside the repository

```bash
bazel test //...
```

### Test Python AGW services

To run only the Python unit tests, run the following command from inside the repository

```bash
bazel test //lte/gateway/python/... //orc8r/gateway/python/...
```

To run unit tests of an arbitrary service or library, run commands of the following form.
E.g. for the magmad unit tests, run

```bash
bazel test //orc8r/gateway/python/magma/magmad/...
```

and for the unit tests of the common library, run

```bash
bazel test //orc8r/gateway/python/magma/common/...
```

### Test C/C++ AGW services

We have several C/C++ services that live in `lte/gateway/c/`. We will list some of the useful commands here, but please refer to the [Bazel user guide](https://docs.bazel.build/versions/main/guide.html) for a complete overview.

From inside the repository, run

```bash
bazel test //lte/gateway/c/session_manager/...:* # to test all targets under lte/gateway/c/session_manager
bazel test //orc8r/gateway/c/...:* //lte/gateway/c/...:* # to test all C/C++ targets
```

### Test Go AGW services

We have several Go implementations of AGW services that live in `orc8r/gateway/go`.
To test any changes, run the following from inside the magma-dev VM

```bash
[VM] cd magma/orc8r/gateway/go
[VM] go test ./...
```

## Format AGW

### Format Python

Docker is required for running the steps below.
To use the `--diff` flag, the script will have to be on your host machine where the Magma repository lives.
Refer to the script at `lte/gateway/python/precommit.py` for all available commands, but the main ones are as follows.

```bash
cd $MAGMA/lte/gateway/python

# run the flake8 linter by specifying paths
./precommit.py --lint -p PATH1 PATH2
# run the flake8 linter on all modified files in HEAD vs master
# this command can only be run on your host
./precommit.py --lint --diff

# run all available formatters by specifying paths
./precommit.py --format -p PATH1 PATH2
# run all available formatters on all modified files in HEAD vs master
# this command can only be run on your host
./precommit.py --format --diff
```

### Format C/C++

To run formatting for each C/C++ service, run the following from inside the magma-dev VM

```bash
[VM] cd magma/dev_tools
[VM] ./clang_format.sh
```

#### Apply IWYU

> This tool currently only works inside the devcontainer environment and does not support fixups for `lte/gateway/c/core`.

[include-what-you-use](https://include-what-you-use.org/) is a tool developed by Google to analyze C++ files to help ensure source files include all headers used.
We have added a utility script that uses Bazel to generate a compilation database, then uses two scripts provided by IWYU, `iwyu_tool.py` and `fix_includes.py`, to apply changes.

To use the script, run

```bash
# Recommended: Run IWYU for a specific directory, run `apply-iwyu.sh <PATH>`
[DevContainer] $MAGMA_ROOT/dev_tools/apply-iwyu.sh orc8r/gateway/c
[DevContainer] $MAGMA_ROOT/dev_tools/apply-iwyu.sh lte/gateway/c/session_manager
# Run IWYU for all C/C++ files
# Note: The script currently does not work for lte/gateway/c/core, so you may need to revert changes for that directory
[DevContainer] $MAGMA_ROOT/dev_tools/apply-iwyu.sh
```

### Format Bazel BUILD files

To format all Bazel related files, run

```bash
cd $MAGMA_ROOT
./bazel/scripts/run_buildifier.sh format
```
