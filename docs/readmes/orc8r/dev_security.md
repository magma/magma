---
id: dev_security
title: Security
hide_title: true
---

# Security Debugging

This document describes some tools for understanding and debugging Orc8r security. Read the [Security Overview](./architecture_security) first for
a top-level understanding of the security architecture.

## Architecture

By default, TLS is mutually terminated at the edge of both Orc8r and gateways

- Orc8r terminates at Orc8r proxy
    - Server certificate: `controller.crt`
    - Client-validation certificate: `certifier.pem`
- Gateways terminate at their `control_proxy` service
    - Server-validation certificate: `rootCA.pem`
    - Client certificate: `gateway.crt`

In the default configuration, with a registered and fully bootstrapped gateway, TLS termination looks like this

![Orc8r proxying](assets/orc8r/orc8r_proxying.png)

### Certificates

Starting point for understanding the set of certificates related to Orc8r

- General
    - `rootCA.pem` CA certificate which signed `controller.crt` (root of trust)
        - When Orc8r self-signs its certificates, the `rootCA.key` will also be included in the deployment
    - `admin_operator.{key.pem,pem}` admin operator certificates for full access to Orc8r's northbound interface by NMS and CLIs
- Orc8r
    - `controller.{key,crt}` server-validation certificate for Orc8r and NMS, signed by `rootCA.pem`
    - `certifier.pem` Orc8r's client-validation certificate (root of trust)
    - `bootstrapper.key` Orc8r's signing key used in the bootstrap process
    - `fluentd.{key,pem}` Orc8r's certificates for its fluentd endpoints (fluentd is currently outside Orc8r proxy)
- Gateway
    - `gw_challenge.key` gateway's long-term key used for the bootstrap process
    - `gateway.{key,crt}` gateway's session certificates, used as client certificates to Orc8r proxy

There are two roots of trust in Orc8r, as identified above

- `rootCA.pem`
    - Root of trust for server validation
    - `controller.crt`, signed by `rootCA.pem`, is presented by Orc8r during TLS handshakes
- `certifier.pem`
    - Root of trust for client validation
    - `gateway.crt`, signed by `certifier.pem`, is presented by gateways during TLS handshakes

### Endpoints

Default Orc8r endpoints, as exposed on Orc8r proxy

- `:80` K8s liveness probe (open)
- `:443` northbound interface exposing a REST API (mutually-authenticated)
    - Orc8r proxy exposes this over `:9443`, which Helm then maps to `:443`
- `:8443` general southbound interface (mutually-authenticated)
- `:8444` bootstrapper southbound interface (server-authenticated)

## Gateway bootstrap

The gateway bootstrap process allows registered gateways to securely request signed session certificates from Orc8r, without Orc8r needing
to maintain a stateful view of the bootstrap process.

Before the bootstrap process, operators provision a gateway with Orc8r by including the gateway's hardware ID and the public version
of the gateway's `gw_challenge.key`. Then, during the two-phase bootstrap process, the gateway requests a challenge, signs and returns the
challenge along with a CSR, then receives back signed session certificates.

### Get challenge

Performed via the `GetChallenge` RPC, where a gateway requests a challenge over a server-validated endpoint.

- Define `nonce := randomBytes || timestamp`
- Define `signature := sign(sha256_digest( randomBytes || timestamp ))`, computed using Orc8r's `bootstrapper.key`
- Define `challenge := nonce || signature`
- Orc8r returns `challenge` to gateway

### Get signed session certificates

Performed via the `RequestSign` RPC, where a gateway presents the gateway-signed challenge and a CSR, and receives back the
Orc8r-signed CSR.

- Define `signed_challenge := sign(sha256_digest( challenge ))`, computed using gateway's `gw_challenge.key`
- Gateway sends `challenge`, `signed_challenge`, and a CSR
- Orc8r validates request
    - Validate `challenge` by verifying the original embedded signature and the validity window
    - Validate `singed_challenge` by verifying the gateway's signature of `challenge`
- Orc8r returns the CSR as the gateway's session certificates, signed by `certifier.pem`

### Implementation details

Some implementation-specific details, useful for debugging

- Challenge is valid for short period of time (5 minutes)
- Signed CSR includes capped validity duration (4 days)
    - Gateways proactively reach out before expiration to re-bootstrap
- Orc8r supports the `echo` bootstrap mode on a per-gateway basis, where, for testing purposes, gateways simply echo back the challenge
  to receive session certs
- Since Orc8r retains no state across the `GetChallenge` and `RequestSign` calls, gateways can make multiple `RequestSign` requests with
  the same challenge. This is admissible because session certificates are only granted if the requesting gateway possesses the
  requisite private key.

## Debug tools

This section includes a selection of notes and commands useful for debugging Orc8r from a security perspective.

### Notes

