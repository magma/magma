#!/usr/bin/env bash

echo "#######################"
echo "Creating certifier CA.."
echo "#######################"
openssl genrsa -out certifier.key 2048
openssl req -x509 -new -nodes -key certifier.key -sha256 -days 3650 \
      -out certifier.pem -subj "/C=US/CN=certifier.$1"

echo "###########################"
echo "Creating bootstrapper key.."
echo "###########################"
openssl genrsa -out bootstrapper.key 2048


echo "########################"
echo "Creating fluentd certs.."
echo "########################"
openssl genrsa -out fluentd.key 2048
openssl req -x509 -new -nodes -key fluentd.key -sha256 -days 3650 \
      -out fluentd.pem -subj "/C=US/CN=fluentd.$1"
