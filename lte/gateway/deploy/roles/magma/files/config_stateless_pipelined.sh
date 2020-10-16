#!/bin/bash

if [[ $1 == "check" ]]; then
  #check service config
  if ! grep -q "clean_restart.*false" /etc/magma/pipelined.yml; then
    echo "Pipelined config file is stateful."
    exit 1
  fi

  echo "Pipelined service is stateless."
  exit 0
elif [[ $1 == "disable" ]]; then
  echo "Disabling stateless pipelined config"
  # change clean_restart setting in pipelined.yml
  sed -e '/clean_restart/ s/false/true/' -i /etc/magma/pipelined.yml
elif [[ $1 == "enable" ]]; then
  echo "Enabling stateless pipelined config"
  # change clean_restart setting in pipelined.yml
  sed -e '/clean_restart/ s/true/false/' -i /etc/magma/pipelined.yml
else
  echo "Invalid argument. Use one of the following"
  echo "check: Run a check whether Pipelined is stateless or not"
  echo "enable: Enable stateless mode, do nothing if already stateless"
  echo "disable: Disable stateless mode, do nothing if already stateful"
  exit 0
fi
