---
id: version-1.4.0-architecture_security
title: Security
hide_title: true
original_id: architecture_security
---

# Security Overview

This document gives an overview of Orc8r's security architecture, with specific focus on the northbound (REST API) and southbound (gateway)
interfaces. For more in-depth information, see [Security Debugging](./dev_security.md).

## Architecture

Orc8r security is enforced in a two-tiered approach

- mutual TLS termination at the Orc8r proxy
- application-level access enforcement

Every request to Orc8r is mutually-authenticated via client-validated TLS. This validation occurs at the Orc8r proxy, where the request's
client certificate is converted to a set of trusted metadata in the proxied request. From there, individual services can make
application-level decisions on access enforcement, often by making requests to the `accessd` Orc8r service.

The one exception to client-validated endpoints is the southbound bootstrapper endpoint. Orc8r provides a mechanism for registered gateways
to securely bootstrap their authentication process via an open (server-validated only) service endpoint granting access to the
`bootstrapper` service. This endpoint allows gateways to periodically refresh their session certificates, which grant access to the full
southbound interface.

For the northbound interface (NMS and direct usage of the REST API), requests must include a provisioned client certificate. This
certificate can be generated using the `accessc` CLI from any Orc8r pod.

## Gateway bootstrap

Gateways generate, or can be provisioned with, two pieces of identifying information

- *hardware ID* is the gateway's unique identifier
- *challenge key* is a long-term keypair used for the bootstrap process

Operators must include two additional pieces of information to allow the gateway to connect to Orc8r

- `control_proxy.yml` is a config file providing the location of the Orc8r
- `rootCA.pem` is a CA certificate providing the root of trust for validating Orc8r during TLS handshakes

Operators use the `show_gateway_info.py` script to read the hardware ID and public part of the challenge key. They then register the gateway
with the cloud (e.g. via NMS), including these pieces of information, and provision the gateway with the correct values in
`control_proxy.yml` and `rootCA.pem`. From here, the gateway can begin the bootstrap process.

The bootstrap process involves two sequential requests to the `bootstrapper` service

- `GetChallenge` gateway requests a challenge from Orc8r
- `RequestSign` gateway sends signed challenge and a CSR

If the bootstrap process is successful, Orc8r signs and returns the CSR from the `RequestSign` request, which the gateway can now use as
its session certificate. Orc8r enforces a max validity period, so gateways must periodically re-bootstrap to receive updated session
certificates.

Operators can use the `checkin_cli.py` script to debug issues with this process of checking in and bootstrapping.

## Application-layer access control

Once a request passes TLS termination at the Orc8r proxy, Orc8r also supports finer-grained application-level access control. In Orc8r,
an ACL defines a mapping between identities, where each mapping includes read/write permissions. There are 3 types of identities

- *Gateway* identities represent a gateway
- *Operator* identities represent a network operator
- *Network* identities represent access to a particular network

At the Orc8r proxy, client certificates are translated to trusted request metadata. Middleware can then convert the metadata, via calls to
the `certifier` service, to an identity.

At the *southbound interface*, only gateway identities are allowed to make requests. Unregistered gateways are blocked.

At the the *northbound interface*, only operator identities are allowed to make requests. Operators must have ACLs collectively covering
the requested REST resources.

For more information, see the following servicer definitions
- [`access.proto`](https://github.com/magma/magma/blob/master/orc8r/cloud/go/services/accessd/protos/access.proto) access control
- [`certifier.proto`](https://github.com/magma/magma/blob/master/orc8r/cloud/go/services/certifier/protos/certifier.proto) certificate management
