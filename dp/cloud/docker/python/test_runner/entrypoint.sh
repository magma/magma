#!/bin/bash
set -e

if [ "$TESTS_DEV" = true ] ; then
    sleep infinity
else
    python "$@"
fi
