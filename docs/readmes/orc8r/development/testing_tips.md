---
id: testing_tips
title: Testing Tips
hide_title: true
---

# Testing Tips

This guide covers tips for quickly validating Orchestrator changes.

## About the build process

The build process follows a subproject-based architecture. The top-level build
file, `orc8r/cloud/Makefile`, defines a set of generic commands, which each
subproject implements, e.g. `orc8r/cloud/go/Makefile`.

The set of subprojects is determined by the `MAGMA_MODULES` environment
variable, which is defined during the container image build process. The
default subprojects include `orc8r`, `lte`, `feg`, etc.

## How to test

The normal way to run Orchestrator unit tests is `build.py -t`. This builds
and brings up a test and DB container, then runs the full set of unit tests.
Testing this way is effective, but can be heavyweight. Instead, you can also
run tests directly from your IDE. We'll describe how to run all tests directly
from IntelliJ in one click.

### Install prereqs

Our tests are not fully isolated from their environment. To set up your
environment for testing, run

```bash
cd $MAGMA_ROOT/orc8r/cloud/ && make tools  # install tools
cd $MAGMA_ROOT/orc8r/cloud/docker && ./run.py  # bring up postgres_test
```

### Configure IntelliJ run profiles

For each Orchestrator subproject, create an IntelliJ run configuration to
run the subproject's full set of tests. Name them `go test orc8r`,
`go test lte`, etc. Set the following environment variables for each config

- `TEST_DATABASE_HOST=localhost`
- `TEST_DATABASE_PORT_MARIA=3307`
- `TEST_DATABASE_PORT_POSTGRES=5433`

This should looks something like

![intellij_subproject_configs](assets/orc8r/intellij_subproject_configs.png)

### Configure IntelliJ Multirun profile

Install the
[Multirun plugin for IntelliJ](https://plugins.jetbrains.com/plugin/7248-multirun),
create the `go test all` Multirun config, and add all the per-subproject
configs to this new Multirun config.

This should looks something like

![intellij_multirun](assets/orc8r/intellij_multirun.png)

### Run tests

Now you should be able to run, without rebuilding any container images, the
full set of tests for

- a particular subproject, e.g. `go test lte`
- the entire codebase, e.g. `go test all`
