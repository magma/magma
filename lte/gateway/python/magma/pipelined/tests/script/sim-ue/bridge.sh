dir=$1
op=$2
tun_id_in=$3
ue_ip=$4
eth1_ip=$5
tun_id=$6

LOC="$dir/$tun_id"

tgt_ip="192.168.128.1"

function stop() {
  ovs-vsctl --db=unix:$LOC/db.sock del-br br0

  rm -f $LOC/conf.db
  rm -f $LOC/ovs-vswitchd.log

  if test -f $LOC/ovsdb.pid; then
	pkill -F $LOC/ovsdb.pid
  fi

  if test -f $LOC/vswitchd.pid; then
	  pkill -F $LOC/vswitchd.pid
  fi
}

function init() {
  mkdir -p $LOC
  stop

  sleep 1
  ovsdb-tool create $LOC/conf.db /usr/share/openvswitch/vswitch.ovsschema

  sleep 1

  ovsdb-server $LOC/conf.db \
    --remote=punix:$LOC/db.sock \
    --remote=db:Open_vSwitch,Open_vSwitch,manager_options \
    --private-key=db:Open_vSwitch,SSL,private_key \
    --certificate=db:Open_vSwitch,SSL,certificate \
    --bootstrap-ca-cert=db:Open_vSwitch,SSL,ca_cert --detach --log-file=$LOC/ovsdb.log  --pidfile=$LOC/ovsdb.pid

  ovs-vsctl --no-wait  --db=unix:$LOC/db.sock  init
  ovs-vswitchd  unix:$LOC/db.sock  --detach --log-file=$LOC/vswitchd.log --pidfile=$LOC/vswitchd.pid
  utilities/ovs-vsctl --db=unix:$LOC/db.sock show

  ovs-vsctl --db=unix:$LOC/db.sock add-br br0


  ifconfig br0 $ue_ip/24 up
  ovs-vsctl --db=unix:$LOC/db.sock add-port br0 gtp1 -- set Interface gtp1 type=gtp options:remote_ip=$eth1_ip option:key=flow
  ovs-vsctl --db=unix:$LOC/db.sock add-port br0 gtp0 -- set Interface gtp0 type=gtp options:remote_ip=flow option:key=flow

  ovs-ofctl del-flows br0
  br0_mac=$(ip l sh br0|grep link|awk '{print $2}')

  ovs-ofctl add-flow  br0  "in_port=LOCAL ip action=load:$tun_id->NXM_NX_TUN_ID[],output:gtp1"
  ovs-ofctl add-flow  br0  "in_port=gtp1 ip,tun_id=$tun_id_in action=mod_dl_dst:$br0_mac,output:local"
  ovs-ofctl add-flow  br0  "tun_id=$tun_id_in,ip action=mod_dl_dst:$br0_mac,output:local"

  ip neigh replace $tgt_ip  lladdr 00:12:34:56:78:aa dev br0
}


function sh() {
  ovs-vsctl --db=unix:$LOC/db.sock show 
  ovs-dpctl show
  ovs-dpctl dump-flows
}

$op
