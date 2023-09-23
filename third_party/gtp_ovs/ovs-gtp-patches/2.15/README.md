These are patches for OVS code base to add support for GTP tunnel.
These patches will be removed once GTP patches are upstreamed to OVS repo.

There are couple of difference from existing GTP module:

1. GTP kernel module is included as part of OVS module. So no need
   to insert gtp.ko
2. GTP tunnel type is changed to 'gtpu'

These patches works on ubuntu 20.04 kernels.

1. Use build.sh to build OVS debian packages.

2. To setup DEV environment run following command on magma dev VM
`sudo bash ~/magma/third_party/gtp_ovs/ovs-gtp-patches/2.15/dev.sh setup`

3. To Run OVS GTP tests on OVS kernel datapath:
`sudo bash ~/magma/third_party/gtp_ovs/ovs-gtp-patches/2.15/dev.sh build_test`

4. For OVS development, use following ovs sources

```
sudo su -
cd ovs-build/ovs/
bash /home/vagrant/magma/third_party/gtp_ovs/ovs-gtp-patches/2.15/dev.sh build_test
```
