---
id: contribute_conventions
title: Contributing Conventions
hide_title: true
---

# Contributing Conventions

This document describes contributing conventions for the Magma project. The goal of this style guide is to

- Steer contributors away from experienced landmines
- Align on coding styles in support of a uniform Magma codebase, with the aim to improve developer productivity and codebase approachability

## General

Follow these conventions when making changes to the Magma codebase.

### Leave it better than you found it

- The project's principal convention is the [boy scout rule](https://www.oreilly.com/library/view/97-things-every/9780596809515/ch08.html): leave it better than you found it

### Add tests for your changes

- Tests should cover, at minimum, the mainline of the new feature. For reference, that usually ends up around 50-70% coverage.
- Unit tests should default to being placed in the same directory as the files they're testing, except for the following
    - Python: directly-adjacent `tests` directory
    - C/C++: directly-adjacent `tests` directory
    - JavaScript: directly-adjacent `__tests__` directory
- Integration tests should be placed as close to the code-under-test as possible
- If you're not sure how to test a change, reach out on the community Slack workspace for input

### Separate cleanup PRs from functional changes

- Keeps PRs small and understandable
- Exception: if the area of the codebase you're editing needs a cleanup PR, but you don't have bandwidth to add one, default to mimicking the surrounding code

### Scope component responsibilities

- Functions, and components in general, should be narrowly scoped to a single, specific role
- When writing a function over 100 lines long, consider extracting a helper functions out of the intermediate logical steps

### Prefer immutability and idempotency

- Prefer immutable state
- When mutability is necessary, consider the following
    - Prefer to set a component's state entirely in its constructor
    - Mutate a component's state as close to construction-time as possible
    - Perform mutations as high in the call chain as possible
- Prefer side-effect-free functions
    - When side-effects are necessary, move them as high in the call chain as possible

### Prefer composition over inheritance

- Avoid inheritance as a design pattern, in favor of composition and dependency injection
- If complex logic begins bleeding into test case setup, consider pulling that logic into a dependency interface
- Build a complex component as a composition of multiple simpler components with clear interfaces
- Avoid non-trivial static functions: pull interfaces out of the static functions and inject them into depending components

### Use simple constructors

- Split complex logic and side-effect-inducing functionality out of the constructor and into an initialization method
- If desired, can also use a static factory function to construct the component and call its initialization method

### Comment with why, not what

- Good code is self-documenting
- Instead of defaulting to inline comments, focus on
    - Concise and descriptive identifier names
    - Intelligent types and pervasive typing information
    - High-quality docstrings on functions, components, and top-level identifiers
- Avoid "topic sentence" comments
    - E.g. "this block of code does X ... this block of code does Y", when there's no value-add other than summarizing the next few lines
    - Instead, code paragraphs should be skimmable
    - Consider breaking dense code paragraphs out into private functions.
- Save comments for code blocks that require non-obvious context, e.g. answering why an idiosyncratic or non-obvious decision was made

### Follow style conventions

