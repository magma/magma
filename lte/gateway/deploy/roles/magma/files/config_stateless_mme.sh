#!/bin/bash

if [[ $1 == "check" ]]; then
  #check service config
  if ! grep -q "use_stateless.*true" /etc/magma/mme.yml; then
    echo "MME config file is stateful."
    exit 1
  fi
  echo "MME service is stateless."
  exit 0
elif [[ $1 == "disable" ]]; then
  echo "Disabling stateless MME config"
  # change use_stateless setting in mme.yml
  sed -e '/use_stateless/ s/true/false/' -i /etc/magma/mme.yml
elif [[ $1 == "enable" ]]; then
  echo "Enabling stateless MME config"
  # change use_stateless setting in mme.yml
  sed -e '/use_stateless/ s/false/true/' -i /etc/magma/mme.yml
else
  echo "Invalid argument. Use one of the following"
  echo "check: Run a check whether MME is stateless or not"
  echo "enable: Enable stateless mode, do nothing if already stateless"
  echo "disable: Disable stateless mode, do nothing if already stateful"
  exit 0
fi
