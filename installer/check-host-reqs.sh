#!/bin/bash
#
# Provided by IRSOLS Inc
# Check Host requirements for installing Magma
# version 0.5
# last modified 08/15/2019

export RECOMMENDED_MEM="12"
export RECOMMENDED_CPU="4"
export CURR_MEM=`lsmem | grep online | grep Total | awk '{print $4}' | sed -e s/G//g`
export CURR_CPU=`lscpu | grep "CPU(s):" | head -1 | awk '{print $2}'`
export VIRTS=`lscpu | grep -i VT | awk '{print $2}'`

# echo " Checking Memory
if [ "$CURR_MEM" -lt "$RECOMMENDED_MEM" ]; then
echo "You do not have enough mem, exiting .."
exit
else 
echo "You have more than minimum recommended mem, proceeding .."
fi
#echo "Checking CPU"
if [ "$CURR_CPU" -lt "$RECOMMENDED_CPU" ]; then
echo "You do not have enough CPUs, exiting .."
exit
else 
echo "You have more than minimum recommended CPUs, proceeding .."
fi
#echo "Checking Virtualization Capabilities"
if [ -z "$VIRTS" ]
 then
 echo "You dont have nested virtualization enabled, please enable, exiting .."
 exit
else
 echo "You have Virtualization enabled , proceeding"
fi

