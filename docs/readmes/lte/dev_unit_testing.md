---
id: dev_unit_testing
title: Test AGW
hide_title: true
---

# Test Access Gateway

This guide covers tips for quickly validating AGW changes.

## Run all unit tests on the dev VM

In general, all unit testing for AGW is done on the magma dev VM.
To SSH into the VM, run

```bash
[HOST] cd $MAGMA_ROOT/lte/gateway
[HOST] vagrant up magma
[HOST] vagrant ssh magma
```

To run all existing unit tests, run

```bash
[VM] cd magma # or any subdirectory inside magma
[VM] bazel test //lte/gateway/...
```

Note: Running all unit tests can take close to 15 minutes.

In short, to run tests with Bazel you just need to provide the directory that the test targets reside in:

```bash
[VM] bazel test //... # to test all targets in magma
[VM] bazel test //lte/gateway/c/session_manager/... # to test all targets under lte/gateway/c/session_manager 
[VM] bazel test //orc8r/gateway/c/... //lte/gateway/c/... # to test all C/C++ targets
[VM] bazel test //lte/gateway/python/magma/enodebd/... # to test all enodebd test targets
[VM] bazel query "kind(py_test, //lte/gateway/python/magma/... union //orc8r/gateway/python/magma/...) except attr(tags, manual, kind(py_test, //lte/gateway/python/...))" # to list all python service unit tests
```

### Test Go AGW services

We have several Go implementations of AGW services that live in `orc8r/gateway/go`.
To test any changes, run

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
# build the base image
./precommit.py --build

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

To run formatting for each C/C++ service, run

```bash
[VM] cd magma/lte/gateway
[VM] make format_all
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
[VM] cd magma # or any subdirectory inside magma
[VM] bazel run //:buildifier
```
