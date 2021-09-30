---
id: configure_sentry
title: Sentry Integration
hide_title: true
---

# Sentry + Magma AGW

[Sentry](https://sentry.io/welcome/) is an error monitoring platform.

As part of v1.5, we have integrated Sentry's [Python](https://docs.sentry.io/platforms/python/) SDK and [Native](https://docs.sentry.io/platforms/native/) (C/C++) SDK to enable basic error / stacktrace collection.

## Setting up a Sentry platform

There are two main ways to deploy the platform: [self-hosting](https://develop.sentry.dev/self-hosted/) and paying for [Sentry.io](https://sentry.io/pricing/).

## Enabling error reporting on an AGW

Reporting for Python services, MME, and SessionD will *only* be enabled if the corresponding URL fields are non-empty.

The URL fields can be set in a network-wide configuration through the Orc8r's Swagger endpoint at `/networks/{network_id}`. You can also set or override Sentry configuration for specific gateways through the Orc8r's `/networks/{network_id}/gateways/{gateway_id}` endpoint. Finally, you can set them locally on a gateway in the `control_proxy.yml` file by filling out the following fields.

```yaml
# Sentry related configs
# If set, the Sentry Native SDK will be initialized for MME and SessionD
sentry_url_native: ""
# If set, /var/log/mme.log will be uploaded along MME crashreports
sentry_upload_mme_log: false
# Rate at which we want to sample Python error events
# Should be a number between 0 (0% of errors sent) and 1 (100% of errors sent)
sentry_sample_rate: 1.0

# [Experimental] If set, the Sentry Python SDK will be initialized for all python services
sentry_url_python: ""
```

Note that Sentry configuration set in `control_proxy.yml` takes highest precedence followed by gateway-specific Orc8r configuration and, finally network-wide Orc8r configuration.

## Uploading debug information files for SessionD and MME

To fully enable stacktrace analysis for any services running with the Native SDK, you will need to upload the corresponding debug information files for the stacktraces to be readable in Sentry.
For more information, please refer to [this documentation provided by Sentry](https://docs.sentry.io/platforms/android/data-management/debug-files/).

In order for the build artifacts to have a debug section, they have to be built with the `-g` compiler flag. The flag is set by default for all C/C++ services, so no additional change is necessary for this.

The follwing script outlines the necessary steps.

```bash
#!/bin/bash

### TODO: fill in all variables!
# Doc on auth tokens: https://docs.sentry.io/product/cli/configuration/
export SENTRY_AUTH_TOKEN=""
ORG=""
PROJECT=""

# Absolute paths to executables separated by a space
# Ex: EXECS_PATHS="/usr/local/bin/mme /usr/local/bin/sessiond"
EXECS_PATHS=""

if ! command -v sentry-cli &> /dev/null
then
    curl -sL https://sentry.io/get-cli/ | bash
fi

TMP_DIR=`mktemp -d -t symbol_upload.XXXXXX`

for EXEC_PATH in $EXECS_PATHS
do
    EXEC_NAME=`basename $EXEC_PATH`
    TMP_EXEC="$TMP_DIR/$EXEC_NAME"
    cp "$EXEC_PATH" "$TMP_EXEC"

    # Strip artifacts
    objcopy --only-keep-debug "$TMP_EXEC" "$TMP_EXEC".debug
    objcopy --strip-debug --strip-unneeded "$TMP_EXEC"
    objcopy --add-gnu-debuglink="$EXEC".debug "$TMP_EXEC"

    # [Optional] Log included debug information
    sentry-cli difutil check "$TMP_EXEC"
    sentry-cli difutil check "$TMP_EXEC".debug

    # Upload the debug artifact with `symtab`, `debug`, and `sources`
    sentry-cli upload-dif --log-level=info --org="$ORG" --project="$PROJECT" --include-sources  "$TMP_EXEC".debug
    # Upload the stripped executable with `unwind`
    sentry-cli upload-dif --log-level=info --org="$ORG" --project="$PROJECT" "$TMP_EXEC"
done

rm -rf $TMP_DIR
```
