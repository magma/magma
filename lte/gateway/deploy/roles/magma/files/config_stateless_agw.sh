#!/bin/bash
SRC_DIR=$MAGMA_ROOT/lte/gateway/deploy/roles/magma/files
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
    sudo "$SRC_DIR/config_stateless_$service_name.sh" disable
  done
elif [[ $1 == "enable" ]]; then
  if check_stateless_agw; then      # Checks whether return was success, i.e. 0
    echo "Nothing to enable, AGW is stateless"
    exit $RETURN_STATELESS
  fi
  echo "Enabling stateless AGW config"
  for service_name in "${SERVICE_LIST[@]}"
  do
    sudo "$SRC_DIR/config_stateless_$service_name.sh" enable
  done
else
  echo "Invalid argument. Use one of the following"
  echo "check: Run a check whether AGW is stateless or not"
  echo "enable: Enable stateless mode, do nothing if already stateless"
  echo "disable: Disable stateless mode, do nothing if already stateful"
  exit $RETURN_INVALID
fi

sudo service magma@* stop

# Remove all stored state from redis
sudo service magma@redis start
redis-cli -p 6380 FLUSHALL
sudo service magma@redis stop

# force restart sctpd so that eNB connections are reset
sudo service sctpd restart

sudo service magma@magmad start
# Sleep for a bit so OVS and Magma services come up before proceeding
sleep 15
echo "Config complete"

check_stateless_agw; ret_check=$?
exit $ret_check
