

apt update

# rename interfaces
sed -i 's/enp0s3/eth0/g' /etc/netplan/50-cloud-init.yaml
sed -i 's/GRUB_CMDLINE_LINUX=""/GRUB_CMDLINE_LINUX="net.ifnames=0 biosdevname=0"/g' /etc/default/grub
grub-mkconfig -o /boot/grub/grub.cfg

# device management
apt install ifupdown

echo "auto eth0
iface eth0 inet dhcp" >> /etc/network/interfaces
# configuring eth1
echo "auto eth1
iface eth1 inet static
address 10.10.2.1
netmask 255.255.255.0" >> /etc/network/interfaces


# name server config
ln -sf /var/run/systemd/resolve/resolv.conf /etc/resolv.conf

# get rid of netplan
systemctl unmask networking
systemctl enable networking
sleep 5

apt-get --assume-yes purge nplan netplan.i

