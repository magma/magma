#!/bin/bash

set -e

if [[ -z "${MAGMA_ROOT}" ]]; then
  echo "Must set env var MAGMA_ROOT"
  exit 1
fi

mkdir -p "${MAGMA_ROOT}"/.cache/test_certs
export CERT_DIR="${MAGMA_ROOT}"/.cache/test_certs

# First, rootCA
openssl genrsa -out "${CERT_DIR}"/rootCA.key 2048
openssl req -x509 -new -nodes -key "${CERT_DIR}"/rootCA.key \
      -sha256 -days 3650 -out "${CERT_DIR}"/rootCA.pem \
            -subj "/C=US/ST=CA/L=Menlo Park/O=Facebook/OU=Magma/CN=rootca.magma.test/emailAddress=admin@magma.test"

# Then, controller key and crt
openssl genrsa -out "${CERT_DIR}"/controller.key 2048
openssl req -new -key "${CERT_DIR}"/controller.key -out "${CERT_DIR}"/controller.csr \
      -subj "/C=US/ST=CA/L=Menlo Park/O=Facebook/OU=Magma/CN=*.magma.test/emailAddress=admin@magma.test"
openssl x509 -req -in "${CERT_DIR}"/controller.csr -CA "${CERT_DIR}"/rootCA.pem \
      -CAkey "${CERT_DIR}"/rootCA.key -CAcreateserial -out "${CERT_DIR}"/controller.crt \
      -days 365 -sha256

# Remove unneeded intermediate files
ls "${CERT_DIR}"
rm "${CERT_DIR}"/controller.csr "${CERT_DIR}"/rootCA.key

# Now, certiifier key and pem and bootstrapper key
openssl genrsa -out "${CERT_DIR}"/certifier.key 2048
openssl req -x509 -new -nodes -key "${CERT_DIR}"/certifier.key -sha256 \
      -days 3650 -out "${CERT_DIR}"/certifier.pem \
            -subj "/C=US/ST=CA/L=Menlo Park/O=Facebook/OU=Magma/CN=certifier.magma.test/emailAddress=admin@magma.test"
openssl genrsa -out "${CERT_DIR}"/bootstrapper.key 2048

ls "${CERT_DIR}"
