#!/bin/bash

padlimit=60

printf -v line '%*s' "$padlimit"
printf -v pad '%*s' "$padlimit"
line=+${line// /-}+

function print_data {
  printf "| %s %0*s |\n" "$1" $((padlimit - ${#1} - 3 )) " "
}

function print_data_padded {
  printf "| %0*s %0*s |\n" $2 "$1" $((padlimit - ${#1} - 3 - $(( $2 - ${#1} )) )) " "
}

function print_center {
  if [ $((${#1}/2*2)) -eq ${#1} ]
  then
    extra=0
  else
    extra=1
  fi

  count=$(((padlimit -${#1} -3)/2))
  printf "| %0*s %s %0*s |\n" $(($count-$extra)) " " "$1" $count " "
}

echo ""
echo $line
print_center "Magma Access Gateway"
echo $line
print_data "Gateway ID:       $(cat /etc/snowflake)"
echo $line
for interface in $(ip a | grep -oP "eth\d|uplink_br\d" | uniq | sort);
do
  ip=$(ip -4 addr | grep $interface | grep -oP '(?<=inet\s)\d+(\.\d+){3}(\/\d+)?')
  print_data "$interface: $ip"
done
echo $line
print_data "Jumphost ports: $(systemctl list-units | grep -oP '(?<=revssh@)(\d*)(?=\.service)' | tr '\n' ' ')"
echo $line
print_data "Free RAM:     $(free -h | grep Mem | awk '{print $4}')"
print_data "Free SWAP:    $(free -h | grep Swap | awk '{print $4}')"
print_data "Free storage: $(df -h | grep -oP "\d*[A-Z](?=\s*\d*%\s*\/$)" | sed 's/  / /g')"
echo $line
echo ""
echo $line
print_center "Magma status"
print_center "$(apt-cache policy magma | grep -oP '(?<=Installed: )[^\s]+')"
echo $line
print_data "Subscribers connected: $(mobility_cli.py get_subscriber_table | tail -n +2 | wc -l)"
echo $line
for service in $(systemctl list-units --all | grep -oP "(?<=magma@)([a-z_]*)(?=\.service)")
do
  print_data_padded "$service: $(systemctl is-active magma@$service)" 25
done
echo $line
echo ""
echo $line
print_center  "SCTP port Status"
echo $line
ip=$(ip -4 add | grep eth1 | grep -oP '(?<=inet\s)\d+(\.\d+){3}(\/\d+)?' | cut -d/ -f1)
if ( netstat -l | grep sctp | grep $ip:36412 > /dev/null ); 
then 
  print_data "sctp    $ip:36412   LISTEN"
else 
  print_data "sctp    $ip:36412   DOWN"
fi
if ( netstat -l | grep sctp | grep $ip:38412 > /dev/null ); 
then 
  print_data "sctp    $ip:38412   LISTEN"
else 
  print_data "sctp    $ip:38412   DOWN"
fi
echo $line
echo ""