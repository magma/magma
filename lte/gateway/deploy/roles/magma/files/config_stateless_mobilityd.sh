#!/bin/bash

if [[ $1 == "check" ]]; then
  # check dependency in systemd files of other services
  if ! grep -q "^#PartOf=magma@mobilityd" /etc/systemd/system/magma@mme.service
  then
    echo "MME will restart with mobilityd, i.e. stateful."
    exit 1
  fi

  #check service config
  if ! grep -q "persist_to_redis.*true" /etc/magma/mobilityd.yml; then
    echo "Mobilityd config file is stateful."
    exit 1
  fi

  echo "Mobilityd service is stateless."
  exit 0
elif [[ $1 == "disable" ]]; then
  echo "Disabling stateless mobilityd config"
  # restore restart dependencies between mobilityd and other services
  sudo sed -e '/PartOf=magma@mobilityd/ s/^#*//' -i \
    /etc/systemd/system/magma@mme.service

  # change persist_to_redis setting in mobilityd.yml
  sed -e '/persist_to_redis/ s/true/false/' -i /etc/magma/mobilityd.yml
elif [[ $1 == "enable" ]]; then
  echo "Enabling stateless mobilityd config"
  # remove restart dependencies between mobilityd and other services
  sudo sed -e '/PartOf=magma@mobilityd/ s/^#*/#/' -i \
    /etc/systemd/system/magma@mme.service

  # change persist_to_redis setting in mobilityd.yml
  sed -e '/persist_to_redis/ s/false/true/' -i /etc/magma/mobilityd.yml
else
  echo "Invalid argument. Use one of the following"
  echo "check: Run a check whether Mobilityd is stateless or not"
  echo "enable: Enable stateless mode, do nothing if already stateless"
  echo "disable: Disable stateless mode, do nothing if already stateful"
  exit 0
fi

# reload systemd config
sudo systemctl daemon-reload
