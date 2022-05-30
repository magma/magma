# github.com/magma/magma Go module

This is the home for all new Golang code for all magma deployment targets.

Goals:

1. Single go.mod defines all dependencies in one place and keeps versions in
   sync.
2. Code can easily be shared with natual package paths that match file paths.
   No replace directives necessary in go.mod.
3. Make it easy to find shared libs/patterns for use across magma targets.
4. Other projects can easily import github.com/magma/magma/src/go code for use in
   derivative works.

Notes:

1. The go.mod is placed at $MAGMA_ROOT/src/go to avoid name collisions with existing Go directories.
2. We eventually want a canonical golang import path of magmacore.org/magma



## Bazel integration

### Generating Bazel build files

When updating or adding a new dependency, run

```sh
bazel run //:gazelle -- update-repos -from_file=src/go/go.mod -to_macro=bazel/go_repositories.bzl%go_repositories
```

To generate BUILD.bazel files, run

```sh
bazel run //:gazelle
```
### Running tests
To run all the tests, run
```sh
bazel test //src/go/...
```

## Dependencies

Only pull in dependencies when they are worth their weight in tech debt.
Strongly review libraries before adding as a dependency along the following
considerations:

1. Does the dependency have a license compatible with BSD-3?
2. How large is the dependency and its transitive dependencies? Weigh the total
   increase in go.sum against the provided functionality.
3. Does the dependency have good test coverage and is it well maintained?
4. What alternatives were considered?

Add the considerations in this README when adding the dependency to go.mod.

### Protobuf

Google/Go authors BSD-3 license.

A large/heavy dependency, but provides a solid, well-tested, and production-
hardened IDL. Good semantics for forwards/backwards compatibility and highly
performant. Front-runner for de facto standard IDL/serdes library in modern
software engineering projects.

Considered alternatives: Thrift, Swagger, JSON Schema.

### gRPC

Apache License 2.0

A large/heavy dependency that provides both RPC service IDL and base
client/server implementation. Tightly coupled with protobuf.

Considered alternatives: Thrift, Swagger, JSON-RPC

### Uber Zap

MIT license.

Medium weight, pulls in a few transitive dependencies, but has a large overlap
with testify. Provides
scoped/named loggers and parameterized logging while achieving better
performance than other libraries proven via benchmarks.

Very well tested: 98% code coverage; performance benchmarks included in CI

Considered alternatives: built-in log, glog, klog

### testify

MIT license.

Lightweight, for test purposes only. Good test coverage and widely adopted.

Downside is main usage eschews compile-time type checking in favor of runtime
reflection-based equality checking.

Considered alternatives: regular Golang code

### GoMock

Apache License 2.0

Lightweight test tool to codegen verifiable mocks from interfaces. Generated
EXPECT() functions are type-defined, which engages compile-time checking.
Strict call count and ordering expectation. These explicit requirements make
GoMock superior to mockery.

Considered alternatives: mockery, hand-coded mock implementations

## Dependency Injection

TBD -- manual DI or library e.g. Uber dig/fx

### External Dependencies / SDKs

In most cases, we will want to abstract dependencies via a clean interface
package. In this context, "clean" means the interface package is defined
entirely with Golang built-ins and other clean magma interface packages. i.e.
there are no 3rd party library dependencies, except for test helpers (e.g.
gomock).
