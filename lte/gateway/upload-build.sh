#!/bin/bash

C_BUILD=/home/vagrant/build/c
EXECS="session_manager/sessiond oai/oai_mme/mme sctpd/sctpd connection_tracker/connectiond"
for EXEC in $EXECS
do
    EXEC_PATH="${C_BUILD}/${EXEC}"
    echo "Uploading debug artifacts for $EXEC_PATH"

    # Strip artifacts
    objcopy --only-keep-debug "$C_BUILD/$EXEC" "$C_BUILD/$EXEC".debug
    objcopy --strip-debug --strip-unneeded $C_BUILD/$EXEC
    objcopy --add-gnu-debuglink=$C_BUILD/$EXEC.debug $C_BUILD/$EXEC

    # Upload artifacts to Sentry
    sentry-cli upload-dif --log-level=info --org=magma-sentry-testing --project=sentry-native-testing $C_BUILD/$EXEC.debug
    sentry-cli upload-dif --log-level=info --org=magma-sentry-testing --project=sentry-native-testing $C_BUILD/$EXEC
done
