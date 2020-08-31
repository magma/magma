#!/bin/bash
source /usr/local/bin/config_stateless_utils.sh

MME_DEPS=("magma@mobilityd" "magma@pipelined" "magma@sessiond" "sctpd")

if [[ $1 == "check" ]]; then
  # check dependency in systemd files of other services
  for dep_service_name in "${MME_DEPS[@]}"
  do
    check_systemd_file "magma@mme" "$dep_service_name"
  done

  # check service config
  check_stateless_flag mme use_stateless; ret_check=$?
  if [[ $ret_check -eq $RETURN_STATEFUL ]]; then
    echo "MME config file is stateful."
    exit 1
  fi

  echo "MME service is stateless."
  exit 0
elif [[ $1 == "disable" ]]; then
  echo "Disabling stateless MME config"
  # force other services to restart when MME restarts
  for dep_service_name in "${MME_DEPS[@]}"
  do
    remove_systemd_override "$dep_service_name"
  done

  # change use_stateless setting in mme.yml
  disable_stateless_flag mme use_stateless
elif [[ $1 == "enable" ]]; then
  echo "Enabling stateless MME config"
  # stop other services from restarting when MME restarts
  for dep_service_name in "${MME_DEPS[@]}"
  do
    add_systemd_override "magma@mme" "$dep_service_name"
  done

  # change use_stateless setting in mme.yml
  enable_stateless_flag mme use_stateless true
else
  echo "Invalid argument. Use one of the following"
  echo "check: Run a check whether MME is stateless or not"
  echo "enable: Enable stateless mode, do nothing if already stateless"
  echo "disable: Disable stateless mode, do nothing if already stateful"
  exit 0
fi

# reload systemd config
sudo systemctl daemon-reload
