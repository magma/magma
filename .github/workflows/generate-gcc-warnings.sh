#!/bin/bash

# List of files or paths to grep for in the compilation log
# Ex: .github/workflows/gcc-problems.yml,.github/workflows/generate-gcc-warnings.sh,lte/gateway/c/sctpd/src/sctpd.cpp,lte/gateway/c/session_manager/AAAClient.cpp
FILES=$1

# We want the full compilation log everytime the script is run
bazel clean

# shellcheck disable=SC2086
bazel build --color=no //orc8r/gateway/c/... //lte/gateway/c/... --config=max_gcc_warnings --profile=Bazel_build_gcc_problems_profile 2>&1 | tee compile.log

rm -f filtered-compile.log
echo "$FILES" | tr , '\n' | while read f
    do echo "$f"; grep "$f" compile.log >> filtered-compile.log;
done;

exit 0;
