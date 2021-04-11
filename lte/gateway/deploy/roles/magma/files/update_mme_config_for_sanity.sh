#!/bin/bash
RETURN_SUCCESS=0
RETURN_INVALID=1
RETURN_CONFIG_MISSING=2
RETURN_BACKUP_MISSING=3
mme_config_file=$MAGMA_ROOT/lte/gateway/configs/templates/mme.conf.template
mme_config_backup_file=$mme_config_file".bak"

function create_backup_or_restore_mme_config {
  # This function creates a backup of default MME configuration file,
  # which can later be used to restore the original configuration
  # In case the backup file is already present, it means there was failure
  # in last sanity run and current configuration file is already modified.
  # Hence, this function will restore the same backup file before modifying
  # it again, otherwise MME will crash in reading configuration from file
  if [[ -f $mme_config_backup_file ]]; then
    cp "$mme_config_backup_file" "$mme_config_file"
  else
    cp -n "$mme_config_file" "$mme_config_backup_file"
  fi
}

function configure_multiple_plmn_tac {
  # Remove default PLMN and TAC from MME configuration file
  sed -i -e '/GUMMEI_LIST/{n;d}' -e '/TAI_LIST/{n;N;N;N;N;N;N;d}' \
    -e '/TAC_LIST/{n;N;N;d}' "$mme_config_file"

  # Configure multiple PLMNs and TACs in MME configuration file
  gummei_config=(
    '{ MCC: "001"; MNC: "01"; MME_GID: "1"; MME_CODE: "1" }'
    '{ MCC: "001"; MNC: "02"; MME_GID: "1"; MME_CODE: "1" }'
    '{ MCC: "001"; MNC: "03"; MME_GID: "1"; MME_CODE: "1" }'
    '{ MCC: "001"; MNC: "04"; MME_GID: "1"; MME_CODE: "1" }'
    '{ MCC: "001"; MNC: "05"; MME_GID: "1"; MME_CODE: "1" }'
  )
  gummei_cmd_str=""
  for config in "${gummei_config[@]}"
  do
    gummei_cmd_str="$gummei_cmd_str\ \ \ \ \ \ \ \ $config,\n"
  done
  gummei_cmd_str=${gummei_cmd_str::-3}

  tac_config=(
    '{ MCC: "001"; MNC: "01"; TAC: "1" }'
    '{ MCC: "001"; MNC: "02"; TAC: "2" }'
    '{ MCC: "001"; MNC: "03"; TAC: "3" }'
    '{ MCC: "001"; MNC: "04"; TAC: "4" }'
    '{ MCC: "001"; MNC: "05"; TAC: "5" }'
  )
  tac_cmd_str=""
  for config in "${tac_config[@]}"
  do
    tac_cmd_str="$tac_cmd_str\ \ \ \ \ \ \ \ $config,\n"
  done
  tac_cmd_str=${tac_cmd_str::-3}

  sed -i -e "/GUMMEI_LIST/a $gummei_cmd_str" \
    -e "/TAI_LIST/a $tac_cmd_str" \
    -e "/TAC_LIST/a $tac_cmd_str" \
    "$mme_config_file"
}

function reduce_mobile_reachability_timer_value {
  # Reduce the mobile reachability timer to 1 minute, so that it can
  # quickly be tested as part of Sanity. The current default value of
  # Mobile Reachability Timer is 54 minutes
  sed -i '/^        T3412/s/54/1/' "$mme_config_file"
}

function configure_restricted_plmn {
  # Remove default restricted PLMN from MME configuration file
  sed -i -e '/RESTRICTED_PLMN_LIST/{n;N;N;N;N;N;N;d}' \
    "$mme_config_file"

  # Configure restricted PLMN/s in MME configuration file
  restricted_plmn_config=(
    '{ MCC= "123"; MNC= "450"}'
  )
  restricted_plmn_cmd_str=""
  for config in "${restricted_plmn_config[@]}"
  do
    restricted_plmn_cmd_str="$restricted_plmn_cmd_str\ \ \ \ \ \ \ \ $config,\n"
  done
  restricted_plmn_cmd_str=${restricted_plmn_cmd_str::-3}

  sed -i -e "/RESTRICTED_PLMN_LIST/a $restricted_plmn_cmd_str" \
    "$mme_config_file"
}

function configure_blocked_imei {
  # Remove default blocked imei(s) from MME configuration file
  sed -i -e '/BLOCKED_IMEI_LIST/{n;N;N;N;N;N;N;N;N;N;N;d}' \
    "$mme_config_file"

  # Configure blocked imei(s) in MME configuration file
  blocked_imei_config=(
    '{ IMEI_TAC="99000482"; SNR="351037" }'
    '{ IMEI_TAC="99333821"; }'
  )
  blocked_imei_cmd_str=""
  for config in "${blocked_imei_config[@]}"
  do
    blocked_imei_cmd_str="$blocked_imei_cmd_str\ \ \ \ \ \ \ \ $config,\n"
  done
  blocked_imei_cmd_str=${blocked_imei_cmd_str::-3}

  sed -i -e "/BLOCKED_IMEI_LIST/a $blocked_imei_cmd_str" \
    "$mme_config_file"
}

function restore_mme_config {
  # Restore the MME configuration from the backup configuration file and
  # delete the backup configuration file, so that MME will use latest
  # configuration file in next sanity runs
  if [[ -f $mme_config_backup_file ]]; then
    cp "$mme_config_backup_file" "$mme_config_file"
    rm -f "$mme_config_backup_file"
  else
    exit $RETURN_BACKUP_MISSING
  fi
}

if [[ $1 == "modify" ]]; then
  # Modify the MME configuration file so that all the sanity test cases pass
  if [[ ! -f $mme_config_file ]]; then
    exit $RETURN_CONFIG_MISSING
  fi
  create_backup_or_restore_mme_config
  configure_multiple_plmn_tac
  reduce_mobile_reachability_timer_value
  configure_restricted_plmn
  configure_blocked_imei
elif [[ $1 == "restore" ]]; then
  # Restore the MME configuration file from the backup config file
  restore_mme_config
else
  exit $RETURN_INVALID
fi

exit $RETURN_SUCCESS
