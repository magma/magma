#!/bin/bash

echo "get controller ip"
[[ -z "${CTRL_IP}" ]] && CtrlIP="$(getent hosts ofproxy | awk '{ print $1 }')" || CtrlIP="${CTRL_IP}"

echo "Running in $CONNECTION_MODE"
CONNECTION_MODE=${CONNECTION_MODE:=tcp}
if [ "$CONNECTION_MODE" == "ssl" ]
then
  CtrlIP="$(getent hosts tls-termination | awk '{ print $1 }')" || CtrlIP="${CTRL_IP}"
fi

echo "start ovs-ctl"
/usr/share/openvswitch/scripts/ovs-ctl start --system-id=random --no-ovs-vswitchd
/usr/share/openvswitch/scripts/ovs-ctl stop
echo "start db server"
ovsdb-server --pidfile /etc/openvswitch/conf.db -vconsole:emer -vsyslog:err -vfile:info \
--remote=punix:/var/run/openvswitch/db.sock --private-key=db:Open_vSwitch,SSL,private_key \
--certificate=db:Open_vSwitch,SSL,certificate --bootstrap-ca-cert=db:Open_vSwitch,SSL,ca_cert --log-file=/var/log/openvswitch/ovsdb-server.log --no-chdir &
    ovs-vswitchd --pidfile -vconsole:emer -vsyslog:err -vfile:info --mlockall --no-chdir --log-file=/var/log/openvswitch/ovs-vswitchd.log &

# Copy files to /etc/magma it must be here and not in dockerfile because the volume
# are shared and may be taint on the local host
echo "copy config file"
cp cwf/gateway/configs/* /etc/magma/
cp xwf/gateway/configs/* /etc/magma/
cp orc8r/gateway/configs/templates/* /etc/magma/

echo "get xwfwhoami"
curl -X POST  https://graph.expresswifi.com/openflow/configxwfm?access_token=$ACCESSTOKEN | jq -r .configxwfm > /etc/xwfwhoami
sed -i '/^uplink_if/d'  /etc/xwfwhoami # TODO: remove this

echo "run XWF ansible"
ANSIBLE_CONFIG=xwf/gateway/ansible.cfg ansible-playbook -e xwf_ctrl_ip="${CtrlIP} connection_mode=$CONNECTION_MODE" \
xwf/gateway/deploy/xwf.yml -i "localhost," --skip-tags "install,install_docker,no_ci" -c local -v

echo "run DNS server"
dnsmasq

echo "run DHCP server"
/usr/sbin/dhcpd -f -cf /etc/dhcp/dhcpd.conf -user dhcpd -group dhcpd --no-pid gw0 &

echo "loop forever"
tail -f /dev/null
