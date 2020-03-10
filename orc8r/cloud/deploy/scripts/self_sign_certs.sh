#!/usr/bin/env bash

echo "##################"
echo "Creating root CA.."
echo "##################"
openssl genrsa -out rootCA.key 2048
openssl req -x509 -new -nodes -key rootCA.key -sha256 -days 365000 \
      -out rootCA.pem -subj "/C=US/CN=rootca.$1"

echo "##########################"
echo "Creating controller cert.."
echo "##########################"
openssl genrsa -out controller.key 2048
openssl req -new -key controller.key -out controller.csr \
      -subj "/C=US/CN=*.$1"
openssl x509 -req -in controller.csr -CA rootCA.pem -CAkey rootCA.key \
      -CAcreateserial -out controller.crt -days 36400 -sha256

rm -f controller.csr rootCA.key rootCA.srl
