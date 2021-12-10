---
id: version-1.5.0-dev_testing
title: Testing Tips
hide_title: true
original_id: dev_testing
---

# Testing Tips

This guide covers tips for quickly validating Orchestrator changes.

## Build process overview

The build process follows a subproject-based architecture. The top-level build
file, `orc8r/cloud/Makefile`, defines a set of generic commands, which each
subproject implements, e.g. `orc8r/cloud/go/Makefile`.

The set of subprojects is determined by the `MAGMA_MODULES` environment
variable, which is defined during the container image build process. The
default subprojects include `orc8r`, `lte`, `feg`, etc.

By default, all subprojects are included in the build. If desired, you can
use e.g. `build.py --deployment orc8r` to build a deployment-specific set of
subprojects. This can reduce build time by up to 50%.

## Run tests on the host

The normal way to run Orchestrator unit tests is `build.py --tests`. This
builds and brings up a test and DB container, then runs the full set of unit
tests. Testing this way is effective, but can be heavyweight.

Instead, you can also run tests directly from your IDE. We provide a default
set of IntelliJ run configurations to support running all tests in one click,
directly from IntelliJ.

### Default run configurations

The default run configurations are checked in under `.run/`. You'll also need
to install the
[Multirun plugin](https://plugins.jetbrains.com/plugin/7248-multirun).

`Go Test` configurations should look something like

![intellij_subproject_configs](assets/orc8r/intellij_subproject_configs.png)

`Multirun` configurations should include `go test all`, and look something like

![intellij_multirun](assets/orc8r/intellij_multirun.png)

### Install prereqs

Our tests are not fully isolated from their environment. To set up your
environment for testing, run

```bash
cd ${MAGMA_ROOT}/orc8r/cloud/ && make tools  # install tools
cd ${MAGMA_ROOT}/orc8r/cloud/docker && ./run.py  # bring up postgres_test
```

### Run tests

Now you should be able to run the full set of tests, without rebuilding any
container images, for

- a particular subproject e.g. `go test lte`
- the entire codebase `go test all`

## Custom run configurations

You can also manually create your own run configurations. Depending on the
test, you may need to include the following environment variables

- `TEST_DATABASE_HOST=localhost`
- `TEST_DATABASE_PORT_MARIA=3307`
- `TEST_DATABASE_PORT_POSTGRES=5433`
