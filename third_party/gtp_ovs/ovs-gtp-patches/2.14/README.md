These are patches for OVS code base to add support for GTP tunnel.
These patches will be removed once GTP patches are upstreamed to OVS repo.

There are couple of difference from existing GTP module:
1. GTP kernel module is included as part of OVS module. So no need
   to insert gtp.ko
2. GTP tunnel type is changed to 'gtpu'

These patches should work on following kernel version:
1. 4.9.214
2. 4.14.171
3. 4.19.110
4. 5.4.50
5. 5.6.19

Steps to build package from source code.
1. checkout magma repo, set MAGMA_ROOT to the repo root dir.
2. cd $MAGMA_ROOT/third_party/gtp_ovs/ovs/2.14/
3. git clone https://github.com/openvswitch/ovs
4. cd ovs/
5. git checkout branch-2.14 /* Checkout ovs2.14.1 */
6. git am ../00*
7. DEB_BUILD_OPTIONS='parallel=8 nocheck' fakeroot debian/rules binary
8. Packages are copied in parent (..) dir
