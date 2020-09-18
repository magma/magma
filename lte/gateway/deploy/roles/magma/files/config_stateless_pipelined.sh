#!/bin/bash
source /usr/local/bin/config_stateless_utils.sh

if [[ $1 == "check" ]]; then
  # check dependency in systemd files of other services
  check_systemd_file "magma@pipelined" "magma@mme"

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
  remove_systemd_override "magma@mme"

  # change clean_restart setting in pipelined.yml
  sed -e '/clean_restart/ s/false/true/' -i /etc/magma/pipelined.yml
elif [[ $1 == "enable" ]]; then
  echo "Enabling stateless pipelined config"
  # remove restart dependencies between pipelined and other services
  add_systemd_override "magma@pipelined" "magma@mme"

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
