#!/usr/bin/env bash

# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

# create_application_certs.sh generates application certs for orc8r from
# existing certificates.

usage() {
  echo "Usage: $0 DOMAIN_NAME"
  exit 2
}

domain="$1"
if [[ ${domain} == "" ]]; then
    usage
fi

echo "#######################"
echo "Creating certifier CA.."
echo "#######################"
openssl genrsa -out certifier.key 2048
openssl req -x509 -new -nodes -key certifier.key -sha256 -days 3650 \
      -out certifier.pem -subj "/C=US/CN=certifier.$domain"

echo "###########################"
echo "Creating bootstrapper key.."
echo "###########################"
openssl genrsa -out bootstrapper.key 2048


echo "########################"
echo "Creating fluentd certs.."
echo "########################"
openssl genrsa -out fluentd.key 2048
openssl req -x509 -new -nodes -key fluentd.key -sha256 -days 3650 \
      -out fluentd.pem -subj "/C=US/CN=fluentd.$domain"
