#!/bin/bash

# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.
#

# prints specified texts from syslog from AWS machines using fab get_machines and grep
# if <message to include> is empty, include all; if <message to exclude is empty>, do not exclude any
# usage: ./scripts/grep_syslog.sh tier service '<message to include>' '<message to exclude>'
# example: ./scripts/grep_syslog.sh prod controller 'Starting User Manager for UID' '1000'
# Has to be run under magma/platform/cloud/, where the get_machines fabfile.py is.

tier=$1
service=$2
inc=$3
exc=$4


if [ "$tier" = "" ]; then
    echo "usage: ./grep_syslog.sh tier service '<message to include>' '<message to exclude>'"
    exit 1
fi
if [ "$service" != "" ]; then
    get_machines_output=$(fab "$tier" get_machines:"$service")
else
    get_machines_output=$(fab "$tier" get_machines)
fi
machines=()
while read -r line; do
    machines+=("$line")
done <<< "$get_machines_output"
for host in "${machines[@]}"
do
    if [[ $host != ec2* ]]; then
        continue
    fi
    echo "$host"
    if [ "$exc" != "" ]; then
        ssh ubuntu@"$host" "grep -E '$inc' /var/log/syslog | grep -vE '$exc'"
    else  # if <message to exclude> is empty, then do not exclude anything
        ssh ubuntu@"$host" "grep -E '$inc' /var/log/syslog"
    fi
done