- [List of additional HTTP/2 debug tools](https://blog.cloudflare.com/tools-for-debugging-testing-and-using-http-2/)
- Certificate serial number is expected to be unique, but
  [uniqueness is *not* guaranteed](https://stackoverflow.com/questions/9104108/is-serial-number-a-unique-key-for-x509-certificate)


### Gateway information

- `show_gateway_info.py` show gateway's hardware ID and (public part of) `gw_challenge.key`
- `checkin_cli.py` attempt to check in to Orc8r, outputting debug information on failure

### CLI commands

These commands have been useful for debugging Orc8r's northbound and southbound interfaces. These commands are presented targeting a local
dev Orc8r spun up by the default docker-compose file, but with minor modifications they can target arbitrary production deployments.

Utilized tools

- [`curl`](https://curl.se/)
- [`grpcurl`](https://github.com/fullstorydev/grpcurl)
- [`s_client`](https://www.openssl.org/docs/man1.0.2/man1/openssl-s_client.html)

View certificates and ACLs

```bash
# These commands are run from one of Orc8r's application containers/pods

# List all registered operators, along with their ACLs and certificate serial numbers
/var/opt/magma/bin/accessc list

# List all registered certificates, along with their associated identities
/var/opt/magma/bin/accessc list-certs
```

Debug northbound interface (REST API)

```bash
# Make verbose request to REST API ("list networks" endpoint)
curl \
    --insecure \
    --verbose \
    --key ${MAGMA_ROOT}/.cache/test_certs/admin_operator.key.pem \
    --cert ${MAGMA_ROOT}/.cache/test_certs/admin_operator.pem \
    'https://localhost:9443/magma/v1/networks' \
    -H 'accept: application/json'

# Open connection to REST API, logging certificate and handshake info
# See below for some HTTP requests to use once the connection is open
openssl s_client \
    -showcerts \
    -CAfile ${MAGMA_ROOT}/.cache/test_certs/rootCA.pem \
    -key ${MAGMA_ROOT}/.cache/test_certs/admin_operator.key.pem \
    -cert ${MAGMA_ROOT}/.cache/test_certs/admin_operator.pem \
    -connect localhost:9443
```

Debug southbound interface (gRPC)

```bash
# Emulate a full gateway bootstrap
#
# Places generated session secrets at /tmp/magma_protos/gateway.{key,crt}
${MAGMA_ROOT}/orc8r/tools/scripts/bootstrap.bash YOUR_PROVISIONED_GATEWAY_HWID

# Emulate gateway request to bootstrapper service's GetChallenge endpoint (unprotected)
grpcurl \
    -insecure \
    -authority bootstrapper-controller.magma.test \
    -proto /tmp/magma_protos/orc8r/protos/bootstrapper.proto \
    -import-path /tmp/magma_protos/ \
    -d '{"id": "YOUR_PROVISIONED_GATEWAY_HWID"}' \
    localhost:7444 \
    magma.orc8r.Bootstrapper/GetChallenge

# Emulate gateway request to state service's ReportStates endpoint (protected)
#
# gateway.{key,crt} copied from target gateway
grpcurl \
    -insecure \
    -key /tmp/magma_protos/gateway.key \
    -cert /tmp/magma_protos/gateway.crt \
    -authority state-controller.magma.test \
    -protoset /tmp/magma_protos/out.protoset \
    -d '{"states": [{"type": "test"}]}' \
    localhost:7443 \
    magma.orc8r.StateService/ReportStates

# Emulate gateway request to state service's GetStates endpoint (protected)
#
# gateway.{key,crt} copied from target gateway
#
# Note: this is an Orc8r-internal endpoint. To run this outside an Orc8r
# application container, you may have to temporarily mutate some access
# enforcment configuration in your dev setup. Alternatively, you could
# run this command from inside an Orc8r application pod.
grpcurl \
    -insecure \
    -key /tmp/magma_protos/gateway.key \
    -cert /tmp/magma_protos/gateway.crt \
    -authority state-controller.magma.test \
    -protoset /tmp/magma_protos/out.protoset \
    -d '{"networkID": "test"}' \
    localhost:7443 \
    magma.orc8r.StateService/GetStates
```

Some HTTP requests to use with `s_client`

```
GET / HTTP/1.1
Host: localhost:9443

GET /magma/v1/networks HTTP/1.1
Host: localhost:9443
```

Miscellaneous commands

```bash
# Examine a certificate copied to clipboard (macOS)
openssl x509 -text -noout -in <(pbpaste)

# Verify controller's certificate
openssl verify \
    -CAfile ${MAGMA_ROOT}/.cache/test_certs/rootCA.pem \
    ${MAGMA_ROOT}/.cache/test_certs/controller.crt

# Verify gateway's session certificate
openssl verify \
    -CAfile ${MAGMA_ROOT}/.cache/test_certs/certifier.pem \
    gateway.crt  # copied from target gateway

# Consolidate Magma and its imported protos, copying to /tmp/magma_protos/
${MAGMA_ROOT}/orc8r/tools/scripts/consolidate_protos.sh

# Generate compiled proto definitions
cd /tmp/magma_protos/ && \
protoc \
    -I . \
    -I orc8r/protos/prometheus/ \
    --descriptor_set_out out.protoset \
    --include_imports \
    orc8r/cloud/go/services/state/protos/indexer.proto  # or target proto
```
