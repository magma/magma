#!/usr/bin/env python3

"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import argparse
import subprocess

from magma.configuration.service_configs import load_service_config


def is_uplink(direction):
    return direction == 'OUT'


def is_downlink(direction):
    return direction == 'IN'


def _get_common_args(args):
    common_args = []
    if is_uplink(args.direction):
        common_args.append('in_port=%d' % args.in_port)
        common_args.append('tun_src=%s' % args.tun_src)
        common_args.append('tun_dst=%s' % args.tun_dst)
        if args.tun_id:
            common_args.append('tun_id=%s' % args.tun_id)
        common_args.append('dl_src=%s' % args.ue_mac)
        common_args.append('dl_dst=%s' % args.uplink_mac)
        common_args.append('ip_src=%s' % args.ue_ip_addr)
        common_args.append('ip_dst=%s' % args.uplink_ip_addr)
    else:
        common_args.append('in_port=%s' % args.in_port)
        common_args.append('dl_src=%s' % args.uplink_mac)
        common_args.append('dl_dst=%s' % args.ue_mac)
        common_args.append('ip_src=%s' % args.uplink_ip_addr)
        common_args.append('ip_dst=%s' % args.ue_ip_addr)
    return common_args


def _trace(bridge_name, pkt_args):
    cmd = ['find', '/var/run/openvswitch/', '-name', 'ovs-vswitchd*.ctl']
    output = subprocess.check_output(cmd)
    ovs_ctl_str = str(output, 'utf-8').strip()

    pkt_args_string = ','.join(pkt_args)
    cmd = [
        "ovs-appctl", '-t', ovs_ctl_str, "ofproto/trace", bridge_name,
        pkt_args_string,
    ]
    output = subprocess.check_output(cmd)
    output_str = str(output, 'utf-8').strip()
    print("Executing:")
    print(' '.join(cmd).strip())
    print("\nOutput:")
    print(output_str)


def _trace_tcp_pkt(args):
    pkt_args = ['tcp']
    pkt_args.extend(_get_common_args(args))

    if is_uplink(args.direction):
        pkt_args.append('tcp_src=%s' % args.ue_port)
        pkt_args.append('tcp_dst=%s' % args.uplink_port)
    else:
        pkt_args.append('tcp_src=%s' % args.uplink_port)
        pkt_args.append('tcp_dst=%s' % args.ue_port)
    _trace(args.bridge_name, pkt_args)


def _trace_icmp_pkt(args):
    pkt_args = ['icmp']
    pkt_args.extend(_get_common_args(args))
    _trace(args.bridge_name, pkt_args)


def create_parser():
    """
    Creates the argparse parser with all the arguments.
    """
    pipelined_config = load_service_config('pipelined')
    bridge_name = pipelined_config['bridge_name']
    ovs_gtp_port_number = pipelined_config['ovs_gtp_port_number']
    parser = argparse.ArgumentParser(
        description='CLI wrapper around ovs-appctl ofproto/trace.\n'
                    'Note: this won\'t be able to trace packets in userspace',
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )
    parser.add_argument(
        '-br', '--bridge_name', default=bridge_name,
        help='OVS bridge name(gtp_br0 or cwag_br0 or...)',
    )

    subparsers = parser.add_subparsers(title='direction', dest='cmd')
    uplink = subparsers.add_parser('uplink', help='Trace packet going from UE')
    uplink.add_argument(
        '-p', '--in_port', default=ovs_gtp_port_number,
        help='gtp port number',
    )
    dlink = subparsers.add_parser(
        'dlink',
        help='Trace packet coming from UPLINK',
    )
    dlink.add_argument(
        '-p', '--in_port', default='eth2',
        help='In port name',
    )

    uplink.add_argument('tun_src', help='Tunnel src ip', default='')
    uplink.add_argument(
        '-tun_dst', help='Tunnel dst ip',
        default='192.168.128.1',
    )
    uplink.add_argument('-tun_id', help='Tunnel id(key)', default='')
    uplink.set_defaults(direction='OUT')
    dlink.set_defaults(direction='IN')

    for dir_p in [uplink, dlink]:
        dir_p.add_argument(
            'ue_mac', help='UE mac address(unused in LTE)',
            default='0a:00:27:00:00:03',
        )
        dir_p.add_argument(
            '-ue_ip_addr', help='UE IP address',
            default='172.168.10.13',
        )
        dir_p.add_argument(
            '-uplink_mac', help='UPLINK mac address',
            default='00:27:12:00:33:00',
        )
        dir_p.add_argument(
            '--uplink_ip_addr', help='UPLINK IP address',
            default='104.28.26.94',
        )
        proto_sp = dir_p.add_subparsers(title='protos', dest='cmd')
        tcp = proto_sp.add_parser('tcp', help='Specify tcp packet')
        tcp.add_argument('-ue_port', type=int, help='', default=3372)
        tcp.add_argument('-uplink_port', type=int, help='', default=80)
        tcp.add_argument('-seq', type=int, default=10001)
        tcp.add_argument('-ack', type=int, default=10002)
        tcp.add_argument(
            '-tf', '--tcp_flags', help='Specify TCP flags',
            default='',
        )
        icmp = proto_sp.add_parser('icmp', help='Specify icmp packet')
        icmp.add_argument('-t', '--type', type=int, help='ICMP type', default=0)
        icmp.add_argument('-c', '--code', type=int, help='ICMP code', default=0)
        tcp.set_defaults(func=_trace_tcp_pkt)
        icmp.set_defaults(func=_trace_icmp_pkt)

    return parser


def main():
    parser = create_parser()
    # Parse the args
    args = parser.parse_args()
    if not args.cmd:
        parser.print_usage()
        exit(1)

    # Execute the subcommand function
    args.func(args)


if __name__ == "__main__":
    main()
