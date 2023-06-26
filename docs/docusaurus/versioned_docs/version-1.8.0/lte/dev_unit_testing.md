---
id: version-1.8.0-dev_unit_testing
title: Test AGW
hide_title: true
original_id: dev_unit_testing
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
[VM] cd magma/lte/gateway
[VM] make test
```

Note: Running all unit tests can take close to 15 minutes.

### Test Python AGW services

To run only the Python unit tests, run

```bash
[VM] cd magma/lte/gateway
[VM] make test_python
```

The list of services to test are configured in the following files.

- `orc8r/gateway/python/defs.mk`
- `lte/gateway/python/defs.mk`

To run unit tests for a single Python service, select a name from the list of services and run

```bash
[VM] cd magma/lte/gateway
[VM] make test_python_service MAGMA_SERVICE=<service_name>
```

In the case that unit tests for a single Python service are started multiple times, it is preferable to avoid the upstream installation process of the virtual environment by adding `DONT_BUILD_ENV=1` to the command, run

```bash
[VM] cd magma/lte/gateway
[VM] make test_python_service MAGMA_SERVICE=<service_name> DONT_BUILD_ENV=1
```

To run unit tests of an arbitrary directory, run

```bash
[VM] cd magma/lte/gateway
[VM] make test_python_service UT_PATH=<path_of_the_test_folder>
```

### Test C/C++ AGW services

We have several C/C++ services that live in `lte/gateway/c/`.
To run tests for those services, run

```bash
[VM] cd magma/lte/gateway
[VM] make test_<service_directory_name> # Ex: make test_session_manager
```

A subset of the AGW C/C++ directories are in the process of being migrated to use Bazel as the default build system. We will list out some of the useful commands here, but please refer to the [Bazel user guide](https://docs.bazel.build/versions/main/guide.html) for a complete overview.

```bash
[VM] cd magma # or any subdirectory inside magma
[VM] bazel test //... # to test all targets
[VM] bazel test //lte/gateway/c/session_manager/...:* # to test all targets under lte/gateway/c/session_manager 
[VM] bazel test //orc8r/gateway/c/...:* //lte/gateway/c/...:* # to test all C/C++ targets
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
