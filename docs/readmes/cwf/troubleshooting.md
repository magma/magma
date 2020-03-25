## Pipelined Packet Tracer Debugging


This is what the trace should look like when everything is setup properly. There might be additional tables based on a specific setup but the primary objective is for the packer to hit table 20 and to go out of either the uplink bridge(uplink_br0) or the gre port(32768, trace will say `output to kernel tunnel`)

```
root@cwf02-cff6689f7-pglcw:/# /usr/local/bin/packet_tracer_cli.py uplink  10.22.3.50 1.1.1.1 0 8c:f5:a3:bd:ad:5a 1.1.1.1 8c:f5:a3:bd:ad:ff tcp
Executing:
ovs-appctl -t /var/run/openvswitch/ovs-vswitchd.12394.ctl ofproto/trace cwag_br0 tcp,in_port=32768,tun_src=10.22.3.50,tun_dst=1.1.1.1,tun_id=0,dl_src=8c:f5:a3:bd:ad:5a,dl_dst=8c:f5:a3:bd:ad:ff,ip_src=1.1.1.1,ip_dst=104.28.26.94,tcp_src=3372,tcp_dst=80
Output:
Flow: tcp,tun_src=10.22.3.50,tun_dst=1.1.1.1,tun_ipv6_src=::,tun_ipv6_dst=::,tun_gbp_id=0,tun_gbp_flags=0,tun_tos=0,tun_ttl=0,tun_erspan_ver=0,tun_flags=0,in_port=32768,vlan_tci=0x0000,dl_src=8c:f5:a3:bd:ad:5a,dl_dst=8c:f5:a3:bd:ad:ff,nw_src=1.1.1.1,nw_dst=104.28.26.94,nw_tos=0,nw_ecn=0,nw_ttl=0,tp_src=3372,tp_dst=80,tcp_flags=0
bridge("cwag_br0")
------------------
 0. dl_src=8c:f5:a3:bd:ad:5a, priority 12
    set_field:0x75945885e55->metadata
    resubmit(,1)
     1. in_port=32768, priority 10
            set_field:0x1->reg1
            resubmit(,2)
             2. priority 0
                    resubmit(,3)
                     3. reg1=0x1, priority 0
                            resubmit(,25)
                            25. reg1=0x1,tun_id=0,tun_src=10.22.3.50, priority 10
                                    resubmit(,4)
                                     4. reg1=0x1, priority 0
                                            learn(table=21,priority=10,NXM_OF_ETH_DST[]=NXM_OF_ETH_SRC[],reg1=0x10,load:NXM_NX_TUN_ID[0..31]->NXM_NX_TUN_ID[0..31],load:NXM_NX_TUN_IPV4_SRC[]->NXM_NX_TUN_IPV4_DST[],load:NXM_NX_TUN_IPV4_DST[]->NXM_NX_TUN_IPV4_SRC[])
                                             >> suppressing side effects, so learn action ignored
                                            resubmit(,10)
                                            10. priority 12
                                                    sample(probability=30000,collector_set_id=3,obs_domain_id=1,obs_point_id=1,apn_mac_addr=0a:00:27:00:00:05,msisdn=magmaIsTheBest,apn_name=big_tower123,sampling_port=32768)
                                                    resubmit(,11)
                                                11. ip,reg1=0x1, priority 0
                                                    resubmit(,12)
                                                    12. ip,reg1=0x1,metadata=0x75945885e55, priority 65395, cookie 0x1
                                                            note:34.4d.42.5f.31.44.61.79.00.00.00.00.00.00
                                                            set_field:0x1->reg2
                                                            set_field:0x3->reg4
                                                            resubmit(,24)
                                                        24. reg1=0x1,reg2=0x1,reg4=0x3,metadata=0x75945885e55, priority 10, cookie 0x1
                                                            resubmit(,20)
                                                            20. reg1=0x1, priority 0
                                                                    output:1
                                                                bridge("uplink_br0")
                                                                --------------------
                                                                 0. priority 0
                                                                    NORMAL
                                                                     -> no learned MAC for destination, flooding
                                                            set_field:0->reg0
                                                    set_field:0->reg0
                                            set_field:0->reg0
                                    set_field:0->reg0
                            set_field:0->reg0
                    set_field:0->reg0
            set_field:0->reg0
    set_field:0->reg0
Final flow: tcp,reg1=0x1,reg2=0x1,reg4=0x3,tun_src=10.22.3.50,tun_dst=1.1.1.1,tun_ipv6_src=::,tun_ipv6_dst=::,tun_gbp_id=0,tun_gbp_flags=0,tun_tos=0,tun_ttl=0,tun_erspan_ver=0,tun_flags=0,metadata=0x75945885e55,in_port=32768,vlan_tci=0x0000,dl_src=8c:f5:a3:bd:ad:5a,dl_dst=8c:f5:a3:bd:ad:ff,nw_src=1.1.1.1,nw_dst=104.28.26.94,nw_tos=0,nw_ecn=0,nw_ttl=0,tp_src=3372,tp_dst=
80,tcp_flags=0
Megaflow: recirc_id=0,eth,tcp,tun_id=0,tun_src=10.22.3.50,tun_dst=1.1.1.1,tun_tos=0,tun_flags=-df-csum-key,in_port=32768,vlan_tci=0x0000/0x1fff,dl_src=8c:f5:a3:bd:ad:5a,dl_dst=8c:f5:a3:bd:ad:ff,nw_src=0.0.0.0/2,nw_dst=104.24.0.0/13,nw_frag=no,tp_dst=0x40/0xffc0
Datapath actions: sample(sample=45.8%,actions(userspace(pid=4242786699,flow_sample(probability=30000,collector_set_id=3,obs_domain_id=1,obs_point_id=1,output_port=4294967295)))),5,3
```

### To understand table assignment use the `pipelined_cli.py debug table_assignment` script
```
root@cwf02-cff6689f7-pglcw:/# /usr/local/bin/pipelined_cli.py debug table_assignment
App                      Main Table          Scratch Tables           
----------------------------------------------------------------------
ue_mac                   0                   [28]
ingress                  1                   []
arpd                     2                   []
access_control           3                   [21]
tunnel_learn             4                   [29]
vlan_learn               5                   [24, 25, 26, 27]
middle                   10                  []
check_quota              11                  []
enforcement              12                  [22]
enforcement_stats        12                  [23]
egress                   20                  []
```

The packet tracer usage is straight forward, you just need to fill out all the protocol fields. There are only a few points where the packets are usually dropped. The most frequent being the Subscriber is out of quota.
 - Using the packet tracer cli you would notice the packet being dropped at the enforcement controller
 - There is an issue with a gre tunnel, the packet tracer would show that the packets get dropped because they don't meet the access control tunnel filtering
 - If subscriber is not at all attached the traffic will be dropped at table 0 (no flow set by the UE controller)
 - For all other issues, note the table the packet is being dropped and check the flow dump of that table, get logs for the controller responsible for that table.
