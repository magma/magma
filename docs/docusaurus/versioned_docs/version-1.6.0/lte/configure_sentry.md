---
id: version-1.6.0-configure_sentry
title: Sentry Integration
hide_title: true
original_id: configure_sentry
---

# Sentry + Magma AGW
[Sentry](https://sentry.io/welcome/) is an error monitoring platform. Members of the Magma community are able to request access to link one access gateway in their lab to the Linux Foundation Sentry page. Alternatively, users can set up their own instance of Sentry outside of the Magma community.

As part of v1.5, we have integrated Sentry's [Native](https://docs.sentry.io/platforms/native/) (C/C++) SDK to enable basic error / stacktrace collection.

## Connect to the Linux Foundation Sentry Account
The first step to setting up your account is to request access to Linux Foundation [Sentry.io](http://Sentry.io) instance by submitting your email address and company through the [community form](https://docs.google.com/forms/d/e/1FAIpQLSeJMWecw9An5-aYv0US8Fc_PDO7kUMx4Pky13S_3LhFJkge_g/viewform). You should also create a team for your company on the Sentry site. This team name should match your company and will be used to enable group access to your project. Once you are granted access, you will need to login and navigate to the “Projects” tab and locate your team’s project name. You will need this information to update the CircleCI config.yml file on the Magma Github.

## Configuration

To configure Sentry, you will need to create a pull request to update the config.yml on the Magma Github page. If you skip this step, all C/C++ will be unreadable for your Sentry instance. This file is located [here](https://github.com/magma/magma/blob/master/.circleci/config.yml). You will need to navigate to the "sentry-create-and-upload-artifacts" section of the file and create a new a new sentry upload in the following format:

```bash
sentry-upload:
executable_name: << parameters.executable_name >>
project: [fill in your project here]
org: lf-9c
```

## Enabling error reporting on an AGW
Reporting for Python services, MME, and SessionD will *only* be enabled if the corresponding URL fields are non-empty.

Fill out the following fields in `control_proxy.yml` to enable Sentry reporting.

```bash
# [Experimental] Sentry related configs
# If set, the Sentry Python SDK will be initialized for all python services
sentry_url_python: ""
# If set, the Sentry Native SDK will be initialized for MME and SessionD
sentry_url_native: ""
# If set, /var/log/mme.log will be uploaded along MME crash reports
sentry_upload_mme_log: false
# A sampling rate to apply to events. A value of 0.0 will send no
# events, and a value of 1.0 will send all events (default).
sentry_sample_rate: 1.0
```
## Monitoring
Once your Sentry instance is established we recommend setting up alerts and your dashboard to monitor your results. This will allow you to understand the health of your network and also enable you to view key metrics such as crash rate per release.

Follow the Sentry dashboard instructions [here](https://docs.sentry.io/product/dashboards/) to understand the dashboard page, and use the instructions [here](https://docs.sentry.io/product/dashboards/custom-dashboards/) to set up a custom dashboard. In order to set up alerts to Slack or email, follow the Sentry instructions [here](https://docs.sentry.io/product/alerts/).

## Non-Linux Foundation Sentry Accounts: Uploading debug information files for SessionD and MME
To fully enable stacktrace analysis for any services running with the Native SDK, you will need to upload the corresponding debug information files for the stacktraces to be readable in Sentry.

For more information, please refer to [this documentation provided by Sentry](https://docs.sentry.io/platforms/android/data-management/debug-files/).

In order for the build artifacts to have a debug section, they have to be built with the `-g` compiler flag. The flag is set by default for all C/C++ services, so no additional change is necessary for this.

The following script outlines the necessary steps.

```bash
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

   # [Optional] Log included debug information
    sentry-cli difutil check "$EXEC"
    sentry-cli difutil check "$EXEC".debug


    # Upload the debug artifact with `symtab`, `debug`, and `sources`
    sentry-cli upload-dif --log-level=info --org="$ORG" --project="$PROJECT" --include-sources  "$EXEC".debug
    # Upload the stripped executable with `unwind`
    sentry-cli upload-dif --log-level=info --org="$ORG" --project="$PROJECT" "$EXEC"
done
```
