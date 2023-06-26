---
id: configure_sentry
title: Sentry Integration
hide_title: true
---

# Sentry + Magma AGW

[Sentry](https://sentry.io/welcome/) is an error monitoring platform. Members of the Magma community are able to request access to link one access gateway in their lab to the Linux Foundation Sentry page. Alternatively, users can set up their own instance of Sentry outside of the Magma community.

As part of v1.5, we have integrated Sentry's [Native](https://docs.sentry.io/platforms/native/) (C/C++) SDK to enable basic error / stacktrace collection.

## Connect to the Linux Foundation Sentry Account

The first step to setting up your account is to request access to Linux Foundation [Sentry.io](https://sentry.io) instance by submitting your email address and company through the [community form](https://docs.google.com/forms/d/e/1FAIpQLSeJMWecw9An5-aYv0US8Fc_PDO7kUMx4Pky13S_3LhFJkge_g/viewform). You should also create a team for your company on the Sentry site. This team name should match your company and will be used to enable group access to your project. Once you are granted access, you will need to login and navigate to the “Projects” tab and locate your team’s project name. You will need this information to update a CI job in the next section.

## Configuration

To configure Sentry, you will need to create a pull request to update the config.yml on the Magma Github page. If you skip this step, all C/C++ will be unreadable for your Sentry instance. This file is located [here](https://github.com/magma/magma/blob/master/.github/workflows/composite/sentry-create-and-upload-artifacts/action.yml#L19). You will need to navigate to the "PROJECTS:" section of the file and add your project name in the following format:

```yaml
  PROJECTS:
    required: false
    default: ('lab-agws-native' 'magma-staging-native' 'NEW_PROJECT_NAME')
```

## Enabling error reporting on an AGW

Reporting for Python services, MME, and SessionD will *only* be enabled if the corresponding URL fields are non-empty.

The URL fields can be set in a network-wide configuration through the Orc8r's Swagger endpoint at `/networks/{network_id}` or `/networks/{network_id}/sentry`. You can also set them locally on a gateway in the `control_proxy.yml` file by filling out the following fields, which overrides the network-wide configuration.

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

## Client-side filters for error messages

Another configuration option is a list of error messages that should be filtered and not be sent to Sentry. This is done to avoid spamming Sentry with errors that can occur frequently but that are not bugs to be monitored.

This is similar to Sentry's [server-side filters](https://docs.sentry.io/product/accounts/quotas/#inbound-data-filters) that can be configured per project. Filtering on the client side can help avoid network traffic and load on the server.

The filters are configured by giving a list of regex exclusion patterns in Orc8r's Swagger endpoints (at `/networks/{network_id}` or `/networks/{network_id}/sentry`). If exclusion patterns are not included in the call, a list of default patterns is used. In contrast to the other configuration options for Sentry, the exclusion patterns cannot be overridden in `control_proxy.yml`.

Example payload for `/networks/{network_id}/sentry`:

```json
{
  "url_native": "https://something@o0.ingest.sentry.io/...",
  "url_python": "https://another@o0.ingest.sentry.io/...",
  "exclusion_patterns": [
    "\\[SyncRPC\\]",
    "Metrics upload error",
    "Streaming from the cloud failed!"
  ],
  "sample_rate": 1,
  "upload_mme_log": false,
  "number_of_lines_in_log": 0
}
```

## Log file support

There are two configuration options for sending log files to Sentry after a C++ service crash occurs. Both can be activated and deactivated via the Swagger UI. The first option is `upload_mme_log`. In the case it is set to `true` like the following example, the MME service log file located in `/var/log/mme.log` will be submitted within the crash report.

```json
"upload_mme_log": true,
```

The second option `number_of_lines_in_log` supports the transmission of the journal syslog file `/var/log/syslog`. It selects the last `n` entries in the file and sends them to Sentry attached to the crash report. The sent file considers log entries of all services marked with `magma@` and additionally of the service `sctpd`. In the subsequent example, `n` is 1000. If `n` is set to 0, `number_of_lines_in_log` is disabled, and no log history will be sent.

```json
"number_of_lines_in_log": 1000, 
```

Note: Only `upload_mme_log` can be overridden in the `control_proxy.yml`.

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
