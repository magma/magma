#!/bin/bash
SRC_DIR=/usr/local/bin
SERVICE_LIST=("mme" "mobilityd" "pipelined" "sessiond")
RETURN_STATELESS=0
RETURN_STATEFUL=1
RETURN_CORRUPT=2
RETURN_INVALID=3

function check_stateless_agw {
  echo "Checking stateless AGW config"
  num_stateful=0
  for service_name in "${SERVICE_LIST[@]}"
  do
    if ! "$SRC_DIR/config_stateless_$service_name.sh" check; then
      num_stateful=$((num_stateful + 1))
    fi
  done
  if [[ $num_stateful -eq 0 ]]; then
    return $RETURN_STATELESS
  elif [[ $num_stateful -eq ${#SERVICE_LIST[@]} ]]; then
    return $RETURN_STATEFUL
  fi
  return $RETURN_CORRUPT
}

function stop_and_clear_state {
  # Remove all stored state from Redis
  sudo service magma@* stop
  sudo service magma@redis start
  redis-cli -p 6380 FLUSHALL
  sudo service magma@redis stop
}

if [[ $1 == "check" ]]; then
  check_stateless_agw; ret_check=$?
  if [[ $ret_check -eq 1 ]]; then
    echo "AGW is stateful."
    exit $RETURN_STATEFUL
  elif [[ $ret_check -eq 2 ]]; then
    echo "Some AGW services are stateless, while others are not."
    echo "Please run with enable/disable command to fix it."
    exit $RETURN_CORRUPT
  fi

  echo "AGW is stateless."
  exit $RETURN_STATELESS
elif [[ $1 == "disable" ]]; then
  check_stateless_agw; ret_check=$?
  if [[ $ret_check -eq 1 ]]; then
    echo "Nothing to disable, AGW is stateful"
  exit $RETURN_STATEFUL
  fi

  echo "Disabling stateless AGW config"
  for service_name in "${SERVICE_LIST[@]}"
  do
    sudo -E "$SRC_DIR/config_stateless_$service_name.sh" disable
  done
elif [[ $1 == "enable" ]]; then
  if check_stateless_agw; then      # Checks whether return was success, i.e. 0
    echo "Nothing to enable, AGW is stateless"
    exit $RETURN_STATELESS
  fi
  echo "Enabling stateless AGW config"
  for service_name in "${SERVICE_LIST[@]}"
  do
    sudo -E "$SRC_DIR/config_stateless_$service_name.sh" enable
  done
elif [[ $1 == "sctpd_pre" ]]; then
  # In stateless mode, clear Redis state before sctpd starts
  check_stateless_agw; ret_check=$?
  if [[ $ret_check -eq 1 ]]; then
    echo "AGW is stateful."
    exit 0
  fi
  stop_and_clear_state
  exit 0
elif [[ $1 == "sctpd_post" ]]; then
  # In stateless mode, start magmad after sctpd starts
  check_stateless_agw; ret_check=$?
  if [[ $ret_check -eq 1 ]]; then
    echo "AGW is stateful."
    exit 0
  fi
  sudo service magma@magmad start
  exit 0
else
  echo "Invalid argument. Use one of the following"
  echo "check: Run a check whether AGW is stateless or not"
  echo "enable: Enable stateless mode, do nothing if already stateless"
  echo "disable: Disable stateless mode, do nothing if already stateful"
  exit $RETURN_INVALID
fi

# Force restart sctpd so that eNB connections are reset and
# local state is cleared before sctpd starts
sudo service sctpd restart

echo "Config complete"

check_stateless_agw; ret_check=$?

if [[ $ret_check -eq $RETURN_STATEFUL ]]; then
  # For stateless AGW, Magmad is started as part of sctpd_post
  sudo service magma@magmad start
fi

# Sleep for a bit so OVS and Magma services come up before proceeding
sleep 60

exit $ret_check
