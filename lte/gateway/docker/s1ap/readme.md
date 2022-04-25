# Build and publish s1aptester images

1. Move into AGW docker directory in the repo and run build script

```
cd lte/gateway/docker
s1ap/build-s1ap.sh
```

2. Publish images to your registry

```
s1ap/publish.sh yourregistry.com/yourrepo/
```

# Run s1aptester

1. Create a host that has 3 interfaces. eth0 which is your SGi interface, eth1 which is your S1 interface, and a third interface which is your ssh management interface to connect to the instance while doing the tests. You could skip the third interface if you have some kind of serial console access.

Due to the IP address contraints of how s1aptester is written, we have to hardcode the addresses of eth0 and eth1. Here is an example netplan configuration:

```
network:
    ethernets:
        # SGi
        eth0:
            dhcp6: false
            dhcp4: no
            addresses: [192.168.60.142/24]
            nameservers:
              addresses: [8.8.8.8,8.8.4.4]
        # S1
        eth1:
            dhcp6: false
            dhcp4: no
            addresses: [192.168.129.1/24]
            nameservers:
              addresses: [8.8.8.8,8.8.4.4]
        # SSH Management
        eth2:
            dhcp4: true
            dhcp6: false
    version: 2
```

Add a password to the ubuntu user and enable password authentication for the ssh server on the AGW host.

```
sed -i 's/PasswordAuthentication no/PasswordAuthentication yes/' /etc/ssh/sshd_config
systemctl restart ssh

passwd ubuntu
```

2. Move into AGW docker directory on the host and run start script. Make sure that your `.env` file points to your registry with your AGW and s1aptester images.
```
cd /var/opt/magma/docker
s1ap/start-s1ap.sh
```

2. This will drop you into a shell that you can start to run tests from, or run the full suite of tests.
```
root@472f8708ec12:/magma/lte/gateway/python/integ_tests#
# Run individual test(s)
make integ_test TESTS=s1aptests/test_attach_detach.py
# Run full suite
make integ_test
```

# Stop s1aptester

Move into AGW docker directory on the host and run stop script.
```
cd /var/opt/magma/docker
s1ap/stop-s1ap.sh
```

If inside of container, CTRL+d or exit from container and run stop script
```
root@472f8708ec12:/magma/lte/gateway/python/integ_tests# exit
s1ap/stop-s1ap.sh
```
