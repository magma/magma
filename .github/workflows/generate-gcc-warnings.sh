#!/bin/bash

# List of files or paths to grep for in the compilation log
# Ex: .github/workflows/gcc-problems.yml,.github/workflows/generate-gcc-warnings.sh,lte/gateway/c/sctpd/src/sctpd.cpp,lte/gateway/c/session_manager/AAAClient.cpp
FILES=$1
# Applied to both C and C++. 
# Ex: "-Wextra -Wshadow -Wimplicit-fallthrough"
C_CPP_FLAGS=$2
# Applied to C only. 
# Ex: "-Wjump-misses-init"
C_ONLY_FLAGS=$3

# We want the full compilation log everytime the script is run
bazel clean

# Massage some of our external build files so that we can bypass -Werror
bazel fetch //orc8r/gateway/c/... //lte/gateway/c/...
sed -i 's/"-Werror",/#"-Werror",/g' /tmp/bazel/external/boringssl/BUILD

COPTS=${C_CPP_FLAGS//'-W'/'--copt=-W'}
CONLY_OPTS=${C_ONLY_FLAGS//'-W'/'--conlyopt=-W'}

# shellcheck disable=SC2086
bazel build --color=no //orc8r/gateway/c/... //lte/gateway/c/... $COPTS $CONLY_OPTS 2>&1 | tee compile.log

rm -f filtered-compile.log
echo "$FILES" | tr , '\n' | while read f
    do echo "$f"; grep "$f" compile.log >> filtered-compile.log;
done;

exit 0;
