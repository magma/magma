#!/bin/bash
ansible-playbook /opt/xwfm/deploy/xwf.yml --tags=run_pipelined -e "upgrade=true"

