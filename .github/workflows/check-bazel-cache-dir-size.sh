#!/bin/bash

# This script should be run from $MAGMA_ROOT

# Relative path to cache directory from $MAGMA_ROOT
RELATIVE_CACHE_DIR=$1
CUTOFF_MB=$2

# See https://stackoverflow.com/a/27485157 for reference.
CACHE_SIZE_MB=$(du -smc "$RELATIVE_CACHE_DIR" | grep "$RELATIVE_CACHE_DIR" | cut -f1)
echo "Total size of Bazel cache (rounded up to MBs): $CACHE_SIZE_MB"

if [[ "$CACHE_SIZE_MB" -gt "$CUTOFF_MB" ]]; then
    echo "Cache exceeds cut-off; resetting it (will result in a slow build)"
    rm -rf "$RELATIVE_CACHE_DIR"
fi
