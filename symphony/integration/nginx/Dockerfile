# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# Fetch base alpine image
FROM alpine AS builder

# Install openssl package
RUN apk add --no-cache openssl bash

# Temporary directory for certificate generation
WORKDIR /ssl

# The following 5 steps generate a self-signed Root Certificate Authority
# (Root CA) certificate and use that cert to sign a certificate for nginx.
# The Root CA cert can then be added to Android devices/emulators or Chrome
# to allow TLS with certs signed with a self-signed Root CA.

# Step 1: Create a Root Key
SHELL ["/bin/bash", "-c"]
RUN openssl genrsa -out rootCA.key 4096

# Step 2: Create a self-signed Root Certificate
# New certificate requirements for macOS 10.15: https://support.apple.com/en-us/HT210176
#   max days is 2 years
RUN openssl req -x509 -new -nodes -key rootCA.key -sha256 -days 729 -out rootCA.crt \
  -subj '/CN=localhost'

# Step 3: Create a certificate key
RUN openssl genrsa -out ssl.key 2048

# Step 4: Create a certificate signing request
RUN openssl req -new -sha256 \
  -key ssl.key \
  -subj '/CN=localhost' \
  -extensions EXT \
  -out ssl.csr \
  -config <( \
  printf "[dn]\nCN=localhost\n[req]\ndistinguished_name = dn\n[EXT]\nsubjectAltName=DNS:localhost\nkeyUsage=digitalSignature\nextendedKeyUsage=serverAuth")

# Step 5: Generate the self-signed certificate
# New certificate requirements for macOS 10.15: https://support.apple.com/en-us/HT210176
#   max days is 2 years
#   must have extendedKeyUsage (EKU)
#   must use subjectAltName
RUN openssl x509 -req -in ssl.csr -CA rootCA.crt -CAkey rootCA.key \
  -CAcreateserial -out ssl.crt -days 729 -sha256 -extfile <( \
    printf "authorityKeyIdentifier=keyid,issuer\nbasicConstraints=CA:FALSE\nkeyUsage=digitalSignature,nonRepudiation,keyEncipherment,dataEncipherment\nextendedKeyUsage=serverAuth,clientAuth\nsubjectAltName=@alt_names\n[alt_names]\nDNS.1=localhost\nDNS.2=*.localhost")

# Create reverse proxy image
FROM nginx:alpine

# Install runtime dependencies
RUN apk add --no-cache curl

# Copy generated certificate
COPY --from=builder /ssl /etc/nginx/ssl
RUN chmod 700 /etc/nginx/ssl

# Copy reverse proxy config
ARG config=default.nginx
COPY ${config} /etc/nginx/conf.d/default.conf
