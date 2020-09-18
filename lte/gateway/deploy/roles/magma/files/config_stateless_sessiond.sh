#!/bin/bash
source /usr/local/bin/config_stateless_utils.sh

if [[ $1 == "check" ]]; then
  # check dependency in systemd files of other services
  check_systemd_file "magma@sessiond" "magma@mme"

  #check service config
  if ! grep -q "support_stateless.*true" /etc/magma/sessiond.yml; then
    echo "Sessiond config file is stateful."
    exit 1
  fi

  echo "Sessiond service is stateless."
  exit 0
elif [[ $1 == "disable" ]]; then
  echo "Disabling stateless sessiond config"
  # restore restart dependencies between sessiond and other services
  remove_systemd_override "magma@mme"

  # change support_stateless setting in sessiond.yml
  sed -e '/support_stateless/ s/true/false/' -i /etc/magma/sessiond.yml
elif [[ $1 == "enable" ]]; then
  echo "Enabling stateless sessiond config"
  # remove restart dependencies between sessiond and other services
  add_systemd_override "magma@sessiond" "magma@mme"

  # change support_stateless setting in sessiond.yml
  sed -e '/support_stateless/ s/false/true/' -i /etc/magma/sessiond.yml
else
  echo "Invalid argument. Use one of the following"
  echo "check: Run a check whether Sessiond is stateless or not"
  echo "enable: Enable stateless mode, do nothing if already stateless"
  echo "disable: Disable stateless mode, do nothing if already stateful"
  exit 0
fi

# reload systemd config
sudo systemctl daemon-reload
