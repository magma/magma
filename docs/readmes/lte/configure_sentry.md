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

```
# [Experimental] Sentry related configs
# If set, the Sentry Python SDK will be initialized for all python services
sentry_url_python: ""
# If set, the Sentry Native SDK will be initialized for MME and SessionD
sentry_url_native: ""
# If set, /var/log/mme.log will be uploaded along MME crashreports
sentry_upload_mme_log: false
# A sampling rate to apply to events. A value of 0.0 will send no
# events, and a value of 1.0 will send all events (default).
sentry_sample_rate: 1.0
```

Note that Sentry configuration set in `control_proxy.yml` takes highest precedence followed by gateway-specific Orc8r configuration and, finally network-wide Orc8r configuration.


## Uploading debug information files for SessionD and MME
To fully enable stacktrace analysis for any services running with the Native SDK, you will need to upload the corresponding debug information files for the stacktraces to be readable in Sentry.
For more information, please refer to [this documentation provided by Sentry](https://docs.sentry.io/platforms/android/data-management/debug-files/).

In order for the build artifacts to have a debug section, they have to be built with the `-g` compiler flag. The flag is set by default for all C/C++ services, so no additional change is necessary for this.

The follwing script outlines the necessary steps.

```
#!/bin/bash
# To install sentry-cli, run `curl -sL https://sentry.io/get-cli/ | bash`

TODO: change me!
EXECS_PATHS=""
ORG=""
PROJECT=""

for EXEC in $EXECS_PATHS
do
    # Strip artifacts
    objcopy --only-keep-debug "$EXEC" "$EXEC".debug
    objcopy --strip-debug --strip-unneeded "$EXEC"
    objcopy --add-gnu-debuglink="$EXEC".debug "$EXEC"

    # [Optional] Log included debug information
    sentry-cli difutil check "$EXEC"
    sentry-cli difutil check "$EXEC".debug


    # Upload the debug artifact with `symtab`, `debug`, and `sources`
    sentry-cli upload-dif --log-level=info --org="$ORG" --project="$PROJECT" --include-sources  "$EXEC".debug
    # Upload the stripped executable with `unwind`
    sentry-cli upload-dif --log-level=info --org="$ORG" --project="$PROJECT" "$EXEC"
done
```
