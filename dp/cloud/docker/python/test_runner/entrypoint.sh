#!/bin/bash
set -e

if [ "$TESTS_DEV" = true ] ; then
    echo "This container will sleep infinitely when TESTS_DEV env is set to true, it's for tests development"
    sleep infinity
else
    python "$@"
fi
