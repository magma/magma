#!/usr/bin/env bash

echo "#######################"
echo "Creating certifier CA.."
echo "#######################"
openssl genrsa -out certifier.key 2048
openssl req -x509 -new -nodes -key certifier.key -sha256 -days 365000 \
      -out certifier.pem -subj "/C=US/CN=certifier.$1"

echo "###########################"
echo "Creating bootstrapper key.."
echo "###########################"
openssl genrsa -out bootstrapper.key 2048
