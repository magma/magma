#!/bin/bash
source /usr/local/bin/config_stateless_utils.sh

MOBILITYD_DEPS=("magma@mme" "magma@pipelined")

if [[ $1 == "check" ]]; then
  # check dependency in systemd files of other services
  for dep_service_name in "${MOBILITYD_DEPS[@]}"
  do
    check_systemd_file "magma@mobilityd" "$dep_service_name"
  done

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
  for dep_service_name in "${MOBILITYD_DEPS[@]}"
  do
    remove_systemd_override "$dep_service_name"
  done

  # change persist_to_redis setting in mobilityd.yml
  sed -e '/persist_to_redis/ s/true/false/' -i /etc/magma/mobilityd.yml
elif [[ $1 == "enable" ]]; then
  echo "Enabling stateless mobilityd config"
  # remove restart dependencies between mobilityd and other services
  for dep_service_name in "${MOBILITYD_DEPS[@]}"
  do
    add_systemd_override "magma@mobilityd" "$dep_service_name"
  done

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
