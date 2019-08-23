import subprocess

import fire
from lte.protos.pipelined_pb2 import SerializedRyuPacket
from lte.protos.pipelined_pb2_grpc import PipelinedStub
from magma.common.service_registry import ServiceRegistry
from ryu.lib.packet import ethernet, arp, ipv4, icmp
from ryu.lib.packet.ether_types import ETH_TYPE_ARP, ETH_TYPE_IP
from ryu.lib.packet.packet import Packet


def exec_commandline(command, exception=None):
    try:
        p = subprocess.Popen(command.split(),
                             stdin=subprocess.PIPE,
                             stdout=subprocess.PIPE,
                             stderr=subprocess.PIPE)
        out, err = p.communicate()

        if p.returncode and exception:
            raise exception
        return p.returncode, out, err
    except OSError:
        raise exception


class PacketTracerCLI:
    BRIDGE_NAME = "gtp_br0"

    def dump_ports(self, port=None):
        command = 'ovs-ofctl dump-ports {} {}'.format(self.BRIDGE_NAME, port) \
            if port else 'ovs-ofctl dump-ports {}'.format(self.BRIDGE_NAME)
        print(exec_commandline(command)[1].decode('utf-8'))

    def dump_flows(self):
        flows = exec_commandline('ovs-ofctl dump-flows {}'.format(
            self.BRIDGE_NAME
        ))[1].decode('utf-8')
        flows = flows.split('\n')
        flows = [' '.join([i for i in flow.split()
                           if 'duration' not in i
                           and 'idle' not in i])
                 for flow in flows]
        print('\n'.join(flows))

    def raw(self, data):
        data = bytes(data)
        pkt = Packet(data)

        # Send the packet to a grpc service
        chan = ServiceRegistry.get_rpc_channel('pipelined',
                                               ServiceRegistry.LOCAL)
        client = PipelinedStub(chan)

        print('Sending: {}'.format(pkt))
        table_id = client.TracePacket(SerializedRyuPacket(pkt=data))

        if table_id == -1:
            print('Successfully passed through all the tables!')
        else:
            print('Dropped by table: {}'.format(table_id))

    def icmp(self):
        pkt = ethernet.ethernet(dst='5e:cc:cc:b1:49:4b') / \
              ipv4.ipv4(src='192.168.70.2',
                        dst='192.168.70.3',
                        proto=1) / \
              icmp.icmp()
        pkt.serialize()
        self.raw(data=pkt.data)

    def arp(self):
        pkt = ethernet.ethernet(ethertype=ETH_TYPE_ARP,
                                src='fe:ee:ee:ee:ee:ef',
                                dst='ff:ff:ff:ff:ff:ff') / \
              arp.arp(hwtype=arp.ARP_HW_TYPE_ETHERNET, proto=ETH_TYPE_IP,
                      hlen=6, plen=4,
                      opcode=arp.ARP_REQUEST,
                      src_mac='fe:ee:ee:ee:ee:ef', src_ip='192.168.70.2',
                      dst_mac='00:00:00:00:00:00', dst_ip='192.168.70.3')
        pkt.serialize()
        self.raw(data=pkt.data)


if __name__ == '__main__':
    cli = PacketTracerCLI()
    fire.Fire(cli)
