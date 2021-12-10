---
id: version-1.0.1-rest_apis
title: Swagger UI for REST APIs
hide_title: true
original_id: rest_apis
---
# Swagger UI for REST APIs
We use [Swagger](https://swagger.io/) for defining the north bound REST APIs, and the APIs can be viewed and tested using the Swagger UI. To use the UI:

```console
HOST$ open magma/.cache/test_certs
```

This will open up a finder window. Double-click the `admin_operator.pfx` cert
in this directory, which will open up Keychain to import the cert. The
password for the cert is `magma`. If you use Chrome or Safari, this is all you
need to do. If you use Firefox, copy this file to your desktop, then go to
`Preferences -> PrivacyAndSecurity -> View Certificates -> Import` and select
it.

Linux/Windows users should replace the above steps with the system-appropriate
method to import a client cert.

You can access the orchestrator REST API at https://127.0.0.1:9443/apidocs.
The SSL cert is self-signed, so click through any security warnings your
browser gives you. You should be prompted for a client cert, at which point
you should select the `admin_operator` cert that you added to Keychain above.
