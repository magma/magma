iptables -t mangle -I PREROUTING -p tcp --dport 80 -j MARK --set-mark 1
iptables -t mangle -I PREROUTING -p tcp --sport 80 -j MARK --set-mark 1

ip rule add fwmark 1 lookup 100
ip route add local 0.0.0.0/0 dev lo table 100

sysctl -w net.ipv4.conf.all.rp_filter=0
sysctl -w net.ipv4.conf.all.route_localnet=1
