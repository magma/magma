# Remove MME dependencies on Pipelined and Sessiond
sed '/magma@pipelined/d' -i /etc/systemd/system/magma@mme.service
sed '/magma@sessiond/d' -i /etc/systemd/system/magma@mme.service

# Remove the openvswitch gtp bridge
ifdown gtp_br0

# Remove the oai-gtp package which installs the custom gtp module
apt-get -y remove oai-gtp

# Remove the custom gtp module needed by openvswitch
rmmod vport_gtp
rmmod openvswitch
rmmod gtp

# Install the default kernel gtp module
insmod /lib/modules/`uname -r`/kernel/drivers/net/gtp.ko

# Update the module symbols
depmod -a

# Enable NATing on the SGI interface, i.e. eth2
iptables -t nat -A POSTROUTING -o eth2 -j MASQUERADE

# Install libgtpnl
bash /home/vagrant/magma/third_party/libgtpnl/install.sh
