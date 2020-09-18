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
