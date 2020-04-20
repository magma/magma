#!/bin/bash

if [[ $1 == "check" ]]; then
  # check dependency in systemd files of other services
  if ! grep -q "^#PartOf=magma@pipelined" /etc/systemd/system/magma@mme.service
  then
    echo "MME will restart with pipelined, i.e. stateful."
    exit 1
  fi

  #check service config
  if ! grep -q "clean_restart.*false" /etc/magma/pipelined.yml; then
    echo "Pipelined config file is stateful."
    exit 1
  fi

  echo "Pipelined service is stateless."
  exit 0
elif [[ $1 == "disable" ]]; then
  echo "Disabling stateless pipelined config"
  # restore restart dependencies between pipelined and other services
  sudo sed -e '/PartOf=magma@pipelined/ s/^#*//' -i \
    /etc/systemd/system/magma@mme.service

  # change clean_restart setting in pipelined.yml
  sed -e '/clean_restart/ s/false/true/' -i /etc/magma/pipelined.yml
elif [[ $1 == "enable" ]]; then
  echo "Enabling stateless pipelined config"
  # remove restart dependencies between pipelined and other services
  sudo sed -e '/PartOf=magma@pipelined/ s/^#*/#/' -i \
    /etc/systemd/system/magma@mme.service

  # change clean_restart setting in pipelined.yml
  sed -e '/clean_restart/ s/true/false/' -i /etc/magma/pipelined.yml
else
  echo "Invalid argument. Use one of the following"
  echo "check: Run a check whether Pipelined is stateless or not"
  echo "enable: Enable stateless mode, do nothing if already stateless"
  echo "disable: Disable stateless mode, do nothing if already stateful"
  exit 0
fi

# reload systemd config
sudo systemctl daemon-reload