- Use [Go-style doc comments](https://golang.org/doc/effective_go#commentary), where the doc comment is prefixed by the name of the object being documented
- Use [Americanized spellings](https://en.wikipedia.org/wiki/Wikipedia:List_of_spelling_variants)
    - marshaling, marshaled
    - canceling, canceled, cancellation
- Use alphabetized [metasyntactic variables](https://en.wikipedia.org/wiki/Metasyntactic_variable)
    - Good: `apple, banana, cherry, date, egg`
    - Bad: `foo, bar, baz, quz, soap`
- Prefer underscores over hyphens
    - File, directory names
    - YAML, JSON
    - Exception: in certain parts of K8s, underscores are disallowed. In this case, hyphens are preferred, and translation between hyphens and underscores is acceptable.
- Omit trailing slash of directory paths, except where semantically meaningful
- Don't terminate new service names with `d`

## Documentation

### All documentation

The end-goal is uniform, approachable documentation, especially in the Docusaurus

- Write in [plain English](https://www.plainenglish.co.uk/how-to-write-in-plain-english.html)
    - Short sentences
    - Active verbs
    - Use "you" and "we"
    - Avoid nominalisations
    - Use lists
- [Use descriptive hyperlink text](https://www.nngroup.com/articles/writing-links/)
    - Don't use "here" as the text for a hyperlink
- Consistent capitalization and spelling
    - Magma-specific
        - Orc8r, Orchestrator
        - NMS
        - Magma
        - Mconfig
    - Magma-adjacent
        - K8s, Kubernetes
        - Helm
        - gRPC
        - eNodeB
    - Everyday words
        - metadata
        - use-case
    - Magma service names
        - state, subscriberdb, lte, etc.
- Use long form of CLI flags
    - `--deployment` rather than `-d`.

### Markdown

Consider the following Markdown conventions

- Don't wrap long lines
- Use `#` for headers, rather than underlining
- Place an empty line before and after a list block
- Don't preface lists with punctuation
    - Good: `Consider the following Markdown conventions`
    - Bad: `Consider the following Markdown conventions:`
    - Bad: `Consider the following Markdown conventions,`
    - Bad: `Consider the following Markdown conventions.`
- Default to sentence-casing listables
    - Good: `- Default to sentence-casing listables`
    - Bad: `- default to sentence-casing listables`
- Use hyphens for unordered lists, not asterisks
    - Good: `- Magma`
    - Bad: `* Magma`
- Use flat apostrophes and quotes, not curly
    - Good: `Magma's`, `"Magma"`
    - Bad: `Magma’s`, `“Magma”`
- Title-case H1 headers, sentence-case H2 and lower headers
    - Good: `# Code Conventions in Magma`, `## Code conventions in Magma`
    - Bad: `# Code conventions in Magma`, `## Code Conventions in Magma`
- Use space-padded double-hyphen to approximate an [em dash](https://www.thepunctuationguide.com/em-dash.html)
    - Good: `Magma -- an open source project -- has code conventions`
    - Bad: `Magma--an open source project--has code conventions`
    - Bad: `Magma - an open source project - has code conventions`

## Languages

### Go

Orc8r's cloud code has some basic [CI lint checks](https://github.com/magma/magma/blob/master/.golangci.yml). The [Go style guide](https://github.com/golang/go/wiki/CodeReviewComments) and anything concrete from [Effective Go](https://golang.org/doc/effective_go) are authoritative. Aside from those, consider the following conventions

**General**
- Familiarize yourself with these [3 common Go landmines](https://gist.github.com/lavalamp/4bd23295a9f32706a48f)
- Check in generated code
- Avoid init functions in Magma code
    - Exception: generated code and imported code
    - If a new init function absolutely must be added, it must be idempotent, contained to its package, and not cause global state mutation
- Avoid global state
    - This includes global registries
    - Preferred alternative: pass instance of object around directly
    - Acceptable alternative: singleton instances of a private object accessed by public functions
- Default to using separate `_test` package for tests
    - I.e. [same directory, different package](https://medium.com/@matryer/5-simple-tips-and-tricks-for-writing-unit-tests-in-golang-619653f90742)
    - Example: [`orc8r/cloud/go/services/state/indexer/indexer_test.go`](https://github.com/magma/magma/blob/51843e3245e0b785a738d991f62657c2cac328b6/orc8r/cloud/go/services/state/indexer/indexer_test.go#L14)
    - In almost all cases, the code-under-test should be re-writable into something that can be tested from an external test package
    - Only use same-package tests when absolutely necessary, and in that case put them in a separate test file
- When returning an error, all other returns should contain their [zero value](https://yourbasic.org/golang/default-zero-value/)

**Logging**
- Use [the `golang/glog` package](https://pkg.go.dev/github.com/golang/glog) for all logging
    - Default to `-v=0` for all services
- Deciding when to log errors
    - There are two types of errors: *platform* errors, where there's something wrong with Orc8r, and *client* errors, where a client made an invalid request. The former must be logged to communicate the error. The latter can be logged as a high-verbosity info log, as a helpful debugging tool -- that is, client errors should *not* be logged as Orc8r errors
    - Prefer returning an error rather than error logging -- error logging should occur as high in the call stack as possible
    - Add an explanatory comment when swallowing errors (i.e. when choosing to log an error and not bubble it up the call stack)
- Log level
    - `Fatal` *conservatively*, and only when the service has degraded to the point of inability to function. Fatal-ing on service startup is a useful pattern, but fatal-ing after service initialization should be avoided unless absolutely necessary.
    - `Error` when something is definitely wrong, e.g. a violated invariant.
    - `Warning` when something is probably wrong, but it's not possible to be sure it's an error. This is an infrequent use-case, prefer error.
    - `Info` for everything else, with appropriate verbosity.

**Style**
- Verbify function names
  - Exception: composable method names with well-understood functionality, e.g. `foo.Filter(...).Keys().Sorted()`
  - Exception: using `new*` or `New*` when instantiating new objects
- When import aliasing is required, prefer to alias with `snake_case` rather than `camelCase`
- Prefer readable code over rigid adherence to max line lengths. Capping around 140 characters feels about right.

### Python

The [PEP 8 style guide](https://www.python.org/dev/peps/pep-0008/) is authoritative.

**Type annotations**

- All new code should be fully type-annotated
  - For reference, please look at this [type hints cheat sheet for Python 3](https://mypy.readthedocs.io/en/stable/cheat_sheet_py3.html)

**Documentation**

- Document all public functions and *keep those docs up to date* when you make changes
- We use [Google style docstrings](https://google.github.io/styleguide/pyguide.html#383-functions-and-methods) in our codebase
  - For VSCode users, [Python Docstring Generator](https://marketplace.visualstudio.com/items?itemName=njpwerner.autodocstring) plugin is recommended
  - For IntelliJ users, you can configure a doc string format via `Preferences->Tools->Python Integrated Tools->Docstring format`

Example:
```
def foo(arg1: str) -> int:
    """Returns the length of arg1.

    Args:
        arg1 (str): string to calculate the length of

    Returns: the length of the provided parameter
    """
    return len(arg1)
```

**Logging**
- Use the [logging](https://docs.python.org/3/library/logging.html) module for all logging
- Refer to the Go logging section for deciding between log levels


**Linter**

- For mandatory lint checks, we have a unit test that runs [Pylint](https://pypi.org/project/pylint/) on all gateway services
  - On CI, the check gets run as part of the `lte-test` job
- Additionally, we have a [Reviewdog](https://github.com/reviewdog/reviewdog) linter using [wemake-python-styleguide](https://wemake-python-stylegui.de/en/latest/) enabled to aid the code review process
  - To run the linter locally, use the [precommit script](https://github.com/magma/magma/blob/master/lte/gateway/python/precommit.py)

**Formatters**

- We recommend [autopep8](https://pypi.org/project/autopep8/) as it conforms to [pep8](https://www.python.org/dev/peps/pep-0008/)
  - The above-mentioned [precommit script](https://github.com/magma/magma/blob/master/lte/gateway/python/precommit.py) also has an option to format your changes with
  [isort](https://pypi.org/project/isort/), [autopep8](https://pypi.org/project/autopep8/), and [add-trailing-comma](https://pypi.org/project/add-trailing-comma/)
- We do *not* recommend other formatters such as [black](https://black.readthedocs.io/en/stable/installation_and_usage.html), as it diverges from pep8 on basic things like line length, etc.


### C++

The [Google C++ Style Guide](https://google.github.io/styleguide/cppguide.html) is authoritative.

**Documentation**
- Always document your functions and classes in the header files over source files
- Use Doxygen style documentation for functions
  - For VSCode users, the [doxygen documentation generator plugin](https://marketplace.visualstudio.com/items?itemName=cschlosser.doxdocgen) is recommended

**Types**
- Be mindful when choosing input / output types when writing new functions
  - Opt for return values over output parameters
  - Non-optional input parameters should be values or const references
  - Use non-const pointers to represent optional outputs and optional input/output parameters

**Headers**
- Always include what you use
  - [cpplint](https://github.com/cpplint/cpplint) has include-what-you-use warnings

**Linter**
- We recommend Google's [cpplint](https://github.com/cpplint/cpplint) to lint your changes locally
  - For VSCode users, the [cpplint plugin](https://marketplace.visualstudio.com/items?itemName=mine.cpplint) is recommended
- We also have a [Reviewdog](https://github.com/reviewdog/reviewdog) annotatoer that runs [cpplint](https://github.com/cpplint/cpplint) to aid the code review process

**Logging**
- For non-OAI C++ services, use the `MLOG` macros defined in `orc8r/gateway/c/common/logging/magma_logging.h`
- For OAI, use the `OAILOG_*` macros defined in `lte/gateway/c/oai/common/log/h`
- Refer to the Go logging section for deciding between log levels

### Shell

- Shell script names should be suffixed with the proper file extension
    - Reference: [`sh` vs. `bash`](https://medium.com/100-days-of-linux/bash-vs-sh-whats-your-take-3e886e4c1cbc)
    - `.sh` for POSIX-compliant shell
    - `.bash` for bash
    - Default to `.bash` except with specific reason
- When a shell script passes around 100 lines, it's time to re-write it in Python or Go

### Javascript

**General**
- Import order should be: types, default components, and named components with each separated by a newline
- Use `TitleCase` for component file names
- Favor `map` and `forEach` functions in place of regular `for` loops
- Refrain from using literal strings/numbers without defining them
- Use strict equality(`===`) when comparing values

**Type annotations**
- All new code should be fully type anotatated
- We have a mandatory unit test that runs the [Flow](https://docs.magmacore.org/docs/next/nms/dev_testing#pre-commit-tests#flow) type checker on the NMS codebase
    - On CI, the check gets run as part of the `nms-flow-test` job

**Documentation**
- Document all the components and functions you create
- Keep the docs up to date when you make changes
- Use [JSDoc](https://jsdoc.app/index.html) tags to document your code
- Use `//` to comment line of code, if it is not clear in what it is doing

**Linter**
- Run [ESLint](https://docs.magmacore.org/docs/next/nms/dev_testing#pre-commit-tests#eslint) to lint your changes locally
- For mandatory lint checks, we have a unit test that runs `eslint` on JS code
    - On CI, the check gets run as part of the `eslint` job

## Tools

### gRPC

The [Google Protocol Buffer style guide](https://developers.google.com/protocol-buffers/docs/style) is authoritative. We also follow a subset of the [Uber Protocol Buffer style guide](https://github.com/uber/prototool/blob/dev/style/README.md). Consider the selection below

- [Streaming RPCs are strongly discouraged](https://github.com/uber/prototool/blob/dev/style/README.md#streaming-rpcs)
- When deprecating a field, [use the `deprecated` option instead of the `reserved` keyword](https://github.com/uber/prototool/blob/dev/style/README.md#reserved-keyword)
- RPC request and return definitions should be unique to the RPC
    - E.g. `rpc GetTrip(GetTripRequest) returns (GetTripResponse);`
    - This is especially relevant for servicer definitions at the Orc8r-gateway interface
- Uniform file structure [(example)](https://github.com/protocolbuffers/protobuf/blob/master/examples/addressbook.proto)
    - License
    - File overview
    - Syntax
    - Package
    - Imports (sorted)
    - File options
    - Everything else
        - Define *services first*, then their constituent objects
- Use `PascalCase` for message names and `snake_case` for field names
- Two-space indents

### Swagger

- Routes always return an object (forward compatibility)

### YAML

- Use the casing convention that is idiomatic for the code that will be reading the YAML file
    - Rationale: facilitates automatically unmarshaling the file to native object
    - Example: for Go config files, use `camelCase`
    - When a YAML file may be read by multiple languages, default to `snake_case`

### CLIs

- Consolidate related functionality into a single CLI
