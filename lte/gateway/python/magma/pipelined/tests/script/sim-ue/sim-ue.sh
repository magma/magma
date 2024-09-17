op=$1
tun_id=$2
ue_ip=$3
tun_id_in=$4

eth1_ip=192.168.60.142
enb_ip=192.168.60.141

LOC="/var/ovs/"

function s() {
  	tun_id=$1
  	ue_ns="ue_ns_$tun_id"
  	ue_dev="ue_dev_$tun_id"

	ip netns add $ue_ns

	ip link add $ue_dev  type veth peer name  "$ue_dev"_ns

	ip link set dev "$ue_dev"_ns  netns $ue_ns

	ifconfig "$ue_dev" up
	ip route add $enb_ip dev $ue_dev

	ip netns exec  $ue_ns     ifconfig lo up
	ip netns exec  $ue_ns     ifconfig "$ue_dev"_ns $enb_ip/24 up
	ip netns exec  $ue_ns     ip route add default via $eth1_ip 

	ip netns exec  $ue_ns  bash -x  bridge.sh $LOC init $tun_id_in $ue_ip $eth1_ip $tun_id
}

function d() {
    	tun_id=$1
    	ue_ns="ue_ns_$tun_id"
    	ue_dev="ue_dev_$tun_id"

	ip netns exec  $ue_ns bash -x  bridge.sh $LOC stop "d1" "d2" "d3" $tun_id
	ip a flush $ue_dev
  	ip link del $ue_dev 
  	sleep 1
  	ip netns delete $ue_ns
}

# destroy all
function da() {
	for file in $(ls $LOC)
  	do
    		tun_id=$file
    		d $tun_id
  	done
}

function sh() {
    	tun_id=$1
    	ue_ns="ue_ns_$tun_id"
 
	ip netns exec $ue_ns bash bridge.sh $LOC sh "d1" "d2" "d3" $tun_id
}

$op $tun_id

