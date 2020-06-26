#!/usr/bin/env bash

# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

# self_sign_certs.sh generates a set of keys and self-signed certificates.

usage() {
  echo "Usage: $0 DOMAIN_NAME"
  exit 2
}

domain="$1"
if [[ ${domain} == "" ]]; then
    usage
fi

echo "##################"
echo "Creating root CA.."
echo "##################"
openssl genrsa -out rootCA.key 2048
openssl req -x509 -new -nodes -key rootCA.key -sha256 -days 3650 \
      -out rootCA.pem -subj "/C=US/CN=rootca.$domain"

echo "##########################"
echo "Creating controller cert.."
echo "##########################"
openssl genrsa -out controller.key 2048
openssl req -new -key controller.key -out controller.csr \
      -subj "/C=US/CN=*.$domain"
openssl x509 -req -in controller.csr -CA rootCA.pem -CAkey rootCA.key \
      -CAcreateserial -out controller.crt -days 3650 -sha256

rm -f controller.csr rootCA.key rootCA.srl
