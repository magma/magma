---
id: agw_with_bazel
title: AGW with Bazel
hide_title: true
---

# AGW with Bazel

This page provides a Bazel-specific overview on building and testing the Magma Access Gateway (AGW).
This documentation is aimed mainly at Magma developers. For general Bazel support at Magma there is also the Slack channel [#bazel](https://magmacore.slack.com/archives/C033S1CEAUU).

Questions related to Bazel may not be specific to Magma and the [official Bazel reference](https://bazel.build/reference) and the `bazel help` command are both very useful in these cases.

> Note: A guide on how the AGW Make commands translate to Bazel can be found in the [GitHub Wiki](https://github.com/magma/magma/wiki/Bazel-vs.-Make-commands-dictionary).

## What is Bazel?

[Bazel](https://bazel.build/about) is an open-source build system that is intended for multi-language mono-repos like Magma. Bazel is focused on fast, scalable, parallel and reproducible builds. A core strength of Bazel is an extensive [caching framework](#caching) based on a strict and fine grained build graph.

### Structure and files

The basic structure of the Bazel files in Magma is as follows:

```bash
.
├── WORKSPACE.bazel
├── .bazelrc
├── BUILD.bazel
├── bazel
│   ├── bazelrcs
│   │   ├── vm.bazelrc
│   │   └── ...
│   ├── external
│   │   ├── BUILD.bazel
│   │   ├── sentry_native.BUILD
│   │   ├── requirements.in
│   │   ├── requirements.txt
│   │   └── ...
│   ├── scripts
│   │   ├── bazel_diff.sh
│   │   ├── run_buildifier.sh
│   │   └── ...
│   ├── BUILD.bazel
│   ├── cpp_repositories.bzl
│   ├── go_repositories.bzl
│   ├── python_repositories.bzl
│   ├── python_swagger.bzl
│   └── ...
├── lte
│   ├── gateway
│   │   └── python
│   │       ├── magma
│   │       │   ├── mobilityd
│   │       │   │   ├── tests
│   │       │   │   │   └── BUILD.bazel
│   │       │   │   └── BUILD.bazel
│   │       │   └── policydb
│   │       │       ├── servicers
│   │       │       │   └── BUILD.bazel
│   │       │       ├── tests
│   │       │       │   └── BUILD.bazel
│   │       │       └── BUILD.bazel
│   │       └── BUILD.bazel
│   ├── protos
│   │   ├── oai
│   │   │   └── BUILD.bazel
│   │   └── BUILD.bazel
│   └── swagger
│       └── BUILD.bazel
.
```

- [`WORKSPACE.bazel`](https://bazel.build/concepts/build-ref#workspace): Defines a Bazel project. Configuration of [Bazel rules](https://bazel.build/extending/rules) to be used in the project.
- [`.bazelrc`](https://bazel.build/run/bazelrc): General configuration of flags for commands.
- [`BUILD.bazel`](https://bazel.build/concepts/build-files): A file, e.g. `lte/gateway/python/magma/mobilityd/BUILD.bazel`, in which sources are defined in binaries, libraries and tests. A `BUILD.bazel` file defines a [Bazel package](https://bazel.build/concepts/build-ref#packages).
- `bazel/`: Centralized folder for custom Bazel definitions and configuration of external sources.
- `bazel/scripts/`: Centralized folder for Bash wrapper scripts for Bazel, e.g. for running the [Starlark formatter](#starlark-formatter), [Bazel-diff](#bazel-diff) or the [LTE integration tests](#lte-integration-tests).
- `bazel/bazelrcs/name.bazelrc`: Configuration of flags for commands in different environments.
- `bazel/name.bzl`
    - Custom Bazel code - these are custom [Bazel rules](https://bazel.build/extending/rules), configuration of the tool chain and external dependencies.
    - In `cpp`, `go` and `python_repositories.bzl`, external repositories are specified. These can be local repositories (folder in the environment) or git repositories.
- `bazel/external/name.BUILD`: `.BUILD` files for external repositories that are not bazelified.

#### Starlark formatter

To ensure the standardized formatting of all `BUILD.bazel` and `.bzl` files the Starlark formatter [buildifier](https://github.com/bazelbuild/buildtools/tree/master/buildifier) is used and run in [CI](#workflows).

You can also run the Starlark format check locally, from your host or from any [Magma environment](#environments).

To run the Starlark format check locally, run

```bash
cd $MAGMA_ROOT
./bazel/scripts/run_buildifier.sh check
```

To run the Starlark formatter and to fix most errors that occur during the check, run

```bash
cd $MAGMA_ROOT
./bazel/scripts/run_buildifier.sh format
```

### Environments

The following Magma environments currently support Bazel:

- Docker
    - **devcontainer**
        - The devcontainer is defined in `.devcontainer/Dockerfile` and meant for local development. The docker image is published at [`ghcr.io/magma/magma/devcontainer:latest`](https://github.com/magma/magma/pkgs/container/magma%2Fdevcontainer).
    - **Bazel-base**
        - The Bazel-base docker container is a minimal container where only required dependencies are installed. While this image is mainly intended to be used in CI, it can also be useful for local testing. The Bazel-base image is defined in `.devcontainer/bazel-base/Dockerfile` and published at [`ghcr.io/magma/magma/bazel-base:latest`](https://github.com/magma/magma/pkgs/container/magma%2Fbazel-base).
- Vagrant VMs
    - **magma-dev**
        - The magma-dev virtual machine is intended for local development and for extended testing in CI. It can be used to run the Magma AGW services, as well as e.g. the [Python sudo tests](#python-sudo-tests).
    - **magma-test**
        - The magma-test virtual machine is used in CI to run the [LTE integration tests](#lte-integration-tests).

### Caching

Bazel has a very extensive and granular system of caching intermediate build outputs, artifacts and even test results. Builds that have many cache hits are much faster than builds without cache hits.

The number of cache hits is reported at the end of each build, e.g. as `INFO: 4449 processes: 703 disk cache hit, 2377 internal, 1367 processwrapper-sandbox, 2 worker.`. Without code changes it should be possible to cache essentially all processes.

Whenever cacheable results are produced, a local cache entry is generated in the [disk cache](https://bazel.build/remote/caching#disk-cache). Some of the cache locations can be found by running `bazel info` and looking e.g. at the `repository_cache` entry.

> Info: A useful guide on [Bazel cache analysis](https://github.com/magma/magma/wiki/Bazel-cache-analysis) can be found on the GitHub Wiki pages.

In CI, caches have to be downloaded from elsewhere, this can be done either via [Docker images that contain pre-built caches](#docker-bazel-cache), [GitHub caching](#github-cache) or by using a [remote cache](#bazel-remote-cache). [Remote caching](#bazel-remote-cache) is the recommended way of caching for Bazel in CI. Remote caches can be used in any environment (Docker, Vagrant, etc.) and the remote cache is always up-to-date, as cache entries are constantly being uploaded. Remote caching can also be simpler than other caching approaches, in particular when using the [Google cloud backend](#google-cloud-storage).

#### Docker Bazel cache

The Docker Bazel caches are Docker containers that contain pre-built Bazel caches for all relevant [targets](#targets) and [configurations](#configurations). The containers are created in the [workflow](#workflows) `docker-builder-devcontainer.yml` on a regular basis. The Docker Bazel cache containers are based on the [Bazel-base](#environments) Docker container with the addition of the Bazel caches. The images are published at:

- `ghcr.io/magma/magma/bazel-cache-plain:latest`
- `ghcr.io/magma/magma/bazel-cache-asan:latest`
- `ghcr.io/magma/magma/bazel-cache-prod:latest`

The Docker Bazel caches were introduced in [#14562](https://github.com/magma/magma/issues/14562) - this issue also contains a detailed slide with information on the Docker Bazel cache setup.

#### GitHub cache

GitHub Actions provides a generic [built-in caching solution](https://github.com/actions/cache), which can be used to store the Bazel disk cache. Possible downsides of this caching approach, as with the [Docker caches](#docker-bazel-cache), are that the cache is not managed by Bazel and has to be kept up-to-date and below a reasonable cache size. The GitHub cache has a 10GB limit and is also used for other purposes at Magma, such as caching Vagrant base images.

#### Bazel remote cache

[Remote caching](https://bazel.build/remote/caching) is an integrated caching mechanism that is part of Bazel. Cache entries are up-/downloaded to/from a backend server. There are many different backends that can be used, such as [Google Cloud Storage](#google-cloud-storage) or the [buchgr/bazel-remote](#buchgrbazel-remote). Various backends are also documented in the [official documentation](https://bazel.build/remote/caching). To optimize the build times, the backend should be as close as possible to the CI runners, ideally in the same data center.

##### Google Cloud Storage

The simplest and best supported backend for Bazel remote caching is Google Cloud Storage. With Google Cloud Storage a simple storage bucket is sufficient to operate the remote cache - for details see the [official documentation](https://bazel.build/remote/caching#cloud-storage).

##### buchgr/bazel-remote

The [buchgr/bazel-remote](https://github.com/buchgr/bazel-remote/) is a self-hosted Bazel remote caching backend that runs as a Docker container.

> Info: As of the 4th of January 2023 the bazel-remote cache for Magma is **deprecated**. The tear-down is documented in the issue [#14796](https://github.com/magma/magma/issues/14796).

Information on the bazel-remote caching infrastructure for Magma can be found in the [magma/ci-infra/bazel/remote_caching](https://github.com/magma/ci-infra/blob/master/bazel/remote_caching/Readme.md) repository (CI codeowner access only). This information in particular details how the bazel-remote cache can be deployed on AWS using Terraform. General information on Bazel remote caching can be found in the [official Bazel documentation](https://bazel.build/remote/caching).

Furthermore, there are GitHub Wiki pages with information on [how to use a deployment of the bazel-remote cache](https://github.com/magma/magma/wiki/Bazel-remote-caching) and on [how to recover from a bazel-remote cache failure](https://github.com/magma/magma/wiki/Bazel-remote-caching-disaster-recovery-plan).

### Commands

The Bazel commands are structured around verbs like [`build`](#building), [`test`](#testing) or [`query`](#querying). The complete list of Bazel commands can be found in the [official command reference](https://bazel.build/run/build#available-commands). The following sections go into more detail on the most common commands.

## Building

To build something with Bazel, use the command [`bazel build`](https://bazel.build/run/build) and then specify the [targets](#targets) that you want to build, e.g. `bazel build //lte/gateway/python/magma/mobilityd:mobilityd`.

Building happens in a sandbox, i.e. various host properties are not present at build time. This includes environment variables and stdin. The build artifacts can be found in `$MAGMA_ROOT/bazel-bin`, e.g. `$MAGMA_ROOT/bazel-bin/lte/gateway/python/magma/mobilityd/mobilityd`.

> Info: Information on debugging Bazel builds can be found in the section [debugging](#debugging).

### Targets

The Bazel project structure is organized around targets in the build files, written in the [Starlark language](https://github.com/bazelbuild/starlark). Targets usually have at least the attributes `name`, `srcs`, `deps`, but there may be many more attributes. The list of possible attributes depends on the Bazel rule.

A `cc_library` target, for example, could look like this:

```starlark
cc_library(
	name = "my_library",
	srcs = ["my_code.c"],
	hdrs = ["my_header.h"],
	deps = ["//other/path:other_library"],
)
```

This example is a C-library with one source file `my_code.c`, header `my_header.h` and a build dependency to another library `//other/path:other_library`. The attribute name `my_library` is the name of the Bazel target. It can be used for building the target directly, by running `bazel build //lib/path:my_library`, or for referencing the target as a dependency in another target.

The official Bazel documentation provides a complete list of all attributes of the [C/C++ rules](https://bazel.build/reference/be/c-cpp), the [Python rules](https://bazel.build/reference/be/python), the [Go rules](https://github.com/bazelbuild/rules_go) and many others.

Some Bazel rules that are used in the Magma AGW are:

- `cc_library`, `py_library` and `go_library` are used to group C, C++, Python and Go code respectively.
- `cc_binary`, `py_binary`, `go_binary` are used for the definitions of the C, C++, Python and Go AGW services and scripts, e.g. the `//lte/gateway/c/core:agw_of` target for MME, the `//lte/gateway/python/magma/pipelined` target for pipelined or the `//feg/gateway/services/envoy_controller` target for the envoy controller.
- `cc_test` and `go_test` are used for C, C++ and Go unit tests.
- `pytest_test` is a wrapper around `py_test`. This rule type is used for Python unit tests, the [LTE integration tests](#lte-integration-tests) and the [Python sudo tests](#python-sudo-tests).
- `proto_library`, `cpp_proto_library`, `python_proto_library`, `cpp_grpc_library`, `python_grpc_library` are used to group proto code, as well as for providing proto and grpc code for C++ and Python respectively.

#### Target visibility

Bazel targets have a [visibility attribute](https://bazel.build/concepts/visibility). Per default a target's visibility is set to private and the target is only accessible for targets within the same package. Possible values include (among others) a public visibility `//visibility:public` (accessible from all packages) or package specific visibility `//foo/bar:__pkg__`, to grant targets in package `//foo/bar` access. The visibility can be set explicitly for a target or globally for the whole package.

The following coding guidelines should be applied when adding new targets. One important motivation for the guidelines below is based on encapsulation, i.e. keeping a clean architecture by restricting access to targets. For example, service libraries should only be accessed by the service itself. Proto and utility libraries can be accessed by any service (with certain restrictions if appropriate - e.g. the `all_python_proto` target should only be accessed by the `state` service).

- Protos `BUILD.bazel` files: The global package visibility is set to public as compiled protos libraries need to be accessed from a lot of packages. The visibility of native proto libraries is set to private or is restricted to package visibility for other protos packages.
- Service `BUILD.bazel` files: The global package visibility is set to private (i.e. default). The binaries should explicitly be set to private for readability reasons, even if this means redundancy. Libraries should explicitly be set to private or package visibility (e.g. for testing) when needed. The public visibility should be used only if necessary.
- Utility Packages: In the `orc8r/configuration` and `orc8r/common` `BUILD.bazel` files, the global package visibility is set to public.

### Target command line syntax

A complete description of the [syntax for specifying targets](https://bazel.build/run/build) can be found in the official Bazel documentation. Here we provide a list of some useful examples:

- To build all targets in the Bazel workspace, excluding the ones [tagged as manual](https://bazel.build/reference/be/common-definitions#common.tags), use `bazel build ...` or `bazel build //...`.

> Warning: Building all targets for the first time, i.e. without any caches, can overwhelm some systems, due to the amount of parallelization. The number of parallel jobs that Bazel runs in parallel can be limited with the `--jobs` flag, e.g. on a system with 8 physical CPU cores, one could use `--jobs 7` for the initial build.

- To build all targets inside of a folder, e.g. `lte/gateway`, use `bazel build //lte/gateway/...`.
- To build one target, use e.g. `bazel build //lte/gateway/python/magma/mobilityd:mobilityd`. Here the `//` indicates the start of the path to the target and the `:` indicates the start of the target name.
    - The same target can be built with `bazel build //lte/gateway/python/magma/mobilityd` as this expands to the target with the same name as the folder.
- To specify multiple targets or multiple folders, simply add them to the command line, e.g. `bazel build //lte/gateway/python/magma/mobilityd:mobilityd //orc8r/gateway/...`. This example will build the mobilityd service target, as well as all targets inside the `orc8r/gateway` folder.

When specifying targets on the command line, the leading `//` are optional, but they make it clearer that a target is referenced. The prefix `//` is required in build files when a dependency is referenced.

A complete list of [command line arguments](https://bazel.build/reference/command-line-reference) can be found in the official Bazel documentation.

#### Alias targets

For developer convenience there are [alias targets](https://bazel.build/reference/be/general#alias) for all AGW services. An alias is simply a more convenient label for a target.

The service aliases for the AGW are defined in the top level `BUILD.bazel` file, in order to keep the labels as short as possible. An alias target, e.g. for the mobilityd service, looks like this:

```Starlark
alias(
    name = "mobilityd",
    actual = "//lte/gateway/python/magma/mobilityd:mobilityd",
)
```

Instead of using the label `//lte/gateway/python/magma/mobilityd:mobilityd`, the shorter, more convenient label `//:mobilityd` can be used to refer to the target on the command line.

### Configurations

Targets can be built with different configurations. Configurations are specified on the command line by adding the command line flag `--config=<config_name>`. Definitions of configurations can be found in the `.bazelrc` file. A configuration is specified by a Bazel command followed by a colon, the name of the configuration and the content, e.g. `build:example --announce_rc --color=yes`.

 The most important configurations for Magma are:

- Plain, i.e. no special configuration, this is the development version of the targets.
- `production`, which is the production version of the targets used for releases.
- `asan`, a configuration with the address sanitizer enabled, used for finding memory corruption bugs.
- `lsan`, a configuration with the leak sanitizer enabled, used for finding run-time memory leaks.

#### Inheritance

When editing configurations and options in the `.bazelrc` file, one has to be aware of the [precedence and inheritance](https://bazel.build/run/bazelrc#option-defaults) of the bazel commands, as well as the effects of the [options](https://bazel.build/reference/command-line-reference).

For example, the options specified on the command line take precedence over options specified in the `.bazelrc` file. The commands `test`, `run`, `clean` and `info` inherit all options from `build`. The command `coverage` inherits from `test`. All commands inherit options from the specifier `common`.

Inversely, this also means that when options, such as `--verbose_test_summary`, are set for the `test` command, they are not used by the `build` command. In general, options should be set at the highest possible command level, i.e. `build` would be preferable over `test`, and `common` should be used whenever possible.

When options are added to the `test` command that are marked as `affects_outputs` in the [Bazel command line reference](https://bazel.build/reference/command-line-reference#test-options), this can lead to slower builds and less cache hits if the option is not set at the `build` or `common` level. This is because the cache entries for the two commands, for the same target, will be different. A detailed description of this problem, with possible solutions, can also be found in the issue [#13073](https://github.com/magma/magma/issues/13073).

### Dependencies

The `deps` attribute is used to add a dependency to a Bazel target. The `deps` attribute exists for all library, binary and test rules. Dependencies are specified as a list of strings. There are several distinct types of dependency:

- [Local dependencies](#local-dependencies)
- [External dependencies](#external-dependencies)
- [Pip dependencies](#pip-dependencies-python).

A simple example that illustrates these is as follows:

```starlark
py_library(
  name = "my_lib",
  srcs = ["my_lib.py"],
  deps = [
  	":my_package_dependency", # This is a local dependency.
  	"//orc8r/magma/common:my_util", # This is a local dependency.
  	"@aioh2_repo//:aioh2", # This is an external dependency.
  	requirement("grpcio"), # This is a pip dependency.
  ],
)
```

It is [best practice](https://bazel.build/configure/best-practices) to use fine-grained dependencies and to keep targets small, as this allows for better caching and parallelism. Closely connected source files should be contained in the same target for maintainability.

Furthermore, Bazel does not allow circular dependencies among targets. If a circular dependency is present, the build will fail with an error that a cycle in the dependency graph has been detected. The MME service currently does not allow for a finer granularity because circular dependencies are present and therefore larger targets need to be used.

#### Local dependencies

Dependencies on targets in the **same package**, i.e. the same `BUILD.bazel` file, are specified by the target name, prefixed by `:`. For example, as `deps = [":my_package_dependency"]`.

Dependencies on targets in **other packages within Magma** are additionally prefixed by `//` and the path of the respective package relative to the repository root. For example, as `deps = ["//orc8r/magma/common:my_util"]`.

#### External dependencies

External dependencies can be added to `$MAGMA_ROOT/bazel/<language>_repositories.bzl`. The target is either defined in the external source (if it is a bazelified dependency) or in the respective build file in `$MAGMA_ROOT/bazel/external` (if it is a dependency that is not bazelified).

In the `deps` attribute of the target that needs the external dependency, the name of the external dependency repository is prefixed by `@` and followed by the Bazel target. For example, as `deps = ["@aioh2_repo//:aioh2"]`.

#### Pip dependencies (Python)

In most cases, external Python dependencies can be added as pip dependencies. For this, the dependency needs to be added to the `$MAGMA_ROOT/bazel/external/requirements.in` file. If a new dependency is added, the `requirements.txt` file needs to be rebuilt by

```bash
$MAGMA_ROOT/bazel/external $ pip-compile --allow-unsafe --generate-hashes --output-file=requirements.txt requirements.in
```

which requires `pip-tools` to be installed.

> Info: See `bazel/external/requirements_README.md` for more information regarding the `requirements.in` file.

The pip dependency, e.g. `grpcio`, can now be added to a target by adding a requirement statement `deps = [requirement("grpcio")]`.

### Cleaning

The Bazel output, i.e. the artifacts, can be cleaned by running [`bazel clean`](https://bazel.build/docs/user-manual#cleaning-build-outputs), which will remove all outputs from all previous builds for all configurations. The entire working tree created by Bazel can be removed by running `bazel clean --expunge`.

Usually, cleaning should not be needed as Bazel can build incrementally. However, cleaning can be used to reclaim disk space or to recover a consistent state when [debugging](#debugging).

> Info: The `bazel clean` command does **not** remove the disk caches. To remove the disk caches, the cache folders, e.g. `$MAGMA_ROOT/.bazel-cache` and `$MAGMA_ROOT/.bazel-cache-repo`, need to be removed manually. The exact locations of the disk caches may depend on the environment.

## Testing

There are four categories of tests that concern the AGW, each of these categories is covered in the following sections:

- [Unit tests](#unit-tests)
- [LTE integration tests](#lte-integration-tests)
- [Python sudo tests](#python-sudo-tests)
- [Load tests](#load-tests)

> Warning: Note that not all Bazel test targets in the Magma AGW can be executed directly with Bazel. These special targets are tagged in their attributes and the following sections explain how they can be run.

### Unit tests

See the [notes on unit testing](../lte/dev_unit_testing.md) for more detailed information and examples on unit testing.

The unit tests can be run on any of the Docker [environments](#environments) or on the magma-dev VM.

To run unit tests with Bazel, use the command `bazel test` and then specify the [Bazel targets](#targets) using the same [command line syntax](#target-command-line-syntax) described in the build section.

For example, to run all unit tests, that are not tagged as manual, run `bazel test //...` and to run a specific unit test specify the test target directly, e.g. `bazel test //lte/gateway/c/core/oai/test/pipelined_client:pipelined_client_test`.

> Warning: Testing all targets for the first time, i.e. without any caches, can overwhelm some systems, due to the amount of parallelization. The number of parallel jobs that Bazel runs in parallel can be limited with the `--jobs` flag, e.g. on a system with 8 physical CPU cores, one could use `--jobs 7` for the initial build.

Some useful command line options for unit testing include:

- `--cache_test_results=no` disables the caching of test results and forces a re-run even if there are no changes.
- `--runs_per_test=<integer>` executes all requested test targets exactly `<integer>`-many times, even if runs are successful. This flag can be used e.g. to detect flaky tests.
- `--flaky_test_attempts=<integer>` executes all test targets that have the `flaky = True` attribute set, up to `<integer>`-many times. If there are failures, the test is marked as flaky in the results. Test targets with the attribute `flaky = False` (default) are unaffected.

A complete list of [command line arguments](https://bazel.build/reference/command-line-reference) can be found in the official Bazel documentation.

### LTE integration tests

Information on how to run the LTE integration tests with Bazel is documented on the [s1ap tests](../lte/s1ap_tests.md) page.

> **Warning**: If you try to execute a LTE integration test with `bazel test` you might get a permission error. Do not execute Bazel commands as sudo/root user, because it can break your local Bazel setup!

The LTE integration tests cannot be executed directly with Bazel because the tests need extensive access to system resources that are not available during the Bazel runtime. In CI, the LTE integration tests are executed directly with pytest. Locally, there is a wrapper script that can be used - for instructions on how to use this script, please refer to the [s1ap tests](../lte/s1ap_tests.md) page.

### Python sudo tests

The Python sudo tests are a group of tests for the pipelined and mobilityd services that require root permissions during the test runtime.

> **Warning**: Do not execute Bazel commands as sudo/root user, because it can break your local Bazel setup!

To run all Python sudo tests, start the magma-dev VM and run the script:

```bash
cd $MAGMA_ROOT
bazel/scripts/run_sudo_tests.sh
```

To run individual tests, specify the test target, for example

```bash
cd $MAGMA_ROOT
bazel/scripts/run_sudo_tests.sh //lte/gateway/python/magma/pipelined/tests:test_classifier_mme_flow_dl
```

To display the full list of supported options, run

```bash
cd $MAGMA_ROOT
bazel/scripts/run_sudo_tests.sh --help
```

### Load tests

> **Info**: For detailed information on the load tests, please refer to the README file `lte/gateway/python/load_tests/README.md`.

To run every load test that is defined inside the `LOAD_TEST_LIST` of the script `bazel/scripts/run_load_tests.sh`, run the script on the magma-dev VM from anywhere inside the Magma folder.

For example, for mobilityd's `AllocateIPRequest`, there is a `load_test_mobilityd.py` which contains an `allocate` command (run as `load_test_mobilityd.py allocate`).

The load tests are defined in Bazel as `py_binary` targets, because they run as scripts and not as tests.

### Generating unit test coverage

To generate coverage data for unit tests with Bazel, use the command `bazel coverage` and then specify the unit test targets. Generating unit test coverage can be done on any of the Docker [environments](#environments) or on the magma-dev VM.

To generate test coverage data for all C and C++ unit tests, run

```bash
bazel coverage //orc8r/gateway/c/...:* //lte/gateway/c/...:*
```

Similarly, to generate test coverage data for all Python unit tests, run

```bash
bazel coverage //orc8r/gateway/python/...:* //lte/gateway/python/...:*
```

The coverage data can then be found in the folder `$MAGMA_ROOT/bazel-out/_coverage/_coverage_report.dat`.

To generate a HTML report, run

```bash
cd $MAGMA_ROOT
genhtml bazel-out/_coverage/_coverage_report.dat --output-directory <path/to/output/directory>
```

The HTML report can be opened via `<path/to/output/directory>/index.html`.

For general information on the `bazel coverage` command, see the [official Bazel documentation](https://bazel.build/configure/coverage).

## Running

The `py_binary` targets tagged as either `service` or `util_script` should not be run with the `bazel run` command, because their use requires special wrappers.

### Services

The AGW services are running on the magma-dev VM as systemd services, therefore the Bazel-built executables need to be wrapped in a systemd service file. The file templates can be found in the folder `$MAGMA_ROOT/lte/gateway/deploy/roles/magma/files/systemd_bazel/`.

To learn more about the services, go to the page [AGW Services/Sub-Components](../lte/readme_agw.md).

### Utility scripts

The AGW utility scripts can be executed independently and they are sometimes called by the services. These scripts are expected to work from any folder on the magma-dev VM, therefore the Bazel-built artifacts are linked to `/usr/local/bin/[filename].py`. To link the scripts, run

```bash
cd $MAGMA_ROOT
bazel/scripts/link_scripts_for_bazel_integ_tests.sh
```

### How to run services and scripts

The process for building and running the services and scripts is already automated and documented in the [Quick Start Guide](../basics/quick_start_guide.md). To build all services, on the magma-dev VM, you can use the alias `magma-build-agw`. For the detailed commands, you can look at the script `$MAGMA_ROOT/bazel/scripts/build_and_run_bazelified_agw.sh`.

To run an individual utility script, after it was built and linked, execute the name of the file on the magma-dev VM in any directory.

## Querying

Bazel is equipped with a powerful query language, which can be used to analyze information about the targets and packages. The output of a Bazel query can be used as the input for other Bazel commands, this can be achieved, e.g. by using a sub-shell.

Detailed information on the `bazel query` command, including many examples, can be found on the pages

- [Query quickstart guide](https://bazel.build/query/quickstart)
- [Query guide](https://bazel.build/query/guide)
- [Bazel query reference](https://bazel.build/query/language)

of the official Bazel documentation.

The remainder of this section contains some Magma AGW specific examples.

To find all test targets that use the Python test rule `py_test`, run

```bash
bazel query 'kind(py_test, //...)'
```

To exclude the Python tests that are tagged as `manual`, e.g. integration tests, from the output, run

```bash
bazel query 'kind(py_test, //...) except attr(tags, manual, //...)'
```

To list all scripts and services, run

```bash
bazel query "attr(tags, 'service|util_script', kind(.*_binary, //...)"
```

To list all C or C++ services, run

```bash
bazel query "attr(tags, 'service', kind(cc_binary, //...)"
```

To create a graph with all dependencies of, e.g. the envoy controller service, excluding external packages, run

```bash
bazel query 'filter("(^[^@])", filter("^(?!.*external)", deps(//feg/gateway/services/envoy_controller )))' --output graph | dot -Tsvg > envoy_controller_deps.svg
```

Creating the image requires the [graphviz](https://graphviz.org/download/) package to be installed.

## Bazel CI

In CI, there are several [workflows](#workflows) that use Bazel and they are briefly introduced in the following section.

### Workflows

- `agw-coverage.yml`: This workflow runs on pull requests and merges to master. It creates and uploads the unit test [coverage](#generating-unit-test-coverage) data for C, C++ and Python tests.
- `bazel.yml`: This is the main workflow for jobs that use Bazel. The workflow runs on pull requests and merges to master. In order to only run the jobs that are really needed on pull requests, the [Bazel-diff](#bazel-diff) tool is used. The workflow contains the
    - [starlark format check](#starlark-formatter).
    - C, C++, Go and Python [unit tests](#unit-tests) with the [configurations](#configurations):
        - native/plain
        - `asan`
        - `production`
    - build of the Magma and SCTP debian packages.
    - runtime check that all Python imports in the AGW services are bazelified.
    - check that all C, C++ and Python files in the repository are bazelified or excluded in the `bazel/scripts/check_py_bazel.sh` script.
- `docker-builder-devcontainer.yml`: The workflow runs on pull requests and merges to master. If necessary, this workflow builds and uploads the [Docker environments](#environments). This includes in particular the Docker images that contain pre-built Bazel caches that can be used in CI or locally. These cache images are rebuilt regularly.
- `lte-integ-test-bazel-magma-deb.yml`: This workflow runs on merges to master. It runs the [LTE integration tests](#lte-integration-tests) on the AGW that have been installed via the debian package.
- `sudo-python-tests.yml`: This workflow runs on merges to master. It runs the [Python sudo tests](#python-sudo-tests).

### Bazel-diff

[Tinder/bazel-diff](https://github.com/Tinder/bazel-diff) is an open-source tool that can determine which Bazel targets are impacted by the changes between two git commits. Bazel-diff is used in the `bazel.yml` [workflow](#workflows) on pull requests to determine if the changes impact

- the debian package/release build
- the Python service import statements
- C, C++, Python and Go AGW unit test targets

The granularity for the unit tests is individual Bazel test targets, which typically include a single file.

For ease of use, the Bazel-diff tool has been wrapped in a Bash script `bazel/scripts/bazel_diff.sh`. To determine the impacted targets from after the commit `<GIT_SHA_PRE>` up to and including the commit `<GIT_SHA_POST>`, run

```bash
cd $MAGMA_ROOT
bazel/scripts/bazel_diff.sh <GIT_SHA_PRE> <GIT_SHA_POST>
```

The script then outputs the list of impacted targets to stdout.

#### Filtering the Bazel-diff output

Using [`bazel query`](#querying), the output of Bazel-diff can be filtered, e.g. for test targets or the release target. In CI, this is done with the script `bazel/scripts/filter_test_targets.sh`, which takes a Bazel target rule name as a mandatory argument and a tag as an optional second argument.

To filter a list of targets, contained in a file `<IMPACTED_TARGETS_FILE>`, for test targets of any language, run

```bash
bazel/scripts/filter_test_targets.sh ".*_test" < "<IMPACTED_TARGETS_FILE>"
```

The script will then print the test targets to stdout.

To filter a list of targets for Python services, run

```bash
bazel/scripts/filter_test_targets.sh "py_binary" "service" < "<IMPACTED_TARGETS_FILE>"
```

To filter for specific individual targets, like the release build, use `grep`, e.g.

```bash
grep --quiet '//lte/gateway/release:release_build' "<IMPACTED_TARGETS_FILE>"
```

## Release build

To build the Magma, SCTPD and DHCP helper CLI debian packages, with [production configuration](#configurations), locally, run

```bash
bazel run //lte/gateway/release:release_build --config=production
```

The target `//lte/gateway/release:release_build` is a wrapper target for creating the Magma, SCTPD and DHCP helper CLI debian artifacts with the proper versions. Creating the debian packages with [`rules_pkg`](https://github.com/bazelbuild/rules_pkg) must be reproducible and a version like `1.9.0-1667381719-ebd3bb56` (`<version>_<timestamp>_<hash>`) would produce artifacts that are not reproducible. This target builds the artifacts and changes the version afterwards.

The debian packages that have been built can be found at `/tmp/packages/*.deb`.

The debian packages are built in CI in the `bazel.yml` [workflow](#workflows).

## Debugging

In this section, we list some steps that can be taken to debug failing Bazel builds.

- Make sure that the code is in the right state and that the correct [environment](#environments) is used.
- When debugging failing builds, the [`bazel info`](https://bazel.build/docs/user-manual#info) command can be helpful to document the current state of the system and the configurations.
- The [`bazel clean`](#cleaning) or `bazel clean --expunge` commands can help to recover a consistent state. Usually, this is not needed between builds as Bazel can build incrementally.
- Permission errors can occur when Bazel commands were accidentally run with root permissions. In this case one should update the permissions of the affected files from the host system.
- The flags `--sandbox_debug --verbose_failures` can be used to preserve the [sandbox](https://bazel.build/docs/sandboxing) that was used for the build and show the commands that were used to set up the sandbox. To find the Bazel output base, run `bazel info output_base`. The sandbox will be located below the output base and could be at a path like `<output_base>/execroot/__main__/bazel-out/k8-fastbuild/`. The option `--spawn_strategy=standalone` can be used to disable sand-boxing.
- To empty the Bazel [caches](#caching), remove the cache folders. They may be located at `$MAGMA_ROOT/.bazel-cache` and `$MAGMA_ROOT/.bazel-cache-repo`. Removing the cache folders will significantly slow down the subsequent builds until the caches are filled again. Removing the caches should usually not be necessary.
- The [Bazel profile](https://bazel.build/rules/performance#performance-profiling) can be used to debug issues related to performance or hardware limitations. The default location for the profile is `<output_base>/command.profile.gz`. To find the Bazel output base, run `bazel info output_base`. The option `--profile` can be used to specify a different path and name for the profile file. Bazel profile files can be opened in the Chromium or Chrome browsers by navigating to `chrome://tracing` and loading the profile. For more details, see the [official Bazel documentation](https://bazel.build/rules/performance#performance-profiling).
