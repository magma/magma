# VS Code Container Dev Environment

For a comprehensive developer environnment supporting AGW c / c++ workflows (extensible with ~effort to other Magma uses) - see the instructions below.

## Setup

- Install Docker
- Configure Docker
  - In some operating systems (e.g. MAC) docker runs in a VM and you must allocate memory / CPU / etc
  - Some targets in Magma can be built with e.g. 7 Gigs of VM RAM
    - e.g. `make test_oai` and `make build_oai` within `lte/gateway/`
  - Others require 12 Gigs of VM RAM
    - e.g. `make test_connection_manager` within `lte/gateway/`
- Install VS Code
  - Install Microsoft's VS Code [Remote - Containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) Extension
- `F1` For command search
  - Type `Remote-Containers: Clone Repository in Containers Volume...`
  - In dialog box that pops up, paste `https://github.com/electronjoe/magma/tree/experimental-vscode-dockerized`
    - Once this PR is committed, you could instead point at any Magma Github {branch or PR}
  - When prompted, select `Create a Unique Volume`
    - As our goal is to only work on this one Git target
  - Sit back and relax for awhile (10-15 min)
    - VS Code will build Dockerfile
    - Checkout Repo
    - Install plugins within Docker instance
  - Test your setup by opening a terminal and testing OAI
    - `Terminal->New Terminal`
    - `cd lte/gateway`
    - `make test_oai`

## Using Clang-Tidy in VS Code

While we could commit / build a compilation database (compile_commands.json), at this time we are not doing so. Instead we will explicitly call out the Include paths for VS Code / Clang-Tidy as follows (open to ammending this if a clean mechanism of generating compile commands is found):

You should not have to ask for Clang-Tidy to run on the open file, it should do so automatically. But if not try `F1` and type `Clang-Tidy: Lint File`. You can monitor behavior by `View->Output` tab in VS code -> select `Clang-Tidy` from the Drop down menu top right. 

For any header include that is undiscovered by Clang-Tidy (e.g. `error: 'bstrlib.h' file not found [clang-diagnostic-error]`):

**note! note to be confused with `cannot open source file` of type `C/C++` -> this indicates that Intellisense is missing includes, see [below](#Fixing-Intellisense) for how to fix.**

- Fix locally in docker and test
  - Update of `clang-tidy.compilerArgs` in your **container** settings.json file
    - Press `F1` and then search for `Preferences: Open Remote Settings (Dev Container)` Hit enter
    - Update the `clang-tidy.compilerArgs` with additional include directory paths
      - These may be in Magma codebase, or ...
      - environmental within docker (`find / -name <header_name.h> to discover`)
  - Update of settings.json for Dev Container should immediately take effect for Clang-Tidy
- Once tested locally, add to .devcontainer for future users
  - Open `.devcontainer/devcontainer.json` and update its `clang-tidy.compilerArgs`

### Fixing Intellisense

If you see `View->Problems` pane outputs of the shape `cannot open source file "bstrlib.h" C/C++`, this is an Intellisense warning.  We want to work on cleaning these up as well.

To address these, we want to modify the repo's `.vscode/c_cpp_properties.json` file to add the additional Include paths.  Test the change in your Docker container, then upstream to the repo!

### Using Clang-Tidy broadly

To run it over the entire code base, it's advised to use [run-clang-tidy.py](https://github.com/llvm-mirror/clang-tools-extra/blob/master/clang-tidy/tool/run-clang-tidy.py) provided by LLVM team. This handles concurrency of analysis and edits to the many files.

This tool wants a `compile_commands.json` list of compiler directives, which we can achieve by setting an environment variable for CMake: [`CMAKE_EXPORT_COMPILE_COMMANDS`](https://cmake.org/cmake/help/latest/variable/CMAKE_EXPORT_COMPILE_COMMANDS.html).

```shell
cd lte/gateway
make build_oai
cd /build/c/
wget https://raw.githubusercontent.com/llvm-mirror/clang-tools-extra/master/clang-tidy/tool/run-clang-tidy.py
python run-clang-tidy.py -checks='-*,clang-analyzer-security*,android-*,cert-*,clang-analyzer-*,concurrency-*,misc-*,bugprone-*' 2>&1 | tee clang-tidy.findings | grep warning:
```

Note that we are only using clang-tidy fimndings that are possibly bugs / alarming (or at least htat is the intent) - and are diabling the readability / formatting / code structure best practices.  We should come back and look at what automatic fixes clang-tidy can apply in these domains.

Detailed outputs will be available in clang-tidy.findings.

For an explanation of all checks see [Clang-Tidy Documentation](https://clang.llvm.org/extra/clang-tidy/checks/list.html).

Note that some builds are not generating build_commands.json - even with CMake flags set in CMakeFiles.txt. I have not yet found an explanation - but one solution is to manually build the projects like so:

```bash
cd lte/gateway/c/session_manager
mkdir build
cd build
cmake ../
vim CMakeCache.txt -> edit the CMAKE_EXPORT_COMPILE_COMMANDS to set it ON
cmake --build .
```

### Converting Clang-Tidy Output to CI Friendly things

```shell
root@ecee08edef4b:/build/c/oai# cat clang.findings | egrep "android|bugprone|cert|clang|concurrency|misc" | awk -F'[][]' '{print $2}' | sort | uniq -c
   1315 
      3 android-cloexec-fopen
      1 android-cloexec-open
     34 bugprone-branch-clone
     43 bugprone-macro-parentheses
    338 bugprone-narrowing-conversions
    219 bugprone-reserved-identifier,cert-dcl37-c,cert-dcl51-cpp
      9 bugprone-signed-char-misuse,cert-str34-c
     92 bugprone-sizeof-expression
      9 bugprone-suspicious-string-compare
      9 bugprone-too-small-loop-variable
      9 bugprone-undefined-memory-manipulation
      3 cert-dcl16-c
      7 cert-env33-c
     66 cert-err34-c
     64 cert-err58-cpp
      1 cert-msc30-c,cert-msc50-cpp
      1 cert-msc32-c,cert-msc51-cpp
     59 clang-analyzer-core.CallAndMessage
      6 clang-analyzer-core.NonNullParamChecker
     56 clang-analyzer-core.NullDereference
      3 clang-analyzer-core.UndefinedBinaryOperatorResult
      4 clang-analyzer-core.uninitialized.Assign
     69 clang-analyzer-deadcode.DeadStores
      3 clang-analyzer-optin.cplusplus.VirtualCall
      1 clang-analyzer-optin.performance.Padding
      4 clang-analyzer-optin.portability.UnixAPI
     23 clang-analyzer-security.insecureAPI.strcpy
     12 clang-analyzer-unix.Malloc
     16 clang-analyzer-unix.MallocSizeof
      5 misc-misplaced-const
     11 misc-no-recursion
      1 misc-non-private-member-variables-in-classes
      5 misc-redundant-expression
     15 misc-unused-using-decls
```



## Sharp Edges

- Your Docker Container Volume contains your work!
  - Be careful not to delete the container until outstanding work is pushed to your personal Github Fork
  - Note cleanup for disk space below - and the danger here
  - TODO: How to re-attach a Volume if you close VS Code?
- Docker running out of disk on Mac VM
  - Docker just leaks resources like mad
  - Eventually you'll get a "out of space" error
  - Or an error about apt-get failure
  - TODO: find good docker command to wipe everything

## TODOs

- Get the Docker image up to DockerHub?
  - Want easy way to build from Dockerfile locally though..