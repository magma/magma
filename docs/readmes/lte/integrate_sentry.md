---
id: integrate_sentry
title: Sentry Integration
hide_title: true
---

# Sentry + Magma AGW

[Sentry](https://sentry.io/welcome/) is an error monitoring platform.

As part of v1.5, we have integrated Sentry's [Python](https://docs.sentry.io/platforms/python/) SDK and [Native](https://docs.sentry.io/platforms/native/) (C/C++) SDK to enable basic error / stacktrace collection. 

## Setting up a Sentry platform

There are two main ways to deploy the platform: [self-hosting](https://develop.sentry.dev/self-hosted/) and paying for [Sentry.io](https://sentry.io/pricing/). 

## Enabling error reporting on an AGW

Fill out the following fields in [`control_proxy.yml`](https://github.com/magma/magma/blob/master/lte/gateway/configs/control_proxy.yml#L43) to enable Sentry reporting.
Reporting for Python services, MME, and SessionD will *only* be enabled if the corresponding URL fields are non-empty.
```
# [Experimental] Sentry related configs
# If set, the Sentry Python SDK will be initialized for all python services
sentry_url_python: ""
# If set, the Sentry Native SDK will be initialized for MME and SessionD
sentry_url_native: ""
# If set, /var/log/mme.log will be uploaded along MME crashreports
sentry_upload_mme_log: false
sentry_sample_rate: 1.0
```

## Uploading debug information files for SessionD and MME
To fully enable stacktrace analysis for any services running with the Native SDK, you will need to upload the corresponding debug information files for the stacktraces to be readable in Sentry. 
For more information, please refer to [this documentation provided by Sentry](https://docs.sentry.io/platforms/android/data-management/debug-files/). 
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


