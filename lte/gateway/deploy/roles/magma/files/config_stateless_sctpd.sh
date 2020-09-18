#!/bin/bash
SRC_DIR=/usr/local/bin
PRE_START_CMD="ExecStartPre=$SRC_DIR/config_stateless_agw.sh\ sctpd_pre"
POST_START_CMD="ExecStartPost=$SRC_DIR/config_stateless_agw.sh\ sctpd_post"
SYSTEMD_FILE=/lib/systemd/system/sctpd.service
OVERRIDE_FILE=/etc/systemd/system/sctpd.service
RETURN_STATELESS=0
RETURN_STATEFUL=1
RETURN_CORRUPT=2

function check_stateless_sctpd {
  if [ -f "$OVERRIDE_FILE" ]; then
    sys_file=$OVERRIDE_FILE
  else
    sys_file=$SYSTEMD_FILE
  fi
  if ! grep -q "$PRE_START_CMD" $sys_file; then
   if ! grep -q "$POST_START_CMD" $sys_file; then
     echo "Sctpd is stateful"
     return $RETURN_STATEFUL
   else
     echo "Sctpd systemd file is corrupted"
     return $RETURN_CORRUPT
   fi
 elif ! grep -q "$POST_START_CMD" $sys_file; then
   echo "Sctpd systemd file is corrupted"
   return $RETURN_CORRUPT
 fi
 echo "Sctpd is stateless"
 return $RETURN_STATELESS
}

if [[ $1 == "check" ]]; then
  # check if the pre start and post start commands are in systemd file
  check_stateless_sctpd; ret_check=$?
  exit $ret_check
elif [[ $1 == "enable" ]]; then
  check_stateless_sctpd; ret_check=$?
  if [[ $ret_check -eq $RETURN_STATELESS ]]; then
    exit $RETURN_STATELESS
  fi
  echo "Enabling stateless Sctpd"
  # create override systemd service file, if it does not exist
  sudo /bin/cp -n $SYSTEMD_FILE $OVERRIDE_FILE
  # add a rule to clear Redis state whenever sctpd restarts
  sudo sed -i '/^ExecStart=.*/i '"$PRE_START_CMD" $OVERRIDE_FILE
  sudo sed -i '/^ExecStart=.*/a '"$POST_START_CMD" $OVERRIDE_FILE
elif [[ $1 == "disable" ]]; then
  check_stateless_sctpd; ret_check=$?
  if [[ $ret_check -eq $RETURN_STATEFUL ]]; then
    exit $RETURN_STATEFUL
  fi
  echo "Disabling stateless Sctpd"
  # remove the clear redis state command from sctpd system file
  sudo rm -f $OVERRIDE_FILE
else
  echo "Invalid argument. Use one of the following"
  echo "check: Run a check whether Sctpd is stateless or not"
  echo "enable: Enable stateless mode, do nothing if already stateless"
  echo "disable: Disable stateless mode, do nothing if already stateful"
  exit 0

fi

# reload systemd config
sudo systemctl daemon-reload
