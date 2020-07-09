#!/usr/bin/env bash

# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

# self_sign_certs.sh generates a set of keys and self-signed certificates.
#
# Generated secrets
#   - rootCA.key, rootCA.pem -- certs for trusted root CA (rootCA.key is
#     subsequently deleted)
#   - controller.key, controller.crt -- certs for orc8r deployment's public
#     domain name, signed by rootCA.key
#
# NOTE: extension naming for certs and keys is non-normalized here due to
# expected naming conventions from other parts of the code. For reference,
# all outputs are PEM-encoded, *.key is a private key, and *.crt and missing
# indicator are both certificates.

usage() {
  echo "Usage: $0 DOMAIN_NAME"
  exit 2
}

domain="$1"
if [[ ${domain} == "" ]]; then
    usage
fi

echo ""
echo "################"
echo "Creating root CA"
echo "################"
openssl genrsa -out rootCA.key 2048
openssl req -x509 -new -nodes -key rootCA.key -sha256 -days 3650 -out rootCA.pem -subj "/C=US/CN=rootca.$domain"

echo ""
echo "########################"
echo "Creating controller cert"
echo "########################"
openssl genrsa -out controller.key 2048
openssl req -new -key controller.key -out controller.csr -subj "/C=US/CN=*.$domain"
openssl x509 -req -in controller.csr -CA rootCA.pem -CAkey rootCA.key -CAcreateserial -out controller.crt -days 3650 -sha256

echo ""
echo "###########################"
echo "Deleting intermediate files"
echo "###########################"
rm -f controller.csr rootCA.srl

echo ""
echo "####################"
echo "Deleting root CA key"
echo "####################"
rm -f rootCA.key
