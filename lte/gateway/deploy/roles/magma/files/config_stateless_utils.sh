#!/bin/bash
RETURN_STATELESS=0
RETURN_STATEFUL=1

function check_systemd_file {
  service_name=$1
  dep_service_name=$2
  systemd_file=/lib/systemd/system/"$dep_service_name".service
  override_file=/etc/systemd/system/"$dep_service_name".service
  if [ -f "$override_file" ]; then
    systemd_file=$override_file
  fi

  if grep -q "^ *PartOf=$service_name" "$systemd_file"; then
    echo "The $dep_service_name service will restart with $service_name," \
      "i.e. stateful."
    exit $RETURN_STATEFUL
  else
    return $RETURN_STATELESS
  fi
}

function add_systemd_override {
  service_name=$1
  dep_service_name=$2
  # create override systemd service file, if it does not exist
  sudo /bin/cp -n /lib/systemd/system/"$dep_service_name".service \
    /etc/systemd/system/"$dep_service_name".service
  sudo sed -e "/PartOf=$service_name/ s/^#*/#/" -i \
    /etc/systemd/system/"$dep_service_name".service
}

function remove_systemd_override {
  dep_service_name=$1
  sudo rm -f /etc/systemd/system/"$dep_service_name".service
}

function check_stateless_flag {
  service_name=$1
  flag=$2
  value=false  # stateful
  override_file=/var/opt/magma/configs/$service_name.yml
  if [[ $service_name == "pipelined" ]]; then
    value=true # pipelined has a reverse logic compared to other services
  fi

  # check if stateless_flag config exists in override file
  if grep -s -q "$flag" "$override_file"; then
    if grep -q "$flag: $value" "$override_file"; then
      return $RETURN_STATEFUL
    else
      return $RETURN_STATELESS
    fi
  fi

  # check in regular file
  if grep -q "$flag: $value" /etc/magma/"$service_name".yml; then
    return $RETURN_STATEFUL
  else
    return $RETURN_STATELESS
  fi
}

function enable_stateless_flag {
  service_name=$1
  flag=$2
  value=$3
  override_file=/var/opt/magma/configs/$service_name.yml

  # check if override file exists
  if [ ! -f "$override_file" ]; then
    sudo mkdir -p /var/opt/magma/configs && touch "$override_file"
  fi

  # check if the setting exists and set the correct value
  if ! grep -q "$flag" "$override_file"; then
    echo "$flag: $value" >> "$override_file"
  elif ! grep -q "$flag: $value" "$override_file"; then
    sed -e "/$flag/ s/.*/$flag:\ $value/" -i "$override_file"
  fi
}

function disable_stateless_flag {
  service_name=$1
  flag=$2
  override_file=/var/opt/magma/configs/$service_name.yml

  # check if override file exists
  if [ -f "$override_file" ]; then
    # remove the stateless config override as default is always stateful
    sed -i "/$flag/d" "$override_file"

    # delete override file if it is empty
    if [ ! -s "$override_file" ]; then
      sudo rm -f "$override_file"
    fi
  fi
}
