#!/bin/bash
rm -rf test
fab upgrade_teravm_agw:setup_2 >> test
cat test| grep "Upgrade Status" |cut -d':' -f 1,2,3
rm -rf test
fab upgrade_teravm_feg:setup_2 >> test
cat test| grep "Upgrade Status" |cut -d':' -f 1,2,3
rm -rf test
#fab run_3gpp_tests:setup_2
